package brokers_test

import (
	"encoding/json"
	"gcp-service-broker/brokerapi/brokers"
	. "gcp-service-broker/brokerapi/brokers"
	"gcp-service-broker/brokerapi/brokers/broker_base"
	"gcp-service-broker/brokerapi/brokers/cloudsql"
	"gcp-service-broker/brokerapi/brokers/models"
	"gcp-service-broker/brokerapi/brokers/models/modelsfakes"
	"gcp-service-broker/brokerapi/brokers/name_generator"
	"gcp-service-broker/brokerapi/brokers/pubsub"
	"gcp-service-broker/brokerapi/brokers/spanner"
	"gcp-service-broker/db_service"
	"net/http"
	"os"

	"code.cloudfoundry.org/lager"

	"gcp-service-broker/fakes"

	"github.com/jinzhu/gorm"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Brokers", func() {
	var (
		gcpBroker                *GCPAsyncServiceBroker
		err                      error
		logger                   lager.Logger
		serviceNameToId          map[string]string = make(map[string]string)
		bqProvisionDetails       models.ProvisionDetails
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
		name_generator.New()

		os.Setenv("ROOT_SERVICE_ACCOUNT_JSON", `{
			"type": "service_account",
			"project_id": "foo",
			"private_key_id": "something",
			"private_key": "foobar",
			"client_email": "example@gmail.com",
			"client_id": "1",
			"auth_uri": "somelink",
			"token_uri": "somelink",
			"auth_provider_x509_cert_url": "somelink",
			"client_x509_cert_url": "somelink"
		      }`)
		os.Setenv("SECURITY_USER_NAME", "username")
		os.Setenv("SECURITY_USER_PASSWORD", "password")
		os.Setenv("SERVICES", fakes.Services)
		os.Setenv("PRECONFIGURED_PLANS", fakes.PreconfiguredPlans)

		os.Setenv("CLOUDSQL_CUSTOM_PLANS", fakes.TestCloudSQLPlan)
		os.Setenv("BIGTABLE_CUSTOM_PLANS", fakes.TestBigtablePlan)
		os.Setenv("SPANNER_CUSTOM_PLANS", fakes.TestSpannerPlan)

		instanceId = "newid"
		bindingId = "newbinding"

		gcpBroker, err = brokers.New(logger)
		if err != nil {
			logger.Error("error", err)
		}

		var someBigQueryPlanId string
		var someCloudSQLPlanId string
		var someStoragePlanId string
		for _, service := range *gcpBroker.Catalog {
			serviceNameToId[service.Name] = service.ID
			if service.Name == models.BigqueryName {
				someBigQueryPlanId = service.Plans[0].ID
			}
			if service.Name == models.CloudsqlName {

				someCloudSQLPlanId = service.Plans[0].ID
			}
			if service.Name == models.StorageName {
				someStoragePlanId = service.Plans[0].ID
			}
		}

		for k, _ := range gcpBroker.ServiceBrokerMap {
			async := false
			if k == serviceNameToId[models.CloudsqlName] {
				async = true
			}
			gcpBroker.ServiceBrokerMap[k] = &modelsfakes.FakeServiceBrokerHelper{
				ProvisionsAsyncStub:   func() bool { return async },
				DeprovisionsAsyncStub: func() bool { return async },
				ProvisionStub: func(instanceId string, details models.ProvisionDetails, plan models.PlanDetails) (models.ServiceInstanceDetails, error) {
					return models.ServiceInstanceDetails{ID: instanceId, OtherDetails: "{\"mynameis\": \"instancename\"}"}, nil
				},
				BindStub: func(instanceID, bindingID string, details models.BindDetails) (models.ServiceBindingCredentials, error) {
					return models.ServiceBindingCredentials{OtherDetails: "{\"foo\": \"bar\"}"}, nil
				},
			}
		}

		bqProvisionDetails = models.ProvisionDetails{
			ServiceID: serviceNameToId[models.BigqueryName],
			PlanID:    someBigQueryPlanId,
		}

		cloudSqlProvisionDetails = models.ProvisionDetails{
			ServiceID: serviceNameToId[models.CloudsqlName],
			PlanID:    someCloudSQLPlanId,
		}

		storageBindDetails = models.BindDetails{
			ServiceID: serviceNameToId[models.StorageName],
			PlanID:    someStoragePlanId,
		}

		storageUnbindDetails = models.UnbindDetails{
			ServiceID: serviceNameToId[models.StorageName],
			PlanID:    someStoragePlanId,
		}

	})

	Describe("Broker init", func() {
		It("should have 9 services in sevices map", func() {
			Expect(len(gcpBroker.ServiceBrokerMap)).To(Equal(9))
		})

		It("should have a default client", func() {
			Expect(gcpBroker.GCPClient).NotTo(Equal(&http.Client{}))
		})

		It("should have loaded credentials correctly and have a project id", func() {
			Expect(gcpBroker.RootGCPCredentials.ProjectId).To(Equal("foo"))
		})
	})

	Describe("getting broker catalog", func() {
		It("should have 9 services available", func() {
			Expect(len(gcpBroker.Services())).To(Equal(9))
		})

		It("should have 3 storage plans available", func() {
			serviceList := gcpBroker.Services()
			for _, s := range serviceList {
				if s.ID == serviceNameToId[models.StorageName] {
					Expect(len(s.Plans)).To(Equal(3))
				}
			}

		})

		It("should have 1 cloudsql plan available", func() {
			serviceList := gcpBroker.Services()
			for _, s := range serviceList {
				if s.ID == serviceNameToId[models.CloudsqlName] {
					Expect(len(s.Plans)).To(Equal(1))
				}
			}

		})

		It("should have 1 bigtable plan available", func() {
			serviceList := gcpBroker.Services()
			for _, s := range serviceList {
				if s.ID == serviceNameToId[models.BigtableName] {
					Expect(len(s.Plans)).To(Equal(1))
				}
			}

		})

		It("should have 1 debugger plan available", func() {
			serviceList := gcpBroker.Services()
			for _, s := range serviceList {
				if s.ID == serviceNameToId[models.StackdriverDebuggerName] {
					Expect(len(s.Plans)).To(Equal(1))
				}
			}
		})

		It("should error if plan ids are not supplied", func() {
			os.Setenv("PRECONFIGURED_PLANS", fakes.PlanNoId)
			_, err := brokers.New(logger)
			Expect(err).To(HaveOccurred())
			os.Setenv("PRECONFIGURED_PLANS", fakes.PreconfiguredPlans)
		})
	})

	Describe("updating broker catalog", func() {

		It("should update cloudsql custom plans with different names on startup", func() {

			os.Setenv("CLOUDSQL_CUSTOM_PLANS", `{
				"newPlan": {
					"id": "some-other-cloudsql-plan",
					"name": "newPlan",
					"description": "testplan",
					"tier": "D8",
					"pricing_plan": "athing",
					"max_disk_size": "15",
					"display_name": "FOOBAR",
					"service": "4bc59b9a-8520-409f-85da-1c7552315863"
				}
			}`)

			newBroker, err := brokers.New(logger)

			serviceList := newBroker.Services()
			for _, s := range serviceList {
				if s.ID == serviceNameToId[models.CloudsqlName] {
					Expect(s.Plans[0].Name).To(Equal("newPlan"))
					Expect(len(s.Plans)).To(Equal(1))
					plan := models.PlanDetails{}
					if err := db_service.DbConnection.Where("service_id = ?", "4bc59b9a-8520-409f-85da-1c7552315863").First(&plan).Error; err != nil {
						panic("The provided plan does not exist " + err.Error())
					}
					var planDetails map[string]string
					if err = json.Unmarshal([]byte(plan.Features), &planDetails); err != nil {
						panic("Error unmarshalling plan features: " + err.Error())
					}
					Expect(planDetails["tier"]).To(Equal("D8"))
					Expect(planDetails["max_disk_size"]).To(Equal("15"))
				}
			}

		})

		It("should update cloudsql custom plans with the same name on startup", func() {

			os.Setenv("CLOUDSQL_CUSTOM_PLANS", `{
				"test_plan": {
					"id": "some-other-cloudsql-plan",
					"name": "test_plan",
					"description": "testplan",
					"tier": "D8",
					"pricing_plan": "athing",
					"max_disk_size": "15",
					"display_name": "FOOBAR",
					"service": "4bc59b9a-8520-409f-85da-1c7552315863"
				}
			}`)

			newBroker, err := brokers.New(logger)

			serviceList := newBroker.Services()
			for _, s := range serviceList {
				if s.ID == serviceNameToId[models.CloudsqlName] {
					Expect(len(s.Plans)).To(Equal(1))
					plan := models.PlanDetails{}
					if err := db_service.DbConnection.Where("service_id = ?", "4bc59b9a-8520-409f-85da-1c7552315863").First(&plan).Error; err != nil {
						panic("The provided plan does not exist " + err.Error())
					}
					var planDetails map[string]string
					if err = json.Unmarshal([]byte(plan.Features), &planDetails); err != nil {
						panic("Error unmarshalling plan features: " + err.Error())
					}
					Expect(planDetails["tier"]).To(Equal("D8"))
				}
			}

		})
	})

	Describe("provision", func() {
		Context("when the bigquery service id is provided", func() {
			It("should call bigquery provisioning", func() {
				bqId := serviceNameToId[models.BigqueryName]
				_, err := gcpBroker.Provision(instanceId, bqProvisionDetails, true)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(gcpBroker.ServiceBrokerMap[bqId].(*modelsfakes.FakeServiceBrokerHelper).ProvisionCallCount()).To(Equal(1))
			})

		})

		Context("when too many services are provisioned", func() {
			It("should return an error", func() {
				gcpBroker.InstanceLimit = 0
				_, err := gcpBroker.Provision(instanceId, bqProvisionDetails, true)
				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(models.ErrInstanceLimitMet))
			})
		})

		Context("when an unrecognized service is provisioned", func() {
			It("should return an error", func() {
				_, err = gcpBroker.Provision(instanceId, models.ProvisionDetails{
					ServiceID: "nope",
					PlanID:    "nope",
				}, true)
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when an unrecognized plan is provisioned", func() {
			It("should return an error", func() {
				_, err = gcpBroker.Provision(instanceId, models.ProvisionDetails{
					ServiceID: serviceNameToId[models.BigqueryName],
					PlanID:    "nope",
				}, true)
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when duplicate services are provisioned", func() {
			It("should return an error", func() {
				_, err = gcpBroker.Provision(instanceId, bqProvisionDetails, true)
				Expect(err).NotTo(HaveOccurred())
				_, err := gcpBroker.Provision(instanceId, bqProvisionDetails, true)
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when async provisioning isn't allowed but the service requested requires it", func() {
			It("should return an error", func() {
				_, err := gcpBroker.Provision(instanceId, cloudSqlProvisionDetails, false)
				Expect(err).To(HaveOccurred())
			})
		})

	})

	Describe("deprovision", func() {
		Context("when the bigquery service id is provided", func() {
			It("should call bigquery deprovisioning", func() {
				bqId := serviceNameToId[models.BigqueryName]
				_, err := gcpBroker.Provision(instanceId, bqProvisionDetails, true)
				Expect(err).NotTo(HaveOccurred())
				_, err = gcpBroker.Deprovision(instanceId, models.DeprovisionDetails{
					ServiceID: bqId,
				}, true)
				Expect(err).NotTo(HaveOccurred())
				Expect(gcpBroker.ServiceBrokerMap[bqId].(*modelsfakes.FakeServiceBrokerHelper).DeprovisionCallCount()).To(Equal(1))
			})
		})

		Context("when the service doesn't exist", func() {
			It("should return an error", func() {
				_, err := gcpBroker.Deprovision(instanceId, models.DeprovisionDetails{
					ServiceID: serviceNameToId[models.BigqueryName],
				}, true)
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when async provisioning isn't allowed but the service requested requires it", func() {
			It("should return an error", func() {
				_, err := gcpBroker.Deprovision(instanceId, models.DeprovisionDetails{
					ServiceID: serviceNameToId[models.CloudsqlName],
				}, false)
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("bind", func() {
		Context("when bind is called on storage", func() {
			It("it should call storage bind", func() {
				_, err = gcpBroker.Provision(instanceId, bqProvisionDetails, true)
				Expect(err).NotTo(HaveOccurred())
				_, err = gcpBroker.Bind(instanceId, bindingId, storageBindDetails)
				Expect(err).NotTo(HaveOccurred())
				Expect(gcpBroker.ServiceBrokerMap[serviceNameToId[models.StorageName]].(*modelsfakes.FakeServiceBrokerHelper).BindCallCount()).To(Equal(1))
			})
		})

		Context("when bind is called more than once on the same id", func() {
			It("it should throw an error", func() {
				_, err = gcpBroker.Provision(instanceId, bqProvisionDetails, true)
				Expect(err).NotTo(HaveOccurred())
				_, err = gcpBroker.Bind(instanceId, bindingId, storageBindDetails)
				Expect(err).NotTo(HaveOccurred())
				_, err = gcpBroker.Bind(instanceId, bindingId, storageBindDetails)
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when bind is called", func() {
			It("it should update credentials with instance information", func() {
				_, err = gcpBroker.Provision(instanceId, bqProvisionDetails, true)
				Expect(err).NotTo(HaveOccurred())
				_, err := gcpBroker.Bind(instanceId, bindingId, storageBindDetails)
				Expect(err).NotTo(HaveOccurred())
				Expect(gcpBroker.ServiceBrokerMap[serviceNameToId[models.StorageName]].(*modelsfakes.FakeServiceBrokerHelper).BuildInstanceCredentialsCallCount()).To(Equal(1))
			})
		})

	})

	Describe("unbind", func() {
		Context("when unbind is called on storage", func() {
			It("it should call storage unbind", func() {
				_, err = gcpBroker.Provision(instanceId, bqProvisionDetails, true)
				Expect(err).NotTo(HaveOccurred())
				_, err = gcpBroker.Bind(instanceId, bindingId, storageBindDetails)
				Expect(err).NotTo(HaveOccurred())
				err = gcpBroker.Unbind(instanceId, bindingId, storageUnbindDetails)
				Expect(err).NotTo(HaveOccurred())
				Expect(gcpBroker.ServiceBrokerMap[serviceNameToId[models.StorageName]].(*modelsfakes.FakeServiceBrokerHelper).UnbindCallCount()).To(Equal(1))
			})
		})

		Context("when unbind is called more than once on the same id", func() {
			It("it should throw an error", func() {
				_, err = gcpBroker.Provision(instanceId, bqProvisionDetails, true)
				Expect(err).NotTo(HaveOccurred())
				_, err = gcpBroker.Bind(instanceId, bindingId, storageBindDetails)
				Expect(err).NotTo(HaveOccurred())
				err = gcpBroker.Unbind(instanceId, bindingId, storageUnbindDetails)
				Expect(err).NotTo(HaveOccurred())
				err = gcpBroker.Unbind(instanceId, bindingId, storageUnbindDetails)
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("lastOperation", func() {
		Context("when last operation is called on a service that doesn't exist", func() {
			It("should throw an error", func() {
				_, err = gcpBroker.LastOperation("somethingnonexistant")
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when last operation is called on a service that is provisioned synchronously", func() {
			It("should throw an error", func() {
				_, err = gcpBroker.Provision(instanceId, bqProvisionDetails, true)
				Expect(err).NotTo(HaveOccurred())
				_, err = gcpBroker.LastOperation(instanceId)
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when last operation is called on an asynchronous service", func() {
			It("should call PollInstance", func() {
				_, err = gcpBroker.Provision(instanceId, cloudSqlProvisionDetails, true)
				Expect(err).NotTo(HaveOccurred())
				_, err = gcpBroker.LastOperation(instanceId)
				Expect(err).NotTo(HaveOccurred())
				Expect(gcpBroker.ServiceBrokerMap[serviceNameToId[models.CloudsqlName]].(*modelsfakes.FakeServiceBrokerHelper).PollInstanceCallCount()).To(Equal(1))
			})
		})

	})

	AfterEach(func() {
		os.Remove(models.AppCredsFileName)
		os.Remove("test.db")
	})
})

var _ = Describe("AccountManagers", func() {

	var (
		logger         lager.Logger
		iamStyleBroker models.ServiceBrokerHelper
		spannerBroker  models.ServiceBrokerHelper
		cloudsqlBroker models.ServiceBrokerHelper
		accountManager modelsfakes.FakeAccountManager
		err            error
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
		name_generator.New()

		accountManager = modelsfakes.FakeAccountManager{}

		iamStyleBroker = &pubsub.PubSubBroker{
			Logger: logger,
			BrokerBase: broker_base.BrokerBase{
				AccountManager: &accountManager,
			},
		}

		cloudsqlBroker = &cloudsql.CloudSQLBroker{
			Logger:         logger,
			AccountManager: &accountManager,
		}

		spannerBroker = &spanner.SpannerBroker{
			Logger:         logger,
			AccountManager: &accountManager,
		}
	})

	Describe("bind", func() {
		Context("when bind is called on an iam-style broker", func() {
			It("should call the account manager create account in google method", func() {
				_, err = iamStyleBroker.Bind("foo", "bar", models.BindDetails{})
				Expect(err).NotTo(HaveOccurred())
				Expect(accountManager.CreateAccountInGoogleCallCount()).To(Equal(1))
			})
		})

		Context("when bind is called on a cloudsql broker after provision", func() {
			It("should call the account manager create account in google method", func() {
				db_service.DbConnection.Save(&models.ServiceInstanceDetails{ID: "foo"})
				_, err = iamStyleBroker.Bind("foo", "bar", models.BindDetails{})
				Expect(err).NotTo(HaveOccurred())
				Expect(accountManager.CreateAccountInGoogleCallCount()).To(Equal(1))
			})
		})

		Context("when bind is called on a cloudsql broker on a missing service instance", func() {
			It("should throw an error", func() {
				_, err = cloudsqlBroker.Bind("foo", "bar", models.BindDetails{})
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when bind is called on a cloudsql broker with no username/password after provision", func() {
			It("should return a generated username and password", func() {
				db_service.DbConnection.Create(&models.ServiceInstanceDetails{ID: "foo"})

				_, err := cloudsqlBroker.Bind("foo", "bar", models.BindDetails{})

				Expect(err).NotTo(HaveOccurred())
				Expect(accountManager.CreateAccountInGoogleCallCount()).To(Equal(1))
				_, _, details, _ := accountManager.CreateAccountInGoogleArgsForCall(0)
				Expect(details.Parameters).NotTo(BeEmpty())

				username, usernameOk := details.Parameters["username"].(string)
				password, passwordOk := details.Parameters["password"].(string)

				Expect(usernameOk).To(BeTrue())
				Expect(passwordOk).To(BeTrue())
				Expect(username).NotTo(BeEmpty())
				Expect(password).NotTo(BeEmpty())
			})
		})

		Context("when MergeCredentialsAndInstanceInfo is called on a broker", func() {
			It("should call MergeCredentialsAndInstanceInfo on the account manager", func() {
				_ = iamStyleBroker.BuildInstanceCredentials(make(map[string]string), make(map[string]string))
				Expect(accountManager.BuildInstanceCredentialsCallCount()).To(Equal(1))
			})
		})
	})

	Describe("unbind", func() {
		Context("when unbind is called on the broker", func() {
			It("it should call the account manager delete account from google method", func() {
				err = iamStyleBroker.Unbind(models.ServiceBindingCredentials{})
				Expect(err).NotTo(HaveOccurred())
				Expect(accountManager.DeleteAccountFromGoogleCallCount()).To(Equal(1))
			})
		})

		Context("when unbind is called on a cloudsql broker", func() {
			It("it should call the account manager delete account from google method", func() {
				err = cloudsqlBroker.Unbind(models.ServiceBindingCredentials{})
				Expect(err).NotTo(HaveOccurred())
				Expect(accountManager.DeleteAccountFromGoogleCallCount()).To(Equal(1))
			})
		})
	})

	Describe("async", func() {
		Context("with a pubsub broker", func() {
			It("should return false", func() {
				Expect(iamStyleBroker.ProvisionsAsync()).To(Equal(false))
				Expect(iamStyleBroker.DeprovisionsAsync()).To(Equal(false))
			})
		})

		Context("with a cloudsql broker", func() {
			It("should return true", func() {
				Expect(cloudsqlBroker.ProvisionsAsync()).To(Equal(true))
				Expect(cloudsqlBroker.DeprovisionsAsync()).To(Equal(true))
			})
		})

		Context("with a spanner broker", func() {
			It("should return true", func() {
				Expect(spannerBroker.ProvisionsAsync()).To(Equal(true))
				Expect(spannerBroker.DeprovisionsAsync()).To(Equal(false))
			})
		})
	})

	AfterEach(func() {
		os.Remove(models.AppCredsFileName)
		os.Remove("test.db")
	})

})
