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

package brokers

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"

	"code.cloudfoundry.org/lager"
	"google.golang.org/api/googleapi"

	"gcp-service-broker/brokerapi/brokers/account_managers"
	"gcp-service-broker/brokerapi/brokers/api_service"
	"gcp-service-broker/brokerapi/brokers/bigquery"
	"gcp-service-broker/brokerapi/brokers/bigtable"
	"gcp-service-broker/brokerapi/brokers/broker_base"
	"gcp-service-broker/brokerapi/brokers/cloudsql"
	"gcp-service-broker/brokerapi/brokers/models"
	"gcp-service-broker/brokerapi/brokers/pubsub"
	"gcp-service-broker/brokerapi/brokers/spanner"
	"gcp-service-broker/brokerapi/brokers/stackdriver_debugger"
	"gcp-service-broker/brokerapi/brokers/stackdriver_trace"
	"gcp-service-broker/brokerapi/brokers/storage"
	"gcp-service-broker/db_service"
	"gcp-service-broker/utils"
	"golang.org/x/oauth2/jwt"
	"gopkg.in/validator.v2"
)

type GCPServiceBroker struct {
	RootGCPCredentials *models.GCPCredentials
	HttpConfig         *jwt.Config
	Catalog            *[]models.Service
	ServiceBrokerMap   map[string]models.ServiceBrokerHelper

	InstanceLimit int

	Logger lager.Logger
}

type GCPAsyncServiceBroker struct {
	GCPServiceBroker
	ShouldProvisionAsync bool
}

// returns a new service broker and nil if no errors occur else nil and the error
func New(Logger lager.Logger) (*GCPAsyncServiceBroker, error) {
	var err error

	self := GCPAsyncServiceBroker{}
	self.Logger = Logger
	self.ShouldProvisionAsync = false

	// hardcoding this for now so we don't have to delete the nice built-in quota code, but also don't have to
	// handle that as a config option.
	self.InstanceLimit = math.MaxInt32

	// save credentials to broker object
	rootCreds, err := GetCredentialsFromEnv()
	if err != nil {
		return nil, fmt.Errorf("Error initializing GCP credentials: %s", err)
	}
	self.RootGCPCredentials = &rootCreds

	// set up GCP client with root gcp credentials
	cfg, err := utils.GetAuthedConfig()
	if err != nil {
		return nil, fmt.Errorf("Error getting authorized http client: %s", err)
	}
	self.HttpConfig = cfg

	// save catalog to broker object

	cat, err := InitCatalogFromEnv()
	if err != nil {
		return nil, fmt.Errorf("Error initializing catalog: %s", err)
	}
	self.Catalog = &cat

	saManager := &account_managers.ServiceAccountManager{
		HttpConfig: self.HttpConfig,
		ProjectId:  self.RootGCPCredentials.ProjectId,
	}

	sqlManager := &account_managers.SqlAccountManager{
		HttpConfig: self.HttpConfig,
		ProjectId:  self.RootGCPCredentials.ProjectId,
	}

	// map service specific brokers to general broker
	self.ServiceBrokerMap = map[string]models.ServiceBrokerHelper{
		models.StorageName: &storage.StorageBroker{
			HttpConfig: self.HttpConfig,
			ProjectId:  self.RootGCPCredentials.ProjectId,
			Logger:     self.Logger,
			BrokerBase: broker_base.BrokerBase{
				AccountManager: saManager,
			},
		},
		models.PubsubName: &pubsub.PubSubBroker{
			HttpConfig: self.HttpConfig,
			ProjectId:  self.RootGCPCredentials.ProjectId,
			Logger:     self.Logger,
			BrokerBase: broker_base.BrokerBase{
				AccountManager: saManager,
			},
		},
		models.StackdriverDebuggerName: &stackdriver_debugger.StackdriverDebuggerBroker{
			HttpConfig:            self.HttpConfig,
			ProjectId:             self.RootGCPCredentials.ProjectId,
			Logger:                self.Logger,
			ServiceAccountManager: saManager,
			BrokerBase: broker_base.BrokerBase{
				AccountManager: saManager,
			},
		},
		models.StackdriverTraceName: &stackdriver_trace.StackdriverTraceBroker{
			HttpConfig:            self.HttpConfig,
			ProjectId:             self.RootGCPCredentials.ProjectId,
			Logger:                self.Logger,
			ServiceAccountManager: saManager,
			BrokerBase: broker_base.BrokerBase{
				AccountManager: saManager,
			},
		},
		models.BigqueryName: &bigquery.BigQueryBroker{
			HttpConfig: self.HttpConfig,
			ProjectId:  self.RootGCPCredentials.ProjectId,
			Logger:     self.Logger,
			BrokerBase: broker_base.BrokerBase{
				AccountManager: saManager,
			},
		},
		models.MlName: &api_service.ApiServiceBroker{
			HttpConfig: self.HttpConfig,
			ProjectId:  self.RootGCPCredentials.ProjectId,
			Logger:     self.Logger,
			BrokerBase: broker_base.BrokerBase{
				AccountManager: saManager,
			},
		},
		models.CloudsqlName: &cloudsql.CloudSQLBroker{
			HttpConfig:     self.HttpConfig,
			ProjectId:      self.RootGCPCredentials.ProjectId,
			Logger:         self.Logger,
			AccountManager: sqlManager,
		},
		models.BigtableName: &bigtable.BigTableBroker{
			HttpConfig: self.HttpConfig,
			ProjectId:  self.RootGCPCredentials.ProjectId,
			Logger:     self.Logger,
			BrokerBase: broker_base.BrokerBase{
				AccountManager: saManager,
			},
		},
		models.SpannerName: &spanner.SpannerBroker{
			HttpConfig: self.HttpConfig,
			ProjectId:  self.RootGCPCredentials.ProjectId,
			Logger:     self.Logger,
			BrokerBase: broker_base.BrokerBase{
				AccountManager: saManager,
			},
		},
	}
	// replace the mapping from name to a mapping from id
	for _, service := range *self.Catalog {
		self.ServiceBrokerMap[service.ID] = self.ServiceBrokerMap[service.Name]
		delete(self.ServiceBrokerMap, service.Name)
	}

	return &self, nil

}

