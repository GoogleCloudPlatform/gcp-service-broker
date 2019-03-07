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
	"fmt"
	"net/http"
	"time"

	"code.cloudfoundry.org/lager"
	"github.com/GoogleCloudPlatform/gcp-service-broker/db_service/models"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/broker"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/validation"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/varcontext"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/jwt"
	cloudres "google.golang.org/api/cloudresourcemanager/v1"
	"google.golang.org/api/googleapi"
	iam "google.golang.org/api/iam/v1"
)

const (
	roleResourcePrefix    = "roles/"
	saResourcePrefix      = "serviceAccount:"
	projectResourcePrefix = "projects/"
)

type ServiceAccountManager struct {
	ProjectId  string
	HttpConfig *jwt.Config
	Logger     lager.Logger
}

// If roleWhitelist is specified, then the extracted role is validated against it and an error is returned if
// the role is not contained within the whitelist
func (sam *ServiceAccountManager) CreateCredentials(ctx context.Context, vc *varcontext.VarContext) (map[string]interface{}, error) {
	role := vc.GetString("role")
	accountId := vc.GetString("service_account_name")
	displayName := vc.GetString("service_account_display_name")

	if err := vc.Error(); err != nil {
		return nil, err
	}

	sam.Logger.Info("create-service-account", lager.Data{
		"role":                         role,
		"service_account_name":         accountId,
		"service_account_display_name": displayName,
	})

	// create and save account
	newSA, err := sam.createServiceAccount(ctx, accountId, displayName)
	if err != nil {
		return nil, err
	}

	// adjust account permissions
	// roles defined here: https://cloud.google.com/iam/docs/understanding-roles?hl=en_US#curated_roles
	if err := sam.grantRoleToAccount(ctx, role, newSA); err != nil {
		return nil, err
	}

	// create and save key
	newSAKey, err := sam.createServiceAccountKey(ctx, newSA)
	if err != nil {
		return nil, fmt.Errorf("Error creating new service account key: %s", err)
	}

	newSAInfo := ServiceAccountInfo{
		Name:           newSA.DisplayName,
		Email:          newSA.Email,
		UniqueId:       newSA.UniqueId,
		PrivateKeyData: newSAKey.PrivateKeyData,
		ProjectId:      sam.ProjectId,
	}

	return varcontext.Builder().MergeStruct(newSAInfo).BuildMap()
}

// deletes the given service account from Google
func (sam *ServiceAccountManager) DeleteCredentials(ctx context.Context, binding models.ServiceBindingCredentials) error {

	var saCreds ServiceAccountInfo
	if err := json.Unmarshal([]byte(binding.OtherDetails), &saCreds); err != nil {
		return fmt.Errorf("Error unmarshalling credentials: %s", err)
	}

	iamService, err := iam.New(sam.HttpConfig.Client(ctx))
	if err != nil {
		return fmt.Errorf("Error creating IAM service: %s", err)
	}

	resourceName := projectResourcePrefix + sam.ProjectId + "/serviceAccounts/" + saCreds.UniqueId
	if _, err := iam.NewProjectsServiceAccountsService(iamService).Delete(resourceName).Do(); err != nil {
		if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == http.StatusNotFound {
			return nil
		}

		return fmt.Errorf("error deleting service account: %s", err)
	}

	return nil
}

func (sam *ServiceAccountManager) createServiceAccount(ctx context.Context, accountId, displayName string) (*iam.ServiceAccount, error) {
	client := sam.HttpConfig.Client(ctx)
	iamService, err := iam.New(client)
	if err != nil {
		return nil, fmt.Errorf("Error creating new IAM service: %s", err)
	}

	resourceName := projectResourcePrefix + sam.ProjectId

	// create and save account
	newSARequest := iam.CreateServiceAccountRequest{
		AccountId: accountId,
		ServiceAccount: &iam.ServiceAccount{
			DisplayName: displayName,
		},
	}

	return iam.NewProjectsServiceAccountsService(iamService).Create(resourceName, &newSARequest).Do()
}

func (sam *ServiceAccountManager) createServiceAccountKey(ctx context.Context, account *iam.ServiceAccount) (*iam.ServiceAccountKey, error) {
	client := sam.HttpConfig.Client(ctx)
	iamService, err := iam.New(client)
	if err != nil {
		return nil, fmt.Errorf("Error creating new IAM service: %s", err)
	}

	saKeyService := iam.NewProjectsServiceAccountsKeysService(iamService)
	return saKeyService.Create(account.Name, &iam.CreateServiceAccountKeyRequest{}).Do()
}

