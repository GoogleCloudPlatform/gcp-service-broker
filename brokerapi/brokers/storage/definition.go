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

package storage

import "github.com/GoogleCloudPlatform/gcp-service-broker/pkg/broker"
import accountmanagers "github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/account_managers"

var roleWhitelist = []string{
	"storage.objectCreator",
	"storage.objectViewer",
	"storage.objectAdmin",
}

func init() {
	bs := &broker.BrokerService{
		Name: "google-storage",
		DefaultServiceDefinition: `{
	        "id": "b9e4332e-b42b-4680-bda5-ea1506797474",
	        "description": "A Powerful, Simple and Cost Effective Object Storage Service",
	        "name": "google-storage",
	        "bindable": true,
	        "plan_updateable": false,
	        "metadata": {
	          "displayName": "Google Cloud Storage",
	          "longDescription": "A Powerful, Simple and Cost Effective Object Storage Service",
	          "documentationUrl": "https://cloud.google.com/storage/docs/overview",
	          "supportUrl": "https://cloud.google.com/support/",
	          "imageUrl": "https://cloud.google.com/_static/images/cloud/products/logos/svg/storage.svg"
	        },
	        "tags": ["gcp", "storage"],
	        "plans": [
	          {
	            "id": "e1d11f65-da66-46ad-977c-6d56513baf43",
	            "service_id": "b9e4332e-b42b-4680-bda5-ea1506797474",
	            "name": "standard",
	            "display_name": "Standard",
	            "description": "Standard storage class",
	            "service_properties": {"storage_class": "STANDARD"}
	          },
	          {
	            "id": "a42c1182-d1a0-4d40-82c1-28220518b360",
	            "service_id": "b9e4332e-b42b-4680-bda5-ea1506797474",
	            "name": "nearline",
	            "display_name": "Nearline",
	            "description": "Nearline storage class",
	            "service_properties": {"storage_class": "NEARLINE"}
	          },
	          {
	            "id": "1a1f4fe6-1904-44d0-838c-4c87a9490a6b",
	            "service_id": "b9e4332e-b42b-4680-bda5-ea1506797474",
	            "name": "reduced-availability",
	            "display_name": "Durable Reduced Availability",
	            "description": "Durable Reduced Availability storage class",
	            "service_properties": {"storage_class": "DURABLE_REDUCED_AVAILABILITY"}
	          }
	        ]
	      }`,
		ProvisionInputVariables: []broker.BrokerVariable{
			broker.BrokerVariable{
				FieldName: "name",
				Type:      broker.JsonTypeString,
				Details:   "The name of the bucket. There is a single global namespace shared by all buckets so it MUST be unique.",
				Default:   "a generated value",
			},
			broker.BrokerVariable{
				FieldName: "location",
				Type:      broker.JsonTypeString,
				Default:   "US",
				Details:   `The location of the bucket. Object data for objects in the bucket resides in physical storage within this region. See https://cloud.google.com/storage/docs/bucket-locations`,
			},
		},
		ServiceAccountRoleWhitelist: roleWhitelist,
		BindInputVariables:          accountmanagers.ServiceAccountBindInputVariables(roleWhitelist),
		BindOutputVariables: append(accountmanagers.ServiceAccountBindOutputVariables(),
			broker.BrokerVariable{
				FieldName: "bucket_name",
				Type:      broker.JsonTypeString,
				Details:   "Name of the bucket this binding is for",
			},
		),

		Examples: []broker.ServiceExample{
			broker.ServiceExample{
				Name:            "Basic Configuration",
				Description:     "Create a nearline bucket with a service account that can create/read/delete the objects in it.",
				PlanId:          "a42c1182-d1a0-4d40-82c1-28220518b360",
				ProvisionParams: map[string]interface{}{"location": "us"},
				BindParams: map[string]interface{}{
					"role": "storage.objectAdmin",
				},
			},
		},
	}

	broker.Register(bs)
}