// CORE SERVICE BROKER API METHODS

// cf marketplace
// lists services in the broker's catalog
func (gcpBroker *GCPServiceBroker) Services() []models.Service {

	return *gcpBroker.Catalog
}

func (gcpBroker *GCPServiceBroker) GetPlanFromId(serviceId, planId string) (models.ServicePlan, error) {
	for _, s := range *gcpBroker.Catalog {
		if s.ID == serviceId {
			for _, p := range s.Plans {
				if p.ID == planId {
					return p, nil
				}
			}
		}
	}
	return models.ServicePlan{}, fmt.Errorf("Plan with id %s and serviceId %s not found", planId, serviceId)
}

// cf create-service
// creates a new service instance. What a "new service instance" means varies based on the service type
// CloudSQL: a new database instance and database
// BigQuery: a new dataset
// Storage: a new bucket
// PubSub: a new topic
// Bigtable: a new instance
func (gcpBroker *GCPAsyncServiceBroker) Provision(instanceID string, details models.ProvisionDetails, asyncAllowed bool) (models.ProvisionedServiceSpec, error) {
	var err error

	// first make sure we're not over quota
	provisionedInstancesCount, err := db_service.GetServiceInstanceTotal()
	if err != nil {
		return models.ProvisionedServiceSpec{}, fmt.Errorf("Database error checking for instance count: %s", err)
	} else {
		if provisionedInstancesCount >= gcpBroker.InstanceLimit {
			return models.ProvisionedServiceSpec{}, models.ErrInstanceLimitMet
		}
	}

	// make sure that instance hasn't already been provisioned
	count, err := db_service.GetServiceInstanceCount(instanceID)
	if err != nil {
		return models.ProvisionedServiceSpec{}, fmt.Errorf("Database error checking for existing instance: %s", err)
	}
	if count > 0 {
		return models.ProvisionedServiceSpec{}, models.ErrInstanceAlreadyExists
	}

	serviceId := details.ServiceID

	plan, err := gcpBroker.GetPlanFromId(serviceId, details.PlanID)
	if err != nil {
		return models.ProvisionedServiceSpec{}, err
	}

	// verify async provisioning is allowed if it is required
	gcpBroker.ShouldProvisionAsync = gcpBroker.ServiceBrokerMap[serviceId].ProvisionsAsync()
	if gcpBroker.ShouldProvisionAsync && !asyncAllowed {
		return models.ProvisionedServiceSpec{}, models.ErrAsyncRequired
	}

	// get instance details
	instanceDetails, err := gcpBroker.ServiceBrokerMap[serviceId].Provision(instanceID, details, plan)
	if err != nil {
		return models.ProvisionedServiceSpec{}, err
	}

	// save instance details
	instanceDetails.ServiceId = serviceId
	instanceDetails.ID = instanceID
	instanceDetails.PlanId = details.PlanID
	instanceDetails.SpaceGuid = details.SpaceGUID
	instanceDetails.OrganizationGuid = details.OrganizationGUID

	err = db_service.DbConnection.Create(&instanceDetails).Error
	if err != nil {
		return models.ProvisionedServiceSpec{}, fmt.Errorf("Error saving instance details to database: %s. WARNING: this instance cannot be deprovisioned through cf. Contact your operator for cleanup", err)
	}

	// save provision request details
	pr := models.ProvisionRequestDetails{
		ServiceInstanceId: instanceID,
		RequestDetails:    string(details.RawParameters),
	}
	if err = db_service.DbConnection.Create(&pr).Error; err != nil {
		return models.ProvisionedServiceSpec{}, fmt.Errorf("Error saving provision request details to database: %s. Services relying on async provisioning will not be able to complete provisioning", err)
	}

	return models.ProvisionedServiceSpec{IsAsync: gcpBroker.ShouldProvisionAsync, DashboardURL: ""}, nil
}

