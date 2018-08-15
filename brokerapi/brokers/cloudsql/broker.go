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

package cloudsql

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/account_managers"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/models"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/name_generator"
	"github.com/GoogleCloudPlatform/gcp-service-broker/db_service"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/broker"
	"github.com/pivotal-cf/brokerapi"

	"context"

	"code.cloudfoundry.org/lager"
	"golang.org/x/oauth2/jwt"
	googlecloudsql "google.golang.org/api/sqladmin/v1beta4"
)

type CloudSQLBroker struct {
	HttpConfig       *jwt.Config
	ProjectId        string
	Logger           lager.Logger
	AccountManager   models.AccountManager
	SaAccountManager *account_managers.ServiceAccountManager
}

type InstanceInformation struct {
	InstanceName string `json:"instance_name"`
	DatabaseName string `json:"database_name"`
	Host         string `json:"host"`
	Region       string `json:"region"`

	LastMasterOperationId string `json:"last_master_operation_id"`
}

const SecondGenPricingPlan string = "PER_USE"
const PostgresDefaultVersion string = "POSTGRES_9_6"
const MySqlFirstGenDefaultVersion string = "MYSQL_5_6"
const MySqlSecondGenDefaultVersion string = "MYSQL_5_7"

// Creates a new CloudSQL instance
//
// optional custom parameters:
//   - database_name (generated and returned in ServiceInstanceDetails.OtherDetails.DatabaseName)
//   - instance_name (generated and returned in ServiceInstanceDetails.Name if not provided)
//   - version (defaults to 5.6 for 1st gen mysql, 5.7 for 2nd gen mysql, 9.6 for postgres)
//   - disk_size in GB (only for 2nd gen, defaults to 10)
//   - region (defaults to us-central)
//   - zone (for 2nd gen)
//   - disk_type (for 2nd gen, defaults to ssd)
//   - failover_replica_name (only for 2nd gen, if specified creates a failover replica, defaults to "")
//   - maintenance_window_day (for 2nd gen only)
//   - defaults to 1 (Sunday))
//   - maintenance_window_hour (for 2nd gen only, defaults to 0)
//   - backups_enabled (defaults to true)
//   - backup_start_time (defaults to 06:00)
//   - binlog (defaults to false for 1st gen, true for 2nd gen)
//   - activation_policy (defaults to on demand)
//   - replication_type (defaults to synchronous)
//   - auto_resize (2nd gen only, defaults to false)
//
// for more information, see: https://cloud.google.com/sql/docs/admin-api/v1beta4/instances/insert
func (b *CloudSQLBroker) Provision(instanceId string, details brokerapi.ProvisionDetails, plan models.ServicePlan) (models.ServiceInstanceDetails, error) {
	// validate parameters
	var params map[string]string
	var err error

	// validate parameters

	if len(details.RawParameters) == 0 {
		params = map[string]string{}
	} else if err = json.Unmarshal(details.RawParameters, &params); err != nil {
		return models.ServiceInstanceDetails{}, fmt.Errorf("Error unmarshalling parameters: %s", err)
	}

	if v, ok := params["instance_name"]; !ok || v == "" {
		params["instance_name"] = name_generator.Sql.InstanceName()
	}

	instanceName := params["instance_name"]

	// set default parameters or cast strings to proper values
	firstGenTiers := []string{"d0", "d1", "d2", "d4", "d8", "d16", "d32"}
	isFirstGen := false
	for _, a := range firstGenTiers {
		if a == strings.ToLower(plan.ServiceProperties["tier"]) {
			isFirstGen = true
		}
	}

	var binlogEnabled = false

	svc, err := broker.GetServiceById(details.ServiceID)
	if err != nil {
		return models.ServiceInstanceDetails{}, err
	}

	_, versionOk := params["version"]
	// set default parameters or cast strings to proper values
	if svc.Name == models.CloudsqlPostgresName {
		if !versionOk {
			params["version"] = PostgresDefaultVersion
		}
	} else {
		if !versionOk {
			params["version"] = MySqlFirstGenDefaultVersion
		}

		if !isFirstGen {
			binlogEnabled = true
			if !versionOk {
				params["version"] = MySqlSecondGenDefaultVersion
			}
		}
		binlog, binlogOk := params["binlog"]
		if binlogOk {
			binlogEnabled, err = strconv.ParseBool(binlog)
			if err != nil {
				return models.ServiceInstanceDetails{}, fmt.Errorf("%s is not a valid value for binlog", binlog)
			}
		}
	}

	openAcls := []*googlecloudsql.AclEntry{}
	aclsParamDetails, aclsParamOk := params["authorized_networks"]
	if aclsParamOk && aclsParamDetails != "" {
		authorizedNetworks := strings.Split(aclsParamDetails, ",")
		for _, v := range authorizedNetworks {
			openAcl := googlecloudsql.AclEntry{
				Value: v,
			}
			openAcls = append(openAcls, &openAcl)
		}
	}

	backupsEnabled := true
	if params["backups_enabled"] == "false" {
		backupsEnabled = false
	}

	backupStartTime := "06:00"
	if startTime, ok := params["backup_start_time"]; ok {
		backupStartTime = startTime
	}

	var di googlecloudsql.DatabaseInstance
	if isFirstGen {

		di = createFirstGenRequest(plan.ServiceProperties, params)
	} else {
		di, err = createInstanceRequest(plan.ServiceProperties, params)
		if err != nil {
			return models.ServiceInstanceDetails{}, err
		}

	}
	di.Name = instanceName
	di.Settings.BackupConfiguration = &googlecloudsql.BackupConfiguration{
		Enabled:          backupsEnabled,
		StartTime:        backupStartTime,
		BinaryLogEnabled: binlogEnabled,
	}
	di.Settings.IpConfiguration.AuthorizedNetworks = openAcls

	// init sqladmin service
	sqlService, err := googlecloudsql.New(b.HttpConfig.Client(context.Background()))
	if err != nil {
		return models.ServiceInstanceDetails{}, fmt.Errorf("Error creating new CloudSQL Client: %s", err)
	}
	sqlService.UserAgent = models.CustomUserAgent

	// make insert request
	op, err := sqlService.Instances.Insert(b.ProjectId, &di).Do()
	if err != nil {
		return models.ServiceInstanceDetails{}, fmt.Errorf("Error creating new CloudSQL instance: %s", err)
	}

	// save new cloud operation
	if err = createCloudOperation(op, instanceId, details.ServiceID); err != nil {
		return models.ServiceInstanceDetails{}, err
	}

	// update instance information on instancedetails object
	ii := InstanceInformation{
		InstanceName:          instanceName,
		LastMasterOperationId: op.Name,
	}

	otherDetails, err := json.Marshal(ii)
	if err != nil {
		return models.ServiceInstanceDetails{}, fmt.Errorf("Error marshalling instance information: %s", err)
	}
	b.Logger.Debug(fmt.Sprintf("UPDATING OTHER DETAILS FROM %v to %s", "nothing", string(otherDetails)))
	i := models.ServiceInstanceDetails{
		Name:         params["instance_name"],
		Url:          "",
		Location:     "",
		OtherDetails: string(otherDetails),
	}

	return i, nil

}

