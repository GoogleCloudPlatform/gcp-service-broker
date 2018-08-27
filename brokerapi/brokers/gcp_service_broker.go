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
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/account_managers"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/api_service"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/bigquery"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/bigtable"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/broker_base"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/cloudsql"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/stackdriver_profiler"
	"github.com/pivotal-cf/brokerapi"
	"google.golang.org/api/googleapi"

	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/config"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/datastore"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/models"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/pubsub"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/spanner"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/stackdriver_debugger"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/stackdriver_trace"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/storage"
	"github.com/GoogleCloudPlatform/gcp-service-broker/db_service"
)

type GCPServiceBroker struct {
	Catalog          map[string]models.Service
	ServiceBrokerMap map[string]models.ServiceBrokerHelper

	Logger lager.Logger
}

type GCPAsyncServiceBroker struct {
	GCPServiceBroker
}

// returns a new service broker and nil if no errors occur else nil and the error
func New(cfg *config.BrokerConfig, Logger lager.Logger) (*GCPAsyncServiceBroker, error) {

	self := GCPAsyncServiceBroker{}
	self.Logger = Logger
	self.Catalog = cfg.Catalog

	saManager := &account_managers.ServiceAccountManager{
		HttpConfig: cfg.HttpConfig,
		ProjectId:  cfg.ProjectId,
	}

	sqlManager := &account_managers.SqlAccountManager{
		HttpConfig: cfg.HttpConfig,
		ProjectId:  cfg.ProjectId,
	}

	bb := broker_base.BrokerBase{
		AccountManager: saManager,
		HttpConfig:     cfg.HttpConfig,
		ProjectId:      cfg.ProjectId,
		Logger:         self.Logger,
	}

	// map service specific brokers to general broker
	self.ServiceBrokerMap = map[string]models.ServiceBrokerHelper{
		models.StorageName: &storage.StorageBroker{
			BrokerBase: bb,
		},
		models.PubsubName: &pubsub.PubSubBroker{
			BrokerBase: bb,
		},
		models.StackdriverDebuggerName: &stackdriver_debugger.StackdriverDebuggerBroker{
			BrokerBase: bb,
		},
		models.StackdriverProfilerName: &stackdriver_profiler.StackdriverProfilerBroker{
			BrokerBase: bb,
		},
		models.StackdriverTraceName: &stackdriver_trace.StackdriverTraceBroker{
			BrokerBase: bb,
		},
		models.BigqueryName: &bigquery.BigQueryBroker{
			BrokerBase: bb,
		},
		models.MlName: &api_service.ApiServiceBroker{
			BrokerBase: bb,
		},
		models.CloudsqlMySQLName: &cloudsql.CloudSQLBroker{
			HttpConfig:       cfg.HttpConfig,
			ProjectId:        cfg.ProjectId,
			Logger:           self.Logger,
			AccountManager:   sqlManager,
			SaAccountManager: saManager,
		},
		models.CloudsqlPostgresName: &cloudsql.CloudSQLBroker{
			HttpConfig:       cfg.HttpConfig,
			ProjectId:        cfg.ProjectId,
			Logger:           self.Logger,
			AccountManager:   sqlManager,
			SaAccountManager: saManager,
		},
		models.BigtableName: &bigtable.BigTableBroker{
			BrokerBase: bb,
		},
		models.SpannerName: &spanner.SpannerBroker{
			BrokerBase: bb,
		},
		models.DatastoreName: &datastore.DatastoreBroker{
			BrokerBase: bb,
		},
	}

	// replace the mapping from name to a mapping from id
	for _, service := range self.Catalog {
		self.ServiceBrokerMap[service.ID] = self.ServiceBrokerMap[service.Name]
		delete(self.ServiceBrokerMap, service.Name)
	}

	return &self, nil
}

// CORE SERVICE BROKER API METHODS

// cf marketplace
// lists services in the broker's catalog
func (gcpBroker *GCPServiceBroker) Services(ctx context.Context) ([]brokerapi.Service, error) {
	svcs := []brokerapi.Service{}

	for _, svc := range gcpBroker.Catalog {
		svcs = append(svcs, svc.ToPlain())
	}

	return svcs, nil
}

