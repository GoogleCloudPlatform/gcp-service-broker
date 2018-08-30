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
	accountmanagers "github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/account_managers"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/broker"
)

func init() {
	roleWhitelist := []string{
		"ml.developer",
		"ml.viewer",
		"ml.modelOwner",
		"ml.modelUser",
		"ml.jobOwner",
		"ml.operationOwner",
	}

	bs := &broker.BrokerService{
		Name: "google-ml-apis",
		DefaultServiceDefinition: `
		{
      "id": "5ad2dce0-51f7-4ede-8b46-293d6df1e8d4",
      "description": "Machine Learning Apis including Vision, Translate, Speech, and Natural Language",
      "name": "google-ml-apis",
      "bindable": true,
      "plan_updateable": false,
      "metadata": {
        "displayName": "Google Machine Learning APIs",
        "longDescription": "Machine Learning Apis including Vision, Translate, Speech, and Natural Language",
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
         "description": "Machine Learning api default plan",
         "service_properties": {}
        }
      ]
    }
		`,
		ProvisionInputVariables: []broker.BrokerVariable{},
		DefaultRoleWhitelist:    roleWhitelist,
		BindInputVariables:      accountmanagers.ServiceAccountBindInputVariables(roleWhitelist),
		BindOutputVariables:     accountmanagers.ServiceAccountBindOutputVariables(),
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
	}

	broker.Register(bs)
}
