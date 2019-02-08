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
	"fmt"

	"code.cloudfoundry.org/lager"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/broker"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/providers/tf/wrapper"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/validation"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/varcontext"
	"github.com/GoogleCloudPlatform/gcp-service-broker/utils"
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
	Plans             []TfServiceDefinitionV1Plan `yaml:"plans" validate:"required,dive"`
	ProvisionSettings TfServiceDefinitionV1Action `yaml:"provision" validate:"required,dive"`
	BindSettings      TfServiceDefinitionV1Action `yaml:"bind" validate:"required,dive"`
	Examples          []broker.ServiceExample     `yaml:"examples" validate:"required,dive"`

	// Internal SHOULD be set to true for Google maintained services.
	Internal bool `yaml:"-"`
}

// TfServiceDefinitionV1Plan represents a service plan in a human-friendly format
// that can be converted into an OSB compatible plan.
type TfServiceDefinitionV1Plan struct {
	Name        string            `yaml:"name" validate:"required"`
	Id          string            `yaml:"id" validate:"required,uuid"`
	Description string            `yaml:"description" validate:"required"`
	DisplayName string            `yaml:"display_name" validate:"required"`
	Bullets     []string          `yaml:"bullets,omitempty"`
	Free        bool              `yaml:"free,omitempty"`
	Properties  map[string]string `yaml:"properties" validate:"required"`
}

// Converts this plan definition to a broker.ServicePlan.
func (plan *TfServiceDefinitionV1Plan) ToPlan() broker.ServicePlan {
	masterPlan := brokerapi.ServicePlan{
		ID:          plan.Id,
		Description: plan.Description,
		Name:        plan.Name,
		Free:        brokerapi.FreeValue(plan.Free),
		Metadata: &brokerapi.ServicePlanMetadata{
			Bullets:     plan.Bullets,
			DisplayName: plan.DisplayName,
		},
	}

	return broker.ServicePlan{
		ServicePlan:       masterPlan,
		ServiceProperties: plan.Properties,
	}
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

// ValidateTemplateIO makes sure that the inputs supplied by the user are a
// superset of the inputs needed by the Terraform template, and the template
// outputs match the outputs.
func (action *TfServiceDefinitionV1Action) ValidateTemplateIO() error {
	if err := action.validateTemplateInputs(); err != nil {
		return err
	}

	return action.validateTemplateOutputs()
}

// validateTemplateInputs checks that all the inputs of the Terraform template
// are defined by the service.
func (action *TfServiceDefinitionV1Action) validateTemplateInputs() error {
	inputs := utils.NewStringSet()

	for _, in := range action.PlanInputs {
		inputs.Add(in.FieldName)
	}

	for _, in := range action.UserInputs {
		inputs.Add(in.FieldName)
	}

	for _, in := range action.Computed {
		inputs.Add(in.Name)
	}

	tfModule := wrapper.ModuleDefinition{Definition: action.Template}
	tfIn, err := tfModule.Inputs()
	if err != nil {
		return err
	}

	missingFields := utils.NewStringSet(tfIn...).Minus(inputs).ToSlice()
	if len(missingFields) > 0 {
		return fmt.Errorf("The Terraform template requires the fields %v which are missing from the declared inputs.", missingFields)
	}

	return nil
}

// validateTemplateOutputs checks that the Terraform template outputs match
// the names of the defined outputs.
func (action *TfServiceDefinitionV1Action) validateTemplateOutputs() error {
	definedOutputs := utils.NewStringSet()

	for _, in := range action.Outputs {
		definedOutputs.Add(in.FieldName)
	}

	tfModule := wrapper.ModuleDefinition{Definition: action.Template}
	tfOut, err := tfModule.Outputs()
	if err != nil {
		return err
	}

	if !definedOutputs.Equals(utils.NewStringSet(tfOut...)) {
		return fmt.Errorf("The Terraform template outputs %v MUST match the service declared outputs %v.", tfOut, definedOutputs)
	}

	return nil
}

// Validate checks the service definition for semantic errors.
func (tfb *TfServiceDefinitionV1) Validate() error {
	if err := validation.ValidateStruct(tfb); err != nil {
		return err
	}

	if err := tfb.ProvisionSettings.ValidateTemplateIO(); err != nil {
		return err
	}

	return tfb.BindSettings.ValidateTemplateIO()
}

// ToService converts the flat TfServiceDefinitionV1 into a broker.ServiceDefinition
// that the registry can use.
func (tfb *TfServiceDefinitionV1) ToService(executor wrapper.TerraformExecutor) (*broker.ServiceDefinition, error) {
	if err := tfb.Validate(); err != nil {
		return nil, err
	}

	var rawPlans []broker.ServicePlan
	for _, plan := range tfb.Plans {
		rawPlans = append(rawPlans, plan.ToPlan())
	}

	return &broker.ServiceDefinition{
		Id:               tfb.Id,
		Name:             tfb.Name,
		Description:      tfb.Description,
		Bindable:         true,
		PlanUpdateable:   false,
		DisplayName:      tfb.DisplayName,
		DocumentationUrl: tfb.DocumentationUrl,
		SupportUrl:       tfb.SupportUrl,
		ImageUrl:         tfb.ImageUrl,
		Tags:             tfb.Tags,
		Plans:            rawPlans,

		ProvisionInputVariables: tfb.ProvisionSettings.UserInputs,
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
			jobRunner := NewTfJobRunnerForProject(projectId)
			jobRunner.Executor = executor
			return NewTerraformProvider(jobRunner, logger, constDefn)
		},
	}, nil
}

