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

package models

import (
	"context"

	"github.com/pivotal-cf/brokerapi"
)

//go:generate counterfeiter . ServiceBrokerHelper
//go:generate counterfeiter . ServiceAccountManager

type ServiceBrokerHelper interface {
	Provision(ctx context.Context, instanceId string, details brokerapi.ProvisionDetails, plan ServicePlan) (ServiceInstanceDetails, error)
	// Bind creates credentials for accessing the service and stores information necessary to
	// access the service _and_ delete the binding in the returned map.
	Bind(ctx context.Context, instanceID, bindingID string, details brokerapi.BindDetails) (map[string]interface{}, error)
	BuildInstanceCredentials(ctx context.Context, bindRecord ServiceBindingCredentials, instanceRecord ServiceInstanceDetails) (map[string]interface{}, error)
	Unbind(ctx context.Context, details ServiceBindingCredentials) error
	// Deprovision deprovisions the service.
	// If the deprovision is asynchronous (results in a long-running job), then operationId is returned.
	// If no error and no operationId are returned, then the deprovision is expected to have been completed successfully.
	Deprovision(ctx context.Context, instance ServiceInstanceDetails, details brokerapi.DeprovisionDetails) (operationId *string, err error)
	PollInstance(ctx context.Context, instance ServiceInstanceDetails) (bool, error)
	ProvisionsAsync() bool
	DeprovisionsAsync() bool

	// UpdateInstanceDetails updates the ServiceInstanceDetails with the most recent state from GCP.
	// This function is optional, but will be called after async provisions, updates, and possibly
	// on broker version changes.
	// Return a nil error if you choose not to implement this function.
	UpdateInstanceDetails(ctx context.Context, instance *ServiceInstanceDetails) error
}

type ServiceAccountManager interface {
	CreateCredentials(ctx context.Context, instanceID string, bindingID string, details brokerapi.BindDetails, instance ServiceInstanceDetails) (map[string]interface{}, error)
	DeleteCredentials(ctx context.Context, creds ServiceBindingCredentials) error
	CreateAccountWithRoles(ctx context.Context, bindingID string, roles []string) (map[string]interface{}, error)
}

// This custom user agent string is added to provision calls so that Google can track the aggregated use of this tool
// We can better advocate for devoting resources to supporting cloud foundry and this service broker if we can show
// good usage statistics for it, so if you feel the need to fork this repo, please leave this string in place!
var CustomUserAgent = "cf-gcp-service-broker-test 4.0.0"

func ProductionizeUserAgent() {
	CustomUserAgent = "cf-gcp-service-broker 4.0.0"
}

const CloudPlatformScope = "https://www.googleapis.com/auth/cloud-platform"
const StorageName = "google-storage"
const BigqueryName = "google-bigquery"
const BigtableName = "google-bigtable"
const CloudsqlMySQLName = "google-cloudsql-mysql"
const CloudsqlPostgresName = "google-cloudsql-postgres"
const PubsubName = "google-pubsub"
const MlName = "google-ml-apis"
const SpannerName = "google-spanner"
const StackdriverTraceName = "google-stackdriver-trace"
const StackdriverDebuggerName = "google-stackdriver-debugger"
const StackdriverProfilerName = "google-stackdriver-profiler"
const DatastoreName = "google-datastore"
