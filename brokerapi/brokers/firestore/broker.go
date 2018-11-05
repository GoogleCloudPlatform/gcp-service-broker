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

package firestore

import (
	"context"

	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/broker_base"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/models"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/varcontext"
	"github.com/pivotal-cf/brokerapi"
)

// FirestoreBroker is the service-broker back-end for creating and binding Firestore clients.
type FirestoreBroker struct {
	broker_base.BrokerBase
}

// Provision is a no-op call because only service accounts need to be bound/unbound for Firestore.
func (b *FirestoreBroker) Provision(ctx context.Context, provisionContext *varcontext.VarContext) (models.ServiceInstanceDetails, error) {
	return models.ServiceInstanceDetails{}, nil
}

// Deprovision is a no-op call because only service accounts need to be bound/unbound for Firestore.
func (b *FirestoreBroker) Deprovision(ctx context.Context, instance models.ServiceInstanceDetails, details brokerapi.DeprovisionDetails) (*string, error) {
	return nil, nil
}
