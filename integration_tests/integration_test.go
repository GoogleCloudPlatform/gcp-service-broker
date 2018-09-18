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

package integration_tests

import (
	"time"

	. "github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/broker"
	"github.com/pivotal-cf/brokerapi"

	"golang.org/x/net/context"

	"fmt"
	"hash/crc32"
	"net/http"
	"os"

	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/models"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/name_generator"
	"github.com/GoogleCloudPlatform/gcp-service-broker/db_service"
	"github.com/GoogleCloudPlatform/gcp-service-broker/fakes"

	googlepubsub "cloud.google.com/go/pubsub"

	"encoding/json"

	googlebigtable "cloud.google.com/go/bigtable"
	googlestorage "cloud.google.com/go/storage"
	"code.cloudfoundry.org/lager"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/config"
	"github.com/jinzhu/gorm"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	googlebigquery "google.golang.org/api/bigquery/v2"
	iam "google.golang.org/api/iam/v1"
	"google.golang.org/api/option"
)

const timeout = 60

type genericService struct {
	serviceId              string
	planId                 string
	bindingId              string
	rawBindingParams       json.RawMessage
	instanceId             string
	serviceExistsFn        func(bool) bool
	cleanupFn              func()
	serviceMetadataSavedFn func(string) (bool, error)
}

type iamService struct {
	bindingId string
	serviceId string
	planId    string
}

func getAndUnmarshalInstanceDetails(instanceId string) (map[string]string, error) {
	instanceRecord, _ := db_service.GetServiceInstanceDetailsById(instanceId)
	return instanceRecord.GetOtherDetails()
}

func testGenericService(brokerConfig *config.BrokerConfig, gcpBroker *GCPServiceBroker, params *genericService) {
	// If the service already exists (eg, failed previous test), clean it up before the run
	if params.serviceExistsFn != nil && params.serviceExistsFn(false) {
		params.cleanupFn()
	}
	//
	// Provision
	//
	provisionDetails := brokerapi.ProvisionDetails{
		ServiceID: params.serviceId,
		PlanID:    params.planId,
	}

	_, err := gcpBroker.Provision(context.Background(), params.instanceId, provisionDetails, true)
	Expect(err).ToNot(HaveOccurred())

	// Provision is registered in the database
	count, err := db_service.CountServiceInstanceDetailsById(params.instanceId)
	Expect(err).NotTo(HaveOccurred())
	Expect(count).To(Equal(1))

	if params.serviceExistsFn != nil {
		Expect(params.serviceExistsFn(true)).To(BeTrue())
	}
	metadataSaved, err := params.serviceMetadataSavedFn(params.instanceId)
	Expect(err).NotTo(HaveOccurred())
	Expect(metadataSaved).To(BeTrue())

	//
	// Bind
	//
	bindDetails := brokerapi.BindDetails{
		ServiceID:     params.serviceId,
		PlanID:        params.planId,
		RawParameters: params.rawBindingParams,
	}
	creds, err := gcpBroker.Bind(context.Background(), params.instanceId, params.bindingId, bindDetails)
	Expect(err).ToNot(HaveOccurred())

	count, err = db_service.CountServiceBindingCredentialsByBindingId(params.bindingId)
	Expect(err).ToNot(HaveOccurred())
	Expect(count).To(Equal(1))

	iamService, err := iam.New(brokerConfig.HttpConfig.Client(context.Background()))
	Expect(err).ToNot(HaveOccurred())
	saService := iam.NewProjectsServiceAccountsService(iamService)
	resourceName := "projects/" + brokerConfig.ProjectId + "/serviceAccounts/" + creds.Credentials.(map[string]string)["UniqueId"]
	_, err = saService.Get(resourceName).Do()
	Expect(err).ToNot(HaveOccurred())

	//
	// Unbind
	//
	unbindDetails := brokerapi.UnbindDetails{
		ServiceID: params.serviceId,
		PlanID:    params.planId,
	}
	err = gcpBroker.Unbind(context.Background(), params.instanceId, params.bindingId, unbindDetails)
	Expect(err).ToNot(HaveOccurred())

	deleted, err := db_service.CheckDeletedServiceBindingCredentialsByBindingId(params.bindingId)
	Expect(err).ToNot(HaveOccurred())
	Expect(deleted).To(BeTrue())

	// wait because services don't always show as deleted right away
	time.Sleep(5 * time.Second)
	_, err = saService.Get(resourceName).Do()
	Expect(err).NotTo(BeNil())

	//
	// Deprovision

	deprovisionDetails := brokerapi.DeprovisionDetails{
		ServiceID: params.serviceId,
		PlanID:    params.planId,
	}
	_, err = gcpBroker.Deprovision(context.Background(), params.instanceId, deprovisionDetails, true)
	Expect(err).ToNot(HaveOccurred())

	deleted, err = db_service.CheckDeletedServiceInstanceDetailsById(params.instanceId)
	Expect(err).NotTo(HaveOccurred())
	Expect(deleted).To(BeTrue())

	if params.serviceExistsFn != nil {
		Expect(params.serviceExistsFn(false)).To(BeFalse())
	}
}

