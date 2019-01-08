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
	"errors"
	"os"
	"reflect"
	"testing"

	. "github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/api_service"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/bigquery"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/bigtable"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/broker_base"
	brokerbasefakes "github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/broker_base/broker_basefakes"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/cloudsql"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/dataflow"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/datastore"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/dialogflow"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/firestore"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/models"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/pubsub"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/spanner"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/stackdriver"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/storage"
	"github.com/GoogleCloudPlatform/gcp-service-broker/db_service"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/broker"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/broker/brokerfakes"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/providers/builtin"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/varcontext"
	"github.com/GoogleCloudPlatform/gcp-service-broker/utils"
	"github.com/pivotal-cf/brokerapi"

	"code.cloudfoundry.org/lager"

	"github.com/jinzhu/gorm"
	"golang.org/x/oauth2/jwt"
)

type lifecyclePoint int

const (
	None lifecyclePoint = iota
	Provisioned
	Bound
	Unbound
	Deprovisioned
)

const (
	fakeInstanceId = "newid"
	fakeBindingId  = "newbinding"
)

type serviceStub struct {
	ServiceId         string
	PlanId            string
	Provider          *brokerfakes.FakeServiceProvider
	ServiceDefinition *broker.ServiceDefinition

	realProvider broker.ServiceProvider
}

func (s *serviceStub) ProvisionDetails() brokerapi.ProvisionDetails {
	return brokerapi.ProvisionDetails{
		ServiceID: s.ServiceId,
		PlanID:    s.PlanId,
	}
}

func (s *serviceStub) DeprovisionDetails() brokerapi.DeprovisionDetails {
	return brokerapi.DeprovisionDetails{
		ServiceID: s.ServiceId,
		PlanID:    s.PlanId,
	}
}

func (s *serviceStub) BindDetails() brokerapi.BindDetails {
	return brokerapi.BindDetails{
		ServiceID: s.ServiceId,
		PlanID:    s.PlanId,
	}
}

func (s *serviceStub) UnbindDetails() brokerapi.UnbindDetails {
	return brokerapi.UnbindDetails{
		ServiceID: s.ServiceId,
		PlanID:    s.PlanId,
	}
}

func fakeService(isAsync bool) *serviceStub {
	defn := storage.ServiceDefinition()
	svc, err := defn.CatalogEntry()
	if err != nil {
		panic(err)
	}

	stub := serviceStub{
		ServiceId:         svc.ID,
		PlanId:            svc.Plans[0].ID,
		ServiceDefinition: defn,

		Provider: &brokerfakes.FakeServiceProvider{
			ProvisionsAsyncStub:   func() bool { return isAsync },
			DeprovisionsAsyncStub: func() bool { return isAsync },
			ProvisionStub: func(ctx context.Context, vc *varcontext.VarContext) (models.ServiceInstanceDetails, error) {
				return models.ServiceInstanceDetails{OtherDetails: "{\"mynameis\": \"instancename\"}"}, nil
			},
			BindStub: func(ctx context.Context, vc *varcontext.VarContext) (map[string]interface{}, error) {
				return map[string]interface{}{"foo": "bar"}, nil
			},
		},
	}

	stub.ServiceDefinition.ProviderBuilder = func(projectId string, auth *jwt.Config, logger lager.Logger) broker.ServiceProvider {
		return stub.Provider
	}

	return &stub
}

func newStubbedBroker(t *testing.T, registry broker.BrokerRegistry) (broker *GCPServiceBroker, closer func()) {
	// Set up database
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		t.Fatalf("couldn't create database: %v", err)
	}
	db_service.RunMigrations(db)
	db_service.DbConnection = db

	closer = func() {
		db.Close()
		os.Remove("test.db")
	}

	config := &BrokerConfig{
		ProjectId: "stub-project",
		Registry:  registry,
	}

	broker, err = New(config, utils.NewLogger("brokers-test"))
	if err != nil {
		t.Fatalf("couldn't create broker: %v", err)
	}

	return
}

func failIfErr(t *testing.T, action string, err error) {
	t.Helper()

	if err != nil {
		t.Fatalf("Expected no error while %s, got: %v", action, err)
	}
}

func assertEqual(t *testing.T, message string, expected, actual interface{}) {
	t.Helper()

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Error: %s Expected: %#v Actual: %#v", message, expected, actual)
	}
}

type BrokerEndpointTestCase struct {
	AsyncService bool
	ServiceState lifecyclePoint
	Check        func(t *testing.T, broker *GCPServiceBroker, stub *serviceStub)
}

