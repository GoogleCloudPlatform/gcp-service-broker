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

package spanner

import "github.com/GoogleCloudPlatform/gcp-service-broker/pkg/broker"
import accountmanagers "github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/account_managers"

func init() {
	roleWhitelist := []string{
		"spanner.databaseAdmin",
		"spanner.databaseReader",
		"spanner.databaseUser",
		"spanner.viewer",
	}

	bs := &broker.BrokerService{
		Name: "google-spanner",
		DefaultServiceDefinition: `
		{
			"id": "51b3e27e-d323-49ce-8c5f-1211e6409e82",
			"description": "The first horizontally scalable, globally consistent, relational database service",
			"name": "google-spanner",
			"bindable": true,
			"plan_updateable": false,
			"metadata": {
				"displayName": "Google Spanner",
				"longDescription": "The first horizontally scalable, globally consistent, relational database service",
				"documentationUrl": "https://cloud.google.com/spanner/",
				"supportUrl": "https://cloud.google.com/support/",
				"imageUrl": "https://cloud.google.com/_static/images/cloud/products/logos/svg/spanner.svg"
			},
			"tags": ["gcp", "spanner"],
			"plans": [
				{
					"id": "44828436-cfbd-47ae-b4bc-48854564347b",
					"name": "sandbox",
					"description": "Useful for testing, not eligible for SLA",
					"free": false,
					"service_properties": {"num_nodes": "1"}
				},
				{
					"id": "0752b1ad-a784-4dcc-96eb-64149089a1c9",
					"name": "minimal-production",
					"description": "A minimal production level Spanner setup eligible for 99.99% SLA. Each node can provide up to 10,000 QPS of reads or 2,000 QPS of writes (writing single rows at 1KB data per row), and 2 TiB storage.",
					"free": false,
					"service_properties": {"num_nodes": "3"}
				}
			]
		}`,
		ProvisionInputVariables: []broker.BrokerVariable{
			broker.BrokerVariable{
				FieldName: "name",
				Type:      broker.JsonTypeString,
				Details:   "The name of the instance.",
				Default:   "a generated value",
			},
			broker.BrokerVariable{
				FieldName: "display_name",
				Type:      broker.JsonTypeString,
				Details:   "A human-readable name for the instance.",
				Default:   "a generated value",
			},
			broker.BrokerVariable{
				FieldName: "location",
				Type:      broker.JsonTypeString,
				Default:   "regional-us-central1",
				Details:   `The location of the Spanner instance.`,
			},
		},
		DefaultRoleWhitelist: roleWhitelist,
		BindInputVariables:          accountmanagers.ServiceAccountBindInputVariables(roleWhitelist),
		BindOutputVariables: append(accountmanagers.ServiceAccountBindOutputVariables(),
			broker.BrokerVariable{
				FieldName: "instance_id",
				Type:      broker.JsonTypeString,
				Details:   "Name of the spanner instance the account can connect to.",
			},
		),
		PlanVariables: []broker.BrokerVariable{
			broker.BrokerVariable{
				FieldName: "num_nodes",
				Type:      broker.JsonTypeString,
				Details:   "Number of Nodes, A minimum of 3 nodes is recommended for production environments. (see https://cloud.google.com/spanner/pricing for more information)",
				Default:   "1",
				Required:  true,
			},
		},
		Examples: []broker.ServiceExample{
			broker.ServiceExample{
				Name:            "Basic Configuration",
				Description:     "Create a sandbox environment with a database admin account",
				PlanId:          "44828436-cfbd-47ae-b4bc-48854564347b",
				ProvisionParams: map[string]interface{}{"name": "auth-database"},
				BindParams:      map[string]interface{}{"role": "spanner.databaseAdmin"},
			},
		},
	}

	broker.Register(bs)
}
