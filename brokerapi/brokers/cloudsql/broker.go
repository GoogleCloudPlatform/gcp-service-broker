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

	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/models"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/name_generator"
	"github.com/GoogleCloudPlatform/gcp-service-broker/db_service"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/broker"
	"github.com/GoogleCloudPlatform/gcp-service-broker/utils"
	"github.com/pivotal-cf/brokerapi"

	"context"

	"code.cloudfoundry.org/lager"
	multierror "github.com/hashicorp/go-multierror"
	"golang.org/x/oauth2/jwt"
	googlecloudsql "google.golang.org/api/sqladmin/v1beta4"
)

const (
	secondGenPricingPlan         string = "PER_USE"
	postgresDefaultVersion       string = "POSTGRES_9_6"
	mySqlFirstGenDefaultVersion  string = "MYSQL_5_6"
	mySqlSecondGenDefaultVersion string = "MYSQL_5_7"
)

// CloudSQLBroker is the service-broker back-end for creating and binding CloudSQL instances.
type CloudSQLBroker struct {
	HttpConfig       *jwt.Config
	ProjectId        string
	Logger           lager.Logger
	AccountManager   models.AccountManager
	SaAccountManager models.AccountManager
}

// InstanceInformation holds the details needed to bind a service account to a CloudSQL instance after it has been provisioned.
type InstanceInformation struct {
	InstanceName string `json:"instance_name"`
	DatabaseName string `json:"database_name"`
	Host         string `json:"host"`
	Region       string `json:"region"`

	LastMasterOperationId string `json:"last_master_operation_id"`
}

// Provision creates a new CloudSQL instance from the settings in the user-provided details and service plan.
func (b *CloudSQLBroker) Provision(ctx context.Context, instanceId string, details brokerapi.ProvisionDetails, plan models.ServicePlan) (models.ServiceInstanceDetails, error) {
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
	labels := utils.ExtractDefaultLabels(instanceId, details)

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
			params["version"] = postgresDefaultVersion
		}
	} else {
		if !versionOk {
			params["version"] = mySqlFirstGenDefaultVersion
		}

		if !isFirstGen {
			binlogEnabled = true
			if !versionOk {
				params["version"] = mySqlSecondGenDefaultVersion
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

		di = createFirstGenRequest(plan.ServiceProperties, params, labels)
	} else {
		di, err = createInstanceRequest(plan.ServiceProperties, params, labels)
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
	sqlService, err := b.createClient(ctx)
	if err != nil {
		return models.ServiceInstanceDetails{}, err
	}

	// make insert request
	op, err := sqlService.Instances.Insert(b.ProjectId, &di).Do()
	if err != nil {
		return models.ServiceInstanceDetails{}, fmt.Errorf("Error creating new CloudSQL instance: %s", err)
	}

	if v, ok := params["database_name"]; !ok || v == "" {
		params["database_name"] = name_generator.Sql.DatabaseName()
	}

	// update instance information on instancedetails object
	ii := InstanceInformation{
		InstanceName: instanceName,
		DatabaseName: params["database_name"],
	}

	otherDetails, err := json.Marshal(ii)
	if err != nil {
		return models.ServiceInstanceDetails{}, fmt.Errorf("Error marshalling instance information: %s", err)
	}
	b.Logger.Debug("updating details", lager.Data{"from": "{}", "to": otherDetails})
	i := models.ServiceInstanceDetails{
		Name:         params["instance_name"],
		Url:          "",
		Location:     "",
		OtherDetails: string(otherDetails),

		OperationType: models.ProvisionOperationType,
		OperationId:   op.Name,
	}

	return i, nil

}

func createFirstGenRequest(planDetails, params, labels map[string]string) googlecloudsql.DatabaseInstance {
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
			UserLabels:       labels,
		},
		DatabaseVersion: params["version"],
		Region:          params["region"],
	}
}

