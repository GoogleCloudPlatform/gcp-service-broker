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

package firestore

import (
	"code.cloudfoundry.org/lager"
	accountmanagers "github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/account_managers"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/broker_base"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/broker"
	"golang.org/x/oauth2/jwt"
)

func init() {
	// NOTE(jlewisiii) Firestore has some intentional differences from other services.
	// First, it doesn't require legacy compatibility so we won't allow operators to override the whitelist.
	// Second, Firestore uses the old datastore IAM role model so the roles will look strange.
	bs := &broker.ServiceDefinition{
		Name: "google-firestore",
		DefaultServiceDefinition: `{
      "id": "a2b7b873-1e34-4530-8a42-902ff7d66b43",
      "description": "Cloud Firestore is a fast, fully managed, serverless, cloud-native NoSQL document database that simplifies storing, syncing, and querying data for your mobile, web, and IoT apps at global scale.",
      "name": "google-firestore",
      "bindable": true,
      "plan_updateable": false,
      "metadata": {
        "displayName": "Google Cloud Firestore",
        "longDescription": "Cloud Firestore is a fast, fully managed, serverless, cloud-native NoSQL document database that simplifies storing, syncing, and querying data for your mobile, web, and IoT apps at global scale.",
        "documentationUrl": "https://cloud.google.com/firestore/docs/",
        "supportUrl": "https://cloud.google.com/support/",
        "imageUrl": "https://cloud.google.com/_static/images/cloud/products/logos/svg/firestore.svg"
      },
      "tags": ["gcp", "firestore", "preview", "beta"],
      "plans": [
        {
         "id": "64403af0-4413-4ef3-a813-37f0306ef498",
         "name": "default",
         "display_name": "Default",
         "description": "Firestore default plan.",
         "service_properties": {}
        }
      ]
    }`,
		ProvisionInputVariables: []broker.BrokerVariable{},
		BindInputVariables:      accountmanagers.ServiceAccountWhitelistWithDefault([]string{"datastore.user", "datastore.viewer"}, "datastore.user"),
		BindOutputVariables:     accountmanagers.ServiceAccountBindOutputVariables(),
		BindComputedVariables:   accountmanagers.ServiceAccountBindComputedVariables(),
		Examples: []broker.ServiceExample{
			{
				Name:            "Reader Writer",
				Description:     "Creates a general Firestore user and grants it permission to read and write entities.",
				PlanId:          "64403af0-4413-4ef3-a813-37f0306ef498",
				ProvisionParams: map[string]interface{}{},
				BindParams:      map[string]interface{}{},
			},

			{
				Name:            "Read Only",
				Description:     "Creates a Firestore user that can only view entities.",
				PlanId:          "64403af0-4413-4ef3-a813-37f0306ef498",
				ProvisionParams: map[string]interface{}{},
				BindParams:      map[string]interface{}{"role": "datastore.user"},
			},
		},
		ProviderBuilder: func(projectId string, auth *jwt.Config, logger lager.Logger) broker.ServiceProvider {
			bb := broker_base.NewBrokerBase(projectId, auth, logger)
			return &FirestoreBroker{BrokerBase: bb}
		},
	}

	broker.Register(bs)
}
