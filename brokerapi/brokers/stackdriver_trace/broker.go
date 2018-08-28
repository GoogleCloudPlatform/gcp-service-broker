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

package stackdriver_trace

import (
	"context"

	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/broker_base"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/models"
	"github.com/pivotal-cf/brokerapi"
)

type StackdriverTraceBroker struct {
	broker_base.BrokerBase
}

// No-op, no serivce is required for Stackdriver Trace
func (b *StackdriverTraceBroker) Provision(instanceId string, details brokerapi.ProvisionDetails, plan models.ServicePlan) (models.ServiceInstanceDetails, error) {
	return models.ServiceInstanceDetails{}, nil
}

// No-op, no serivce is required for Stackdriver Trace
func (b *StackdriverTraceBroker) Deprovision(ctx context.Context, instance models.ServiceInstanceDetails, details brokerapi.DeprovisionDetails) error {
	return nil
}

// Creates a service account with access to Stackdriver Trace
func (b *StackdriverTraceBroker) Bind(instanceID, bindingID string, details brokerapi.BindDetails) (models.ServiceBindingCredentials, error) {
	return b.AccountManager.CreateAccountWithRoles(bindingID, []string{"cloudtrace.agent"})
}
