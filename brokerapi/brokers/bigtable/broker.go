// Copyright 2018 the Service Broker Project Authors.
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

package bigtable

import (
	"encoding/json"
	"errors"
	"fmt"

	googlebigtable "cloud.google.com/go/bigtable"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/broker_base"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/models"
	"github.com/pivotal-cf/brokerapi"
	"golang.org/x/net/context"
	"google.golang.org/api/option"
)

// BigTableBroker is the service-broker back-end for creating and binding BigTable instances.
type BigTableBroker struct {
	broker_base.BrokerBase
}

// InstanceInformation holds the details needed to bind a service account to a
// BigTable instance after it has been provisioned.
type InstanceInformation struct {
	InstanceId string `json:"instance_id"`
}

// storageTypes holds the valid value mapping for string storage types to their
// REST call equivalent.
var storageTypes = map[string]googlebigtable.StorageType{
	"SSD": googlebigtable.SSD,
	"HDD": googlebigtable.HDD,
}

// Provision creates a new Bigtable instance from the settings in the user-provided details and service plan.
func (b *BigTableBroker) Provision(ctx context.Context, instanceId string, details brokerapi.ProvisionDetails, plan models.ServicePlan) (models.ServiceInstanceDetails, error) {
	provisionContext, err := serviceDefinition().ProvisionVariables(instanceId, details, plan)
	if err != nil {
		return models.ServiceInstanceDetails{}, err
	}

	instanceName := provisionContext.GetString("name")

	ic := googlebigtable.InstanceConf{
		InstanceId:  instanceName,
		ClusterId:   provisionContext.GetString("cluster_id"),
		NumNodes:    int32(provisionContext.GetInt("num_nodes")),
		StorageType: storageTypes[provisionContext.GetString("storage_type")],
		Zone:        provisionContext.GetString("zone"),
		DisplayName: provisionContext.GetString("display_name"),
	}

	if err := provisionContext.Error(); err != nil {
		return models.ServiceInstanceDetails{}, err
	}

	// custom constraints
	if instanceName == "" {
		return models.ServiceInstanceDetails{}, errors.New("name must not be empty")
	}

	service, err := b.createClient(ctx)
	if err != nil {
		return models.ServiceInstanceDetails{}, err
	}

	if err := service.CreateInstance(ctx, &ic); err != nil {
		return models.ServiceInstanceDetails{}, fmt.Errorf("Error creating new Bigtable instance: %s", err)
	}

	ii := InstanceInformation{
		InstanceId: instanceName,
	}

	otherDetails, err := json.Marshal(ii)
	if err != nil {
		return models.ServiceInstanceDetails{}, fmt.Errorf("Error marshalling other details: %s", err)
	}

	return models.ServiceInstanceDetails{
		Name:         instanceName,
		Url:          "",
		Location:     "",
		OtherDetails: string(otherDetails),
	}, nil
}

// Deprovision deletes the BigTable associated with the given instance.
func (b *BigTableBroker) Deprovision(ctx context.Context, instance models.ServiceInstanceDetails, details brokerapi.DeprovisionDetails) (*string, error) {
	service, err := b.createClient(ctx)
	if err != nil {
		return nil, err
	}

	if err := service.DeleteInstance(ctx, instance.Name); err != nil {
		return nil, fmt.Errorf("Error deleting Bigtable instance: %s", err)
	}

	return nil, nil
}

func (b *BigTableBroker) createClient(ctx context.Context) (*googlebigtable.InstanceAdminClient, error) {
	co := option.WithUserAgent(models.CustomUserAgent)
	ct := option.WithTokenSource(b.HttpConfig.TokenSource(ctx))
	client, err := googlebigtable.NewInstanceAdminClient(ctx, b.ProjectId, ct, co)
	if err != nil {
		return nil, fmt.Errorf("Couldn't instantiate Bigtable API client: %s", err)
	}

	return client, nil
}
