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
	"strings"
	"time"

	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/broker_base"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/models"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/name_generator"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/broker"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/varcontext"
	"github.com/GoogleCloudPlatform/gcp-service-broker/utils"
	"github.com/pivotal-cf/brokerapi"
	"github.com/spf13/cast"

	"context"

	"code.cloudfoundry.org/lager"
	multierror "github.com/hashicorp/go-multierror"
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
	broker_base.BrokerBase
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
	di, ii, err := createProvisionRequest(instanceId, details, plan)
	if err != nil {
		return models.ServiceInstanceDetails{}, err
	}

	// init sqladmin service
	sqlService, err := b.createClient(ctx)
	if err != nil {
		return models.ServiceInstanceDetails{}, err
	}

	// make insert request
	op, err := sqlService.Instances.Insert(b.ProjectId, di).Do()
	if err != nil {
		return models.ServiceInstanceDetails{}, fmt.Errorf("Error creating new CloudSQL instance: %s", err)
	}

	otherDetails, err := json.Marshal(ii)
	if err != nil {
		return models.ServiceInstanceDetails{}, fmt.Errorf("Error marshalling instance information: %s", err)
	}

	b.Logger.Debug("updating details", lager.Data{"from": "{}", "to": otherDetails})
	return models.ServiceInstanceDetails{
		Name:         di.Name,
		Url:          "",
		Location:     "",
		OtherDetails: string(otherDetails),

		OperationType: models.ProvisionOperationType,
		OperationId:   op.Name,
	}, nil
}

func createProvisionRequest(instanceId string, details brokerapi.ProvisionDetails, plan models.ServicePlan) (*googlecloudsql.DatabaseInstance, *InstanceInformation, error) {
	svc, err := broker.GetServiceById(details.ServiceID)
	if err != nil {
		return nil, nil, err
	}

	vars, err := svc.ProvisionVariables(instanceId, details, plan)
	if err != nil {
		return nil, nil, err
	}

	// set up database information
	var di *googlecloudsql.DatabaseInstance
	if vars.GetBool("is_first_gen") {
		di = createFirstGenRequest(vars)
	} else {
		di = createInstanceRequest(vars)
	}

	instanceName := vars.GetString("instance_name")

	di.Name = instanceName
	di.Settings.BackupConfiguration = &googlecloudsql.BackupConfiguration{
		Enabled:          vars.GetBool("backups_enabled"),
		StartTime:        vars.GetString("backup_start_time"),
		BinaryLogEnabled: vars.GetBool("binlog"),
	}
	di.Settings.IpConfiguration.AuthorizedNetworks = varctxGetAcls(vars)
	di.Settings.UserLabels = utils.ExtractDefaultLabels(instanceId, details)

	// Set up instance information
	ii := InstanceInformation{
		InstanceName: instanceName,
		DatabaseName: vars.GetString("database_name"),
	}

	return di, &ii, vars.Error()
}

func varctxGetAcls(vars *varcontext.VarContext) []*googlecloudsql.AclEntry {
	openAcls := []*googlecloudsql.AclEntry{}
	authorizedNetworkCsv := vars.GetString("authorized_networks")
	if authorizedNetworkCsv == "" {
		return openAcls
	}

	for _, v := range strings.Split(authorizedNetworkCsv, ",") {
		openAcls = append(openAcls, &googlecloudsql.AclEntry{Value: v})
	}

	return openAcls
}

func createFirstGenRequest(vars *varcontext.VarContext) *googlecloudsql.DatabaseInstance {
	// set up instance resource
	return &googlecloudsql.DatabaseInstance{
		Settings: &googlecloudsql.Settings{
			IpConfiguration: &googlecloudsql.IpConfiguration{
				RequireSsl:  true,
				Ipv4Enabled: true,
			},
			Tier:             vars.GetString("tier"),
			PricingPlan:      vars.GetString("pricing_plan"),
			ActivationPolicy: vars.GetString("activation_policy"),
			ReplicationType:  vars.GetString("replication_type"),
		},
		DatabaseVersion: vars.GetString("version"),
		Region:          vars.GetString("region"),
	}
}

