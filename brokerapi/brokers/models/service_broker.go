// Copyright the Service Broker Project Authors.
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
//
////////////////////////////////////////////////////////////////////////////////
//

package models

import (
	"encoding/json"
	"errors"
)

type ServiceBrokerHelper interface {
	Provision(instanceId string, details ProvisionDetails, plan PlanDetails) (ServiceInstanceDetails, error)
	Bind(instanceID, bindingID string, details BindDetails) (ServiceBindingCredentials, error)
	BuildInstanceCredentials(bindDetails map[string]string, instanceDetails map[string]string) map[string]string
	Unbind(details ServiceBindingCredentials) error
	Deprovision(instanceID string, details DeprovisionDetails) error
	PollInstance(instanceID string) (bool, error)
	LastOperationWasDelete(instanceID string) (bool, error)
	ProvisionsAsync() bool
	DeprovisionsAsync() bool
}

type ServiceBroker interface {
	Services() []Service

	Provision(instanceID string, details ProvisionDetails, asyncAllowed bool) (ProvisionedServiceSpec, error)
	Deprovision(instanceID string, details DeprovisionDetails, asyncAllowed bool) (IsAsync, error)

	Bind(instanceID, bindingID string, details BindDetails) (Binding, error)
	Unbind(instanceID, bindingID string, details UnbindDetails) error

	Update(instanceID string, details UpdateDetails, asyncAllowed bool) (IsAsync, error)

	LastOperation(instanceID string) (LastOperation, error)
}

type AccountManager interface {
	CreateAccountInGoogle(instanceID string, bindingID string, details BindDetails, instance ServiceInstanceDetails) (ServiceBindingCredentials, error)
	DeleteAccountFromGoogle(creds ServiceBindingCredentials) error
	BuildInstanceCredentials(bindDetails map[string]string, instanceDetails map[string]string) map[string]string
}

type GCPCredentials struct {
	Type                string `json:"type"`
	ProjectId           string `json:"project_id"`
	PrivateKeyId        string `json:"private_key_id"`
	PrivateKey          string `json:"private_key"`
	ClientEmail         string `json:"client_email"`
	ClientId            string `json:"client_id"`
	AuthUri             string `json:"auth_uri"`
	TokenUri            string `json:"token_uri"`
	AuthProviderCertUrl string `json:"auth_provider_x509_cert_url"`
	ClientCertUrl       string `json:"client_x509_cert_url"`
}

type IsAsync bool

type ProvisionDetails struct {
	ServiceID        string          `json:"service_id"`
	PlanID           string          `json:"plan_id"`
	OrganizationGUID string          `json:"organization_guid"`
	SpaceGUID        string          `json:"space_guid"`
	RawParameters    json.RawMessage `json:"parameters,omitempty"`
}

type ProvisionedServiceSpec struct {
	IsAsync      bool
	DashboardURL string
}

type BindDetails struct {
	AppGUID      string                 `json:"app_guid"`
	PlanID       string                 `json:"plan_id"`
	ServiceID    string                 `json:"service_id"`
	BindResource *BindResource          `json:"bind_resource,omitempty"`
	Parameters   map[string]interface{} `json:"parameters,omitempty"`
}

type BindResource struct {
	AppGuid string `json:"app_guid,omitempty"`
	Route   string `json:"route,omitempty"`
}

type UnbindDetails struct {
	PlanID    string `json:"plan_id"`
	ServiceID string `json:"service_id"`
}

type DeprovisionDetails struct {
	PlanID    string `json:"plan_id"`
	ServiceID string `json:"service_id"`
}

type UpdateDetails struct {
	ServiceID      string                 `json:"service_id"`
	PlanID         string                 `json:"plan_id"`
	Parameters     map[string]interface{} `json:"parameters"`
	PreviousValues PreviousValues         `json:"previous_values"`
}

type PreviousValues struct {
	PlanID    string `json:"plan_id"`
	ServiceID string `json:"service_id"`
	OrgID     string `json:"organization_id"`
	SpaceID   string `json:"space_id"`
}

type LastOperation struct {
	State       LastOperationState
	Description string
}

type LastOperationState string

const (
	InProgress LastOperationState = "in progress"
	Succeeded  LastOperationState = "succeeded"
	Failed     LastOperationState = "failed"
)

type Binding struct {
	Credentials     interface{} `json:"credentials"`
	SyslogDrainURL  string      `json:"syslog_drain_url,omitempty"`
	RouteServiceURL string      `json:"route_service_url,omitempty"`
}

var (
	ErrInstanceAlreadyExists  = errors.New("instance already exists")
	ErrInstanceDoesNotExist   = errors.New("instance does not exist")
	ErrInstanceLimitMet       = errors.New("instance limit for this service has been reached")
	ErrPlanQuotaExceeded      = errors.New("The quota for this service plan has been exceeded. Please contact your Operator for help.")
	ErrBindingAlreadyExists   = errors.New("binding already exists")
	ErrBindingDoesNotExist    = errors.New("binding does not exist")
	ErrAsyncRequired          = errors.New("This service plan requires client support for asynchronous service operations.")
	ErrServiceIsNotAsync      = errors.New("This service is not provisioned asynchronously; this operation does not apply.")
	ErrPlanChangeNotSupported = errors.New("The requested plan migration cannot be performed")
	ErrRawParamsInvalid       = errors.New("The format of the parameters is not valid JSON")
	ErrAppGuidNotProvided     = errors.New("app_guid is a required field but was not provided")
)

// This custom user agent string is added to provision calls so that Google can track the aggregated use of this tool
// We can better advocate for devoting resources to supporting cloud foundry and this service broker if we can show
// good usage statistics for it, so if you feel the need to fork this repo, please leave this string in place!
var CustomUserAgent = "cf-gcp-service-broker-test 3.3.0"

func ProductionizeUserAgent() {
	CustomUserAgent = "cf-gcp-service-broker 3.3.0"
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
const DatastoreName = "google-datastore"
const AppCredsEnvVar = "GOOGLE_APPLICATION_CREDENTIALS"
const AppCredsFileName = "application-default-credentials.json"
const RootSaEnvVar = "ROOT_SERVICE_ACCOUNT_JSON"
