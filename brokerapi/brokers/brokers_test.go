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
	"github.com/GoogleCloudPlatform/gcp-service-broker/db_service"
	"github.com/GoogleCloudPlatform/gcp-service-broker/db_service/models"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/broker"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/broker/brokerfakes"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/providers/builtin"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/providers/builtin/base"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/providers/builtin/storage"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/varcontext"
	"github.com/GoogleCloudPlatform/gcp-service-broker/utils"
	"github.com/pivotal-cf/brokerapi"
	"google.golang.org/api/googleapi"

	"code.cloudfoundry.org/lager"

	"github.com/jinzhu/gorm"
	"golang.org/x/oauth2/jwt"
)

// InstanceState holds the lifecycle state of a provisioned service instance.
// It goes None -> Provisioned -> Bound -> Unbound -> Deprovisioned
type InstanceState int

const (
	StateNone InstanceState = iota
	StateProvisioned
	StateBound
	StateUnbound
	StateDeprovisioned
)

const (
	fakeInstanceId = "newid"
	fakeBindingId  = "newbinding"
)

// serviceStub holds a stubbed out ServiceDefinition with easy access to
// its ID, a valid plan ID, and the mock provider.
type serviceStub struct {
	ServiceId         string
	PlanId            string
	Provider          *brokerfakes.FakeServiceProvider
	ServiceDefinition *broker.ServiceDefinition
}

// ProvisionDetails creates a brokerapi.ProvisionDetails object valid for
// the given service.
func (s *serviceStub) ProvisionDetails() brokerapi.ProvisionDetails {
	return brokerapi.ProvisionDetails{
		ServiceID: s.ServiceId,
		PlanID:    s.PlanId,
	}
}

// DeprovisionDetails creates a brokerapi.DeprovisionDetails object valid for
// the given service.
func (s *serviceStub) DeprovisionDetails() brokerapi.DeprovisionDetails {
	return brokerapi.DeprovisionDetails{
		ServiceID: s.ServiceId,
		PlanID:    s.PlanId,
	}
}

// BindDetails creates a brokerapi.BindDetails object valid for
// the given service.
func (s *serviceStub) BindDetails() brokerapi.BindDetails {
	return brokerapi.BindDetails{
		ServiceID: s.ServiceId,
		PlanID:    s.PlanId,
	}
}

// UnbindDetails creates a brokerapi.UnbindDetails object valid for
// the given service.
func (s *serviceStub) UnbindDetails() brokerapi.UnbindDetails {
	return brokerapi.UnbindDetails{
		ServiceID: s.ServiceId,
		PlanID:    s.PlanId,
	}
}

// fakeService creates a ServiceDefinition with a mock ServiceProvider and
// references to some important properties.
func fakeService(t *testing.T, isAsync bool) *serviceStub {
	defn := storage.ServiceDefinition()
	svc, err := defn.CatalogEntry()
	if err != nil {
		t.Fatal(err)
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
			BuildInstanceCredentialsStub: func(ctx context.Context, bc models.ServiceBindingCredentials, id models.ServiceInstanceDetails) (*brokerapi.Binding, error) {
				mixin := base.MergedInstanceCredsMixin{}
				return mixin.BuildInstanceCredentials(ctx, bc, id)
			},
		},
	}

	stub.ServiceDefinition.ProviderBuilder = func(projectId string, auth *jwt.Config, logger lager.Logger) broker.ServiceProvider {
		return stub.Provider
	}

	return &stub
}

// newStubbedBroker creates a new GCPServiceBroker with a dummy database for the given registry.
// It returns the broker and a callback used to clean up the database when done with it.
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

// failIfErr is a test helper function which stops the test immediately if the
// error is set.
func failIfErr(t *testing.T, action string, err error) {
	t.Helper()

	if err != nil {
		t.Fatalf("Expected no error while %s, got: %v", action, err)
	}
}

// assertEqual does a reflect.DeepEqual on the values and if they're different
// reports the message and the values.
func assertEqual(t *testing.T, message string, expected, actual interface{}) {
	t.Helper()

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Error: %s Expected: %#v Actual: %#v", message, expected, actual)
	}
}

// BrokerEndpointTestCase is the base test used for testing any
// brokerapi.ServiceBroker endpoint.
type BrokerEndpointTestCase struct {
	// The following properties are used to set up the environment for your test
	// to run in.
	AsyncService bool
	ServiceState InstanceState

	// Check is used to validate the state of the world and is where you should
	// put your test cases.
	Check func(t *testing.T, broker *GCPServiceBroker, stub *serviceStub)
}

// BrokerEndpointTestSuite holds a set of tests for a single endpoint.
type BrokerEndpointTestSuite map[string]BrokerEndpointTestCase

