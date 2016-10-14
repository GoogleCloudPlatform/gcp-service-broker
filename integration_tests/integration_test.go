package integration_tests

import (
	. "gcp-service-broker/brokerapi/brokers"

	"code.cloudfoundry.org/lager"
	"gcp-service-broker/brokerapi/brokers"
	"gcp-service-broker/brokerapi/brokers/models"
	"gcp-service-broker/db_service"
	"github.com/jinzhu/gorm"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	googlebigquery "google.golang.org/api/bigquery/v2"
	"net/http"
	"os"
)

var _ = Describe("LiveIntegrationTests", func() {
	var (
		gcpBroker                *GCPAsyncServiceBroker
		err                      error
		logger                   lager.Logger
		serviceNameToId          map[string]string = make(map[string]string)
		someBigQueryPlanId       string
		cloudSqlProvisionDetails models.ProvisionDetails
		storageBindDetails       models.BindDetails
		storageUnbindDetails     models.UnbindDetails
		instanceId               string
		bindingId                string
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

		instanceId = "newid"
		bindingId = "newbinding"

		gcpBroker, err = brokers.New(logger)
		println("inited the broker!")
		if err != nil {
			logger.Error("error", err)
		}

		var someCloudSQLPlanId string
		var someStoragePlanId string
		for _, service := range *gcpBroker.Catalog {
			serviceNameToId[service.Name] = service.ID
			if service.Name == BigqueryName {
				someBigQueryPlanId = service.Plans[0].ID
			}
			if service.Name == CloudsqlName {
				someCloudSQLPlanId = service.Plans[0].ID
			}
			if service.Name == StorageName {
				someStoragePlanId = service.Plans[0].ID
			}
		}

		cloudSqlProvisionDetails = models.ProvisionDetails{
			ServiceID: serviceNameToId[brokers.CloudsqlName],
			PlanID:    someCloudSQLPlanId,
		}

		storageBindDetails = models.BindDetails{
			ServiceID: serviceNameToId[brokers.StorageName],
			PlanID:    someStoragePlanId,
		}

		storageUnbindDetails = models.UnbindDetails{
			ServiceID: serviceNameToId[brokers.StorageName],
			PlanID:    someStoragePlanId,
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
			Expect(gcpBroker.RootGCPCredentials.ProjectId).To(Equal("gcp-service-broker-testing"))
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

		var (
			bqProvisionDetails   models.ProvisionDetails
			bqDeprovisionDetails models.DeprovisionDetails
			service              *googlebigquery.Service
			datasetName          string
		)

		BeforeEach(func() {
			datasetName = "integration_test_dataset"

			bqProvisionDetails = models.ProvisionDetails{
				ServiceID:     serviceNameToId[brokers.BigqueryName],
				PlanID:        someBigQueryPlanId,
				RawParameters: []byte("{\"name\": \"integration_test_dataset\"}"),
			}

			bqDeprovisionDetails = models.DeprovisionDetails{
				ServiceID: serviceNameToId[brokers.BigqueryName],
				PlanID:    someBigQueryPlanId,
			}

			service, err = googlebigquery.New(gcpBroker.GCPClient)
			if err != nil {
				panic("error creating bigquery client for testing")
			}
		})

		Context("bigquery provision and deprovision", func() {

			It("should make a bigquery dataset on provision and delete it on deprovision, and maintain db records", func() {
				_, err := gcpBroker.Provision(instanceId, bqProvisionDetails, true)
				Expect(err).ToNot(HaveOccurred())
				_, err = service.Datasets.Get(gcpBroker.RootGCPCredentials.ProjectId, datasetName).Do()
				Expect(err).ToNot(HaveOccurred())

				var count int
				db_service.DbConnection.Model(&models.ServiceInstanceDetails{}).Where("id = ?", instanceId).Count(&count)
				Expect(count).To(Equal(1))

				_, err = gcpBroker.Deprovision(instanceId, bqDeprovisionDetails, true)
				Expect(err).ToNot(HaveOccurred())
				_, err = service.Datasets.Get(gcpBroker.RootGCPCredentials.ProjectId, datasetName).Do()
				Expect(err).To(HaveOccurred())

			})

		})

		//Context("when too many services are provisioned", func() {
		//	It("should return an error", func() {
		//		gcpBroker.InstanceLimit = 0
		//		_, err := gcpBroker.Provision("something", bqProvisionDetails, true)
		//		Expect(err).To(HaveOccurred())
		//		Expect(err).To(Equal(models.ErrInstanceLimitMet))
		//	})
		//})
		//
		//Context("when an unrecognized service is provisioned", func() {
		//	It("should return an error", func() {
		//		_, err = gcpBroker.Provision("something", models.ProvisionDetails{
		//			ServiceID: "nope",
		//			PlanID:    "nope",
		//		}, true)
		//		Expect(err).To(HaveOccurred())
		//	})
		//})
		//
		//Context("when an unrecognized plan is provisioned", func() {
		//	It("should return an error", func() {
		//		_, err = gcpBroker.Provision("something", models.ProvisionDetails{
		//			ServiceID: serviceNameToId[BigqueryName],
		//			PlanID:    "nope",
		//		}, true)
		//		Expect(err).To(HaveOccurred())
		//	})
		//})
		//
		//Context("when duplicate services are provisioned", func() {
		//	It("should return an error", func() {
		//		_, err = gcpBroker.Provision("something", bqProvisionDetails, true)
		//		Expect(err).NotTo(HaveOccurred())
		//		_, err := gcpBroker.Provision("something", bqProvisionDetails, true)
		//		Expect(err).To(HaveOccurred())
		//	})
		//})
		//
		//Context("when async provisioning isn't allowed but the service requested requires it", func() {
		//	It("should return an error", func() {
		//		_, err := gcpBroker.Provision("something", cloudSqlProvisionDetails, false)
		//		Expect(err).To(HaveOccurred())
		//	})
		//})

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

	AfterEach(func() {
		os.Remove(brokers.AppCredsFileName)
		os.Remove("test.db")
	})
})
