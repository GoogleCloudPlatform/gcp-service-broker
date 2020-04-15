package cloudsql

import (
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/broker"
	"github.com/pivotal-cf/brokerapi"
)

const postgresVPCID = "c90ea118-605a-47e8-8f63-57fc09c113f1"

// PostgresVPCServiceDefinition creates a new ServiceDefinition object for the
// Postgres service on a VPC.
func PostgresVPCServiceDefinition() *broker.ServiceDefinition {
	definition := buildDatabase(cloudSQLOptions{
		DatabaseType:                 postgresDatabaseType,
		CustomizableActivationPolicy: false,
		AdminControlsTier:            false,
		AdminControlsMaxDiskSize:     false,
		VPCNetwork:                   true,
	})
	definition.Id = postgresVPCID
	definition.Name = "google-cloudsql-postgres-vpc"
	definition.Plans = []broker.ServicePlan{
		{
			ServicePlan: brokerapi.ServicePlan{
				ID:          "60f0b6c0-c48f-4f84-baab-57836611e013",
				Name:        "default",
				Description: "PostgreSQL attached to a VPC",
				Free:        brokerapi.FreeValue(false),
			},
			ServiceProperties: map[string]string{},
		},
	}

	definition.Examples = []broker.ServiceExample{
		{
			Name:        "Dedicated Machine Sandbox",
			Description: "A low end PostgreSQL sandbox that uses a dedicated machine.",
			PlanId:      "60f0b6c0-c48f-4f84-baab-57836611e013",
			ProvisionParams: map[string]interface{}{
				"tier":            "db-custom-1-3840",
				"backups_enabled": "false",
				"disk_size":       "25",
			},
			BindParams: map[string]interface{}{
				"role": "cloudsql.editor",
			},
		},
		{
			Name:        "HA Instance",
			Description: "A regionally available database with automatic failover.",
			PlanId:      "60f0b6c0-c48f-4f84-baab-57836611e013",
			ProvisionParams: map[string]interface{}{
				"tier":              "db-custom-1-3840",
				"backups_enabled":   "true",
				"disk_size":         "25",
				"availability_type": "REGIONAL",
			},
			BindParams: map[string]interface{}{
				"role": "cloudsql.editor",
			},
		},
	}

	return definition
}
