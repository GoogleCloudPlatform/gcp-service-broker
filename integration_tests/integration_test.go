package integration_tests

import (
	. "gcp-service-broker/brokerapi/brokers"
	"golang.org/x/net/context"

	"fmt"
	"gcp-service-broker/brokerapi/brokers"
	"gcp-service-broker/brokerapi/brokers/models"
	"gcp-service-broker/brokerapi/brokers/name_generator"
	"gcp-service-broker/db_service"
	"gcp-service-broker/fakes"
	"hash/crc32"
	"net/http"
	"os"

	googlepubsub "cloud.google.com/go/pubsub"

	googlestorage "cloud.google.com/go/storage"
	"code.cloudfoundry.org/lager"
	"encoding/json"
	"github.com/jinzhu/gorm"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	googlebigquery "google.golang.org/api/bigquery/v2"
	"google.golang.org/api/iam/v1"
	"google.golang.org/api/option"
	"strconv"
	"time"
)

const timeout = 60

type genericService struct {
	serviceId              string
	planId                 string
	bindingId              string
	rawBindingParams       map[string]interface{}
	instanceId             string
	serviceExistsFn        func(bool) bool
	cleanupFn              func()
	serviceMetadataSavedFn func(string) bool
}

func getAndUnmarshalInstanceDetails(instanceId string) map[string]string {
	var instanceRecord models.ServiceInstanceDetails
	db_service.DbConnection.Find(&instanceRecord).Where("id = ?", instanceId)
	var instanceDetails map[string]string
	json.Unmarshal([]byte(instanceRecord.OtherDetails), &instanceDetails)
	return instanceDetails
}

