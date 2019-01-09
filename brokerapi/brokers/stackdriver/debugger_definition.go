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

package stackdriver

import (
	accountmanagers "github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/account_managers"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/broker"
	"github.com/pivotal-cf/brokerapi"
)

// StackdriverDebuggerServiceDefinition creates a new ServiceDefinition object
// for the Stackdriver Debugger service.
func StackdriverDebuggerServiceDefinition() *broker.ServiceDefinition {
	return &broker.ServiceDefinition{
		Id:               "83837945-1547-41e0-b661-ea31d76eed11",
		Name:             "google-stackdriver-debugger",
		Description:      "Stackdriver Debugger",
		DisplayName:      "Stackdriver Debugger",
		ImageUrl:         "https://cloud.google.com/_static/images/cloud/products/logos/svg/debugger.svg",
		DocumentationUrl: "https://cloud.google.com/debugger/docs/",
		SupportUrl:       "https://cloud.google.com/stackdriver/docs/getting-support",
		Tags:             []string{"gcp", "stackdriver", "debugger"},
		Bindable:         true,
		PlanUpdateable:   false,
		Plans: []broker.ServicePlan{
			{
				ServicePlan: brokerapi.ServicePlan{
					ID:          "10866183-a775-49e8-96e3-4e7a901e4a79",
					Name:        "default",
					Description: "Stackdriver Debugger default plan.",
					Free:        brokerapi.FreeValue(false),
					Metadata: &brokerapi.ServicePlanMetadata{
						DisplayName: "Default",
					},
				},
				ServiceProperties: map[string]string{},
			},
		},
		ProvisionInputVariables: []broker.BrokerVariable{},
		BindInputVariables:      []broker.BrokerVariable{},
		BindComputedVariables:   accountmanagers.FixedRoleBindComputedVariables("clouddebugger.agent"),
		BindOutputVariables:     accountmanagers.ServiceAccountBindOutputVariables(),
		Examples: []broker.ServiceExample{
			{
				Name:            "Basic Configuration",
				Description:     "Creates an account with the permission `clouddebugger.agent`.",
				PlanId:          "10866183-a775-49e8-96e3-4e7a901e4a79",
				ProvisionParams: map[string]interface{}{},
				BindParams:      map[string]interface{}{},
			},
		},
		ProviderBuilder: NewStackdriverAccountProvider,
		IsBuiltin:       true,
	}
}
