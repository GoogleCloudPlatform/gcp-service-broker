// Copyright 2019 the Service Broker Project Authors.
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
	"encoding/json"
	"fmt"

	"github.com/pivotal-cf/brokerapi"
	"github.com/spf13/viper"
)

// ServiceConfigProperty holds the Viper property for the service config map.
const ServiceConfigProperty = "service-config"

// ServiceConfigMap is a mapping of Service GUID -> ServiceConfig objects.
type ServiceConfigMap map[string]ServiceConfig

// ServiceConfig holds user-defined configuration for a specific service.
type ServiceConfig struct {
	// Notes contains user-defined notes about this config.
	Notes             string                 `json:"//,omitempty"`
	Disabled          bool                   `json:"disabled"`
	ProvisionDefaults map[string]interface{} `json:"provision_defaults"`
	BindDefaults      map[string]interface{} `json:"bind_defaults"`
	CustomPlans       []CustomPlan           `json:"custom_plans"`
}

// CustomPlan holds operator defined variables for each service.
type CustomPlan struct {
	GUID               string                 `json:"guid" validate:"required,uuid"`
	Name               string                 `json:"name" validate:"required"`
	DisplayName        string                 `json:"display_name" validate:"required"`
	Description        string                 `json:"description" validate:"required"`
	Properties         map[string]string      `json:"properties"`
	ProvisionOverrides map[string]interface{} `json:"provision_overrides"`
	BindOverrides      map[string]interface{} `json:"bind_overrides"`
}

// ToServicePlan converts the CustomPlan to a ServicePlan.
func (c *CustomPlan) ToServicePlan() ServicePlan {
	return ServicePlan{
		ServicePlan: brokerapi.ServicePlan{
			Description: c.Description,
			Name:        c.Name,
			ID:          c.GUID,
			Metadata: &brokerapi.ServicePlanMetadata{
				DisplayName: c.DisplayName,
			},
		},
		ServiceProperties:  c.Properties,
		ProvisionOverrides: c.ProvisionOverrides,
		BindOverrides:      c.BindOverrides,
	}
}

// NewServiceConfigMapFromEnv reads viper for the value at ServiceConfigProperty
// and deserializes it into a ServiceConfigMap.
func NewServiceConfigMapFromEnv() (ServiceConfigMap, error) {
	out := ServiceConfigMap{}
	sources := viper.GetString(ServiceConfigProperty)
	if len(sources) == 0 {
		return out, nil
	}

	if err := json.Unmarshal([]byte(sources), &out); err != nil {
		return nil, fmt.Errorf("couldn't deserialize ServiceConfigMap: %v", err)
	}

	return out, nil
}