// Run executes every test case, setting up a new environment for each and
// tearing it down afterward.
func (cases BrokerEndpointTestSuite) Run(t *testing.T) {
	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			stub := fakeService(t, tc.AsyncService)

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

// initService creates a new service and brings it up to the lifecycle state given
// by state.
func initService(t *testing.T, state InstanceState, broker *GCPServiceBroker, stub *serviceStub) {
	if state >= StateProvisioned {
		_, err := broker.Provision(context.Background(), fakeInstanceId, stub.ProvisionDetails(), true)
		failIfErr(t, "provisioning", err)
	}

	if state >= StateBound {
		_, err := broker.Bind(context.Background(), fakeInstanceId, fakeBindingId, stub.BindDetails(), true)
		failIfErr(t, "binding", err)
	}

	if state >= StateUnbound {
		_, err := broker.Unbind(context.Background(), fakeInstanceId, fakeBindingId, stub.UnbindDetails(), true)
		failIfErr(t, "unbinding", err)
	}

	if state >= StateDeprovisioned {
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
			ServiceState: StateProvisioned,
			Check: func(t *testing.T, broker *GCPServiceBroker, stub *serviceStub) {

				assertEqual(t, "provision calls should match", 1, stub.Provider.ProvisionCallCount())
			},
		},
		"duplicate-request": {
			ServiceState: StateProvisioned,
			Check: func(t *testing.T, broker *GCPServiceBroker, stub *serviceStub) {
				_, err := broker.Provision(context.Background(), fakeInstanceId, stub.ProvisionDetails(), true)
				assertEqual(t, "errors should match", brokerapi.ErrInstanceAlreadyExists, err)
			},
		},
		"requires-async": {
			AsyncService: true,
			ServiceState: StateNone,
			Check: func(t *testing.T, broker *GCPServiceBroker, stub *serviceStub) {
				// false for async support
				_, err := broker.Provision(context.Background(), fakeInstanceId, stub.ProvisionDetails(), false)
				assertEqual(t, "errors should match", brokerapi.ErrAsyncRequired, err)
			},
		},
		"unknown-service-id": {
			ServiceState: StateNone,
			Check: func(t *testing.T, broker *GCPServiceBroker, stub *serviceStub) {
				req := stub.ProvisionDetails()
				req.ServiceID = "bad-service-id"
				_, err := broker.Provision(context.Background(), fakeInstanceId, req, true)
				assertEqual(t, "errors should match", errors.New("Unknown service ID: \"bad-service-id\""), err)
			},
		},
		"unknown-plan-id": {
			ServiceState: StateNone,
			Check: func(t *testing.T, broker *GCPServiceBroker, stub *serviceStub) {
				req := stub.ProvisionDetails()
				req.PlanID = "bad-plan-id"
				_, err := broker.Provision(context.Background(), fakeInstanceId, req, true)
				assertEqual(t, "errors should match", errors.New("Plan ID \"bad-plan-id\" could not be found"), err)
			},
		},
		"bad-request-json": {
			ServiceState: StateNone,
			Check: func(t *testing.T, broker *GCPServiceBroker, stub *serviceStub) {
				req := stub.ProvisionDetails()
				req.RawParameters = json.RawMessage("{invalid json")
				_, err := broker.Provision(context.Background(), fakeInstanceId, req, true)
				assertEqual(t, "errors should match", ErrInvalidUserInput, err)
			},
		},
	}

	cases.Run(t)
}

func TestGCPServiceBroker_Deprovision(t *testing.T) {
	cases := BrokerEndpointTestSuite{
		"good-request": {
			ServiceState: StateProvisioned,
			Check: func(t *testing.T, broker *GCPServiceBroker, stub *serviceStub) {
				_, err := broker.Deprovision(context.Background(), fakeInstanceId, stub.DeprovisionDetails(), true)
				failIfErr(t, "deprovisioning", err)

				assertEqual(t, "deprovision calls should match", 1, stub.Provider.DeprovisionCallCount())
			},
		},
		"duplicate-deprovision": {
			ServiceState: StateDeprovisioned,
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
			ServiceState: StateProvisioned,
			Check: func(t *testing.T, broker *GCPServiceBroker, stub *serviceStub) {
				_, err := broker.Deprovision(context.Background(), fakeInstanceId, stub.DeprovisionDetails(), false)
				assertEqual(t, "async required should be returned if not supported", brokerapi.ErrAsyncRequired, err)
			},
		},
		"async-deprovision-returns-operation": {
			AsyncService: true,
			ServiceState: StateProvisioned,
			Check: func(t *testing.T, broker *GCPServiceBroker, stub *serviceStub) {
				operationId := "my-operation-id"
				stub.Provider.DeprovisionReturns(&operationId, nil)
				resp, err := broker.Deprovision(context.Background(), fakeInstanceId, stub.DeprovisionDetails(), true)
				failIfErr(t, "deprovisioning", err)

				assertEqual(t, "operationid should be set as the data", operationId, resp.OperationData)
				assertEqual(t, "IsAsync should be set", true, resp.IsAsync)
			},
		},

		"async-deprovision-updates-db": {
			AsyncService: true,
			ServiceState: StateProvisioned,
			Check: func(t *testing.T, broker *GCPServiceBroker, stub *serviceStub) {
				operationId := "my-operation-id"
				stub.Provider.DeprovisionReturns(&operationId, nil)
				_, err := broker.Deprovision(context.Background(), fakeInstanceId, stub.DeprovisionDetails(), true)
				failIfErr(t, "deprovisioning", err)

				details, err := db_service.GetServiceInstanceDetailsById(context.Background(), fakeInstanceId)
				failIfErr(t, "looking up details", err)

				assertEqual(t, "OperationId should be set as the data", operationId, details.OperationId)
				assertEqual(t, "OperationType should be set as Deprovision", models.DeprovisionOperationType, details.OperationType)
			},
		},
	}

	cases.Run(t)
}

