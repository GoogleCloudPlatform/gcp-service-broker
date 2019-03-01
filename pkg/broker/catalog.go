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

import (
	"github.com/pivotal-cf/brokerapi"
)

// Service overrides the canonical Service Broker service type using a custom
// type for Plans, everything else is the same.
type Service struct {
	brokerapi.Service

	Plans []ServicePlan `json:"plans"`
}

// ToPlain converts this service to a plain PCF Service definition.
func (s Service) ToPlain() brokerapi.Service {
	plain := s.Service
	plainPlans := []brokerapi.ServicePlan{}

	for _, plan := range s.Plans {
		plainPlans = append(plainPlans, plan.ServicePlan)
	}

	plain.Plans = plainPlans

	return plain
}

// ServicePlan extends the OSB ServicePlan by including a map of key/value
// pairs that can be used to pass additional information to the back-end.
type ServicePlan struct {
	brokerapi.ServicePlan

	ServiceProperties  map[string]string      `json:"service_properties"`
	ProvisionOverrides map[string]interface{} `json:"provision_overrides,omitempty"`
	BindOverrides      map[string]interface{} `json:"bind_overrides,omitempty"`
}

// GetServiceProperties gets the plan settings variables as a string->interface map.
func (sp *ServicePlan) GetServiceProperties() map[string]interface{} {
	props := make(map[string]interface{})
	for k, v := range sp.ServiceProperties {
		props[k] = v
	}

	return props
}