func createFirstGenRequest(planDetails, params map[string]string) googlecloudsql.DatabaseInstance {
	// set up instance resource
	return googlecloudsql.DatabaseInstance{
		Settings: &googlecloudsql.Settings{
			IpConfiguration: &googlecloudsql.IpConfiguration{
				RequireSsl:  true,
				Ipv4Enabled: true,
			},
			Tier:             planDetails["tier"],
			PricingPlan:      planDetails["pricing_plan"],
			ActivationPolicy: params["activation_policy"],
			ReplicationType:  params["replication_type"],
		},
		DatabaseVersion: params["version"],
		Region:          params["region"],
	}
}

func createInstanceRequest(planDetails, params map[string]string) (googlecloudsql.DatabaseInstance, error) {
	var err error

	diskSize, err := getDiskSize(params, planDetails)
	if err != nil {
		return googlecloudsql.DatabaseInstance{}, err
	}

	mw, err := getMaintenanceWindow(params)
	if err != nil {
		return googlecloudsql.DatabaseInstance{}, err
	}

	autoResize := false
	if params["auto_resize"] == "true" {
		autoResize = true
	}

	// set up instance resource
	return googlecloudsql.DatabaseInstance{
		Settings: &googlecloudsql.Settings{
			IpConfiguration: &googlecloudsql.IpConfiguration{
				RequireSsl:  true,
				Ipv4Enabled: true,
			},
			Tier:           planDetails["tier"],
			DataDiskSizeGb: diskSize,
			LocationPreference: &googlecloudsql.LocationPreference{
				Zone: params["zone"],
			},
			DataDiskType:      params["disk_type"],
			MaintenanceWindow: mw,
			PricingPlan:       SecondGenPricingPlan,
			ActivationPolicy:  params["activation_policy"],
			ReplicationType:   params["replication_type"],
			StorageAutoResize: &autoResize,
		},
		DatabaseVersion: params["version"],
		Region:          params["region"],
		FailoverReplica: &googlecloudsql.DatabaseInstanceFailoverReplica{
			Name: params["failover_replica_name"],
		},
	}, nil
}

