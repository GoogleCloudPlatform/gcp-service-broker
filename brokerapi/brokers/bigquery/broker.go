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

package bigquery

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/broker_base"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/models"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/name_generator"
	"github.com/GoogleCloudPlatform/gcp-service-broker/db_service"
	"github.com/pivotal-cf/brokerapi"
	googlebigquery "google.golang.org/api/bigquery/v2"
)

// BigQueryBroker is the service-broker back-end for creating and binding BigQuery instances
type BigQueryBroker struct {
	broker_base.BrokerBase
}

// InstanceInformation holds the details needed to bind a service account to a BigQuery instance after it has been provisioned
type InstanceInformation struct {
	DatasetId string `json:"dataset_id"`
}

// Provision creates a new BigQuery dataset identified by the name provided in details.RawParameters.name and optional location
// (possible values are "US" or "EU", defaults to "US")
func (b *BigQueryBroker) Provision(instanceId string, details brokerapi.ProvisionDetails, plan models.ServicePlan) (models.ServiceInstanceDetails, error) {
	var err error
	var params map[string]string

	if len(details.RawParameters) == 0 {
		params = map[string]string{}
	} else if err = json.Unmarshal(details.RawParameters, &params); err != nil {
		return models.ServiceInstanceDetails{}, fmt.Errorf("Error unmarshalling parameters: %s", err)
	}

	// Ensure there is a name for this instance
	if _, ok := params["name"]; !ok {
		params["name"] = name_generator.Basic.InstanceName()
	}

	service, err := googlebigquery.New(b.HttpConfig.Client(context.Background()))
	if err != nil {
		return models.ServiceInstanceDetails{}, fmt.Errorf("Error creating bigquery client: %s", err)
	}
	service.UserAgent = models.CustomUserAgent

	loc := "US"
	userLoc, locOk := params["location"]
	if locOk {
		loc = userLoc
	}
	d := googlebigquery.Dataset{
		Location: loc,
		DatasetReference: &googlebigquery.DatasetReference{
			DatasetId: params["name"],
		},
	}
	new_dataset, err := service.Datasets.Insert(b.ProjectId, &d).Do()
	if err != nil {
		return models.ServiceInstanceDetails{}, fmt.Errorf("Error inserting new dataset: %s", err)
	}

	ii := InstanceInformation{
		DatasetId: params["name"],
	}

	otherDetails, err := json.Marshal(ii)
	if err != nil {
		return models.ServiceInstanceDetails{}, fmt.Errorf("Error marshalling other details: %s", err)
	}

	i := models.ServiceInstanceDetails{
		Name:         new_dataset.DatasetReference.DatasetId,
		Url:          new_dataset.SelfLink,
		Location:     new_dataset.Location,
		OtherDetails: string(otherDetails),
	}

	return i, nil
}

// Deprovision deletes the dataset associated with the given instanceID string
// note that all tables in the dataset must be deleted prior to deprovisioning
func (b *BigQueryBroker) Deprovision(instanceID string, details brokerapi.DeprovisionDetails) error {
	var err error
	service, err := googlebigquery.New(b.HttpConfig.Client(context.Background()))
	if err != nil {
		return fmt.Errorf("Error creating BigQuery client: %s", err)
	}

	dataset := models.ServiceInstanceDetails{}
	if err = db_service.DbConnection.Where("ID = ?", instanceID).First(&dataset).Error; err != nil {
		return brokerapi.ErrInstanceDoesNotExist
	}

	if err = service.Datasets.Delete(b.ProjectId, dataset.Name).Do(); err != nil {
		return fmt.Errorf("Error deleting dataset: %s", err)
	}

	return nil
}
