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
	accountmanagers "github.com/GoogleCloudPlatform/gcp-service-broker/pkg/providers/builtin/account_managers"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/validation"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/varcontext"
)

const (
	passwordTemplate   = "${rand.base64(32)}"
	usernameTemplate   = `sb${str.truncate(14, time.nano())}`
	identifierTemplate = `sb-${counter.next()}-${time.nano()}`
)

func roleWhitelist() []string {
	return []string{
		"cloudsql.editor",
		"cloudsql.viewer",
		"cloudsql.client",
	}
}

func commonBindVariables() []broker.BrokerVariable {
	return append(accountmanagers.ServiceAccountWhitelistWithDefault(roleWhitelist(), "cloudsql.client"),
		broker.BrokerVariable{
			FieldName: "jdbc_uri_format",
			Type:      broker.JsonTypeString,
			Details:   "If `true`, `uri` field will contain a JDBC formatted URI.",
			Default:   "false",
			Enum: map[interface{}]string{
				"true":  "return a JDBC formatted URI",
				"false": "return a SQL formatted URI",
			},
		},
		broker.BrokerVariable{
			FieldName: "username",
			Type:      broker.JsonTypeString,
			Details:   "The SQL username for the account.",
			Default:   `sb${str.truncate(14, time.nano())}`,
		},
		broker.BrokerVariable{
			FieldName: "password",
			Type:      broker.JsonTypeString,
			Details:   "The SQL password for the account.",
			Default:   "${rand.base64(32)}",
		},
	)
}

func commonBindComputedVariables() []varcontext.DefaultVariable {
	serviceAccountComputed := accountmanagers.ServiceAccountBindComputedVariables()
	sqlComputed := []varcontext.DefaultVariable{
		// legacy behavior dictates that empty values get defaults
		{Name: "password", Default: `${password == "" ? "` + passwordTemplate + `" : password}`, Overwrite: true},
		{Name: "username", Default: `${username == "" ? "` + usernameTemplate + `" : username}`, Overwrite: true},

		// necessary additions
		{Name: "certname", Default: `${str.truncate(10, request.binding_id)}cert`, Overwrite: true},
		{Name: "db_name", Default: `${instance.name}`, Overwrite: true},
	}

	return append(serviceAccountComputed, sqlComputed...)
}

func commonProvisionVariables() []broker.BrokerVariable {
	return []broker.BrokerVariable{
		{
			FieldName: "region",
			Type:      broker.JsonTypeString,
			Details:   "The geographical region. See the instance locations list https://cloud.google.com/sql/docs/mysql/instance-locations for which regions support which databases.",
			Default:   "us-central",
			Constraints: validation.NewConstraintBuilder().
				Pattern("^[A-Za-z][-a-z0-9A-Z]+$").
				Examples("northamerica-northeast1", "southamerica-east1", "us-east1").
				Build(),
		},
		{
			FieldName: "disk_size",
			Type:      broker.JsonTypeString,
			Details:   "In GB.",
			Default:   "10",
			Constraints: validation.NewConstraintBuilder().
				Pattern("^[1-9][0-9]+$").
				MaxLength(5).
				Examples("10", "500", "10230").
				Build(),
		},
		{
			FieldName: "database_flags",
			Type:      broker.JsonTypeString,
			Details:   "The database flags passed to the instance at startup (comma separated list of flags, e.g. general_log=on,skip_show_database=off).",
			Default:   "",
			Constraints: validation.NewConstraintBuilder().
				Pattern(`^(|([a-z_]+=[a-zA-Z0-9\.\+\:-]+)(,[a-z_]+=[a-zA-Z0-9\.\+\:-]+)*)$`).
				Examples("long_query_time=10", "general_log=on,skip_show_database=off").
				Build(),
		},
		{
			FieldName: "zone",
			Type:      broker.JsonTypeString,
			Details:   "Optional, the specific zone in the region to run the instance.",
			Default:   "",
			Constraints: validation.NewConstraintBuilder().
				Pattern("^(|[A-Za-z][-a-z0-9A-Z]+)$").
				Build(),
		},
		{
			FieldName: "disk_type",
			Type:      broker.JsonTypeString,
			Details:   "The type of disk backing the database.",
			Default:   "PD_SSD",
			Enum: map[interface{}]string{
				"PD_SSD": "flash storage drive",
				"PD_HDD": "magnetic hard drive",
			},
		},
		{
			FieldName: "maintenance_window_day",
			Type:      broker.JsonTypeString,
			Details:   "The day of week a CloudSQL instance should preferably be restarted for system maintenance purposes. (1-7), starting on Monday.",
			Default:   "1",
			Enum: map[interface{}]string{
				"1": "Monday",
				"2": "Tuesday",
				"3": "Wednesday",
				"4": "Thursday",
				"5": "Friday",
				"6": "Saturday",
				"7": "Sunday",
			},
		},
		{
			FieldName: "maintenance_window_hour",
			Type:      broker.JsonTypeString,
			Details:   "The hour of the day when disruptive updates (updates that require an instance restart) to this CloudSQL instance can be made. Hour of day 0-23.",
			Default:   "0",
			Constraints: validation.NewConstraintBuilder().
				Pattern("^([0-9]|1[0-9]|2[0-3])$").
				Build(),
		},
		{
			FieldName: "backups_enabled",
			Type:      broker.JsonTypeString,
			Details:   "Should daily backups be enabled for the service?",
			Default:   "true",
			Enum: map[interface{}]string{
				"true":  "enable daily backups",
				"false": "do not enable daily backups",
			},
		},
		{
			FieldName: "backup_start_time",
			Type:      broker.JsonTypeString,
			Details:   "Start time for the daily backup configuration in UTC timezone in the 24 hour format - HH:MM.",
			Default:   "06:00",
			Constraints: validation.NewConstraintBuilder().
				Pattern("^(0[0-9]|1[0-9]|2[0-3]):[0-5][0-9]$").
				Build(),
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
			Details:   "Configuration to increase storage size automatically.",
			Default:   "false",
			Enum: map[interface{}]string{
				"true":  "increase storage size automatically",
				"false": "do not increase storage size automatically",
			},
		},
		{
			FieldName: "auto_resize_limit",
			Type:      broker.JsonTypeString,
			Details:   "The maximum size to which storage capacity can be automatically increased.",
			Default:   "0",
			Constraints: validation.NewConstraintBuilder().
				Pattern("^[0-9][0-9]*$").
				MaxLength(5).
				Examples("10", "500", "10230").
				Build(),
		},
		{
			FieldName: "availability_type",
			Type:      broker.JsonTypeString,
			Details:   "Availability type specifies whether the instance serves data from multiple zones.",
			Default:   "ZONAL",
			Enum: map[interface{}]string{
				"ZONAL":    "The instance serves data from only one zone (NOT highly available).",
				"REGIONAL": "The instance serves data zones in a region (highly available).",
			},
		},
	}
}