func testGenericService(gcpBroker *GCPAsyncServiceBroker, params *genericService) {
	// If the service already exists (eg, failed previous test), clean it up before the run
	if params.serviceExistsFn(false) {
		params.cleanupFn()
	}
	//
	// Provision
	//
	provisionDetails := models.ProvisionDetails{
		ServiceID: params.serviceId,
		PlanID:    params.planId,
	}

	_, err := gcpBroker.Provision(params.instanceId, provisionDetails, true)
	Expect(err).ToNot(HaveOccurred())

	// Provision is registered in the database
	var count int
	db_service.DbConnection.Model(&models.ServiceInstanceDetails{}).Where("id = ?", params.instanceId).Count(&count)
	Expect(count).To(Equal(1))

	Expect(params.serviceExistsFn(true)).To(BeTrue())
	Expect(params.serviceMetadataSavedFn(params.instanceId)).To(BeTrue())

	//
	// Bind
	//
	bindDetails := models.BindDetails{
		ServiceID:  params.serviceId,
		PlanID:     params.planId,
		Parameters: params.rawBindingParams,
	}
	creds, err := gcpBroker.Bind(params.instanceId, params.bindingId, bindDetails)
	Expect(err).ToNot(HaveOccurred())

	db_service.DbConnection.Model(&models.ServiceBindingCredentials{}).Where("binding_id = ?", params.bindingId).Count(&count)
	Expect(count).To(Equal(1))

	iamService, err := iam.New(gcpBroker.GCPClient)
	Expect(err).ToNot(HaveOccurred())
	saService := iam.NewProjectsServiceAccountsService(iamService)
	resourceName := "projects/" + gcpBroker.RootGCPCredentials.ProjectId + "/serviceAccounts/" + creds.Credentials.(map[string]string)["UniqueId"]
	_, err = saService.Get(resourceName).Do()
	Expect(err).ToNot(HaveOccurred())

	//
	// Unbind
	//
	unbindDetails := models.UnbindDetails{
		ServiceID: params.serviceId,
		PlanID:    params.planId,
	}
	err = gcpBroker.Unbind(params.instanceId, params.bindingId, unbindDetails)
	Expect(err).ToNot(HaveOccurred())

	binding := models.ServiceBindingCredentials{}
	if err := db_service.DbConnection.Unscoped().Where("binding_id = ?", params.bindingId).First(&binding).Error; err != nil {
		panic("error checking for binding details: " + err.Error())
	}
	Expect(binding.DeletedAt).NotTo(Equal(nil))

	_, err = saService.Get(resourceName).Do()
	Expect(err).To(HaveOccurred())

	//
	// Deprovision
	//
	deprovisionDetails := models.DeprovisionDetails{
		ServiceID: params.serviceId,
		PlanID:    params.planId,
	}
	_, err = gcpBroker.Deprovision(params.instanceId, deprovisionDetails, true)
	Expect(err).ToNot(HaveOccurred())
	instance := models.ServiceInstanceDetails{}
	if err := db_service.DbConnection.Unscoped().Where("ID = ?", params.instanceId).First(&instance).Error; err != nil {
		panic("error checking for service instance details: " + err.Error())
	}
	Expect(instance.DeletedAt).NotTo(Equal(nil))

	Expect(params.serviceExistsFn(false)).To(BeFalse())
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
func generateInstanceName(projectId string) string {
	hashed := crc32.ChecksumIEEE([]byte(projectId))
	return fmt.Sprintf("pcf_sb_1_%d", hashed)
}

var _ = Describe("LiveIntegrationTests", func() {
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
		os.Setenv("SERVICES", `[
			{
			  "id": "b9e4332e-b42b-4680-bda5-ea1506797474",
			  "description": "A Powerful, Simple and Cost Effective Object Storage Service",
			  "name": "google-storage",
			  "bindable": true,
			  "plan_updateable": false,
			  "metadata": {
			    "displayName": "Google Cloud Storage",
			    "longDescription": "A Powerful, Simple and Cost Effective Object Storage Service",
			    "documentationUrl": "https://cloud.google.com/storage/docs/overview",
			    "supportUrl": "https://cloud.google.com/support/"
			  },
			  "tags": ["gcp", "storage"]
			},
			{
			  "id": "628629e3-79f5-4255-b981-d14c6c7856be",
			  "description": "A global service for real-time and reliable messaging and streaming data",
			  "name": "google-pubsub",
			  "bindable": true,
			  "plan_updateable": false,
			  "metadata": {
			    "displayName": "Google PubSub",
			    "longDescription": "A global service for real-time and reliable messaging and streaming data",
			    "documentationUrl": "https://cloud.google.com/pubsub/docs/",
			    "supportUrl": "https://cloud.google.com/support/"
			  },
			  "tags": ["gcp", "pubsub"]
			},
			{
			  "id": "f80c0a3e-bd4d-4809-a900-b4e33a6450f1",
			  "description": "A fast, economical and fully managed data warehouse for large-scale data analytics",
			  "name": "google-bigquery",
			  "bindable": true,
			  "plan_updateable": false,
			  "metadata": {
			    "displayName": "Google BigQuery",
			    "longDescription": "A fast, economical and fully managed data warehouse for large-scale data analytics",
			    "documentationUrl": "https://cloud.google.com/bigquery/docs/",
			    "supportUrl": "https://cloud.google.com/support/"
			  },
			  "tags": ["gcp", "bigquery"]
			},
			{
			  "id": "4bc59b9a-8520-409f-85da-1c7552315863",
			  "description": "Google Cloud SQL is a fully-managed MySQL database service",
			  "name": "google-cloudsql",
			  "bindable": true,
			  "plan_updateable": false,
			  "metadata": {
			    "displayName": "Google CloudSQL",
			    "longDescription": "Google Cloud SQL is a fully-managed MySQL database service",
			    "documentationUrl": "https://cloud.google.com/sql/docs/",
			    "supportUrl": "https://cloud.google.com/support/"
			  },
			  "tags": ["gcp", "cloudsql"]
			},
			{
			  "id": "5ad2dce0-51f7-4ede-8b46-293d6df1e8d4",
			  "description": "Machine Learning Apis including Vision, Translate, Speech, and Natural Language",
			  "name": "google-ml-apis",
			  "bindable": true,
			  "plan_updateable": false,
			  "metadata": {
			    "displayName": "Google Machine Learning APIs",
			    "longDescription": "Machine Learning Apis including Vision, Translate, Speech, and Natural Language",
			    "documentationUrl": "https://cloud.google.com/ml/",
			    "supportUrl": "https://cloud.google.com/support/"
			  },
			  "tags": ["gcp", "ml"]
			}
		      ]`)
		os.Setenv("PRECONFIGURED_PLANS", `[
			{
			  "service_id": "b9e4332e-b42b-4680-bda5-ea1506797474",
			  "name": "standard",
			  "display_name": "Standard",
			  "description": "Standard storage class",
			  "features": {"storage_class": "STANDARD"}
			},
			{
			  "service_id": "b9e4332e-b42b-4680-bda5-ea1506797474",
			  "name": "nearline",
			  "display_name": "Nearline",
			  "description": "Nearline storage class",
			  "features": {"storage_class": "NEARLINE"}
			},
			{
			  "service_id": "b9e4332e-b42b-4680-bda5-ea1506797474",
			  "name": "reduced_availability",
			  "display_name": "Durable Reduced Availability",
			  "description": "Durable Reduced Availability storage class",
			  "features": {"storage_class": "DURABLE_REDUCED_AVAILABILITY"}
			},
			{
			  "service_id": "628629e3-79f5-4255-b981-d14c6c7856be",
			  "name": "default",
			  "display_name": "Default",
			  "description": "PubSub Default plan",
			  "features": ""
			},
			{ "service_id": "f80c0a3e-bd4d-4809-a900-b4e33a6450f1",
			  "name": "default",
			  "display_name": "Default",
			  "description": "BigQuery default plan",
			  "features": ""
			},
			{
			  "service_id": "5ad2dce0-51f7-4ede-8b46-293d6df1e8d4",
			  "name": "default",
			  "display_name": "Default",
			  "description": "Machine Learning api default plan",
			  "features": ""
			}
		      ]`)

		os.Setenv("CLOUDSQL_CUSTOM_PLANS", `{
			"test_plan": {
				"guid": "foo",
				"name": "bar",
				"description": "testplan",
				"tier": "D4",
				"pricing_plan": "PER_USE",
				"max_disk_size": "20",
				"display_name": "FOOBAR",
				"service": "4bc59b9a-8520-409f-85da-1c7552315863"
			}
		}`)

		var creds models.GCPCredentials
		creds, err = brokers.InitCredentialsFromEnv()
		if err != nil {
			logger.Error("error", err)
		}
		instance_name = generateInstanceName(creds.ProjectId)
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

	Describe("Broker init", func() {
		It("should have 5 services in sevices map", func() {
			Expect(len(gcpBroker.ServiceBrokerMap)).To(Equal(5))
		})

		It("should have a default client", func() {
			Expect(gcpBroker.GCPClient).NotTo(Equal(&http.Client{}))
		})

		It("should have loaded credentials correctly and have a project id", func() {
			Expect(gcpBroker.RootGCPCredentials.ProjectId).To(Not(BeEmpty()))
		})
	})

	Describe("getting broker catalog", func() {
		It("should have 5 services available", func() {
			Expect(len(gcpBroker.Services())).To(Equal(5))
		})

		It("should have 3 storage plans available", func() {
			serviceList := gcpBroker.Services()
			for _, s := range serviceList {
				if s.ID == serviceNameToId[StorageName] {
					Expect(len(s.Plans)).To(Equal(3))
				}
			}

		})
	})

	Describe("bigquery", func() {
		It("can provision/bind/unbind/deprovision", func() {
			service, err := googlebigquery.New(gcpBroker.GCPClient)
			Expect(err).NotTo(HaveOccurred())

			params := &genericService{
				serviceId:  serviceNameToId[brokers.BigqueryName],
				planId:     serviceNameToPlanId[brokers.BigqueryName],
				bindingId:  "integration_test_bind",
				instanceId: "integration_test_dataset",
				rawBindingParams: map[string]interface{}{
					"role": "bigquery.admin",
				},
				serviceExistsFn: func(expected bool) bool {
					_, err = service.Datasets.Get(gcpBroker.RootGCPCredentials.ProjectId, instance_name).Do()

					return err == nil
				},
				serviceMetadataSavedFn: func(instanceId string) bool {
					instanceDetails := getAndUnmarshalInstanceDetails(instanceId)
					return instanceDetails["name"] != ""
				},
				cleanupFn: func() {
					err := service.Datasets.Delete(gcpBroker.RootGCPCredentials.ProjectId, instance_name).Do()
					Expect(err).NotTo(HaveOccurred())
				},
			}
			testGenericService(gcpBroker, params)
		}, timeout)
	})

	//Describe("deprovision", func() {
	//	Context("when the bigquery service id is provided", func() {
	//		It("should call bigquery deprovisioning", func() {
	//			bqId := serviceNameToId[brokers.BigqueryName]
	//			_, err := gcpBroker.Provision("something", bqProvisionDetails, true)
	//			Expect(err).NotTo(HaveOccurred())
	//			_, err = gcpBroker.Deprovision("something", models.DeprovisionDetails{
	//				ServiceID: bqId,
	//			}, true)
	//			Expect(err).NotTo(HaveOccurred())
	//			Expect(gcpBroker.ServiceBrokerMap[bqId].(*modelsfakes.FakeServiceBrokerHelper).DeprovisionCallCount()).To(Equal(1))
	//		})
	//	})
	//
	//	Context("when the service doesn't exist", func() {
	//		It("should return an error", func() {
	//			_, err := gcpBroker.Deprovision("something", models.DeprovisionDetails{
	//				ServiceID: serviceNameToId[brokers.BigqueryName],
	//			}, true)
	//			Expect(err).To(HaveOccurred())
	//		})
	//	})
	//
	//	Context("when async provisioning isn't allowed but the service requested requires it", func() {
	//		It("should return an error", func() {
	//			_, err := gcpBroker.Deprovision("something", models.DeprovisionDetails{
	//				ServiceID: serviceNameToId[brokers.CloudsqlName],
	//			}, false)
	//			Expect(err).To(HaveOccurred())
	//		})
	//	})
	//})

	//Describe("bind", func() {
	//	Context("when bind is called on storage", func() {
	//		It("it should call storage bind", func() {
	//			_, err = gcpBroker.Provision("storagething", bqProvisionDetails, true)
	//			Expect(err).NotTo(HaveOccurred())
	//			_, err = gcpBroker.Bind(instanceId, "newbinding", storageBindDetails)
	//			Expect(err).NotTo(HaveOccurred())
	//			Expect(gcpBroker.ServiceBrokerMap[serviceNameToId[StorageName]].(*modelsfakes.FakeServiceBrokerHelper).BindCallCount()).To(Equal(1))
	//		})
	//	})
	//
	//	Context("when bind is called more than once on the same id", func() {
	//		It("it should throw an error", func() {
	//			_, err = gcpBroker.Provision("storagething", bqProvisionDetails, true)
	//			Expect(err).NotTo(HaveOccurred())
	//			_, err = gcpBroker.Bind(instanceId, bindingId, storageBindDetails)
	//			Expect(err).NotTo(HaveOccurred())
	//			_, err = gcpBroker.Bind(instanceId, bindingId, storageBindDetails)
	//			Expect(err).To(HaveOccurred())
	//		})
	//	})
	//})
	//
	//Describe("unbind", func() {
	//	Context("when unbind is called on storage", func() {
	//		It("it should call storage unbind", func() {
	//			_, err = gcpBroker.Provision("storagething", bqProvisionDetails, true)
	//			Expect(err).NotTo(HaveOccurred())
	//			_, err = gcpBroker.Bind(instanceId, bindingId, storageBindDetails)
	//			Expect(err).NotTo(HaveOccurred())
	//			err = gcpBroker.Unbind(instanceId, bindingId, storageUnbindDetails)
	//			Expect(err).NotTo(HaveOccurred())
	//			Expect(gcpBroker.ServiceBrokerMap[serviceNameToId[StorageName]].(*modelsfakes.FakeServiceBrokerHelper).UnbindCallCount()).To(Equal(1))
	//		})
	//	})
	//
	//	Context("when unbind is called more than once on the same id", func() {
	//		It("it should throw an error", func() {
	//			_, err = gcpBroker.Provision("storagething", bqProvisionDetails, true)
	//			Expect(err).NotTo(HaveOccurred())
	//			_, err = gcpBroker.Bind(instanceId, bindingId, storageBindDetails)
	//			Expect(err).NotTo(HaveOccurred())
	//			err = gcpBroker.Unbind(instanceId, bindingId, storageUnbindDetails)
	//			Expect(err).NotTo(HaveOccurred())
	//			err = gcpBroker.Unbind(instanceId, bindingId, storageUnbindDetails)
	//			Expect(err).To(HaveOccurred())
	//		})
	//	})
	//})
	//
	//Describe("lastOperation", func() {
	//	Context("when last operation is called on a service that doesn't exist", func() {
	//		It("should throw an error", func() {
	//			_, err = gcpBroker.LastOperation("somethingnonexistant")
	//			Expect(err).To(HaveOccurred())
	//		})
	//	})
	//
	//	Context("when last operation is called on a service that is provisioned synchronously", func() {
	//		It("should throw an error", func() {
	//			_, err = gcpBroker.Provision(instanceId, bqProvisionDetails, true)
	//			Expect(err).NotTo(HaveOccurred())
	//			_, err = gcpBroker.LastOperation(instanceId)
	//			Expect(err).To(HaveOccurred())
	//		})
	//	})
	//
	//	Context("when last operation is called on an asynchronous service", func() {
	//		It("should call PollInstance", func() {
	//			_, err = gcpBroker.Provision(instanceId, cloudSqlProvisionDetails, true)
	//			Expect(err).NotTo(HaveOccurred())
	//			_, err = gcpBroker.LastOperation(instanceId)
	//			Expect(err).NotTo(HaveOccurred())
	//			Expect(gcpBroker.ServiceBrokerMap[serviceNameToId[CloudsqlName]].(*modelsfakes.FakeServiceBrokerHelper).PollInstanceCallCount()).To(Equal(1))
	//		})
	//	})

	//})

	Describe("cloud storage", func() {
		It("can provision/bind/unbind/deprovision", func() {
			service, err := googlestorage.NewClient(context.Background(), option.WithUserAgent(models.CustomUserAgent))
			Expect(err).NotTo(HaveOccurred())

			params := &genericService{
				serviceId:  serviceNameToId[brokers.StorageName],
				planId:     serviceNameToPlanId[brokers.StorageName],
				instanceId: "integration_test_bucket",
				bindingId:  "integration_test_bucket_binding",
				rawBindingParams: map[string]interface{}{
					"role": "storage.admin",
				},
				serviceExistsFn: func(bool) bool {
					bucket := service.Bucket(instance_name)
					_, err = bucket.Attrs(context.Background())

					return err == nil
				},
				serviceMetadataSavedFn: func(instanceId string) bool {
					instanceDetails := getAndUnmarshalInstanceDetails(instanceId)
					return instanceDetails["name"] != ""
				},
				cleanupFn: func() {
					bucket := service.Bucket(instance_name)
					bucket.Delete(context.Background())
				},
			}

			testGenericService(gcpBroker, params)
		}, timeout)
	})

	Describe("pub sub", func() {
		It("can provision/bind/unbind/deprovision", func() {
			service, err := googlepubsub.NewClient(context.Background(), gcpBroker.RootGCPCredentials.ProjectId, option.WithUserAgent(models.CustomUserAgent))
			Expect(err).NotTo(HaveOccurred())

			topic := service.Topic(instance_name)

			params := &genericService{
				serviceId:  serviceNameToId[brokers.PubsubName],
				planId:     serviceNameToPlanId[brokers.PubsubName],
				instanceId: "integration_test_topic",
				bindingId:  "integration_test_topic_bindingId",
				rawBindingParams: map[string]interface{}{
					"role": "pubsub.admin",
				},
				serviceExistsFn: func(bool) bool {
					exists, err := topic.Exists(context.Background())
					return exists && err == nil
				},
				serviceMetadataSavedFn: func(instanceId string) bool {
					instanceDetails := getAndUnmarshalInstanceDetails(instanceId)
					fmt.Printf("%v", instanceDetails)
					return instanceDetails["topic_name"] != ""
				},
				cleanupFn: func() {
					err := topic.Delete(context.Background())
					Expect(err).NotTo(HaveOccurred())
				},
			}

			testGenericService(gcpBroker, params)
		}, timeout)
	})

	AfterEach(func() {
		os.Remove(brokers.AppCredsFileName)
		os.Remove("test.db")
	})
})
