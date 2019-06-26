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
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/broker"
	accountmanagers "github.com/GoogleCloudPlatform/gcp-service-broker/pkg/providers/builtin/account_managers"
	"github.com/pivotal-cf/brokerapi"
)

// StackdriverTraceServiceDefinition creates a new ServiceDefinition object
// for the Stackdriver Trace service.
func StackdriverTraceServiceDefinition() *broker.ServiceDefinition {
	return &broker.ServiceDefinition{
		Id:               "c5ddfe15-24d9-47f8-8ffe-f6b7daa9cf4a",
		Name:             "google-stackdriver-trace",
		Description:      "A real-time distributed tracing system.",
		DisplayName:      "Stackdriver Trace",
		ImageUrl:         "https://cloud.google.com/_static/images/cloud/products/logos/svg/trace.svg",
		DocumentationUrl: "https://cloud.google.com/trace/docs/",
		SupportUrl:       "https://cloud.google.com/stackdriver/docs/getting-support",
		Tags:             []string{"gcp", "stackdriver", "trace"},
		Bindable:         true,
		PlanUpdateable:   false,
		Plans: []broker.ServicePlan{
			{
				ServicePlan: brokerapi.ServicePlan{
					ID:          "ab6c2287-b4bc-4ff4-a36a-0575e7910164",
					Name:        "default",
					Description: "Stackdriver Trace default plan.",
					Free:        brokerapi.FreeValue(false),
				},
				ServiceProperties: map[string]string{},
			},
		},
		ProvisionInputVariables: []broker.BrokerVariable{},
		BindInputVariables:      []broker.BrokerVariable{},
		BindComputedVariables:   accountmanagers.FixedRoleBindComputedVariables("cloudtrace.agent"),
		BindOutputVariables:     accountmanagers.ServiceAccountBindOutputVariables(),
		Examples: []broker.ServiceExample{
			{
				Name:            "Basic Configuration",
				Description:     "Creates an account with the permission `cloudtrace.agent`.",
				PlanId:          "ab6c2287-b4bc-4ff4-a36a-0575e7910164",
				ProvisionParams: map[string]interface{}{},
				BindParams:      map[string]interface{}{},
			},
		},
		ProviderBuilder: NewStackdriverAccountProvider,
		IsBuiltin:       true,
	}
}
