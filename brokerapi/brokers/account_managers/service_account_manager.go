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

package account_managers

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/models"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/broker"
	"github.com/GoogleCloudPlatform/gcp-service-broker/utils"
	"github.com/pivotal-cf/brokerapi"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/jwt"
	cloudres "google.golang.org/api/cloudresourcemanager/v1"
	iam "google.golang.org/api/iam/v1"
)

const roleResourcePrefix = "roles/"
const saResourcePrefix = "serviceAccount:"
const saPrefix = "pcf-binding-"
const projectResourcePrefix = "projects/"

type ServiceAccountManager struct {
	ProjectId  string
	HttpConfig *jwt.Config
}

// If roleWhitelist is specified, then the extracted role is validated against it and an error is returned if
// the role is not contained within the whitelist
func (sam *ServiceAccountManager) CreateCredentials(bindingID string, details brokerapi.BindDetails, roleWhitelist []string) (models.ServiceBindingCredentials, error) {
	role, err := sam.ExtractRole(details)
	if err != nil {
		return models.ServiceBindingCredentials{}, err
	}

	if !whitelistAllows(roleWhitelist, role) {
		return models.ServiceBindingCredentials{}, fmt.Errorf("The role %s is not allowed for this service. You must use one of %v.", role, roleWhitelist)
	}

	return sam.CreateAccountWithRoles(bindingID, []string{role})
}

func (sam *ServiceAccountManager) ExtractRole(details brokerapi.BindDetails) (string, error) {
	bindParameters := map[string]interface{}{}
	if err := json.Unmarshal(details.RawParameters, &bindParameters); err != nil {
		return "", err
	}

	role, ok := bindParameters["role"].(string)
	if !ok {
		return "", errors.New("Error getting role as string from request")
	}

	return role, nil
}

// CreateAccountWithRoles creates a service account with a name based on bindingID, JSON key and grants it zero or more roles
// the roles MUST be missing the roles/ prefix.
func (sam *ServiceAccountManager) CreateAccountWithRoles(bindingID string, roles []string) (models.ServiceBindingCredentials, error) {
	// create and save account
	newSA, err := sam.createServiceAccount(bindingID)
	if err != nil {
		return models.ServiceBindingCredentials{}, err
	}

	// adjust account permissions
	// roles defined here: https://cloud.google.com/iam/docs/understanding-roles?hl=en_US#curated_roles
	for _, role := range roles {
		if err := sam.grantRoleToAccount(role, newSA); err != nil {
			return models.ServiceBindingCredentials{}, err
		}
	}

	// create and save key
	newSAKey, err := sam.createServiceAccountKey(newSA)
	if err != nil {
		return models.ServiceBindingCredentials{}, fmt.Errorf("Error creating new service account key: %s", err)
	}

	newSAInfo := ServiceAccountInfo{
		Name:           newSA.DisplayName,
		Email:          newSA.Email,
		UniqueId:       newSA.UniqueId,
		PrivateKeyData: newSAKey.PrivateKeyData,
		ProjectId:      sam.ProjectId,
	}

	saBytes, err := json.Marshal(&newSAInfo)
	if err != nil {
		return models.ServiceBindingCredentials{}, fmt.Errorf("Error marshalling new service account key %s", err)
	}

	newCreds := models.ServiceBindingCredentials{
		OtherDetails: string(saBytes),
	}

	return newCreds, nil
}

// deletes the given service account from Google
func (sam *ServiceAccountManager) DeleteCredentials(binding models.ServiceBindingCredentials) error {

	var saCreds ServiceAccountInfo
	if err := json.Unmarshal([]byte(binding.OtherDetails), &saCreds); err != nil {
		return fmt.Errorf("Error unmarshalling credentials: %s", err)
	}

	iamService, err := iam.New(sam.HttpConfig.Client(context.Background()))
	if err != nil {
		return fmt.Errorf("Error creating IAM service: %s", err)
	}
	saService := iam.NewProjectsServiceAccountsService(iamService)

	var resourceName = projectResourcePrefix + sam.ProjectId + "/serviceAccounts/" + saCreds.UniqueId

	_, err = saService.Delete(resourceName).Do()
	if err != nil {
		return fmt.Errorf("error deleting service account: %s", err)
	}
	return nil
}

func (b *ServiceAccountManager) BuildInstanceCredentials(bindRecord models.ServiceBindingCredentials, instanceRecord models.ServiceInstanceDetails) (map[string]string, error) {
	bindDetails := bindRecord.GetOtherDetails()
	instanceDetails := instanceRecord.GetOtherDetails()
	return utils.MergeStringMaps(bindDetails, instanceDetails), nil
}

// XXX names are truncated to 20 characters because of a bug in the IAM service
func ServiceAccountName(bindingId string) string {
	name := saPrefix + bindingId
	if len(name) > 20 {
		return name[:20]
	} else {
		return name
	}
}

