// Copyright 2018 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
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

// StackdriverProfilerServiceDefinition creates a new ServiceDefinition object
// for the Stackdriver Profiler service.
func StackdriverProfilerServiceDefinition() *broker.ServiceDefinition {
	return &broker.ServiceDefinition{
		Id:               "00b9ca4a-7cd6-406a-a5b7-2f43f41ade75",
		Name:             "google-stackdriver-profiler",
		Description:      "Stackdriver Profiler",
		DisplayName:      "Stackdriver Profiler",
		ImageUrl:         "https://cloud.google.com/_static/images/cloud/products/logos/svg/stackdriver.svg",
		DocumentationUrl: "https://cloud.google.com/profiler/docs/",
		SupportUrl:       "https://cloud.google.com/stackdriver/docs/getting-support",
		Tags:             []string{"gcp", "stackdriver", "profiler"},
		Bindable:         true,
		PlanUpdateable:   false,
		Plans: []broker.ServicePlan{
			{
				ServicePlan: brokerapi.ServicePlan{
					ID:          "594627f6-35f5-462f-9074-10fb033fb18a",
					Name:        "default",
					Description: "Stackdriver Profiler default plan.",
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
		BindComputedVariables:   accountmanagers.FixedRoleBindComputedVariables("cloudprofiler.agent"),
		BindOutputVariables:     accountmanagers.ServiceAccountBindOutputVariables(),
		Examples: []broker.ServiceExample{
			{
				Name:            "Basic Configuration",
				Description:     "Creates an account with the permission `cloudprofiler.agent`.",
				PlanId:          "594627f6-35f5-462f-9074-10fb033fb18a",
				ProvisionParams: map[string]interface{}{},
				BindParams:      map[string]interface{}{},
			},
		},
		ProviderBuilder: NewStackdriverAccountProvider,
		IsBuiltin:       true,
	}
}
