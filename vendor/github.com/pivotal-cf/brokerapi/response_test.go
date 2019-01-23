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

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf/brokerapi"
)

var _ = Describe("Catalog Response", func() {
	Describe("JSON encoding", func() {
		It("has a list of services", func() {
			catalogResponse := brokerapi.CatalogResponse{
				Services: []brokerapi.Service{},
			}
			jsonString := `{"services":[]}`

			Expect(json.Marshal(catalogResponse)).To(MatchJSON(jsonString))
		})
	})
})

var _ = Describe("Provisioning Response", func() {
	Describe("JSON encoding", func() {
		Context("when the dashboard URL is not present", func() {
			It("does not return it in the JSON", func() {
				provisioningResponse := brokerapi.ProvisioningResponse{}
				jsonString := `{}`

				Expect(json.Marshal(provisioningResponse)).To(MatchJSON(jsonString))
			})
		})

		Context("when the dashboard URL is present", func() {
			It("returns it in the JSON", func() {
				provisioningResponse := brokerapi.ProvisioningResponse{
					DashboardURL: "http://example.com/broker",
				}
				jsonString := `{"dashboard_url":"http://example.com/broker"}`

				Expect(json.Marshal(provisioningResponse)).To(MatchJSON(jsonString))
			})
		})
	})
})

var _ = Describe("Update Response", func() {
	Describe("JSON encoding", func() {
		Context("when the dashboard URL is not present", func() {
			It("does not return it in the JSON", func() {
				updateResponse := brokerapi.UpdateResponse{}
				jsonString := `{}`

				Expect(json.Marshal(updateResponse)).To(MatchJSON(jsonString))
			})
		})

		Context("when the dashboard URL is present", func() {
			It("returns it in the JSON", func() {
				updateResponse := brokerapi.UpdateResponse{
					DashboardURL: "http://example.com/broker_updated",
				}
				jsonString := `{"dashboard_url":"http://example.com/broker_updated"}`

				Expect(json.Marshal(updateResponse)).To(MatchJSON(jsonString))
			})
		})
	})
})

var _ = Describe("Binding Response", func() {
	Describe("JSON encoding", func() {
		It("has a credentials object", func() {
			binding := brokerapi.BindingResponse{}
			jsonString := `{"credentials":null}`

			Expect(json.Marshal(binding)).To(MatchJSON(jsonString))
		})
	})
})

var _ = Describe("Error Response", func() {
	Describe("JSON encoding", func() {
		It("has a description field", func() {
			errorResponse := brokerapi.ErrorResponse{
				Description: "a bad thing happened",
			}
			jsonString := `{"description":"a bad thing happened"}`

			Expect(json.Marshal(errorResponse)).To(MatchJSON(jsonString))
		})
	})
})
