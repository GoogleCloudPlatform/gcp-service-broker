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

package brokers

import (
	"context"
	"fmt"

	"code.cloudfoundry.org/lager"
	"github.com/pivotal-cf/brokerapi"
	"golang.org/x/oauth2/jwt"
	"google.golang.org/api/googleapi"

	"encoding/json"

	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/models"
	"github.com/GoogleCloudPlatform/gcp-service-broker/db_service"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/broker"

	// import the brokers to register them
	_ "github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/api_service"
	_ "github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/bigquery"
	_ "github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/bigtable"
	_ "github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/cloudsql"
	_ "github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/dataflow"
	_ "github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/datastore"
	_ "github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/dialogflow"
	_ "github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/firestore"
	_ "github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/pubsub"
	_ "github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/spanner"
	_ "github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/stackdriver"
	_ "github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/storage"
	_ "github.com/GoogleCloudPlatform/gcp-service-broker/pkg/providers/tf"
)

// GCPServiceBroker is a brokerapi.ServiceBroker that can be used to generate an OSB compatible service broker.
type GCPServiceBroker struct {
	enableInputValidation bool
	registry              broker.BrokerRegistry
	jwtConfig             *jwt.Config
	projectId             string

	Logger lager.Logger
}

// New creates a GCPServiceBroker.
// Exactly one of GCPServiceBroker or error will be nil when returned.
func New(cfg *BrokerConfig, logger lager.Logger) (*GCPServiceBroker, error) {
	return &GCPServiceBroker{
		enableInputValidation: cfg.EnableInputValidation,
		registry:              cfg.Registry,
		jwtConfig:             cfg.HttpConfig,
		projectId:             cfg.ProjectId,
		Logger:                logger,
	}, nil
}

// Services lists services in the broker's catalog.
// It is called through the `GET /v2/catalog` endpoint or the `cf marketplace` command.
func (gcpBroker *GCPServiceBroker) Services(ctx context.Context) ([]brokerapi.Service, error) {
	svcs := []brokerapi.Service{}

	enabledServices, err := broker.GetEnabledServices()
	if err != nil {
		return nil, err
	}

	for _, service := range enabledServices {
		entry, err := service.CatalogEntry()
		if err != nil {
			return svcs, err
		}
		svcs = append(svcs, entry.ToPlain())
	}

	return svcs, nil
}

func (gcpBroker *GCPServiceBroker) getDefinitionAndProvider(serviceId string) (*broker.ServiceDefinition, broker.ServiceProvider, error) {
	defn, err := gcpBroker.registry.GetServiceById(serviceId)
	if err != nil {
		return nil, nil, err
	}

	providerBuilder := defn.ProviderBuilder(gcpBroker.projectId, gcpBroker.jwtConfig, gcpBroker.Logger)
	return defn, providerBuilder, nil
}