func (gcpBroker *GCPServiceBroker) GetPlanFromId(serviceId, planId string) (models.ServicePlan, error) {
	if _, sidOk := gcpBroker.Catalog[serviceId]; !sidOk {
		return models.ServicePlan{}, fmt.Errorf("serviceId %s not found", serviceId)
	}

	for _, plan := range gcpBroker.Catalog[serviceId].Plans {
		if plan.ID == planId {
			return plan, nil
		}
	}

	return models.ServicePlan{}, fmt.Errorf("planId %s not found", planId)
}

// cf create-service
// creates a new service instance. What a "new service instance" means varies based on the service type
// CloudSQL: a new database instance and database
// BigQuery: a new dataset
// Storage: a new bucket
// PubSub: a new topic
// Bigtable: a new instance
//
func (gcpBroker *GCPAsyncServiceBroker) Provision(ctx context.Context, instanceID string, details brokerapi.ProvisionDetails, asyncAllowed bool) (brokerapi.ProvisionedServiceSpec, error) {
	gcpBroker.Logger.Info("Provisioning", lager.Data{
		"instance_id":  instanceID,
		"asyncAllowed": asyncAllowed,
		"details":      details,
	})

	// make sure that instance hasn't already been provisioned
	count, err := db_service.GetServiceInstanceCount(instanceID)
	if err != nil {
		return brokerapi.ProvisionedServiceSpec{}, fmt.Errorf("Database error checking for existing instance: %s", err)
	}
	if count > 0 {
		return brokerapi.ProvisionedServiceSpec{}, brokerapi.ErrInstanceAlreadyExists
	}

	serviceId := details.ServiceID

	plan, err := gcpBroker.GetPlanFromId(serviceId, details.PlanID)
	if err != nil {
		return brokerapi.ProvisionedServiceSpec{}, err
	}

	// verify async provisioning is allowed if it is required
	shouldProvisionAsync := gcpBroker.ServiceBrokerMap[serviceId].ProvisionsAsync()
	if shouldProvisionAsync && !asyncAllowed {
		return brokerapi.ProvisionedServiceSpec{}, brokerapi.ErrAsyncRequired
	}

	// get instance details
	instanceDetails, err := gcpBroker.ServiceBrokerMap[serviceId].Provision(instanceID, details, plan)
	if err != nil {
		return brokerapi.ProvisionedServiceSpec{}, err
	}

	// save instance details
	instanceDetails.ServiceId = serviceId
	instanceDetails.ID = instanceID
	instanceDetails.PlanId = details.PlanID
	instanceDetails.SpaceGuid = details.SpaceGUID
	instanceDetails.OrganizationGuid = details.OrganizationGUID

	err = db_service.DbConnection.Create(&instanceDetails).Error
	if err != nil {
		return brokerapi.ProvisionedServiceSpec{}, fmt.Errorf("Error saving instance details to database: %s. WARNING: this instance cannot be deprovisioned through cf. Contact your operator for cleanup", err)
	}

	// save provision request details
	pr := models.ProvisionRequestDetails{
		ServiceInstanceId: instanceID,
		RequestDetails:    string(details.RawParameters),
	}
	if err = db_service.DbConnection.Create(&pr).Error; err != nil {
		return brokerapi.ProvisionedServiceSpec{}, fmt.Errorf("Error saving provision request details to database: %s. Services relying on async provisioning will not be able to complete provisioning", err)
	}

	return brokerapi.ProvisionedServiceSpec{IsAsync: shouldProvisionAsync, DashboardURL: ""}, nil
}

