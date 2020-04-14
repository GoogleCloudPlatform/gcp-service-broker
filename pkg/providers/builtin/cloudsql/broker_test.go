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
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/gcp-service-broker/db_service/models"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/broker"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/varcontext"
	"github.com/pivotal-cf/brokerapi"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
	googlecloudsql "google.golang.org/api/sqladmin/v1beta4"
	"k8s.io/utils/diff"
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
	defer viper.Reset()

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
			Service:    MysqlServiceDefinition(),
			PlanId:     mysqlSecondGenPlan,
			UserParams: `{"instance_name":""}`,
			Validate: func(t *testing.T, di googlecloudsql.DatabaseInstance, ii InstanceInformation) {
				if len(di.Name) == 0 {
					t.Errorf("instance name wasn't generated")
				}
			},
		},

		"tiers matching (D|d)\\d+ get firstgen outputs for MySQL": {
			Service:     MysqlServiceDefinition(),
			PlanId:      mysqlFirstgenPlan,
			UserParams:  `{"name":""}`,
			ErrContains: "First generation support will end March 25th, 2020, please use a second gen machine type",
		},

		"second-gen MySQL defaults": {
			Service:    MysqlServiceDefinition(),
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
			Service:    PostgresServiceDefinition(),
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
			Service:    MysqlServiceDefinition(),
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
			Service:    MysqlServiceDefinition(),
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
			Service:    MysqlServiceDefinition(),
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
			Service:    MysqlServiceDefinition(),
			PlanId:     mysqlSecondGenPlan,
			UserParams: `{"database_name":""}`,
			Validate: func(t *testing.T, di googlecloudsql.DatabaseInstance, ii InstanceInformation) {
				if len(ii.DatabaseName) == 0 {
					t.Error("Expected DatabaseName to not be blank.")
				}
			},
		},

		"instance info has name and db name ": {
			Service:    MysqlServiceDefinition(),
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
			Service:     MysqlServiceDefinition(),
			PlanId:      mysqlSecondGenPlan,
			UserParams:  `{"disk_size":"99999"}`,
			Validate:    func(t *testing.T, di googlecloudsql.DatabaseInstance, ii InstanceInformation) {},
			ErrContains: "disk size",
		},

		"postgres disk size greater than operator specified max fails": {
			Service:     PostgresServiceDefinition(),
			PlanId:      postgresPlan,
			UserParams:  `{"disk_size":"99999"}`,
			Validate:    func(t *testing.T, di googlecloudsql.DatabaseInstance, ii InstanceInformation) {},
			ErrContains: "disk size",
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			details := brokerapi.ProvisionDetails{RawParameters: json.RawMessage(tc.UserParams), ServiceID: tc.Service.Id}
			plan, err := tc.Service.GetPlanById(tc.PlanId)
			if err != nil {
				t.Fatalf("got error trying to find plan %s %v", tc.PlanId, err)
			}
			if plan == nil {
				t.Fatalf("Expected plan with id %s to not be nil", tc.PlanId)
			}

			vars, err := tc.Service.ProvisionVariables("instance-id-here", details, *plan)
			if err != nil {
				if tc.ErrContains != "" && strings.Contains(err.Error(), tc.ErrContains) {
					return
				}

				t.Fatalf("got error trying to get provision details %s %v", tc.PlanId, err)
			}

			request, instanceInfo, err := createProvisionRequest(vars)
			if err != nil {

				t.Fatalf("got unexpected error while creating provision request: %v", err)
			}
			if tc.ErrContains != "" {
				t.Fatalf("Expected error containing %q, but got none.", tc.ErrContains)
			}

			tc.Validate(t, *request, *instanceInfo)
		})
	}
}

func TestPostgresCustomMachineTypes(t *testing.T) {
	for _, plan := range PostgresServiceDefinition().Plans {
		t.Run(plan.Name, func(t *testing.T) {
			props := plan.ServiceProperties

			tier := props["tier"]
			if !strings.HasPrefix(tier, "db-custom") {
				return
			}

			splitTier := strings.Split(tier, "-")
			if len(splitTier) != 4 {
				t.Errorf("Expected custom machine type to be in format db-custom-[NCPU]-[MEM_MB]")
			}

			cpu := cast.ToInt(splitTier[2])
			mem := cast.ToInt(splitTier[3])

			// Rules for custom machines
			// https://cloud.google.com/compute/docs/instances/creating-instance-with-custom-machine-type#create
			if cpu != 1 && cpu%2 != 0 {
				t.Errorf("Only machine types with 1 vCPU or an even number of vCPUs can be created got %d", cpu)
			}

			memPerCpu := float64(mem) / float64(cpu)
			if memPerCpu < (.9*1024) || memPerCpu > (6.5*1024) {
				t.Errorf("Memory must be between 0.9 GB per vCPU, up to 6.5 GB per vCPU got %f MB/CPU", memPerCpu)
			}

			if mem%256 != 0 {
				t.Errorf("The total memory of the instance must be a multiple of 256 MB, got: %d MB", mem)
			}

			if cpu > 64 {
				t.Errorf("The maximum number of CPUs allowed are 64, got %d", cpu)
			}
		})
	}
}

