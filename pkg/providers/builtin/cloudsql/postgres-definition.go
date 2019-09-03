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

package cloudsql

import (
	"code.cloudfoundry.org/lager"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/providers/builtin/base"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/broker"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/validation"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/varcontext"
	"github.com/pivotal-cf/brokerapi"
	"golang.org/x/oauth2/jwt"
)

const PostgresServiceId = "cbad6d78-a73c-432d-b8ff-b219a17a803a"

// PostgresServiceDefinition creates a new ServiceDefinition object for the PostgreSQL service.
func PostgresServiceDefinition() *broker.ServiceDefinition {
	return &broker.ServiceDefinition{
		Id:               PostgresServiceId,
		Name:             "google-cloudsql-postgres",
		Description:      "Google CloudSQL for PostgreSQL is a fully-managed PostgreSQL database service.",
		DisplayName:      "Google CloudSQL for PostgreSQL",
		ImageUrl:         "https://cloud.google.com/_static/images/cloud/products/logos/svg/sql.svg",
		DocumentationUrl: "https://cloud.google.com/sql/docs/",
		SupportUrl:       "https://cloud.google.com/support/",
		Tags:             []string{"gcp", "cloudsql", "postgres"},
		Bindable:         true,
		PlanUpdateable:   false,
		Plans: []broker.ServicePlan{
			{
				ServicePlan: brokerapi.ServicePlan{
					ID:          "2513d4d9-684b-4c3c-add4-6404969006de",
					Name:        "postgres-db-f1-micro",
					Description: "PostgreSQL on a db-f1-micro (Shared CPUs, 0.6 GB/RAM, 3062 GB/disk, 250 Connections)",
					Free:        brokerapi.FreeValue(false),
				},
				ServiceProperties: map[string]string{"max_disk_size": "3062", "tier": "db-f1-micro"},
			},
			{
				ServicePlan: brokerapi.ServicePlan{
					ID:          "6c1174d8-243c-44d1-b7a8-e94a779f67f5",
					Name:        "postgres-db-g1-small",
					Description: "PostgreSQL on a db-g1-small (Shared CPUs, 1.7 GB/RAM, 3062 GB/disk, 1,000 Connections)",
					Free:        brokerapi.FreeValue(false),
				},
				ServiceProperties: map[string]string{"tier": "db-g1-small", "max_disk_size": "3062"},
			},
			{
				ServicePlan: brokerapi.ServicePlan{
					ID:          "c4e68ab5-34ca-4d02-857d-3e6b3ab079a7",
					Name:        "postgres-db-n1-standard-1",
					Description: "PostgreSQL with 1 CPU, 3.75 GB/RAM, 10230 GB/disk, supporting 4,000 connections.",
					Free:        brokerapi.FreeValue(false),
				},
				ServiceProperties: map[string]string{"tier": "db-custom-1-3840", "max_disk_size": "10230"},
			},
			{
				ServicePlan: brokerapi.ServicePlan{
					ID:          "3f578ecf-885c-4b60-b38b-60272f34e00f",
					Name:        "postgres-db-n1-standard-2",
					Description: "PostgreSQL with 2 CPUs, 7.5 GB/RAM, 10230 GB/disk, supporting 4,000 connections.",
					Free:        brokerapi.FreeValue(false),
				},
				ServiceProperties: map[string]string{"tier": "db-custom-2-7680", "max_disk_size": "10230"},
			},
			{
				ServicePlan: brokerapi.ServicePlan{
					ID:          "b7fcab5d-d66d-4e82-af16-565e84cef7f9",
					Name:        "postgres-db-n1-standard-4",
					Description: "PostgreSQL with 4 CPUs, 15 GB/RAM, 10230 GB/disk, supporting 4,000 connections.",
					Free:        brokerapi.FreeValue(false),
				},
				ServiceProperties: map[string]string{"tier": "db-custom-4-15360", "max_disk_size": "10230"},
			},
			{
				ServicePlan: brokerapi.ServicePlan{
					ID:          "4b2fa14a-caf1-42e0-bd8c-3342502008a8",
					Name:        "postgres-db-n1-standard-8",
					Description: "PostgreSQL with 8 CPUs, 30 GB/RAM, 10230 GB/disk, supporting 4,000 connections.",
					Free:        brokerapi.FreeValue(false),
				},
				ServiceProperties: map[string]string{"tier": "db-custom-8-30720", "max_disk_size": "10230"},
			},
			{
				ServicePlan: brokerapi.ServicePlan{
					ID:          "ca2e770f-bfa5-4fb7-a249-8b943c3474ca",
					Name:        "postgres-db-n1-standard-16",
					Description: "PostgreSQL with 16 CPUs, 60 GB/RAM, 10230 GB/disk, supporting 4,000 connections.",
					Free:        brokerapi.FreeValue(false),
				},
				ServiceProperties: map[string]string{"tier": "db-custom-16-61440", "max_disk_size": "10230"},
			},
			{
				ServicePlan: brokerapi.ServicePlan{
					ID:          "b44f8294-b003-4a50-80c2-706858073f44",
					Name:        "postgres-db-n1-standard-32",
					Description: "PostgreSQL with 32 CPUs, 120 GB/RAM, 10230 GB/disk, supporting 4,000 connections.",
					Free:        brokerapi.FreeValue(false),
				},
				ServiceProperties: map[string]string{"max_disk_size": "10230", "tier": "db-custom-32-122880"},
			},
			{
				ServicePlan: brokerapi.ServicePlan{
					ID:          "d97326e0-5af2-4da5-b970-b4772d59cded",
					Name:        "postgres-db-n1-standard-64",
					Description: "PostgreSQL with 64 CPUs, 240 GB/RAM, 10230 GB/disk, supporting 4,000 connections.",
					Free:        brokerapi.FreeValue(false),
				},
				ServiceProperties: map[string]string{"tier": "db-custom-64-245760", "max_disk_size": "10230"},
			},
			{
				ServicePlan: brokerapi.ServicePlan{
					ID:          "c10f8691-02f5-44eb-989f-7217393012ca",
					Name:        "postgres-db-n1-highmem-2",
					Description: "PostgreSQL with 2 CPUs, 13 GB/RAM, 10230 GB/disk, supporting 4,000 connections.",
					Free:        brokerapi.FreeValue(false),
				},
				ServiceProperties: map[string]string{"tier": "db-custom-2-13312", "max_disk_size": "10230"},
			},
			{
				ServicePlan: brokerapi.ServicePlan{
					ID:          "610cc78d-d26a-41a9-90b7-547a44517f03",
					Name:        "postgres-db-n1-highmem-4",
					Description: "PostgreSQL with 4 CPUs, 26 GB/RAM, 10230 GB/disk, supporting 4,000 connections.",
					Free:        brokerapi.FreeValue(false),
				},
				ServiceProperties: map[string]string{"tier": "db-custom-4-26624", "max_disk_size": "10230"},
			},
			{
				ServicePlan: brokerapi.ServicePlan{
					ID:          "2a351e8d-958d-4c4f-ae46-c984fec18740",
					Name:        "postgres-db-n1-highmem-8",
					Description: "PostgreSQL with 8 CPUs, 52 GB/RAM, 10230 GB/disk, supporting 4,000 connections.",
					Free:        brokerapi.FreeValue(false),
				},
				ServiceProperties: map[string]string{"tier": "db-custom-8-53248", "max_disk_size": "10230"},
			},
			{
				ServicePlan: brokerapi.ServicePlan{
					ID:          "51d3ca0c-9d21-447d-a395-3e0dc0659775",
					Name:        "postgres-db-n1-highmem-16",
					Description: "PostgreSQL with 16 CPUs, 104 GB/RAM, 10230 GB/disk, supporting 4,000 connections.",
					Free:        brokerapi.FreeValue(false),
				},
				ServiceProperties: map[string]string{"tier": "db-custom-16-106496", "max_disk_size": "10230"},
			},
			{
				ServicePlan: brokerapi.ServicePlan{
					ID:          "2e72b386-f7ce-4f0d-a149-9f9a851337d4",
					Name:        "postgres-db-n1-highmem-32",
					Description: "PostgreSQL with 32 CPUs, 208 GB/RAM, 10230 GB/disk, supporting 4,000 connections.",
					Free:        brokerapi.FreeValue(false),
				},
				ServiceProperties: map[string]string{"tier": "db-custom-32-212992", "max_disk_size": "10230"},
			},
			{
				ServicePlan: brokerapi.ServicePlan{
					ID:          "82602649-e4ac-4a2f-b80d-dacd745aed6a",
					Name:        "postgres-db-n1-highmem-64",
					Description: "PostgreSQL with 64 CPUs, 416 GB/RAM, 10230 GB/disk, supporting 4,000 connections.",
					Free:        brokerapi.FreeValue(false),
				},
				ServiceProperties: map[string]string{"tier": "db-custom-64-425984", "max_disk_size": "10230"},
			},
		},
		ProvisionInputVariables: append([]broker.BrokerVariable{
			{
				FieldName: "instance_name",
				Type:      broker.JsonTypeString,
				Details:   "Name of the CloudSQL instance.",
				Expression:   identifierTemplate,
				Constraints: validation.NewConstraintBuilder().
					Pattern("^[a-z][a-z0-9-]+$").
					MaxLength(75).
					Build(),
			},
			{
				FieldName: "database_name",
				Type:      broker.JsonTypeString,
				Details:   "Name of the database inside of the instance. Must be a valid identifier for your chosen database type.",
				Expression:   identifierTemplate,
			},
			{
				FieldName: "version",
				Type:      broker.JsonTypeString,
				Details:   "The database engine type and version.",
				Default:   "POSTGRES_9_6",
				Enum: map[interface{}]string{
					"POSTGRES_9_6": "PostgreSQL 9.6.X",
				},
			},
			{
				FieldName: "failover_replica_name",
				Type:      broker.JsonTypeString,
				Details:   "(only for 2nd generation instances) If specified, creates a failover replica with the given name.",
				Default:   "",
				Constraints: validation.NewConstraintBuilder().
					Pattern("^(|[a-z][a-z0-9-]+)$").
					MaxLength(75).
					Build(),
			},
			{
				FieldName: "activation_policy",
				Type:      broker.JsonTypeString,
				Details:   "The activation policy specifies when the instance is activated; it is applicable only when the instance state is RUNNABLE.",
				Default:   "ALWAYS",
				Enum: map[interface{}]string{
					"ALWAYS": "Always, instance is always on.",
					"NEVER":  "Never, instance does not turn on if a request arrives.",
				},
			},
		}, commonProvisionVariables()...),
		ProvisionComputedVariables: []varcontext.DefaultVariable{
			{Name: "labels", Expression: `${json.marshal(request.default_labels)}`, Overwrite: true},

			// legacy behavior dictates that empty values get defaults
			{Name: "instance_name", Expression: `${instance_name == "" ? "` + identifierTemplate + `" : instance_name}`, Overwrite: true},
			{Name: "database_name", Expression: `${database_name == "" ? "` + identifierTemplate + `" : database_name}`, Overwrite: true},

			// these variables are fixed for PostgreSQL
			{Name: "is_first_gen", Default: `false`, Overwrite: true},
			{Name: "binlog", Default: `false`, Overwrite: true},

			// validation
			{Name: "_", Expression: `${assert(disk_size <= max_disk_size, "disk size (${disk_size}) is greater than max allowed disk size for this plan (${max_disk_size})")}`, Overwrite: true},
		},

		DefaultRoleWhitelist:  roleWhitelist(),
		BindInputVariables:    commonBindVariables(),
		BindOutputVariables:   commonBindOutputVariables(),
		BindComputedVariables: commonBindComputedVariables(),
		PlanVariables: []broker.BrokerVariable{
			{
				FieldName: "tier",
				Type:      broker.JsonTypeString,
				Details:   "A string of the form db-custom-[CPUS]-[MEMORY_MBS], where memory is at least 3840.",
				Required:  true,
			},
			{
				FieldName: "pricing_plan",
				Type:      broker.JsonTypeString,
				Details:   "The pricing plan.",
				Enum: map[interface{}]string{
					"PER_USE": "Per-Use",
				},
				Required: true,
			},
			{
				FieldName: "max_disk_size",
				Type:      broker.JsonTypeString,
				Details:   "Maximum disk size in GB, 10 is the minimum.",
				Default:   "10",
				Required:  true,
			},
		},
		Examples: []broker.ServiceExample{
			{
				Name:        "Dedicated Machine Sandbox",
				Description: "A low end PostgreSQL sandbox that uses a dedicated machine.",
				PlanId:      "c4e68ab5-34ca-4d02-857d-3e6b3ab079a7",
				ProvisionParams: map[string]interface{}{
					"backups_enabled": "false",
					"binlog":          "false",
					"disk_size":       "25",
				},
				BindParams: map[string]interface{}{
					"role": "cloudsql.editor",
				},
			},
			{
				Name:        "Development Sandbox",
				Description: "An inexpensive PostgreSQL sandbox for developing with no backups.",
				PlanId:      "2513d4d9-684b-4c3c-add4-6404969006de",
				ProvisionParams: map[string]interface{}{
					"backups_enabled": "false",
					"binlog":          "false",
					"disk_size":       "10",
				},
				BindParams: map[string]interface{}{
					"role": "cloudsql.editor",
				},
			},
		},
		ProviderBuilder: func(projectId string, auth *jwt.Config, logger lager.Logger) broker.ServiceProvider {
			bb := base.NewBrokerBase(projectId, auth, logger)
			return &CloudSQLBroker{BrokerBase: bb}
		},
		IsBuiltin: true,
	}
}
