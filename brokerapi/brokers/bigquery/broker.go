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
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/varcontext"
	"github.com/GoogleCloudPlatform/gcp-service-broker/utils"
	"github.com/pivotal-cf/brokerapi"
	googlebigquery "google.golang.org/api/bigquery/v2"
)

// BigQueryBroker is the service-broker back-end for creating and binding BigQuery instances.
type BigQueryBroker struct {
	broker_base.BrokerBase
}

// InstanceInformation holds the details needed to bind a service account to a BigQuery instance after it has been provisioned.
type InstanceInformation struct {
	DatasetId string `json:"dataset_id"`
}

// Provision creates a new BigQuery dataset from the settings in the user-provided details and service plan.
func (b *BigQueryBroker) Provision(ctx context.Context, instanceId string, details brokerapi.ProvisionDetails, plan models.ServicePlan) (models.ServiceInstanceDetails, error) {
	variableContext, err := varcontext.Builder().
		MergeDefaults(serviceDefinition().ProvisionInputVariables).
		MergeJsonObject(details.GetRawParameters()).
		MergeMap(plan.GetServiceProperties()).Build()

	if err != nil {
		return models.ServiceInstanceDetails{}, err
	}

	service, err := b.createClient(ctx)
	if err != nil {
		return models.ServiceInstanceDetails{}, err
	}

	d := googlebigquery.Dataset{
		Location: variableContext.GetString("location"),
		DatasetReference: &googlebigquery.DatasetReference{
			DatasetId: variableContext.GetString("name"),
		},
		Labels: utils.ExtractDefaultLabels(instanceId, details),
	}

	if variableContext.HasErrors() {
		return models.ServiceInstanceDetails{}, variableContext.Error()
	}

	newDataset, err := service.Datasets.Insert(b.ProjectId, &d).Do()
	if err != nil {
		return models.ServiceInstanceDetails{}, fmt.Errorf("Error inserting new dataset: %s", err)
	}

	ii := InstanceInformation{
		DatasetId: newDataset.DatasetReference.DatasetId,
	}

	otherDetails, err := json.Marshal(ii)
	if err != nil {
		return models.ServiceInstanceDetails{}, fmt.Errorf("Error marshalling other details: %s", err)
	}

	i := models.ServiceInstanceDetails{
		Name:         newDataset.DatasetReference.DatasetId,
		Url:          newDataset.SelfLink,
		Location:     newDataset.Location,
		OtherDetails: string(otherDetails),
	}

	return i, nil
}

// Deprovision deletes the dataset associated with the given instance.
// Note: before deprovisioning you must delete all the tables in the dataset.
func (b *BigQueryBroker) Deprovision(ctx context.Context, dataset models.ServiceInstanceDetails, details brokerapi.DeprovisionDetails) error {
	service, err := b.createClient(ctx)
	if err != nil {
		return err
	}

	if err = service.Datasets.Delete(b.ProjectId, dataset.Name).Do(); err != nil {
		return fmt.Errorf("Error deleting dataset: %s", err)
	}

	return nil
}

func (b *BigQueryBroker) createClient(ctx context.Context) (*googlebigquery.Service, error) {
	service, err := googlebigquery.New(b.HttpConfig.Client(ctx))
	if err != nil {
		return nil, fmt.Errorf("Couldn't instantiate BigQuery API client: %s", err)
	}
	service.UserAgent = models.CustomUserAgent

	return service, nil
}
