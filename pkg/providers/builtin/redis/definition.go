// Copyright 2019 the Service Broker Project Authors.
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

package redis

import (
	"code.cloudfoundry.org/lager"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/broker"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/providers/builtin/base"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/varcontext"
	"github.com/pivotal-cf/brokerapi"
	"golang.org/x/oauth2/jwt"
)

func ServiceDefinition() *broker.ServiceDefinition {
	return &broker.ServiceDefinition{
		Id:               "3ea92b54-838c-4fe1-b75d-9bda513380aa",
		Name:             "google-memorystore-redis",
		Description:      "Creates and manages Redis instances on the Google Cloud Platform.",
		DisplayName:      "Google Cloud Memorystore for Redis API",
		ImageUrl:         "https://cloud.google.com/_static/images/cloud/products/logos/svg/cache.svg",
		DocumentationUrl: "https://cloud.google.com/memorystore/docs/redis",
		SupportUrl:       "https://cloud.google.com/memorystore/docs/redis/support",
		Tags:             []string{"gcp", "memorystore", "redis"},
		Bindable:         true,
		PlanUpdateable:   false,
		Plans: []broker.ServicePlan{
			{
				ServicePlan: brokerapi.ServicePlan{
					ID:          "df10762e-6ef1-44e3-84c2-07e9358ceb1f",
					Name:        "default",
					Description: "Lets you chose your own values for all properties.",
					Free:        brokerapi.FreeValue(false),
				},
				ProvisionOverrides: map[string]interface{}{},
			},
			{
				ServicePlan: brokerapi.ServicePlan{
					ID:          "dd1923b6-ac26-4697-83d6-b3a0c05c2c94",
					Name:        "basic",
					Description: "Provides a standalone Redis instance. Use this tier for applications that require a simple Redis cache.",
					Free:        brokerapi.FreeValue(false),
				},
				ProvisionOverrides: map[string]interface{}{"service_tier": "BASIC"},
			},
			{
				ServicePlan: brokerapi.ServicePlan{
					ID:          "41771881-b456-4940-9081-34b6424744c6",
					Name:        "standard_ha",
					Description: "Provides a highly available Redis instance.",
					Free:        brokerapi.FreeValue(false),
				},
				ProvisionOverrides: map[string]interface{}{"service_tier": "STANDARD_HA"},
			},
		},
		ProvisionInputVariables: []broker.BrokerVariable{
			base.InstanceID(1, 40, base.ProjectArea),
			base.AuthorizedNetwork(),
			base.Region("us-east1", "https://cloud.google.com/memorystore/docs/redis/regions"),
			{
				FieldName: "memory_size_gb",
				Type:      broker.JsonTypeInteger,
				Details:   "Redis memory size in GiB.",
				Default:   4,
			},
			{
				FieldName: "tier",
				Type:      broker.JsonTypeString,
				Details:   "The performance tier.",
				Default:   "BASIC",
				Enum: map[interface{}]string{
					"BASIC":       "Standalone instance, good for caching.",
					"STANDARD_HA": "Highly available primary/replica, good for databases.",
				},
			},
		},
		ProvisionComputedVariables: []varcontext.DefaultVariable{
			{Name: "labels", Default: "${json.marshal(request.default_labels)}", Overwrite: true},
		},
		BindInputVariables: []broker.BrokerVariable{},
		BindOutputVariables: []broker.BrokerVariable{
			{
				FieldName: "authorized_network",
				Type:      broker.JsonTypeString,
				Details:   "Name of the VPC network the instance is attached to.",
			},
			{
				FieldName: "reserved_ip_range",
				Type:      broker.JsonTypeString,
				Details:   "Range of IP addresses reserved for the instance.",
			},
			{
				FieldName: "redis_version",
				Type:      broker.JsonTypeString,
				Details:   "The version of Redis software.",
			},
			{
				FieldName: "memory_size_gb",
				Type:      broker.JsonTypeInteger,
				Details:   "Redis memory size in GiB.",
			},
			{
				FieldName: "host",
				Type:      broker.JsonTypeString,
				Details:   "Hostname or IP address of the exposed Redis endpoint used by clients to connect to the service.",
			},
			{
				FieldName: "port",
				Type:      broker.JsonTypeInteger,
				Details:   "The port number of the exposed Redis endpoint.",
			},
			{
				FieldName: "uri",
				Type:      broker.JsonTypeString,
				Details:   "URI of the instance.",
			},
		},
		PlanVariables: []broker.BrokerVariable{},
		Examples: []broker.ServiceExample{
			{
				Name:            "Standard Redis Configuration",
				Description:     "Create a Redis instance with standard service tier.",
				PlanId:          "dd1923b6-ac26-4697-83d6-b3a0c05c2c94",
				ProvisionParams: map[string]interface{}{},
				BindParams:      map[string]interface{}{},
			},
			{
				Name:            "HA Redis Configuration",
				Description:     "Create a Redis instance with high availability.",
				PlanId:          "41771881-b456-4940-9081-34b6424744c6",
				ProvisionParams: map[string]interface{}{},
				BindParams:      map[string]interface{}{},
			},
		},
		ProviderBuilder: func(projectID string, auth *jwt.Config, logger lager.Logger) broker.ServiceProvider {
			return &Broker{
				PeeredNetworkServiceBase: base.NewPeeredNetworkServiceBase(projectID, auth, logger),
			}
		},
		IsBuiltin: true,
	}
}