func getMaintenanceWindow(params map[string]string) (*googlecloudsql.MaintenanceWindow, error) {
	var mw *googlecloudsql.MaintenanceWindow
	day, dayOk := params["maintenance_window_day"]
	hour, hourOk := params["maintenance_window_hour"]
	if dayOk && hourOk {
		intDay, err := strconv.Atoi(day)
		if err != nil {
			return &googlecloudsql.MaintenanceWindow{}, fmt.Errorf("Error converting maintenance_window_day string to int: %s", err)
		}
		intHour, err := strconv.Atoi(hour)
		if err != nil {
			return &googlecloudsql.MaintenanceWindow{}, fmt.Errorf("Error converting maintenance_window_hour string to int: %s", err)
		}
		mw = &googlecloudsql.MaintenanceWindow{
			Day:         int64(intDay),
			Hour:        int64(intHour),
			UpdateTrack: "stable",
		}
	}
	return mw, nil
}

func getDiskSize(params, planDetails map[string]string) (int64, error) {
	var err error
	diskSize := 10
	if _, diskSizeOk := params["disk_size"]; diskSizeOk {
		diskSize, err = strconv.Atoi(params["disk_size"])
		if err != nil {
			return 0, fmt.Errorf("Error converting disk_size gb string to int: %s", err)
		}
	}
	maxDiskSize, err := strconv.Atoi(planDetails["max_disk_size"])
	if err != nil {
		return 0, fmt.Errorf("Error converting max_disk_size gb string to int: %s", err)
	}
	if diskSize > maxDiskSize {
		return 0, errors.New("disk size is greater than max allowed disk size for this plan")
	}
	return int64(diskSize), nil
}

// Completes the second step in provisioning a cloudsql instance, namely, creating the db.
// executing this "synchronously" even though technically db creation returns an op - but it's just a db call, so
// it should be quick and not actually async.
func (b *CloudSQLBroker) FinishProvisioning(instanceId string, params map[string]string) error {
	var err error

	instance := models.ServiceInstanceDetails{}
	if err = db_service.DbConnection.Where("ID = ?", instanceId).First(&instance).Error; err != nil {
		return brokerapi.ErrInstanceDoesNotExist
	}

	sqlService, err := googlecloudsql.New(b.HttpConfig.Client(context.Background()))
	if err != nil {
		return fmt.Errorf("Error creating new CloudSQL Client: %s", err)
	}

	dbService := googlecloudsql.NewInstancesService(sqlService)
	clouddb, err := dbService.Get(b.ProjectId, instance.Name).Do()
	if err != nil {
		return fmt.Errorf("Error getting instance from api: %s", err)
	}

	//create actual database entry

	if v, ok := params["database_name"]; !ok || v == "" {
		params["database_name"] = name_generator.Sql.DatabaseName()
	}

	d := googlecloudsql.Database{
		Name: params["database_name"],
	}

	op, err := sqlService.Databases.Insert(b.ProjectId, clouddb.Name, &d).Do()
	if err != nil {
		return fmt.Errorf("Error creating database: %s", err)
	}

	// Create new operation entry for the database insert
	if err = createCloudOperation(op, instanceId, instance.ServiceId); err != nil {
		return err
	}

	// Save new operation id and database name to instance data
	if err = updateOperationId(instance, op.Name); err != nil {
		return err
	}

	//poll for the database creation operation to be completed
	// TODO(cbriant): consider changing this. It isn't strictly needed, though it is unlikely to hurt either.
	err = b.pollOperationUntilDone(op, b.ProjectId)
	// XXX: return this error exactly as is from the google api
	if err != nil {
		return err
	}

	// update db information
	instance.Url = clouddb.SelfLink
	instance.Location = clouddb.Region

	// update instance information
	var ii InstanceInformation
	if err := json.Unmarshal([]byte(instance.OtherDetails), &ii); err != nil {
		return fmt.Errorf("Error unmarshalling instance information.")
	}

	ii.Host = clouddb.IpAddresses[0].IpAddress
	ii.DatabaseName = params["database_name"]
	ii.Region = instance.Location
	otherDetails, err := json.Marshal(ii)
	if err != nil {
		return fmt.Errorf("Error marshalling instance information: %s.", err)
	}
	b.Logger.Debug(fmt.Sprintf("UPDATING OTHER DETAILS FROM %v to %s", instance.OtherDetails, string(otherDetails)))
	instance.OtherDetails = string(otherDetails)

	if err = db_service.DbConnection.Save(&instance).Error; err != nil {
		return fmt.Errorf(`Error saving instance details to database: %s. WARNING: this instance cannot be deprovisioned through cf.
		Please contact your operator for cleanup`, err)
	}

	return nil
}