// For services that only create a service account and bind those credentials.
func testIamBasedService(brokerConfig *config.BrokerConfig, gcpBroker *GCPServiceBroker, params *iamService) {
	genericServiceParams := &genericService{
		serviceId:        params.serviceId,
		planId:           params.planId,
		instanceId:       "iam-instance",
		bindingId:        "iam-instance",
		rawBindingParams: json.RawMessage{},
		serviceMetadataSavedFn: func(instanceId string) (bool, error) {
			// Metadata should be empty, there is no additional information required
			instanceDetails, err := getAndUnmarshalInstanceDetails(instanceId)
			return len(instanceDetails) == 0, err
		},
	}

	testGenericService(brokerConfig, gcpBroker, genericServiceParams)
}

// Instance Name is used to name every instance created in GCP (eg, a storage bucket)
// The name should be consistent between runs to ensure there's bounds to the resources it creates
// and to have some insurance that they are properly destroyed.
//
// Why:
// - If we allow it to generate a random instance name every time the test will
//   not fail if the resource existed before hand.
// - If we always use a static one, globally named resources (eg, a storage bucket)
//   would fail to create when two different projects run these tests.
func generateInstanceName(projectId string, sep string) string {
	hashed := crc32.ChecksumIEEE([]byte(projectId))
	if sep == "" {
		sep = "_"
	}
	return fmt.Sprintf("pcf%ssb%s1%s%d", sep, sep, sep, hashed)
}