func createInstanceRequest(vars *varcontext.VarContext) *googlecloudsql.DatabaseInstance {
	autoResize := vars.GetBool("auto_resize")

	// set up instance resource
	return &googlecloudsql.DatabaseInstance{
		Settings: &googlecloudsql.Settings{
			IpConfiguration: &googlecloudsql.IpConfiguration{
				RequireSsl:  true,
				Ipv4Enabled: true,
			},
			Tier:           vars.GetString("tier"),
			DataDiskSizeGb: int64(vars.GetInt("disk_size")),
			LocationPreference: &googlecloudsql.LocationPreference{
				Zone: vars.GetString("zone"),
			},
			DataDiskType: vars.GetString("disk_type"),
			MaintenanceWindow: &googlecloudsql.MaintenanceWindow{
				Day:             int64(vars.GetInt("maintenance_window_day")),
				Hour:            int64(vars.GetInt("maintenance_window_hour")),
				UpdateTrack:     "stable",
				ForceSendFields: []string{"Day", "Hour"},
			},
			PricingPlan:       secondGenPricingPlan,
			ActivationPolicy:  vars.GetString("activation_policy"),
			ReplicationType:   vars.GetString("replication_type"),
			StorageAutoResize: &autoResize,
		},
		DatabaseVersion: vars.GetString("version"),
		Region:          vars.GetString("region"),
		FailoverReplica: &googlecloudsql.DatabaseInstanceFailoverReplica{
			Name: vars.GetString("failover_replica_name"),
		},
	}
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
func (b *CloudSQLBroker) Bind(ctx context.Context, instance models.ServiceInstanceDetails, bindingID string, details brokerapi.BindDetails) (map[string]interface{}, error) {
	// get context before trying to create anything to catch errors early
	params := make(map[string]interface{})
	if err := json.Unmarshal(details.RawParameters, &params); err != nil {
		return nil, fmt.Errorf("Error unmarshalling parameters: %s", err)
	}

	if err := b.ensureUsernamePassword(instance.ID, bindingID, &details); err != nil {
		return nil, err
	}

	combinedCreds := varcontext.Builder()

	// Create the service account
	saCreds, err := b.BrokerBase.Bind(ctx, instance, bindingID, details)
	if err != nil {
		return saCreds, err
	}
	combinedCreds.MergeMap(saCreds)

	sqlCreds, err := b.createSqlCredentials(ctx, bindingID, details, instance)
	if err != nil {
		return saCreds, err
	}
	combinedCreds.MergeMap(sqlCreds)

	uriPrefix := ""
	if cast.ToBool(params["jdbc_uri_format"]) {
		uriPrefix = "jdbc:"
	}
	combinedCreds.MergeMap(map[string]interface{}{"UriPrefix": uriPrefix})

	return combinedCreds.BuildMap()
}

func (b *CloudSQLBroker) BuildInstanceCredentials(ctx context.Context, bindRecord models.ServiceBindingCredentials, instanceRecord models.ServiceInstanceDetails) (map[string]interface{}, error) {
	service, err := broker.GetServiceById(instanceRecord.ServiceId)
	if err != nil {
		return nil, err
	}
	uriFormat := ""
	switch service.Name {
	case models.CloudsqlMySQLName:
		uriFormat = `${str.queryEscape(UriPrefix)}mysql://${str.queryEscape(Username)}:${str.queryEscape(Password)}@${str.queryEscape(host)}/${str.queryEscape(database_name)}?ssl_mode=required`
	case models.CloudsqlPostgresName:
		uriFormat = `${str.queryEscape(UriPrefix)}postgres://${str.queryEscape(Username)}:${str.queryEscape(Password)}@${str.queryEscape(host)}/${str.queryEscape(database_name)}?sslmode=require&sslcert=${str.queryEscape(ClientCert)}&sslkey=${str.queryEscape(ClientKey)}&sslrootcert=${str.queryEscape(CaCert)}`
	default:
		return map[string]interface{}{}, errors.New("Unknown service")
	}

	combinedCreds, err := b.BrokerBase.BuildInstanceCredentials(ctx, bindRecord, instanceRecord)
	if err != nil {
		return nil, err
	}

	return varcontext.Builder().
		MergeMap(combinedCreds).
		MergeEvalResult("uri", uriFormat).
		BuildMap()
}

// Unbind deletes the database user, service account and invalidates the ssl certs associated with this binding.
func (b *CloudSQLBroker) Unbind(ctx context.Context, instance models.ServiceInstanceDetails, binding models.ServiceBindingCredentials) error {
	var accumulator error

	if err := b.deleteSqlSslCert(ctx, binding, instance); err != nil {
		accumulator = multierror.Append(accumulator, err)
	}

	if err := b.deleteSqlUserAccount(ctx, binding, instance); err != nil {
		accumulator = multierror.Append(accumulator, err)
	}

	if err := b.BrokerBase.Unbind(ctx, instance, binding); err != nil {
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
		// Update the instance information from the server side before
		// creating the database. The modification happens _only_ to
		// this instance of the details and is not persisted to the db.
		if err := b.UpdateInstanceDetails(ctx, &instance); err != nil {
			return true, err
		}

		return true, b.createDatabase(ctx, &instance)
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
func (b *CloudSQLBroker) Deprovision(ctx context.Context, instance models.ServiceInstanceDetails, details brokerapi.DeprovisionDetails) (*string, error) {
	sqlService, err := b.createClient(ctx)
	if err != nil {
		return nil, err
	}

	// delete the instance from google
	op, err := sqlService.Instances.Delete(b.ProjectId, instance.Name).Do()
	if err != nil {
		return nil, fmt.Errorf("Error deleting instance: %s", err)
	}

	return &op.Name, nil
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