type BrokerEndpointTestSuite map[string]BrokerEndpointTestCase

func (cases BrokerEndpointTestSuite) Run(t *testing.T) {
	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			stub := fakeService(tc.AsyncService)

			t.Log("Creating broker")
			registry := broker.BrokerRegistry{}
			registry.Register(stub.ServiceDefinition)

			broker, closer := newStubbedBroker(t, registry)
			defer closer()

			initService(t, tc.ServiceState, broker, stub)

			t.Log("Running check")
			tc.Check(t, broker, stub)
		})
	}
}

func initService(t *testing.T, point lifecyclePoint, broker *GCPServiceBroker, stub *serviceStub) {
	if point >= Provisioned {
		_, err := broker.Provision(context.Background(), fakeInstanceId, stub.ProvisionDetails(), true)
		failIfErr(t, "provisioning", err)
	}

	if point >= Bound {
		_, err := broker.Bind(context.Background(), fakeInstanceId, fakeBindingId, stub.BindDetails())
		failIfErr(t, "binding", err)
	}

	if point >= Unbound {
		err := broker.Unbind(context.Background(), fakeInstanceId, fakeBindingId, stub.UnbindDetails())
		failIfErr(t, "unbinding", err)
	}

	if point >= Deprovisioned {
		_, err := broker.Deprovision(context.Background(), fakeInstanceId, stub.DeprovisionDetails(), true)
		failIfErr(t, "deprovisioning", err)
	}
}

func TestGCPServiceBroker_Services(t *testing.T) {
	registry := builtin.BuiltinBrokerRegistry()
	broker, closer := newStubbedBroker(t, registry)
	defer closer()

	services, err := broker.Services(context.Background())
	failIfErr(t, "getting services", err)
	assertEqual(t, "service count should be the same", len(registry), len(services))
}

func TestGCPServiceBroker_Provision(t *testing.T) {
	cases := BrokerEndpointTestSuite{
		"good-request": {
			ServiceState: Provisioned,
			Check: func(t *testing.T, broker *GCPServiceBroker, stub *serviceStub) {

				assertEqual(t, "provision calls should match", 1, stub.Provider.ProvisionCallCount())
			},
		},
		"duplicate-request": {
			ServiceState: Provisioned,
			Check: func(t *testing.T, broker *GCPServiceBroker, stub *serviceStub) {
				_, err := broker.Provision(context.Background(), fakeInstanceId, stub.ProvisionDetails(), true)
				assertEqual(t, "errors should match", brokerapi.ErrInstanceAlreadyExists, err)
			},
		},
		"requires-async": {
			AsyncService: true,
			ServiceState: None,
			Check: func(t *testing.T, broker *GCPServiceBroker, stub *serviceStub) {
				// false for async support
				_, err := broker.Provision(context.Background(), fakeInstanceId, stub.ProvisionDetails(), false)
				assertEqual(t, "errors should match", brokerapi.ErrAsyncRequired, err)
			},
		},
		"unknown-service-id": {
			ServiceState: None,
			Check: func(t *testing.T, broker *GCPServiceBroker, stub *serviceStub) {
				req := stub.ProvisionDetails()
				req.ServiceID = "bad-service-id"
				_, err := broker.Provision(context.Background(), fakeInstanceId, req, true)
				assertEqual(t, "errors should match", errors.New("Unknown service ID: \"bad-service-id\""), err)
			},
		},
		"unknown-plan-id": {
			ServiceState: None,
			Check: func(t *testing.T, broker *GCPServiceBroker, stub *serviceStub) {
				req := stub.ProvisionDetails()
				req.PlanID = "bad-plan-id"
				_, err := broker.Provision(context.Background(), fakeInstanceId, req, true)
				assertEqual(t, "errors should match", errors.New("Plan ID \"bad-plan-id\" could not be found"), err)
			},
		},
	}

	cases.Run(t)
}

