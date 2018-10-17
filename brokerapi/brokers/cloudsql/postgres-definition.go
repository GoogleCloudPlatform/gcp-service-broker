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
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/broker"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/validation"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/varcontext"
)

func init() {
	broker.Register(postgresServiceDefinition())
}

func postgresServiceDefinition() *broker.BrokerService {
	return &broker.BrokerService{
		Name: "google-cloudsql-postgres",
		DefaultServiceDefinition: `{
        "id": "cbad6d78-a73c-432d-b8ff-b219a17a803a",
        "description": "Google Cloud SQL is a fully-managed PostgreSQL database service.",
        "name": "google-cloudsql-postgres",
        "bindable": true,
        "plan_updateable": false,
        "metadata": {
	          "displayName": "Google CloudSQL PostgreSQL",
	          "longDescription": "Google Cloud SQL is a fully-managed MySQL database service.",
	          "documentationUrl": "https://cloud.google.com/sql/docs/",
	          "supportUrl": "https://cloud.google.com/support/",
	          "imageUrl": "https://cloud.google.com/_static/images/cloud/products/logos/svg/sql.svg"
        },
        "tags": ["gcp", "cloudsql", "postgres"],
        "plans":[
				    {
				        "service_properties": { "tier": "db-f1-micro", "max_disk_size": "3062" },
				        "description": "PostgreSQL on a db-f1-micro (Shared CPUs, 0.6 GB/RAM, 3062 GB/disk, 250 Connections)",
				        "id": "2513d4d9-684b-4c3c-add4-6404969006de",
				        "name": "postgres-db-f1-micro"
				    },
				    {
				        "service_properties": { "tier": "db-g1-small", "max_disk_size": "3062" },
				        "description": "PostgreSQL on a db-g1-small (Shared CPUs, 1.7 GB/RAM, 3062 GB/disk, 1,000 Connections)",
				        "id": "6c1174d8-243c-44d1-b7a8-e94a779f67f5",
				        "name": "postgres-db-g1-small"
				    },
				    {
				        "service_properties": { "tier": "db-n1-standard-1", "max_disk_size": "10230" },
				        "description": "PostgreSQL on a db-n1-standard-1 (1 CPUs, 3.75 GB/RAM, 10230 GB/disk, 4,000 Connections)",
				        "id": "c4e68ab5-34ca-4d02-857d-3e6b3ab079a7",
				        "name": "postgres-db-n1-standard-1"
				    },
				    {
				        "service_properties": { "tier": "db-n1-standard-2", "max_disk_size": "10230" },
				        "description": "PostgreSQL on a db-n1-standard-2 (2 CPUs, 7.5 GB/RAM, 10230 GB/disk, 4,000 Connections)",
				        "id": "3f578ecf-885c-4b60-b38b-60272f34e00f",
				        "name": "postgres-db-n1-standard-2"
				    },
				    {
				        "service_properties": { "tier": "db-n1-standard-4", "max_disk_size": "10230" },
				        "description": "PostgreSQL on a db-n1-standard-4 (4 CPUs, 15 GB/RAM, 10230 GB/disk, 4,000 Connections)",
				        "id": "b7fcab5d-d66d-4e82-af16-565e84cef7f9",
				        "name": "postgres-db-n1-standard-4"
				    },
				    {
				        "service_properties": { "tier": "db-n1-standard-8", "max_disk_size": "10230" },
				        "description": "PostgreSQL on a db-n1-standard-8 (8 CPUs, 30 GB/RAM, 10230 GB/disk, 4,000 Connections)",
				        "id": "4b2fa14a-caf1-42e0-bd8c-3342502008a8",
				        "name": "postgres-db-n1-standard-8"
				    },
				    {
				        "service_properties": { "tier": "db-n1-standard-16", "max_disk_size": "10230" },
				        "description": "PostgreSQL on a db-n1-standard-16 (16 CPUs, 60 GB/RAM, 10230 GB/disk, 4,000 Connections)",
				        "id": "ca2e770f-bfa5-4fb7-a249-8b943c3474ca",
				        "name": "postgres-db-n1-standard-16"
				    },
				    {
				        "service_properties": { "tier": "db-n1-standard-32", "max_disk_size": "10230" },
				        "description": "PostgreSQL on a db-n1-standard-32 (32 CPUs, 120 GB/RAM, 10230 GB/disk, 4,000 Connections)",
				        "id": "b44f8294-b003-4a50-80c2-706858073f44",
				        "name": "postgres-db-n1-standard-32"
				    },
				    {
				        "service_properties": { "tier": "db-n1-standard-64", "max_disk_size": "10230" },
				        "description": "PostgreSQL on a db-n1-standard-64 (64 CPUs, 240 GB/RAM, 10230 GB/disk, 4,000 Connections)",
				        "id": "d97326e0-5af2-4da5-b970-b4772d59cded",
				        "name": "postgres-db-n1-standard-64"
				    },
				    {
				        "service_properties": {"tier": "db-n1-highmem-2", "max_disk_size": "10230" },
				        "description": "PostgreSQL on a db-n1-highmem-2 (2 CPUs, 13 GB/RAM, 10230 GB/disk, 4,000 Connections)",
				        "id": "c10f8691-02f5-44eb-989f-7217393012ca",
				        "name": "postgres-db-n1-highmem-2"
				    },
				    {
				        "service_properties": { "tier": "db-n1-highmem-4", "max_disk_size": "10230" },
				        "description": "PostgreSQL on a db-n1-highmem-4 (4 CPUs, 26 GB/RAM, 10230 GB/disk, 4,000 Connections)",
				        "id": "610cc78d-d26a-41a9-90b7-547a44517f03",
				        "name": "postgres-db-n1-highmem-4"
				    },
				    {
				        "service_properties": { "tier": "db-n1-highmem-8",  "max_disk_size": "10230" },
				        "description": "PostgreSQL on a db-n1-highmem-8 (8 CPUs, 52 GB/RAM, 10230 GB/disk, 4,000 Connections)",
				        "id": "2a351e8d-958d-4c4f-ae46-c984fec18740",
				        "name": "postgres-db-n1-highmem-8"
				    },
				    {
				        "service_properties": { "tier": "db-n1-highmem-16", "max_disk_size": "10230" },
				        "description": "PostgreSQL on a db-n1-highmem-16 (16 CPUs, 104 GB/RAM, 10230 GB/disk, 4,000 Connections)",
				        "id": "51d3ca0c-9d21-447d-a395-3e0dc0659775",
				        "name": "postgres-db-n1-highmem-16"
				    },
				    {
				        "service_properties": { "tier": "db-n1-highmem-32", "max_disk_size": "10230" },
				        "description": "PostgreSQL on a db-n1-highmem-32 (32 CPUs, 208 GB/RAM, 10230 GB/disk, 4,000 Connections)",
				        "id": "2e72b386-f7ce-4f0d-a149-9f9a851337d4",
				        "name": "postgres-db-n1-highmem-32"
				    },
				    {
				        "service_properties": { "tier": "db-n1-highmem-64", "max_disk_size": "10230" },
				        "description": "PostgreSQL on a db-n1-highmem-64 (64 CPUs, 416 GB/RAM, 10230 GB/disk, 4,000 Connections)",
				        "id": "82602649-e4ac-4a2f-b80d-dacd745aed6a",
				        "name": "postgres-db-n1-highmem-64"
				    }
				]
    }`,
		ProvisionInputVariables: append([]broker.BrokerVariable{
			{
				FieldName: "instance_name",
				Type:      broker.JsonTypeString,
				Details:   "Name of the CloudSQL instance.",
				Default:   "pcf-sb-${counter.next()}-${time.nano()}",
				Constraints: validation.NewConstraintBuilder().
					Pattern("^[a-z][a-z0-9-]+$").
					MaxLength(75).
					Build(),
			},
			{
				FieldName: "database_name",
				Type:      broker.JsonTypeString,
				Details:   "Name of the database inside of the instance. Must be a valid identifier for your chosen database type.",
				Default:   "pcf-sb-${counter.next()}-${time.nano()}",
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
					Pattern("^[a-z][a-z0-9-]+$").
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
			// legacy behavior dictates that empty values get defaults
			{Name: "instance_name", Default: `${instance_name == "" ? "pcf-sb-${counter.next()}-${time.nano()}" : instance_name}`, Overwrite: true},
			{Name: "database_name", Default: `${database_name == "" ? "pcf-sb-${counter.next()}-${time.nano()}" : database_name}`, Overwrite: true},

			{Name: "is_first_gen", Default: `false`, Overwrite: true},

			// validation
			{Name: "_", Default: `${assert(disk_size <= max_disk_size, "disk size (${disk_size}) is greater than max allowed disk size for this plan (${max_disk_size})")}`, Overwrite: true},
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
	}
}
