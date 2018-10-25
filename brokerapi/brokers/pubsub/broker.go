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
	"errors"

	googlepubsub "cloud.google.com/go/pubsub"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/models"
	"github.com/GoogleCloudPlatform/gcp-service-broker/utils"
	"github.com/pivotal-cf/brokerapi"
	"golang.org/x/net/context"

	"fmt"
	"time"

	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/broker_base"
	"google.golang.org/api/option"
)

// PubSubBroker is the service-broker back-end for creating Google Pub/Sub
// topics, subscriptions, and accounts.
type PubSubBroker struct {
	broker_base.BrokerBase
}

// InstanceInformation holds the details needed to connect to a PubSub instance
// after it has been provisioned.
type InstanceInformation struct {
	TopicName string `json:"topic_name"`

	// SubscriptionName is optional, if non-empty then a susbcription was created.
	SubscriptionName string `json:"subscription_name"`
}

// Provision creates a new Pub/Sub topic from the settings in the user-provided details and service plan.
// If a subscription name is supplied, the function will also create a subscription for the topic.
func (b *PubSubBroker) Provision(ctx context.Context, instanceId string, details brokerapi.ProvisionDetails, plan models.ServicePlan) (models.ServiceInstanceDetails, error) {
	variableContext, err := serviceDefinition().ProvisionVariables(instanceId, details, plan)
	if err != nil {
		return models.ServiceInstanceDetails{}, err
	}

	// Extract and validate all the params exist and are the right types
	defaultLabels := utils.ExtractDefaultLabels(instanceId, details)
	topicName := variableContext.GetString("topic_name")
	subscriptionName := variableContext.GetString("subscription_name")
	endpoint := ""
	if variableContext.GetBool("is_push") {
		endpoint = variableContext.GetString("endpoint")
	}

	subscriptionConfig := googlepubsub.SubscriptionConfig{
		PushConfig: googlepubsub.PushConfig{
			Endpoint: endpoint,
		},
		AckDeadline: time.Duration(variableContext.GetInt("ack_deadline")) * time.Second,
		Labels:      defaultLabels,
	}

	if err := variableContext.Error(); err != nil {
		return models.ServiceInstanceDetails{}, err
	}

	// Check special-cases
	if topicName == "" {
		return models.ServiceInstanceDetails{}, errors.New("topic_name must not be blank")
	}

	// Create
	pubsubClient, err := b.createClient(ctx)
	if err != nil {
		return models.ServiceInstanceDetails{}, err
	}

	topic, err := pubsubClient.CreateTopic(ctx, topicName)
	if err != nil {
		return models.ServiceInstanceDetails{}, fmt.Errorf("Error creating new Pub/Sub topic: %s", err)
	}

	// This service needs labels to be set after creation
	if _, err := topic.Update(ctx, googlepubsub.TopicConfigToUpdate{Labels: defaultLabels}); err != nil {
		return models.ServiceInstanceDetails{}, fmt.Errorf("Error setting labels on new Pub/Sub topic: %s", err)
	}

	if subscriptionName != "" {
		subscriptionConfig.Topic = topic

		if _, err := pubsubClient.CreateSubscription(ctx, subscriptionName, subscriptionConfig); err != nil {
			return models.ServiceInstanceDetails{}, fmt.Errorf("Error creating subscription: %s", err)
		}
	}

	ii := InstanceInformation{
		TopicName:        topicName,
		SubscriptionName: subscriptionName,
	}

	id := models.ServiceInstanceDetails{
		Name:     topicName,
		Url:      "",
		Location: "",
	}

	if err := id.SetOtherDetails(ii); err != nil {
		return models.ServiceInstanceDetails{}, err
	}

	return id, err
}

// Deprovision deletes the topic and subscription associated with the given instance.
func (b *PubSubBroker) Deprovision(ctx context.Context, topic models.ServiceInstanceDetails, details brokerapi.DeprovisionDetails) (*string, error) {
	service, err := b.createClient(ctx)
	if err != nil {
		return nil, err
	}

	if err := service.Topic(topic.Name).Delete(ctx); err != nil {
		return nil, fmt.Errorf("Error deleting pubsub topic: %s", err)
	}

	otherDetails := InstanceInformation{}
	if err := topic.GetOtherDetails(&otherDetails); err != nil {
		return nil, err
	}

	if otherDetails.SubscriptionName != "" {
		if err := service.Subscription(otherDetails.SubscriptionName).Delete(ctx); err != nil {
			return nil, fmt.Errorf("Error deleting subscription: %s", err)
		}
	}

	return nil, nil
}

func (b *PubSubBroker) createClient(ctx context.Context) (*googlepubsub.Client, error) {
	co := option.WithUserAgent(models.CustomUserAgent)
	ct := option.WithTokenSource(b.HttpConfig.TokenSource(ctx))
	client, err := googlepubsub.NewClient(ctx, b.ProjectId, co, ct)
	if err != nil {
		return nil, fmt.Errorf("Couldn't instantiate Pub/Sub API client: %s", err)
	}

	return client, nil
}
