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

package bigquery

import (
	"code.cloudfoundry.org/lager"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/providers/builtin/base"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/broker"
	accountmanagers "github.com/GoogleCloudPlatform/gcp-service-broker/pkg/providers/builtin/account_managers"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/validation"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/varcontext"
	"github.com/pivotal-cf/brokerapi"
	"golang.org/x/oauth2/jwt"
)

// ServiceDefinition creates a new ServiceDefinition object for the BigQuery service.
func ServiceDefinition() *broker.ServiceDefinition {
	roleWhitelist := []string{
		"bigquery.dataViewer",
		"bigquery.dataEditor",
		"bigquery.dataOwner",
		"bigquery.user",
		"bigquery.jobUser",
	}

	return &broker.ServiceDefinition{
		Id:               "f80c0a3e-bd4d-4809-a900-b4e33a6450f1",
		Name:             "google-bigquery",
		Description:      "A fast, economical and fully managed data warehouse for large-scale data analytics.",
		DisplayName:      "Google BigQuery",
		ImageUrl:         "https://cloud.google.com/_static/images/cloud/products/logos/svg/bigquery.svg",
		DocumentationUrl: "https://cloud.google.com/bigquery/docs/",
		SupportUrl:       "https://cloud.google.com/bigquery/support",
		Tags:             []string{"gcp", "bigquery"},
		Bindable:         true,
		PlanUpdateable:   false,
		Plans: []broker.ServicePlan{
			{
				ServicePlan: brokerapi.ServicePlan{
					ID:          "10ff4e72-6e84-44eb-851f-bdb38a791914",
					Name:        "default",
					Description: "BigQuery default plan.",
					Free:        brokerapi.FreeValue(false),
				},
				ServiceProperties: map[string]string{},
			},
		},
		ProvisionInputVariables: []broker.BrokerVariable{
			{
				FieldName: "name",
				Type:      broker.JsonTypeString,
				Details:   "The name of the BigQuery dataset.",
				Expression:   "pcf_sb_${counter.next()}_${time.nano()}",
				Constraints: validation.NewConstraintBuilder().
					Pattern("^[A-Za-z0-9_]+$").
					MaxLength(1024).
					Build(),
			},
			{
				FieldName: "location",
				Type:      broker.JsonTypeString,
				Details:   "The location of the BigQuery instance.",
				Default:   "US",
				Constraints: validation.NewConstraintBuilder().
					Pattern("^[A-Za-z][-a-z0-9A-Z]+$").
					Examples("US", "EU", "asia-northeast1").
					Build(),
			},
		},
		ProvisionComputedVariables: []varcontext.DefaultVariable{
			{Name: "labels", Expression: "${json.marshal(request.default_labels)}", Overwrite: true},
		},
		DefaultRoleWhitelist:  roleWhitelist,
		BindInputVariables:    accountmanagers.ServiceAccountWhitelistWithDefault(roleWhitelist, "bigquery.user"),
		BindComputedVariables: accountmanagers.ServiceAccountBindComputedVariables(),
		BindOutputVariables: append(accountmanagers.ServiceAccountBindOutputVariables(),
			broker.BrokerVariable{
				FieldName: "dataset_id",
				Type:      broker.JsonTypeString,
				Details:   "The name of the BigQuery dataset.",
				Required:  true,
				Constraints: validation.NewConstraintBuilder().
					Pattern("^[A-Za-z0-9_]+$").
					MaxLength(1024).
					Build(),
			},
		),
		Examples: []broker.ServiceExample{
			{
				Name:        "Basic Configuration",
				Description: "Create a dataset and account that can manage and query the data.",
				PlanId:      "10ff4e72-6e84-44eb-851f-bdb38a791914",
				ProvisionParams: map[string]interface{}{
					"name": "orders_1997",
				},
				BindParams: map[string]interface{}{
					"role": "bigquery.user",
				},
			},
		},
		ProviderBuilder: func(projectId string, auth *jwt.Config, logger lager.Logger) broker.ServiceProvider {
			bb := base.NewBrokerBase(projectId, auth, logger)
			return &BigQueryBroker{BrokerBase: bb}
		},
		IsBuiltin: true,
	}
}