// Provision creates a new instance of a service.
// It is bound to the `PUT /v2/service_instances/:instance_id` endpoint and can be called using the `cf create-service` command.
func (gcpBroker *GCPServiceBroker) Provision(ctx context.Context, instanceID string, details brokerapi.ProvisionDetails, clientSupportsAsync bool) (brokerapi.ProvisionedServiceSpec, error) {
	gcpBroker.Logger.Info("Provisioning", lager.Data{
		"instanceId":         instanceID,
		"accepts_incomplete": clientSupportsAsync,
		"details":            details,
	})

	// make sure that instance hasn't already been provisioned
	count, err := db_service.CountServiceInstanceDetailsById(ctx, instanceID)
	if err != nil {
		return brokerapi.ProvisionedServiceSpec{}, fmt.Errorf("Database error checking for existing instance: %s", err)
	}
	if count > 0 {
		return brokerapi.ProvisionedServiceSpec{}, brokerapi.ErrInstanceAlreadyExists
	}

	brokerService, serviceHelper, err := gcpBroker.getDefinitionAndProvider(details.ServiceID)
	if err != nil {
		return brokerapi.ProvisionedServiceSpec{}, err
	}

	// verify the service exists and the plan exists
	plan, err := brokerService.GetPlanById(details.PlanID)
	if err != nil {
		return brokerapi.ProvisionedServiceSpec{}, err
	}

	// verify async provisioning is allowed if it is required
	shouldProvisionAsync := serviceHelper.ProvisionsAsync()
	if shouldProvisionAsync && !clientSupportsAsync {
		return brokerapi.ProvisionedServiceSpec{}, brokerapi.ErrAsyncRequired
	}

	if gcpBroker.enableInputValidation {
		// validate parameters meet the service's schema
		if err := gcpBroker.validateProvisionVariables(details); err != nil {
			return brokerapi.ProvisionedServiceSpec{}, err
		}
	}

	vars, err := brokerService.ProvisionVariables(instanceID, details, *plan)
	if err != nil {
		return brokerapi.ProvisionedServiceSpec{}, err
	}

	// get instance details
	instanceDetails, err := serviceHelper.Provision(ctx, vars)
	if err != nil {
		return brokerapi.ProvisionedServiceSpec{}, err
	}

	// save instance details
	instanceDetails.ServiceId = details.ServiceID
	instanceDetails.ID = instanceID
	instanceDetails.PlanId = details.PlanID
	instanceDetails.SpaceGuid = details.SpaceGUID
	instanceDetails.OrganizationGuid = details.OrganizationGUID

	err = db_service.CreateServiceInstanceDetails(ctx, &instanceDetails)
	if err != nil {
		return brokerapi.ProvisionedServiceSpec{}, fmt.Errorf("Error saving instance details to database: %s. WARNING: this instance cannot be deprovisioned through cf. Contact your operator for cleanup", err)
	}

	// save provision request details
	pr := models.ProvisionRequestDetails{
		ServiceInstanceId: instanceID,
		RequestDetails:    string(details.RawParameters),
	}
	if err = db_service.CreateProvisionRequestDetails(ctx, &pr); err != nil {
		return brokerapi.ProvisionedServiceSpec{}, fmt.Errorf("Error saving provision request details to database: %s. Services relying on async provisioning will not be able to complete provisioning", err)
	}

	return brokerapi.ProvisionedServiceSpec{IsAsync: shouldProvisionAsync, DashboardURL: "", OperationData: instanceDetails.OperationId}, nil
}

func (gcpBroker *GCPServiceBroker) validateProvisionVariables(details brokerapi.ProvisionDetails) error {
	serviceDefinition, err := gcpBroker.registry.GetServiceById(details.ServiceID)
	if err != nil {
		return err
	}

	params := make(map[string]interface{})
	if len(details.RawParameters) > 0 {
		if err := json.Unmarshal([]byte(details.RawParameters), &params); err != nil {
			return err
		}
	}

	return broker.ValidateVariables(params, serviceDefinition.ProvisionInputVariables)
}

func (gcpBroker *GCPServiceBroker) validateBindVariables(details brokerapi.BindDetails) error {
	serviceDefinition, err := gcpBroker.registry.GetServiceById(details.ServiceID)
	if err != nil {
		return err
	}

	params := make(map[string]interface{})
	if len(details.RawParameters) > 0 {
		if err := json.Unmarshal([]byte(details.RawParameters), &params); err != nil {
			return err
		}
	}

	return broker.ValidateVariables(params, serviceDefinition.BindInputVariables)
}

// Deprovision destroys an existing instance of a service.
// It is bound to the `DELETE /v2/service_instances/:instance_id` endpoint and can be called using the `cf delete-service` command.
// If a deprovision is asynchronous, the returned DeprovisionServiceSpec will contain the operation ID for tracking its progress.
func (gcpBroker *GCPServiceBroker) Deprovision(ctx context.Context, instanceID string, details brokerapi.DeprovisionDetails, clientSupportsAsync bool) (response brokerapi.DeprovisionServiceSpec, err error) {
	gcpBroker.Logger.Info("Deprovisioning", lager.Data{
		"instance_id":        instanceID,
		"accepts_incomplete": clientSupportsAsync,
		"details":            details,
	})

	// make sure that instance actually exists
	instance, err := db_service.GetServiceInstanceDetailsById(ctx, instanceID)
	if err != nil {
		return response, brokerapi.ErrInstanceDoesNotExist
	}

	_, serviceProvider, err := gcpBroker.getDefinitionAndProvider(instance.ServiceId)
	if err != nil {
		return response, err
	}

	// if async deprovisioning isn't allowed but this service needs it, throw an error
	if serviceProvider.DeprovisionsAsync() && !clientSupportsAsync {
		return response, brokerapi.ErrAsyncRequired
	}

	operationId, err := serviceProvider.Deprovision(ctx, *instance, details)
	if err != nil {
		return response, err
	}

	if operationId == nil {
		// soft-delete instance details from the db if this is a synchronous operation
		// if it's an async operation we can't delete from the db until we're sure delete succeeded, so this is
		// handled internally to LastOperation
		if err := db_service.DeleteServiceInstanceDetailsById(ctx, instanceID); err != nil {
			return response, fmt.Errorf("Error deleting instance details from database: %s. WARNING: this instance will remain visible in cf. Contact your operator for cleanup", err)
		}
		return response, nil
	} else {
		response.IsAsync = true
		response.OperationData = *operationId

		instance.OperationType = models.DeprovisionOperationType
		instance.OperationId = *operationId
		if err := db_service.SaveServiceInstanceDetails(ctx, instance); err != nil {
			return response, fmt.Errorf("Error saving instance details to database: %s. WARNING: this instance will remain visible in cf. Contact your operator for cleanup.", err)
		}
		return response, nil
	}
}

