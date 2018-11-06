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

package tf

import (
	"context"
	"errors"

	"code.cloudfoundry.org/lager"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/broker_base"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/models"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/broker"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/varcontext"
	"github.com/pivotal-cf/brokerapi"
	"golang.org/x/oauth2/jwt"
)

func NewTerraformProvider(projectId string, auth *jwt.Config, logger lager.Logger, serviceDefinition TfServiceDefinitionV1) broker.ServiceProvider {
	return &terraformProvider{
		BrokerBase:        broker_base.NewBrokerBase(projectId, auth, logger),
		serviceDefinition: serviceDefinition,
	}
}

type terraformProvider struct {
	broker_base.BrokerBase

	serviceDefinition TfServiceDefinitionV1
}

// Provision creates the necessary resources that an instance of this service
// needs to operate.
func (provider *terraformProvider) Provision(ctx context.Context, provisionContext *varcontext.VarContext) (models.ServiceInstanceDetails, error) {
	return models.ServiceInstanceDetails{}, nil
}

// Deprovision deprovisions the service.
// If the deprovision is asynchronous (results in a long-running job), then operationId is returned.
// If no error and no operationId are returned, then the deprovision is expected to have been completed successfully.
func (provider *terraformProvider) Deprovision(ctx context.Context, instance models.ServiceInstanceDetails, details brokerapi.DeprovisionDetails) (operationId *string, err error) {
	return nil, nil
}

func (provider *terraformProvider) PollInstance(ctx context.Context, instance models.ServiceInstanceDetails) (bool, error) {
	return false, nil
}

func (provider *terraformProvider) ProvisionsAsync() bool {
	return false
}

func (provider *terraformProvider) DeprovisionsAsync() bool {
	return false
}

// UpdateInstanceDetails updates the ServiceInstanceDetails with the most recent state from GCP.
// This function is optional, but will be called after async provisions, updates, and possibly
// on broker version changes.
// Return a nil error if you choose not to implement this function.
func (provider *terraformProvider) UpdateInstanceDetails(ctx context.Context, instance *models.ServiceInstanceDetails) error {
	return errors.New("Update is not supported.")
}