// generate a new username, password if not provided
func (b *CloudSQLBroker) ensureUsernamePassword(instanceID, bindingID string, details *brokerapi.BindDetails) error {
	if details.RawParameters == nil {
		details.RawParameters = []byte("{}")
	}

	tempParams := map[string]interface{}{}
	err := json.Unmarshal(details.RawParameters, &tempParams)
	if err != nil {
		return err
	}

	if v, ok := tempParams["username"].(string); !ok || v == "" {
		username, err := name_generator.Sql.GenerateUsername(instanceID, bindingID)
		if err != nil {
			return err
		}
		tempParams["username"] = username
	}
	if v, ok := tempParams["password"].(string); !ok || v == "" {
		password, err := name_generator.Sql.GeneratePassword()
		if err != nil {
			return err
		}
		tempParams["password"] = password
	}

	details.RawParameters, err = json.Marshal(tempParams)
	return err
}

// creates a new username, password, and set of ssl certs for the given instance
// may be slow to return because CloudSQL operations are async. Timeout may need to be raised to 90 or 120 seconds
func (b *CloudSQLBroker) Bind(instanceID, bindingID string, details brokerapi.BindDetails) (models.ServiceBindingCredentials, error) {

	cloudDb := models.ServiceInstanceDetails{}
	if err := db_service.DbConnection.Where("ID = ?", instanceID).First(&cloudDb).Error; err != nil {
		return models.ServiceBindingCredentials{}, brokerapi.ErrInstanceDoesNotExist
	}

	if err := b.ensureUsernamePassword(instanceID, bindingID, &details); err != nil {
		return models.ServiceBindingCredentials{}, err
	}

	sqlCredBytes, err := b.AccountManager.CreateCredentials(instanceID, bindingID, details, cloudDb)
	if err != nil {
		return models.ServiceBindingCredentials{}, err
	}

	saCredBytes, err := b.SaAccountManager.CreateCredentials(instanceID, bindingID, details, models.ServiceInstanceDetails{})
	if err != nil {
		return models.ServiceBindingCredentials{}, err
	}

	credsJSON, err := combineServiceBindingCreds(sqlCredBytes, saCredBytes)

	if err != nil {
		return models.ServiceBindingCredentials{}, err
	}

	params := make(map[string]interface{})
	if err := json.Unmarshal(details.RawParameters, &params); err != nil {
		return models.ServiceBindingCredentials{}, fmt.Errorf("Error unmarshalling parameters: %s", err)
	}

	jdbcUriFormat, jdbcUriFormatOk := params["jdbc_uri_format"].(string)
	credsJSON["UriPrefix"] = ""
	if jdbcUriFormatOk && jdbcUriFormat == "true" {
		credsJSON["UriPrefix"] = "jdbc:"
	}

	credBytes, err := json.Marshal(&credsJSON)

	if err != nil {
		return models.ServiceBindingCredentials{}, err
	}

	newBinding := models.ServiceBindingCredentials{
		OtherDetails: string(credBytes),
	}

	return newBinding, nil
}