func TestGCPServiceBroker_Bind(t *testing.T) {
	cases := BrokerEndpointTestSuite{
		"good-request": {
			ServiceState: StateBound,
			Check: func(t *testing.T, broker *GCPServiceBroker, stub *serviceStub) {
				assertEqual(t, "BindCallCount should match", 1, stub.Provider.BindCallCount())
				assertEqual(t, "BuildInstanceCredentialsCallCount should match", 1, stub.Provider.BuildInstanceCredentialsCallCount())
			},
		},
		"duplicate-request": {
			ServiceState: StateBound,
			Check: func(t *testing.T, broker *GCPServiceBroker, stub *serviceStub) {
				_, err := broker.Bind(context.Background(), fakeInstanceId, fakeBindingId, stub.BindDetails(), true)
				assertEqual(t, "errors should match", brokerapi.ErrBindingAlreadyExists, err)
			},
		},
		"bad-bind-call": {
			ServiceState: StateProvisioned,
			Check: func(t *testing.T, broker *GCPServiceBroker, stub *serviceStub) {
				req := stub.BindDetails()
				req.RawParameters = json.RawMessage(`{"role":"project.admin"}`)

				expectedErr := "1 error(s) occurred: role: role must be one of the following: \"storage.objectAdmin\", \"storage.objectCreator\", \"storage.objectViewer\""
				_, err := broker.Bind(context.Background(), fakeInstanceId, "bad-bind-call", req, true)
				assertEqual(t, "errors should match", expectedErr, err.Error())
			},
		},
		"bad-request-json": {
			ServiceState: StateProvisioned,
			Check: func(t *testing.T, broker *GCPServiceBroker, stub *serviceStub) {
				req := stub.BindDetails()
				req.RawParameters = json.RawMessage("{invalid json")
				_, err := broker.Bind(context.Background(), fakeInstanceId, fakeBindingId, req, true)
				assertEqual(t, "errors should match", ErrInvalidUserInput, err)
			},
		},
	}

	cases.Run(t)
}

func TestGCPServiceBroker_Unbind(t *testing.T) {
	cases := BrokerEndpointTestSuite{
		"good-request": {
			ServiceState: StateBound,
			Check: func(t *testing.T, broker *GCPServiceBroker, stub *serviceStub) {

				_, err := broker.Unbind(context.Background(), fakeInstanceId, fakeBindingId, stub.UnbindDetails(), true)
				failIfErr(t, "unbinding", err)
			},
		},
		"multiple-unbinds": {
			ServiceState: StateUnbound,
			Check: func(t *testing.T, broker *GCPServiceBroker, stub *serviceStub) {
				_, err := broker.Unbind(context.Background(), fakeInstanceId, fakeBindingId, stub.UnbindDetails(), true)
				assertEqual(t, "errors should match", brokerapi.ErrBindingDoesNotExist, err)
			},
		},
	}

	cases.Run(t)
}

