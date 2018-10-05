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
	accountmanagers "github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/account_managers"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/broker"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/validation"
)

func init() {
	broker.Register(serviceDefinition())
}

func serviceDefinition() *broker.BrokerService {
	roleWhitelist := []string{
		"bigquery.dataViewer",
		"bigquery.dataEditor",
		"bigquery.dataOwner",
		"bigquery.user",
		"bigquery.jobUser",
	}

	return &broker.BrokerService{
		Name: "google-bigquery",
		DefaultServiceDefinition: `{
        "id": "f80c0a3e-bd4d-4809-a900-b4e33a6450f1",
        "description": "A fast, economical and fully managed data warehouse for large-scale data analytics.",
        "name": "google-bigquery",
        "bindable": true,
        "plan_updateable": false,
        "metadata": {
          "displayName": "Google BigQuery",
          "longDescription": "A fast, economical and fully managed data warehouse for large-scale data analytics.",
          "documentationUrl": "https://cloud.google.com/bigquery/docs/",
          "supportUrl": "https://cloud.google.com/support/",
          "imageUrl": "https://cloud.google.com/_static/images/cloud/products/logos/svg/bigquery.svg"
        },
        "tags": ["gcp", "bigquery"],
        "plans": [
          {
            "id": "10ff4e72-6e84-44eb-851f-bdb38a791914",
            "service_id": "f80c0a3e-bd4d-4809-a900-b4e33a6450f1",
            "name": "default",
            "display_name": "Default",
            "description": "BigQuery default plan.",
            "service_properties": {}
          }
        ]
      }`,
		ProvisionInputVariables: []broker.BrokerVariable{
			{
				FieldName: "name",
				Type:      broker.JsonTypeString,
				Details:   "The name of the BigQuery dataset.",
				Default:   "pcf-sb-${counter.next()}-${time.nano()}",
				Constraints: validation.NewConstraintBuilder().
					Pattern("^[A-Za-z0-9_]+$").
					MaxLength(1024).
					Build(),
			},
			{
				FieldName: "location",
				Type:      broker.JsonTypeString,
				Details:   "The location of the BigQuery instance.",
				Default:   "US",
				Constraints: validation.NewConstraintBuilder().
					Pattern("^[A-Za-z][-a-z0-9A-Z]+$").
					Examples("US", "EU", "asia-northeast1").
					Build(),
			},
		},
		DefaultRoleWhitelist: roleWhitelist,
		BindInputVariables:   accountmanagers.ServiceAccountBindInputVariables(roleWhitelist),
		BindOutputVariables: append(accountmanagers.ServiceAccountBindOutputVariables(),
			broker.BrokerVariable{
				FieldName: "dataset_id",
				Type:      broker.JsonTypeString,
				Details:   "The name of the BigQuery dataset.",
			},
		),
		Examples: []broker.ServiceExample{
			{
				Name:        "Basic Configuration",
				Description: "Create a dataset and account that can manage and query the data.",
				PlanId:      "10ff4e72-6e84-44eb-851f-bdb38a791914",
				ProvisionParams: map[string]interface{}{
					"name": "orders_1997",
				},
				BindParams: map[string]interface{}{
					"role": "bigquery.user",
				},
			},
		},
	}
}
