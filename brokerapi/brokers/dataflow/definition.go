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

package dataflow

import (
	"code.cloudfoundry.org/lager"
	accountmanagers "github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/account_managers"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/broker_base"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/broker"
	"golang.org/x/oauth2/jwt"
)

// ServiceDefinition creates a new ServiceDefinition object for the Dataflow service.
func ServiceDefinition() *broker.ServiceDefinition {
	roleWhitelist := []string{"dataflow.viewer", "dataflow.developer"}

	return &broker.ServiceDefinition{
		Name: "google-dataflow",
		DefaultServiceDefinition: `{
      "id": "3e897eb3-9062-4966-bd4f-85bda0f73b3d",
      "description": "A managed service for executing a wide variety of data processing patterns built on Apache Beam.",
      "name": "google-dataflow",
      "bindable": true,
      "plan_updateable": false,
      "metadata": {
        "displayName": "Google Cloud Dataflow",
        "longDescription": "A managed service for executing a wide variety of data processing patterns built on Apache Beam.",
        "documentationUrl": "https://cloud.google.com/dataflow/docs/",
        "supportUrl": "https://cloud.google.com/dataflow/docs/support",
        "imageUrl": "https://cloud.google.com/_static/images/cloud/products/logos/svg/dataflow.svg"
      },
      "tags": ["gcp", "dataflow", "preview"],
      "plans": [
        {
         "id": "8e956dd6-8c0f-470c-9a11-065537d81872",
         "name": "default",
         "display_name": "Default",
         "description": "Dataflow default plan.",
         "service_properties": {},
         "free": false
        }
      ]
    }`,
		ProvisionInputVariables: []broker.BrokerVariable{},
		BindInputVariables:      accountmanagers.ServiceAccountWhitelistWithDefault(roleWhitelist, "dataflow.developer"),
		BindComputedVariables:   accountmanagers.ServiceAccountBindComputedVariables(),
		BindOutputVariables:     accountmanagers.ServiceAccountBindOutputVariables(),
		Examples: []broker.ServiceExample{
			{
				Name:            "Developer",
				Description:     "Creates a Dataflow user and grants it permission to create, drain and cancel jobs.",
				PlanId:          "8e956dd6-8c0f-470c-9a11-065537d81872",
				ProvisionParams: map[string]interface{}{},
				BindParams:      map[string]interface{}{},
			},
			{
				Name:            "Viewer",
				Description:     "Creates a Dataflow user and grants it permission to create, drain and cancel jobs.",
				PlanId:          "8e956dd6-8c0f-470c-9a11-065537d81872",
				ProvisionParams: map[string]interface{}{},
				BindParams:      map[string]interface{}{"role": "dataflow.viewer"},
			},
		},
		ProviderBuilder: func(projectId string, auth *jwt.Config, logger lager.Logger) broker.ServiceProvider {
			bb := broker_base.NewBrokerBase(projectId, auth, logger)
			return &DataflowBroker{BrokerBase: bb}
		},
		IsBuiltin: true,
	}
}
