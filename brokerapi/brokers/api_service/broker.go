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

package api_service

import (
	"code.cloudfoundry.org/lager"
	"gcp-service-broker/brokerapi/brokers/broker_base"
	"gcp-service-broker/brokerapi/brokers/models"
	"net/http"
)

type ApiServiceBroker struct {
	Client         *http.Client
	ProjectId      string
	Logger         lager.Logger
	AccountManager models.AccountManager

	broker_base.BrokerBase
}

func (b *ApiServiceBroker) Provision(instanceId string, details models.ProvisionDetails, plan models.PlanDetails) (models.ServiceInstanceDetails, error) {

	return models.ServiceInstanceDetails{}, nil
}

func (b *ApiServiceBroker) Deprovision(instanceID string, details models.DeprovisionDetails) error {

	return nil
}
