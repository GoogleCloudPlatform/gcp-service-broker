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
		"spanner.databaseAdmin",
		"spanner.databaseReader",
		"spanner.databaseUser",
		"spanner.viewer",
	}

	return &broker.ServiceDefinition{
		Name: models.SpannerName,
		DefaultServiceDefinition: `
		{
			"id": "51b3e27e-d323-49ce-8c5f-1211e6409e82",
			"description": "The first horizontally scalable, globally consistent, relational database service.",
			"name": "google-spanner",
			"bindable": true,
			"plan_updateable": false,
			"metadata": {
				"displayName": "Google Spanner",
				"longDescription": "The first horizontally scalable, globally consistent, relational database service.",
				"documentationUrl": "https://cloud.google.com/spanner/",
				"supportUrl": "https://cloud.google.com/support/",
				"imageUrl": "https://cloud.google.com/_static/images/cloud/products/logos/svg/spanner.svg",
				"shareable": "true"
			},
			"tags": ["gcp", "spanner"],
			"plans": [
				{
					"id": "44828436-cfbd-47ae-b4bc-48854564347b",
					"name": "sandbox",
					"description": "Useful for testing, not eligible for SLA.",
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
			{
				FieldName: "name",
				Type:      broker.JsonTypeString,
				Details:   "A unique identifier for the instance, which cannot be changed after the instance is created.",
				Default:   "pcf-sb-${counter.next()}-${time.nano()}",
				Constraints: validation.NewConstraintBuilder().
					MinLength(6).
					MaxLength(30).
					Pattern("^[a-z][-a-z0-9]*[a-z0-9]$").
					Build(),
			},
			{
				FieldName: "display_name",
				Type:      broker.JsonTypeString,
				Details:   "The name of this instance configuration as it appears in UIs.",
				Default:   "${name}",
				Constraints: validation.NewConstraintBuilder().
					MinLength(4).
					MaxLength(30).
					Build(),
			},
			{
				FieldName: "location",
				Type:      broker.JsonTypeString,
				Default:   "regional-us-central1",
				Details: `A configuration for a Cloud Spanner instance.
				 Configurations define the geographic placement of nodes and their replication and are slightly different from zones.
				 There are single region configurations, multi-region configurations, and multi-continent configurations.
				 See the instance docs https://cloud.google.com/spanner/docs/instances for a list of configurations.`,
				Constraints: validation.NewConstraintBuilder().
					Examples("regional-asia-east1", "nam3", "nam-eur-asia1").
					Pattern("^[a-z][-a-z0-9]*[a-z0-9]$").
					Build(),
			},
		},
		ProvisionComputedVariables: []varcontext.DefaultVariable{
			{Name: "labels", Default: "${json.marshal(request.default_labels)}", Overwrite: true},
		},
		DefaultRoleWhitelist: roleWhitelist,
		BindInputVariables:   accountmanagers.ServiceAccountBindInputVariables(models.SpannerName, roleWhitelist),
		BindOutputVariables: append(accountmanagers.ServiceAccountBindOutputVariables(),
			broker.BrokerVariable{
				FieldName: "instance_id",
				Type:      broker.JsonTypeString,
				Details:   "Name of the Spanner instance the account can connect to.",
				Required:  true,
				Constraints: validation.NewConstraintBuilder().
					MinLength(6).
					MaxLength(30).
					Pattern("^[a-z][-a-z0-9]*[a-z0-9]$").
					Build(),
			},
		),
		PlanVariables: []broker.BrokerVariable{
			{
				FieldName: "num_nodes",
				Type:      broker.JsonTypeString,
				Details:   "Number of nodes, a minimum of 3 nodes is recommended for production environments. See: https://cloud.google.com/spanner/pricing for more information.",
				Default:   "1",
				Required:  true,
			},
		},
		Examples: []broker.ServiceExample{
			{
				Name:            "Basic Configuration",
				Description:     "Create a sandbox environment with a database admin account.",
				PlanId:          "44828436-cfbd-47ae-b4bc-48854564347b",
				ProvisionParams: map[string]interface{}{"name": "auth-database"},
				BindParams:      map[string]interface{}{"role": "spanner.databaseAdmin"},
			},
			{
				Name:            "99.999% availability",
				Description:     "Create a spanner instance spanning North America.",
				PlanId:          "44828436-cfbd-47ae-b4bc-48854564347b",
				ProvisionParams: map[string]interface{}{"name": "auth-database", "location": "nam3"},
				BindParams:      map[string]interface{}{"role": "spanner.databaseAdmin"},
			},
		},
		ProviderBuilder: func(projectId string, auth *jwt.Config, logger lager.Logger) broker.ServiceProvider {
			bb := broker_base.NewBrokerBase(projectId, auth, logger)
			return &SpannerBroker{BrokerBase: bb}
		},
	}
}
