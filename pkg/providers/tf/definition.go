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
	Version           int                         `yaml:"version"`
	Name              string                      `yaml:"name"`
	Id                string                      `yaml:"id"`
	Description       string                      `yaml:"description"`
	DisplayName       string                      `yaml:"display_name"`
	ImageUrl          string                      `yaml:"image_url"`
	DocumentationUrl  string                      `yaml:"documentation_url"`
	SupportUrl        string                      `yaml:"support_url"`
	Tags              []string                    `yaml:"tags,flow"`
	Plans             []TfServiceDefinitionV1Plan `yaml:"plans"`
	ProvisionSettings TfServiceDefinitionV1Action `yaml:"provision"`
	BindSettings      TfServiceDefinitionV1Action `yaml:"bind"`
	Examples          []broker.ServiceExample     `yaml:"examples"`

	// Internal SHOULD be set to true for Google maintained services.
	Internal bool `yaml:"-"`
}

// TfServiceDefinitionV1Plan represents a service plan in a human-friendly format
// that can be converted into an OSB compatible plan.
type TfServiceDefinitionV1Plan struct {
	Name               string                 `yaml:"name"`
	Id                 string                 `yaml:"id"`
	Description        string                 `yaml:"description"`
	DisplayName        string                 `yaml:"display_name"`
	Bullets            []string               `yaml:"bullets,omitempty"`
	Free               bool                   `yaml:"free,omitempty"`
	Properties         map[string]string      `yaml:"properties"`
	ProvisionOverrides map[string]interface{} `yaml:"provision_overrides,omitempty"`
	BindOverrides      map[string]interface{} `yaml:"bind_overrides,omitempty"`
}

var _ validation.Validatable = (*TfServiceDefinitionV1Plan)(nil)

// Validate implements validation.Validatable.
func (plan *TfServiceDefinitionV1Plan) Validate() (errs *validation.FieldError) {
	return errs.Also(
		validation.ErrIfBlank(plan.Name, "name"),
		validation.ErrIfNotUUID(plan.Id, "id"),
		validation.ErrIfBlank(plan.Description, "description"),
		validation.ErrIfBlank(plan.DisplayName, "display_name"),
	)
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
		ServicePlan:        masterPlan,
		ServiceProperties:  plan.Properties,
		ProvisionOverrides: plan.ProvisionOverrides,
		BindOverrides:      plan.BindOverrides,
	}
}

// TfServiceDefinitionV1Action holds information needed to process user inputs
// for a single provision or bind call.
type TfServiceDefinitionV1Action struct {
	PlanInputs []broker.BrokerVariable      `yaml:"plan_inputs"`
	UserInputs []broker.BrokerVariable      `yaml:"user_inputs"`
	Computed   []varcontext.DefaultVariable `yaml:"computed_inputs"`
	Template   string                       `yaml:"template"`
	Outputs    []broker.BrokerVariable      `yaml:"outputs"`
}

var _ validation.Validatable = (*TfServiceDefinitionV1Action)(nil)

// Validate implements validation.Validatable.
func (action *TfServiceDefinitionV1Action) Validate() (errs *validation.FieldError) {
	for i, v := range action.PlanInputs {
		errs = errs.Also(v.Validate().ViaFieldIndex("plan_inputs", i))
	}

	for i, v := range action.UserInputs {
		errs = errs.Also(v.Validate().ViaFieldIndex("user_inputs", i))
	}

	for i, v := range action.Computed {
		errs = errs.Also(v.Validate().ViaFieldIndex("computed_inputs", i))
	}

	errs = errs.Also(
		validation.ErrIfNotHCL(action.Template, "template"),
		action.validateTemplateInputs().ViaField("template"),
		action.validateTemplateOutputs().ViaField("template"),
	)

	for i, v := range action.Outputs {
		errs = errs.Also(v.Validate().ViaFieldIndex("outputs", i))
	}

	return errs
}

func (action *TfServiceDefinitionV1Action) ValidateTemplateIO() (errs *validation.FieldError) {
	return errs.Also(
		action.validateTemplateInputs().ViaField("template"),
		action.validateTemplateOutputs().ViaField("template"),
	)
}

// validateTemplateInputs checks that all the inputs of the Terraform template
// are defined by the service.
func (action *TfServiceDefinitionV1Action) validateTemplateInputs() (errs *validation.FieldError) {
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
		return &validation.FieldError{
			Message: err.Error(),
		}
	}

	missingFields := utils.NewStringSet(tfIn...).Minus(inputs).ToSlice()
	if len(missingFields) > 0 {
		return &validation.FieldError{
			Message: "fields used but not declared",
			Paths:   missingFields,
		}
	}

	return nil
}

