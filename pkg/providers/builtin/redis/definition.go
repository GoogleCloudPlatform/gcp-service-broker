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
	accountmanagers "github.com/GoogleCloudPlatform/gcp-service-broker/pkg/providers/builtin/account_managers"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/providers/builtin/base"
	. "github.com/GoogleCloudPlatform/gcp-service-broker/pkg/providers/builtin/common"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/validation"
	"github.com/pivotal-cf/brokerapi"
	"golang.org/x/oauth2/jwt"
)

func ServiceDefinition() *broker.ServiceDefinition {
	roleWhitelist := []string{
		"redis.editor",
		"redis.viewer",
	}

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
					ID:          "dd1923b6-ac26-4697-83d6-b3a0c05c2c94",
					Name:        "basic",
					Description: "Provides a standalone Redis instance. Use this tier for applications that require a simple Redis cache.",
					Free:        brokerapi.FreeValue(false),
				},
				ServiceProperties: map[string]string{"service_tier": "basic"},
			},
			{
				ServicePlan: brokerapi.ServicePlan{
					ID:          "41771881-b456-4940-9081-34b6424744c6",
					Name:        "standard_ha",
					Description: "Provides a highly available Redis instance that includes automatically enabled cross-zone replication and automatic failover. Use this tier for applications that require high availability for a Redis instance.",
					Free:        brokerapi.FreeValue(false),
				},
				ServiceProperties: map[string]string{"service_tier": "standard_ha"},
			},
		},
		ProvisionInputVariables: []broker.BrokerVariable{
			{
				FieldName: "authorized_network",
				Type:      broker.JsonTypeString,
				Details:   "Optional. The full name of the Google Compute Engine network to which the instance is connected.",
				Default:   "",
			},
			{
				FieldName: "memory_size_gb",
				Type:      broker.JsonTypeString,
				Details:   "The Redis instance's provisioned capacity in GB. See: https://cloud.google.com/memorystore/pricing for more information.",
				Default:   "4",
				Constraints: validation.NewConstraintBuilder().
					MinLength(1).
					MaxLength(10).
					Pattern("[1-9][0-9]*").
					Build(),
			},
			{
				FieldName: "instance_id",
				Type:      broker.JsonTypeString,
				Details:   "The name of the Redis instance.",
				Default:   "gsb-${counter.next()}-${time.nano()}",
				Constraints: validation.NewConstraintBuilder().
					MinLength(1).
					MaxLength(40).
					Pattern("^[a-z]([-0-9a-z]*[a-z0-9]$)*").
					Build(),
			},
			{
				FieldName: "region",
				Type:      broker.JsonTypeString,
				Details:   "The region in which to provision the Redis instance. See: https://cloud.google.com/memorystore/docs/redis/regions for supported regions.",
				Default:   UsEast1.Region(),
				Enum: map[interface{}]string{
					AsiaEast1.Region():              AsiaEast1.Region(),
					AsiaEast2.Region():              AsiaEast2.Region(),
					AsiaNorthEast1.Region():         AsiaNorthEast1.Region(),
					AsiaSouth1.Region():             AsiaSouth1.Region(),
					AsiaSouthEast1.Region():         AsiaSouthEast1.Region(),
					AustraliaSouthEast1.Region():    AustraliaSouthEast1.Region(),
					EuropeNorth1.Region():           EuropeNorth1.Region(),
					EuropeWest1.Region():            EuropeWest1.Region(),
					EuropeWest2.Region():            EuropeWest2.Region(),
					EuropeWest3.Region():            EuropeWest3.Region(),
					EuropeWest4.Region():            EuropeWest4.Region(),
					NorthAmericaNorthEast1.Region(): NorthAmericaNorthEast1.Region(),
					SouthAmericaEast1.Region():      SouthAmericaEast1.Region(),
					UsCentral1.Region():             UsCentral1.Region(),
					UsEast1.Region():                UsEast1.Region(),
					UsEast4.Region():                UsEast4.Region(),
					UsWest1.Region():                UsWest1.Region(),
					UsWest2.Region():                UsWest2.Region(),
				},
			},
			{
				FieldName: "display_name",
				Type:      broker.JsonTypeString,
				Details:   "The human-readable display name of the Redis instance.",
				Default:   "${instance_id}",
				Constraints: validation.NewConstraintBuilder().
					MinLength(4).
					MaxLength(30).
					Build(),
			},
		},
		DefaultRoleWhitelist: roleWhitelist,
		BindInputVariables:   accountmanagers.ServiceAccountWhitelistWithDefault(roleWhitelist, "redis.viewer"),
		BindOutputVariables: append(accountmanagers.ServiceAccountBindOutputVariables(),
			broker.BrokerVariable{
				FieldName: "redis_version",
				Type:      broker.JsonTypeString,
				Details:   "The version of Redis software.",
				Required:  true,
			},
			broker.BrokerVariable{
				FieldName: "host",
				Type:      broker.JsonTypeString,
				Details:   "Hostname or IP address of the exposed Redis endpoint used by clients to connect to the service.",
				Required:  true,
			},
			broker.BrokerVariable{
				FieldName: "port",
				Type:      broker.JsonTypeString,
				Details:   "The port number of the exposed Redis endpoint.",
				Required:  true,
			},
			broker.BrokerVariable{
				FieldName: "memory_size_gb",
				Type:      broker.JsonTypeInteger,
				Details:   "Redis memory size in GiB.",
				Required:  true,
			},
		),
		BindComputedVariables: accountmanagers.ServiceAccountBindComputedVariables(),
		PlanVariables: []broker.BrokerVariable{
			{
				FieldName: "service_tier",
				Type:      broker.JsonTypeString,
				Details:   "Either BASIC or STANDARD_HA. See: https://cloud.google.com/memorystore/pricing for more information.",
				Default:   "basic",
				Required:  true,
			},
		},
		Examples: []broker.ServiceExample{
			{
				Name:            "Basic Redis Configuration",
				Description:     "Create a Redis instance with basic service tier.",
				PlanId:          "dd1923b6-ac26-4697-83d6-b3a0c05c2c94",
				ProvisionParams: map[string]interface{}{},
				BindParams: map[string]interface{}{
					"role": "redis.viewer",
				},
			},
		},
		ProviderBuilder: func(projectId string, auth *jwt.Config, logger lager.Logger) broker.ServiceProvider {
			bb := base.NewBrokerBase(projectId, auth, logger)
			return &RedisBroker{BrokerBase: bb}
		},
		IsBuiltin: true,
	}
}
