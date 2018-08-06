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
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/broker_base"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/models"
	"github.com/GoogleCloudPlatform/gcp-service-broker/utils"
	"github.com/pivotal-cf/brokerapi"
)

type StackdriverDebuggerBroker struct {
	broker_base.BrokerBase
}

type InstanceInformation struct {
}

// No-op, no service is required for the Debugger
func (b *StackdriverDebuggerBroker) Provision(instanceId string, details brokerapi.ProvisionDetails, plan models.ServicePlan) (models.ServiceInstanceDetails, error) {
	return models.ServiceInstanceDetails{}, nil
}

// No-op, no service is required for the Debugger
func (b *StackdriverDebuggerBroker) Deprovision(instanceID string, details brokerapi.DeprovisionDetails) error {
	return nil
}

// Creates a service account with access to Stackdriver Debugger
func (b *StackdriverDebuggerBroker) Bind(instanceID, bindingID string, details brokerapi.BindDetails) (models.ServiceBindingCredentials, error) {
	out, err := utils.SetParameter(details.RawParameters, "role", "clouddebugger.agent")
	if err != nil {
		return models.ServiceBindingCredentials{}, err
	}
	details.RawParameters = out

	// Create account
	newBinding, err := b.AccountManager.CreateCredentials(instanceID, bindingID, details, models.ServiceInstanceDetails{})

	if err != nil {
		return models.ServiceBindingCredentials{}, err
	}

	return newBinding, nil
}