func (sam *ServiceAccountManager) createServiceAccount(bindingID string) (*iam.ServiceAccount, error) {
	client := sam.HttpConfig.Client(context.Background())
	iamService, err := iam.New(client)
	if err != nil {
		return nil, fmt.Errorf("Error creating new IAM service: %s", err)
	}

	someName := ServiceAccountName(bindingID)
	resourceName := projectResourcePrefix + sam.ProjectId

	// create and save account
	newSARequest := iam.CreateServiceAccountRequest{
		AccountId: someName,
		ServiceAccount: &iam.ServiceAccount{
			DisplayName: someName,
		},
	}

	return iam.NewProjectsServiceAccountsService(iamService).Create(resourceName, &newSARequest).Do()
}

func (sam *ServiceAccountManager) createServiceAccountKey(account *iam.ServiceAccount) (*iam.ServiceAccountKey, error) {
	client := sam.HttpConfig.Client(context.Background())
	iamService, err := iam.New(client)
	if err != nil {
		return nil, fmt.Errorf("Error creating new IAM service: %s", err)
	}

	saKeyService := iam.NewProjectsServiceAccountsKeysService(iamService)
	return saKeyService.Create(account.Name, &iam.CreateServiceAccountKeyRequest{}).Do()
}

func (sam *ServiceAccountManager) grantRoleToAccount(role string, account *iam.ServiceAccount) error {
	client := sam.HttpConfig.Client(context.Background())

	cloudresService, err := cloudres.New(client)
	if err != nil {
		return fmt.Errorf("Error creating new cloud resource management service: %s", err)
	}

	currPolicy, err := cloudresService.Projects.GetIamPolicy(sam.ProjectId, &cloudres.GetIamPolicyRequest{}).Do()
	if err != nil {
		return fmt.Errorf("Error getting current project iam policy: %s", err)
	}

	// seems not really necessary, but collapse the bindings into single role entries just in case.
	var existingBinding *cloudres.Binding

	for _, binding := range currPolicy.Bindings {
		if binding.Role == roleResourcePrefix+role {
			existingBinding = binding
		}
	}

	if existingBinding != nil {
		existingBinding.Members = append(existingBinding.Members, saResourcePrefix+account.Email)
	} else {
		existingBinding = &cloudres.Binding{
			Members: []string{saResourcePrefix + account.Email},
			Role:    roleResourcePrefix + role,
		}
		b := append(currPolicy.Bindings, existingBinding)
		currPolicy.Bindings = b
	}

	newPolicyRequest := cloudres.SetIamPolicyRequest{
		Policy: currPolicy,
	}
	_, err = cloudresService.Projects.SetIamPolicy(sam.ProjectId, &newPolicyRequest).Do()
	if err != nil {
		return fmt.Errorf("Error assigning policy to service account: %s", err)
	}

	return nil
}

type ServiceAccountInfo struct {
	// the bits to save
	Name      string
	Email     string
	UniqueId  string
	ProjectId string

	// the bit to return
	PrivateKeyData string
}

func ServiceAccountBindInputVariables(roleWhitelist []string) []broker.BrokerVariable {
	whitelistMap := make(map[interface{}]string)
	for _, role := range roleWhitelist {
		whitelistMap[role] = role
	}

	return []broker.BrokerVariable{
		broker.BrokerVariable{
			Required:  true,
			FieldName: "role",
			Type:      broker.JsonTypeString,
			Details:   `The role for the account without the "roles/" prefix. See https://cloud.google.com/iam/docs/understanding-roles for available roles.`,
			Enum:      whitelistMap,
		},
	}
}

// Variables output by all brokers that return service account info
func ServiceAccountBindOutputVariables() []broker.BrokerVariable {
	return []broker.BrokerVariable{
		broker.BrokerVariable{
			FieldName: "Email",
			Type:      broker.JsonTypeString,
			Details:   "Email address of the service account",
		},
		broker.BrokerVariable{
			FieldName: "Name",
			Type:      broker.JsonTypeString,
			Details:   "The name of the service account",
		},
		broker.BrokerVariable{
			FieldName: "PrivateKeyData",
			Type:      broker.JsonTypeString,
			Details:   "Service account private key data. Base-64 encoded JSON.",
		},
		broker.BrokerVariable{
			FieldName: "ProjectId",
			Type:      broker.JsonTypeString,
			Details:   "ID of the project that owns the service account",
		},
		broker.BrokerVariable{
			FieldName: "UniqueId",
			Type:      broker.JsonTypeString,
			Details:   "Unique and stable id of the service account",
		},
	}
}

func whitelistAllows(whitelist []string, needle string) bool {
	for _, s := range whitelist {
		if s == needle {
			return true
		}
	}

	return len(whitelist) > 0
}
