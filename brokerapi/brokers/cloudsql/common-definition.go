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
	accountmanagers "github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/account_managers"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/broker"
)

var roleWhitelist = []string{
	"cloudsql.editor",
	"cloudsql.viewer",
	"cloudsql.client",
}

var commonBindVariables = append(accountmanagers.ServiceAccountBindInputVariables(roleWhitelist),
	broker.BrokerVariable{
		FieldName: "jdbc_uri_format",
		Type:      broker.JsonTypeString,
		Details:   "If `true`, `uri` field will contain a JDBC formatted URI.",
		Default:   "false",
	},
	broker.BrokerVariable{
		FieldName: "username",
		Type:      broker.JsonTypeString,
		Details:   "The SQL username for the account.",
		Default:   "a generated value",
	},
	broker.BrokerVariable{
		FieldName: "password",
		Type:      broker.JsonTypeString,
		Details:   "The SQL password for the account.",
		Default:   "a generated value",
	},
)

var commonProvisionVariables = []broker.BrokerVariable{
	{
		FieldName: "instance_name",
		Type:      broker.JsonTypeString,
		Details:   "Name of the Cloud SQL instance.",
		Default:   "a generated value",
	},
	{
		FieldName: "database_name",
		Type:      broker.JsonTypeString,
		Details:   "Name of the database inside of the instance.",
		Default:   "a generated value",
	},
	{
		FieldName: "version",
		Type:      broker.JsonTypeString,
		Details:   "The database engine type and version. Defaults to `MYSQL_5_6` for 1st gen MySQL instances, `MYSQL_5_7` for 2nd gen MySQL instances, or `POSTGRES_9_6` for PostgreSQL instances.",
	},
	{
		FieldName: "disk_size",
		Type:      broker.JsonTypeString,
		Details:   "In GB (only for 2nd generation instances).",
		Default:   "10",
	},
	{
		FieldName: "region",
		Type:      broker.JsonTypeString,
		Details:   "The geographical region.",
		Default:   "us-central",
	},
	{
		FieldName: "zone",
		Type:      broker.JsonTypeString,
		Details:   "(only for 2nd generation instances)",
	},
	{
		FieldName: "disk_type",
		Type:      broker.JsonTypeString,
		Details:   "(only for 2nd generation instances)",
		Default:   "ssd",
	},
	{
		FieldName: "failover_replica_name",
		Type:      broker.JsonTypeString,
		Details:   "(only for 2nd generation instances) If specified, creates a failover replica with the given name.",
		Default:   "",
	},
	{
		FieldName: "maintenance_window_day",
		Type:      broker.JsonTypeString,
		Details:   "(only for 2nd generation instances) The day when disruptive updates (updates that require an instance restart) to this CloudSQL instance can be made. Day of week (1-7), starting on Monday.",
		Default:   "1",
	},
	{
		FieldName: "maintenance_window_hour",
		Type:      broker.JsonTypeString,
		Details:   "(only for 2nd generation instances) The hour of the day when disruptive updates (updates that require an instance restart) to this CloudSQL instance can be made. Hour of day 0-23.",
		Default:   "0",
	},
	{
		FieldName: "backups_enabled",
		Type:      broker.JsonTypeString,
		Details:   "Should daily backups be enabled for the service?",
		Default:   "true",
	},
	{
		FieldName: "backup_start_time",
		Type:      broker.JsonTypeString,
		Details:   "Start time for the daily backup configuration in UTC timezone in the 24 hour format - HH:MM.",
		Default:   "06:00",
	},
	{
		FieldName: "binlog",
		Type:      broker.JsonTypeString,
		Details:   "Whether binary log is enabled. If backup configuration is disabled, binary log must be disabled as well. Defaults: `false` for 1st gen, `true` for 2nd gen, set to `true` to use.",
	},
	{
		FieldName: "activation_policy",
		Type:      broker.JsonTypeString,
		Details:   "The activation policy specifies when the instance is activated; it is applicable only when the instance state is RUNNABLE.",
		Default:   "ON_DEMAND",
		Enum: map[interface{}]string{
			"ALWAYS":    "Always, instance is always on.",
			"NEVER":     "Never, instance does not turn on if a request arrives.",
			"ON_DEMAND": "On Demand, instance responses to incoming requests and turns off when not in use.",
		},
	},
	{
		FieldName: "authorized_networks",
		Type:      broker.JsonTypeString,
		Details:   "A comma separated list without spaces.",
		Default:   "none",
	},
	{
		FieldName: "replication_type",
		Type:      broker.JsonTypeString,
		Details:   "The type of replication this instance uses. This can be either ASYNCHRONOUS or SYNCHRONOUS.",
		Default:   "SYNCHRONOUS",
		Enum: map[interface{}]string{
			"ASYNCHRONOUS": "Asynchronous Replication",
			"SYNCHRONOUS":  "Synchronous Replication",
		},
	},
	{
		FieldName: "auto_resize",
		Type:      broker.JsonTypeString,
		Details:   "(only for 2nd generation instances) Configuration to increase storage size automatically.",
		Default:   "false",
	},
}

var commonBindOutputVariables = append(accountmanagers.ServiceAccountBindOutputVariables(), []broker.BrokerVariable{
	// Certificate
	{FieldName: "CaCert", Type: broker.JsonTypeString, Details: "The server Certificate Authority's certificate."},
	{FieldName: "ClientCert", Type: broker.JsonTypeString, Details: "The client certificate. For First Generation instances, the new certificate does not take effect until the instance is restarted."},
	{FieldName: "ClientKey", Type: broker.JsonTypeString, Details: "The client certificate key."},
	{FieldName: "Sha1Fingerprint", Type: broker.JsonTypeString, Details: "The SHA1 fingerprint of the client certificate."},

	// Connection URI
	{FieldName: "UriPrefix", Type: broker.JsonTypeString, Details: "The connection prefix e.g. `mysql` or `postgres`"},
	{FieldName: "Username", Type: broker.JsonTypeString, Details: "The name of the SQL user provisioned"},
	{FieldName: "database_name", Type: broker.JsonTypeString, Details: "The name of the database on the instance"},
	{FieldName: "host", Type: broker.JsonTypeString, Details: "The hostname or ip of the database instance"},
	{FieldName: "instance_name", Type: broker.JsonTypeString, Details: "The name of the database instance"},
	{FieldName: "uri", Type: broker.JsonTypeString, Details: "A database connection string"},

	{FieldName: "last_master_operation_id", Type: broker.JsonTypeString, Details: "(GCP internals) The id of the last operation on the database."},
	{FieldName: "region", Type: broker.JsonTypeString, Details: "The region the database is in."},
}...)
