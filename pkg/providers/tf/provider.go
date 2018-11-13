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
	"encoding/json"

	"code.cloudfoundry.org/lager"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/models"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/broker"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/providers/tf/wrapper"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/varcontext"
	"github.com/GoogleCloudPlatform/gcp-service-broker/utils"
	"github.com/pivotal-cf/brokerapi"
	"golang.org/x/oauth2/jwt"
)

func NewTerraformProvider(projectId string, auth *jwt.Config, logger lager.Logger, serviceDefinition TfServiceDefinitionV1) broker.ServiceProvider {
	return &terraformProvider{
		serviceDefinition: serviceDefinition,
		jobRunner:         TfJobRunner{ProjectId: projectId, ServiceAccount: utils.GetServiceAccountJson()},
		logger:            logger.Session("terraform-" + serviceDefinition.Name),
	}
}

type terraformProvider struct {
	// broker_base.BrokerBase
	logger            lager.Logger
	jobRunner         TfJobRunner
	serviceDefinition TfServiceDefinitionV1
}

// Provision creates the necessary resources that an instance of this service
// needs to operate.
func (provider *terraformProvider) Provision(ctx context.Context, provisionContext *varcontext.VarContext) (models.ServiceInstanceDetails, error) {
	provider.logger.Info("provision", lager.Data{
		"context": provisionContext.ToMap(),
	})

	tfId, err := provider.create(ctx, provisionContext, provider.serviceDefinition.ProvisionSettings)
	if err != nil {
		return models.ServiceInstanceDetails{}, err
	}

	return models.ServiceInstanceDetails{
		OperationId:   tfId,
		OperationType: models.ProvisionOperationType,
	}, nil
}

func (provider *terraformProvider) Bind(ctx context.Context, bindContext *varcontext.VarContext) (map[string]interface{}, error) {
	provider.logger.Info("bind", lager.Data{
		"context": bindContext.ToMap(),
	})

	tfId, err := provider.create(ctx, bindContext, provider.serviceDefinition.BindSettings)
	if err != nil {
		return nil, err
	}

	if err := provider.jobRunner.Wait(ctx, tfId); err != nil {
		return nil, err
	}

	return provider.jobRunner.Outputs(ctx, tfId, wrapper.DefaultInstanceName)
}

func (provider *terraformProvider) create(ctx context.Context, vars *varcontext.VarContext, action *TfServiceDefinitionV1Action) (string, error) {
	tfId := vars.GetString("tf_id")
	if err := vars.Error(); err != nil {
		return "", err
	}

	workspace, err := wrapper.NewWorkspace(vars.ToMap(), action.Template)
	if err != nil {
		return tfId, err
	}

	if err := provider.jobRunner.StageJob(ctx, tfId, workspace); err != nil {
		return tfId, err
	}

	return tfId, provider.jobRunner.Create(ctx, tfId)
}

func (provider *terraformProvider) BuildInstanceCredentials(ctx context.Context, bindRecord models.ServiceBindingCredentials, instanceRecord models.ServiceInstanceDetails) (map[string]interface{}, error) {
	vc, err := varcontext.Builder().
		MergeJsonObject(json.RawMessage(bindRecord.OtherDetails)).
		MergeJsonObject(json.RawMessage(instanceRecord.OtherDetails)).
		Build()
	if err != nil {
		return nil, err
	}

	return vc.ToMap(), nil
}

func (provider *terraformProvider) Unbind(ctx context.Context, instanceRecord models.ServiceInstanceDetails, bindRecord models.ServiceBindingCredentials) error {
	tfId := generateTfId(instanceRecord.ID, bindRecord.BindingId)
	provider.logger.Info("unbind", lager.Data{
		"instance": instanceRecord.ID,
		"binding":  bindRecord.ID,
		"tfId":     tfId,
	})

	if err := provider.jobRunner.Destroy(ctx, tfId); err != nil {
		return err
	}

	return provider.jobRunner.Wait(ctx, tfId)
}

// Deprovision deprovisions the service.
// If the deprovision is asynchronous (results in a long-running job), then operationId is returned.
// If no error and no operationId are returned, then the deprovision is expected to have been completed successfully.
func (provider *terraformProvider) Deprovision(ctx context.Context, instance models.ServiceInstanceDetails, details brokerapi.DeprovisionDetails) (operationId *string, err error) {
	provider.logger.Info("terraform-deprovision", lager.Data{
		"instance": instance.ID,
	})

	tfId := generateTfId(instance.ID, "")
	if err := provider.jobRunner.Destroy(ctx, tfId); err != nil {
		return nil, err
	}

	return &tfId, nil
}

func (provider *terraformProvider) PollInstance(ctx context.Context, instance models.ServiceInstanceDetails) (bool, error) {
	return provider.jobRunner.Status(ctx, generateTfId(instance.ID, ""))
}

func (provider *terraformProvider) ProvisionsAsync() bool {
	return true
}

func (provider *terraformProvider) DeprovisionsAsync() bool {
	return true
}

// UpdateInstanceDetails updates the ServiceInstanceDetails with the most recent state from GCP.
// This function is optional, but will be called after async provisions, updates, and possibly
// on broker version changes.
// Return a nil error if you choose not to implement this function.
func (provider *terraformProvider) UpdateInstanceDetails(ctx context.Context, instance *models.ServiceInstanceDetails) error {
	tfId := generateTfId(instance.ID, "")

	outs, err := provider.jobRunner.Outputs(ctx, tfId, wrapper.DefaultInstanceName)
	if err != nil {
		return err
	}

	return instance.SetOtherDetails(outs)
}