// validateTemplateOutputs checks that the Terraform template outputs match
// the names of the defined outputs.
func (action *TfServiceDefinitionV1Action) validateTemplateOutputs() (errs *validation.FieldError) {
	definedOutputs := utils.NewStringSet()

	for _, in := range action.Outputs {
		definedOutputs.Add(in.FieldName)
	}

	tfModule := wrapper.ModuleDefinition{Definition: action.Template}
	tfOut, err := tfModule.Outputs()
	if err != nil {
		return &validation.FieldError{
			Message: err.Error(),
		}
	}

	if !definedOutputs.Equals(utils.NewStringSet(tfOut...)) {
		return &validation.FieldError{
			Message: fmt.Sprintf("template outputs %v must match declared outputs %v", tfOut, definedOutputs),
		}
	}

	return nil
}

var _ validation.Validatable = (*TfServiceDefinitionV1)(nil)

// Validate checks the service definition for semantic errors.
func (tfb *TfServiceDefinitionV1) Validate() (errs *validation.FieldError) {

	if tfb.Version != 1 {
		errs = errs.Also(validation.ErrInvalidValue(tfb.Version, "version"))
	}

	errs = errs.Also(
		validation.ErrIfBlank(tfb.Name, "name"),
		validation.ErrIfNotUUID(tfb.Id, "id"),
		validation.ErrIfBlank(tfb.Description, "description"),
		validation.ErrIfBlank(tfb.DisplayName, "display_name"),
		validation.ErrIfNotURL(tfb.ImageUrl, "image_url"),
		validation.ErrIfNotURL(tfb.DocumentationUrl, "documentation_url"),
		validation.ErrIfNotURL(tfb.SupportUrl, "support_url"),
	)

	for i, v := range tfb.Plans {
		errs = errs.Also(v.Validate().ViaFieldIndex("plans", i))
	}

	errs = errs.Also(tfb.ProvisionSettings.Validate().ViaField("provision"))
	errs = errs.Also(tfb.BindSettings.Validate().ViaField("bind"))

	for i, v := range tfb.Examples {
		errs = errs.Also(v.Validate().ViaFieldIndex("examples", i))
	}

	return errs
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

	// Bindings get special computed properties because the broker didn't
	// originally support injecting plan variables into a binding
	// to fix that, we auto-inject the properties from the plan to make it look
	// like they were to the TF template.
	bindComputed := []varcontext.DefaultVariable{}
	for _, pi := range tfb.BindSettings.PlanInputs {
		bindComputed = append(bindComputed, varcontext.DefaultVariable{
			Name:      pi.FieldName,
			Default:   fmt.Sprintf("${request.plan_properties[%q]}", pi.FieldName),
			Overwrite: true,
			Type:      string(pi.Type),
		})
	}

	bindComputed = append(bindComputed, tfb.BindSettings.Computed...)
	bindComputed = append(bindComputed, varcontext.DefaultVariable{
		Name:      "tf_id",
		Default:   "tf:${request.instance_id}:${request.binding_id}",
		Overwrite: true,
	})

	constDefn := *tfb
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
		BindInputVariables:    tfb.BindSettings.UserInputs,
		BindComputedVariables: bindComputed,
		BindOutputVariables:   append(tfb.ProvisionSettings.Outputs, tfb.BindSettings.Outputs...),
		PlanVariables:         append(tfb.ProvisionSettings.PlanInputs, tfb.BindSettings.PlanInputs...),
		Examples:              tfb.Examples,
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
					"domain":                 "example.com",
					"password_special_chars": `@/ \"?`,
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
			PlanInputs: []broker.BrokerVariable{
				{
					FieldName: "password_special_chars",
					Type:      broker.JsonTypeString,
					Details:   "Supply your own list of special characters to use for string generation.",
					Required:  true,
				},
			},
			Computed: []varcontext.DefaultVariable{
				{Name: "domain", Default: `${request.plan_properties["domain"]}`, Overwrite: true},
				{Name: "address", Default: `${instance.details["email"]}`, Overwrite: true},
			},
			Template: `
			variable domain {type = "string"}
			variable address {type = "string"}
			variable password_special_chars {type = "string"}

			resource "random_string" "password" {
			  length = 16
			  special = true
				override_special = "${var.password_special_chars}"
			}

			output uri {value = "smtp://${var.address}:${random_string.password.result}@smtp.${var.domain}"}
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
