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

package brokers_test

import (
	"context"
	"encoding/json"
	"os"

	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers"
	. "github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/broker_base"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/cloudsql"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/models"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/models/modelsfakes"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/name_generator"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/pubsub"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/spanner"
	"github.com/GoogleCloudPlatform/gcp-service-broker/db_service"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/broker"
	"github.com/pivotal-cf/brokerapi"

	"code.cloudfoundry.org/lager"

	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/config"
	"github.com/jinzhu/gorm"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"golang.org/x/oauth2/jwt"
)

var _ = Describe("Brokers", func() {
	var (
		gcpBroker                *GCPServiceBroker
		brokerConfig             *config.BrokerConfig
		err                      error
		logger                   lager.Logger
		serviceNameToId          map[string]string = make(map[string]string)
		bqProvisionDetails       brokerapi.ProvisionDetails
		cloudSqlProvisionDetails brokerapi.ProvisionDetails
		storageBindDetails       brokerapi.BindDetails
		storageUnbindDetails     brokerapi.UnbindDetails
		instanceId               string
		bindingId                string
	)

	BeforeEach(func() {
		logger = lager.NewLogger("brokers_test")
		logger.RegisterSink(lager.NewWriterSink(GinkgoWriter, lager.DEBUG))

		testDb, err := gorm.Open("sqlite3", "test.db")
		Expect(err).NotTo(HaveOccurred())
		db_service.RunMigrations(testDb)
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

		for k := range gcpBroker.ServiceBrokerMap {
			async := false
			if k == serviceNameToId[models.CloudsqlMySQLName] {
				async = true
			}
			gcpBroker.ServiceBrokerMap[k] = &modelsfakes.FakeServiceBrokerHelper{
				ProvisionsAsyncStub:   func() bool { return async },
				DeprovisionsAsyncStub: func() bool { return async },
				ProvisionStub: func(ctx context.Context, instanceId string, details brokerapi.ProvisionDetails, plan models.ServicePlan) (models.ServiceInstanceDetails, error) {
					return models.ServiceInstanceDetails{ID: instanceId, OtherDetails: "{\"mynameis\": \"instancename\"}"}, nil
				},
				BindStub: func(ctx context.Context, instanceID, bindingID string, details brokerapi.BindDetails) (models.ServiceBindingCredentials, error) {
					return models.ServiceBindingCredentials{OtherDetails: "{\"foo\": \"bar\"}"}, nil
				},
			}
		}

		bqProvisionDetails = brokerapi.ProvisionDetails{
			ServiceID: serviceNameToId[models.BigqueryName],
			PlanID:    someBigQueryPlanId,
		}

		cloudSqlProvisionDetails = brokerapi.ProvisionDetails{
			ServiceID: serviceNameToId[models.CloudsqlMySQLName],
			PlanID:    someCloudSQLPlanId,
		}

		storageBindDetails = brokerapi.BindDetails{
			ServiceID: serviceNameToId[models.StorageName],
			PlanID:    someStoragePlanId,
		}

		storageUnbindDetails = brokerapi.UnbindDetails{
			ServiceID: serviceNameToId[models.StorageName],
			PlanID:    someStoragePlanId,
		}

	})

	Describe("Broker init", func() {
		It("should have enabled services in sevices map", func() {
			Expect(len(gcpBroker.ServiceBrokerMap)).To(Equal(len(broker.GetEnabledServices())))
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
			serviceList, err := gcpBroker.Services(context.Background())
			Expect(err).ToNot(HaveOccurred())

			Expect(len(serviceList)).To(Equal(len(broker.GetEnabledServices())))
		})

		It("should have 4 storage plans available", func() {
			serviceList, err := gcpBroker.Services(context.Background())
			Expect(err).ToNot(HaveOccurred())
			for _, s := range serviceList {
				if s.ID == serviceNameToId[models.StorageName] {
					Expect(len(s.Plans)).To(Equal(4))
				}
			}

		})

		It("should have 15 cloudsql plans available", func() {
			serviceList, err := gcpBroker.Services(context.Background())
			Expect(err).ToNot(HaveOccurred())
			for _, s := range serviceList {
				if s.ID == serviceNameToId[models.CloudsqlMySQLName] {
					Expect(len(s.Plans)).To(Equal(15))
				}
			}

		})

		It("should have 2 bigtable plans available", func() {
			serviceList, err := gcpBroker.Services(context.Background())
			Expect(err).ToNot(HaveOccurred())
			for _, s := range serviceList {
				if s.ID == serviceNameToId[models.BigtableName] {
					Expect(len(s.Plans)).To(Equal(2))
				}
			}

		})

		It("should have 1 debugger plan available", func() {
			serviceList, err := gcpBroker.Services(context.Background())
			Expect(err).ToNot(HaveOccurred())
			for _, s := range serviceList {
				if s.ID == serviceNameToId[models.StackdriverDebuggerName] {
					Expect(len(s.Plans)).To(Equal(1))
				}
			}
		})

		It("should have 1 profiler plan available", func() {
			serviceList, err := gcpBroker.Services(context.Background())
			Expect(err).ToNot(HaveOccurred())
			for _, s := range serviceList {
				if s.ID == serviceNameToId[models.StackdriverProfilerName] {
					Expect(len(s.Plans)).To(Equal(1))
				}
			}
		})

		It("should have 1 datastore plan available", func() {
			serviceList, err := gcpBroker.Services(context.Background())
			Expect(err).ToNot(HaveOccurred())
			for _, s := range serviceList {
				if s.ID == serviceNameToId[models.DatastoreName] {
					Expect(len(s.Plans)).To(Equal(1))
				}
			}
		})
	})

	Describe("provision", func() {
		Context("when the bigquery service id is provided", func() {
			It("should call bigquery provisioning", func() {
				bqId := serviceNameToId[models.BigqueryName]
				_, err := gcpBroker.Provision(context.Background(), instanceId, bqProvisionDetails, true)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(gcpBroker.ServiceBrokerMap[bqId].(*modelsfakes.FakeServiceBrokerHelper).ProvisionCallCount()).To(Equal(1))
			})

		})

		Context("when an unrecognized service is provisioned", func() {
			It("should return an error", func() {
				_, err = gcpBroker.Provision(context.Background(), instanceId, brokerapi.ProvisionDetails{
					ServiceID: "nope",
					PlanID:    "nope",
				}, true)
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when an unrecognized plan is provisioned", func() {
			It("should return an error", func() {
				_, err = gcpBroker.Provision(context.Background(), instanceId, brokerapi.ProvisionDetails{
					ServiceID: serviceNameToId[models.BigqueryName],
					PlanID:    "nope",
				}, true)
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when duplicate services are provisioned", func() {
			It("should return an error", func() {
				_, err = gcpBroker.Provision(context.Background(), instanceId, bqProvisionDetails, true)
				Expect(err).NotTo(HaveOccurred())
				_, err := gcpBroker.Provision(context.Background(), instanceId, bqProvisionDetails, true)
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when async provisioning isn't allowed but the service requested requires it", func() {
			It("should return an error", func() {
				_, err := gcpBroker.Provision(context.Background(), instanceId, cloudSqlProvisionDetails, false)
				Expect(err).To(HaveOccurred())
			})
		})

	})

	Describe("deprovision", func() {
		Context("when the bigquery service id is provided", func() {
			It("should call bigquery deprovisioning", func() {
				bqId := serviceNameToId[models.BigqueryName]
				_, err := gcpBroker.Provision(context.Background(), instanceId, bqProvisionDetails, true)
				Expect(err).NotTo(HaveOccurred())
				_, err = gcpBroker.Deprovision(context.Background(), instanceId, brokerapi.DeprovisionDetails{
					ServiceID: bqId,
				}, true)
				Expect(err).NotTo(HaveOccurred())
				Expect(gcpBroker.ServiceBrokerMap[bqId].(*modelsfakes.FakeServiceBrokerHelper).DeprovisionCallCount()).To(Equal(1))
			})
		})

		Context("when the service doesn't exist", func() {
			It("should return an error", func() {
				_, err := gcpBroker.Deprovision(context.Background(), instanceId, brokerapi.DeprovisionDetails{
					ServiceID: serviceNameToId[models.BigqueryName],
				}, true)
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when async provisioning isn't allowed but the service requested requires it", func() {
			It("should return an error", func() {
				_, err := gcpBroker.Deprovision(context.Background(), instanceId, brokerapi.DeprovisionDetails{
					ServiceID: serviceNameToId[models.CloudsqlMySQLName],
				}, false)
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("bind", func() {
		Context("when bind is called on storage", func() {
			It("it should call storage bind", func() {
				_, err = gcpBroker.Provision(context.Background(), instanceId, bqProvisionDetails, true)
				Expect(err).NotTo(HaveOccurred())
				_, err = gcpBroker.Bind(context.Background(), instanceId, bindingId, storageBindDetails)
				Expect(err).NotTo(HaveOccurred())
				Expect(gcpBroker.ServiceBrokerMap[serviceNameToId[models.StorageName]].(*modelsfakes.FakeServiceBrokerHelper).BindCallCount()).To(Equal(1))
			})
		})

		Context("when bind is called more than once on the same id", func() {
			It("it should throw an error", func() {
				_, err = gcpBroker.Provision(context.Background(), instanceId, bqProvisionDetails, true)
				Expect(err).NotTo(HaveOccurred())
				_, err = gcpBroker.Bind(context.Background(), instanceId, bindingId, storageBindDetails)
				Expect(err).NotTo(HaveOccurred())
				_, err = gcpBroker.Bind(context.Background(), instanceId, bindingId, storageBindDetails)
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when bind is called", func() {
			It("it should update credentials with instance information", func() {
				_, err = gcpBroker.Provision(context.Background(), instanceId, bqProvisionDetails, true)
				Expect(err).NotTo(HaveOccurred())
				_, err := gcpBroker.Bind(context.Background(), instanceId, bindingId, storageBindDetails)
				Expect(err).NotTo(HaveOccurred())
				Expect(gcpBroker.ServiceBrokerMap[serviceNameToId[models.StorageName]].(*modelsfakes.FakeServiceBrokerHelper).BuildInstanceCredentialsCallCount()).To(Equal(1))
			})
		})

	})

	Describe("unbind", func() {
		Context("when unbind is called on storage", func() {
			It("it should call storage unbind", func() {
				_, err = gcpBroker.Provision(context.Background(), instanceId, bqProvisionDetails, true)
				Expect(err).NotTo(HaveOccurred())
				_, err = gcpBroker.Bind(context.Background(), instanceId, bindingId, storageBindDetails)
				Expect(err).NotTo(HaveOccurred())
				err = gcpBroker.Unbind(context.Background(), instanceId, bindingId, storageUnbindDetails)
				Expect(err).NotTo(HaveOccurred())
				Expect(gcpBroker.ServiceBrokerMap[serviceNameToId[models.StorageName]].(*modelsfakes.FakeServiceBrokerHelper).UnbindCallCount()).To(Equal(1))
			})
		})

		Context("when unbind is called more than once on the same id", func() {
			It("it should throw an error", func() {
				_, err = gcpBroker.Provision(context.Background(), instanceId, bqProvisionDetails, true)
				Expect(err).NotTo(HaveOccurred())
				_, err = gcpBroker.Bind(context.Background(), instanceId, bindingId, storageBindDetails)
				Expect(err).NotTo(HaveOccurred())
				err = gcpBroker.Unbind(context.Background(), instanceId, bindingId, storageUnbindDetails)
				Expect(err).NotTo(HaveOccurred())
				err = gcpBroker.Unbind(context.Background(), instanceId, bindingId, storageUnbindDetails)
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("lastOperation", func() {
		Context("when last operation is called on a service that doesn't exist", func() {
			It("should throw an error", func() {
				_, err = gcpBroker.LastOperation(context.Background(), "somethingnonexistant", "operationtoken")
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when last operation is called on a service that is provisioned synchronously", func() {
			It("should throw an error", func() {
				_, err = gcpBroker.Provision(context.Background(), instanceId, bqProvisionDetails, true)
				Expect(err).NotTo(HaveOccurred())
				_, err = gcpBroker.LastOperation(context.Background(), instanceId, "operationtoken")
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when last operation is called on an asynchronous service", func() {
			It("should call PollInstance", func() {
				_, err = gcpBroker.Provision(context.Background(), instanceId, cloudSqlProvisionDetails, true)
				Expect(err).NotTo(HaveOccurred())
				_, err = gcpBroker.LastOperation(context.Background(), instanceId, "operationtoken")
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
		accountManager    modelsfakes.FakeServiceAccountManager
		sqlAccountManager modelsfakes.FakeAccountManager
		err               error
		testCtx           context.Context
	)

	BeforeEach(func() {
		testCtx = context.Background()
		logger = lager.NewLogger("brokers_test")
		logger.RegisterSink(lager.NewWriterSink(GinkgoWriter, lager.DEBUG))

		testDb, err := gorm.Open("sqlite3", "test.db")
		Expect(err).NotTo(HaveOccurred())
		db_service.RunMigrations(testDb)
		db_service.DbConnection = testDb
		name_generator.New()

		accountManager = modelsfakes.FakeServiceAccountManager{
			CreateCredentialsStub: func(ctx context.Context, instanceID string, bindingID string, details brokerapi.BindDetails, instance models.ServiceInstanceDetails) (models.ServiceBindingCredentials, error) {
				return models.ServiceBindingCredentials{OtherDetails: "{}"}, nil
			},
		}
		sqlAccountManager = modelsfakes.FakeAccountManager{
			CreateCredentialsStub: func(ctx context.Context, instanceID string, bindingID string, details brokerapi.BindDetails, instance models.ServiceInstanceDetails) (models.ServiceBindingCredentials, error) {
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
				_, err = iamStyleBroker.Bind(context.Background(), "foo", "bar", brokerapi.BindDetails{})
				Expect(err).NotTo(HaveOccurred())
				Expect(accountManager.CreateCredentialsCallCount()).To(Equal(1))
			})
		})

		Context("when bind is called on a cloudsql broker after provision", func() {
			It("should call the account manager create account in google method", func() {
				db_service.SaveServiceInstanceDetails(testCtx, &models.ServiceInstanceDetails{ID: "foo"})
				_, err = iamStyleBroker.Bind(context.Background(), "foo", "bar", brokerapi.BindDetails{})
				Expect(err).NotTo(HaveOccurred())
				Expect(accountManager.CreateCredentialsCallCount()).To(Equal(1))
			})
		})

		Context("when bind is called on a cloudsql broker on a missing service instance", func() {
			It("should throw an error", func() {
				_, err = cloudsqlBroker.Bind(context.Background(), "foo", "bar", brokerapi.BindDetails{})
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when bind is called on a cloudsql broker with no username/password after provision", func() {
			It("should return a generated username and password", func() {
				db_service.CreateServiceInstanceDetails(testCtx, &models.ServiceInstanceDetails{ID: "foo"})

				_, err := cloudsqlBroker.Bind(context.Background(), "foo", "bar", brokerapi.BindDetails{})

				Expect(err).NotTo(HaveOccurred())
				Expect(accountManager.CreateCredentialsCallCount()).To(Equal(1))
				Expect(sqlAccountManager.CreateCredentialsCallCount()).To(Equal(1))
				_, _, _, details, _ := accountManager.CreateCredentialsArgsForCall(0)

				rawparams := details.GetRawParameters()
				params := make(map[string]interface{})
				err = json.Unmarshal(rawparams, &params)
				Expect(err).NotTo(HaveOccurred())

				Expect(params).NotTo(BeEmpty())

				username, usernameOk := params["username"].(string)
				password, passwordOk := params["password"].(string)

				Expect(usernameOk).To(BeTrue())
				Expect(passwordOk).To(BeTrue())
				Expect(username).NotTo(BeEmpty())
				Expect(password).NotTo(BeEmpty())
			})

		})

		Context("when MergeCredentialsAndInstanceInfo is called on a broker", func() {
			It("should call MergeCredentialsAndInstanceInfo on the account manager", func() {
				_, err = iamStyleBroker.BuildInstanceCredentials(context.Background(), models.ServiceBindingCredentials{}, models.ServiceInstanceDetails{})
				Expect(err).ToNot(HaveOccurred())
				Expect(accountManager.BuildInstanceCredentialsCallCount()).To(Equal(1))
			})
		})
	})

	Describe("unbind", func() {
		Context("when unbind is called on the broker", func() {
			It("it should call the account manager delete account from google method", func() {
				err = iamStyleBroker.Unbind(context.Background(), models.ServiceBindingCredentials{})
				Expect(err).NotTo(HaveOccurred())
				Expect(accountManager.DeleteCredentialsCallCount()).To(Equal(1))
			})
		})

		Context("when unbind is called on a cloudsql broker", func() {
			It("it should call the account manager delete account from google method", func() {
				err = cloudsqlBroker.Unbind(context.Background(), models.ServiceBindingCredentials{})
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