func commonBindOutputVariables() []broker.BrokerVariable {
	return append(accountmanagers.ServiceAccountBindOutputVariables(), []broker.BrokerVariable{
		// Certificate
		{
			FieldName: "CaCert",
			Type:      broker.JsonTypeString,
			Details:   "The server Certificate Authority's certificate.",
			Required:  true,
			Constraints: validation.NewConstraintBuilder().
				Examples("-----BEGIN CERTIFICATE-----BASE64 Certificate Text-----END CERTIFICATE-----").
				Build(),
		},
		{
			FieldName: "ClientCert",
			Type:      broker.JsonTypeString,
			Details:   "The client certificate.",
			Required:  true,
			Constraints: validation.NewConstraintBuilder().
				Examples("-----BEGIN CERTIFICATE-----BASE64 Certificate Text-----END CERTIFICATE-----").
				Build(),
		},
		{
			FieldName: "ClientKey",
			Type:      broker.JsonTypeString,
			Details:   "The client certificate key.",
			Required:  true,
			Constraints: validation.NewConstraintBuilder().
				Examples("-----BEGIN RSA PRIVATE KEY-----BASE64 Key Text-----END RSA PRIVATE KEY-----").
				Build(),
		},
		{
			FieldName: "Sha1Fingerprint",
			Type:      broker.JsonTypeString,
			Details:   "The SHA1 fingerprint of the client certificate.",
			Required:  true,
			Constraints: validation.NewConstraintBuilder().
				Examples("e6d0c68f35032c6c2132217d1f1fb06b12ed32e2").
				Pattern(`^[0-9a-f]{40}$`).
				Build(),
		},

		// Connection URI
		{
			FieldName:   "UriPrefix",
			Type:        broker.JsonTypeString,
			Details:     "The connection prefix.",
			Required:    false,
			Constraints: validation.NewConstraintBuilder().Examples("jdbc:", "").Build(),
		},
		{
			FieldName:   "Username",
			Type:        broker.JsonTypeString,
			Details:     "The name of the SQL user provisioned.",
			Required:    true,
			Constraints: validation.NewConstraintBuilder().Examples("sb15404128767777").Build(),
		},
		{
			FieldName:   "Password",
			Type:        broker.JsonTypeString,
			Details:     "The database password for the SQL user.",
			Required:    true,
			Constraints: validation.NewConstraintBuilder().Examples("N-JPz7h2RHPZ81jB5gDHdnluddnIFMWG4nd5rKjR_8A=").Build(),
		},
		{
			FieldName:   "database_name",
			Type:        broker.JsonTypeString,
			Details:     "The name of the database on the instance.",
			Required:    true,
			Constraints: validation.NewConstraintBuilder().Examples("sb-2-1540412407295372465").Build(),
		},
		{
			FieldName:   "host",
			Type:        broker.JsonTypeString,
			Details:     "The hostname or IP address of the database instance.",
			Required:    true,
			Constraints: validation.NewConstraintBuilder().Examples("127.0.0.1").Build(),
		},
		{
			FieldName: "instance_name",
			Type:      broker.JsonTypeString,
			Details:   "The name of the database instance.",
			Required:  true,
			Constraints: validation.NewConstraintBuilder().
				Examples("sb-1-1540412407295273023").
				Pattern("^[a-z][a-z0-9-]+$").
				MaxLength(84).
				Build(),
		},
		{
			FieldName:   "uri",
			Type:        broker.JsonTypeString,
			Details:     "A database connection string.",
			Required:    true,
			Constraints: validation.NewConstraintBuilder().Examples("mysql://user:pass@127.0.0.1/sb-2-1540412407295372465?ssl_mode=required").Build(),
		},
		{
			FieldName:   "last_master_operation_id",
			Type:        broker.JsonTypeString,
			Details:     "(deprecated) The id of the last operation on the database.",
			Required:    false,
			Constraints: validation.NewConstraintBuilder().Examples("mysql://user:pass@127.0.0.1/sb-2-1540412407295372465?ssl_mode=required").Build(),
		},
		{
			FieldName: "region",
			Type:      broker.JsonTypeString,
			Details:   "The region the database is in.",
			Required:  true,
			Constraints: validation.NewConstraintBuilder().
				Pattern("^[A-Za-z][-a-z0-9A-Z]+$").
				Examples("northamerica-northeast1", "southamerica-east1", "us-east1").
				Build(),
		},
	}...)
}