func TestGCPServiceBroker_Deprovision(t *testing.T) {
	cases := BrokerEndpointTestSuite{
		"good-request": {
			ServiceState: Provisioned,
			Check: func(t *testing.T, broker *GCPServiceBroker, stub *serviceStub) {
				_, err := broker.Deprovision(context.Background(), fakeInstanceId, stub.DeprovisionDetails(), true)
				failIfErr(t, "deprovisioning", err)

				assertEqual(t, "deprovision calls should match", 1, stub.Provider.DeprovisionCallCount())
			},
		},
		"duplicate-deprovision": {
			ServiceState: Deprovisioned,
			Check: func(t *testing.T, broker *GCPServiceBroker, stub *serviceStub) {
				_, err := broker.Deprovision(context.Background(), fakeInstanceId, stub.DeprovisionDetails(), true)
				assertEqual(t, "duplicate deprovision should lead to DNE", brokerapi.ErrInstanceDoesNotExist, err)
			},
		},
		"instance-does-not-exist": {
			Check: func(t *testing.T, broker *GCPServiceBroker, stub *serviceStub) {
				_, err := broker.Deprovision(context.Background(), fakeInstanceId, stub.DeprovisionDetails(), true)
				assertEqual(t, "instance does not exist should be set", brokerapi.ErrInstanceDoesNotExist, err)
			},
		},
		"async-required": {
			AsyncService: true,
			ServiceState: Provisioned,
			Check: func(t *testing.T, broker *GCPServiceBroker, stub *serviceStub) {
				_, err := broker.Deprovision(context.Background(), fakeInstanceId, stub.DeprovisionDetails(), false)
				assertEqual(t, "async required should be returned if not supported", brokerapi.ErrAsyncRequired, err)
			},
		},
	}

	cases.Run(t)
}

func TestGCPServiceBroker_Bind(t *testing.T) {
	cases := BrokerEndpointTestSuite{
		"good-request": {
			ServiceState: Bound,
			Check: func(t *testing.T, broker *GCPServiceBroker, stub *serviceStub) {
				assertEqual(t, "BindCallCount should match", 1, stub.Provider.BindCallCount())
				assertEqual(t, "BuildInstanceCredentialsCallCount should match", 1, stub.Provider.BuildInstanceCredentialsCallCount())
			},
		},
		"duplicate-request": {
			ServiceState: Bound,
			Check: func(t *testing.T, broker *GCPServiceBroker, stub *serviceStub) {
				_, err := broker.Bind(context.Background(), fakeInstanceId, fakeBindingId, stub.BindDetails())
				assertEqual(t, "errors should match", brokerapi.ErrBindingAlreadyExists, err)
			},
		},
		"bad-bind-call": {
			ServiceState: Provisioned,
			Check: func(t *testing.T, broker *GCPServiceBroker, stub *serviceStub) {
				req := stub.BindDetails()
				req.RawParameters = json.RawMessage(`{"role":"project.admin"}`)

				expectedErr := "1 error(s) occurred: role: role must be one of the following: \"storage.objectAdmin\", \"storage.objectCreator\", \"storage.objectViewer\""
				_, err := broker.Bind(context.Background(), fakeInstanceId, "bad-bind-call", req)
				assertEqual(t, "errors should match", expectedErr, err.Error())
			},
		},
	}

	cases.Run(t)
}

func TestGCPServiceBroker_Unbind(t *testing.T) {
	cases := BrokerEndpointTestSuite{
		"good-request": {
			ServiceState: Bound,
			Check: func(t *testing.T, broker *GCPServiceBroker, stub *serviceStub) {

				err := broker.Unbind(context.Background(), fakeInstanceId, fakeBindingId, stub.UnbindDetails())
				failIfErr(t, "unbinding", err)
			},
		},
		"multiple-unbinds": {
			ServiceState: Unbound,
			Check: func(t *testing.T, broker *GCPServiceBroker, stub *serviceStub) {
				err := broker.Unbind(context.Background(), fakeInstanceId, fakeBindingId, stub.UnbindDetails())
				assertEqual(t, "errors should match", brokerapi.ErrBindingDoesNotExist, err)
			},
		},
	}

	cases.Run(t)
}

func TestGCPServiceBroker_LastOperation(t *testing.T) {
	cases := BrokerEndpointTestSuite{
		"missing-instance": {ServiceState: Provisioned,

			Check: func(t *testing.T, broker *GCPServiceBroker, stub *serviceStub) {
				_, err := broker.LastOperation(context.Background(), "invalid-instance-id", "operationtoken")
				assertEqual(t, "errors should match", brokerapi.ErrInstanceDoesNotExist, err)
			},
		},
		"called-on-synchronous-service": {
			ServiceState: Provisioned,
			Check: func(t *testing.T, broker *GCPServiceBroker, stub *serviceStub) {
				_, err := broker.LastOperation(context.Background(), fakeInstanceId, "operationtoken")
				assertEqual(t, "errors should match", brokerapi.ErrAsyncRequired, err)
			},
		},
		"called-on-async-service": {
			AsyncService: true,
			ServiceState: Provisioned,
			Check: func(t *testing.T, broker *GCPServiceBroker, stub *serviceStub) {
				_, err := broker.LastOperation(context.Background(), fakeInstanceId, "operationtoken")
				failIfErr(t, "shouldn't be called on async service", err)

				assertEqual(t, "PollInstanceCallCount should match", 1, stub.Provider.PollInstanceCallCount())
			},
		},
	}

	cases.Run(t)
}

