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
	"context"
	"encoding/json"
	"fmt"

	"github.com/GoogleCloudPlatform/gcp-service-broker/db_service/models"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/varcontext"
	googlecloudsql "google.golang.org/api/sqladmin/v1beta4"
)

// inserts a new user into the database and creates new ssl certs
func (broker *CloudSQLBroker) createSqlCredentials(ctx context.Context, vars *varcontext.VarContext) (map[string]interface{}, error) {

	userAccount, err := broker.createSqlUserAccount(ctx, vars)
	if err != nil {
		return nil, err
	}

	sslCert, err := broker.createSqlSslCert(ctx, vars)
	if err != nil {
		return nil, err
	}

	return varcontext.Builder().MergeStruct(userAccount).MergeStruct(sslCert).BuildMap()
}

type sqlUserAccount struct {
	Username string `json:"Username"`
	Password string `json:"Password"`
}

type sqlSslCert struct {
	CaCert          string `json:"CaCert"`
	ClientCert      string `json:"ClientCert"`
	ClientKey       string `json:"ClientKey"`
	Sha1Fingerprint string `json:"Sha1Fingerprint"`
}

func (broker *CloudSQLBroker) createSqlUserAccount(ctx context.Context, vars *varcontext.VarContext) (*sqlUserAccount, error) {
	request := &googlecloudsql.User{
		Name:     vars.GetString("username"),
		Password: vars.GetString("password"),
	}
	instanceName := vars.GetString("db_name")

	if err := vars.Error(); err != nil {
		return nil, err
	}

	// create username, pw with grants
	client, err := broker.createClient(ctx)
	if err != nil {
		return nil, err
	}

	op, err := client.Users.Insert(broker.ProjectId, instanceName, request).Do()
	if err != nil {
		return nil, fmt.Errorf("Error creating new database user: %s", err)
	}

	// poll for the user creation operation to be completed
	if err := broker.pollOperationUntilDone(ctx, op, broker.ProjectId); err != nil {
		return nil, fmt.Errorf("Error encountered waiting for operation %q to finish: %s", op.Name, err)
	}

	return &sqlUserAccount{
		Username: request.Name,
		Password: request.Password,
	}, nil
}

func (broker *CloudSQLBroker) deleteSqlUserAccount(ctx context.Context, binding models.ServiceBindingCredentials, instance models.ServiceInstanceDetails) error {
	var creds sqlUserAccount
	if err := json.Unmarshal([]byte(binding.OtherDetails), &creds); err != nil {
		return fmt.Errorf("Error unmarshalling credentials: %s", err)
	}

	client, err := broker.createClient(ctx)
	if err != nil {
		return err
	}

	userList, err := client.Users.List(broker.ProjectId, instance.Name).Do()
	if err != nil {
		return fmt.Errorf("Error fetching users to delete: %s", err)
	}

	// XXX: CloudSQL used to allow deleting users without specifying the host,
	// however that no longer works. They also no longer accept a blank string
	// which _is_ a valid host, so we expand to a single space string if the
	// user we're trying to delete doesn't have some other host specified.
	hostToDelete := ""
	foundUser := false
	for _, user := range userList.Items {
		if user.Name == creds.Username {
			hostToDelete = user.Host
			foundUser = true
			break
		}
	}

	// XXX: If the user was already deleted, don't fail here because it could
	// block deprovisioning.
	if !foundUser {
		return nil
	}

	if hostToDelete == "" {
		hostToDelete = " "
	}

	op, err := client.Users.Delete(broker.ProjectId, instance.Name, hostToDelete, creds.Username).Do()
	if err != nil {
		return fmt.Errorf("Error deleting user: %s", err)
	}

	if err := broker.pollOperationUntilDone(ctx, op, broker.ProjectId); err != nil {
		return fmt.Errorf("Error encountered waiting for operation %q to finish: %s", op.Name, err)
	}

	return nil
}

func (broker *CloudSQLBroker) createSqlSslCert(ctx context.Context, vars *varcontext.VarContext) (*sqlSslCert, error) {
	request := &googlecloudsql.SslCertsInsertRequest{
		CommonName: vars.GetString("certname"),
	}
	instanceName := vars.GetString("db_name")

	if err := vars.Error(); err != nil {
		return nil, err
	}

	// create username, pw with grants
	client, err := broker.createClient(ctx)
	if err != nil {
		return nil, err
	}

	newCert, err := client.SslCerts.Insert(broker.ProjectId, instanceName, request).Do()
	if err != nil {
		return nil, fmt.Errorf("Error creating SSL certs: %s", err)
	}

	// poll for the user creation operation to be completed
	return &sqlSslCert{
		Sha1Fingerprint: newCert.ClientCert.CertInfo.Sha1Fingerprint,
		CaCert:          newCert.ServerCaCert.Cert,
		ClientCert:      newCert.ClientCert.CertInfo.Cert,
		ClientKey:       newCert.ClientCert.CertPrivateKey,
	}, nil
}

func (broker *CloudSQLBroker) deleteSqlSslCert(ctx context.Context, binding models.ServiceBindingCredentials, instance models.ServiceInstanceDetails) error {
	var creds sqlSslCert
	if err := json.Unmarshal([]byte(binding.OtherDetails), &creds); err != nil {
		return fmt.Errorf("Error unmarshalling credentials: %s", err)
	}

	client, err := broker.createClient(ctx)
	if err != nil {
		return err
	}

	// If we didn't generate SSL certs for this binding, then we cannot delete them
	if creds.CaCert == "" {
		return nil
	}

	op, err := client.SslCerts.Delete(broker.ProjectId, instance.Name, creds.Sha1Fingerprint).Do()
	if err != nil {
		return fmt.Errorf("Error deleting ssl cert: %s", err)
	}

	if err := broker.pollOperationUntilDone(ctx, op, broker.ProjectId); err != nil {
		return fmt.Errorf("Error encountered waiting for operation %q to finish: %s", op.Name, err)
	}

	return nil
}
