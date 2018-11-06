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
)

func init() {
	bs := &broker.ServiceDefinition{
		Name: "google-stackdriver-debugger",
		DefaultServiceDefinition: `{
		      "id": "83837945-1547-41e0-b661-ea31d76eed11",
		      "description": "Stackdriver Debugger",
		      "name": "google-stackdriver-debugger",
		      "bindable": true,
		      "plan_updateable": false,
		      "metadata": {
		        "displayName": "Stackdriver Debugger",
		        "longDescription": "Stackdriver Debugger is a feature of the Google Cloud Platform that lets you inspect the state of an application at any code location without using logging statements and without stopping or slowing down your applications. Your users are not impacted during debugging. Using the production debugger you can capture the local variables and call stack and link it back to a specific line location in your source code.",
		        "documentationUrl": "https://cloud.google.com/debugger/docs/",
		        "supportUrl": "https://cloud.google.com/stackdriver/docs/getting-support",
		        "imageUrl": "https://cloud.google.com/_static/images/cloud/products/logos/svg/debugger.svg"
		      },
		      "tags": ["gcp", "stackdriver", "debugger"],
		      "plans": [
		        {
		          "id": "10866183-a775-49e8-96e3-4e7a901e4a79",
		          "service_id": "83837945-1547-41e0-b661-ea31d76eed11",
		          "name": "default",
		          "display_name": "Default",
		          "description": "Stackdriver Debugger default plan.",
		          "service_properties": {}
		        }
		      ]
				}
		`,
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
	}

	broker.Register(bs)
}