func TestGCPServiceBroker_LastOperation(t *testing.T) {
	cases := BrokerEndpointTestSuite{
		"missing-instance": {
			ServiceState: StateProvisioned,
			Check: func(t *testing.T, broker *GCPServiceBroker, stub *serviceStub) {
				_, err := broker.LastOperation(context.Background(), "invalid-instance-id", brokerapi.PollDetails{OperationData: "operationtoken"})
				assertEqual(t, "errors should match", brokerapi.ErrInstanceDoesNotExist, err)
			},
		},
		"called-on-synchronous-service": {
			ServiceState: StateProvisioned,
			Check: func(t *testing.T, broker *GCPServiceBroker, stub *serviceStub) {
				_, err := broker.LastOperation(context.Background(), fakeInstanceId, brokerapi.PollDetails{OperationData: "operationtoken"})
				assertEqual(t, "errors should match", brokerapi.ErrAsyncRequired, err)
			},
		},
		"called-on-async-service": {
			AsyncService: true,
			ServiceState: StateProvisioned,
			Check: func(t *testing.T, broker *GCPServiceBroker, stub *serviceStub) {
				_, err := broker.LastOperation(context.Background(), fakeInstanceId, brokerapi.PollDetails{OperationData: "operationtoken"})
				failIfErr(t, "shouldn't be called on async service", err)

				assertEqual(t, "PollInstanceCallCount should match", 1, stub.Provider.PollInstanceCallCount())
			},
		},
		"poll-returns-retryable-error": {
			AsyncService: true,
			ServiceState: StateProvisioned,
			Check: func(t *testing.T, broker *GCPServiceBroker, stub *serviceStub) {
				stub.Provider.PollInstanceReturns(false, &googleapi.Error{Code: 503})
				status, _ := broker.LastOperation(context.Background(), fakeInstanceId, brokerapi.PollDetails{OperationData: "operationtoken"})
				assertEqual(t, "retryable errors should result in in-progress state", brokerapi.InProgress, status.State)
			},
		},
		"poll-returns-failure": {
			AsyncService: true,
			ServiceState: StateProvisioned,
			Check: func(t *testing.T, broker *GCPServiceBroker, stub *serviceStub) {
				stub.Provider.PollInstanceReturns(false, errors.New("not-retryable"))
				status, _ := broker.LastOperation(context.Background(), fakeInstanceId, brokerapi.PollDetails{OperationData: "operationtoken"})
				assertEqual(t, "non-retryable errors should result in a failure state", brokerapi.Failed, status.State)
			},
		},
		"poll-returns-not-done": {
			AsyncService: true,
			ServiceState: StateProvisioned,
			Check: func(t *testing.T, broker *GCPServiceBroker, stub *serviceStub) {
				stub.Provider.PollInstanceReturns(false, nil)
				status, err := broker.LastOperation(context.Background(), fakeInstanceId, brokerapi.PollDetails{OperationData: "operationtoken"})
				failIfErr(t, "checking last operation", err)
				assertEqual(t, "polls that return no error should result in an in-progress state", brokerapi.InProgress, status.State)
			},
		},
		"poll-returns-success": {
			AsyncService: true,
			ServiceState: StateProvisioned,
			Check: func(t *testing.T, broker *GCPServiceBroker, stub *serviceStub) {
				stub.Provider.PollInstanceReturns(true, nil)
				status, err := broker.LastOperation(context.Background(), fakeInstanceId, brokerapi.PollDetails{OperationData: "operationtoken"})
				failIfErr(t, "checking last operation", err)
				assertEqual(t, "polls that return finished should result in a succeeded state", brokerapi.Succeeded, status.State)
			},
		},
	}

	cases.Run(t)
}

func TestGCPServiceBroker_GetBinding(t *testing.T) {
	cases := BrokerEndpointTestSuite{
		"called-on-bound": {
			ServiceState: StateBound,
			Check: func(t *testing.T, broker *GCPServiceBroker, stub *serviceStub) {
				_, err := broker.GetBinding(context.Background(), fakeInstanceId, fakeBindingId)

				assertEqual(t, "expect get binding not supported err", ErrGetBindingsUnsupported, err)
			},
		},
	}

	cases.Run(t)
}

func TestGCPServiceBroker_GetInstance(t *testing.T) {
	cases := BrokerEndpointTestSuite{
		"called-while-provisioned": {
			ServiceState: StateProvisioned,
			Check: func(t *testing.T, broker *GCPServiceBroker, stub *serviceStub) {
				_, err := broker.GetInstance(context.Background(), fakeInstanceId)

				assertEqual(t, "expect get instances not supported err", ErrGetInstancesUnsupported, err)
			},
		},
	}

	cases.Run(t)
}

func TestGCPServiceBroker_LastBindingOperation(t *testing.T) {
	cases := BrokerEndpointTestSuite{
		"called-while-bound": {
			ServiceState: StateProvisioned,
			AsyncService: true,
			Check: func(t *testing.T, broker *GCPServiceBroker, stub *serviceStub) {
				_, err := broker.LastBindingOperation(context.Background(), fakeInstanceId, fakeBindingId, brokerapi.PollDetails{})

				assertEqual(t, "expect last binding to return async required", brokerapi.ErrAsyncRequired, err)
			},
		},
	}

	cases.Run(t)
}
