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

package dialogflow

import (
	"code.cloudfoundry.org/lager"
	accountmanagers "github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/account_managers"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/broker_base"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/broker"
	"golang.org/x/oauth2/jwt"
)

// ServiceDefinition creates a new ServiceDefinition object for the Dialogflow service.
func ServiceDefinition() *broker.ServiceDefinition {
	return &broker.ServiceDefinition{
		Name: "google-dialogflow",
		DefaultServiceDefinition: `{
      "id": "e84b69db-3de9-4688-8f5c-26b9d5b1f129",
      "description": "Dialogflow is an end-to-end, build-once deploy-everywhere development suite for creating conversational interfaces for websites, mobile applications, popular messaging platforms, and IoT devices.",
      "name": "google-dialogflow",
      "bindable": true,
      "plan_updateable": false,
      "metadata": {
        "displayName": "Google Cloud Dialogflow",
        "longDescription": "Dialogflow is an end-to-end, build-once deploy-everywhere development suite for creating conversational interfaces for websites, mobile applications, popular messaging platforms, and IoT devices.",
        "documentationUrl": "https://cloud.google.com/dialogflow-enterprise/docs/",
        "supportUrl": "https://cloud.google.com/dialogflow-enterprise/docs/support",
        "imageUrl": "https://cloud.google.com/_static/images/cloud/products/logos/svg/dialogflow-enterprise.svg"
      },
      "tags": ["gcp", "dialogflow", "preview"],
      "plans": [
        {
         "id": "3ac4e1bd-b22d-4a99-864b-d3a3ac582348",
         "name": "default",
         "display_name": "Default",
         "description": "Dialogflow default plan.",
         "service_properties": {}
        }
      ]
    }`,
		ProvisionInputVariables: []broker.BrokerVariable{},
		BindInputVariables:      []broker.BrokerVariable{},
		BindComputedVariables:   accountmanagers.FixedRoleBindComputedVariables("dialogflow.client"),
		BindOutputVariables:     accountmanagers.ServiceAccountBindOutputVariables(),
		Examples: []broker.ServiceExample{
			{
				Name:            "Reader",
				Description:     "Creates a Dialogflow user and grants it permission to detect intent and read/write session properties (contexts, session entity types, etc.).",
				PlanId:          "3ac4e1bd-b22d-4a99-864b-d3a3ac582348",
				ProvisionParams: map[string]interface{}{},
				BindParams:      map[string]interface{}{},
			},
		},
		ProviderBuilder: func(projectId string, auth *jwt.Config, logger lager.Logger) broker.ServiceProvider {
			bb := broker_base.NewBrokerBase(projectId, auth, logger)
			return &DialogflowBroker{BrokerBase: bb}
		},
	}
}