func createInstanceRequest(planDetails, params, labels map[string]string) (googlecloudsql.DatabaseInstance, error) {
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
			PricingPlan:       secondGenPricingPlan,
			ActivationPolicy:  params["activation_policy"],
			ReplicationType:   params["replication_type"],
			StorageAutoResize: &autoResize,
			UserLabels:        labels,
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

// Bind creates a new username, password, and set of ssl certs for the given instance.
// The function may be slow to return because CloudSQL operations are asynchronous.
// The default PCF service broker timeout may need to be raised to 90 or 120 seconds to accommodate the long bind time.
func (b *CloudSQLBroker) Bind(ctx context.Context, instanceID, bindingID string, details brokerapi.BindDetails) (models.ServiceBindingCredentials, error) {

	cloudDb, err := db_service.GetServiceInstanceDetailsById(ctx, instanceID)
	if err != nil {
		return models.ServiceBindingCredentials{}, brokerapi.ErrInstanceDoesNotExist
	}

	if err := b.ensureUsernamePassword(instanceID, bindingID, &details); err != nil {
		return models.ServiceBindingCredentials{}, err
	}

	sqlCredBytes, err := b.AccountManager.CreateCredentials(ctx, instanceID, bindingID, details, *cloudDb)
	if err != nil {
		return models.ServiceBindingCredentials{}, err
	}

	saCredBytes, err := b.SaAccountManager.CreateCredentials(ctx, instanceID, bindingID, details, models.ServiceInstanceDetails{})
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
	sqlCredsJSON, err := sqlCreds.GetOtherDetails()
	if err != nil {
		return map[string]string{}, err
	}

	saCredsJSON, err := saCreds.GetOtherDetails()
	if err != nil {
		return map[string]string{}, err
	}

	sqlCredsJSON["PrivateKeyData"] = saCredsJSON["PrivateKeyData"]
	sqlCredsJSON["ProjectId"] = saCredsJSON["ProjectId"]
	sqlCredsJSON["Email"] = saCredsJSON["Email"]
	sqlCredsJSON["UniqueId"] = saCredsJSON["UniqueId"]

	return sqlCredsJSON, nil
}

func (b *CloudSQLBroker) BuildInstanceCredentials(ctx context.Context, bindDetails models.ServiceBindingCredentials, instanceDetails models.ServiceInstanceDetails) (map[string]string, error) {
	return b.AccountManager.BuildInstanceCredentials(ctx, bindDetails, instanceDetails)
}

// Unbind deletes the database user, service account and invalidates the ssl certs associated with this binding.
func (b *CloudSQLBroker) Unbind(ctx context.Context, creds models.ServiceBindingCredentials) error {
	var accumulator error

	if err := b.AccountManager.DeleteCredentials(ctx, creds); err != nil {
		accumulator = multierror.Append(accumulator, err)
	}

	if err := b.SaAccountManager.DeleteCredentials(ctx, creds); err != nil {
		accumulator = multierror.Append(accumulator, err)
	}

	return accumulator
}

// PollInstance gets the last operation for this instance and checks its status.
func (b *CloudSQLBroker) PollInstance(ctx context.Context, instance models.ServiceInstanceDetails) (bool, error) {
	b.Logger.Info("PollInstance", lager.Data{
		"instance":       instance.Name,
		"operation_type": instance.OperationType,
		"operation_id":   instance.OperationId,
	})

	if instance.OperationType == "" {
		return false, errors.New("Couldn't find any pending operations for this CloudSQL instance.")
	}

	result, err := b.pollOperation(ctx, instance.OperationId)
	if result == false || err != nil {
		return result, err
	}

	if instance.OperationType == models.ProvisionOperationType {
		instancePtr := &instance
		if err := b.UpdateInstanceDetails(ctx, instancePtr); err != nil {
			return true, err
		}

		return true, b.createDatabase(ctx, instancePtr)
	}

	return true, nil
}

// refreshServiceInstanceDetails fetches the settings for the instance from GCP
// and upates the provided instance with the refreshed info.
func (b *CloudSQLBroker) UpdateInstanceDetails(ctx context.Context, instance *models.ServiceInstanceDetails) error {
	var instanceInfo InstanceInformation
	if err := json.Unmarshal([]byte(instance.OtherDetails), &instanceInfo); err != nil {
		return fmt.Errorf("Error unmarshalling instance information.")
	}

	client, err := b.createClient(ctx)
	if err != nil {
		return err
	}

	clouddb, err := googlecloudsql.NewInstancesService(client).Get(b.ProjectId, instance.Name).Do()
	if err != nil {
		return fmt.Errorf("Error getting instance from API: %s", err)
	}

	// update db information
	instance.Url = clouddb.SelfLink
	instance.Location = clouddb.Region

	// update instance information
	instanceInfo.Host = clouddb.IpAddresses[0].IpAddress
	instanceInfo.Region = clouddb.Region
	otherDetails, err := json.Marshal(instanceInfo)
	if err != nil {
		return fmt.Errorf("Error marshalling instance information: %s.", err)
	}
	instance.OtherDetails = string(otherDetails)

	return nil
}

// createDatabase creates tha database on the instance referenced by ServiceInstanceDetails.
func (b *CloudSQLBroker) createDatabase(ctx context.Context, instance *models.ServiceInstanceDetails) error {
	var instanceInfo InstanceInformation
	if err := json.Unmarshal([]byte(instance.OtherDetails), &instanceInfo); err != nil {
		return fmt.Errorf("Error unmarshalling instance information.")
	}

	client, err := b.createClient(ctx)
	if err != nil {
		return err
	}

	d := googlecloudsql.Database{Name: instanceInfo.DatabaseName}
	op, err := client.Databases.Insert(b.ProjectId, instance.Name, &d).Do()
	if err != nil {
		return fmt.Errorf("Error creating database: %s", err)
	}

	// poll for the database creation operation to be completed
	// XXX: return this error exactly as is from the google api
	return b.pollOperationUntilDone(ctx, op, b.ProjectId)
}

func (b *CloudSQLBroker) pollOperation(ctx context.Context, opterationId string) (bool, error) {
	client, err := b.createClient(ctx)
	if err != nil {
		return false, err
	}

	// get the status of the operation
	operation, err := googlecloudsql.NewOperationsService(client).Get(b.ProjectId, opterationId).Do()
	if err != nil {
		return false, err
	}

	if operation.Status == "DONE" {
		if operation.Error == nil {
			return true, nil
		} else {
			errs := ""

			for _, err := range operation.Error.Errors {
				errs += fmt.Sprintf("%s: %q; ", err.Code, err.Message)
			}

			return true, errors.New(errs)
		}
	}

	return false, nil
}

// pollOperationUntilDone loops and waits until a cloudsql operation is done, returning an error if any is encountered
// XXX: note that for this function in particular, we are being explicit to return errors from the google api exactly
// as we get them, because further up the stack these errors will be evaluated differently and need to be preserved
func (b *CloudSQLBroker) pollOperationUntilDone(ctx context.Context, op *googlecloudsql.Operation, projectId string) error {
	sqlService, err := b.createClient(ctx)
	if err != nil {
		return err
	}

	opsService := googlecloudsql.NewOperationsService(sqlService)
	for {
		status, err := opsService.Get(projectId, op.Name).Do()
		if err != nil {
			return err
		}

		if status.EndTime != "" {
			return nil
		}

		b.Logger.Info("waiting for operation", lager.Data{"operation": op.Name, "status": status.Status})
		// sleep for 1 second between polling so we don't hit our rate limit
		time.Sleep(time.Second)
	}
}

// Deprovision issues a delete call on the database instance.
func (b *CloudSQLBroker) Deprovision(ctx context.Context, instance models.ServiceInstanceDetails, details brokerapi.DeprovisionDetails) error {
	sqlService, err := b.createClient(ctx)
	if err != nil {
		return err
	}

	// delete the instance from google
	op, err := sqlService.Instances.Delete(b.ProjectId, instance.Name).Do()
	if err != nil {
		return fmt.Errorf("Error deleting instance: %s", err)
	}

	instance.OperationType = models.DeprovisionOperationType
	instance.OperationId = op.Name

	if err := db_service.SaveServiceInstanceDetails(ctx, &instance); err != nil {
		return fmt.Errorf("Error saving delete instance operation: %s", err)
	}

	return nil
}

// ProvisionsAsync indicates that CloudSQL uses asynchronous provisioning.
func (b *CloudSQLBroker) ProvisionsAsync() bool {
	return true
}

// DeprovisionsAsync indicates that CloudSQL uses asynchronous deprovisioning.
func (b *CloudSQLBroker) DeprovisionsAsync() bool {
	return true
}

func (b *CloudSQLBroker) createClient(ctx context.Context) (*googlecloudsql.Service, error) {
	client, err := googlecloudsql.New(b.HttpConfig.Client(ctx))
	if err != nil {
		return nil, fmt.Errorf("Couldn't instantiate CloudSQL API client: %s", err)
	}

	client.UserAgent = models.CustomUserAgent
	return client, nil
}
