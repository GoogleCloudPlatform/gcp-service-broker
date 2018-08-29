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
	"encoding/json"

	googlepubsub "cloud.google.com/go/pubsub"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/models"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/name_generator"
	"github.com/pivotal-cf/brokerapi"
	"golang.org/x/net/context"

	"fmt"
	"strconv"
	"time"

	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/broker_base"
	"google.golang.org/api/option"
)

// PubSubBroker is the service-broker back-end for creating Google Pub/Sub
// topics, subscriptions, and accounts
type PubSubBroker struct {
	broker_base.BrokerBase
}

// InstanceInformation holds the details needed to connect to a PubSub instance
// after it has been provisioned
type InstanceInformation struct {
	TopicName        string `json:"topic_name"`
	SubscriptionName string `json:"subscription_name"`
}

// Provision creates a new PubSub topic with the name given in details.topic_name
// if subscription_name is supplied, will also create a subscription for this topic with optional config parameters
// is_push (defaults to "false"; i.e. pull), endpoint (defaults to nil), ack_deadline (seconds, defaults to 10, 600 max)
func (b *PubSubBroker) Provision(instanceId string, details brokerapi.ProvisionDetails, plan models.ServicePlan) (models.ServiceInstanceDetails, error) {

	var err error
	var params map[string]string
	if len(details.RawParameters) == 0 {
		params = map[string]string{}
	} else if err := json.Unmarshal(details.RawParameters, &params); err != nil {
		return models.ServiceInstanceDetails{}, fmt.Errorf("Error unmarshalling provision details: %s", err)
	}

	// Ensure there is a name for this topic
	if _, ok := params["topic_name"]; !ok {
		params["topic_name"] = name_generator.Basic.InstanceName()
	}

	ctx := context.Background()
	co := option.WithUserAgent(models.CustomUserAgent)
	ct := option.WithTokenSource(b.HttpConfig.TokenSource(context.Background()))
	pubsubClient, err := googlepubsub.NewClient(ctx, b.ProjectId, co, ct)

	if err != nil {
		return models.ServiceInstanceDetails{}, fmt.Errorf("Error creating new pubsub client: %s", err)
	}

	t, err := pubsubClient.CreateTopic(ctx, params["topic_name"])
	if err != nil {
		return models.ServiceInstanceDetails{}, fmt.Errorf("Error creating new pubsub topic: %s", err)
	}

	i := models.ServiceInstanceDetails{
		Name:         params["topic_name"],
		Url:          "",
		Location:     "",
		OtherDetails: "{}",
	}
	ii := InstanceInformation{
		TopicName: params["topic_name"],
	}

	if sub_name, ok := params["subscription_name"]; ok {
		var pushConfig googlepubsub.PushConfig
		var ackDeadline = 10

		if ackd, ok := params["ack_deadline"]; ok {
			ackDeadline, err = strconv.Atoi(ackd)
			if err != nil {
				return models.ServiceInstanceDetails{}, fmt.Errorf("Error converting ack deadline to int: %s", err)
			}
		}

		if isPush, ok := params["is_push"]; ok {
			if isPush == "true" {
				pushConfig = googlepubsub.PushConfig{
					Endpoint: params["endpoint"],
				}
			}
		}

		subsConfig := googlepubsub.SubscriptionConfig{
			PushConfig:  pushConfig,
			Topic:       t,
			AckDeadline: time.Duration(ackDeadline) * time.Second,
		}

		_, err = pubsubClient.CreateSubscription(ctx, sub_name, subsConfig)
		if err != nil {
			return models.ServiceInstanceDetails{}, fmt.Errorf("Error creating subscription: %s", err)
		}

		ii.SubscriptionName = params["subscription_name"]
	}

	otherDetails, err := json.Marshal(ii)
	if err != nil {
		return models.ServiceInstanceDetails{}, fmt.Errorf("Error marshalling json: %s", err)
	}
	i.OtherDetails = string(otherDetails)

	return i, nil
}

// Deprovision deletes the topic associated with the given instanceID
func (b *PubSubBroker) Deprovision(ctx context.Context, topic models.ServiceInstanceDetails, details brokerapi.DeprovisionDetails) error {
	ct := option.WithTokenSource(b.HttpConfig.TokenSource(ctx))
	service, err := googlepubsub.NewClient(ctx, b.ProjectId, ct)
	if err != nil {
		return fmt.Errorf("Error creating new pubsub client: %s", err)
	}

	err = service.Topic(topic.Name).Delete(ctx)
	if err != nil {
		return fmt.Errorf("Error deleting pubsub topic: %s", err)
	}

	otherD := make(map[string]string)
	err = json.Unmarshal([]byte(topic.OtherDetails), &otherD)
	if err != nil {
		return fmt.Errorf("Error unmarshalling service instance other details: %s", err)
	}

	if subscriptionName := otherD["subscription_name"]; subscriptionName != "" {
		err = service.Subscription(subscriptionName).Delete(ctx)
		if err != nil {
			return fmt.Errorf("Error deleting subscription: %s", err)
		}
	}

	return nil
}
