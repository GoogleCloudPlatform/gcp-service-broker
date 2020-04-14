package cloudsql

import (
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/broker"
	"github.com/pivotal-cf/brokerapi"
)

const mySQLVPCID = "b48d2a6b-b1b0-499f-8389-57ba33bfbb19"

// MySQLVPCServiceDefinition creates a new ServiceDefinition object for the MySQL service
// on a VPC.
func MySQLVPCServiceDefinition() *broker.ServiceDefinition {
	definition := buildDatabase(cloudSQLOptions{
		DatabaseType:                 mySQLDatabaseType,
		CustomizableActivationPolicy: false,
		AdminControlsTier:            false,
		AdminControlsMaxDiskSize:     false,
		VPCNetwork:                   true,
	})
	definition.Id = mySQLVPCID
	definition.Plans = []broker.ServicePlan{
		{
			ServicePlan: brokerapi.ServicePlan{
				ID:          "89e2c84e-4d5c-457c-ad14-329dcf44b806",
				Name:        "default",
				Description: "MySQL attached to a VPC",
				Free:        brokerapi.FreeValue(false),
			},
			ServiceProperties: map[string]string{},
		},
	}

	definition.Examples = []broker.ServiceExample{
		{
			Name:        "HA Instance",
			Description: "A regionally available database with automatic failover.",
			PlanId:      "89e2c84e-4d5c-457c-ad14-329dcf44b806",
			ProvisionParams: map[string]interface{}{
				"tier":              "db-n1-standard-1",
				"backups_enabled":   "true",
				"binlog":            "true",
				"availability_type": "REGIONAL",
			},
			BindParams: map[string]interface{}{
				"role": "cloudsql.editor",
			},
		},
		{
			Name:        "Development Sandbox",
			Description: "An inexpensive MySQL sandbox for developing with no backups.",
			PlanId:      "89e2c84e-4d5c-457c-ad14-329dcf44b806",
			ProvisionParams: map[string]interface{}{
				"tier":            "db-n1-standard-1",
				"backups_enabled": "false",
				"binlog":          "false",
				"disk_size":       "10",
			},
			BindParams: map[string]interface{}{
				"role": "cloudsql.editor",
			},
		},
	}

	return definition
}