func TestBuildInstanceCredentials(t *testing.T) {
	pgProvider := PostgresServiceDefinition().ProviderBuilder("project-id", nil, nil)
	mysqlProvider := MysqlServiceDefinition().ProviderBuilder("project-id", nil, nil)

	cases := map[string]struct {
		serviceID       string
		bindDetails     string
		instanceDetails string
		provider        broker.ServiceProvider

		expectedCreds map[string]interface{}
	}{
		"no prefix mysql": {
			serviceID: MySqlServiceId,
			instanceDetails: `{
				"database_name": "sb-2-1543346570614873901",
				"host": "35.202.18.12"
			}`,
			bindDetails: `{
				"Username": "sb15433468744175",
				"Password": "pass=",
				"UriPrefix": ""
			}`,
			provider: mysqlProvider,
			expectedCreds: map[string]interface{}{
				"Password":      "pass=",
				"UriPrefix":     "",
				"Username":      "sb15433468744175",
				"database_name": "sb-2-1543346570614873901",
				"host":          "35.202.18.12",
				"uri":           "mysql://sb15433468744175:pass%3D@35.202.18.12/sb-2-1543346570614873901?ssl_mode=required",
			},
		},

		"no prefix postgres": {
			serviceID: PostgresServiceId,
			instanceDetails: `{
				"database_name": "sb-2-1543346570614873901",
				"host": "35.202.18.12"
			}`,
			bindDetails: `{
				"ClientCert": "@clientcert",
				"ClientKey": "@clientkey",
				"CaCert": "@cacert",
				"Username": "sb15433468744175",
				"Password": "pass=",
				"UriPrefix": ""
			}`,
			provider: pgProvider,
			expectedCreds: map[string]interface{}{
				"ClientCert":    "@clientcert",
				"ClientKey":     "@clientkey",
				"CaCert":        "@cacert",
				"Password":      "pass=",
				"UriPrefix":     "",
				"Username":      "sb15433468744175",
				"database_name": "sb-2-1543346570614873901",
				"host":          "35.202.18.12",
				"uri":           "postgres://sb15433468744175:pass%3D@35.202.18.12/sb-2-1543346570614873901?sslmode=require&sslcert=%40clientcert&sslkey=%40clientkey&sslrootcert=%40cacert",
			},
		},

		"prefix mysql": {
			serviceID: MySqlServiceId,
			instanceDetails: `{
				"database_name": "sb-2-1543346570614873901",
				"host": "35.202.18.12"
			}`,
			bindDetails: `{
				"Username": "sb15433468744175",
				"Password": "pass=",
				"UriPrefix": "jdbc:"
			}`,
			provider: mysqlProvider,
			expectedCreds: map[string]interface{}{
				"Password":      "pass=",
				"UriPrefix":     "jdbc:",
				"Username":      "sb15433468744175",
				"database_name": "sb-2-1543346570614873901",
				"host":          "35.202.18.12",
				"uri":           "jdbc:mysql://sb15433468744175:pass%3D@35.202.18.12/sb-2-1543346570614873901?ssl_mode=required",
			},
		},

		"prefix postgres": {
			serviceID: PostgresServiceId,
			instanceDetails: `{
				"database_name": "sb-2-1543346570614873901",
				"host": "35.202.18.12"
			}`,
			bindDetails: `{
				"ClientCert": "@clientcert",
				"ClientKey": "@clientkey",
				"CaCert": "@cacert",
				"Username": "sb15433468744175",
				"Password": "pass=",
				"UriPrefix": "jdbc:"
			}`,
			provider: pgProvider,
			expectedCreds: map[string]interface{}{
				"ClientCert":    "@clientcert",
				"ClientKey":     "@clientkey",
				"CaCert":        "@cacert",
				"Password":      "pass=",
				"UriPrefix":     "jdbc:",
				"Username":      "sb15433468744175",
				"database_name": "sb-2-1543346570614873901",
				"host":          "35.202.18.12",
				"uri":           "jdbc:postgres://sb15433468744175:pass%3D@35.202.18.12/sb-2-1543346570614873901?sslmode=require&sslcert=%40clientcert&sslkey=%40clientkey&sslrootcert=%40cacert",
			},
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {

			bindRecord := models.ServiceBindingCredentials{
				OtherDetails: tc.bindDetails,
			}

			instanceRecord := models.ServiceInstanceDetails{
				ServiceId:    tc.serviceID,
				OtherDetails: tc.instanceDetails,
			}

			binding, err := tc.provider.BuildInstanceCredentials(context.Background(), bindRecord, instanceRecord)
			if err != nil {
				t.Error("expected no error, got:", err)
				return
			}

			if !reflect.DeepEqual(binding.Credentials, tc.expectedCreds) {
				t.Errorf("Expected credentials %#v, got %#v", tc.expectedCreds, binding.Credentials)
			}
		})
	}
}

