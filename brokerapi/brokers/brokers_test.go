package brokers_test

import (
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
	"os"

	"code.cloudfoundry.org/lager"

	"gcp-service-broker/fakes"

	"gcp-service-broker/brokerapi/brokers/config"
	"github.com/jinzhu/gorm"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"golang.org/x/oauth2/jwt"
)

var _ = Describe("Brokers", func() {
	var (
		gcpBroker                *GCPAsyncServiceBroker
		brokerConfig             *config.BrokerConfig
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

		fakes.SetUpTestServices()

		brokerConfig, err = config.NewBrokerConfigFromEnv()
		if err != nil {
			logger.Error("error", err)
		}

		instanceId = "newid"
		bindingId = "newbinding"

		gcpBroker, err = brokers.New(brokerConfig, logger)
		if err != nil {
			logger.Error("error", err)
		}

		var someBigQueryPlanId string
		var someCloudSQLPlanId string
		var someStoragePlanId string
		for _, service := range gcpBroker.Catalog {
			serviceNameToId[service.Name] = service.ID
			if service.Name == models.BigqueryName {
				someBigQueryPlanId = service.Plans[0].ID
			}
			if service.Name == models.CloudsqlMySQLName {

				someCloudSQLPlanId = service.Plans[0].ID
			}
			if service.Name == models.StorageName {
				someStoragePlanId = service.Plans[0].ID
			}
		}

		for k, _ := range gcpBroker.ServiceBrokerMap {
			async := false
			if k == serviceNameToId[models.CloudsqlMySQLName] {
				async = true
			}
			gcpBroker.ServiceBrokerMap[k] = &modelsfakes.FakeServiceBrokerHelper{
				ProvisionsAsyncStub:   func() bool { return async },
				DeprovisionsAsyncStub: func() bool { return async },
				ProvisionStub: func(instanceId string, details models.ProvisionDetails, plan models.ServicePlan) (models.ServiceInstanceDetails, error) {
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
			ServiceID: serviceNameToId[models.CloudsqlMySQLName],
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
		It("should have 11 services in sevices map", func() {
			Expect(len(gcpBroker.ServiceBrokerMap)).To(Equal(11))
		})

		It("should have a default client", func() {
			Expect(brokerConfig.HttpConfig).NotTo(Equal(&jwt.Config{}))
		})

		It("should have loaded credentials correctly and have a project id", func() {
			Expect(brokerConfig.ProjectId).To(Equal("foo"))
		})
	})

	Describe("getting broker catalog", func() {
		It("should have 11 services available", func() {
			Expect(len(gcpBroker.Services())).To(Equal(11))
		})

		It("should record api fields in the catalog", func() {
			serviceList := gcpBroker.Services()

			for _, s := range serviceList {
				if s.ID == serviceNameToId[models.StorageName] || s.ID == serviceNameToId[models.CloudsqlMySQLName] {
					Expect(len(s.Plans[0].ServiceProperties)).ToNot(Equal(0))
				}
			}

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
				if s.ID == serviceNameToId[models.CloudsqlMySQLName] {
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
			os.Setenv("GOOGLE_STACKDRIVER_TRACE", fakes.PlanNoId)
			_, err := config.NewBrokerConfigFromEnv()
			Expect(err).To(HaveOccurred())
		})

		It("should have 1 datastore plan available", func() {
			serviceList := gcpBroker.Services()
			for _, s := range serviceList {
				if s.ID == serviceNameToId[models.DatastoreName] {
					Expect(len(s.Plans)).To(Equal(1))
				}
			}
		})
	})

	Describe("updating broker catalog", func() {

		It("should update plans on startup", func() {

			os.Setenv("GOOGLE_CLOUDSQL_MYSQL", fakes.CloudSqlNewPlan)

			newcfg, err := config.NewBrokerConfigFromEnv()
			Expect(err).ToNot(HaveOccurred())
			newBroker, err := brokers.New(newcfg, logger)
			Expect(err).ToNot(HaveOccurred())

			serviceList := newBroker.Services()
			for _, s := range serviceList {
				if s.ID == serviceNameToId[models.CloudsqlMySQLName] {
					Expect(s.Plans[0].Name).To(Equal("newPlan"))
					Expect(len(s.Plans)).To(Equal(1))

					Expect(err).ToNot(HaveOccurred())
					Expect(s.Plans[0].ServiceProperties["tier"]).To(Equal("D8"))
					Expect(s.Plans[0].ServiceProperties["max_disk_size"]).To(Equal("15"))

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
					ServiceID: serviceNameToId[models.CloudsqlMySQLName],
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
				Expect(gcpBroker.ServiceBrokerMap[serviceNameToId[models.CloudsqlMySQLName]].(*modelsfakes.FakeServiceBrokerHelper).PollInstanceCallCount()).To(Equal(1))
			})
		})

	})

	AfterEach(func() {
		os.Remove("test.db")
	})
})

var _ = Describe("AccountManagers", func() {

	var (
		logger            lager.Logger
		iamStyleBroker    models.ServiceBrokerHelper
		spannerBroker     models.ServiceBrokerHelper
		cloudsqlBroker    models.ServiceBrokerHelper
		accountManager    modelsfakes.FakeAccountManager
		sqlAccountManager modelsfakes.FakeAccountManager
		err               error
	)

	BeforeEach(func() {
		logger = lager.NewLogger("brokers_test")
		logger.RegisterSink(lager.NewWriterSink(GinkgoWriter, lager.DEBUG))

		testDb, _ := gorm.Open("sqlite3", "test.db")
		testDb.CreateTable(models.ServiceInstanceDetails{})
		testDb.CreateTable(models.ServiceBindingCredentials{})
		testDb.CreateTable(models.ProvisionRequestDetails{})

		db_service.DbConnection = testDb
		name_generator.New()

		accountManager = modelsfakes.FakeAccountManager{
			CreateAccountInGoogleStub: func(instanceID string, bindingID string, details models.BindDetails, instance models.ServiceInstanceDetails) (models.ServiceBindingCredentials, error) {
				return models.ServiceBindingCredentials{OtherDetails: "{}"}, nil
			},
		}
		sqlAccountManager = modelsfakes.FakeAccountManager{
			CreateAccountInGoogleStub: func(instanceID string, bindingID string, details models.BindDetails, instance models.ServiceInstanceDetails) (models.ServiceBindingCredentials, error) {
				return models.ServiceBindingCredentials{OtherDetails: "{}"}, nil
			},
		}

		iamStyleBroker = &pubsub.PubSubBroker{
			BrokerBase: broker_base.BrokerBase{
				AccountManager: &accountManager,
			},
		}

		cloudsqlBroker = &cloudsql.CloudSQLBroker{
			Logger:           logger,
			AccountManager:   &sqlAccountManager,
			SaAccountManager: &accountManager,
		}

		spannerBroker = &spanner.SpannerBroker{
			BrokerBase: broker_base.BrokerBase{
				AccountManager: &accountManager,
			},
		}
	})

	Describe("bind", func() {
		Context("when bind is called on an iam-style broker", func() {
			It("should call the account manager create account in google method", func() {
				_, err = iamStyleBroker.Bind("foo", "bar", models.BindDetails{})
				Expect(err).NotTo(HaveOccurred())
				Expect(accountManager.CreateCredentialsCallCount()).To(Equal(1))
			})
		})

		Context("when bind is called on a cloudsql broker after provision", func() {
			It("should call the account manager create account in google method", func() {
				db_service.DbConnection.Save(&models.ServiceInstanceDetails{ID: "foo"})
				_, err = iamStyleBroker.Bind("foo", "bar", models.BindDetails{})
				Expect(err).NotTo(HaveOccurred())
				Expect(accountManager.CreateCredentialsCallCount()).To(Equal(1))
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
				Expect(accountManager.CreateCredentialsCallCount()).To(Equal(1))
				Expect(sqlAccountManager.CreateCredentialsCallCount()).To(Equal(1))
				_, _, details, _ := accountManager.CreateCredentialsArgsForCall(0)
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
				_, err = iamStyleBroker.BuildInstanceCredentials(models.ServiceBindingCredentials{}, models.ServiceInstanceDetails{})
				Expect(err).ToNot(HaveOccurred())
				Expect(accountManager.BuildInstanceCredentialsCallCount()).To(Equal(1))
			})
		})
	})

	Describe("unbind", func() {
		Context("when unbind is called on the broker", func() {
			It("it should call the account manager delete account from google method", func() {
				err = iamStyleBroker.Unbind(models.ServiceBindingCredentials{})
				Expect(err).NotTo(HaveOccurred())
				Expect(accountManager.DeleteCredentialsCallCount()).To(Equal(1))
			})
		})

		Context("when unbind is called on a cloudsql broker", func() {
			It("it should call the account manager delete account from google method", func() {
				err = cloudsqlBroker.Unbind(models.ServiceBindingCredentials{})
				Expect(err).NotTo(HaveOccurred())
				Expect(accountManager.DeleteCredentialsCallCount()).To(Equal(1))
				Expect(sqlAccountManager.DeleteCredentialsCallCount()).To(Equal(1))
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
		os.Remove("test.db")
	})

})