func combineServiceBindingCreds(sqlCreds models.ServiceBindingCredentials, saCreds models.ServiceBindingCredentials) (map[string]string, error) {
	var sqlCredsJSON map[string]string

	if err := json.Unmarshal([]byte(sqlCreds.OtherDetails), &sqlCredsJSON); err != nil {
		return map[string]string{}, err
	}

	var saCredsJSON map[string]string

	if err := json.Unmarshal([]byte(saCreds.OtherDetails), &saCredsJSON); err != nil {
		return map[string]string{}, err
	}

	sqlCredsJSON["PrivateKeyData"] = saCredsJSON["PrivateKeyData"]
	sqlCredsJSON["ProjectId"] = saCredsJSON["ProjectId"]
	sqlCredsJSON["Email"] = saCredsJSON["Email"]
	sqlCredsJSON["UniqueId"] = saCredsJSON["UniqueId"]

	return sqlCredsJSON, nil
}

func (b *CloudSQLBroker) BuildInstanceCredentials(bindDetails models.ServiceBindingCredentials, instanceDetails models.ServiceInstanceDetails) (map[string]string, error) {
	return b.AccountManager.BuildInstanceCredentials(bindDetails, instanceDetails)
}

// Deletes the user and invalidates the ssl certs associated with this binding
// CloudFoundry doesn't seem to support async unbinding, so hopefully this works all the time even though technically
// some of these operations are async
func (b *CloudSQLBroker) Unbind(creds models.ServiceBindingCredentials) error {

	err := b.AccountManager.DeleteCredentials(creds)

	if err != nil {
		return err
	}

	err = b.SaAccountManager.DeleteCredentials(creds)

	if err != nil {
		return err
	}

	return nil
}

// gets the last operation for this instance and polls the status of it
func (b *CloudSQLBroker) PollInstance(instanceId string) (bool, error) {
	var instance models.ServiceInstanceDetails

	if err := db_service.DbConnection.Where("id = ?", instanceId).First(&instance).Error; err != nil {
		return false, brokerapi.ErrInstanceDoesNotExist
	}

	op, err := db_service.GetLastOperation(instanceId)
	if err != nil {
		return false, err
	}

	return b.PollOperation(instance, op)
}

// Checks the status of the given CloudSQL operation and determines if it is ready to continue provisioning. If so,
// calls the finish method and returns a bool indicating if provisioning is done or not, and any error
// TODO(cbriant): at least rename, if not restructure, this function
// XXX: note that for this function in particular, we are being explicit to return errors from the google api exactly
// as we get them, because further up the stack these errors will be evaluated differently and need to be preserved
func (b *CloudSQLBroker) PollOperation(instance models.ServiceInstanceDetails, op models.CloudOperation) (bool, error) {

	var err error

	// create operation service
	sqlService, err := googlecloudsql.New(b.HttpConfig.Client(context.Background()))
	if err != nil {
		return false, err
	}

	opsService := googlecloudsql.NewOperationsService(sqlService)

	// get the status of the operation
	opStatus, err := opsService.Get(b.ProjectId, op.Name).Do()
	if err != nil {
		return false, err
	}

	// update the operation status if it has changed
	if op.Status != opStatus.Status {
		op.Status = opStatus.Status
		var opErr string
		if opStatus.Error != nil {
			opErrBytes, _ := opStatus.Error.MarshalJSON()
			opErr = string(opErrBytes)
		} else {
			opErr = ""
		}
		op.ErrorMessage = string(opErr)
		db_service.DbConnection.Save(&op)
	}

	// we were provisioning and finished the first step
	if opStatus.Status == "DONE" && op.OperationType == "CREATE" {
		pr := models.ProvisionRequestDetails{}
		if err = db_service.DbConnection.Where("service_instance_id = ?", instance.ID).First(&pr).Error; err != nil {
			return false, brokerapi.ErrInstanceDoesNotExist
		}

		var pd map[string]string
		if len(pr.RequestDetails) == 0 {
			pd = map[string]string{}
		} else if err = json.Unmarshal([]byte(pr.RequestDetails), &pd); err != nil {
			return false, fmt.Errorf("Error unmarshalling request details: %s", err)
		}

		// XXX: return error exactly as is from google api
		err = b.FinishProvisioning(instance.ID, pd)
		if err != nil {
			return false, err
		}

	}

	return opStatus.Status == "DONE", nil
}