var _ = Describe("LiveIntegrationTests", func() {
	var (
		brokerConfig        *config.BrokerConfig
		gcpBroker           *GCPServiceBroker
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
		db_service.RunMigrations(testDb)
		db_service.DbConnection = testDb

		os.Setenv("SECURITY_USER_NAME", "username")
		os.Setenv("SECURITY_USER_PASSWORD", "password")

		fakes.SetUpTestServices()

		brokerConfig, err = config.NewBrokerConfigFromEnv()
		if err != nil {
			logger.Error("error", err)
		}

		instance_name = generateInstanceName(brokerConfig.ProjectId, "")
		name_generator.Basic = &fakes.StaticNameGenerator{Val: instance_name}

		gcpBroker, err = brokers.New(brokerConfig, logger)
		if err != nil {
			logger.Error("error", err)
		}

		for _, service := range gcpBroker.Catalog {
			serviceNameToId[service.Name] = service.ID
			serviceNameToPlanId[service.Name] = service.Plans[0].ID
		}
	})

	Describe("Broker init", func() {
		It("should have all enabled services in sevices map", func() {
			Expect(len(gcpBroker.ServiceBrokerMap)).To(Equal(len(broker.GetEnabledServices())))
		})

		It("should have a default client", func() {
			Expect(brokerConfig.HttpConfig).NotTo(Equal(&http.Client{}))
		})

		It("should have loaded credentials correctly and have a project id", func() {
			Expect(brokerConfig.ProjectId).ToNot(BeEmpty())
		})
	})

	Describe("getting broker catalog", func() {
		It("should have all enabled services available", func() {
			serviceList, err := gcpBroker.Services(context.Background())

			Expect(err).ToNot(HaveOccurred())
			Expect(len(serviceList)).To(Equal(len(broker.GetEnabledServices())))
		})

		It("should have 3 storage plans available", func() {
			serviceList, err := gcpBroker.Services(context.Background())
			Expect(err).ToNot(HaveOccurred())

			for _, s := range serviceList {
				if s.ID == serviceNameToId[models.StorageName] {
					Expect(len(s.Plans)).To(Equal(3))
				}
			}
		})
	})

	Describe("bigquery", func() {
		It("can provision/bind/unbind/deprovision", func() {
			service, err := googlebigquery.New(brokerConfig.HttpConfig.Client(context.Background()))
			Expect(err).NotTo(HaveOccurred())

			params := &genericService{
				serviceId:        serviceNameToId[models.BigqueryName],
				planId:           serviceNameToPlanId[models.BigqueryName],
				bindingId:        "integration_test_bind",
				instanceId:       "integration_test_dataset",
				rawBindingParams: []byte(`{"role": "bigquery.dataOwner"}`),
				serviceExistsFn: func(expected bool) bool {
					_, err = service.Datasets.Get(brokerConfig.ProjectId, instance_name).Do()

					return err == nil
				},
				serviceMetadataSavedFn: func(instanceId string) (bool, error) {
					instanceDetails, err := getAndUnmarshalInstanceDetails(instanceId)
					return instanceDetails["dataset_id"] != "", err
				},
				cleanupFn: func() {
					err := service.Datasets.Delete(brokerConfig.ProjectId, instance_name).Do()
					Expect(err).NotTo(HaveOccurred())
				},
			}
			testGenericService(brokerConfig, gcpBroker, params)
		}, timeout)
	})

	Describe("bigtable", func() {
		var bigtableInstanceName string
		BeforeEach(func() {
			bigtableInstanceName = generateInstanceName(brokerConfig.ProjectId, "-")
			name_generator.Basic = &fakes.StaticNameGenerator{Val: bigtableInstanceName}
		})

		AfterEach(func() {
			name_generator.Basic = &fakes.StaticNameGenerator{Val: instance_name}
		})

		It("can provision/bind/unbind/deprovision", func() {

			ctx := context.Background()
			co := option.WithUserAgent(models.CustomUserAgent)
			ct := option.WithTokenSource(brokerConfig.HttpConfig.TokenSource(context.Background()))
			service, err := googlebigtable.NewInstanceAdminClient(ctx, brokerConfig.ProjectId, co, ct)
			Expect(err).NotTo(HaveOccurred())

			params := &genericService{
				serviceId:        serviceNameToId[models.BigtableName],
				planId:           serviceNameToPlanId[models.BigtableName],
				bindingId:        "integration_test_bind",
				instanceId:       "integration_test_instance",
				rawBindingParams: []byte(`{"role": "bigtable.user"}`),
				serviceExistsFn: func(expected bool) bool {
					instances, err := service.Instances(ctx)

					return err == nil && len(instances) == 1 && instances[0].Name == bigtableInstanceName
				},
				serviceMetadataSavedFn: func(instanceId string) (bool, error) {
					instanceDetails, err := getAndUnmarshalInstanceDetails(instanceId)
					return instanceDetails["instance_id"] != "", err
				},
				cleanupFn: func() {
					err := service.DeleteInstance(ctx, bigtableInstanceName)
					Expect(err).NotTo(HaveOccurred())
				},
			}
			testGenericService(brokerConfig, gcpBroker, params)
		}, timeout)
	})

	Describe("cloud storage", func() {
		It("can provision/bind/unbind/deprovision", func() {
			co := option.WithUserAgent(models.CustomUserAgent)
			ct := option.WithTokenSource(brokerConfig.HttpConfig.TokenSource(context.Background()))
			service, err := googlestorage.NewClient(context.Background(), co, ct)
			Expect(err).NotTo(HaveOccurred())

			params := &genericService{
				serviceId:        serviceNameToId[models.StorageName],
				planId:           serviceNameToPlanId[models.StorageName],
				instanceId:       "integration_test_bucket",
				bindingId:        "integration_test_bucket_binding",
				rawBindingParams: []byte(`{"role": "storage.objectAdmin"}`),
				serviceExistsFn: func(bool) bool {
					bucket := service.Bucket(instance_name)
					_, err = bucket.Attrs(context.Background())

					return err == nil
				},
				serviceMetadataSavedFn: func(instanceId string) (bool, error) {
					instanceDetails, err := getAndUnmarshalInstanceDetails(instanceId)
					return instanceDetails["bucket_name"] != "", err
				},
				cleanupFn: func() {
					bucket := service.Bucket(instance_name)
					bucket.Delete(context.Background())
				},
			}

			testGenericService(brokerConfig, gcpBroker, params)
		}, timeout)
	})

	Describe("pub sub", func() {
		It("can provision/bind/unbind/deprovision", func() {
			co := option.WithUserAgent(models.CustomUserAgent)
			ct := option.WithTokenSource(brokerConfig.HttpConfig.TokenSource(context.Background()))
			service, err := googlepubsub.NewClient(context.Background(), brokerConfig.ProjectId, co, ct)
			Expect(err).NotTo(HaveOccurred())

			topic := service.Topic(instance_name)

			params := &genericService{
				serviceId:        serviceNameToId[models.PubsubName],
				planId:           serviceNameToPlanId[models.PubsubName],
				instanceId:       "integration_test_topic",
				bindingId:        "integration_test_topic_bindingId",
				rawBindingParams: []byte(`{"role": "pubsub.editor"}`),
				serviceExistsFn: func(bool) bool {
					exists, err := topic.Exists(context.Background())
					return exists && err == nil
				},
				serviceMetadataSavedFn: func(instanceId string) (bool, error) {
					instanceDetails, err := getAndUnmarshalInstanceDetails(instanceId)
					return instanceDetails["topic_name"] != "", err
				},
				cleanupFn: func() {
					err := topic.Delete(context.Background())
					Expect(err).NotTo(HaveOccurred())
				},
			}

			testGenericService(brokerConfig, gcpBroker, params)
		}, timeout)
	})

	Describe("stackdriver debugger", func() {
		It("can provision/bind/unbind/deprovision", func() {
			params := &iamService{
				serviceId: serviceNameToId[models.StackdriverDebuggerName],
				planId:    serviceNameToPlanId[models.StackdriverDebuggerName],
			}
			testIamBasedService(brokerConfig, gcpBroker, params)
		}, timeout)
	})

	Describe("stackdriver trace", func() {
		It("can provision/bind/unbind/deprovision", func() {
			params := &iamService{
				serviceId: serviceNameToId[models.StackdriverTraceName],
				planId:    serviceNameToPlanId[models.StackdriverTraceName],
			}
			testIamBasedService(brokerConfig, gcpBroker, params)
		}, timeout)
	})

	Describe("datastore", func() {
		It("can provision/bind/unbind/deprovision", func() {
			params := &iamService{
				serviceId: serviceNameToId[models.DatastoreName],
				planId:    serviceNameToPlanId[models.DatastoreName],
			}
			testIamBasedService(brokerConfig, gcpBroker, params)
		}, timeout)
	})

	AfterEach(func() {
		os.Remove("test.db")
	})
})
