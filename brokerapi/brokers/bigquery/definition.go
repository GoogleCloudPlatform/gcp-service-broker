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
)

func init() {
	bs := &broker.BrokerService{
		Name: "google-bigquery",
		DefaultServiceDefinition: `{
        "id": "f80c0a3e-bd4d-4809-a900-b4e33a6450f1",
        "description": "A fast, economical and fully managed data warehouse for large-scale data analytics",
        "name": "google-bigquery",
        "bindable": true,
        "plan_updateable": false,
        "metadata": {
          "displayName": "Google BigQuery",
          "longDescription": "A fast, economical and fully managed data warehouse for large-scale data analytics",
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
            "description": "BigQuery default plan",
            "service_properties": {}
          }
        ]
      }`,
		ProvisionInputVariables: []broker.BrokerVariable{
			broker.BrokerVariable{
				FieldName: "name",
				Type:      broker.JsonTypeString,
				Details:   "The name of the BigQuery dataset. Must be alphanumeric (plus underscores) and must be at most 1024 characters long.",
				Default:   "a generated value",
			},
		},
		BindInputVariables: accountmanagers.ServiceAccountBindInputVariables(),
		BindOutputVariables: append(accountmanagers.ServiceAccountBindOutputVariables(),
			broker.BrokerVariable{
				FieldName: "dataset_id",
				Type:      broker.JsonTypeString,
				Details:   "The name of the BigQuery dataset.",
			},
		),
		Examples: []broker.ServiceExample{
			broker.ServiceExample{
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

	broker.Register(bs)
}
