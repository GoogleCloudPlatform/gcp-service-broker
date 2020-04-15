// Copyright 2020 the Service Broker Project Authors.
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
	"fmt"
	"strings"

	"code.cloudfoundry.org/lager"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/broker"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/providers/builtin/base"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/validation"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/varcontext"
	"golang.org/x/oauth2/jwt"
)

type databaseType struct {
	Name               string
	URIFormat          string
	CustomizableBinlog bool
	DefaultVersion     string
	Versions           map[interface{}]string
	InstanceNameLength int
	Tags               []string
}

var mySQLDatabaseType = databaseType{
	Name:               "MySQL",
	URIFormat:          `${UriPrefix}mysql://${str.queryEscape(Username)}:${str.queryEscape(Password)}@${str.queryEscape(host)}/${str.queryEscape(database_name)}?ssl_mode=required`,
	CustomizableBinlog: true,
	DefaultVersion:     mySqlSecondGenDefaultVersion,
	Versions: map[interface{}]string{
		"MYSQL_5_6": "MySQL 5.6.X",
		"MYSQL_5_7": "MySQL 5.7.X",
	},
	InstanceNameLength: 84,
	Tags:               []string{"gcp", "cloudsql", "mysql"},
}

var postgresDatabaseType = databaseType{
	Name:               "PostgreSQL",
	URIFormat:          `${UriPrefix}postgres://${str.queryEscape(Username)}:${str.queryEscape(Password)}@${str.queryEscape(host)}/${str.queryEscape(database_name)}?sslmode=require&sslcert=${str.queryEscape(ClientCert)}&sslkey=${str.queryEscape(ClientKey)}&sslrootcert=${str.queryEscape(CaCert)}`,
	CustomizableBinlog: false,
	DefaultVersion:     "POSTGRES_11",
	Versions: map[interface{}]string{
		"POSTGRES_9_6": "PostgreSQL 9.6.X",
		"POSTGRES_10":  "PostgreSQL 10",
		"POSTGRES_11":  "PostgreSQL 11",
		"POSTGRES_12":  "PostgreSQL 12",
	},
	InstanceNameLength: 86,
	Tags:               []string{"gcp", "cloudsql", "postgres"},
}

type cloudSQLOptions struct {
	DatabaseType                 databaseType
	CustomizableActivationPolicy bool
	AdminControlsTier            bool
	AdminControlsMaxDiskSize     bool
	VPCNetwork                   bool
}

