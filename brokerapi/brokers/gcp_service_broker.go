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
	"errors"
	"fmt"
	"net/http"

	"code.cloudfoundry.org/lager"
	"github.com/pivotal-cf/brokerapi"
	"golang.org/x/oauth2/jwt"
	"google.golang.org/api/googleapi"

	"encoding/json"

	"github.com/GoogleCloudPlatform/gcp-service-broker/db_service"
	"github.com/GoogleCloudPlatform/gcp-service-broker/db_service/models"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/broker"
)

var (
	invalidUserInputMsg = "User supplied paramaters must be in the form of a valid JSON map."
	ErrInvalidUserInput = brokerapi.NewFailureResponse(errors.New(invalidUserInputMsg), http.StatusBadRequest, "parsing-user-request")
)

// GCPServiceBroker is a brokerapi.ServiceBroker that can be used to generate an OSB compatible service broker.
type GCPServiceBroker struct {
	registry              broker.BrokerRegistry
	jwtConfig             *jwt.Config
	projectId             string
	requestContextService *RequestContextService

	Logger lager.Logger
}

// New creates a GCPServiceBroker.
// Exactly one of GCPServiceBroker or error will be nil when returned.
func New(cfg *BrokerConfig, logger lager.Logger) (*GCPServiceBroker, error) {
	return &GCPServiceBroker{
		registry:              cfg.Registry,
		jwtConfig:             cfg.HttpConfig,
		projectId:             cfg.ProjectId,
		requestContextService: NewRequestContextService(cfg.Registry),
		Logger:                logger,
	}, nil
}

