// Copyright (C) 2015-Present Pivotal Software, Inc. All rights reserved.

// This program and the accompanying materials are made available under
// the terms of the under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package brokerapi_test

import (
	"encoding/json"
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf/brokerapi"
)

var _ = Describe("Catalog", func() {
	Describe("Service", func() {
		Describe("JSON encoding", func() {
			It("uses the correct keys", func() {
				service := brokerapi.Service{
					ID:            "ID-1",
					Name:          "Cassandra",
					Description:   "A Cassandra Plan",
					Bindable:      true,
					Plans:         []brokerapi.ServicePlan{},
					Metadata:      &brokerapi.ServiceMetadata{},
					Tags:          []string{"test"},
					PlanUpdatable: true,
					DashboardClient: &brokerapi.ServiceDashboardClient{
						ID:          "Dashboard ID",
						Secret:      "dashboardsecret",
						RedirectURI: "the.dashboa.rd",
					},
				}
				jsonString := `{
					"id":"ID-1",
				  	"name":"Cassandra",
					"description":"A Cassandra Plan",
					"bindable":true,
					"plan_updateable":true,
					"tags":["test"],
					"plans":[],
					"dashboard_client":{
						"id":"Dashboard ID",
						"secret":"dashboardsecret",
						"redirect_uri":"the.dashboa.rd"
					},
					"metadata":{

					}
				}`
				Expect(json.Marshal(service)).To(MatchJSON(jsonString))
			})
		})

		It("encodes the optional 'requires' fields", func() {
			service := brokerapi.Service{
				ID:            "ID-1",
				Name:          "Cassandra",
				Description:   "A Cassandra Plan",
				Bindable:      true,
				Plans:         []brokerapi.ServicePlan{},
				Metadata:      &brokerapi.ServiceMetadata{},
				Tags:          []string{"test"},
				PlanUpdatable: true,
				Requires: []brokerapi.RequiredPermission{
					brokerapi.PermissionRouteForwarding,
					brokerapi.PermissionSyslogDrain,
					brokerapi.PermissionVolumeMount,
				},
				DashboardClient: &brokerapi.ServiceDashboardClient{
					ID:          "Dashboard ID",
					Secret:      "dashboardsecret",
					RedirectURI: "the.dashboa.rd",
				},
			}
			jsonString := `{
				"id":"ID-1",
					"name":"Cassandra",
				"description":"A Cassandra Plan",
				"bindable":true,
				"plan_updateable":true,
				"tags":["test"],
				"plans":[],
				"requires": ["route_forwarding", "syslog_drain", "volume_mount"],
				"dashboard_client":{
					"id":"Dashboard ID",
					"secret":"dashboardsecret",
					"redirect_uri":"the.dashboa.rd"
				},
				"metadata":{

				}
			}`
			Expect(json.Marshal(service)).To(MatchJSON(jsonString))
		})
	})

	Describe("ServicePlan", func() {
		Describe("JSON encoding", func() {
			It("uses the correct keys", func() {
				plan := brokerapi.ServicePlan{
					ID:          "ID-1",
					Name:        "Cassandra",
					Description: "A Cassandra Plan",
					Bindable:    brokerapi.BindableValue(true),
					Free:        brokerapi.FreeValue(true),
					Metadata: &brokerapi.ServicePlanMetadata{
						Bullets:     []string{"hello", "its me"},
						DisplayName: "name",
					},
				}
				jsonString := `{
					"id":"ID-1",
					"name":"Cassandra",
					"description":"A Cassandra Plan",
					"free": true,
					"bindable": true,
					"metadata":{
						"bullets":["hello", "its me"],
						"displayName":"name"
					}
				}`

				Expect(json.Marshal(plan)).To(MatchJSON(jsonString))
			})
		})
	})

	Describe("ServicePlanMetadata", func() {
		Describe("JSON encoding", func() {
			It("uses the correct keys", func() {
				metadata := brokerapi.ServicePlanMetadata{
					Bullets:     []string{"test"},
					DisplayName: "Some display name",
				}
				jsonString := `{"bullets":["test"],"displayName":"Some display name"}`

				Expect(json.Marshal(metadata)).To(MatchJSON(jsonString))
			})

			It("encodes the AdditionalMetadata fields in the metadata fields", func() {
				metadata := brokerapi.ServicePlanMetadata{
					Bullets:     []string{"hello", "its me"},
					DisplayName: "name",
					AdditionalMetadata: map[string]interface{}{
						"foo": "bar",
						"baz": 1,
					},
				}
				jsonString := `{
					"bullets":["hello", "its me"],
					"displayName":"name",
					"foo": "bar",
					"baz": 1
				}`

				Expect(json.Marshal(metadata)).To(MatchJSON(jsonString))
			})
		})

		Describe("JSON decoding", func() {
			It("sets the AdditionalMetadata from unrecognized fields", func() {
				metadata := brokerapi.ServicePlanMetadata{}
				jsonString := `{"foo":["test"],"bar":"Some display name"}`

				err := json.Unmarshal([]byte(jsonString), &metadata)
				Expect(err).NotTo(HaveOccurred())
				Expect(metadata.AdditionalMetadata["foo"]).To(Equal([]interface{}{"test"}))
				Expect(metadata.AdditionalMetadata["bar"]).To(Equal("Some display name"))
			})

			It("does not include convention fields into additional metadata", func() {
				metadata := brokerapi.ServicePlanMetadata{}
				jsonString := `{"bullets":["test"],"displayName":"Some display name", "costs": [{"amount": {"usd": 649.0},"unit": "MONTHLY"}]}`

				err := json.Unmarshal([]byte(jsonString), &metadata)
				Expect(err).NotTo(HaveOccurred())
				Expect(metadata.AdditionalMetadata).To(BeNil())
			})
		})
	})

	Describe("ServiceMetadata", func() {
		Describe("JSON encoding", func() {
			It("uses the correct keys", func() {
				shareable := true
				metadata := brokerapi.ServiceMetadata{
					DisplayName:         "Cassandra",
					LongDescription:     "A long description of Cassandra",
					DocumentationUrl:    "doc",
					SupportUrl:          "support",
					ImageUrl:            "image",
					ProviderDisplayName: "display",
					Shareable:           &shareable,
				}
				jsonString := `{
					"displayName":"Cassandra",
					"longDescription":"A long description of Cassandra",
					"documentationUrl":"doc",
					"supportUrl":"support",
					"imageUrl":"image",
					"providerDisplayName":"display",
					"shareable":true
				}`

				Expect(json.Marshal(metadata)).To(MatchJSON(jsonString))
			})

			It("encodes the AdditionalMetadata fields in the metadata fields", func() {
				metadata := brokerapi.ServiceMetadata{
					DisplayName: "name",
					AdditionalMetadata: map[string]interface{}{
						"foo": "bar",
						"baz": 1,
					},
				}
				jsonString := `{
					"displayName":"name",
					"foo": "bar",
					"baz": 1
				}`

				Expect(json.Marshal(metadata)).To(MatchJSON(jsonString))
			})
		})

		Describe("JSON decoding", func() {
			It("sets the AdditionalMetadata from unrecognized fields", func() {
				metadata := brokerapi.ServiceMetadata{}
				jsonString := `{"foo":["test"],"bar":"Some display name"}`

				err := json.Unmarshal([]byte(jsonString), &metadata)
				Expect(err).NotTo(HaveOccurred())
				Expect(metadata.AdditionalMetadata["foo"]).To(Equal([]interface{}{"test"}))
				Expect(metadata.AdditionalMetadata["bar"]).To(Equal("Some display name"))
			})

			It("does not include convention fields into additional metadata", func() {
				metadata := brokerapi.ServiceMetadata{}
				jsonString := `{
					"displayName":"Cassandra",
					"longDescription":"A long description of Cassandra",
					"documentationUrl":"doc",
					"supportUrl":"support",
					"imageUrl":"image",
					"providerDisplayName":"display",
					"shareable":true
				}`
				err := json.Unmarshal([]byte(jsonString), &metadata)
				Expect(err).NotTo(HaveOccurred())
				Expect(metadata.AdditionalMetadata).To(BeNil())
			})
		})
	})

	It("Reflects JSON names from struct", func() {
		type Example1 struct {
			Foo int    `json:"foo"`
			Bar string `yaml:"hello" json:"bar,omitempty"`
			Qux float64
		}

		s := Example1{}
		Expect(brokerapi.GetJsonNames(reflect.ValueOf(&s).Elem())).To(
			ConsistOf([]string{"foo", "bar", "Qux"}))
	})
})
