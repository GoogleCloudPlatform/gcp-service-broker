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

import "encoding/json"

// ModuleInstance represents the configuration of a single instance of a module.
type ModuleInstance struct {
	ModuleName    string
	InstanceName  string
	Configuration map[string]interface{}
}

// MarshalDefinition converts the module instance definition into a JSON
// definition that can be fed to Terraform to be created/destroyed.
func (instance *ModuleInstance) MarshalDefinition() (json.RawMessage, error) {
	instanceConfig := make(map[string]interface{})
	for k, v := range instance.Configuration {
		instanceConfig[k] = v
	}

	instanceConfig["source"] = instance.ModuleName

	defn := map[string]interface{}{
		"module": map[string]interface{}{
			instance.InstanceName: instanceConfig,
		},
	}

	return json.Marshal(defn)
}
