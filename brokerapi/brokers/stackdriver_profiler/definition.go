// Copyright 2018 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package stackdriver_profiler

import (
	"code.cloudfoundry.org/lager"
	accountmanagers "github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/account_managers"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/broker_base"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/broker"
	"golang.org/x/oauth2/jwt"
)

func init() {
	bs := &broker.ServiceDefinition{
		Name: "google-stackdriver-profiler",
		DefaultServiceDefinition: `{
		      "id": "00b9ca4a-7cd6-406a-a5b7-2f43f41ade75",
		      "description": "Stackdriver Profiler",
		      "name": "google-stackdriver-profiler",
		      "bindable": true,
		      "plan_updateable": false,
		      "metadata": {
		        "displayName": "Stackdriver Profiler",
		        "longDescription": "Continuous CPU and heap profiling to improve performance and reduce costs.",
		        "documentationUrl": "https://cloud.google.com/profiler/docs/",
		        "supportUrl": "https://cloud.google.com/support/",
		        "imageUrl": "https://cloud.google.com/_static/images/cloud/products/logos/svg/stackdriver.svg",
		        "shareable": true
		      },
		      "tags": ["gcp", "stackdriver", "profiler"],
		      "plans": [
		        {
		          "id": "594627f6-35f5-462f-9074-10fb033fb18a",
		          "service_id": "00b9ca4a-7cd6-406a-a5b7-2f43f41ade75",
		          "name": "default",
		          "display_name": "Default",
		          "description": "Stackdriver Profiler default plan.",
		          "service_properties": {}
		        }
		      ]
				}
		`,
		ProvisionInputVariables: []broker.BrokerVariable{},
		BindInputVariables:      []broker.BrokerVariable{},
		BindComputedVariables:   accountmanagers.FixedRoleBindComputedVariables("cloudprofiler.agent"),
		BindOutputVariables:     accountmanagers.ServiceAccountBindOutputVariables(),
		Examples: []broker.ServiceExample{
			{
				Name:            "Basic Configuration",
				Description:     "Creates an account with the permission `cloudprofiler.agent`.",
				PlanId:          "594627f6-35f5-462f-9074-10fb033fb18a",
				ProvisionParams: map[string]interface{}{},
				BindParams:      map[string]interface{}{},
			},
		},
		ProviderBuilder: func(projectId string, auth *jwt.Config, logger lager.Logger) broker.ServiceProvider {
			bb := broker_base.NewBrokerBase(projectId, auth, logger)
			return &StackdriverProfilerBroker{BrokerBase: bb}
		},
	}

	broker.Register(bs)
}
