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

package broker_base

import (
	"context"
	"encoding/json"

	"code.cloudfoundry.org/lager"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/models"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/varcontext"

	"github.com/pivotal-cf/brokerapi"
	"golang.org/x/oauth2/jwt"
)

// BrokerBase is the reference bind and unbind implementation for brokers that
// bind and unbind with only Service Accounts.
type BrokerBase struct {
	AccountManager models.ServiceAccountManager
	HttpConfig     *jwt.Config
	ProjectId      string
	Logger         lager.Logger
}

// Bind creates a service account with access to the provisioned resource with
// the given instance.
func (b *BrokerBase) Bind(ctx context.Context, instanceID, bindingID string, details brokerapi.BindDetails) (map[string]interface{}, error) {
	return b.AccountManager.CreateCredentials(ctx, instanceID, bindingID, details, models.ServiceInstanceDetails{})
}

// BuildInstanceCredentials combines the bind credentials with the connection
// information in the instance details to get a full set of connection details.
func (b *BrokerBase) BuildInstanceCredentials(ctx context.Context, bindRecord models.ServiceBindingCredentials, instanceRecord models.ServiceInstanceDetails) (map[string]interface{}, error) {
	vc, err := varcontext.Builder().
		MergeJsonObject(json.RawMessage(bindRecord.OtherDetails)).
		MergeJsonObject(json.RawMessage(instanceRecord.OtherDetails)).
		Build()
	if err != nil {
		return nil, err
	}

	return vc.ToMap(), nil
}

// Unbind deletes the created service account from the GCP Project.
func (b *BrokerBase) Unbind(ctx context.Context, creds models.ServiceBindingCredentials) error {
	return b.AccountManager.DeleteCredentials(ctx, creds)
}

// PollInstance does nothing but return an error because Base services are
// provisioned synchronously so this method should not be called.
func (b *BrokerBase) PollInstance(ctx context.Context, instance models.ServiceInstanceDetails) (bool, error) {
	return true, brokerapi.ErrAsyncRequired
}

// ProvisionsAsync indicates if provisioning must be done asynchronously.
func (b *BrokerBase) ProvisionsAsync() bool {
	return false
}

// DeprovisionsAsync indicates if deprovisioning must be done asynchronously.
func (b *BrokerBase) DeprovisionsAsync() bool {
	return false
}

// UpdateInstanceDetails updates the ServiceInstanceDetails with the most recent state from GCP.
// This instance is a no-op method.
func (b *BrokerBase) UpdateInstanceDetails(ctx context.Context, instance *models.ServiceInstanceDetails) error {
	return nil
}
