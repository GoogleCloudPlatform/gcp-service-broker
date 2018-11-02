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

package datastore

import (
	"code.cloudfoundry.org/lager"
	accountmanagers "github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/account_managers"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/broker_base"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/broker"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/validation"
	"golang.org/x/oauth2/jwt"
)

func init() {
	bs := &broker.ServiceDefinition{
		Name: "google-datastore",
		DefaultServiceDefinition: `{
      "id": "76d4abb2-fee7-4c8f-aee1-bcea2837f02b",
      "description": "Google Cloud Datastore is a NoSQL document database service.",
      "name": "google-datastore",
      "bindable": true,
      "plan_updateable": false,
      "metadata": {
        "displayName": "Google Cloud Datastore",
        "longDescription": "Google Cloud Datastore is a NoSQL document database built for automatic scaling, high performance, and ease of application development.",
        "documentationUrl": "https://cloud.google.com/datastore/docs/",
        "supportUrl": "https://cloud.google.com/support/",
        "imageUrl": "https://cloud.google.com/_static/images/cloud/products/logos/svg/datastore.svg"
      },
      "tags": ["gcp", "datastore"],
      "plans": [
        {
         "id": "05f1fb6b-b5f0-48a2-9c2b-a5f236507a97",
         "service_id": "76d4abb2-fee7-4c8f-aee1-bcea2837f02b",
         "name": "default",
         "display_name": "Default",
         "description": "Datastore default plan.",
         "service_properties": {}
        }
      ]
    }`,
		ProvisionInputVariables: []broker.BrokerVariable{
			{
				FieldName: "namespace",
				Type:      broker.JsonTypeString,
				Details:   "A context for the identifiers in your entity’s dataset. This ensures that different systems can all interpret an entity's data the same way, based on the rules for the entity’s particular namespace. Blank means the default namespace will be used.",
				Default:   "",
				Constraints: validation.NewConstraintBuilder().
					MaxLength(100).
					Pattern("^[A-Za-z0-9_-]*$").
					Build(),
			},
		},
		BindInputVariables:    []broker.BrokerVariable{},
		BindComputedVariables: accountmanagers.FixedRoleBindComputedVariables("datastore.user"),
		BindOutputVariables: append(accountmanagers.ServiceAccountBindOutputVariables(),
			broker.BrokerVariable{
				FieldName: "namespace",
				Type:      broker.JsonTypeString,
				Details:   "A context for the identifiers in your entity’s dataset.",
				Required:  false,
				Constraints: validation.NewConstraintBuilder().
					MaxLength(100).
					Pattern("^[A-Za-z0-9_-]*$").
					Build(),
			},
		),
		Examples: []broker.ServiceExample{
			{
				Name:            "Basic Configuration",
				Description:     "Creates a datastore and a user with the permission `datastore.user`.",
				PlanId:          "05f1fb6b-b5f0-48a2-9c2b-a5f236507a97",
				ProvisionParams: map[string]interface{}{},
				BindParams:      map[string]interface{}{},
			},
			{
				Name:            "Custom Namespace",
				Description:     "Creates a datastore and returns the provided namespace along with bind calls.",
				PlanId:          "05f1fb6b-b5f0-48a2-9c2b-a5f236507a97",
				ProvisionParams: map[string]interface{}{"namespace": "my-namespace"},
				BindParams:      map[string]interface{}{},
			},
		},
		ProviderBuilder: func(projectId string, auth *jwt.Config, logger lager.Logger) broker.ServiceProvider {
			bb := broker_base.NewBrokerBase(projectId, auth, logger)
			return &DatastoreBroker{BrokerBase: bb}
		},
	}

	broker.Register(bs)
}