func (sam *ServiceAccountManager) grantRoleToAccount(ctx context.Context, role string, account *iam.ServiceAccount) error {
	client := sam.HttpConfig.Client(ctx)

	cloudresService, err := cloudres.New(client)
	if err != nil {
		return fmt.Errorf("Error creating new cloud resource management service: %s", err)
	}

	for attempt := 0; attempt < 3; attempt++ {
		currPolicy, err := cloudresService.Projects.GetIamPolicy(sam.ProjectId, &cloudres.GetIamPolicyRequest{}).Do()
		if err != nil {
			return fmt.Errorf("Error getting current project iam policy: %s", err)
		}

		currPolicy.Bindings = mergeBindings(append(currPolicy.Bindings, &cloudres.Binding{
			Members: []string{saResourcePrefix + account.Email},
			Role:    roleResourcePrefix + role,
		}))

		newPolicyRequest := cloudres.SetIamPolicyRequest{
			Policy: currPolicy,
		}
		_, err = cloudresService.Projects.SetIamPolicy(sam.ProjectId, &newPolicyRequest).Do()
		if err == nil {
			return nil
		}

		if isConflictError(err) {
			time.Sleep(5 * time.Second)
			continue
		} else {
			return fmt.Errorf("Error assigning policy to service account: %s", err)
		}
	}

	return err
}

func isConflictError(err error) bool {
	gerr, ok := err.(*googleapi.Error)
	return ok && gerr != nil && gerr.Code == 409
}

type ServiceAccountInfo struct {
	// the bits to save
	Name      string `json:"Name"`
	Email     string `json:"Email"`
	UniqueId  string `json:"UniqueId"`
	ProjectId string `json:"ProjectId"`

	// the bit to return
	PrivateKeyData string `json:"PrivateKeyData"`
}

// ServiceAccountWhitelistWithDefault holds non-overridable whitelists with default values.
func ServiceAccountWhitelistWithDefault(whitelist []string, defaultValue string) []broker.BrokerVariable {
	whitelistEnum := make(map[interface{}]string)
	for _, val := range whitelist {
		whitelistEnum[val] = roleResourcePrefix + val
	}

	return []broker.BrokerVariable{
		{
			FieldName: "role",
			Type:      broker.JsonTypeString,
			Details:   `The role for the account without the "roles/" prefix. See: https://cloud.google.com/iam/docs/understanding-roles for more details.`,
			Enum:      whitelistEnum,
			Default:   defaultValue,
		},
	}
}

// ServiceAccountBindComputedVariables holds computed variables required to provision service accounts, label them and ensure they are unique.
func ServiceAccountBindComputedVariables() []varcontext.DefaultVariable {
	return []varcontext.DefaultVariable{
		// XXX names are truncated to 20 characters because of a bug in the IAM service
		{Name: "service_account_name", Default: `${str.truncate(20, "gsb-binding-${request.binding_id}")}`, Overwrite: true},
		{Name: "service_account_display_name", Default: "${service_account_name}", Overwrite: true},
	}
}

// FixedRoleBindComputedVariables allows you to create a service account with a
// fixed role.
func FixedRoleBindComputedVariables(role string) []varcontext.DefaultVariable {
	fixedRoleVar := varcontext.DefaultVariable{Name: "role", Default: role, Overwrite: true}
	return append(ServiceAccountBindComputedVariables(), fixedRoleVar)
}

// Variables output by all brokers that return service account info
func ServiceAccountBindOutputVariables() []broker.BrokerVariable {
	return []broker.BrokerVariable{
		{
			FieldName: "Email",
			Type:      broker.JsonTypeString,
			Details:   "Email address of the service account.",
			Required:  true,
			Constraints: validation.NewConstraintBuilder().
				Examples("gsb-binding-ex312029@my-project.iam.gserviceaccount.com").
				Pattern(`^gsb-binding-[a-z0-9-]+@.+\.gserviceaccount\.com$`).
				Build(),
		},
		{
			FieldName: "Name",
			Type:      broker.JsonTypeString,
			Details:   "The name of the service account.",
			Required:  true,
			Constraints: validation.NewConstraintBuilder().
				Examples("gsb-binding-ex312029").
				Build(),
		},
		{
			FieldName: "PrivateKeyData",
			Type:      broker.JsonTypeString,
			Details:   "Service account private key data. Base64 encoded JSON.",
			Required:  true,
			Constraints: validation.NewConstraintBuilder().
				MinLength(512).                // absolute lower bound
				Pattern(`^[A-Za-z0-9+/]*=*$`). // very rough Base64 regex
				Build(),
		},
		{
			FieldName: "ProjectId",
			Type:      broker.JsonTypeString,
			Details:   "ID of the project that owns the service account.",
			Required:  true,
			Constraints: validation.NewConstraintBuilder().
				Examples("my-project").
				Pattern(`^[a-z0-9-]+$`).
				MinLength(6).
				MaxLength(30).
				Build(),
		},
		{
			FieldName: "UniqueId",
			Type:      broker.JsonTypeString,
			Details:   "Unique and stable ID of the service account.",
			Required:  true,
			Constraints: validation.NewConstraintBuilder().
				Examples("112447814736626230844").
				Build(),
		},
	}
}