// cf delete-service
// Deletes the given instance
func (gcpBroker *GCPAsyncServiceBroker) Deprovision(instanceID string, details models.DeprovisionDetails, asyncAllowed bool) (models.IsAsync, error) {

	gcpBroker.ShouldProvisionAsync = gcpBroker.ServiceBrokerMap[details.ServiceID].DeprovisionsAsync()

	// make sure that instance actually exists
	count, err := db_service.GetServiceInstanceCount(instanceID)
	if err != nil {
		return models.IsAsync(gcpBroker.ShouldProvisionAsync), fmt.Errorf("Database error checking for existing instance: %s", err)
	}
	if count == 0 {
		return models.IsAsync(gcpBroker.ShouldProvisionAsync), models.ErrInstanceDoesNotExist
	}

	// if async provisioning isn't allowed but this service needs it, throw an error
	if gcpBroker.ShouldProvisionAsync && !asyncAllowed {
		return models.IsAsync(asyncAllowed), models.ErrAsyncRequired
	}

	// deprovision
	err = gcpBroker.ServiceBrokerMap[details.ServiceID].Deprovision(instanceID, details)
	if err != nil {
		return models.IsAsync(gcpBroker.ShouldProvisionAsync), err
	}

	// soft-delete instance details from the db if this is a synchronous operation
	// if it's an async operation we can't delete from the db until we're sure delete succeeded, so this is
	// handled internally to LastOperation
	if !gcpBroker.ShouldProvisionAsync {
		err = db_service.SoftDeleteInstanceDetails(instanceID)
		if err != nil {
			return models.IsAsync(gcpBroker.ShouldProvisionAsync), fmt.Errorf("Error deleting instance details from database: %s. WARNING: this instance will remain visible in cf. Contact your operator for cleanup", err)
		}
	}

	return models.IsAsync(gcpBroker.ShouldProvisionAsync), nil
}