// cf delete-service
// Deletes the given instance
func (gcpBroker *GCPAsyncServiceBroker) Deprovision(ctx context.Context, instanceID string, details brokerapi.DeprovisionDetails, asyncAllowed bool) (brokerapi.DeprovisionServiceSpec, error) {
	gcpBroker.Logger.Info("Deprovisioning", lager.Data{
		"instance_id":  instanceID,
		"asyncAllowed": asyncAllowed,
		"details":      details,
	})

	service := gcpBroker.ServiceBrokerMap[details.ServiceID]
	shouldDeprovisionAsync := service.DeprovisionsAsync()
	response := brokerapi.DeprovisionServiceSpec{IsAsync: shouldDeprovisionAsync}

	// make sure that instance actually exists
	count, err := db_service.GetServiceInstanceCount(instanceID)
	if err != nil {
		return response, fmt.Errorf("Database error checking for existing instance: %s", err)
	}
	if count == 0 {
		return response, brokerapi.ErrInstanceDoesNotExist
	}

	// if async provisioning isn't allowed but this service needs it, throw an error
	if shouldDeprovisionAsync && !asyncAllowed {
		return brokerapi.DeprovisionServiceSpec{IsAsync: asyncAllowed}, brokerapi.ErrAsyncRequired
	}

	// deprovision
	instance, err := db_service.GetServiceInstanceDetailsById(instanceID)
	if err != nil {
		return response, brokerapi.NewFailureResponseBuilder(err, http.StatusInternalServerError, "fetching instance from database").Build()
	}

	if err := service.Deprovision(ctx, *instance, instanceID, details); err != nil {
		return response, err
	}

	// soft-delete instance details from the db if this is a synchronous operation
	// if it's an async operation we can't delete from the db until we're sure delete succeeded, so this is
	// handled internally to LastOperation
	if !shouldDeprovisionAsync {
		err = db_service.SoftDeleteInstanceDetails(instanceID)
		if err != nil {
			return response, fmt.Errorf("Error deleting instance details from database: %s. WARNING: this instance will remain visible in cf. Contact your operator for cleanup", err)
		}
	}

	return response, nil
}

// cf bind-service
// for cloudSql instances, Bind creates a new user and ssl cert
// for all other services, Bind creates a new service account with the IAM role listed in details.Parameters["permissions"]
// a complete list of IAM roles is available here: https://cloud.google.com/iam/docs/understanding-roles
func (gcpBroker *GCPServiceBroker) Bind(ctx context.Context, instanceID, bindingID string, details brokerapi.BindDetails) (brokerapi.Binding, error) {
	gcpBroker.Logger.Info("Binding", lager.Data{
		"instance_id": instanceID,
		"binding_id":  bindingID,
		"details":     details,
	})

	service := gcpBroker.ServiceBrokerMap[details.ServiceID]

	// check for existing binding

	var count int
	var err error

	if err = db_service.DbConnection.Model(&models.ServiceBindingCredentials{}).Where("service_instance_id = ? and binding_id = ?", instanceID, bindingID).Count(&count).Error; err != nil {
		return brokerapi.Binding{}, fmt.Errorf("Error checking for existing binding: %s", err)
	}
	if count > 0 {
		return brokerapi.Binding{}, brokerapi.ErrBindingAlreadyExists
	}

	// create binding
	newCreds, err := service.Bind(instanceID, bindingID, details)
	if err != nil {
		return brokerapi.Binding{}, err
	}

	// save binding to database
	newCreds.ServiceInstanceId = instanceID
	newCreds.BindingId = bindingID
	newCreds.ServiceId = details.ServiceID

	if err := db_service.DbConnection.Create(&newCreds).Error; err != nil {
		return brokerapi.Binding{}, fmt.Errorf("Error saving credentials to database: %s. WARNING: these credentials cannot be unbound through cf. Please contact your operator for cleanup",
			err)
	}

	// get existing service instance details
	var instanceRecord models.ServiceInstanceDetails
	if err = db_service.DbConnection.Where("id = ?", instanceID).First(&instanceRecord).Error; err != nil {
		return brokerapi.Binding{}, fmt.Errorf("Error retrieving service instance details: %s", err)
	}

	updatedCreds, err := service.BuildInstanceCredentials(newCreds, instanceRecord)
	if err != nil {
		return brokerapi.Binding{}, err
	}

	return brokerapi.Binding{
		Credentials:     updatedCreds,
		SyslogDrainURL:  "",
		RouteServiceURL: "",
	}, nil
}

