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

package datastore

import (
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/broker_base"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/models"
	"github.com/pivotal-cf/brokerapi"
)

type DatastoreBroker struct {
	broker_base.BrokerBase
}

// No-op, no service is required for Datastore
func (b *DatastoreBroker) Provision(instanceId string, details brokerapi.ProvisionDetails, plan models.ServicePlan) (models.ServiceInstanceDetails, error) {
	return models.ServiceInstanceDetails{}, nil
}

// No-op, no service is required for Datastore
func (b *DatastoreBroker) Deprovision(instanceID string, details brokerapi.DeprovisionDetails) error {
	return nil
}

// Creates a service account with access to Datastore
func (b *DatastoreBroker) Bind(instanceID, bindingID string, details brokerapi.BindDetails) (models.ServiceBindingCredentials, error) {
	return b.AccountManager.CreateAccountWithRoles(instanceID, []string{"datastore.user"})
}