// cf bind-service
// for cloudSql instances, Bind creates a new user and ssl cert
// for all other services, Bind creates a new service account with the IAM role listed in details.Parameters["permissions"]
// a complete list of IAM roles is available here: https://cloud.google.com/iam/docs/understanding-roles
func (gcpBroker *GCPServiceBroker) Bind(instanceID string, bindingID string, details models.BindDetails) (models.Binding, error) {

	serviceId := details.ServiceID

	// check for existing binding

	var count int
	var err error

	if err = db_service.DbConnection.Model(&models.ServiceBindingCredentials{}).Where("service_instance_id = ? and binding_id = ?", instanceID, bindingID).Count(&count).Error; err != nil {
		return models.Binding{}, fmt.Errorf("Error checking for existing binding: %s", err)
	}
	if count > 0 {
		return models.Binding{}, models.ErrBindingAlreadyExists
	}

	// create binding
	newCreds, err := gcpBroker.ServiceBrokerMap[serviceId].Bind(instanceID, bindingID, details)
	if err != nil {
		return models.Binding{}, err
	}

	// save binding to database
	newCreds.ServiceInstanceId = instanceID
	newCreds.BindingId = bindingID
	newCreds.ServiceId = details.ServiceID

	if err := db_service.DbConnection.Create(&newCreds).Error; err != nil {
		return models.Binding{}, fmt.Errorf("Error saving credentials to database: %s. WARNING: these credentials cannot be unbound through cf. Please contact your operator for cleanup",
			err)
	}

	var creds map[string]string
	if err := json.Unmarshal([]byte(newCreds.OtherDetails), &creds); err != nil {
		return models.Binding{}, err
	}

	// copy provision.otherDetails to creds.
	var instanceRecord models.ServiceInstanceDetails
	if err = db_service.DbConnection.Where("id = ?", instanceID).First(&instanceRecord).Error; err != nil {
		return models.Binding{}, fmt.Errorf("Error retrieving service instance details: %s", err)
	}

	var instanceDetails map[string]string
	// if the instance has access details saved
	if instanceRecord.OtherDetails != "" {
		if err := json.Unmarshal([]byte(instanceRecord.OtherDetails), &instanceDetails); err != nil {
			return models.Binding{}, err
		}
	}

	updatedCreds := gcpBroker.ServiceBrokerMap[serviceId].BuildInstanceCredentials(creds, instanceDetails)

	return models.Binding{
		Credentials:     updatedCreds,
		SyslogDrainURL:  "",
		RouteServiceURL: "",
	}, nil
}

// cf unbind-service
// for cloudSql instances, Unbind deletes the associated user and ssl certs
// for all other services, Unbind deletes the associated service account
func (gcpBroker *GCPServiceBroker) Unbind(instanceID, bindingID string, details models.UnbindDetails) error {

	// validate existence of binding
	var count int
	existingBinding := models.ServiceBindingCredentials{}

	if err := db_service.DbConnection.Where("service_instance_id = ? and binding_id = ?", instanceID, bindingID).Find(&existingBinding).Count(&count).Error; err != nil {
		return models.ErrBindingDoesNotExist
	}

	// remove binding from google
	err := gcpBroker.ServiceBrokerMap[details.ServiceID].Unbind(existingBinding)

	if err != nil {
		return err
	}

	// remove binding from database
	if err := db_service.DbConnection.Delete(&existingBinding).Error; err != nil {
		return fmt.Errorf("Error soft-deleting credentials from database: %s. WARNING: these credentials will remain visible in cf. Contact your operator for cleanup", err)
	}

	return nil
}

