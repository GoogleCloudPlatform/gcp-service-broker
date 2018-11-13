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
	"encoding/json"
	"fmt"

	"code.cloudfoundry.org/lager"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/broker"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/varcontext"
	"github.com/pivotal-cf/brokerapi"
	"golang.org/x/oauth2/jwt"
)

// TfServiceDefinitionV1 is the first version of user defined services.
type TfServiceDefinitionV1 struct {
	Version           int                         `yaml:"version" validate:"required,eq=1"`
	Name              string                      `yaml:"name" validate:"required"`
	Id                string                      `yaml:"id" validate:"required,uuid"`
	Description       string                      `yaml:"description" validate:"required"`
	DisplayName       string                      `yaml:"display_name" validate:"required"`
	ImageUrl          string                      `yaml:"image_url" validate:"url"`
	DocumentationUrl  string                      `yaml:"documentation_url" validate:"url"`
	SupportUrl        string                      `yaml:"support_url" validate:"url"`
	Tags              []string                    `yaml:"tags,flow"`
	Plans             []broker.ServicePlan        `yaml:"plans" validate:"required,dive"`
	ProvisionSettings TfServiceDefinitionV1Action `yaml:"provision" validate:"required,dive"`
	BindSettings      TfServiceDefinitionV1Action `yaml:"bind" validate:"required,dive"`
	Examples          []broker.ServiceExample     `yaml:"examples" validate:"required,dive"`

	// Internal SHOULD be set to true for Google maintained services.
	Internal bool `yaml:"-"`
}

// TfServiceDefinitionV1Action holds information needed to process user inputs
// for a single provision or bind call.
type TfServiceDefinitionV1Action struct {
	PlanInputs []broker.BrokerVariable      `yaml:"plan_inputs" validate:"dive"`
	UserInputs []broker.BrokerVariable      `yaml:"user_inputs" validate:"dive"`
	Computed   []varcontext.DefaultVariable `yaml:"computed_inputs" validate:"dive"`
	Template   string                       `yaml:"template" validate:"hcl"`
	Outputs    []broker.BrokerVariable      `yaml:"outputs" validate:"dive"`
}

// ToService converts the flat TfServiceDefinitionV1 into a broker.ServiceDefinition
// that the registry can use.
func (tfb *TfServiceDefinitionV1) ToService() (*broker.ServiceDefinition, error) {
	osbDefinition := broker.Service{
		Service: brokerapi.Service{
			ID:            tfb.Id,
			Name:          tfb.Name,
			Description:   tfb.Description,
			Bindable:      true,
			PlanUpdatable: false,
			Metadata: &brokerapi.ServiceMetadata{
				DisplayName:      tfb.Name,
				LongDescription:  tfb.Description,
				DocumentationUrl: tfb.DocumentationUrl,
				SupportUrl:       tfb.SupportUrl,
				ImageUrl:         tfb.ImageUrl,
			},
			Tags: tfb.Tags,
		},

		Plans: tfb.Plans,
	}

	defaultServiceDefinition, err := json.Marshal(osbDefinition)
	if err != nil {
		return nil, err
	}

	return &broker.ServiceDefinition{
		Name:                     tfb.Name,
		DefaultServiceDefinition: string(defaultServiceDefinition),
		ProvisionInputVariables:  tfb.ProvisionSettings.UserInputs,
		ProvisionComputedVariables: append(tfb.ProvisionSettings.Computed, varcontext.DefaultVariable{
			Name:      "tf_id",
			Default:   "tf:${request.instance_id}:",
			Overwrite: true,
		}),
		BindInputVariables: tfb.BindSettings.UserInputs,
		BindComputedVariables: append(tfb.BindSettings.Computed, varcontext.DefaultVariable{
			Name:      "tf_id",
			Default:   "tf:${request.instance_id}:${request.binding_id}",
			Overwrite: true,
		}),
		BindOutputVariables: append(tfb.ProvisionSettings.Outputs, tfb.BindSettings.Outputs...),
		PlanVariables:       append(tfb.ProvisionSettings.PlanInputs, tfb.BindSettings.PlanInputs...),
		Examples:            tfb.Examples,
		ProviderBuilder: func(projectId string, auth *jwt.Config, logger lager.Logger) broker.ServiceProvider {
			return NewTerraformProvider(projectId, auth, logger, *tfb)
		},
	}, nil
}

// generateTfId creates a unique id for a given provision/bind combination that
// will be consistent across calls. This ID will be used in LastOperation polls
// as well as to uniquely identify the workspace.
func generateTfId(instanceId, bindingId string) string {
	return fmt.Sprintf("tf:%s:%s", instanceId, bindingId)
}
