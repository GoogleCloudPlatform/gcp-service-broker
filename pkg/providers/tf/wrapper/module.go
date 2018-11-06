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
	"sort"

	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/validation"
	"github.com/hashicorp/hcl"
)

// TerraformModule represents a module in a Terraform workspace.
type TerraformModule struct {
	Name       string `validate:"terraform_identifier,required"`
	Definition string `validate:"hcl"`
}

// Validate checks the validity of the TerraformModule struct.
func (module *TerraformModule) Validate() error {
	return validation.ValidateStruct(module)
}

// Inputs gets the input parameter names for the module.
func (module *TerraformModule) Inputs() ([]string, error) {
	defn := terraformModuleHcl{}
	if err := hcl.Decode(&defn, module.Definition); err != nil {
		return nil, err
	}

	return sortedKeys(defn.Inputs), nil
}

// Outputs gets the output parameter names for the module.
func (module *TerraformModule) Outputs() ([]string, error) {
	defn := terraformModuleHcl{}
	if err := hcl.Decode(&defn, module.Definition); err != nil {
		return nil, err
	}

	return sortedKeys(defn.Outputs), nil
}

func sortedKeys(m map[string]interface{}) []string {
	var keys []string
	for key, _ := range m {
		keys = append(keys, key)
	}

	sort.Slice(keys, func(i int, j int) bool { return keys[i] < keys[j] })
	return keys
}

// terraformModuleHcl is a struct used for marshaling/unmarshaling details about
// Terraform modules.
//
// See https://www.terraform.io/docs/modules/create.html for their structure.
type terraformModuleHcl struct {
	Inputs  map[string]interface{} `hcl:"variable"`
	Outputs map[string]interface{} `hcl:"output"`
}