func buildDatabase(opts cloudSQLOptions) *broker.ServiceDefinition {

	// Initial
	name := strings.ToLower(fmt.Sprintf("google-cloudsql-%s", opts.DatabaseType.Name))
	if opts.VPCNetwork {
		name += "-vpc"
	}

	defn := &broker.ServiceDefinition{
		Name:             name,
		Description:      fmt.Sprintf("Google CloudSQL for %[1]s is a fully-managed %[1]s database service.", opts.DatabaseType.Name),
		DisplayName:      fmt.Sprintf("Google CloudSQL for %s", opts.DatabaseType.Name),
		ImageUrl:         "https://cloud.google.com/_static/images/cloud/products/logos/svg/sql.svg",
		DocumentationUrl: "https://cloud.google.com/sql/docs/",
		SupportUrl:       "https://cloud.google.com/sql/docs/getting-support/",
		Tags:             opts.DatabaseType.Tags,
		Bindable:         true,
		PlanUpdateable:   false,

		DefaultRoleWhitelist:  roleWhitelist(),
		BindInputVariables:    commonBindVariables(),
		BindOutputVariables:   commonBindOutputVariables(),
		BindComputedVariables: commonBindComputedVariables(),

		ProviderBuilder: func(projectId string, auth *jwt.Config, logger lager.Logger) broker.ServiceProvider {
			bb := base.NewBrokerBase(projectId, auth, logger)
			return &CloudSQLBroker{
				BrokerBase: bb,
				uriFormat:  opts.DatabaseType.URIFormat,
			}
		},
		IsBuiltin: true,
	}

	// Database type specific stuff
	defn.ProvisionInputVariables = []broker.BrokerVariable{
		{
			FieldName: "instance_name",
			Type:      broker.JsonTypeString,
			Details:   "Name of the CloudSQL instance.",
			Default:   identifierTemplate,
			Constraints: validation.NewConstraintBuilder().
				Pattern("^[a-z][a-z0-9-]+$").
				MaxLength(opts.DatabaseType.InstanceNameLength).
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
			Details:   "The database engine type and version.",
			Default:   opts.DatabaseType.DefaultVersion,
			Enum:      opts.DatabaseType.Versions,
		},
	}

	defn.ProvisionComputedVariables = []varcontext.DefaultVariable{
		{Name: "labels", Default: `${json.marshal(request.default_labels)}`, Overwrite: true},

		// legacy behavior dictates that empty values get defaults
		{Name: "instance_name", Default: `${instance_name == "" ? "` + identifierTemplate + `" : instance_name}`, Overwrite: true},
		{Name: "database_name", Default: `${database_name == "" ? "` + identifierTemplate + `" : database_name}`, Overwrite: true},
	}

	if opts.CustomizableActivationPolicy {
		defn.ProvisionInputVariables = append(defn.ProvisionInputVariables, broker.BrokerVariable{
			FieldName: "activation_policy",
			Type:      broker.JsonTypeString,
			Details:   "The activation policy specifies when the instance is activated; it is applicable only when the instance state is RUNNABLE.",
			Default:   "ALWAYS",
			Enum: map[interface{}]string{
				"ALWAYS": "Always, instance is always on.",
				"NEVER":  "Never, instance does not turn on if a request arrives.",
			},
		})
	} else {
		defn.ProvisionComputedVariables = append(defn.ProvisionComputedVariables, varcontext.DefaultVariable{
			Name:      "activation_policy",
			Default:   `ALWAYS`,
			Overwrite: true,
		})
	}

	if opts.DatabaseType.CustomizableBinlog {
		defn.ProvisionInputVariables = append(defn.ProvisionInputVariables, broker.BrokerVariable{
			FieldName: "binlog",
			Type:      broker.JsonTypeString,
			Details:   "Whether binary log is enabled. Must be enabled for high availability.",
			Default:   "true",
			Enum: map[interface{}]string{
				"true":  "use binary log",
				"false": "do not use binary log",
			},
		})
	} else {
		defn.ProvisionComputedVariables = append(defn.ProvisionComputedVariables, varcontext.DefaultVariable{
			Name:      "binlog",
			Default:   `false`,
			Overwrite: true,
		})
	}

	if opts.AdminControlsTier {
		defn.PlanVariables = append(defn.PlanVariables, broker.BrokerVariable{
			FieldName: "tier",
			Type:      broker.JsonTypeString,
			Details:   "The machine type the database will run on. MySQL has predefined tiers, other databases use the a string of the form db-custom-[CPUS]-[MEMORY_MBS], where memory is at least 3840.",
			Required:  true,
		})
	} else {
		defn.ProvisionInputVariables = append(defn.ProvisionInputVariables, broker.BrokerVariable{
			FieldName: "tier",
			Type:      broker.JsonTypeString,
			Details:   "The machine type the database will run on. MySQL has predefined tiers, other databases use the a string of the form db-custom-[CPUS]-[MEMORY_MBS], where memory is at least 3840.",
			Constraints: validation.NewConstraintBuilder().
				Pattern("^[A-Za-z][-a-z0-9A-Z]+$").
				Examples("db-n1-standard-1", "db-custom-1-3840").
				Build(),
		})
	}

	if opts.AdminControlsMaxDiskSize {
		defn.PlanVariables = append(defn.PlanVariables, broker.BrokerVariable{
			FieldName: "max_disk_size",
			Type:      broker.JsonTypeString,
			Details:   "Maximum disk size in GB, 10 is the minimum.",
			Default:   "10",
			Required:  true,
		})

		defn.ProvisionComputedVariables = append(defn.ProvisionComputedVariables, varcontext.DefaultVariable{
			Name:      "_",
			Default:   `${assert(disk_size <= max_disk_size, "disk size (${disk_size}) is greater than max allowed disk size for this plan (${max_disk_size})")}`,
			Overwrite: true,
		})
	} else {
		// no-op, max_disk_size is an artificial constraint
	}

	if opts.VPCNetwork {
		defn.ProvisionInputVariables = append(defn.ProvisionInputVariables, broker.BrokerVariable{
			FieldName: "private_network",
			Type:      broker.JsonTypeString,
			Details:   "The private network to attach to. If specified the instance will only be accessible on the VPC.",
			Default:   "default",
			Constraints: validation.NewConstraintBuilder().
				Examples("projects/my-project/global/networks/default").
				Build(),
		})

		defn.ProvisionComputedVariables = append(defn.ProvisionComputedVariables, varcontext.DefaultVariable{
			Name:      "authorized_networks",
			Default:   ``,
			Overwrite: true,
		})
	} else {
		defn.ProvisionInputVariables = append(
			defn.ProvisionInputVariables,
			broker.BrokerVariable{
				FieldName: "authorized_networks",
				Type:      broker.JsonTypeString,
				Details:   "A comma separated list without spaces.",
				Default:   "",
			},
		)

		defn.ProvisionComputedVariables = append(defn.ProvisionComputedVariables, varcontext.DefaultVariable{
			Name:      "private_network",
			Default:   ``,
			Overwrite: true,
		})
	}

	defn.ProvisionInputVariables = append(defn.ProvisionInputVariables, commonProvisionVariables()...)
	return defn
}
