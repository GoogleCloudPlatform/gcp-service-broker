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
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/broker"
	accountmanagers "github.com/GoogleCloudPlatform/gcp-service-broker/pkg/providers/builtin/account_managers"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/providers/builtin/base"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/validation"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/varcontext"
	"github.com/pivotal-cf/brokerapi"
	"golang.org/x/oauth2/jwt"
)

const StorageName = "google-storage"

// ServiceDefinition creates a new ServiceDefinition object for the Cloud Storage service.
func ServiceDefinition() *broker.ServiceDefinition {
	roleWhitelist := []string{
		"storage.objectCreator",
		"storage.objectViewer",
		"storage.objectAdmin",
	}

	return &broker.ServiceDefinition{
		Id:               "b9e4332e-b42b-4680-bda5-ea1506797474",
		Name:             StorageName,
		Description:      "Unified object storage for developers and enterprises. Cloud Storage allows world-wide storage and retrieval of any amount of data at any time.",
		DisplayName:      "Google Cloud Storage",
		ImageUrl:         "https://cloud.google.com/_static/images/cloud/products/logos/svg/storage.svg",
		DocumentationUrl: "https://cloud.google.com/storage/docs/overview",
		SupportUrl:       "https://cloud.google.com/storage/docs/getting-support",
		Tags:             []string{"gcp", "storage"},
		Bindable:         true,
		PlanUpdateable:   false,
		Plans: []broker.ServicePlan{
			{
				ServicePlan: brokerapi.ServicePlan{
					ID:          "e1d11f65-da66-46ad-977c-6d56513baf43",
					Name:        "standard",
					Description: "Standard storage class. Auto-selects either regional or multi-regional based on the location.",
					Free:        brokerapi.FreeValue(false),
				},
				ServiceProperties: map[string]string{"storage_class": "STANDARD"},
			},
			{
				ServicePlan: brokerapi.ServicePlan{
					ID:          "a42c1182-d1a0-4d40-82c1-28220518b360",
					Name:        "nearline",
					Description: "Nearline storage class.",
					Free:        brokerapi.FreeValue(false),
				},
				ServiceProperties: map[string]string{"storage_class": "NEARLINE"},
			},
			{
				ServicePlan: brokerapi.ServicePlan{
					ID:          "1a1f4fe6-1904-44d0-838c-4c87a9490a6b",
					Name:        "reduced-availability",
					Description: "Durable Reduced Availability storage class.",
					Free:        brokerapi.FreeValue(false),
				},
				ServiceProperties: map[string]string{"storage_class": "DURABLE_REDUCED_AVAILABILITY"},
			},
			{
				ServicePlan: brokerapi.ServicePlan{
					ID:          "c8538397-8f15-45e3-a229-8bb349c3a98f",
					Name:        "coldline",
					Description: "Google Cloud Storage Coldline is a very-low-cost, highly durable storage service for data archiving, online backup, and disaster recovery.",
					Free:        brokerapi.FreeValue(false),
				},
				ServiceProperties: map[string]string{"storage_class": "COLDLINE"},
			},
			{
				ServicePlan: brokerapi.ServicePlan{
					ID:          "5e6161d2-0202-48be-80c4-1006cce19b9d",
					Name:        "regional",
					Description: "Data is stored in a narrow geographic region, redundant across availability zones with a 99.99% typical monthly availability.",
					Free:        brokerapi.FreeValue(false),
				},
				ServiceProperties: map[string]string{"storage_class": "REGIONAL"},
			},
			{
				ServicePlan: brokerapi.ServicePlan{
					ID:          "a5e8dfb5-e5ec-472a-8d36-33afcaff2fdb",
					Name:        "multiregional",
					Description: "Data is stored geo-redundantly with >99.99% typical monthly availability.",
					Free:        brokerapi.FreeValue(false),
				},
				ServiceProperties: map[string]string{"storage_class": "MULTI_REGIONAL"},
			},
		},
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
			{
				FieldName:   "only_delete_if_empty",
				Type:        broker.JsonTypeString,
				Default:     "true",
				Details:     `Should this bucket only delete if it's empty?`,
				Constraints: validation.NewConstraintBuilder().Enum("true", "false").Build(),
			},
		},
		ProvisionComputedVariables: []varcontext.DefaultVariable{
			{Name: "labels", Default: "${json.marshal(request.default_labels)}", Overwrite: true},
		},
		DefaultRoleWhitelist: roleWhitelist,
		BindInputVariables:   accountmanagers.ServiceAccountWhitelistWithDefault(roleWhitelist, "storage.objectAdmin"),
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
				Description:     "Create a nearline bucket with a service account that can create/read/list/delete the objects in it.",
				PlanId:          "a42c1182-d1a0-4d40-82c1-28220518b360",
				ProvisionParams: map[string]interface{}{"location": "us"},
				BindParams: map[string]interface{}{
					"role": "storage.objectAdmin",
				},
			},
			{
				Name:            "Cold Storage",
				Description:     "Create a coldline bucket with a service account that can create/read/list/delete the objects in it.",
				PlanId:          "c8538397-8f15-45e3-a229-8bb349c3a98f",
				ProvisionParams: map[string]interface{}{"location": "us"},
				BindParams: map[string]interface{}{
					"role": "storage.objectAdmin",
				},
			},
			{
				Name:            "Regional Storage",
				Description:     "Create a regional bucket with a service account that can create/read/list/delete the objects in it.",
				PlanId:          "5e6161d2-0202-48be-80c4-1006cce19b9d",
				ProvisionParams: map[string]interface{}{"location": "us-west1"},
				BindParams: map[string]interface{}{
					"role": "storage.objectAdmin",
				},
			},
			{
				Name:            "Multi-Regional Storage",
				Description:     "Create a multi-regional bucket with a service account that can create/read/list/delete the objects in it.",
				PlanId:          "a5e8dfb5-e5ec-472a-8d36-33afcaff2fdb",
				ProvisionParams: map[string]interface{}{"location": "us"},
				BindParams: map[string]interface{}{
					"role": "storage.objectAdmin",
				},
			},
		},
		BindComputedVariables: accountmanagers.ServiceAccountBindComputedVariables(),
		ProviderBuilder: func(projectId string, auth *jwt.Config, logger lager.Logger) broker.ServiceProvider {
			bb := base.NewBrokerBase(projectId, auth, logger)
			return &StorageBroker{BrokerBase: bb}
		},
		IsBuiltin: true,
	}
}