// loops and waits until a cloudsql operation is done, returns an error if any is encountered
// XXX: note that for this function in particular, we are being explicit to return errors from the google api exactly
// as we get them, because further up the stack these errors will be evaluated differently and need to be preserved
func (b *CloudSQLBroker) pollOperationUntilDone(op *googlecloudsql.Operation, projectId string) error {
	sqlService, err := googlecloudsql.New(b.HttpConfig.Client(context.Background()))
	opsService := googlecloudsql.NewOperationsService(sqlService)
	done := false
	for done == false {
		status, err := opsService.Get(projectId, op.Name).Do()
		if err != nil {
			return err
		}
		if status.EndTime != "" {
			done = true
		} else {
			println("still waiting for it to be done")
		}
		// sleep for 1 second between polling so we don't hit our rate limit
		time.Sleep(time.Second)
	}
	return err
}

// issue a delete call on the database instance
func (b *CloudSQLBroker) Deprovision(instanceId string, details brokerapi.DeprovisionDetails) error {
	var err error

	// get the service instnace object
	instance := models.ServiceInstanceDetails{}
	if err = db_service.DbConnection.Where("ID = ?", instanceId).First(&instance).Error; err != nil {
		return brokerapi.ErrInstanceDoesNotExist
	}

	sqlService, err := googlecloudsql.New(b.HttpConfig.Client(context.Background()))
	if err != nil {
		return fmt.Errorf("Error creating CloudSQL client: %s", err)
	}

	// delete the instance from google
	op, err := sqlService.Instances.Delete(b.ProjectId, instance.Name).Do()
	if err != nil {
		return fmt.Errorf("Error deleting instance: %s", err)
	}

	// update the service instance state (other details)
	if err = createCloudOperation(op, instanceId, details.ServiceID); err != nil {
		return err
	}

	// Save new operation id to instance data
	if err = updateOperationId(instance, op.Name); err != nil {
		return err
	}

	return nil
}

func createCloudOperation(op *googlecloudsql.Operation, instanceId string, serviceId string) error {
	var err error
	var opErr []byte

	if op.Error != nil {
		opErr, err = op.Error.MarshalJSON()
		if err != nil {
			return fmt.Errorf("Error marshalling operation error: %s", err)
		}
	}

	currentState := models.CloudOperation{
		Name:              op.Name,
		ErrorMessage:      string(opErr),
		InsertTime:        op.InsertTime,
		OperationType:     op.OperationType,
		StartTime:         op.StartTime,
		Status:            op.Status,
		TargetId:          op.TargetId,
		TargetLink:        op.TargetLink,
		ServiceId:         serviceId,
		ServiceInstanceId: instanceId,
	}

	if err = db_service.DbConnection.Create(&currentState).Error; err != nil {
		return fmt.Errorf("Error saving operation details to database: %s. Services relying on async deprovisioning will not be able to complete deprovisioning", err)
	}
	return nil
}

func updateOperationId(instance models.ServiceInstanceDetails, operationId string) error {
	var ii InstanceInformation
	if err := json.Unmarshal([]byte(instance.OtherDetails), &ii); err != nil {
		return fmt.Errorf("Error unmarshalling instance information.")
	}
	ii.LastMasterOperationId = operationId

	otherDetails, err := json.Marshal(ii)
	if err != nil {
		return fmt.Errorf("Error marshalling instance information: %s.", err)
	}
	instance.OtherDetails = string(otherDetails)

	if err = db_service.DbConnection.Save(&instance).Error; err != nil {
		return fmt.Errorf(`Error saving instance details to database: %s. WARNING: this instance cannot be deprovisioned through cf.
		Please contact your operator for cleanup`, err)
	}
	return nil
}

// used during polling of async operations to determine if the workflow is a provision or deprovision flow based off the
// type of the most recent operation
func (b *CloudSQLBroker) LastOperationWasDelete(instanceId string) (bool, error) {
	op, err := db_service.GetLastOperation(instanceId)
	if err != nil {
		return false, err
	}
	return op.OperationType == "DELETE", nil
}

// Indicates that CloudSQL uses asynchronous provisioning
func (b *CloudSQLBroker) ProvisionsAsync() bool {
	return true
}

func (b *CloudSQLBroker) DeprovisionsAsync() bool {
	return true
}
