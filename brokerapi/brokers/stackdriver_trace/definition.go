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

package stackdriver_trace

import (
	accountmanagers "github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/account_managers"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/broker"
)

func init() {
	bs := &broker.BrokerService{
		Name: "google-stackdriver-trace",
		DefaultServiceDefinition: `{
      "id": "c5ddfe15-24d9-47f8-8ffe-f6b7daa9cf4a",
      "description": "Stackdriver Trace",
      "name": "google-stackdriver-trace",
      "bindable": true,
      "plan_updateable": false,
      "metadata": {
        "displayName": "Stackdriver Trace",
        "longDescription": "Stackdriver Trace is a distributed tracing system that collects latency data from your applications and displays it in the Google Cloud Platform Console. You can track how requests propagate through your application and receive detailed near real-time performance insights.",
        "documentationUrl": "https://cloud.google.com/trace/docs/",
        "supportUrl": "https://cloud.google.com/support/",
        "imageUrl": "https://cloud.google.com/_static/images/cloud/products/logos/svg/trace.svg"
      },
      "tags": ["gcp", "stackdriver", "trace"],
      "plans": [
        {
          "id": "ab6c2287-b4bc-4ff4-a36a-0575e7910164",
          "service_id": "c5ddfe15-24d9-47f8-8ffe-f6b7daa9cf4a",
          "name": "default",
          "display_name": "Default",
          "description": "Stackdriver Trace default plan.",
          "service_properties": {}
        }
      ]
    }
		`,
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
	}

	broker.Register(bs)
}
