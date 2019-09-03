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
	"code.cloudfoundry.org/lager"
	accountmanagers "github.com/GoogleCloudPlatform/gcp-service-broker/pkg/providers/builtin/account_managers"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/providers/builtin/base"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/broker"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/validation"
	"github.com/pivotal-cf/brokerapi"
	"golang.org/x/oauth2/jwt"
)

// ServiceDefinition creates a new ServiceDefinition object for the Bigtable service.
func ServiceDefinition() *broker.ServiceDefinition {
	roleWhitelist := []string{
		"bigtable.user",
		"bigtable.reader",
		"bigtable.viewer",
	}

	return &broker.ServiceDefinition{
		Id:               "b8e19880-ac58-42ef-b033-f7cd9c94d1fe",
		Name:             "google-bigtable",
		Description:      "A high performance NoSQL database service for large analytical and operational workloads.",
		DisplayName:      "Google Bigtable",
		ImageUrl:         "https://cloud.google.com/_static/images/cloud/products/logos/svg/bigtable.svg",
		DocumentationUrl: "https://cloud.google.com/bigtable/",
		SupportUrl:       "https://cloud.google.com/bigtable/docs/support/getting-support",
		Tags:             []string{"gcp", "bigtable"},
		Bindable:         true,
		PlanUpdateable:   false,
		Plans: []broker.ServicePlan{
			{
				ServicePlan: brokerapi.ServicePlan{
					ID:          "65a49268-2c73-481e-80f3-9fde5bd5a654",
					Name:        "three-node-production-hdd",
					Description: "BigTable HDD basic production plan: Approx: Reads: 1,500 QPS @ 200ms or Writes: 30,000 QPS @ 50ms or Scans: 540 MB/s, 24TB storage.",
					Free:        brokerapi.FreeValue(false),
				},
				ServiceProperties: map[string]string{"storage_type": "HDD", "num_nodes": "3"},
			},
			{
				ServicePlan: brokerapi.ServicePlan{
					ID:          "38aa0e65-624b-4998-9c06-f9194b56d252",
					Name:        "three-node-production-ssd",
					Description: "BigTable SSD basic production plan: Approx: Reads: 30,000 QPS @ 6ms or Writes: 30,000 QPS @ 6ms or Scans: 660 MB/s, 7.5TB storage.",
					Free:        brokerapi.FreeValue(false),
				},
				ServiceProperties: map[string]string{"storage_type": "SSD", "num_nodes": "3"},
			},
		},
		ProvisionInputVariables: []broker.BrokerVariable{
			{
				FieldName: "name",
				Type:      broker.JsonTypeString,
				Details:   "The name of the Cloud Bigtable instance.",
				Expression:   "pcf-sb-${counter.next()}-${time.nano()}",
				Constraints: validation.NewConstraintBuilder().
					MinLength(6).
					MaxLength(33).
					Pattern("^[a-z][-0-9a-z]+$").
					Build(),
			},
			{
				FieldName: "cluster_id",
				Type:      broker.JsonTypeString,
				Details:   "The ID of the Cloud Bigtable cluster.",
				Expression:   "${str.truncate(20, name)}-cluster",
				Constraints: validation.NewConstraintBuilder().
					MinLength(6).
					MaxLength(30).
					Pattern("^[a-z][-0-9a-z]+[a-z]$").
					Build(),
			},
			{
				FieldName: "display_name",
				Type:      broker.JsonTypeString,
				Details:   "The human-readable display name of the Bigtable instance.",
				Expression:   "${name}",
				Constraints: validation.NewConstraintBuilder().
					MinLength(4).
					MaxLength(30).
					Build(),
			},
			{
				FieldName: "zone",
				Type:      broker.JsonTypeString,
				Details:   "The zone to create the Cloud Bigtable cluster in. Zones that support Bigtable instances are noted on the Cloud Bigtable locations page: https://cloud.google.com/bigtable/docs/locations.",
				Default:   "us-east1-b",
				Constraints: validation.NewConstraintBuilder().
					Pattern("^[A-Za-z][-a-z0-9A-Z]+$").
					Examples("us-central1-a", "europe-west2-b", "asia-northeast1-a", "australia-southeast1-c").
					Build(),
			},
		},
		DefaultRoleWhitelist: roleWhitelist,
		BindInputVariables:   accountmanagers.ServiceAccountWhitelistWithDefault(roleWhitelist, "bigtable.user"),
		BindOutputVariables: append(accountmanagers.ServiceAccountBindOutputVariables(),
			broker.BrokerVariable{
				FieldName: "instance_id",
				Type:      broker.JsonTypeString,
				Details:   "The name of the BigTable dataset.",
				Required:  true,
				Constraints: validation.NewConstraintBuilder().
					MinLength(6).
					MaxLength(33).
					Pattern("^[a-z][-0-9a-z]+$").
					Build(),
			},
		),
		BindComputedVariables: accountmanagers.ServiceAccountBindComputedVariables(),
		PlanVariables: []broker.BrokerVariable{
			{
				FieldName: "storage_type",
				Type:      broker.JsonTypeString,
				Details:   "Either HDD or SSD. See: https://cloud.google.com/bigtable/pricing for more information.",
				Default:   "SSD",
				Required:  true,
				Enum: map[interface{}]string{
					"SSD": "SSD - Solid-state Drive",
					"HDD": "HDD - Hard Disk Drive",
				},
			},
			{
				FieldName: "num_nodes",
				Type:      broker.JsonTypeString,
				Details:   "Number of nodes, between 3 and 30. See: https://cloud.google.com/bigtable/pricing for more information.",
				Default:   "3",
				Required:  true,
			},
		},
		Examples: []broker.ServiceExample{
			{
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
		ProviderBuilder: func(projectId string, auth *jwt.Config, logger lager.Logger) broker.ServiceProvider {
			bb := base.NewBrokerBase(projectId, auth, logger)
			return &BigTableBroker{BrokerBase: bb}
		},
		IsBuiltin: true,
	}
}