// Services lists services in the broker's catalog.
// It is called through the `GET /v2/catalog` endpoint or the `cf marketplace` command.
func (gcpBroker *GCPServiceBroker) Services(ctx context.Context) ([]brokerapi.Service, error) {
	svcs := []brokerapi.Service{}

	enabledServices, err := gcpBroker.registry.GetEnabledServices()
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

func (gcpBroker *GCPServiceBroker) getProvider(definition *broker.ServiceDefinition) broker.ServiceProvider {
	return definition.ProviderBuilder(gcpBroker.projectId, gcpBroker.jwtConfig, gcpBroker.Logger)
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
	exists, err := db_service.ExistsServiceInstanceDetailsById(ctx, instanceID)
	if err != nil {
		return brokerapi.ProvisionedServiceSpec{}, fmt.Errorf("Database error checking for existing instance: %s", err)
	}
	if exists {
		return brokerapi.ProvisionedServiceSpec{}, brokerapi.ErrInstanceAlreadyExists
	}

	rc, err := gcpBroker.requestContextService.ServiceContext(ctx, details.ServiceID, details.PlanID)
	if err != nil {
		return brokerapi.ProvisionedServiceSpec{}, err
	}

	serviceProvider := gcpBroker.getProvider(rc.Definition)

	// verify async provisioning is allowed if it is required
	shouldProvisionAsync := serviceProvider.ProvisionsAsync()
	if shouldProvisionAsync && !clientSupportsAsync {
		return brokerapi.ProvisionedServiceSpec{}, brokerapi.ErrAsyncRequired
	}

	// Give the user a better error message if they give us a bad request
	if !isValidOrEmptyJSON(details.GetRawParameters()) {
		return brokerapi.ProvisionedServiceSpec{}, ErrInvalidUserInput
	}

	// validate parameters meet the service's schema and merge the user vars with
	// the plan's
	vars, err := rc.Definition.ProvisionVariables(instanceID, details, *rc.Plan)
	if err != nil {
		return brokerapi.ProvisionedServiceSpec{}, err
	}

	// get instance details
	instanceDetails, err := serviceProvider.Provision(ctx, vars)
	if err != nil {
		return brokerapi.ProvisionedServiceSpec{}, err
	}

	// save instance details
	instanceDetails.ServiceId = details.ServiceID
	instanceDetails.ID = instanceID
	instanceDetails.PlanId = details.PlanID
	instanceDetails.SpaceGuid = details.SpaceGUID
	instanceDetails.OrganizationGuid = details.OrganizationGUID

	if err := db_service.CreateServiceInstanceDetails(ctx, &instanceDetails); err != nil {
		return brokerapi.ProvisionedServiceSpec{}, fmt.Errorf("Error saving instance details to database: %s. WARNING: this instance cannot be deprovisioned through cf. Contact your operator for cleanup", err)
	}

	// save provision request details
	pr := models.ProvisionRequestDetails{
		ServiceInstanceId: instanceID,
		RequestDetails:    string(details.RawParameters),
	}
	if err := db_service.CreateProvisionRequestDetails(ctx, &pr); err != nil {
		return brokerapi.ProvisionedServiceSpec{}, fmt.Errorf("Error saving provision request details to database: %s. Services relying on async provisioning will not be able to complete provisioning", err)
	}

	return brokerapi.ProvisionedServiceSpec{IsAsync: shouldProvisionAsync, DashboardURL: "", OperationData: instanceDetails.OperationId}, nil
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

	rc, err := gcpBroker.requestContextService.InstanceContext(ctx, instanceID)
	if err != nil {
		return response, err
	}

	instance := rc.ServiceInstance
	serviceProvider := gcpBroker.getProvider(rc.Definition)

	// if async deprovisioning isn't allowed but this service needs it, throw an error
	if serviceProvider.DeprovisionsAsync() && !clientSupportsAsync {
		return response, brokerapi.ErrAsyncRequired
	}

	operationID, err := serviceProvider.Deprovision(ctx, *instance, details)
	if err != nil {
		return response, err
	}

	if operationID == nil {
		// soft-delete instance details from the db if this is a synchronous operation
		// if it's an async operation we can't delete from the db until we're sure delete succeeded, so this is
		// handled internally to LastOperation
		if err := db_service.DeleteServiceInstanceDetailsById(ctx, instanceID); err != nil {
			return response, fmt.Errorf("Error deleting instance details from database: %s. WARNING: this instance will remain visible in cf. Contact your operator for cleanup", err)
		}
		return response, nil
	} else {
		response.IsAsync = true
		response.OperationData = *operationID

		instance.OperationType = models.DeprovisionOperationType
		instance.OperationId = *operationID
		if err := db_service.SaveServiceInstanceDetails(ctx, instance); err != nil {
			return response, fmt.Errorf("Error saving instance details to database: %s. WARNING: this instance will remain visible in cf. Contact your operator for cleanup.", err)
		}
		return response, nil
	}
}

// Bind creates an account with credentials to access an instance of a service.
// It is bound to the `PUT /v2/service_instances/:instance_id/service_bindings/:binding_id` endpoint and can be called using the `cf bind-service` command.
func (gcpBroker *GCPServiceBroker) Bind(ctx context.Context, instanceID, bindingID string, details brokerapi.BindDetails, clientSupportsAsync bool) (brokerapi.Binding, error) {
	gcpBroker.Logger.Info("Binding", lager.Data{
		"instance_id": instanceID,
		"binding_id":  bindingID,
		"details":     details,
	})

	// check for existing binding
	exists, err := db_service.ExistsServiceBindingCredentialsByServiceInstanceIdAndBindingId(ctx, instanceID, bindingID)
	if err != nil {
		return brokerapi.Binding{}, fmt.Errorf("Error checking for existing binding: %s", err)
	}
	if exists {
		return brokerapi.Binding{}, brokerapi.ErrBindingAlreadyExists
	}

	rc, err := gcpBroker.requestContextService.InstanceContext(ctx, instanceID)
	if err != nil {
		return brokerapi.Binding{}, err
	}

	// Give the user a better error message if they give us a bad request
	if !isValidOrEmptyJSON(details.GetRawParameters()) {
		return brokerapi.Binding{}, ErrInvalidUserInput
	}

	// validate parameters meet the service's schema and merge the plan's vars with
	// the user's
	vars, err := rc.Definition.BindVariables(*rc.ServiceInstance, bindingID, details)
	if err != nil {
		return brokerapi.Binding{}, err
	}

	// create binding
	serviceProvider := gcpBroker.getProvider(rc.Definition)
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

	creds, err := gcpBroker.getCredentails(ctx, instanceID, bindingID)
	if err != nil {
		return brokerapi.Binding{}, err
	}

	return brokerapi.Binding{Credentials: creds}, nil
}

// GetBinding fetches an existing service binding.
// GET /v2/service_instances/{instance_id}/service_bindings/{binding_id}
func (broker *GCPServiceBroker) GetBinding(ctx context.Context, instanceID, bindingID string) (brokerapi.GetBindingSpec, error) {
	broker.Logger.Info("GetBinding", lager.Data{
		"instance_id": instanceID,
		"binding_id":  bindingID,
	})

	creds, err := broker.getCredentails(ctx, instanceID, bindingID)
	if err != nil {
		return brokerapi.GetBindingSpec{}, err
	}

	return brokerapi.GetBindingSpec{Credentials: creds}, nil
}

// getCredentials merges the credentials for the service and binding and returns them.
func (broker *GCPServiceBroker) getCredentails(ctx context.Context, instanceID, bindingID string) (map[string]interface{}, error) {
	rc, err := broker.requestContextService.BindingContext(ctx, instanceID, bindingID)
	if err != nil {
		return nil, err
	}

	serviceProvider := broker.getProvider(rc.Definition)
	return serviceProvider.BuildInstanceCredentials(ctx, *rc.BindingInstance, *rc.ServiceInstance)
}

// GetInstance fetches information about a service instance
// GET /v2/service_instances/{instance_id}
func (broker *GCPServiceBroker) GetInstance(ctx context.Context, instanceID string) (brokerapi.GetInstanceDetailsSpec, error) {
	broker.Logger.Info("GetInstance", lager.Data{
		"instance_id": instanceID,
	})

	instance, err := broker.requestContextService.LookupInstanceDetails(ctx, instanceID)
	if err != nil {
		return brokerapi.GetInstanceDetailsSpec{}, err
	}

	return brokerapi.GetInstanceDetailsSpec{
		PlanID:    instance.PlanId,
		ServiceID: instance.ServiceId,
	}, nil
}

// LastBindingOperation fetches last operation state for a service binding.
// GET /v2/service_instances/{instance_id}/service_bindings/{binding_id}/last_operation
//
// NOTE: This functionality is not implemented.
func (broker *GCPServiceBroker) LastBindingOperation(ctx context.Context, instanceID, bindingID string, details brokerapi.PollDetails) (brokerapi.LastOperation, error) {
	broker.Logger.Info("LastBindingOperation", lager.Data{
		"instance_id":    instanceID,
		"binding_id":     bindingID,
		"plan_id":        details.PlanID,
		"service_id":     details.ServiceID,
		"operation_data": details.OperationData,
	})

	return brokerapi.LastOperation{}, brokerapi.ErrAsyncRequired
}

// Unbind destroys an account and credentials with access to an instance of a service.
// It is bound to the `DELETE /v2/service_instances/:instance_id/service_bindings/:binding_id` endpoint and can be called using the `cf unbind-service` command.
func (gcpBroker *GCPServiceBroker) Unbind(ctx context.Context, instanceID, bindingID string, details brokerapi.UnbindDetails, asyncSupported bool) (brokerapi.UnbindSpec, error) {
	gcpBroker.Logger.Info("Unbinding", lager.Data{
		"instance_id": instanceID,
		"binding_id":  bindingID,
		"details":     details,
	})

	rc, err := gcpBroker.requestContextService.BindingContext(ctx, instanceID, bindingID)
	if err != nil {
		return brokerapi.UnbindSpec{}, err
	}
	serviceProvider := gcpBroker.getProvider(rc.Definition)

	// remove binding from Google
	if err := serviceProvider.Unbind(ctx, *rc.ServiceInstance, *rc.BindingInstance); err != nil {
		return brokerapi.UnbindSpec{}, err
	}

	// remove binding from database
	if err := db_service.DeleteServiceBindingCredentials(ctx, rc.BindingInstance); err != nil {
		return brokerapi.UnbindSpec{}, fmt.Errorf("Error soft-deleting credentials from database: %s. WARNING: these credentials will remain visible in cf. Contact your operator for cleanup", err)
	}

	return brokerapi.UnbindSpec{}, nil
}

// LastOperation fetches last operation state for a service instance.
// It is bound to the `GET /v2/service_instances/:instance_id/last_operation` endpoint.
// It is called by `cf create-service` or `cf delete-service` if the operation was asynchronous.
func (gcpBroker *GCPServiceBroker) LastOperation(ctx context.Context, instanceID string, details brokerapi.PollDetails) (brokerapi.LastOperation, error) {
	gcpBroker.Logger.Info("Last Operation", lager.Data{
		"instance_id":    instanceID,
		"plan_id":        details.PlanID,
		"service_id":     details.ServiceID,
		"operation_data": details.OperationData,
	})

	rc, err := gcpBroker.requestContextService.InstanceContext(ctx, instanceID)
	if err != nil {
		return brokerapi.LastOperation{}, err
	}
	serviceProvider := gcpBroker.getProvider(rc.Definition)

	isAsyncService := serviceProvider.ProvisionsAsync() || serviceProvider.DeprovisionsAsync()
	if !isAsyncService {
		return brokerapi.LastOperation{}, brokerapi.ErrAsyncRequired
	}

	done, err := serviceProvider.PollInstance(ctx, *rc.ServiceInstance)
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
	lastOperationType := rc.ServiceInstance.OperationType
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

func isValidOrEmptyJSON(msg json.RawMessage) bool {
	return msg == nil || len(msg) == 0 || json.Valid(msg)
}
