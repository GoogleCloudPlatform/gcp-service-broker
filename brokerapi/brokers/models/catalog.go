// Copyright the Service Broker Project Authors.
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

package models

import pcfosb "github.com/pivotal-cf/brokerapi"

// Service overrides the canonical Service Broker service type using a custom
// type for Plans, everything else is the same.
type Service struct {
	pcfosb.Service

	Plans []ServicePlan `json:"plans"`
}

// Converts this service to a plain PCF Service definition.
func (s Service) ToPlain() pcfosb.Service {
	plain := s.Service
	plainPlans := []pcfosb.ServicePlan{}

	for _, plan := range s.Plans {
		plainPlans = append(plainPlans, plan.ServicePlan)
	}

	plain.Plans = plainPlans

	return plain
}

// ServicePlan extends the OSB ServicePlan by including a map of key/value
// pairs that can be used to pass additional information to the back-end.
type ServicePlan struct {
	pcfosb.ServicePlan

	ServiceProperties map[string]string `json:"service_properties"`
}
