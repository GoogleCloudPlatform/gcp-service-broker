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

package broker

// ServiceExample holds example configurations for a service that _should_
// work.
type ServiceExample struct {
	// Name is a human-readable name of the example.
	Name string `json:"name" yaml:"name" validate:"required"`
	// Description is a long-form description of what this example is about.
	Description string `json:"description" yaml:"description" validate:"required"`
	// PlanId is the plan this example will run against.
	PlanId string `json:"plan_id" yaml:"plan_id" validate:"required"`

	// ProvisionParams is the JSON object that will be passed to provision.
	ProvisionParams map[string]interface{} `json:"provision_params" yaml:"provision_params" validate:"required"`

	// BindParams is the JSON object that will be passed to bind. If nil,
	// this example DOES NOT include a bind portion.
	BindParams map[string]interface{} `json:"bind_params" yaml:"bind_params"`
}
