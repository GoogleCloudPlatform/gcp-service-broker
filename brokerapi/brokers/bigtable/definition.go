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
	accountmanagers "github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/account_managers"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/broker"
)

func init() {
	roleWhitelist := []string{
		"bigtable.user",
		"bigtable.reader",
		"bigtable.viewer",
	}

	bs := &broker.BrokerService{
		Name: "google-bigtable",
		DefaultServiceDefinition: `{
      "id": "b8e19880-ac58-42ef-b033-f7cd9c94d1fe",
      "description": "A high performance NoSQL database service for large analytical and operational workloads",
      "name": "google-bigtable",
      "bindable": true,
      "plan_updateable": false,
      "metadata": {
          "displayName": "Google Bigtable",
          "longDescription": "A high performance NoSQL database service for large analytical and operational workloads",
          "documentationUrl": "https://cloud.google.com/bigtable/",
          "supportUrl": "https://cloud.google.com/support/",
          "imageUrl": "https://cloud.google.com/_static/images/cloud/products/logos/svg/bigtable.svg"
      },
      "tags": ["gcp", "bigtable"],
      "plans": [
        {
          "id": "65a49268-2c73-481e-80f3-9fde5bd5a654",
          "name": "three-node-production-hdd",
          "description": "BigTable HDD basic production plan: Approx: Reads: 1,500 QPS @ 200ms or Writes: 30,000 QPS @ 50ms or Scans: 540 MB/s, 24TB storage",
          "service_properties": {
            "storage_type": "HDD",
            "num_nodes": "3"
          },
          "display_name": "3 Node HDD",
          "service_id": "b8e19880-ac58-42ef-b033-f7cd9c94d1fe"
        },
        {
          "id": "38aa0e65-624b-4998-9c06-f9194b56d252",
          "name": "three-node-production-ssd",
          "description": "BigTable SSD basic production plan: Approx: Reads: 30,000 QPS @ 6ms or Writes: 30,000 QPS @ 6ms or Scans: 660 MB/s, 7.5TB storage",
          "service_properties": {
            "storage_type": "SSD",
            "num_nodes": "3"
          },
          "display_name": "3 Node SSD",
          "service_id": "b8e19880-ac58-42ef-b033-f7cd9c94d1fe"
        }
      ]
    }`,
		ProvisionInputVariables: []broker.BrokerVariable{
			broker.BrokerVariable{
				FieldName: "name",
				Type:      broker.JsonTypeString,
				Details:   "The name of the dataset. Should match [a-z][a-z0-9\\-]+[a-z0-9]",
				Default:   "a generated value",
			},
			broker.BrokerVariable{
				FieldName: "cluster_id",
				Type:      broker.JsonTypeString,
				Details:   "The name of the cluster.",
				Default:   "a generated value",
			},
			broker.BrokerVariable{
				FieldName: "display_name",
				Type:      broker.JsonTypeString,
				Details:   "The human-readable name of the dataset.",
				Default:   "a generated value",
			},
			broker.BrokerVariable{
				FieldName: "zone",
				Type:      broker.JsonTypeString,
				Details:   "The zone the data will reside in.",
				Default:   "us-east1-b",
			},
		},
		ServiceAccountRoleWhitelist: roleWhitelist,
		BindInputVariables:          accountmanagers.ServiceAccountBindInputVariables(roleWhitelist),
		BindOutputVariables: append(accountmanagers.ServiceAccountBindOutputVariables(),
			broker.BrokerVariable{
				FieldName: "instance_id",
				Type:      broker.JsonTypeString,
				Details:   "The name of the BigTable dataset",
			},
		),
		PlanVariables: []broker.BrokerVariable{
			broker.BrokerVariable{
				FieldName: "storage_type",
				Type:      broker.JsonTypeString,
				Details:   "Either HDD or SSD (see https://cloud.google.com/bigtable/pricing for more information)",
				Default:   "SSD",
				Required:  true,
				Enum: map[interface{}]string{
					"SSD": "SSD - Solid-state Drive",
					"HDD": "HDD - Hard Disk Drive",
				},
			},
			broker.BrokerVariable{
				FieldName: "num_nodes",
				Type:      broker.JsonTypeString,
				Details:   "Number of Nodes, Between 3 and 30 (see https://cloud.google.com/bigtable/pricing for more information)",
				Default:   "3",
				Required:  true,
			},
		},
		Examples: []broker.ServiceExample{
			broker.ServiceExample{
				Name:        "Basic Production Configuration",
				Description: "Create an HDD production table and account that can manage and query the data.",
				PlanId:      "65a49268-2c73-481e-80f3-9fde5bd5a654",
				ProvisionParams: map[string]interface{}{
					"name": "orders-table",
				},
				BindParams: map[string]interface{}{
					"role": "bigtable.user",
				},
			},
		},
	}

	broker.Register(bs)
}
