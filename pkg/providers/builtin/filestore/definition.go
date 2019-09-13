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

package filestore

import (
	"code.cloudfoundry.org/lager"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/broker"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/providers/builtin/base"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/varcontext"
	"github.com/pivotal-cf/brokerapi"
	"golang.org/x/oauth2/jwt"
)

// ServiceDefinition creates a new ServiceDefinition object for the Firestore service.
func ServiceDefinition() *broker.ServiceDefinition {
	return &broker.ServiceDefinition{
		Id:          "494eb82e-c4ca-4bed-871d-9c3f02f66e01",
		Name:        "google-filestore",
		Description: "Fully managed NFS file storage with predictable performance.",
		DisplayName: "Google Cloud Filestore",
		// Filestore doesn't have a hex logo so we'll copy storage's.
		ImageUrl:         "https://cloud.google.com/_static/images/cloud/products/logos/svg/storage.svg",
		DocumentationUrl: "https://cloud.google.com/filestore/docs/",
		SupportUrl:       "https://cloud.google.com/filestore/docs/getting-support",
		Tags:             []string{"gcp", "filestore", "nfs"},
		Bindable:         true,
		PlanUpdateable:   false,
		Plans: []broker.ServicePlan{
			{
				ServicePlan: brokerapi.ServicePlan{
					ID:          "e4c83975-e60f-43cf-afde-ebec573c6c2e",
					Name:        "default",
					Description: "Filestore default plan.",
					Free:        brokerapi.FreeValue(false),
				},
				ServiceProperties: map[string]string{},
			},
		},
		ProvisionInputVariables: []broker.BrokerVariable{
			base.InstanceID(1, 63, base.ZoneArea),
			base.Zone("us-west1-a", "https://cloud.google.com/filestore/docs/regions"),
			{
				FieldName: "tier",
				Default:   "STANDARD",
				Type:      broker.JsonTypeString,
				Details:   "The performance tier.",
				Enum: map[interface{}]string{
					"STANDARD": "Standard Tier: 100 MB/s reads, 5000 IOPS",
					"PREMIUM":  "Premium Tier: 1.2 GB/s reads, 60000 IOPS",
				},
			},
			{
				FieldName: "authorized_network",
				Type:      broker.JsonTypeString,
				Details:   "The name of the network to attach the instance to.",
				Default:   "default",
			},
			{
				FieldName: "address_mode",
				Default:   "MODE_IPV4",
				Type:      broker.JsonTypeString,
				Details:   "The address mode of the service.",
				Enum: map[interface{}]string{
					"MODE_IPV4": "IPV4 Addressed",
				},
			},
			{
				FieldName: "capacity_gb",
				Type:      broker.JsonTypeInteger,
				Details:   "The capacity of the Filestore. Standard minimum is 1TiB and Premium is minimum 2.5TiB.",
				Default:   1024,
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
				FieldName: "ip_address",
				Type:      broker.JsonTypeString,
				Details:   "IP address of the service.",
			},
			{
				FieldName: "file_share_name",
				Type:      broker.JsonTypeString,
				Details:   "Name of the share.",
			},
			{
				FieldName: "capacity_gb",
				Type:      broker.JsonTypeInteger,
				Details:   "Capacity of the share in GiB.",
			},
			{
				FieldName: "uri",
				Type:      broker.JsonTypeString,
				Details:   "URI of the instance.",
			},
		},
		BindComputedVariables: []varcontext.DefaultVariable{},
		Examples: []broker.ServiceExample{
			{
				Name:            "Standard",
				Description:     "Creates a standard Filestore.",
				PlanId:          "e4c83975-e60f-43cf-afde-ebec573c6c2e",
				ProvisionParams: map[string]interface{}{},
				BindParams:      map[string]interface{}{},
			},
			{
				Name:        "Premium",
				Description: "Creates a premium Filestore.",
				PlanId:      "e4c83975-e60f-43cf-afde-ebec573c6c2e",
				ProvisionParams: map[string]interface{}{
					"tier":        "PREMIUM",
					"capacity_gb": 2560,
				},
				BindParams: map[string]interface{}{},
			},
		},
		ProviderBuilder: func(projectID string, auth *jwt.Config, logger lager.Logger) broker.ServiceProvider {
			bb := base.NewPeeredNetworkServiceBase(projectID, auth, logger)
			return &Broker{PeeredNetworkServiceBase: bb}
		},
		IsBuiltin: true,
	}
}
