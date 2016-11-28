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
//
////////////////////////////////////////////////////////////////////////////////
//

package pubsub

import (
	googlepubsub "cloud.google.com/go/pubsub"
	"encoding/json"
	"gcp-service-broker/brokerapi/brokers/models"
	"gcp-service-broker/brokerapi/brokers/name_generator"
	"golang.org/x/net/context"
	"net/http"

	"code.cloudfoundry.org/lager"
	"fmt"
	"gcp-service-broker/brokerapi/brokers/broker_base"
	"gcp-service-broker/db_service"
	"google.golang.org/api/option"
	"strconv"
	"time"
)

type PubSubBroker struct {
	Client    *http.Client
	ProjectId string
	Logger    lager.Logger

	broker_base.BrokerBase
}

// Creates a new PubSub topic with the name given in details.topic_name
// if subscription_name is supplied, will also create a subscription for this topic with optional config parameters
// is_push (defaults to "false"; i.e. pull), endpoint (defaults to nil), ack_deadline (seconds, defaults to 10, 600 max)
func (b *PubSubBroker) Provision(instanceId string, details models.ProvisionDetails, plan models.PlanDetails) (models.ServiceInstanceDetails, error) {

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
	service, err := googlepubsub.NewClient(ctx, b.ProjectId, co)

	if err != nil {
		return models.ServiceInstanceDetails{}, fmt.Errorf("Error creating new pubsub client: %s", err)
	}

	t, err := service.NewTopic(ctx, params["topic_name"])
	if err != nil {
		return models.ServiceInstanceDetails{}, fmt.Errorf("Error creating new pubsub topic: %s", err)
	}

	i := models.ServiceInstanceDetails{
		Name:         params["topic_name"],
		Url:          "",
		Location:     "",
		OtherDetails: "{}",
	}

	extraDetails := make(map[string]string)

	if sub_name, ok := params["subscription_name"]; ok {
		var pushConfig *googlepubsub.PushConfig
		var ackDeadline = 10

		if ackd, ok := params["ack_deadline"]; ok {
			ackDeadline, err = strconv.Atoi(ackd)
			if err != nil {
				return models.ServiceInstanceDetails{}, fmt.Errorf("Error converting ack deadline to int: %s", err)
			}
		}

		if isPush, ok := params["is_push"]; ok {
			if isPush == "true" {
				pushConfig = &googlepubsub.PushConfig{
					Endpoint: params["endpoint"],
				}
			}
		}

		_, err = service.NewSubscription(ctx, sub_name, t, time.Duration(ackDeadline)*time.Second, pushConfig)
		if err != nil {
			return models.ServiceInstanceDetails{}, fmt.Errorf("Error creating subscription: %s", err)
		}

		extraDetails["subscription_name"] = params["subscription_name"]
	}

	otherDetails, err := json.Marshal(extraDetails)
	if err != nil {
		return models.ServiceInstanceDetails{}, fmt.Errorf("Error marshalling json: %s", err)
	}
	i.OtherDetails = string(otherDetails)

	return i, nil
}

// Deletes the topic associated with the given instanceID
func (b *PubSubBroker) Deprovision(instanceID string, details models.DeprovisionDetails) error {
	topic := models.ServiceInstanceDetails{}
	if err := db_service.DbConnection.Where("ID = ?", instanceID).First(&topic).Error; err != nil {
		return models.ErrInstanceDoesNotExist
	}

	ctx := context.Background()

	service, err := googlepubsub.NewClient(ctx, b.ProjectId)
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

	if subscriptionName, ok := otherD["subscription_name"]; ok {
		err = service.Subscription(subscriptionName).Delete(ctx)
		if err != nil {
			return fmt.Errorf("Error deleting subscription: %s", err)
		}
	}

	return nil
}
