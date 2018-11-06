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

import (
	"code.cloudfoundry.org/lager"
	accountmanagers "github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/account_managers"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/broker_base"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/models"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/broker"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/validation"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/varcontext"
	"golang.org/x/oauth2/jwt"
)

func init() {
	broker.Register(serviceDefinition())
}

func serviceDefinition() *broker.ServiceDefinition {
	roleWhitelist := []string{
		"storage.objectCreator",
		"storage.objectViewer",
		"storage.objectAdmin",
	}

	return &broker.ServiceDefinition{
		Name: models.StorageName,
		DefaultServiceDefinition: `{
	        "id": "b9e4332e-b42b-4680-bda5-ea1506797474",
	        "description": "Unified object storage for developers and enterprises. Cloud Storage allows world-wide storage and retrieval of any amount of data at any time.",
	        "name": "google-storage",
	        "bindable": true,
	        "plan_updateable": false,
	        "metadata": {
	          "displayName": "Google Cloud Storage",
	          "longDescription": "Unified object storage for developers and enterprises. Cloud Storage allows world-wide storage and retrieval of any amount of data at any time.",
	          "documentationUrl": "https://cloud.google.com/storage/docs/overview",
	          "supportUrl": "https://cloud.google.com/storage/docs/getting-support",
	          "imageUrl": "https://cloud.google.com/_static/images/cloud/products/logos/svg/storage.svg"
	        },
	        "tags": ["gcp", "storage"],
	        "plans": [
	          {
	            "id": "e1d11f65-da66-46ad-977c-6d56513baf43",
	            "service_id": "b9e4332e-b42b-4680-bda5-ea1506797474",
	            "name": "standard",
	            "display_name": "Standard",
	            "description": "Standard storage class.",
	            "service_properties": {"storage_class": "STANDARD"}
	          },
	          {
	            "id": "a42c1182-d1a0-4d40-82c1-28220518b360",
	            "service_id": "b9e4332e-b42b-4680-bda5-ea1506797474",
	            "name": "nearline",
	            "display_name": "Nearline",
	            "description": "Nearline storage class.",
	            "service_properties": {"storage_class": "NEARLINE"}
	          },
	          {
	            "id": "1a1f4fe6-1904-44d0-838c-4c87a9490a6b",
	            "service_id": "b9e4332e-b42b-4680-bda5-ea1506797474",
	            "name": "reduced-availability",
	            "display_name": "Durable Reduced Availability",
	            "description": "Durable Reduced Availability storage class.",
	            "service_properties": {"storage_class": "DURABLE_REDUCED_AVAILABILITY"}
	          },
	          {
	            "id": "c8538397-8f15-45e3-a229-8bb349c3a98f",
	            "name": "coldline",
	            "display_name": "Coldline Storage",
	            "description": "Google Cloud Storage Coldline is a very-low-cost, highly durable storage service for data archiving, online backup, and disaster recovery.",
	            "service_properties": {"storage_class": "COLDLINE"}
	          }
	        ]
	      }`,
		ProvisionInputVariables: []broker.BrokerVariable{
			{
				FieldName: "name",
				Type:      broker.JsonTypeString,
				Details:   "The name of the bucket. There is a single global namespace shared by all buckets so it MUST be unique.",
				Default:   "pcf_sb_${counter.next()}_${time.nano()}",
				Constraints: validation.NewConstraintBuilder(). // https://cloud.google.com/storage/docs/naming
										Pattern("^[A-Za-z0-9_\\.]+$").
										MinLength(3).
										MaxLength(222).
										Build(),
			},
			{
				FieldName: "location",
				Type:      broker.JsonTypeString,
				Default:   "US",
				Details:   `The location of the bucket. Object data for objects in the bucket resides in physical storage within this region. See: https://cloud.google.com/storage/docs/bucket-locations`,
				Constraints: validation.NewConstraintBuilder().
					Pattern("^[A-Za-z][-a-z0-9A-Z]+$").
					Examples("US", "EU", "southamerica-east1").
					Build(),
			},
		},
		ProvisionComputedVariables: []varcontext.DefaultVariable{
			{Name: "labels", Default: "${json.marshal(request.default_labels)}", Overwrite: true},
		},
		DefaultRoleWhitelist: roleWhitelist,
		BindInputVariables:   accountmanagers.ServiceAccountBindInputVariables(models.StorageName, roleWhitelist),
		BindOutputVariables: append(accountmanagers.ServiceAccountBindOutputVariables(),
			broker.BrokerVariable{
				FieldName: "bucket_name",
				Type:      broker.JsonTypeString,
				Details:   "Name of the bucket this binding is for.",
				Required:  true,
				Constraints: validation.NewConstraintBuilder(). // https://cloud.google.com/storage/docs/naming
										Pattern("^[A-Za-z0-9_\\.]+$").
										MinLength(3).
										MaxLength(222).
										Build(),
			},
		),
		PlanVariables: []broker.BrokerVariable{
			{
				FieldName: "storage_class",
				Type:      broker.JsonTypeString,
				Details:   "The storage class of the bucket. See: https://cloud.google.com/storage/docs/storage-classes.",
				Required:  true,
			},
		},
		Examples: []broker.ServiceExample{
			{
				Name:            "Basic Configuration",
				Description:     "Create a nearline bucket with a service account that can create/read/delete the objects in it.",
				PlanId:          "a42c1182-d1a0-4d40-82c1-28220518b360",
				ProvisionParams: map[string]interface{}{"location": "us"},
				BindParams: map[string]interface{}{
					"role": "storage.objectAdmin",
				},
			},
			{
				Name:            "Cold Storage",
				Description:     "Create a coldline bucket with a service account that can create/read/delete the objects in it.",
				PlanId:          "c8538397-8f15-45e3-a229-8bb349c3a98f",
				ProvisionParams: map[string]interface{}{"location": "us"},
				BindParams: map[string]interface{}{
					"role":     "storage.objectAdmin",
					"location": "us-west1",
				},
			},
		},
		ProviderBuilder: func(projectId string, auth *jwt.Config, logger lager.Logger) broker.ServiceProvider {
			bb := broker_base.NewBrokerBase(projectId, auth, logger)
			return &StorageBroker{BrokerBase: bb}
		},
	}
}
