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
	"reflect"
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/broker"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/varcontext"
	"github.com/pivotal-cf/brokerapi"
)

func TestTfServiceDefinitionV1Action_ValidateTemplateIO(t *testing.T) {
	cases := map[string]struct {
		Action      TfServiceDefinitionV1Action
		ErrContains string
	}{
		"nomainal": {
			Action: TfServiceDefinitionV1Action{
				PlanInputs: []broker.BrokerVariable{{FieldName: "storage_class"}},
				UserInputs: []broker.BrokerVariable{{FieldName: "name"}},
				Computed:   []varcontext.DefaultVariable{{Name: "labels"}},
				Template: `
      	variable storage_class {type = "string"}
      	variable name {type = "string"}
      	variable labels {type = "string"}

      	output bucket_name {value = "${var.name}"}
      	`,
				Outputs: []broker.BrokerVariable{{FieldName: "bucket_name"}},
			},
			ErrContains: "",
		},
		"extra inputs okay": {
			Action: TfServiceDefinitionV1Action{
				PlanInputs: []broker.BrokerVariable{{FieldName: "storage_class"}},
				UserInputs: []broker.BrokerVariable{{FieldName: "name"}},
				Computed:   []varcontext.DefaultVariable{{Name: "labels"}},
				Template: `
      	variable storage_class {type = "string"}
      	`,
			},
			ErrContains: "",
		},
		"missing inputs": {
			Action: TfServiceDefinitionV1Action{
				PlanInputs: []broker.BrokerVariable{{FieldName: "storage_class"}},
				UserInputs: []broker.BrokerVariable{{FieldName: "name"}},
				Computed:   []varcontext.DefaultVariable{{Name: "labels"}},
				Template: `
        variable storage_class {type = "string"}
        variable not_defined {type = "string"}
        `,
			},
			ErrContains: "fields used but not declared: template.not_defined",
		},

		"extra template outputs": {
			Action: TfServiceDefinitionV1Action{
				Template: `
        output storage_class {value = "${var.name}"}
        output name {value = "${var.name}"}
        output labels {value = "${var.name}"}
        output bucket_name {value = "${var.name}"}
        `,
				Outputs: []broker.BrokerVariable{{FieldName: "bucket_name"}},
			},
			ErrContains: "template outputs [bucket_name labels name storage_class] must match declared outputs [bucket_name]:",
		},

		"missing template outputs": {
			Action: TfServiceDefinitionV1Action{
				Template: `
        `,
				Outputs: []broker.BrokerVariable{{FieldName: "bucket_name"}},
			},
			ErrContains: "template outputs [] must match declared outputs [bucket_name]:",
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			err := tc.Action.ValidateTemplateIO()
			if err == nil {
				if tc.ErrContains == "" {
					return
				}

				t.Fatalf("Expected error to contain %q, got: <nil>", tc.ErrContains)
			} else {
				if tc.ErrContains == "" {
					t.Fatalf("Expected no error, got: %v", err)
				}

				if !strings.Contains(err.Error(), tc.ErrContains) {
					t.Fatalf("Expected error to contain %q, got: %v", tc.ErrContains, err)
				}
			}
		})
	}
}

func TestNewExampleTfServiceDefinition(t *testing.T) {
	example := NewExampleTfServiceDefinition()

	if err := example.Validate(); err != nil {
		t.Fatalf("example service definition should be valid, but got error: %v", err)
	}
}

func TestTfServiceDefinitionV1Plan_ToPlan(t *testing.T) {
	cases := map[string]struct {
		Definition TfServiceDefinitionV1Plan
		Expected   broker.ServicePlan
	}{
		"full": {
			Definition: TfServiceDefinitionV1Plan{
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
			Expected: broker.ServicePlan{
				ServicePlan: brokerapi.ServicePlan{
					ID:          "00000000-0000-0000-0000-000000000001",
					Name:        "example-email-plan",
					Description: "Builds emails for example.com.",
					Free:        brokerapi.FreeValue(false),
					Metadata: &brokerapi.ServicePlanMetadata{
						Bullets:     []string{"information point 1", "information point 2", "some caveat here"},
						DisplayName: "example.com email builder",
					},
				},
				ServiceProperties: map[string]string{"domain": "example.com"}},
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			actual := tc.Definition.ToPlan()
			if !reflect.DeepEqual(actual, tc.Expected) {
				t.Fatalf("Expected: %v Actual: %v", tc.Expected, actual)
			}
		})
	}
}

