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
	"fmt"
	"os"
	"time"

	googlespanner "cloud.google.com/go/spanner/admin/instance/apiv1"
	"code.cloudfoundry.org/lager"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers"
	. "github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/config"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/models"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/name_generator"
	"github.com/GoogleCloudPlatform/gcp-service-broker/db_service"
	"github.com/GoogleCloudPlatform/gcp-service-broker/fakes"
	"github.com/jinzhu/gorm"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"golang.org/x/net/context"
	"google.golang.org/api/iam/v1"
	"google.golang.org/api/option"
	googlecloudsql "google.golang.org/api/sqladmin/v1beta4"
	instancepb "google.golang.org/genproto/googleapis/spanner/admin/instance/v1"
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
			} else if done.State == models.Succeeded {
				return nil
			}
		}
	}
}

var _ = Describe("AsyncIntegrationTests", func() {
	var (
		gcpBroker            *GCPAsyncServiceBroker
		brokerConfig         *config.BrokerConfig
		err                  error
		logger               lager.Logger
		serviceNameToId      map[string]string = make(map[string]string)
		serviceNameToPlanId  map[string]string = make(map[string]string)
		instance_name        string
		cloudsqlInstanceName string
	)

	BeforeEach(func() {
		logger = lager.NewLogger("brokers_test")
		logger.RegisterSink(lager.NewWriterSink(GinkgoWriter, lager.DEBUG))

		testDb, _ := gorm.Open("sqlite3", "test.db")
		testDb.CreateTable(models.ServiceInstanceDetails{})
		testDb.CreateTable(models.ServiceBindingCredentials{})
		testDb.CreateTable(models.ProvisionRequestDetails{})
		testDb.CreateTable(models.CloudOperation{})

		db_service.DbConnection = testDb

		os.Setenv("SECURITY_USER_NAME", "username")
		os.Setenv("SECURITY_USER_PASSWORD", "password")

		fakes.SetUpTestServices()

		brokerConfig, err = config.NewBrokerConfigFromEnv()

		if err != nil {
			panic(err.Error())
			logger.Error("error", err)
		}

		instance_name = generateInstanceName(brokerConfig.ProjectId, "-")

		// cloudsql instance names need to be random because the recycle time is so long it's untenable to be consistent
		cloudsqlInstanceName = fmt.Sprintf("pcf-sb-test-%d", time.Now().UnixNano())
		name_generator.Basic = &fakes.StaticNameGenerator{Val: instance_name}
		name_generator.Sql = &fakes.StaticSQLNameGenerator{
			StaticNameGenerator: fakes.StaticNameGenerator{Val: cloudsqlInstanceName},
		}

		gcpBroker, err = brokers.New(brokerConfig, logger)
		if err != nil {
			logger.Error("error", err)
		}

		for _, service := range gcpBroker.Catalog {
			serviceNameToId[service.Name] = service.ID
			serviceNameToPlanId[service.Name] = service.Plans[0].ID
		}
	})

	Describe("Cloud SQL - MySQL", func() {

		var dbService *googlecloudsql.InstancesService
		var sslService *googlecloudsql.SslCertsService
		BeforeEach(func() {
			sqlService, err := googlecloudsql.New(brokerConfig.HttpConfig.Client(context.Background()))
			Expect(err).NotTo(HaveOccurred())
			dbService = googlecloudsql.NewInstancesService(sqlService)
			sslService = googlecloudsql.NewSslCertsService(sqlService)
		})

		It("can provision/bind/unbind/deprovision", func() {
			// create the instance
			provisionDetails := models.ProvisionDetails{
				ServiceID: serviceNameToId[models.CloudsqlMySQLName],
				PlanID:    serviceNameToPlanId[models.CloudsqlMySQLName],
			}
			_, err = gcpBroker.Provision("integration_test_instance", provisionDetails, true)
			Expect(err).NotTo(HaveOccurred())
			pollForMaxFiveMins(gcpBroker, "integration_test_instance")

			// make sure it's in the database
			var count int
			db_service.DbConnection.Model(&models.ServiceInstanceDetails{}).Where("id = ?", "integration_test_instance").Count(&count)
			Expect(count).To(Equal(1))

			// make sure we can get it from google
			clouddb, err := dbService.Get(brokerConfig.ProjectId, cloudsqlInstanceName).Do()
			Expect(err).NotTo(HaveOccurred())
			Expect(clouddb.Name).To(Equal(cloudsqlInstanceName))

			// bind the instance
			bindDetails := models.BindDetails{
				ServiceID: serviceNameToId[models.CloudsqlMySQLName],
				PlanID:    serviceNameToPlanId[models.CloudsqlMySQLName],
				Parameters: map[string]interface{}{
					"role": "editor",
				},
			}
			creds, err := gcpBroker.Bind("integration_test_instance", "bind-id", bindDetails)
			Expect(err).NotTo(HaveOccurred())
			credsMap := creds.Credentials.(map[string]string)
			Expect(credsMap["uri"]).To(ContainSubstring("mysql"))

			// make sure we have a username and google has ssl certs
			Expect(credsMap["Username"]).ToNot(Equal(""))
			_, err = sslService.Get(brokerConfig.ProjectId, cloudsqlInstanceName, credsMap["Sha1Fingerprint"]).Do()
			Expect(err).NotTo(HaveOccurred())

			// unbind the instance
			unBindDetails := models.UnbindDetails{
				ServiceID: serviceNameToId[models.CloudsqlMySQLName],
				PlanID:    serviceNameToPlanId[models.CloudsqlMySQLName],
			}
			err = gcpBroker.Unbind("integration_test_instance", "bind-id", unBindDetails)
			Expect(err).NotTo(HaveOccurred())

			// make sure google no longer has certs
			certsList, err := sslService.List(brokerConfig.ProjectId, cloudsqlInstanceName).Do()
			Expect(len(certsList.Items)).To(Equal(0))

			// deprovision the instance
			deprovisionDetails := models.DeprovisionDetails{
				ServiceID: serviceNameToId[models.CloudsqlMySQLName],
				PlanID:    serviceNameToPlanId[models.CloudsqlMySQLName],
			}
			_, err = gcpBroker.Deprovision("integration_test_instance", deprovisionDetails, true)
			Expect(err).NotTo(HaveOccurred())
			pollForMaxFiveMins(gcpBroker, "integration_test_instance")

			// make sure the instance is deleted from the db
			instance := models.ServiceInstanceDetails{}
			if err := db_service.DbConnection.Unscoped().Where("ID = ?", "integration_test_instance").First(&instance).Error; err != nil {
				panic("error checking for service instance details: " + err.Error())
			}
			Expect(instance.DeletedAt).NotTo(BeNil())

			// make sure the instance is deleted from google
			_, err = dbService.Get(brokerConfig.ProjectId, cloudsqlInstanceName).Do()
			Expect(err).To(HaveOccurred())
		})

	})

	Describe("Cloud SQL - PostgreSQL", func() {

		var sqlService *googlecloudsql.Service
		var dbService *googlecloudsql.InstancesService
		var sslService *googlecloudsql.SslCertsService
		BeforeEach(func() {
			sqlService, err := googlecloudsql.New(brokerConfig.HttpConfig.Client(context.Background()))
			Expect(err).NotTo(HaveOccurred())
			dbService = googlecloudsql.NewInstancesService(sqlService)
			sslService = googlecloudsql.NewSslCertsService(sqlService)
		})

		It("can provision/bind/unbind/deprovision", func() {
			// create the instance
			provisionDetails := models.ProvisionDetails{
				ServiceID: serviceNameToId[models.CloudsqlPostgresName],
				PlanID:    serviceNameToPlanId[models.CloudsqlPostgresName],
			}
			_, err = gcpBroker.Provision("integration_test_instance", provisionDetails, true)
			Expect(err).NotTo(HaveOccurred())
			pollForMaxFiveMins(gcpBroker, "integration_test_instance")

			// make sure it's in the database
			var count int
			db_service.DbConnection.Model(&models.ServiceInstanceDetails{}).Where("id = ?", "integration_test_instance").Count(&count)
			Expect(count).To(Equal(1))

			// make sure we can get it from google
			clouddb, err := dbService.Get(brokerConfig.ProjectId, cloudsqlInstanceName).Do()
			Expect(err).NotTo(HaveOccurred())
			Expect(clouddb.Name).To(Equal(cloudsqlInstanceName))

			// bind the instance
			bindDetails := models.BindDetails{
				ServiceID: serviceNameToId[models.CloudsqlPostgresName],
				PlanID:    serviceNameToPlanId[models.CloudsqlPostgresName],
				Parameters: map[string]interface{}{
					"role": "editor",
				},
			}
			creds, err := gcpBroker.Bind("integration_test_instance", "bind-id", bindDetails)
			Expect(err).NotTo(HaveOccurred())
			credsMap := creds.Credentials.(map[string]string)

			// make sure we have a username and google has ssl certs
			Expect(credsMap["Username"]).ToNot(Equal(""))
			_, err = sslService.Get(brokerConfig.ProjectId, cloudsqlInstanceName, credsMap["Sha1Fingerprint"]).Do()
			Expect(err).NotTo(HaveOccurred())
			Expect(credsMap["uri"]).To(ContainSubstring("postgres"))

			// The bind operation finishes before the instance has finished updating
			keepWaiting := true
			for keepWaiting == true {
				resp, err := sqlService.Operations.List(brokerConfig.ProjectId, cloudsqlInstanceName).Do()
				Expect(err).NotTo(HaveOccurred())
				keepWaiting = false
				for _, op := range resp.Items {
					if op.EndTime == "" {
						keepWaiting = true
					}
				}
				// sleep for 1 second between polling so we don't hit our rate limit
				time.Sleep(time.Second)
			}

			// unbind the instance
			unBindDetails := models.UnbindDetails{
				ServiceID: serviceNameToId[models.CloudsqlPostgresName],
				PlanID:    serviceNameToPlanId[models.CloudsqlPostgresName],
			}

			err = gcpBroker.Unbind("integration_test_instance", "bind-id", unBindDetails)
			Expect(err).NotTo(HaveOccurred())

			// make sure google no longer has certs
			certsList, err := sslService.List(brokerConfig.ProjectId, cloudsqlInstanceName).Do()
			Expect(len(certsList.Items)).To(Equal(0))

			// deprovision the instance
			deprovisionDetails := models.DeprovisionDetails{
				ServiceID: serviceNameToId[models.CloudsqlPostgresName],
				PlanID:    serviceNameToPlanId[models.CloudsqlPostgresName],
			}

			_, err = gcpBroker.Deprovision("integration_test_instance", deprovisionDetails, true)
			Expect(err).NotTo(HaveOccurred())
			pollForMaxFiveMins(gcpBroker, "integration_test_instance")

			// make sure the instance is deleted from the db
			instance := models.ServiceInstanceDetails{}
			if err := db_service.DbConnection.Unscoped().Where("ID = ?", "integration_test_instance").First(&instance).Error; err != nil {
				panic("error checking for service instance details: " + err.Error())
			}
			Expect(instance.DeletedAt).NotTo(BeNil())

			// make sure the instance is deleted from google
			_, err = dbService.Get(brokerConfig.ProjectId, cloudsqlInstanceName).Do()
			Expect(err).To(HaveOccurred())
		})

	})

	Describe("Spanner", func() {

		var client *googlespanner.InstanceAdminClient
		BeforeEach(func() {
			co := option.WithUserAgent(models.CustomUserAgent)
			ct := option.WithTokenSource(brokerConfig.HttpConfig.TokenSource(context.Background()))
			client, err = googlespanner.NewInstanceAdminClient(context.Background(), co, ct)
			Expect(err).ToNot(HaveOccurred())
		})

		It("can provision/bind/unbind/deprovision", func() {
			provisionDetails := models.ProvisionDetails{
				ServiceID: serviceNameToId[models.SpannerName],
				PlanID:    serviceNameToPlanId[models.SpannerName],
			}
			_, err = gcpBroker.Provision("integration_test_instance", provisionDetails, true)
			Expect(err).NotTo(HaveOccurred())
			err = pollForMaxFiveMins(gcpBroker, "integration_test_instance")
			Expect(err).NotTo(HaveOccurred())

			var count int
			db_service.DbConnection.Model(&models.ServiceInstanceDetails{}).Where("id = ?", "integration_test_instance").Count(&count)
			Expect(count).To(Equal(1))

			_, err = client.GetInstance(context.Background(), &instancepb.GetInstanceRequest{
				Name: "projects/" + brokerConfig.ProjectId + "/instances/" + instance_name,
			})
			Expect(err).ToNot(HaveOccurred())

			bindDetails := models.BindDetails{
				ServiceID: serviceNameToId[models.SpannerName],
				PlanID:    serviceNameToPlanId[models.SpannerName],
				Parameters: map[string]interface{}{
					"role": "spanner.admin",
				},
			}
			creds, err := gcpBroker.Bind("integration_test_instance", "bind-id", bindDetails)
			Expect(err).NotTo(HaveOccurred())
			credsMap := creds.Credentials.(map[string]string)
			Expect(credsMap["credentials"]).ToNot(BeNil())

			iamService, err := iam.New(brokerConfig.HttpConfig.Client(context.Background()))
			Expect(err).ToNot(HaveOccurred())
			saService := iam.NewProjectsServiceAccountsService(iamService)
			bindResourceName := "projects/" + brokerConfig.ProjectId + "/serviceAccounts/" + creds.Credentials.(map[string]string)["UniqueId"]
			_, err = saService.Get(bindResourceName).Do()
			Expect(err).ToNot(HaveOccurred())

			unBindDetails := models.UnbindDetails{
				ServiceID: serviceNameToId[models.SpannerName],
				PlanID:    serviceNameToPlanId[models.SpannerName],
			}
			err = gcpBroker.Unbind("integration_test_instance", "bind-id", unBindDetails)
			Expect(err).NotTo(HaveOccurred())

			_, err = saService.Get(bindResourceName).Do()
			Expect(err).To(HaveOccurred())

			deprovisionDetails := models.DeprovisionDetails{
				ServiceID: serviceNameToId[models.SpannerName],
				PlanID:    serviceNameToPlanId[models.SpannerName],
			}
			_, err = gcpBroker.Deprovision("integration_test_instance", deprovisionDetails, true)
			Expect(err).NotTo(HaveOccurred())

			instance := models.ServiceInstanceDetails{}
			if err := db_service.DbConnection.Unscoped().Where("ID = ?", "integration_test_instance").First(&instance).Error; err != nil {
				panic("error checking for service instance details: " + err.Error())
			}
			Expect(instance.DeletedAt).NotTo(BeNil())

			_, err = client.GetInstance(context.Background(), &instancepb.GetInstanceRequest{
				Name: "projects/" + brokerConfig.ProjectId + "/instances/" + instance_name,
			})
			Expect(err).To(HaveOccurred())
		})

	})

	AfterEach(func() {
		os.Remove("test.db")
	})
})
