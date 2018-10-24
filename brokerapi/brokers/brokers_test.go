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
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/models"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/models/modelsfakes"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/name_generator"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/pubsub"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/spanner"
	"github.com/GoogleCloudPlatform/gcp-service-broker/db_service"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/broker"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/varcontext"
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
		storageBadBindDetails    brokerapi.BindDetails
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
				BindStub: func(ctx context.Context, vc *varcontext.VarContext) (map[string]interface{}, error) {
					return map[string]interface{}{"foo": "bar"}, nil
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
			ServiceID:     serviceNameToId[models.StorageName],
			PlanID:        someStoragePlanId,
			RawParameters: json.RawMessage(`{"role":"storage.objectAdmin"}`),
		}

		storageBadBindDetails = brokerapi.BindDetails{
			ServiceID:     serviceNameToId[models.StorageName],
			PlanID:        someStoragePlanId,
			RawParameters: json.RawMessage(`{"role":"storage.admin"}`),
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

			It("it should reject bad roles", func() {
				_, err = gcpBroker.Provision(context.Background(), instanceId, bqProvisionDetails, true)
				Expect(err).NotTo(HaveOccurred())
				_, err = gcpBroker.Bind(context.Background(), instanceId, bindingId, storageBadBindDetails)
				Expect(err).To(HaveOccurred())
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
		logger         lager.Logger
		iamStyleBroker models.ServiceBrokerHelper
		spannerBroker  models.ServiceBrokerHelper
		accountManager modelsfakes.FakeServiceAccountManager
		err            error
		testCtx        context.Context
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
			CreateCredentialsStub: func(ctx context.Context, vc *varcontext.VarContext) (map[string]interface{}, error) {
				return map[string]interface{}{}, nil
			},
		}

		iamStyleBroker = &pubsub.PubSubBroker{
			BrokerBase: broker_base.BrokerBase{
				AccountManager: &accountManager,
			},
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
				_, err = iamStyleBroker.Bind(context.Background(), &varcontext.VarContext{})
				Expect(err).NotTo(HaveOccurred())
				Expect(accountManager.CreateCredentialsCallCount()).To(Equal(1))
			})
		})

		Context("when bind is called on an iam-style broker after provision", func() {
			It("should call the account manager create account in google method", func() {
				instance := models.ServiceInstanceDetails{ID: "foo"}
				db_service.SaveServiceInstanceDetails(testCtx, &instance)
				_, err = iamStyleBroker.Bind(context.Background(), &varcontext.VarContext{})
				Expect(err).NotTo(HaveOccurred())
				Expect(accountManager.CreateCredentialsCallCount()).To(Equal(1))
			})
		})
	})

	Describe("unbind", func() {
		Context("when unbind is called on the broker", func() {
			It("it should call the account manager delete account from google method", func() {
				err = iamStyleBroker.Unbind(context.Background(), models.ServiceInstanceDetails{}, models.ServiceBindingCredentials{})
				Expect(err).NotTo(HaveOccurred())
				Expect(accountManager.DeleteCredentialsCallCount()).To(Equal(1))
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
