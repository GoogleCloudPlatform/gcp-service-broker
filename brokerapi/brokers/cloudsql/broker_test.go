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
	"encoding/json"
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/broker"
	"github.com/pivotal-cf/brokerapi"
	"github.com/spf13/viper"
	googlecloudsql "google.golang.org/api/sqladmin/v1beta4"
)

func TestCreateProvisionRequest(t *testing.T) {

	viper.Set("service.google-cloudsql-mysql.plans", `[{
      "tier": "db-n1-standard-1",
      "max_disk_size": "512",
      "id": "00000000-0000-0000-0000-000000000001",
      "name": "second-gen",
      "pricing_plan": "PACKAGE"
  },{
      "tier": "D16",
      "max_disk_size": "512",
      "id": "00000000-0000-0000-0000-000000000002",
      "name": "first-gen",
      "pricing_plan": "PACKAGE"
  }]`)

	viper.Set("service.google-cloudsql-postgres.plans", `[{
      "tier": "db-n1-standard-1",
      "max_disk_size": "512",
      "id": "00000000-0000-0000-0000-000000000003",
      "name": "second-gen",
      "pricing_plan": "PACKAGE"
  }]`)

	mysqlSecondGenPlan := "00000000-0000-0000-0000-000000000001"
	mysqlFirstgenPlan := "00000000-0000-0000-0000-000000000002"
	postgresPlan := "c4e68ab5-34ca-4d02-857d-3e6b3ab079a7"

	cases := map[string]struct {
		Service     *broker.ServiceDefinition
		PlanId      string
		UserParams  string
		Validate    func(t *testing.T, di googlecloudsql.DatabaseInstance, ii InstanceInformation)
		ErrContains string
	}{
		"blank instance names get generated": {
			Service:    mysqlServiceDefinition(),
			PlanId:     mysqlFirstgenPlan,
			UserParams: `{"instance_name":""}`,
			Validate: func(t *testing.T, di googlecloudsql.DatabaseInstance, ii InstanceInformation) {
				if len(di.Name) == 0 {
					t.Errorf("instance name wasn't generated")
				}
			},
		},

		"tiers matching (D|d)\\d+ get firstgen outputs for MySQL": {
			Service:    mysqlServiceDefinition(),
			PlanId:     mysqlFirstgenPlan,
			UserParams: `{"name":""}`,
			Validate: func(t *testing.T, di googlecloudsql.DatabaseInstance, ii InstanceInformation) {
				if di.Settings.LocationPreference != nil {
					t.Error("second-gen instance created for first-gen plan")
				}
			},
		},

		"first-gen MySQL defaults": {
			Service:    mysqlServiceDefinition(),
			PlanId:     mysqlFirstgenPlan,
			UserParams: `{}`,
			Validate: func(t *testing.T, di googlecloudsql.DatabaseInstance, ii InstanceInformation) {
				if di.DatabaseVersion != mySqlFirstGenDefaultVersion {
					t.Errorf("expected version to default to %s for first gen plan got %s", mySqlFirstGenDefaultVersion, di.DatabaseVersion)
				}

				if di.Settings.BackupConfiguration.BinaryLogEnabled == true {
					t.Error("Expected binlog to be off for first gen MySQL")
				}

				if len(di.Name) == 0 {
					t.Error("instance name wasn't generated")
				}

				if di.Settings.MaintenanceWindow != nil {
					t.Error("Expected no maintenance window by default")
				}
			},
		},
		"second-gen MySQL defaults": {
			Service:    mysqlServiceDefinition(),
			PlanId:     mysqlSecondGenPlan,
			UserParams: `{}`,
			Validate: func(t *testing.T, di googlecloudsql.DatabaseInstance, ii InstanceInformation) {
				if di.DatabaseVersion != mySqlSecondGenDefaultVersion {
					t.Errorf("expected version to default to %s for first gen plan got %s", mySqlSecondGenDefaultVersion, di.DatabaseVersion)
				}

				if di.Settings.BackupConfiguration.BinaryLogEnabled == false {
					t.Error("Expected binlog to be on by default for second-gen plans")
				}

				if len(di.Name) == 0 {
					t.Error("instance name wasn't generated")
				}

				if di.Settings.MaintenanceWindow == nil {
					t.Error("Expected maintenance window by default")
				}
			},
		},
		"PostgreSQL defaults": {
			Service:    postgresServiceDefinition(),
			PlanId:     postgresPlan,
			UserParams: `{}`,
			Validate: func(t *testing.T, di googlecloudsql.DatabaseInstance, ii InstanceInformation) {
				if di.DatabaseVersion != postgresDefaultVersion {
					t.Errorf("expected version to default to %s for first gen plan got %s", postgresDefaultVersion, di.DatabaseVersion)
				}

				if di.Settings.BackupConfiguration.BinaryLogEnabled == true {
					t.Error("Expected binlog to be off for postgres")
				}

				if len(di.Name) == 0 {
					t.Error("instance name wasn't generated")
				}

				if di.Settings.MaintenanceWindow == nil {
					t.Error("Expected maintenance window by default")
				}
			},
		},

		"partial maintenance window day": {
			Service:    mysqlServiceDefinition(),
			PlanId:     mysqlSecondGenPlan,
			UserParams: `{"maintenance_window_day":"4"}`,
			Validate: func(t *testing.T, di googlecloudsql.DatabaseInstance, ii InstanceInformation) {
				if di.Settings.MaintenanceWindow == nil {
					t.Error("Expected maintenance window on partial fill")
				}

				if di.Settings.MaintenanceWindow.Day != 4 {
					t.Errorf("Expected maintenance window day to be 4, got %v", di.Settings.MaintenanceWindow.Day)
				}
			},
		},

		"partial maintenance window hour": {
			Service:    mysqlServiceDefinition(),
			PlanId:     mysqlSecondGenPlan,
			UserParams: `{"maintenance_window_hour":"23"}`,
			Validate: func(t *testing.T, di googlecloudsql.DatabaseInstance, ii InstanceInformation) {
				if di.Settings.MaintenanceWindow == nil {
					t.Error("Expected maintenance window on partial fill")
				}

				if di.Settings.MaintenanceWindow.Hour != 23 {
					t.Errorf("Expected maintenance window day to be 4, got %v", di.Settings.MaintenanceWindow.Hour)
				}
			},
		},

		"full maintenance window ": {
			Service:    mysqlServiceDefinition(),
			PlanId:     mysqlSecondGenPlan,
			UserParams: `{"maintenance_window_day":"4","maintenance_window_hour":"23"}`,
			Validate: func(t *testing.T, di googlecloudsql.DatabaseInstance, ii InstanceInformation) {
				if di.Settings.MaintenanceWindow == nil {
					t.Error("Expected maintenance window")
				}

				if di.Settings.MaintenanceWindow.Day != 4 {
					t.Errorf("Expected maintenance window day to be 4, got %v", di.Settings.MaintenanceWindow.Day)
				}

				if di.Settings.MaintenanceWindow.Hour != 23 {
					t.Errorf("Expected maintenance window day to be 4, got %v", di.Settings.MaintenanceWindow.Hour)
				}
			},
		},

		"instance info generates db on blank ": {
			Service:    mysqlServiceDefinition(),
			PlanId:     mysqlSecondGenPlan,
			UserParams: `{"database_name":""}`,
			Validate: func(t *testing.T, di googlecloudsql.DatabaseInstance, ii InstanceInformation) {
				if len(ii.DatabaseName) == 0 {
					t.Error("Expected DatabaseName to not be blank.")
				}
			},
		},

		"instance info has name and db name ": {
			Service:    mysqlServiceDefinition(),
			PlanId:     mysqlSecondGenPlan,
			UserParams: `{"database_name":"foo", "instance_name": "bar"}`,
			Validate: func(t *testing.T, di googlecloudsql.DatabaseInstance, ii InstanceInformation) {
				if ii.DatabaseName != "foo" {
					t.Errorf("Expected DatabaseName to be foo got %s.", ii.DatabaseName)
				}

				if ii.InstanceName != "bar" {
					t.Errorf("Expected InstanceName to be bar got %s.", ii.InstanceName)
				}
			},
		},

		"mysql disk size greater than operator specified max fails": {
			Service:     mysqlServiceDefinition(),
			PlanId:      mysqlSecondGenPlan,
			UserParams:  `{"disk_size":"99999"}`,
			Validate:    func(t *testing.T, di googlecloudsql.DatabaseInstance, ii InstanceInformation) {},
			ErrContains: "disk size",
		},

		"postgres disk size greater than operator specified max fails": {
			Service:     postgresServiceDefinition(),
			PlanId:      postgresPlan,
			UserParams:  `{"disk_size":"99999"}`,
			Validate:    func(t *testing.T, di googlecloudsql.DatabaseInstance, ii InstanceInformation) {},
			ErrContains: "disk size",
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			serviceCatalogEntry, err := tc.Service.CatalogEntry()
			if err != nil {
				t.Fatalf("got error trying to get service catalog %v", err)
			}

			details := brokerapi.ProvisionDetails{RawParameters: json.RawMessage(tc.UserParams), ServiceID: serviceCatalogEntry.ID}
			plan, err := tc.Service.GetPlanById(tc.PlanId)
			if err != nil {
				t.Fatalf("got error trying to find plan %s %v", tc.PlanId, err)
			}
			if plan == nil {
				t.Fatalf("Expected plan with id %s to not be nil", tc.PlanId)
			}
			request, instanceInfo, err := createProvisionRequest("instance-id-here", details, *plan)
			if err != nil {
				if tc.ErrContains != "" && strings.Contains(err.Error(), tc.ErrContains) {
					return
				}

				t.Fatalf("got unexpected error while creating provision request: %v", err)
			}
			if tc.ErrContains != "" {
				t.Fatalf("Expected error containing %q, but got none.", tc.ErrContains)
			}

			tc.Validate(t, *request, *instanceInfo)
		})
	}
}
