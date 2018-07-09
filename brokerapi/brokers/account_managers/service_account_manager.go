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
	"encoding/json"
	"errors"
	"fmt"
	"gcp-service-broker/brokerapi/brokers/models"
	"gcp-service-broker/utils"

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

// creates a new service account for the given binding id with the role listed in details.Parameters["role"]
func (sam *ServiceAccountManager) CreateCredentials(instanceID string, bindingID string, details models.BindDetails, instance models.ServiceInstanceDetails) (models.ServiceBindingCredentials, error) {
	client := sam.HttpConfig.Client(context.Background())

	role, ok := details.Parameters["role"].(string)
	if !ok {
		return models.ServiceBindingCredentials{}, errors.New("Error getting role as string from request")
	}

	someName := ServiceAccountName(bindingID)
	var resourceName = projectResourcePrefix + sam.ProjectId
	var err error

	iamService, err := iam.New(client)
	if err != nil {
		return models.ServiceBindingCredentials{}, fmt.Errorf("Error creating new iam service: %s", err)
	}
	saService := iam.NewProjectsServiceAccountsService(iamService)

	// create and save account
	newSARequest := iam.CreateServiceAccountRequest{
		AccountId: someName,
		ServiceAccount: &iam.ServiceAccount{
			DisplayName: someName,
		},
	}

	newSA, err := saService.Create(resourceName, &newSARequest).Do()
	if err != nil {
		return models.ServiceBindingCredentials{}, fmt.Errorf("Error creating service account: %s", err)
	}

	// adjust account permissions
	// roles defined here: https://cloud.google.com/iam/docs/understanding-roles?hl=en_US#curated_roles
	cloudresService, err := cloudres.New(client)
	if err != nil {
		return models.ServiceBindingCredentials{}, fmt.Errorf("Error creating new cloud resource management service: %s", err)
	}

	currPolicy, err := cloudresService.Projects.GetIamPolicy(sam.ProjectId, &cloudres.GetIamPolicyRequest{}).Do()
	if err != nil {
		return models.ServiceBindingCredentials{}, fmt.Errorf("Error getting current project iam policy: %s", err)
	}

	// seems not really necessary, but collapse the bindings into single role entries just in case.
	var existingBinding *cloudres.Binding

	for _, binding := range currPolicy.Bindings {
		if binding.Role == roleResourcePrefix+role {
			existingBinding = binding
		}
	}

	if existingBinding != nil {
		existingBinding.Members = append(existingBinding.Members, saResourcePrefix+newSA.Email)
	} else {
		existingBinding = &cloudres.Binding{
			Members: []string{saResourcePrefix + newSA.Email},
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
		return models.ServiceBindingCredentials{}, fmt.Errorf("ERROR assigning policy to service account: %s", err)
	}

	// create and save key
	saKeyService := iam.NewProjectsServiceAccountsKeysService(iamService)
	newSAKey, err := saKeyService.Create(newSA.Name, &iam.CreateServiceAccountKeyRequest{}).Do()
	if err != nil {
		return models.ServiceBindingCredentials{}, fmt.Errorf("ERROR creating new service account key: %s", err)
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

type ServiceAccountInfo struct {
	// the bits to save
	Name      string
	Email     string
	UniqueId  string
	ProjectId string

	// the bit to return
	PrivateKeyData string
}
