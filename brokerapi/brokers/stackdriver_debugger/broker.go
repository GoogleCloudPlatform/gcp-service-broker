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

package stackdriver_debugger

import (
	"gcp-service-broker/brokerapi/brokers/account_managers"
	"gcp-service-broker/brokerapi/brokers/broker_base"
	"gcp-service-broker/brokerapi/brokers/models"

	"code.cloudfoundry.org/lager"
	"golang.org/x/oauth2/jwt"
)

type StackdriverDebuggerBroker struct {
	HttpConfig            *jwt.Config
	ProjectId             string
	Logger                lager.Logger
	ServiceAccountManager *account_managers.ServiceAccountManager
	broker_base.BrokerBase
}

type InstanceInformation struct {
}

// No-op, no service is required for the Debugger
func (b *StackdriverDebuggerBroker) Provision(instanceId string, details models.ProvisionDetails, plan models.PlanDetails) (models.ServiceInstanceDetails, error) {
	return models.ServiceInstanceDetails{}, nil
}

// No-op, no service is required for the Debugger
func (b *StackdriverDebuggerBroker) Deprovision(instanceID string, details models.DeprovisionDetails) error {
	return nil
}

// Creates a service account with access to Stackdriver Debugger
func (b *StackdriverDebuggerBroker) Bind(instanceID, bindingID string, details models.BindDetails) (models.ServiceBindingCredentials, error) {
	if details.Parameters == nil {
		b.Logger.Info("the parameters are nil!")
		details.Parameters = make(map[string]interface{})
	}
	details.Parameters["role"] = "clouddebugger.agent"

	// Create account
	newBinding, err := b.ServiceAccountManager.CreateAccountInGoogle(instanceID, bindingID, details, models.ServiceInstanceDetails{})

	if err != nil {
		return models.ServiceBindingCredentials{}, err
	}

	return newBinding, nil
}
