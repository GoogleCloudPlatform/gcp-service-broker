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

package pubsub

import (
	"code.cloudfoundry.org/lager"
	accountmanagers "github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/account_managers"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/broker_base"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/broker"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/validation"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/varcontext"
	"golang.org/x/oauth2/jwt"
)

// ServiceDefinition creates a new ServiceDefinition object for the Pub/Sub service.
func ServiceDefinition() *broker.ServiceDefinition {
	roleWhitelist := []string{
		"pubsub.publisher",
		"pubsub.subscriber",
		"pubsub.viewer",
		"pubsub.editor",
	}

	return &broker.ServiceDefinition{
		Name: "google-pubsub",
		DefaultServiceDefinition: `{
      "id": "628629e3-79f5-4255-b981-d14c6c7856be",
      "description": "A global service for real-time and reliable messaging and streaming data.",
      "name": "google-pubsub",
      "bindable": true,
      "plan_updateable": false,
      "metadata": {
        "displayName": "Google PubSub",
        "longDescription": "A global service for real-time and reliable messaging and streaming data.",
        "documentationUrl": "https://cloud.google.com/pubsub/docs/",
        "supportUrl": "https://cloud.google.com/pubsub/docs/support",
        "imageUrl": "https://cloud.google.com/_static/images/cloud/products/logos/svg/pubsub.svg"
      },
      "tags": ["gcp", "pubsub"],
      "plans": [
        {
          "id": "622f4da3-8731-492a-af29-66a9146f8333",
          "service_id": "628629e3-79f5-4255-b981-d14c6c7856be",
          "name": "default",
          "display_name": "Default",
          "description": "PubSub Default plan.",
          "service_properties": {},
          "free": false
        }
      ]
	  }`,
		ProvisionInputVariables: []broker.BrokerVariable{
			{
				FieldName: "topic_name",
				Type:      broker.JsonTypeString,
				Details:   `Name of the topic. Must not start with "goog".`,
				Default:   "pcf_sb_${counter.next()}_${time.nano()}",
				Constraints: validation.NewConstraintBuilder().
					MinLength(3).
					MaxLength(255).
					Pattern(`^[a-zA-Z][a-zA-Z0-9\d\-_~%\.\+]+$`). // adapted from the Pub/Sub create topic page's validator
					Build(),
			},
			{
				FieldName: "subscription_name",
				Type:      broker.JsonTypeString,
				Details:   `Name of the subscription. Blank means no subscription will be created. Must not start with "goog".`,
				Default:   "",
				Constraints: validation.NewConstraintBuilder().
					MinLength(3).
					MaxLength(255).
					Pattern(`^[a-zA-Z][a-zA-Z0-9\d\-_~%\.\+]+`). // adapted from the Pub/Sub create subscription page's validator
					Build(),
			},
			{
				FieldName: "is_push",
				Type:      broker.JsonTypeString,
				Details:   `Are events handled by POSTing to a URL?`,
				Default:   "false",
				Enum: map[interface{}]string{
					"true":  "The subscription will POST the events to a URL.",
					"false": "Events will be pulled from the subscription.",
				},
			},
			{
				FieldName: "endpoint",
				Type:      broker.JsonTypeString,
				Details:   "If `is_push` == 'true', then this is the URL that will be pushed to.",
				Default:   "",
			},
			{
				FieldName: "ack_deadline",
				Type:      broker.JsonTypeString,
				Details: `Value is in seconds. Max: 600
This is the maximum time after a subscriber receives a message
before the subscriber should acknowledge the message. After message
delivery but before the ack deadline expires and before the message is
acknowledged, it is an outstanding message and will not be delivered
again during that time (on a best-effort basis).
        `,
				Default: "10",
			},
		},
		ProvisionComputedVariables: []varcontext.DefaultVariable{
			{Name: "labels", Default: "${json.marshal(request.default_labels)}", Overwrite: true},
		},
		DefaultRoleWhitelist: roleWhitelist,
		BindInputVariables:   accountmanagers.ServiceAccountWhitelistWithDefault(roleWhitelist, "pubsub.editor"),
		BindOutputVariables: append(accountmanagers.ServiceAccountBindOutputVariables(),
			broker.BrokerVariable{
				FieldName: "subscription_name",
				Type:      broker.JsonTypeString,
				Details:   "Name of the subscription.",
				Required:  false,
				Constraints: validation.NewConstraintBuilder().
					MinLength(0). // subscription name could be blank on return
					MaxLength(255).
					Pattern(`^(|[a-zA-Z][a-zA-Z0-9\d\-_~%\.\+]+)`). // adapted from the Pub/Sub create subscription page's validator
					Build(),
			},
			broker.BrokerVariable{
				FieldName: "topic_name",
				Type:      broker.JsonTypeString,
				Details:   "Name of the topic.",
				Required:  true,
				Constraints: validation.NewConstraintBuilder().
					MinLength(3).
					MaxLength(255).
					Pattern(`^[a-zA-Z][a-zA-Z0-9\d\-_~%\.\+]+$`). // adapted from the Pub/Sub create topic page's validator
					Build(),
			},
		),
		BindComputedVariables: accountmanagers.ServiceAccountBindComputedVariables(),
		Examples: []broker.ServiceExample{
			{
				Name:        "Basic Configuration",
				Description: "Create a topic and a publisher to it.",
				PlanId:      "622f4da3-8731-492a-af29-66a9146f8333",
				ProvisionParams: map[string]interface{}{
					"topic_name":        "example_topic",
					"subscription_name": "example_topic_subscription",
				},
				BindParams: map[string]interface{}{
					"role": "pubsub.publisher",
				},
			},
			{
				Name:        "No Subscription",
				Description: "Create a topic without a subscription.",
				PlanId:      "622f4da3-8731-492a-af29-66a9146f8333",
				ProvisionParams: map[string]interface{}{
					"topic_name": "example_topic",
				},
				BindParams: map[string]interface{}{
					"role": "pubsub.publisher",
				},
			},
			{
				Name:        "Custom Timeout",
				Description: "Create a subscription with a custom deadline for long processess.",
				PlanId:      "622f4da3-8731-492a-af29-66a9146f8333",
				ProvisionParams: map[string]interface{}{
					"topic_name":        "long_deadline_topic",
					"subscription_name": "long_deadline_subscription",
					"ack_deadline":      "200",
				},
				BindParams: map[string]interface{}{
					"role": "pubsub.publisher",
				},
			},
		},
		ProviderBuilder: func(projectId string, auth *jwt.Config, logger lager.Logger) broker.ServiceProvider {
			bb := broker_base.NewBrokerBase(projectId, auth, logger)
			return &PubSubBroker{BrokerBase: bb}
		},
		IsBuiltin: true,
	}
}
