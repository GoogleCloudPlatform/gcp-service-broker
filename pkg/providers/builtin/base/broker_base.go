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

package base

import (
	"context"

	"code.cloudfoundry.org/lager"
	"github.com/GoogleCloudPlatform/gcp-service-broker/db_service/models"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/providers/builtin/account_managers"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/varcontext"

	"golang.org/x/oauth2/jwt"
)

//go:generate counterfeiter . ServiceAccountManager

type ServiceAccountManager interface {
	CreateCredentials(ctx context.Context, vc *varcontext.VarContext) (map[string]interface{}, error)
	DeleteCredentials(ctx context.Context, creds models.ServiceBindingCredentials) error
}

// NewBrokerBase creates a new broker base and account manager it uses from the
// given settings.
func NewBrokerBase(projectId string, auth *jwt.Config, logger lager.Logger) BrokerBase {
	saManager := &account_managers.ServiceAccountManager{
		HttpConfig: auth,
		ProjectId:  projectId,
		Logger:     logger,
	}

	return BrokerBase{
		AccountManager: saManager,
		HttpConfig:     auth,
		ProjectId:      projectId,
		Logger:         logger,
	}
}

// BrokerBase is the reference bind and unbind implementation for brokers that
// bind and unbind with only Service Accounts.
type BrokerBase struct {
	synchronousBase
	MergedInstanceCredsMixin

	AccountManager ServiceAccountManager
	HttpConfig     *jwt.Config
	ProjectId      string
	Logger         lager.Logger
}

// Bind creates a service account with access to the provisioned resource with
// the given instance.
func (b *BrokerBase) Bind(ctx context.Context, vc *varcontext.VarContext) (map[string]interface{}, error) {
	return b.AccountManager.CreateCredentials(ctx, vc)
}

// Unbind deletes the created service account from the GCP Project.
func (b *BrokerBase) Unbind(ctx context.Context, instance models.ServiceInstanceDetails, creds models.ServiceBindingCredentials) error {
	return b.AccountManager.DeleteCredentials(ctx, creds)
}

// UpdateInstanceDetails updates the ServiceInstanceDetails with the most recent state from GCP.
// This instance is a no-op method.
func (b *BrokerBase) UpdateInstanceDetails(ctx context.Context, instance *models.ServiceInstanceDetails) error {
	return nil
}
