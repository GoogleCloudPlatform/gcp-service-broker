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
	"context"

	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/providers/builtin/base"
	"github.com/GoogleCloudPlatform/gcp-service-broker/db_service/models"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/varcontext"
	"github.com/pivotal-cf/brokerapi"
)

// InstanceInformation holds the details needed to bind a service to a DatastoreBroker.
type InstanceInformation struct {
	Namespace string `json:"namespace,omitempty"`
}

// DatastoreBroker is the service-broker back-end for creating and binding Datastore instances.
type DatastoreBroker struct {
	base.BrokerBase
}

// Provision stores the namespace for future reference.
func (b *DatastoreBroker) Provision(ctx context.Context, provisionContext *varcontext.VarContext) (models.ServiceInstanceDetails, error) {
	// return models.ServiceInstanceDetails{}, nil

	ii := InstanceInformation{
		Namespace: provisionContext.GetString("namespace"),
	}

	if err := provisionContext.Error(); err != nil {
		return models.ServiceInstanceDetails{}, err
	}

	details := models.ServiceInstanceDetails{}
	if err := details.SetOtherDetails(ii); err != nil {
		return models.ServiceInstanceDetails{}, err
	}

	return details, nil
}

// Deprovision is a no-op call because only service accounts need to be bound/unbound for Datastore.
func (b *DatastoreBroker) Deprovision(ctx context.Context, instance models.ServiceInstanceDetails, details brokerapi.DeprovisionDetails) (*string, error) {
	return nil, nil
}
