// Copyright the Service Broker Project Authors.
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
//
////////////////////////////////////////////////////////////////////////////////
//

package integration_tests

import (
	"os"
	"gcp-service-broker/brokerapi/brokers"
	"gcp-service-broker/brokerapi/brokers/name_generator"
	"github.com/jinzhu/gorm"
	"gcp-service-broker/db_service"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"code.cloudfoundry.org/lager"
	. "gcp-service-broker/brokerapi/brokers"
	"gcp-service-broker/brokerapi/brokers/models"
	"gcp-service-broker/fakes"
	"time"
)

func pollForMaxFiveMins(gcpb *GCPAsyncServiceBroker, instanceId string) error {
	var err error
	timeout := time.After(5 * time.Minute)
	tick := time.Tick(30 * time.Second)

	// Keep trying until we're timed out or got a result or got an error
	for {
		select {
		case <-timeout:
			return err
		case <-tick:
			done, err := gcpb.LastOperation(instanceId)
			if err != nil {
				return err
			} else if done {
				return nil
			}
		}
	}
}

var _ = Describe("AsyncIntegrationTests", func() {
	var (
		gcpBroker           *GCPAsyncServiceBroker
		err                 error
		logger              lager.Logger
		serviceNameToId     map[string]string = make(map[string]string)
		serviceNameToPlanId map[string]string = make(map[string]string)
		instance_name       string
	)

	BeforeEach(func() {
		logger = lager.NewLogger("brokers_test")
		logger.RegisterSink(lager.NewWriterSink(GinkgoWriter, lager.DEBUG))

		testDb, _ := gorm.Open("sqlite3", "test.db")
		testDb.CreateTable(models.ServiceInstanceDetails{})
		testDb.CreateTable(models.ServiceBindingCredentials{})
		testDb.CreateTable(models.PlanDetails{})
		testDb.CreateTable(models.ProvisionRequestDetails{})

		db_service.DbConnection = testDb

		os.Setenv("SECURITY_USER_NAME", "username")
		os.Setenv("SECURITY_USER_PASSWORD", "password")
		os.Setenv("SERVICES", fakes.Services)
		os.Setenv("PRECONFIGURED_PLANS", fakes.PreconfiguredPlans)

		os.Setenv("CLOUDSQL_CUSTOM_PLANS", `{
			"test_cloudsql_plan": {
				"guid": "foo",
				"name": "bar",
				"description": "test-cloudsql-plan",
				"tier": "D4",
				"pricing_plan": "PER_USE",
				"max_disk_size": "20",
				"display_name": "FOOBAR",
				"service": "4bc59b9a-8520-409f-85da-1c7552315863"
			}
		}`)

		os.Setenv("BIGTABLE_CUSTOM_PLANS", `{
			"test_bigtable_plan": {
				"guid": "foo2",
				"name": "bar2",
				"description": "test-bigtable-plan",
				"storage_type": "SSD",
				"num_nodes": "3",
				"display_name": "FOOBAR2",
				"service": "b8e19880-ac58-42ef-b033-f7cd9c94d1fe"
			}
		}`)

		var creds models.GCPCredentials
		creds, err = brokers.GetCredentialsFromEnv()
		if err != nil {
			logger.Error("error", err)
		}
		instance_name = generateInstanceName(creds.ProjectId, "")
		name_generator.Basic = &fakes.StaticNameGenerator{Val: instance_name}

		gcpBroker, err = brokers.New(logger)
		if err != nil {
			logger.Error("error", err)
		}

		for _, service := range *gcpBroker.Catalog {
			serviceNameToId[service.Name] = service.ID
			serviceNameToPlanId[service.Name] = service.Plans[0].ID
		}
	})
})