func TestTfServiceDefinitionV1_ToService(t *testing.T) {
	definition := TfServiceDefinitionV1{
		Version:     1,
		Id:          "d34705c8-3edf-4ab8-93b3-d97f080da24c",
		Name:        "my-service-name",
		Description: "my-service-description",
		DisplayName: "My Service Name",

		ImageUrl:         "https://example.com/image.png",
		SupportUrl:       "https://example.com/support",
		DocumentationUrl: "https://example.com/docs",
		Plans:            []TfServiceDefinitionV1Plan{},

		ProvisionSettings: TfServiceDefinitionV1Action{
			PlanInputs: []broker.BrokerVariable{
				{
					FieldName: "plan-input-provision",
					Type:      "string",
					Details:   "description",
				},
			},
			UserInputs: []broker.BrokerVariable{
				{
					FieldName: "user-input-provision",
					Type:      "string",
					Details:   "description",
				},
			},
			Computed: []varcontext.DefaultVariable{{Name: "computed-input-provision", Default: ""}},
		},

		BindSettings: TfServiceDefinitionV1Action{
			PlanInputs: []broker.BrokerVariable{
				{
					FieldName: "plan-input-bind",
					Type:      "integer",
					Details:   "description",
				},
			},
			UserInputs: []broker.BrokerVariable{
				{
					FieldName: "user-input-bind",
					Type:      "string",
					Details:   "description",
				},
			},
			Computed: []varcontext.DefaultVariable{{Name: "computed-input-bind", Default: ""}},
		},

		Examples: []broker.ServiceExample{},
	}

	service, err := definition.ToService(nil)
	if err != nil {
		t.Fatal(err)
	}

	expectEqual := func(field string, expected, actual interface{}) {
		if !reflect.DeepEqual(expected, actual) {
			t.Errorf("Expected %q to be equal. Expected: %#v, Actual: %#v", field, expected, actual)
		}
	}

	t.Run("basic-info", func(t *testing.T) {
		expectEqual("Id", definition.Id, service.Id)
		expectEqual("Name", definition.Name, service.Name)
		expectEqual("Description", definition.Description, service.Description)
		expectEqual("Bindable", true, service.Bindable)
		expectEqual("PlanUpdateable", false, service.PlanUpdateable)
		expectEqual("DisplayName", definition.DisplayName, service.DisplayName)
		expectEqual("DocumentationUrl", definition.DocumentationUrl, service.DocumentationUrl)
		expectEqual("SupportUrl", definition.SupportUrl, service.SupportUrl)
		expectEqual("ImageUrl", definition.ImageUrl, service.ImageUrl)
		expectEqual("Tags", definition.Tags, service.Tags)
	})

	t.Run("vars", func(t *testing.T) {
		expectEqual("ProvisionInputVariables", definition.ProvisionSettings.UserInputs, service.ProvisionInputVariables)
		expectEqual("ProvisionComputedVariables", []varcontext.DefaultVariable{
			{
				Name:      "computed-input-provision",
				Default:   "",
				Overwrite: false,
			},
			{
				Name:      "tf_id",
				Default:   "tf:${request.instance_id}:",
				Overwrite: true,
			},
		}, service.ProvisionComputedVariables)
		expectEqual("PlanVariables", append(definition.ProvisionSettings.PlanInputs, definition.BindSettings.PlanInputs...), service.PlanVariables)
		expectEqual("BindInputVariables", definition.BindSettings.UserInputs, service.BindInputVariables)
		expectEqual("BindComputedVariables", []varcontext.DefaultVariable{
			{Name: "plan-input-bind", Default: "${request.plan_properties[\"plan-input-bind\"]}", Overwrite: true, Type: "integer"},
			{Name: "computed-input-bind", Default: "", Overwrite: false, Type: ""},
			{Name: "tf_id", Default: "tf:${request.instance_id}:${request.binding_id}", Overwrite: true, Type: ""},
		}, service.BindComputedVariables)
		expectEqual("BindOutputVariables", append(definition.ProvisionSettings.Outputs, definition.BindSettings.Outputs...), service.BindOutputVariables)
	})

	t.Run("examples", func(t *testing.T) {
		expectEqual("Examples", definition.Examples, service.Examples)
	})

	t.Run("provider-builder", func(t *testing.T) {
		if service.ProviderBuilder == nil {
			t.Fatal("Expected provider builder to not be nil")
		}
	})
}
