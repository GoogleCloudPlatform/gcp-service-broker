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
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/models"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/broker"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/validation"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/varcontext"
)

func init() {
	broker.Register(mysqlServiceDefinition())
}

func mysqlServiceDefinition() *broker.ServiceDefinition {
	return &broker.ServiceDefinition{
		Name: models.CloudsqlMySQLName,
		DefaultServiceDefinition: `{
		    "id": "4bc59b9a-8520-409f-85da-1c7552315863",
		    "description": "Google Cloud SQL is a fully-managed MySQL database service.",
		    "name": "google-cloudsql-mysql",
		    "bindable": true,
		    "plan_updateable": false,
		    "metadata": {
		      "displayName": "Google CloudSQL MySQL",
		      "longDescription": "Google Cloud SQL is a fully-managed MySQL database service.",
		      "documentationUrl": "https://cloud.google.com/sql/docs/",
		      "supportUrl": "https://cloud.google.com/support/",
		      "imageUrl": "https://cloud.google.com/_static/images/cloud/products/logos/svg/sql.svg"
		    },
		    "tags": ["gcp", "cloudsql", "mysql"],
		    "plans": [
				    {
				        "service_properties": {
				            "tier": "db-f1-micro",
				            "max_disk_size": "3062"
				        },
				        "description": "MySQL on a db-f1-micro (Shared CPUs, 0.6 GB/RAM, 3062 GB/disk, 250 Connections)",
				        "id": "7d8f9ade-30c1-4c96-b622-ea0205cc5f0b",
				        "name": "mysql-db-f1-micro"
				    },
				    {
				        "service_properties": {
				            "tier": "db-g1-small",
				            "max_disk_size": "3062"
				        },
				        "description": "MySQL on a db-g1-small (Shared CPUs, 1.7 GB/RAM, 3062 GB/disk, 1,000 Connections)",
				        "id": "b68bf4d8-1636-4121-af2f-087e46189929",
				        "name": "mysql-db-g1-small"
				    },
				    {
				        "service_properties": {
				            "tier": "db-n1-standard-1",
				            "max_disk_size": "10230"
				        },
				        "description": "MySQL on a db-n1-standard-1 (1 CPUs, 3.75 GB/RAM, 10230 GB/disk, 4,000 Connections)",
				        "id": "bdfd8033-c2b9-46e9-9b37-1f3a5889eef4",
				        "name": "mysql-db-n1-standard-1"
				    },
				    {
				        "service_properties": {
				            "tier": "db-n1-standard-2",
				            "max_disk_size": "10230"
				        },
				        "description": "MySQL on a db-n1-standard-2 (2 CPUs, 7.5 GB/RAM, 10230 GB/disk, 4,000 Connections)",
				        "id": "2c99e938-4c1e-4da7-810a-94c9f5b71b57",
				        "name": "mysql-db-n1-standard-2"
				    },
				    {
				        "service_properties": {
				            "tier": "db-n1-standard-4",
				            "max_disk_size": "10230"
				        },
				        "description": "MySQL on a db-n1-standard-4 (4 CPUs, 15 GB/RAM, 10230 GB/disk, 4,000 Connections)",
				        "id": "d520a5f5-7485-4a83-849b-5439f911fe26",
				        "name": "mysql-db-n1-standard-4"
				    },
				    {
				        "service_properties": {
				            "tier": "db-n1-standard-8",
				            "max_disk_size": "10230"
				        },
				        "description": "MySQL on a db-n1-standard-8 (8 CPUs, 30 GB/RAM, 10230 GB/disk, 4,000 Connections)",
				        "id": "7ef42bb4-87e3-4ead-8118-4e88c98ed2e6",
				        "name": "mysql-db-n1-standard-8"
				    },
				    {
				        "service_properties": {
				            "tier": "db-n1-standard-16",
				            "max_disk_size": "10230"
				        },
				        "description": "MySQL on a db-n1-standard-16 (16 CPUs, 60 GB/RAM, 10230 GB/disk, 4,000 Connections)",
				        "id": "200bd90a-4323-46d8-8aa5-afd4601498d0",
				        "name": "mysql-db-n1-standard-16"
				    },
				    {
				        "service_properties": {
				            "tier": "db-n1-standard-32",
				            "max_disk_size": "10230"
				        },
				        "description": "MySQL on a db-n1-standard-32 (32 CPUs, 120 GB/RAM, 10230 GB/disk, 4,000 Connections)",
				        "id": "52305df2-1e64-4cdb-a4c9-bb5dddb33c3e",
				        "name": "mysql-db-n1-standard-32"
				    },
				    {
				        "service_properties": {
				            "tier": "db-n1-standard-64",
				            "max_disk_size": "10230"
				        },
				        "description": "MySQL on a db-n1-standard-64 (64 CPUs, 240 GB/RAM, 10230 GB/disk, 4,000 Connections)",
				        "id": "e45d7c44-4990-4dac-a14d-c5127e9ae0c5",
				        "name": "mysql-db-n1-standard-64"
				    },
				    {
				        "service_properties": {
				            "tier": "db-n1-highmem-2",
				            "max_disk_size": "10230"
				        },
				        "description": "MySQL on a db-n1-highmem-2 (2 CPUs, 13 GB/RAM, 10230 GB/disk, 4,000 Connections)",
				        "id": "07b8a04c-0efe-42d3-8b2c-2c23f7c79583",
				        "name": "mysql-db-n1-highmem-2"
				    },
				    {
				        "service_properties": {
				            "tier": "db-n1-highmem-4",
				            "max_disk_size": "10230"
				        },
				        "description": "MySQL on a db-n1-highmem-4 (4 CPUs, 26 GB/RAM, 10230 GB/disk, 4,000 Connections)",
				        "id": "50fa4baa-e36f-41c3-bbe9-c986d9fbe3c8",
				        "name": "mysql-db-n1-highmem-4"
				    },
				    {
				        "service_properties": {
				            "tier": "db-n1-highmem-8",
				            "max_disk_size": "10230"
				        },
				        "description": "MySQL on a db-n1-highmem-8 (8 CPUs, 52 GB/RAM, 10230 GB/disk, 4,000 Connections)",
				        "id": "6e8e5bc3-bf68-4e57-bda1-d9c9a67faee0",
				        "name": "mysql-db-n1-highmem-8"
				    },
				    {
				        "service_properties": {
				            "tier": "db-n1-highmem-16",
				            "max_disk_size": "10230"
				        },
				        "description": "MySQL on a db-n1-highmem-16 (16 CPUs, 104 GB/RAM, 10230 GB/disk, 4,000 Connections)",
				        "id": "3c83ff6b-165e-47bf-9bba-f4801390d0ff",
				        "name": "mysql-db-n1-highmem-16"
				    },
				    {
				        "service_properties": {
				            "tier": "db-n1-highmem-32",
				            "max_disk_size": "10230"
				        },
				        "description": "MySQL on a db-n1-highmem-32 (32 CPUs, 208 GB/RAM, 10230 GB/disk, 4,000 Connections)",
				        "id": "cbc6d376-8fd3-4a34-9ab5-324311f038f6",
				        "name": "mysql-db-n1-highmem-32"
				    },
				    {
				        "service_properties": {
				            "tier": "db-n1-highmem-64",
				            "max_disk_size": "10230"
				        },
				        "description": "MySQL on a db-n1-highmem-64 (64 CPUs, 416 GB/RAM, 10230 GB/disk, 4,000 Connections)",
				        "id": "b0742cc5-caba-4b8d-98e0-03380ae9522b",
				        "name": "mysql-db-n1-highmem-64"
				    }
				]
		}`,
		ProvisionInputVariables: append([]broker.BrokerVariable{
			{
				FieldName: "instance_name",
				Type:      broker.JsonTypeString,
				Details:   "Name of the Cloud SQL instance.",
				Default:   identifierTemplate,
				Constraints: validation.NewConstraintBuilder().
					Pattern("^[a-z][a-z0-9-]+$").
					MaxLength(84).
					Build(),
			},
			{
				FieldName: "database_name",
				Type:      broker.JsonTypeString,
				Details:   "Name of the database inside of the instance. Must be a valid identifier for your chosen database type.",
				Default:   identifierTemplate,
			},
			{
				FieldName: "version",
				Type:      broker.JsonTypeString,
				Details:   "The database engine type and version. Defaults to `MYSQL_5_6` for 1st gen MySQL instances or `MYSQL_5_7` for 2nd gen MySQL instances.",
				Enum: map[interface{}]string{
					"MYSQL_5_5": "MySQL 5.5.X",
					"MYSQL_5_6": "MySQL 5.6.X",
					"MYSQL_5_7": "MySQL 5.7.X",
				},
			},
			{
				FieldName: "failover_replica_name",
				Type:      broker.JsonTypeString,
				Details:   "(only for 2nd generation instances) If specified, creates a failover replica with the given name.",
				Default:   "",
				Constraints: validation.NewConstraintBuilder().
					Pattern("^[a-z][a-z0-9-]+$").
					MaxLength(84).
					Build(),
			},
			{
				FieldName: "activation_policy",
				Type:      broker.JsonTypeString,
				Details:   "The activation policy specifies when the instance is activated; it is applicable only when the instance state is RUNNABLE.",
				Default:   "ALWAYS",
				Enum: map[interface{}]string{
					"ALWAYS":    "Always, instance is always on.",
					"NEVER":     "Never, instance does not turn on if a request arrives.",
					"ON_DEMAND": "On Demand, instance responds to incoming requests and turns off when not in use.",
				},
			},
		}, commonProvisionVariables()...),
		ProvisionComputedVariables: []varcontext.DefaultVariable{
			{Name: "labels", Default: `${json.marshal(request.default_labels)}`, Overwrite: true},

			// legacy behavior dictates that empty values get defaults
			{Name: "instance_name", Default: `${instance_name == "" ? "` + identifierTemplate + `" : instance_name}`, Overwrite: true},
			{Name: "database_name", Default: `${database_name == "" ? "` + identifierTemplate + `" : database_name}`, Overwrite: true},

			{Name: "is_first_gen", Default: `${regexp.matches("^(d|D)[0-9]+$", tier)}`, Overwrite: true},
			{Name: "version", Default: `${is_first_gen ? "MYSQL_5_6" : "MYSQL_5_7"}`, Overwrite: false},
			{Name: "binlog", Default: `${is_first_gen ? false : true}`, Overwrite: false},

			// validation
			{Name: "_", Default: `${assert(disk_size <= max_disk_size, "disk size (${disk_size}) is greater than max allowed disk size for this plan (${max_disk_size})")}`, Overwrite: true},
		},
		DefaultRoleWhitelist:  roleWhitelist(),
		BindInputVariables:    commonBindVariables(models.CloudsqlMySQLName),
		BindOutputVariables:   commonBindOutputVariables(),
		BindComputedVariables: commonBindComputedVariables(),
		PlanVariables: []broker.BrokerVariable{
			{
				FieldName: "tier",
				Type:      broker.JsonTypeString,
				Details:   "Case-sensitive tier/machine type name (see https://cloud.google.com/sql/pricing for more information).",
				Required:  true,
			},
			{
				FieldName: "pricing_plan",
				Type:      broker.JsonTypeString,
				Details:   "Select a pricing plan (only for 1st generation instances).",
				Default:   "PER_USE",
				Enum: map[interface{}]string{
					"PER_USE": "Per-Use",
					"PACKAGE": "Package",
				},
				Required: true,
			},
			{
				FieldName: "max_disk_size",
				Type:      broker.JsonTypeString,
				Details:   "Maximum disk size in GB (applicable only to Second Generation instances, 10 minimum/default).",
				Default:   "10",
				Required:  true,
			},
		},
		Examples: []broker.ServiceExample{
			{
				Name:        "Development Sandbox",
				Description: "An inexpensive MySQL sandbox for developing with no backups.",
				PlanId:      "7d8f9ade-30c1-4c96-b622-ea0205cc5f0b",
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