// if a service is provisioned asynchronously, LastOperation is called until the provisioning attempt times out
// or success or failure is returned
func (gcpBroker *GCPServiceBroker) LastOperation(instanceID string) (models.LastOperation, error) {

	instance := models.ServiceInstanceDetails{}
	if err := db_service.DbConnection.Where("id = ?", instanceID).First(&instance).Error; err != nil {
		return models.LastOperation{}, models.ErrInstanceDoesNotExist
	}

	if gcpBroker.ServiceBrokerMap[instance.ServiceId].ProvisionsAsync() || gcpBroker.ServiceBrokerMap[instance.ServiceId].DeprovisionsAsync() {
		return gcpBroker.lastOperationAsync(instanceID, instance.ServiceId)

	} else {
		return models.LastOperation{State: models.Succeeded, Description: ""}, errors.New("Can't call LastOperation on a synchronous service")

	}

}

func (gcpBroker *GCPServiceBroker) lastOperationAsync(instanceId, serviceId string) (models.LastOperation, error) {
	done, err := gcpBroker.ServiceBrokerMap[serviceId].PollInstance(instanceId)
	if err != nil {
		// this is a retryable error
		if gerr, ok := err.(*googleapi.Error); ok {
			if gerr.Code == 503 {
				return models.LastOperation{State: models.InProgress, Description: ""}, err
			}
		}
		// This is not a retryable error. Return fail
		return models.LastOperation{State: models.Failed, Description: ""}, err
	}

	if done {
		// no error and we're done! Delete from the SB database if this was a delete flow and return success
		deleteFlow, err := gcpBroker.ServiceBrokerMap[serviceId].LastOperationWasDelete(instanceId)
		if err != nil {
			return models.LastOperation{State: models.Succeeded, Description: ""}, fmt.Errorf("Couldn't determine if provision or deprovision flow, this may leave orphaned resources, contact your operator for cleanup")
		}
		if deleteFlow {
			err = db_service.SoftDeleteInstanceDetails(instanceId)
			if err != nil {
				return models.LastOperation{State: models.Succeeded, Description: ""}, fmt.Errorf("Error deleting instance details from database: %s. WARNING: this instance will remain visible in cf. Contact your operator for cleanup", err)
			}
		}
		return models.LastOperation{State: models.Succeeded, Description: ""}, nil
	} else {
		return models.LastOperation{State: models.InProgress, Description: ""}, nil
	}
}

// updates a service instance plan. This functionality is not implemented and will return an error indicating that plan
// changes are not supported.
func (gcpBroker *GCPServiceBroker) Update(instanceID string, details models.UpdateDetails, asyncAllowed bool) (models.IsAsync, error) {
	return models.IsAsync(asyncAllowed), models.ErrPlanChangeNotSupported
}

// reads the service account json string from the environment variable ROOT_SERVICE_ACCOUNT_JSON, writes it to a file,
// and then exports the file location to the environment variable GOOGLE_APPLICATION_CREDENTIALS, making it visible to
// all google cloud apis
func GetCredentialsFromEnv() (models.GCPCredentials, error) {
	var err error
	g := models.GCPCredentials{}

	rootCreds := os.Getenv(models.RootSaEnvVar)
	if err = json.Unmarshal([]byte(rootCreds), &g); err != nil {
		return models.GCPCredentials{}, fmt.Errorf("Error unmarshalling service account json: %s", err)
	}

	return g, nil
}

// pulls SERVICES, PLANS, and PRECONFIGURED_PLANS environment variables to construct catalog
func InitCatalogFromEnv() ([]models.Service, error) {

	// set up services
	var serviceList []models.Service

	for _, varname := range models.ServiceEnvVarNames {

		var svc models.Service
		if err := json.Unmarshal([]byte(os.Getenv(varname)), &svc); err != nil {
			return []models.Service{}, err
		} else {
			if errs := validator.Validate(svc); errs != nil {
				return []models.Service{}, errs
			} else {
				serviceList = append(serviceList, svc)
			}
		}

	}

	return serviceList, nil
}
