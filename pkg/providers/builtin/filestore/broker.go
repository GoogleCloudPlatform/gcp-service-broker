// Copyright 2019 the Service Broker Project Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package filestore

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/GoogleCloudPlatform/gcp-service-broker/db_service/models"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/broker"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/providers/builtin/base"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/varcontext"
	"github.com/GoogleCloudPlatform/gcp-service-broker/utils"
	"github.com/pivotal-cf/brokerapi"
	filestore "google.golang.org/api/file/v1"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

const (
	// DefaultFileshareName is the default created for fileshares.
	DefaultFileshareName = "filestore"
)

// NewInstanceInformation creates instance information from an instance
func NewInstanceInformation(instance filestore.Instance) (*InstanceInformation, error) {
	if len(instance.Networks) == 0 {
		return nil, errors.New("no networks were defined on the instance")
	}
	network := instance.Networks[0]

	if len(network.IpAddresses) == 0 {
		return nil, errors.New("no IP addresses were defined on the instance")
	}
	ip := network.IpAddresses[0]

	if len(instance.FileShares) == 0 {
		return nil, errors.New("no file shares were defined on the instance")
	}
	share := instance.FileShares[0]

	return &InstanceInformation{
		Network:         network.Network,
		ReservedIPRange: network.ReservedIpRange,
		IPAddress:       ip,
		FileShareName:   share.Name,
		CapacityGB:      share.CapacityGb,
		URI:             fmt.Sprintf("nfs://%s/%s", ip, share.Name),
	}, nil
}

// InstanceInformation holds the details needed to get a Filestore instance
// once it's been created.
type InstanceInformation struct {
	Network         string `json:"authorized_network"`
	ReservedIPRange string `json:"reserved_ip_range"`
	IPAddress       string `json:"ip_address"`
	FileShareName   string `json:"file_share_name"`
	CapacityGB      int64  `json:"capacity_gb"`
	URI             string `json:"uri"`
}

// Broker is the back-end for creating and binding to Google Cloud Filestores.
type Broker struct {
	base.PeeredNetworkServiceBase
	base.NoOpBindMixin
}

var _ (broker.ServiceProvider) = (*Broker)(nil)

// Provision implements ServiceProvider.Provision.
func (b *Broker) Provision(ctx context.Context, provisionContext *varcontext.VarContext) (models.ServiceInstanceDetails, error) {
	details := models.ServiceInstanceDetails{
		Name:     provisionContext.GetString(base.InstanceIDKey),
		Location: provisionContext.GetString(base.ZoneKey),
	}

	instance := &filestore.Instance{
		Labels: provisionContext.GetStringMapString("labels"),
		Tier:   provisionContext.GetString("tier"),
		FileShares: []*filestore.FileShareConfig{
			{
				Name:       DefaultFileshareName,
				CapacityGb: int64(provisionContext.GetInt("capacity_gb")),
			},
		},

		Networks: []*filestore.NetworkConfig{
			{
				Modes: []string{
					provisionContext.GetString("address_mode"),
				},
				Network: provisionContext.GetString(base.AuthorizedNetworkKey),
			},
		},
	}

	if err := provisionContext.Error(); err != nil {
		return models.ServiceInstanceDetails{}, err
	}

	client, err := b.createClient(ctx)
	if err != nil {
		return models.ServiceInstanceDetails{}, err
	}

	op, err := client.Projects.Locations.Instances.
		Create(b.parentPath(details), instance).
		InstanceId(details.Name).
		Do()

	if err != nil {
		return models.ServiceInstanceDetails{}, err
	}

	details.OperationType = models.ProvisionOperationType
	details.OperationId = op.Name

	return details, nil
}

// Deprovision implements ServiceProvider.Deprovision.
func (b *Broker) Deprovision(ctx context.Context, instance models.ServiceInstanceDetails, details brokerapi.DeprovisionDetails) (*string, error) {
	client, err := b.createClient(ctx)
	if err != nil {
		return nil, err
	}

	op, err := client.Projects.Locations.Instances.Delete(b.instancePath(instance)).Do()
	if err != nil {
		// Mark things that have been deleted out of band as gone.
		if gerr, ok := err.(*googleapi.Error); ok {
			if gerr.Code == http.StatusNotFound || gerr.Code == http.StatusGone {
				return nil, nil
			}
		}

		return nil, err
	}

	return &op.Name, nil
}

func (b *Broker) createClient(ctx context.Context) (*filestore.Service, error) {
	co := option.WithUserAgent(utils.CustomUserAgent)
	ct := option.WithTokenSource(b.HTTPConfig.TokenSource(ctx))
	c, err := filestore.NewService(ctx, co, ct)
	if err != nil {
		return nil, fmt.Errorf("couldn't instantiate API client: %s", err)
	}

	return c, nil
}

func (b *Broker) parentPath(instanceDetails models.ServiceInstanceDetails) string {
	return fmt.Sprintf("projects/%s/locations/%s", b.DefaultProjectID, instanceDetails.Location)
}

func (b *Broker) instancePath(instanceDetails models.ServiceInstanceDetails) string {
	return fmt.Sprintf("%s/instances/%s", b.parentPath(instanceDetails), instanceDetails.Name)
}

//
// func (b *Broker) operationPath(instanceDetails models.ServiceInstanceDetails) string {
// 	return fmt.Sprintf("%s/operations/%s", b.parentPath(instanceDetails), instanceDetails.OperationId)
// }

// PollInstance implements ServiceProvider.PollInstance
func (b *Broker) PollInstance(ctx context.Context, instance models.ServiceInstanceDetails) (done bool, err error) {
	if instance.OperationType == "" {
		return false, errors.New("couldn't find any pending operations")
	}

	client, err := b.createClient(ctx)
	if err != nil {
		return false, err
	}

	op, err := client.Projects.Locations.Operations.Get(instance.OperationId).Do()
	if op != nil {
		done = op.Done
	}

	return done, err
}

// UpdateInstanceDetails updates the ServiceInstanceDetails with the most recent state from GCP.
// This instance is a no-op method.
func (b *Broker) UpdateInstanceDetails(ctx context.Context, instance *models.ServiceInstanceDetails) error {
	client, err := b.createClient(ctx)
	if err != nil {
		return err
	}

	actualInstance, err := client.Projects.Locations.Instances.Get(b.instancePath(*instance)).Do()
	if err != nil {
		return err
	}

	instanceInfo, err := NewInstanceInformation(*actualInstance)
	if err != nil {
		return err
	}

	return instance.SetOtherDetails(*instanceInfo)
}

// ProvisionsAsync indicates if provisioning must be done asynchronously.
func (b *Broker) ProvisionsAsync() bool {
	return true
}

// DeprovisionsAsync indicates if deprovisioning must be done asynchronously.
func (b *Broker) DeprovisionsAsync() bool {
	return true
}