func Test_createProvisionRequest(t *testing.T) {
	services := []*broker.ServiceDefinition{
		MysqlServiceDefinition(),
		PostgresServiceDefinition(),
		MySQLVPCServiceDefinition(),
		PostgresVPCServiceDefinition(),
	}

	for _, service := range services {
		t.Run(service.Name, func(t *testing.T) {
			for _, example := range service.Examples {
				t.Run(example.Name, func(t *testing.T) {
					vars := deterministicProvisionVariables(t, service, example)

					req, _, err := createProvisionRequest(vars)
					if err != nil {
						t.Fatalf("couldn't createProvisionRequest: %v", err)
					}

					assertGolden(t, req, map[string]interface{}{
						"example": example,
					})
				})
			}
		})
	}
}

func assertGolden(t *testing.T, object interface{}, context map[string]interface{}) {
	t.Helper()

	contextText, err := json.MarshalIndent(context, "// ", "\t")
	if err != nil {
		t.Errorf("creating golden context: %v", err)
	}

	objectText, err := json.MarshalIndent(object, "", "\t")
	if err != nil {
		t.Errorf("creating golden object: %v", err)
	}

	goldenContents := fmt.Sprintf("// %s\n%s", string(contextText), string(objectText))

	testPath := filepath.FromSlash(t.Name())
	goldenPrefix := append([]string{"testdata", "golden"}, filepath.SplitList(testPath)...)
	goldenPath := filepath.Join(goldenPrefix...)

	if os.Getenv("UPDATE_GOLDEN") == "true" {
		if err := os.MkdirAll(filepath.Dir(goldenPath), 0700); err != nil {
			t.Fatalf("creating parent directory: %v", err)
		}

		if err := ioutil.WriteFile(goldenPath, []byte(goldenContents), 0600); err != nil {
			t.Fatalf("creating golden file: %v", err)
		}
		t.Fatal("golden file updated, run without UPDATE_GOLDEN=true to test")
	} else {
		wantContents, err := ioutil.ReadFile(goldenPath)
		if err != nil {
			t.Fatalf("couldn't read golden file, run with UPDATE_GOLDEN=true to update: %v", err)
		}

		if string(wantContents) != string(goldenContents) {
			diff := diff.StringDiff(string(wantContents), string(goldenContents))
			t.Errorf("actual differed from golden: %s", diff)
		}
	}
}

func deterministicProvisionVariables(t *testing.T, service *broker.ServiceDefinition, example broker.ServiceExample) *varcontext.VarContext {
	t.Helper()

	bytes, err := json.Marshal(example.ProvisionParams)
	if err != nil {
		t.Fatalf("couldn't marshal ProvisionParams for example: %v", err)
	}

	details := brokerapi.ProvisionDetails{
		RawParameters: bytes,
		ServiceID:     service.Id,
	}

	plan, err := service.GetPlanById(example.PlanId)
	if err != nil {
		t.Fatalf("got error trying to find plan %s %v", example.PlanId, err)
	}

	vars1, err := service.ProvisionVariables("instance-id-here", details, *plan)
	if err != nil {
		t.Fatalf("couldn't create ProvisionVariables1: %v", err)
	}

	vars2, err := service.ProvisionVariables("instance-id-here", details, *plan)
	if err != nil {
		t.Fatalf("couldn't create ProvisionVariables2: %v", err)
	}

	vm1 := vars1.ToMap()
	vm2 := vars2.ToMap()

	merged := make(map[string]interface{})

	for k := range vm1 {
		if reflect.DeepEqual(vm1[k], vm2[k]) {
			merged[k] = vm1[k]
		} else {
			// TODO: add additional types if necessary
			merged[k] = "NONDETERMINISTIC"
		}
	}

	context, err := varcontext.Builder().MergeMap(merged).Build()
	if err != nil {
		t.Fatalf("creating varcontext: %v", err)
	}

	return context
}
