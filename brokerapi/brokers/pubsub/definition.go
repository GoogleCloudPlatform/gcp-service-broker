// Copyright the Service Broker Project Authors.
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

import "github.com/GoogleCloudPlatform/gcp-service-broker/pkg/broker"
import accountmanagers "github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/account_managers"

func init() {
	bs := &broker.BrokerService{
		Name: "google-pubsub",
		DefaultServiceDefinition: `{
      "id": "628629e3-79f5-4255-b981-d14c6c7856be",
      "description": "A global service for real-time and reliable messaging and streaming data",
      "name": "google-pubsub",
      "bindable": true,
      "plan_updateable": false,
      "metadata": {
        "displayName": "Google PubSub",
        "longDescription": "A global service for real-time and reliable messaging and streaming data",
        "documentationUrl": "https://cloud.google.com/pubsub/docs/",
        "supportUrl": "https://cloud.google.com/support/",
        "imageUrl": "https://cloud.google.com/_static/images/cloud/products/logos/svg/pubsub.svg"
      },
      "tags": ["gcp", "pubsub"],
      "plans": [
        {
          "id": "622f4da3-8731-492a-af29-66a9146f8333",
          "service_id": "628629e3-79f5-4255-b981-d14c6c7856be",
          "name": "default",
          "display_name": "Default",
          "description": "PubSub Default plan",
          "service_properties": {}
        }
      ]
	  }`,
		ProvisionInputVariables: []broker.BrokerVariable{
			broker.BrokerVariable{
				FieldName: "topic_name",
				Type:      broker.JsonTypeString,
				Details:   "Name of the topic.",
				Default:   "a generated value",
			},
			broker.BrokerVariable{
				Required:  true,
				FieldName: "subscription_name",
				Type:      broker.JsonTypeString,
				Details:   `Name of the subscription.`,
			},
			broker.BrokerVariable{
				FieldName: "is_push",
				Type:      broker.JsonTypeBoolean,
				Details:   `Are events handled by POSTing to a URL?`,
				Default:   false,
			},
			broker.BrokerVariable{
				FieldName: "endpoint",
				Type:      broker.JsonTypeString,
				Details:   "If `is_push` == 'true', then this is the URL that will be pused to.",
				Default:   "",
			},
			broker.BrokerVariable{
				FieldName: "ack_deadline",
				Type:      broker.JsonTypeInteger,
				Details: `Value is in seconds. Max: 600
This is the maximum time after a subscriber receives a message
before the subscriber should acknowledge the message. After message
delivery but before the ack deadline expires and before the message is
acknowledged, it is an outstanding message and will not be delivered
again during that time (on a best-effort basis).
        `,
				Default: 10,
			},
		},
		BindInputVariables: accountmanagers.ServiceAccountBindInputVariables(),
		BindOutputVariables: append(accountmanagers.ServiceAccountBindOutputVariables(),
			broker.BrokerVariable{
				FieldName: "subscription_name",
				Type:      broker.JsonTypeString,
				Details:   "Name of the subscription",
			},
			broker.BrokerVariable{
				FieldName: "topic_name",
				Type:      broker.JsonTypeString,
				Details:   "Name of the topic",
			},
		),

		Examples: []broker.ServiceExample{
			broker.ServiceExample{
				Name:        "Basic Configuration",
				Description: "Create a topic and a publisher to it",
				PlanId:      "622f4da3-8731-492a-af29-66a9146f8333",
				ProvisionParams: map[string]interface{}{
					"topic_name":        "example_topic",
					"subscription_name": "example_topic_subscription",
				},
				BindParams: map[string]interface{}{
					"role": "pubsub.publisher",
				},
			},

			broker.ServiceExample{
				Name:        "Calling a Webhook",
				Description: "Call a webhook with the results and increase timeout for latency.",
				PlanId:      "622f4da3-8731-492a-af29-66a9146f8333",
				ProvisionParams: map[string]interface{}{
					"topic_name":        "pusher",
					"subscription_name": "pusher-subscription",
					"is_push":           true,
					"endpoint":          "https://web.hook/destination",
					"ack_deadline":      120,
				},
				BindParams: map[string]interface{}{
					"role": "pubsub.publisher",
				},
			},
		},
	}

	broker.Register(bs)
}