// Bind creates an account with credentials to access an instance of a service.
// It is bound to the `PUT /v2/service_instances/:instance_id/service_bindings/:binding_id` endpoint and can be called using the `cf bind-service` command.
func (gcpBroker *GCPServiceBroker) Bind(ctx context.Context, instanceID, bindingID string, details brokerapi.BindDetails) (brokerapi.Binding, error) {
	gcpBroker.Logger.Info("Binding", lager.Data{
		"instance_id": instanceID,
		"binding_id":  bindingID,
		"details":     details,
	})

	// check for existing binding
	count, err := db_service.CountServiceBindingCredentialsByServiceInstanceIdAndBindingId(ctx, instanceID, bindingID)
	if err != nil {
		return brokerapi.Binding{}, fmt.Errorf("Error checking for existing binding: %s", err)
	}
	if count > 0 {
		return brokerapi.Binding{}, brokerapi.ErrBindingAlreadyExists
	}

	// get existing service instance details
	instanceRecord, err := db_service.GetServiceInstanceDetailsById(ctx, instanceID)
	if err != nil {
		return brokerapi.Binding{}, fmt.Errorf("Error retrieving service instance details: %s", err)
	}

	serviceDefinition, serviceProvider, err := gcpBroker.getDefinitionAndProvider(instanceRecord.ServiceId)
	if err != nil {
		return brokerapi.Binding{}, err
	}

	if gcpBroker.enableInputValidation {
		// validate parameters meet the service's schema
		if err := gcpBroker.validateBindVariables(details); err != nil {
			return brokerapi.Binding{}, err
		}
	}

	vars, err := serviceDefinition.BindVariables(*instanceRecord, bindingID, details)
	if err != nil {
		return brokerapi.Binding{}, err
	}

	// create binding
	credsDetails, err := serviceProvider.Bind(ctx, vars)
	if err != nil {
		return brokerapi.Binding{}, err
	}

	serializedCreds, err := json.Marshal(credsDetails)
	if err != nil {
		return brokerapi.Binding{}, fmt.Errorf("Error serializing credentials: %s. WARNING: these credentials cannot be unbound through cf. Please contact your operator for cleanup", err)
	}

	// save binding to database
	newCreds := models.ServiceBindingCredentials{
		ServiceInstanceId: instanceID,
		BindingId:         bindingID,
		ServiceId:         details.ServiceID,
		OtherDetails:      string(serializedCreds),
	}

	if err := db_service.CreateServiceBindingCredentials(ctx, &newCreds); err != nil {
		return brokerapi.Binding{}, fmt.Errorf("Error saving credentials to database: %s. WARNING: these credentials cannot be unbound through cf. Please contact your operator for cleanup",
			err)
	}

	updatedCreds, err := serviceProvider.BuildInstanceCredentials(ctx, newCreds, *instanceRecord)
	if err != nil {
		return brokerapi.Binding{}, err
	}

	return brokerapi.Binding{Credentials: updatedCreds}, nil
}

// Unbind destroys an account and credentials with access to an instance of a service.
// It is bound to the `DELETE /v2/service_instances/:instance_id/service_bindings/:binding_id` endpoint and can be called using the `cf unbind-service` command.
func (gcpBroker *GCPServiceBroker) Unbind(ctx context.Context, instanceID, bindingID string, details brokerapi.UnbindDetails) error {
	gcpBroker.Logger.Info("Unbinding", lager.Data{
		"instance_id": instanceID,
		"binding_id":  bindingID,
		"details":     details,
	})

	_, serviceProvider, err := gcpBroker.getDefinitionAndProvider(details.ServiceID)
	if err != nil {
		return err
	}

	// validate existence of binding
	existingBinding, err := db_service.GetServiceBindingCredentialsByServiceInstanceIdAndBindingId(ctx, instanceID, bindingID)
	if err != nil {
		return brokerapi.ErrBindingDoesNotExist
	}

	// get existing service instance details
	instance, err := db_service.GetServiceInstanceDetailsById(ctx, instanceID)
	if err != nil {
		return fmt.Errorf("Error retrieving service instance details: %s", err)
	}

	// remove binding from Google
	if err := serviceProvider.Unbind(ctx, *instance, *existingBinding); err != nil {
		return err
	}

	// remove binding from database
	if err := db_service.DeleteServiceBindingCredentials(ctx, existingBinding); err != nil {
		return fmt.Errorf("Error soft-deleting credentials from database: %s. WARNING: these credentials will remain visible in cf. Contact your operator for cleanup", err)
	}

	return nil
}

