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

package stackdriver_trace

import (
	"gcp-service-broker/brokerapi/brokers/account_managers"
	"gcp-service-broker/brokerapi/brokers/broker_base"
	"gcp-service-broker/brokerapi/brokers/models"
	"net/http"

	"code.cloudfoundry.org/lager"
)

type StackdriverTraceBroker struct {
	Client                *http.Client
	ProjectId             string
	Logger                lager.Logger
	ServiceAccountManager *account_managers.ServiceAccountManager
	broker_base.BrokerBase
}

type InstanceInformation struct {
}

// No-op, no serivce is required for Stackdriver Trace
func (b *StackdriverTraceBroker) Provision(instanceId string, details models.ProvisionDetails, plan models.ServicePlan) (models.ServiceInstanceDetails, error) {
	return models.ServiceInstanceDetails{}, nil
}

// No-op, no serivce is required for Stackdriver Trace
func (b *StackdriverTraceBroker) Deprovision(instanceID string, details models.DeprovisionDetails) error {
	return nil
}

// Creates a service account with access to Stackdriver Trace
func (b *StackdriverTraceBroker) Bind(instanceID, bindingID string, details models.BindDetails) (models.ServiceBindingCredentials, error) {
	if details.Parameters == nil {
		b.Logger.Info("the parameters are nil!")
		details.Parameters = make(map[string]interface{})
	}
	details.Parameters["role"] = "cloudtrace.agent"

	// Create account
	newBinding, err := b.ServiceAccountManager.CreateAccountInGoogle(instanceID, bindingID, details, models.ServiceInstanceDetails{})

	if err != nil {
		return models.ServiceBindingCredentials{}, err
	}

	return newBinding, nil
}
