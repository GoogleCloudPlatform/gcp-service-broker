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

package broker_base

import (
	"code.cloudfoundry.org/lager"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/models"
	"golang.org/x/oauth2/jwt"
)

type BrokerBase struct {
	AccountManager models.AccountManager
	HttpConfig     *jwt.Config
	ProjectId      string
	Logger         lager.Logger
}

func (b *BrokerBase) Bind(instanceID, bindingID string, details models.BindDetails) (models.ServiceBindingCredentials, error) {

	// Create account
	newBinding, err := b.AccountManager.CreateCredentials(instanceID, bindingID, details, models.ServiceInstanceDetails{})

	if err != nil {
		return models.ServiceBindingCredentials{}, err
	}

	return newBinding, nil
}

func (b *BrokerBase) BuildInstanceCredentials(bindDetails models.ServiceBindingCredentials, instanceDetails models.ServiceInstanceDetails) (map[string]string, error) {
	return b.AccountManager.BuildInstanceCredentials(bindDetails, instanceDetails)
}

func (b *BrokerBase) Unbind(creds models.ServiceBindingCredentials) error {

	err := b.AccountManager.DeleteCredentials(creds)
	if err != nil {
		return err
	}

	return nil
}

// Does nothing but return an error because Base services are provisioned synchronously so this method should not be called
func (b *BrokerBase) PollInstance(instanceID string) (bool, error) {
	return true, models.ErrServiceIsNotAsync
}

// Indicates provisioning is done synchronously
func (b *BrokerBase) ProvisionsAsync() bool {
	return false
}

func (b *BrokerBase) DeprovisionsAsync() bool {
	return false
}

// used during polling of async operations to determine if the workflow is a provision or deprovision flow based off the
// type of the most recent operation
func (b *BrokerBase) LastOperationWasDelete(instanceId string) (bool, error) {
	panic("Can't check last operation on a synchronous service")
}