// generateTfId creates a unique id for a given provision/bind combination that
// will be consistent across calls. This ID will be used in LastOperation polls
// as well as to uniquely identify the workspace.
func generateTfId(instanceId, bindingId string) string {
	return fmt.Sprintf("tf:%s:%s", instanceId, bindingId)
}

// NewExampleTfServiceDefinition creates a new service defintition with sample
// values for the service broker suitable to give a user a template to manually
// edit.
func NewExampleTfServiceDefinition() TfServiceDefinitionV1 {
	return TfServiceDefinitionV1{
		Version:          1,
		Name:             "example-service",
		Id:               "00000000-0000-0000-0000-000000000000",
		Description:      "a longer service description",
		DisplayName:      "Example Service",
		ImageUrl:         "https://example.com/icon.jpg",
		DocumentationUrl: "https://example.com",
		SupportUrl:       "https://example.com/support.html",
		Tags:             []string{"gcp", "example", "service"},
		Plans: []TfServiceDefinitionV1Plan{
			{
				Id:          "00000000-0000-0000-0000-000000000001",
				Name:        "example-email-plan",
				DisplayName: "example.com email builder",
				Description: "Builds emails for example.com.",
				Bullets:     []string{"information point 1", "information point 2", "some caveat here"},
				Free:        false,
				Properties: map[string]string{
					"domain": "example.com",
				},
			},
		},
		ProvisionSettings: TfServiceDefinitionV1Action{
			PlanInputs: []broker.BrokerVariable{
				{
					FieldName: "domain",
					Type:      broker.JsonTypeString,
					Details:   "The domain name",
					Required:  true,
				},
			},
			UserInputs: []broker.BrokerVariable{
				{
					FieldName: "username",
					Type:      broker.JsonTypeString,
					Details:   "The username to create",
					Required:  true,
				},
			},
			Template: `
			variable domain {type = "string"}
			variable username {type = "string"}

			output email {value = "${var.username}@${var.domain}"}
			`,
			Outputs: []broker.BrokerVariable{
				{
					FieldName: "email",
					Type:      broker.JsonTypeString,
					Details:   "The combined email address",
					Required:  true,
				},
			},
		},
		BindSettings: TfServiceDefinitionV1Action{
			Computed: []varcontext.DefaultVariable{
				{Name: "address", Default: `${instance.details["email"]}`, Overwrite: true},
			},
			Template: `
			resource "random_string" "password" {
			  length = 16
			  special = true
			  override_special = "/@\" "
			}

			output uri {value = "smtp://${var.address}:${random_string.password.result}@smtp.mycompany.com"}
			`,
			Outputs: []broker.BrokerVariable{
				{
					FieldName: "uri",
					Type:      broker.JsonTypeString,
					Details:   "The uri to use to connect to this service",
					Required:  true,
				},
			},
		},
		Examples: []broker.ServiceExample{
			{
				Name:            "Example",
				Description:     "Examples are used for documenting your service AND as integration tests.",
				PlanId:          "00000000-0000-0000-0000-000000000001",
				ProvisionParams: map[string]interface{}{"username": "my-account"},
				BindParams:      map[string]interface{}{},
			},
		},
	}
}
