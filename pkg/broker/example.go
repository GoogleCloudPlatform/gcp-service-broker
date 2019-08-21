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

import "github.com/GoogleCloudPlatform/gcp-service-broker/pkg/validation"

// ServiceExample holds example configurations for a service that _should_
// work.
type ServiceExample struct {
	// Name is a human-readable name of the example.
	Name string `json:"name" yaml:"name"`
	// Description is a long-form description of what this example is about.
	Description string `json:"description" yaml:"description"`
	// PlanId is the plan this example will run against.
	PlanId string `json:"plan_id" yaml:"plan_id"`

	// ProvisionParams is the JSON object that will be passed to provision.
	ProvisionParams map[string]interface{} `json:"provision_params" yaml:"provision_params"`

	// BindParams is the JSON object that will be passed to bind. If nil,
	// this example DOES NOT include a bind portion.
	BindParams map[string]interface{} `json:"bind_params" yaml:"bind_params"`
}

var _ validation.Validatable = (*ServiceExample)(nil)

// Validate implements validation.Validatable.
func (action *ServiceExample) Validate() (errs *validation.FieldError) {
	return errs.Also(
		validation.ErrIfBlank(action.Name, "name"),
		validation.ErrIfBlank(action.Description, "description"),
		validation.ErrIfBlank(action.PlanId, "plan_id"),
	)
}