// Unbind destroys an account and credentials with access to an instance of a service.
// It is bound to the `GET /v2/service_instances/:instance_id/last_operation` endpoint.
// It is called by `cf create-service` or `cf delete-service` if the operation was asynchronous.
func (gcpBroker *GCPServiceBroker) LastOperation(ctx context.Context, instanceID, operationData string) (brokerapi.LastOperation, error) {
	gcpBroker.Logger.Info("Last Operation", lager.Data{
		"instance_id":    instanceID,
		"operation_data": operationData,
	})

	instance, err := db_service.GetServiceInstanceDetailsById(ctx, instanceID)
	if err != nil {
		return brokerapi.LastOperation{}, brokerapi.ErrInstanceDoesNotExist
	}

	_, serviceProvider, err := gcpBroker.getDefinitionAndProvider(instance.ServiceId)
	if err != nil {
		return brokerapi.LastOperation{}, err
	}

	isAsyncService := serviceProvider.ProvisionsAsync() || serviceProvider.DeprovisionsAsync()
	if !isAsyncService {
		return brokerapi.LastOperation{}, brokerapi.ErrAsyncRequired
	}

	lastOperationType := instance.OperationType

	done, err := serviceProvider.PollInstance(ctx, *instance)
	if err != nil {
		// this is a retryable error
		if gerr, ok := err.(*googleapi.Error); ok {
			if gerr.Code == 503 {
				return brokerapi.LastOperation{State: brokerapi.InProgress}, err
			}
		}
		// This is not a retryable error. Return fail
		return brokerapi.LastOperation{State: brokerapi.Failed}, err
	}

	if !done {
		return brokerapi.LastOperation{State: brokerapi.InProgress}, nil
	}

	// the instance may have been invalidated, so we pass its primary key rather than the
	// instance directly.
	updateErr := gcpBroker.updateStateOnOperationCompletion(ctx, serviceProvider, lastOperationType, instanceID)
	return brokerapi.LastOperation{State: brokerapi.Succeeded}, updateErr
}

// updateStateOnOperationCompletion handles updating/cleaning-up resources that need to be changed
// once lastOperation finishes successfully.
func (gcpBroker *GCPServiceBroker) updateStateOnOperationCompletion(ctx context.Context, service broker.ServiceProvider, lastOperationType, instanceID string) error {
	if lastOperationType == models.DeprovisionOperationType {
		if err := db_service.DeleteServiceInstanceDetailsById(ctx, instanceID); err != nil {
			return fmt.Errorf("Error deleting instance details from database: %s. WARNING: this instance will remain visible in cf. Contact your operator for cleanup", err)
		}

		return nil
	}

	// If the operation was not a delete, clear out the ID and type and update
	// any changed (or finalized) state like IP addresses, selflinks, etc.
	details, err := db_service.GetServiceInstanceDetailsById(ctx, instanceID)
	if err != nil {
		return fmt.Errorf("Error getting instance details from database %v", err)
	}

	if err := service.UpdateInstanceDetails(ctx, details); err != nil {
		return fmt.Errorf("Error getting new instance details from GCP: %v", err)
	}

	details.OperationId = ""
	details.OperationType = models.ClearOperationType
	if err := db_service.SaveServiceInstanceDetails(ctx, details); err != nil {
		return fmt.Errorf("Error saving instance details to database %v", err)
	}

	return nil
}

// Update a service instance plan.
// This functionality is not implemented and will return an error indicating that plan changes are not supported.
func (gcpBroker *GCPServiceBroker) Update(ctx context.Context, instanceID string, details brokerapi.UpdateDetails, asyncAllowed bool) (brokerapi.UpdateServiceSpec, error) {
	return brokerapi.UpdateServiceSpec{}, brokerapi.ErrPlanChangeNotSupported
}
