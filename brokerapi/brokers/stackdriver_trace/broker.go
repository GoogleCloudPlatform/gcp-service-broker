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
	"net/http"
	"code.cloudfoundry.org/lager"
	"gcp-service-broker/brokerapi/brokers/broker_base"
	"gcp-service-broker/brokerapi/brokers/models"
)

type StackdriverTraceBroker struct {
	Client    *http.Client
	ProjectId string
	Logger    lager.Logger
	broker_base.BrokerBase
}

type InstanceInformation struct {
}

// Creates a service account for Stackdriver Trace
func (b *StackdriverTraceBroker) Provision(instanceId string, details models.ProvisionDetails, plan models.PlanDetails) (models.ServiceInstanceDetails, error) {
	return models.ServiceInstanceDetails{}, nil
}

// Deletes the topic associated with the given instanceID
func (b *StackdriverTraceBroker) Deprovision(instanceID string, details models.DeprovisionDetails) error {
	return nil
}

func (b *StackdriverTraceBroker) Bind(instanceID, bindingID string, details models.BindDetails) (models.ServiceBindingCredentials, error) {
	if details.Parameters == nil {
		b.Logger.Info("the parameters are nil!")
		details.Parameters = make(map[string]interface{})
	}
	details.Parameters["role"] = "cloudtrace.agent"
	return b.BrokerBase.Bind(instanceID, bindingID, details)
}
