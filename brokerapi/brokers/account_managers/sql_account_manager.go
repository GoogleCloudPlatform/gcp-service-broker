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

package account_managers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/models"
	"github.com/GoogleCloudPlatform/gcp-service-broker/db_service"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/broker"
	"github.com/GoogleCloudPlatform/gcp-service-broker/utils"
	"github.com/pivotal-cf/brokerapi"
	"golang.org/x/oauth2/jwt"
	googlecloudsql "google.golang.org/api/sqladmin/v1beta4"
)

type SqlAccountManager struct {
	HttpConfig *jwt.Config
	ProjectId  string
}

// inserts a new user into the database and creates new ssl certs
func (sam *SqlAccountManager) CreateCredentials(instanceID string, bindingID string, details brokerapi.BindDetails, instance models.ServiceInstanceDetails) (models.ServiceBindingCredentials, error) {

	bindParameters := map[string]interface{}{}
	if err := json.Unmarshal(details.RawParameters, &bindParameters); err != nil {
		return models.ServiceBindingCredentials{}, err
	}

	var err error
	username, usernameOk := bindParameters["username"].(string)
	password, passwordOk := bindParameters["password"].(string)

	if !passwordOk || !usernameOk {
		return models.ServiceBindingCredentials{}, errors.New("Error binding, missing parameters. Required parameters are username and password")
	}

	// create username, pw with grants
	sqlService, err := googlecloudsql.New(sam.HttpConfig.Client(context.Background()))
	if err != nil {
		return models.ServiceBindingCredentials{}, fmt.Errorf("Error creating CloudSQL client: %s", err)
	}

	op, err := sqlService.Users.Insert(sam.ProjectId, instance.Name, &googlecloudsql.User{
		Name:     username,
		Password: password,
	}).Do()

	if err != nil {
		return models.ServiceBindingCredentials{}, fmt.Errorf("Error inserting new database user: %s", err)
	}

	// poll for the user creation operation to be completed
	err = sam.pollOperationUntilDone(op, sam.ProjectId)
	if err != nil {
		return models.ServiceBindingCredentials{}, fmt.Errorf("Error encountered while polling until operation id %s completes: %s", op.Name, err)
	}

	var creds SqlAccountInfo
	// create ssl certs
	certname := bindingID + "cert"
	if len(bindingID) >= 10 {
		certname = bindingID[:10] + "cert"
	}
	newCert, err := sqlService.SslCerts.Insert(sam.ProjectId, instance.Name, &googlecloudsql.SslCertsInsertRequest{
		CommonName: certname,
	}).Do()
	if err != nil {
		return models.ServiceBindingCredentials{}, fmt.Errorf("Error creating ssl certs: %s", err)
	}

	creds = SqlAccountInfo{
		Username:        username,
		Password:        password,
		Sha1Fingerprint: newCert.ClientCert.CertInfo.Sha1Fingerprint,
		CaCert:          newCert.ServerCaCert.Cert,
		ClientCert:      newCert.ClientCert.CertInfo.Cert,
		ClientKey:       newCert.ClientCert.CertPrivateKey,
	}

	credBytes, err := json.Marshal(&creds)
	if err != nil {
		return models.ServiceBindingCredentials{}, fmt.Errorf("Error marshalling credentials: %s", err)
	}

	newBinding := models.ServiceBindingCredentials{
		OtherDetails: string(credBytes),
	}

	return newBinding, nil
}

// deletes the user from the database and invalidates the associated ssl certs
func (sam *SqlAccountManager) DeleteCredentials(binding models.ServiceBindingCredentials) error {
	var err error

	var sqlCreds SqlAccountInfo
	if err := json.Unmarshal([]byte(binding.OtherDetails), &sqlCreds); err != nil {
		return fmt.Errorf("Error unmarshalling credentials: %s", err)
	}

	var instance models.ServiceInstanceDetails
	if err = db_service.DbConnection.Where("id = ?", binding.ServiceInstanceId).Find(&instance).Error; err != nil {
		return fmt.Errorf("Database error retrieving instance details: %s", err)
	}

	sqlService, err := googlecloudsql.New(sam.HttpConfig.Client(context.Background()))
	if err != nil {
		return fmt.Errorf("Error creating CloudSQL client: %s", err)
	}

	// If we didn't generate ssl certs for this binding, then we cannot delete them
	if sqlCreds.CaCert != "" {
		certOp, err := sqlService.SslCerts.Delete(sam.ProjectId, instance.Name, sqlCreds.Sha1Fingerprint).Do()
		if err != nil {
			return fmt.Errorf("Error deleting ssl cert: %s", err)
		}

		err = sam.pollOperationUntilDone(certOp, sam.ProjectId)
		if err != nil {
			return fmt.Errorf("Error encountered while polling until operation id %s completes: %s", certOp.Name, err)
		}
	}

	// delete our user
	userOp, err := sqlService.Users.Delete(sam.ProjectId, instance.Name, "", sqlCreds.Username).Do()
	if err != nil {
		return fmt.Errorf("Error deleting user: %s", err)
	}

	err = sam.pollOperationUntilDone(userOp, sam.ProjectId)
	if err != nil {
		return fmt.Errorf("Error encountered while polling until operation id %s completes: %s", userOp.Name, err)
	}

	return nil
}

// polls the cloud sql operations service once per second until the given operation is done
// TODO(cbriant): ensure this stays under api call quota
func (sam *SqlAccountManager) pollOperationUntilDone(op *googlecloudsql.Operation, projectId string) error {
	sqlService, err := googlecloudsql.New(sam.HttpConfig.Client(context.Background()))
	if err != nil {
		return fmt.Errorf("Error creating new cloudsql client: %s", err)
	}

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
	return nil
}

func (b *SqlAccountManager) BuildInstanceCredentials(bindRecord models.ServiceBindingCredentials, instanceRecord models.ServiceInstanceDetails) (map[string]string, error) {
	instanceDetails := instanceRecord.GetOtherDetails()
	bindDetails := bindRecord.GetOtherDetails()

	service_to_name := broker.MapServiceIdToName()

	sid := instanceRecord.ServiceId

	combinedCreds := utils.MergeStringMaps(bindDetails, instanceDetails)

	if service_to_name[sid] == models.CloudsqlMySQLName {
		combinedCreds["uri"] = fmt.Sprintf("%smysql://%s:%s@%s/%s?ssl_mode=required",
			combinedCreds["UriPrefix"], url.QueryEscape(combinedCreds["Username"]), url.QueryEscape(combinedCreds["Password"]), combinedCreds["host"], combinedCreds["database_name"])
	} else if service_to_name[sid] == models.CloudsqlPostgresName {
		combinedCreds["uri"] = fmt.Sprintf("%spostgres://%s:%s@%s/%s?sslmode=require&sslcert=%s&sslkey=%s&sslrootcert=%s",
			combinedCreds["UriPrefix"], url.QueryEscape(combinedCreds["Username"]), url.QueryEscape(combinedCreds["Password"]), combinedCreds["host"], combinedCreds["database_name"], url.QueryEscape(combinedCreds["ClientCert"]), url.QueryEscape(combinedCreds["ClientKey"]), url.QueryEscape(combinedCreds["CaCert"]))
	} else {
		return map[string]string{}, errors.New("Unknown service")
	}

	return combinedCreds, nil
}

type SqlAccountInfo struct {
	// the bits to return
	Username   string
	Password   string
	CaCert     string
	ClientCert string
	ClientKey  string

	// the bits to save
	Sha1Fingerprint string
}
