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

package wrapper

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/validation"
	"github.com/hashicorp/hcl2/gohcl"
	"github.com/hashicorp/hcl2/hcl"
	"github.com/hashicorp/hcl2/hclparse"
)

// ModuleDefinition represents a module in a Terraform workspace.
type ModuleDefinition struct {
	Name       string
	Definition string
}

var _ (validation.Validatable) = (*ModuleDefinition)(nil)

// Validate checks the validity of the ModuleDefinition struct.
func (module *ModuleDefinition) Validate() (errs *validation.FieldError) {
	return errs.Also(
		validation.ErrIfBlank(module.Name, "Name"),
		validation.ErrIfNotTerraformIdentifier(module.Name, "Name"),
		validation.ErrIfNotHCL(module.Definition, "Definition"),
	)
}

// Inputs gets the input parameter names for the module.
func (module *ModuleDefinition) Inputs() ([]string, error) {
	defn := terraformModuleHcl{}

	hclFile, err := parseHCL(module.Definition, module.Name)
	if err != nil {
		return nil, err
	}

	diag := gohcl.DecodeBody(hclFile.Body, nil, &defn)
	if diag.HasErrors() {
		return nil, fmt.Errorf("Error decoding hclFile body: %v", diag.Error())
	}

	return sortedVariableNames(defn.Inputs), nil
}

// Outputs gets the output parameter names for the module.
func (module *ModuleDefinition) Outputs() ([]string, error) {
	defn := terraformModuleHcl{}

	hclFile, err := parseHCL(module.Definition, module.Name)
	if err != nil {
		return nil, err
	}

	diag := gohcl.DecodeBody(hclFile.Body, nil, &defn)
	if diag.HasErrors() {
		return nil, fmt.Errorf("Error decoding hclFile body: %v", diag.Error())
	}

	return sortedVariableNames(defn.Outputs), nil
}

func sortedVariableNames(variables []terraformVariableHcl) []string {
	var vars []string
	for _, variable := range variables {
		vars = append(vars, variable.Name)
	}

	// Sort variable names alphabetically
	sort.Strings(vars)
	return vars
}

func parseHCL(value string, field string) (*hcl.File, error) {
	parser := hclparse.NewParser()
	hclFile := &hcl.File{}

	var diag hcl.Diagnostics
	// Check if value is JSON
	var js json.RawMessage
	if json.Unmarshal([]byte(value), &js) == nil {
		// Try to parse JSON terraform syntax
		hclFile, diag = parser.ParseJSON([]byte(value), field)
	} else {
		// Try to parse HCL syntax
		hclFile, diag = parser.ParseHCL([]byte(value), field)
	}

	if diag.HasErrors() {
		return hclFile, fmt.Errorf("Error parsing hcl file: %v", diag.Error())
	}

	return hclFile, nil
}

// terraformModuleHcl is a struct used for marshaling/unmarshaling details about
// Terraform modules.
//
// See https://www.terraform.io/docs/modules/create.html for their structure.
type terraformModuleHcl struct {
	Inputs  []terraformVariableHcl `hcl:"variable,block"`
	Outputs []terraformVariableHcl `hcl:"output,block"`
	Remain  hcl.Body               `hcl:",remain"`
}

type terraformVariableHcl struct {
	Name   string   `hcl:"name,label"`
	Remain hcl.Body `hcl:",remain"`
}