func TestServiceProviderAsync(t *testing.T) {
	cases := map[string]struct {
		AsyncProvisionExpected   bool
		AsyncDeprovisionExpected bool
		Provider                 broker.ServiceProvider
	}{
		"ml": {
			Provider: &api_service.ApiServiceBroker{},
		},
		"bigquery": {
			Provider: &bigquery.BigQueryBroker{},
		},
		"bigtable": {
			Provider: &bigtable.BigTableBroker{},
		},
		"cloudsql": {
			AsyncProvisionExpected:   true,
			AsyncDeprovisionExpected: true,
			Provider:                 &cloudsql.CloudSQLBroker{},
		},
		"dataflow": {
			Provider: &dataflow.DataflowBroker{},
		},
		"datastore": {
			Provider: &datastore.DatastoreBroker{},
		},
		"dialogflow": {
			Provider: &dialogflow.DialogflowBroker{},
		},
		"firestore": {
			Provider: &firestore.FirestoreBroker{},
		},
		"pubsub": {
			Provider: &pubsub.PubSubBroker{},
		},
		"spanner": {
			AsyncProvisionExpected: true,
			Provider:               &spanner.SpannerBroker{},
		},
		"storage": {
			Provider: &storage.StorageBroker{},
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			actualProvisionAsync := tc.Provider.ProvisionsAsync()
			assertEqual(t, "async provision should match", tc.AsyncProvisionExpected, actualProvisionAsync)

			actualDeprovisionAsync := tc.Provider.DeprovisionsAsync()
			assertEqual(t, "async deprovision should match", tc.AsyncDeprovisionExpected, actualDeprovisionAsync)
		})
	}
}

func TestThinWrapperServiceProviders(t *testing.T) {
	cases := map[string]func(broker_base.BrokerBase) broker.ServiceProvider{
		"pubsub": func(brokerBase broker_base.BrokerBase) broker.ServiceProvider {
			return &pubsub.PubSubBroker{BrokerBase: brokerBase}
		},
		"stackdriver": func(brokerBase broker_base.BrokerBase) broker.ServiceProvider {
			return &stackdriver.StackdriverAccountProvider{BrokerBase: brokerBase}
		},
		"ml": func(brokerBase broker_base.BrokerBase) broker.ServiceProvider {
			return &api_service.ApiServiceBroker{BrokerBase: brokerBase}
		},
		"bigquery": func(brokerBase broker_base.BrokerBase) broker.ServiceProvider {
			return &bigquery.BigQueryBroker{BrokerBase: brokerBase}
		},
		"dataflow": func(brokerBase broker_base.BrokerBase) broker.ServiceProvider {
			return &dataflow.DataflowBroker{BrokerBase: brokerBase}
		},
		"datastore": func(brokerBase broker_base.BrokerBase) broker.ServiceProvider {
			return &datastore.DatastoreBroker{BrokerBase: brokerBase}
		},
		"dialogflow": func(brokerBase broker_base.BrokerBase) broker.ServiceProvider {
			return &dialogflow.DialogflowBroker{BrokerBase: brokerBase}
		},
		"firestore": func(brokerBase broker_base.BrokerBase) broker.ServiceProvider {
			return &firestore.FirestoreBroker{BrokerBase: brokerBase}
		},
		"spanner": func(brokerBase broker_base.BrokerBase) broker.ServiceProvider {
			return &spanner.SpannerBroker{BrokerBase: brokerBase}
		},
		"storage": func(brokerBase broker_base.BrokerBase) broker.ServiceProvider {
			return &storage.StorageBroker{BrokerBase: brokerBase}
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			accountManager := brokerbasefakes.FakeServiceAccountManager{}
			brokerBase := broker_base.BrokerBase{
				AccountManager: &accountManager,
			}
			serviceProvider := tc(brokerBase)

			_, err := serviceProvider.Bind(context.Background(), &varcontext.VarContext{})
			failIfErr(t, "binding", err)
			assertEqual(t, "create credentials count should match", 1, accountManager.CreateCredentialsCallCount())

			serviceProvider.Unbind(context.Background(), models.ServiceInstanceDetails{}, models.ServiceBindingCredentials{})
			failIfErr(t, "unbinding", err)
			assertEqual(t, "delete credentials count should match", 1, accountManager.DeleteCredentialsCallCount())
		})
	}
}
