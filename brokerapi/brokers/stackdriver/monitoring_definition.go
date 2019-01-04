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

// StackdriverMonitoringServiceDefinition creates a new ServiceDefinition object
// for the Stackdriver Monitoring service.
func StackdriverMonitoringServiceDefinition() *broker.ServiceDefinition {
	return &broker.ServiceDefinition{
		Name: "google-stackdriver-monitoring",
		DefaultServiceDefinition: `{
      "id": "2bc0d9ed-3f68-4056-b842-4a85cfbc727f",
      "description": "Stackdriver Monitoring",
      "name": "google-stackdriver-monitoring",
      "bindable": true,
      "plan_updateable": false,
      "metadata": {
        "displayName": "Stackdriver Monitoring",
        "longDescription": "Stackdriver Monitoring provides visibility into the performance, uptime, and overall health of cloud-powered applications. ",
        "documentationUrl": "https://cloud.google.com/monitoring/docs/",
        "supportUrl": "https://cloud.google.com/stackdriver/docs/getting-support",
        "imageUrl": "https://cloud.google.com/_static/images/cloud/products/logos/svg/stackdriver.svg"
      },
      "tags": ["gcp", "stackdriver", "monitoring", "preview"],
      "plans": [
        {
          "id": "2e4b85c1-0ce6-46e4-91f5-eebeb373e3f5",
          "name": "default",
          "display_name": "Default",
          "description": "Stackdriver Monitoring default plan.",
          "service_properties": {},
          "free": false
        }
      ]
    }
		`,
		ProvisionInputVariables: []broker.BrokerVariable{},
		BindInputVariables:      []broker.BrokerVariable{},
		BindComputedVariables:   accountmanagers.FixedRoleBindComputedVariables("monitoring.metricWriter"),
		BindOutputVariables:     accountmanagers.ServiceAccountBindOutputVariables(),
		Examples: []broker.ServiceExample{
			{
				Name:            "Basic Configuration",
				Description:     "Creates an account with the permission `monitoring.metricWriter` for writing metrics.",
				PlanId:          "2e4b85c1-0ce6-46e4-91f5-eebeb373e3f5",
				ProvisionParams: map[string]interface{}{},
				BindParams:      map[string]interface{}{},
			},
		},
		ProviderBuilder: NewStackdriverAccountProvider,
		IsBuiltin:       true,
	}
}
