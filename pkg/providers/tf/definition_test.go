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
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/broker"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/varcontext"
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
			ErrContains: "The Terraform template requires the fields [not_defined] which are missing from the declared inputs.",
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
			ErrContains: "MUST match the service declared outputs",
		},

		"missing template outputs": {
			Action: TfServiceDefinitionV1Action{
				Template: `
        `,
				Outputs: []broker.BrokerVariable{{FieldName: "bucket_name"}},
			},
			ErrContains: "MUST match the service declared outputs",
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