// cf unbind-service
// for cloudSql instances, Unbind deletes the associated user and ssl certs
// for all other services, Unbind deletes the associated service account
func (gcpBroker *GCPServiceBroker) Unbind(ctx context.Context, instanceID, bindingID string, details brokerapi.UnbindDetails) error {
	gcpBroker.Logger.Info("Unbinding", lager.Data{
		"instance_id": instanceID,
		"binding_id":  bindingID,
		"details":     details,
	})

	service := gcpBroker.ServiceBrokerMap[details.ServiceID]

	// validate existence of binding
	existingBinding := models.ServiceBindingCredentials{}
	if err := db_service.DbConnection.Where("service_instance_id = ? and binding_id = ?", instanceID, bindingID).Find(&existingBinding).Error; err != nil {
		return brokerapi.ErrBindingDoesNotExist
	}

	// remove binding from Google
	if err := service.Unbind(existingBinding); err != nil {
		return err
	}

	// remove binding from database
	if err := db_service.DeleteServiceBindingCredentials(&existingBinding); err != nil {
		return fmt.Errorf("Error soft-deleting credentials from database: %s. WARNING: these credentials will remain visible in cf. Contact your operator for cleanup", err)
	}

	return nil
}

// if a service is provisioned asynchronously, LastOperation is called until the provisioning attempt times out
// or success or failure is returned
func (gcpBroker *GCPServiceBroker) LastOperation(ctx context.Context, instanceID, operationData string) (brokerapi.LastOperation, error) {

	instance := models.ServiceInstanceDetails{}
	if err := db_service.DbConnection.Where("id = ?", instanceID).First(&instance).Error; err != nil {
		return brokerapi.LastOperation{}, brokerapi.ErrInstanceDoesNotExist
	}

	if gcpBroker.ServiceBrokerMap[instance.ServiceId].ProvisionsAsync() || gcpBroker.ServiceBrokerMap[instance.ServiceId].DeprovisionsAsync() {
		return gcpBroker.lastOperationAsync(instanceID, instance.ServiceId)

	} else {
		return brokerapi.LastOperation{State: brokerapi.Succeeded, Description: ""}, errors.New("Can't call LastOperation on a synchronous service")
	}
}

func (gcpBroker *GCPServiceBroker) lastOperationAsync(instanceId, serviceId string) (brokerapi.LastOperation, error) {
	done, err := gcpBroker.ServiceBrokerMap[serviceId].PollInstance(instanceId)
	if err != nil {
		// this is a retryable error
		if gerr, ok := err.(*googleapi.Error); ok {
			if gerr.Code == 503 {
				return brokerapi.LastOperation{State: brokerapi.InProgress, Description: ""}, err
			}
		}
		// This is not a retryable error. Return fail
		return brokerapi.LastOperation{State: brokerapi.Failed, Description: ""}, err
	}

	if done {
		// no error and we're done! Delete from the SB database if this was a delete flow and return success
		deleteFlow, err := gcpBroker.ServiceBrokerMap[serviceId].LastOperationWasDelete(instanceId)
		if err != nil {
			return brokerapi.LastOperation{State: brokerapi.Succeeded, Description: ""}, fmt.Errorf("Couldn't determine if provision or deprovision flow, this may leave orphaned resources, contact your operator for cleanup")
		}
		if deleteFlow {
			err = db_service.SoftDeleteInstanceDetails(instanceId)
			if err != nil {
				return brokerapi.LastOperation{State: brokerapi.Succeeded, Description: ""}, fmt.Errorf("Error deleting instance details from database: %s. WARNING: this instance will remain visible in cf. Contact your operator for cleanup", err)
			}
		}
		return brokerapi.LastOperation{State: brokerapi.Succeeded, Description: ""}, nil
	} else {
		return brokerapi.LastOperation{State: brokerapi.InProgress, Description: ""}, nil
	}
}

// Update a service instance plan. This functionality is not implemented and will return an error indicating that plan changes are not supported.
func (gcpBroker *GCPServiceBroker) Update(ctx context.Context, instanceID string, details brokerapi.UpdateDetails, asyncAllowed bool) (brokerapi.UpdateServiceSpec, error) {
	return brokerapi.UpdateServiceSpec{}, brokerapi.ErrPlanChangeNotSupported
}
