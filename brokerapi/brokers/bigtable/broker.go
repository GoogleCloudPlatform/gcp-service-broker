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
	"fmt"
	"strconv"

	googlebigtable "cloud.google.com/go/bigtable"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/broker_base"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/models"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/name_generator"
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
	var err error
	var params map[string]string

	if len(details.RawParameters) == 0 {
		params = map[string]string{}
	} else if err = json.Unmarshal(details.RawParameters, &params); err != nil {
		return models.ServiceInstanceDetails{}, fmt.Errorf("Error unmarshalling parameters: %s", err)
	}

	// Ensure there is a name for this instance
	if _, ok := params["name"]; !ok {
		params["name"] = name_generator.Basic.InstanceNameWithSeparator("-")
	}

	service, err := b.createClient(ctx)
	if err != nil {
		return models.ServiceInstanceDetails{}, err
	}

	var clusterId string
	if len(params["name"]) > 20 {
		clusterId = params["name"][:20] + "-cluster"
	} else {
		clusterId = params["name"] + "-cluster"
	}
	if userClusterId, clusterIdOk := params["cluster_id"]; clusterIdOk {
		clusterId = userClusterId
	}

	numNodes, err := strconv.Atoi(plan.ServiceProperties["num_nodes"])
	if err != nil {
		return models.ServiceInstanceDetails{}, fmt.Errorf("Error converting num_nodes to int: %s", err)
	}

	zone := "us-east1-b"
	if userZone, userZoneOk := params["zone"]; userZoneOk {
		zone = userZone
	}

	displayName := params["name"]
	if userDisplayName, userDisplayNameOk := params["display_name"]; userDisplayNameOk {
		displayName = userDisplayName
	}

	ic := googlebigtable.InstanceConf{
		InstanceId:  params["name"],
		ClusterId:   clusterId,
		NumNodes:    int32(numNodes),
		StorageType: storageTypes[plan.ServiceProperties["storage_type"]],
		Zone:        zone,
		DisplayName: displayName,
	}

	err = service.CreateInstance(ctx, &ic)
	if err != nil {
		return models.ServiceInstanceDetails{}, fmt.Errorf("Error creating new instance: %s", err)
	}

	ii := InstanceInformation{
		InstanceId: params["name"],
	}

	otherDetails, err := json.Marshal(ii)
	if err != nil {
		return models.ServiceInstanceDetails{}, fmt.Errorf("Error marshalling other details: %s", err)
	}

	i := models.ServiceInstanceDetails{
		Name:         params["name"],
		Url:          "",
		Location:     "",
		OtherDetails: string(otherDetails),
	}

	return i, nil
}

// Deprovision deletes the BigTable associated with the given instance.
func (b *BigTableBroker) Deprovision(ctx context.Context, instance models.ServiceInstanceDetails, details brokerapi.DeprovisionDetails) error {
	service, err := b.createClient(ctx)
	if err != nil {
		return err
	}

	if err := service.DeleteInstance(ctx, instance.Name); err != nil {
		return fmt.Errorf("Error deleting dataset: %s", err)
	}

	return nil
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
