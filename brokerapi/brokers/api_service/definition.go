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

package api_service

import (
	"code.cloudfoundry.org/lager"
	accountmanagers "github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/account_managers"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/broker_base"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/broker"
	"golang.org/x/oauth2/jwt"
)

// ServiceDefinition creates a new ServiceDefinition object for the ML service.
func ServiceDefinition() *broker.ServiceDefinition {
	roleWhitelist := []string{
		"ml.developer",
		"ml.viewer",
		"ml.modelOwner",
		"ml.modelUser",
		"ml.jobOwner",
		"ml.operationOwner",
	}

	return &broker.ServiceDefinition{
		Name: "google-ml-apis",
		DefaultServiceDefinition: `
		{
      "id": "5ad2dce0-51f7-4ede-8b46-293d6df1e8d4",
      "description": "Machine Learning APIs including Vision, Translate, Speech, and Natural Language.",
      "name": "google-ml-apis",
      "bindable": true,
      "plan_updateable": false,
      "metadata": {
        "displayName": "Google Machine Learning APIs",
        "longDescription": "Machine Learning APIs including Vision, Translate, Speech, and Natural Language.",
        "documentationUrl": "https://cloud.google.com/ml/",
        "supportUrl": "https://cloud.google.com/support/",
        "imageUrl": "https://cloud.google.com/_static/images/cloud/products/logos/svg/machine-learning.svg"
      },
      "tags": ["gcp", "ml"],
      "plans":  [
        {
         "id": "be7954e1-ecfb-4936-a0b6-db35e6424c7a",
         "service_id": "5ad2dce0-51f7-4ede-8b46-293d6df1e8d4",
         "name": "default",
         "display_name": "Default",
         "description": "Machine Learning API default plan.",
         "service_properties": {},
         "free": false
        }
      ]
    }
		`,
		ProvisionInputVariables: []broker.BrokerVariable{},
		DefaultRoleWhitelist:    roleWhitelist,
		BindInputVariables:      accountmanagers.ServiceAccountWhitelistWithDefault(roleWhitelist, "ml.modelUser"),
		BindOutputVariables:     accountmanagers.ServiceAccountBindOutputVariables(),
		BindComputedVariables:   accountmanagers.ServiceAccountBindComputedVariables(),
		Examples: []broker.ServiceExample{
			{
				Name:            "Basic Configuration",
				Description:     "Create an account with developer access to your ML models.",
				PlanId:          "be7954e1-ecfb-4936-a0b6-db35e6424c7a",
				ProvisionParams: map[string]interface{}{},
				BindParams: map[string]interface{}{
					"role": "ml.developer",
				},
			},
		},
		ProviderBuilder: func(projectId string, auth *jwt.Config, logger lager.Logger) broker.ServiceProvider {
			bb := broker_base.NewBrokerBase(projectId, auth, logger)
			return &ApiServiceBroker{BrokerBase: bb}
		},
		IsBuiltin: true,
	}
}
