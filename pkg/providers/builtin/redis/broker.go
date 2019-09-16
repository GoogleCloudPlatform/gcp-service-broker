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

package redis

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/GoogleCloudPlatform/gcp-service-broker/db_service/models"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/providers/builtin/base"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/varcontext"
	"github.com/GoogleCloudPlatform/gcp-service-broker/utils"
	"github.com/pivotal-cf/brokerapi"
	"golang.org/x/net/context"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
	redis "google.golang.org/api/redis/v1"
)

// Broker is the service-broker back-end for creating and binding Redis services.
type Broker struct {
	base.PeeredNetworkServiceBase
	base.NoOpBindMixin
	base.AsynchronousInstanceMixin
}

// NewInstanceInformation creates instance information from an instance
func NewInstanceInformation(instance redis.Instance) InstanceInformation {
	return InstanceInformation{
		Network:         instance.AuthorizedNetwork,
		ReservedIPRange: instance.ReservedIpRange,
		RedisVersion:    instance.RedisVersion,
		MemorySizeGb:    instance.MemorySizeGb,
		Host:            instance.Host,
		Port:            instance.Port,
		URI:             fmt.Sprintf("redis://%s:%d", instance.Host, instance.Port),
	}
}

// InstanceInformation holds the details needed to connect to a Redis instance after it has been provisioned
type InstanceInformation struct {
	// Info for admins to diagnose connection issues
	Network         string `json:"authorized_network"`
	ReservedIPRange string `json:"reserved_ip_range"`

	// Info for developers to diagnose client issues
	RedisVersion string `json:"redis_version"`
	MemorySizeGb int64  `json:"memory_size_gb"`

	// Connection info
	Host string `json:"host"`
	Port int64  `json:"port"`
	URI  string `json:"uri"`
}

// Provision creates a new Redis instance from the settings in the user-provided details and service plan.
func (b *Broker) Provision(ctx context.Context, provisionContext *varcontext.VarContext) (models.ServiceInstanceDetails, error) {
	details := models.ServiceInstanceDetails{
		Name:     provisionContext.GetString(base.InstanceIDKey),
		Location: provisionContext.GetString(base.RegionKey),
	}

	instance := &redis.Instance{
		Labels:            provisionContext.GetStringMapString("labels"),
		Tier:              provisionContext.GetString("tier"),
		MemorySizeGb:      int64(provisionContext.GetInt("memory_size_gb")),
		AuthorizedNetwork: provisionContext.GetString(base.AuthorizedNetworkKey),
	}

	// The API only accepts fully qualified networks.
	// If the user doesn't specify, assume they want the one for the default project.
	if !strings.Contains(instance.AuthorizedNetwork, "/") {
		instance.AuthorizedNetwork = fmt.Sprintf("projects/%s/global/networks/%s", b.DefaultProjectID, instance.AuthorizedNetwork)
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

// Deprovision deletes the Redis instance with the given instance ID
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

func (b *Broker) parentPath(instanceDetails models.ServiceInstanceDetails) string {
	return fmt.Sprintf("projects/%s/locations/%s", b.DefaultProjectID, instanceDetails.Location)
}

func (b *Broker) instancePath(instanceDetails models.ServiceInstanceDetails) string {
	return fmt.Sprintf("%s/instances/%s", b.parentPath(instanceDetails), instanceDetails.Name)
}

func (b *Broker) createClient(ctx context.Context) (*redis.Service, error) {
	co := option.WithUserAgent(utils.CustomUserAgent)
	ct := option.WithTokenSource(b.HTTPConfig.TokenSource(ctx))
	c, err := redis.NewService(ctx, co, ct)
	if err != nil {
		return nil, fmt.Errorf("couldn't instantiate API client: %s", err)
	}

	return c, nil
}

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

	return instance.SetOtherDetails(NewInstanceInformation(*actualInstance))
}
