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
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"

	"code.cloudfoundry.org/lager"
	"code.cloudfoundry.org/lager/lagertest"
	"github.com/drewolson/testflight"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf/brokerapi"
	"github.com/pivotal-cf/brokerapi/fakes"
)

var _ = Describe("Service Broker API", func() {
	var fakeServiceBroker *fakes.FakeServiceBroker
	var brokerAPI http.Handler
	var brokerLogger *lagertest.TestLogger
	var apiVersion string
	var credentials = brokerapi.BrokerCredentials{
		Username: "username",
		Password: "password",
	}

	makeInstanceProvisioningRequest := func(instanceID string, details map[string]interface{}, queryString string) *testflight.Response {
		response := &testflight.Response{}

		testflight.WithServer(brokerAPI, func(r *testflight.Requester) {
			path := "/v2/service_instances/" + instanceID + queryString

			buffer := &bytes.Buffer{}
			json.NewEncoder(buffer).Encode(details)
			request, err := http.NewRequest("PUT", path, buffer)
			Expect(err).NotTo(HaveOccurred())
			request.Header.Add("Content-Type", "application/json")
			if apiVersion != "" {
				request.Header.Add("X-Broker-API-Version", apiVersion)
			}
			request.SetBasicAuth(credentials.Username, credentials.Password)

			response = r.Do(request)
		})
		return response
	}

	makeInstanceProvisioningRequestWithAcceptsIncomplete := func(instanceID string, details map[string]interface{}, acceptsIncomplete bool) *testflight.Response {
		var acceptsIncompleteFlag string

		if acceptsIncomplete {
			acceptsIncompleteFlag = "?accepts_incomplete=true"
		} else {
			acceptsIncompleteFlag = "?accepts_incomplete=false"
		}

		return makeInstanceProvisioningRequest(instanceID, details, acceptsIncompleteFlag)
	}

	lastLogLine := func() lager.LogFormat {
		noOfLogLines := len(brokerLogger.Logs())
		if noOfLogLines == 0 {
			// better way to raise error?
			err := errors.New("expected some log lines but there were none")
			Expect(err).NotTo(HaveOccurred())
		}

		return brokerLogger.Logs()[noOfLogLines-1]
	}

	BeforeEach(func() {
		apiVersion = "2.14"
		fakeServiceBroker = &fakes.FakeServiceBroker{
			InstanceLimit: 3,
			ServiceID:     "0A789746-596F-4CEA-BFAC-A0795DA056E3",
			PlanID:        "plan-id",
		}
		brokerLogger = lagertest.NewTestLogger("broker-api")
		brokerAPI = brokerapi.New(fakeServiceBroker, brokerLogger, credentials)
	})

	Describe("response headers", func() {
		makeRequest := func() *httptest.ResponseRecorder {
			recorder := httptest.NewRecorder()
			request, _ := http.NewRequest("GET", "/v2/catalog", nil)
			request.SetBasicAuth(credentials.Username, credentials.Password)
			brokerAPI.ServeHTTP(recorder, request)
			return recorder
		}

		It("has a Content-Type header", func() {
			response := makeRequest()

			header := response.Header().Get("Content-Type")
			Ω(header).Should(Equal("application/json"))
		})
	})

	Describe("request context", func() {
		var (
			ctx     context.Context
			reqBody string
		)

		makeRequest := func(method, path, body string) *httptest.ResponseRecorder {
			recorder := httptest.NewRecorder()
			request, _ := http.NewRequest(method, path, strings.NewReader(body))
			request.Header.Add("X-Broker-API-Version", "2.13")
			request.SetBasicAuth(credentials.Username, credentials.Password)
			request = request.WithContext(ctx)
			brokerAPI.ServeHTTP(recorder, request)
			return recorder
		}

		BeforeEach(func() {
			ctx = context.WithValue(context.Background(), "test_context", true)
			reqBody = fmt.Sprintf(`{"service_id":"%s","plan_id":"456"}`, fakeServiceBroker.ServiceID)
		})

		Specify("a catalog endpoint which passes the request context to the broker", func() {
			makeRequest("GET", "/v2/catalog", "")
			Expect(fakeServiceBroker.ReceivedContext).To(BeTrue())
		})

		Specify("a provision endpoint which passes the request context to the broker", func() {
			makeRequest("PUT", "/v2/service_instances/instance-id", reqBody)
			Expect(fakeServiceBroker.ReceivedContext).To(BeTrue())
		})

		Specify("a deprovision endpoint which passes the request context to the broker", func() {
			makeRequest("DELETE", "/v2/service_instances/instance-id?service_id=asdf&plan_id=fdsa", "")
			Expect(fakeServiceBroker.ReceivedContext).To(BeTrue())
		})

		Specify("a bind endpoint which passes the request context to the broker", func() {
			makeRequest("PUT", "/v2/service_instances/instance-id/service_bindings/binding-id", reqBody)
			Expect(fakeServiceBroker.ReceivedContext).To(BeTrue())
		})

		Specify("an unbind endpoint which passes the request context to the broker", func() {
			makeRequest("DELETE", "/v2/service_instances/instance-id/service_bindings/binding-id?plan_id=plan-id&service_id=service-id", "")
			Expect(fakeServiceBroker.ReceivedContext).To(BeTrue())
		})

		Specify("an update endpoint which passes the request context to the broker", func() {
			makeRequest("PATCH", "/v2/service_instances/instance-id", `{"service_id":"123"}`)
			Expect(fakeServiceBroker.ReceivedContext).To(BeTrue())
		})

		Specify("a last operation endpoint which passes the request context to the broker", func() {
			makeRequest("GET", "/v2/service_instances/instance-id/last_operation", "{}")
			Expect(fakeServiceBroker.ReceivedContext).To(BeTrue())
		})
	})

	Describe("authentication", func() {
		makeRequestWithoutAuth := func() *testflight.Response {
			response := &testflight.Response{}
			testflight.WithServer(brokerAPI, func(r *testflight.Requester) {
				request, _ := http.NewRequest("GET", "/v2/catalog", nil)
				response = r.Do(request)
			})
			return response
		}

		makeRequestWithAuth := func(username string, password string) *testflight.Response {
			response := &testflight.Response{}
			testflight.WithServer(brokerAPI, func(r *testflight.Requester) {
				request, _ := http.NewRequest("GET", "/v2/catalog", nil)
				request.SetBasicAuth(username, password)

				response = r.Do(request)
			})
			return response
		}

		makeRequestWithUnrecognizedAuth := func() *testflight.Response {
			response := &testflight.Response{}
			testflight.WithServer(brokerAPI, func(r *testflight.Requester) {
				request, _ := http.NewRequest("GET", "/v2/catalog", nil)
				// dXNlcm5hbWU6cGFzc3dvcmQ= is base64 encoding of 'username:password',
				// ie, a correctly encoded basic authorization header
				request.Header["Authorization"] = []string{"NOTBASIC dXNlcm5hbWU6cGFzc3dvcmQ="}

				response = r.Do(request)
			})
			return response
		}

		It("returns 401 when the authorization header has an incorrect password", func() {
			response := makeRequestWithAuth("username", "fake_password")
			Expect(response.StatusCode).To(Equal(401))
		})

		It("returns 401 when the authorization header has an incorrect username", func() {
			response := makeRequestWithAuth("fake_username", "password")
			Expect(response.StatusCode).To(Equal(401))
		})

		It("returns 401 when there is no authorization header", func() {
			response := makeRequestWithoutAuth()
			Expect(response.StatusCode).To(Equal(401))
		})

		It("returns 401 when there is a unrecognized authorization header", func() {
			response := makeRequestWithUnrecognizedAuth()
			Expect(response.StatusCode).To(Equal(401))
		})

		It("does not call through to the service broker when not authenticated", func() {
			makeRequestWithAuth("username", "fake_password")
			Ω(fakeServiceBroker.BrokerCalled).ShouldNot(BeTrue(),
				"broker should not have been hit when authentication failed",
			)
		})
	})

	Describe("catalog endpoint", func() {
		makeCatalogRequest := func(apiVersion string, fail bool) *httptest.ResponseRecorder {
			recorder := httptest.NewRecorder()
			request, _ := http.NewRequest(http.MethodGet, "/v2/catalog", nil)
			if apiVersion != "" {
				request.Header.Add("X-Broker-API-Version", apiVersion)
			}
			request.SetBasicAuth(credentials.Username, credentials.Password)
			ctx := context.Background()
			if fail {
				ctx = context.WithValue(ctx, "fails", true)
			}
			request = request.WithContext(ctx)
			brokerAPI.ServeHTTP(recorder, request)
			return recorder
		}

		It("returns a 200", func() {
			response := makeCatalogRequest("2.14", false)
			Expect(response.Code).To(Equal(200))
		})

		It("returns valid catalog json", func() {
			response := makeCatalogRequest("2.14", false)
			Expect(response.Body).To(MatchJSON(fixture("catalog.json")))
		})

		It("returns a 500", func() {
			response := makeCatalogRequest("2.14", true)
			Expect(response.Code).To(Equal(500))
			Expect(response.Body.String()).To(MatchJSON(`{ "description": "something went wrong!" }`))
		})

		Context("the request is malformed", func() {
			It("missing header X-Broker-API-Version", func() {
				response := makeCatalogRequest("", false)
				Expect(response.Code).To(Equal(412))
				Expect(lastLogLine().Message).To(ContainSubstring(".catalog.broker-api-version-invalid"))
				Expect(lastLogLine().Data["error"]).To(ContainSubstring("X-Broker-API-Version Header not set"))
			})

			It("has wrong version of API", func() {
				response := makeCatalogRequest("1.14", false)
				Expect(response.Code).To(Equal(412))
				Expect(lastLogLine().Message).To(ContainSubstring(".catalog.broker-api-version-invalid"))
				Expect(lastLogLine().Data["error"]).To(ContainSubstring("X-Broker-API-Version Header must be 2.x"))
			})
		})
	})

	Describe("instance lifecycle endpoint", func() {
		makeInstanceDeprovisioningRequestFull := func(instanceID, serviceID, planID, queryString string) *testflight.Response {
			response := &testflight.Response{}
			testflight.WithServer(brokerAPI, func(r *testflight.Requester) {
				path := fmt.Sprintf("/v2/service_instances/%s?plan_id=%s&service_id=%s", instanceID, planID, serviceID)
				if queryString != "" {
					path = fmt.Sprintf("%s&%s", path, queryString)
				}
				request, err := http.NewRequest("DELETE", path, strings.NewReader(""))
				Expect(err).NotTo(HaveOccurred())
				request.Header.Add("Content-Type", "application/json")
				request.SetBasicAuth("username", "password")
				request.Header.Add("X-Broker-API-Version", apiVersion)
				response = r.Do(request)

			})
			return response
		}

		makeInstanceDeprovisioningRequest := func(instanceID, queryString string) *testflight.Response {
			return makeInstanceDeprovisioningRequestFull(instanceID, "service-id", "plan-id", queryString)
		}

		Describe("provisioning", func() {
			var instanceID string
			var provisionDetails map[string]interface{}

			BeforeEach(func() {
				instanceID = uniqueInstanceID()
				provisionDetails = map[string]interface{}{
					"service_id":        fakeServiceBroker.ServiceID,
					"plan_id":           "plan-id",
					"organization_guid": "organization-guid",
					"space_guid":        "space-guid",
				}
			})

			It("calls Provision on the service broker with all params", func() {
				makeInstanceProvisioningRequest(instanceID, provisionDetails, "")
				Expect(fakeServiceBroker.ProvisionDetails).To(Equal(brokerapi.ProvisionDetails{
					ServiceID:        fakeServiceBroker.ServiceID,
					PlanID:           "plan-id",
					OrganizationGUID: "organization-guid",
					SpaceGUID:        "space-guid",
				}))
			})

			It("calls Provision on the service broker with the instance id", func() {
				makeInstanceProvisioningRequest(instanceID, provisionDetails, "")
				Expect(fakeServiceBroker.ProvisionedInstanceIDs).To(ContainElement(instanceID))
			})

			Context("when the broker returns some operation data", func() {
				BeforeEach(func() {
					fakeServiceBroker = &fakes.FakeServiceBroker{
						InstanceLimit:         3,
						OperationDataToReturn: "some-operation-data",
						ServiceID:             fakeServiceBroker.ServiceID,
						PlanID:                fakeServiceBroker.PlanID,
					}
					fakeAsyncServiceBroker := &fakes.FakeAsyncServiceBroker{
						FakeServiceBroker:    *fakeServiceBroker,
						ShouldProvisionAsync: true,
					}
					brokerAPI = brokerapi.New(fakeAsyncServiceBroker, brokerLogger, credentials)
				})

				It("returns the operation data to the cloud controller", func() {
					resp := makeInstanceProvisioningRequest(instanceID, provisionDetails, "")
					Expect(resp.Body).To(MatchJSON(fixture("operation_data_response.json")))
				})
			})

			Context("when there are arbitrary params", func() {
				var rawParams string
				var rawCtx string

				BeforeEach(func() {
					provisionDetails["parameters"] = map[string]interface{}{
						"string": "some-string",
						"number": 1,
						"object": struct{ Name string }{"some-name"},
						"array":  []interface{}{"a", "b", "c"},
					}
					rawParams = `{
					"string":"some-string",
					"number":1,
					"object": { "Name": "some-name" },
					"array": [ "a", "b", "c" ]
				}`
					provisionDetails["context"] = map[string]interface{}{
						"platform":      "fake-platform",
						"serial-number": 12648430,
						"object":        struct{ Name string }{"parameter"},
						"array":         []interface{}{"1", "2", "3"},
					}
					rawCtx = `{
					"platform":"fake-platform",
					"serial-number":12648430,
					"object": {"Name":"parameter"},
					"array":[ "1", "2", "3" ]
				}`
				})

				It("calls Provision on the service broker with all params", func() {
					makeInstanceProvisioningRequest(instanceID, provisionDetails, "")
					Expect(string(fakeServiceBroker.ProvisionDetails.RawParameters)).To(MatchJSON(rawParams))
				})

				It("calls Provision with details with raw parameters", func() {
					makeInstanceProvisioningRequest(instanceID, provisionDetails, "")
					detailsWithRawParameters := brokerapi.DetailsWithRawParameters(fakeServiceBroker.ProvisionDetails)
					rawParameters := detailsWithRawParameters.GetRawParameters()
					Expect(string(rawParameters)).To(MatchJSON(rawParams))
				})

				It("calls Provision with details with raw context", func() {
					makeInstanceProvisioningRequest(instanceID, provisionDetails, "")
					detailsWithRawContext := brokerapi.DetailsWithRawContext(fakeServiceBroker.ProvisionDetails)
					rawContext := detailsWithRawContext.GetRawContext()
					Expect(string(rawContext)).To(MatchJSON(rawCtx))
				})
			})

			Context("when the instance does not exist", func() {
				It("returns a 201", func() {
					response := makeInstanceProvisioningRequest(instanceID, provisionDetails, "")
					Expect(response.StatusCode).To(Equal(201))
				})

				It("returns empty json", func() {
					response := makeInstanceProvisioningRequest(instanceID, provisionDetails, "")
					Expect(response.Body).To(MatchJSON(fixture("provisioning.json")))
				})

				Context("when the broker returns a dashboard URL", func() {
					BeforeEach(func() {
						fakeServiceBroker.DashboardURL = "some-dashboard-url"
					})

					It("returns json with dasboard URL", func() {
						response := makeInstanceProvisioningRequest(instanceID, provisionDetails, "")
						Expect(response.Body).To(MatchJSON(fixture("provisioning_with_dashboard.json")))
					})
				})

				Context("when the instance limit has been reached", func() {
					BeforeEach(func() {
						for i := 0; i < fakeServiceBroker.InstanceLimit; i++ {
							makeInstanceProvisioningRequest(uniqueInstanceID(), provisionDetails, "")
						}
					})

					It("returns a 500", func() {
						response := makeInstanceProvisioningRequest(instanceID, provisionDetails, "")
						Expect(response.StatusCode).To(Equal(500))
					})

					It("returns json with a description field and a useful error message", func() {
						response := makeInstanceProvisioningRequest(instanceID, provisionDetails, "")
						Expect(response.Body).To(MatchJSON(fixture("instance_limit_error.json")))
					})

					It("logs an appropriate error", func() {
						makeInstanceProvisioningRequest(instanceID, provisionDetails, "")

						Expect(lastLogLine().Message).To(ContainSubstring(".provision.instance-limit-reached"))
						Expect(lastLogLine().Data["error"]).To(ContainSubstring("instance limit for this service has been reached"))
					})
				})

				Context("when an unexpected error occurs", func() {
					BeforeEach(func() {
						fakeServiceBroker.ProvisionError = errors.New("broker failed")
					})

					It("returns a 500", func() {
						response := makeInstanceProvisioningRequest(instanceID, provisionDetails, "")
						Expect(response.StatusCode).To(Equal(500))
					})

					It("returns json with a description field and a useful error message", func() {
						response := makeInstanceProvisioningRequest(instanceID, provisionDetails, "")
						Expect(response.Body).To(MatchJSON(`{"description":"broker failed"}`))
					})

					It("logs an appropriate error", func() {
						makeInstanceProvisioningRequest(instanceID, provisionDetails, "")
						Expect(lastLogLine().Message).To(ContainSubstring(".provision.unknown-error"))
						Expect(lastLogLine().Data["error"]).To(ContainSubstring("broker failed"))
					})
				})

				Context("when a custom error occurs", func() {
					BeforeEach(func() {
						fakeServiceBroker.ProvisionError = brokerapi.NewFailureResponse(
							errors.New("I failed in unique and interesting ways"),
							http.StatusTeapot,
							"interesting-failure",
						)
					})

					It("returns status teapot", func() {
						response := makeInstanceProvisioningRequest(instanceID, provisionDetails, "")
						Expect(response.StatusCode).To(Equal(http.StatusTeapot))
					})

					It("returns json with a description field and a useful error message", func() {
						response := makeInstanceProvisioningRequest(instanceID, provisionDetails, "")
						Expect(response.Body).To(MatchJSON(`{"description":"I failed in unique and interesting ways"}`))
					})

					It("logs an appropriate error", func() {
						makeInstanceProvisioningRequest(instanceID, provisionDetails, "")
						Expect(lastLogLine().Message).To(ContainSubstring(".provision.interesting-failure"))
						Expect(lastLogLine().Data["error"]).To(ContainSubstring("I failed in unique and interesting ways"))
					})
				})

				Context("RawParameters are not valid JSON", func() {
					BeforeEach(func() {
						fakeServiceBroker.ProvisionError = brokerapi.ErrRawParamsInvalid
					})

					It("returns a 422", func() {
						response := makeInstanceProvisioningRequest(instanceID, provisionDetails, "")
						Expect(response.StatusCode).To(Equal(http.StatusUnprocessableEntity))
					})

					It("returns json with a description field and a useful error message", func() {
						response := makeInstanceProvisioningRequest(instanceID, provisionDetails, "")
						Expect(response.Body).To(MatchJSON(`{"description":"The format of the parameters is not valid JSON"}`))
					})

					It("logs an appropriate error", func() {
						makeInstanceProvisioningRequest(instanceID, provisionDetails, "")
						Expect(lastLogLine().Message).To(ContainSubstring(".provision.invalid-raw-params"))
						Expect(lastLogLine().Data["error"]).To(ContainSubstring("The format of the parameters is not valid JSON"))
					})
				})

				Context("when we send invalid json", func() {
					makeBadInstanceProvisioningRequest := func(instanceID string) *testflight.Response {
						response := &testflight.Response{}

						testflight.WithServer(brokerAPI, func(r *testflight.Requester) {
							path := "/v2/service_instances/" + instanceID

							body := strings.NewReader("{{{{{")
							request, err := http.NewRequest("PUT", path, body)
							Expect(err).NotTo(HaveOccurred())
							request.Header.Add("Content-Type", "application/json")
							if apiVersion != "" {
								request.Header.Add("X-Broker-Api-Version", apiVersion)
							}
							request.SetBasicAuth(credentials.Username, credentials.Password)

							response = r.Do(request)
						})

						return response
					}

					It("returns a 422 bad request", func() {
						response := makeBadInstanceProvisioningRequest(instanceID)
						Expect(response.StatusCode).Should(Equal(http.StatusUnprocessableEntity))
					})

					It("logs a message", func() {
						makeBadInstanceProvisioningRequest(instanceID)
						Expect(lastLogLine().Message).To(ContainSubstring(".provision.invalid-service-details"))
					})
				})
			})

			Context("when the instance already exists", func() {
				BeforeEach(func() {
					makeInstanceProvisioningRequest(instanceID, provisionDetails, "")
				})

				It("returns a 409", func() {
					response := makeInstanceProvisioningRequest(instanceID, provisionDetails, "")
					Expect(response.StatusCode).To(Equal(409))
				})

				It("returns an empty JSON object", func() {
					response := makeInstanceProvisioningRequest(instanceID, provisionDetails, "")
					Expect(response.Body).To(MatchJSON(`{}`))
				})

				It("logs an appropriate error", func() {
					makeInstanceProvisioningRequest(instanceID, provisionDetails, "")
					Expect(lastLogLine().Message).To(ContainSubstring(".provision.instance-already-exists"))
					Expect(lastLogLine().Data["error"]).To(ContainSubstring("instance already exists"))
				})
			})

			Describe("accepts_incomplete", func() {
				Context("when the accepts_incomplete flag is true", func() {
					It("calls ProvisionAsync on the service broker", func() {
						acceptsIncomplete := true
						makeInstanceProvisioningRequestWithAcceptsIncomplete(instanceID, provisionDetails, acceptsIncomplete)
						Expect(fakeServiceBroker.ProvisionDetails).To(Equal(brokerapi.ProvisionDetails{
							ServiceID:        fakeServiceBroker.ServiceID,
							PlanID:           "plan-id",
							OrganizationGUID: "organization-guid",
							SpaceGUID:        "space-guid",
						}))

						Expect(fakeServiceBroker.ProvisionedInstanceIDs).To(ContainElement(instanceID))
					})

					Context("when the broker chooses to provision asynchronously", func() {
						BeforeEach(func() {
							fakeServiceBroker = &fakes.FakeServiceBroker{
								InstanceLimit: 3,
								ServiceID:     fakeServiceBroker.ServiceID,
								PlanID:        fakeServiceBroker.PlanID,
							}
							fakeAsyncServiceBroker := &fakes.FakeAsyncServiceBroker{
								FakeServiceBroker:    *fakeServiceBroker,
								ShouldProvisionAsync: true,
							}
							brokerAPI = brokerapi.New(fakeAsyncServiceBroker, brokerLogger, credentials)
						})

						It("returns a 202", func() {
							response := makeInstanceProvisioningRequestWithAcceptsIncomplete(instanceID, provisionDetails, true)
							Expect(response.StatusCode).To(Equal(http.StatusAccepted))
						})
					})

					Context("when the broker chooses to provision synchronously", func() {
						BeforeEach(func() {
							fakeServiceBroker = &fakes.FakeServiceBroker{
								InstanceLimit: 3,
								ServiceID:     fakeServiceBroker.ServiceID,
								PlanID:        fakeServiceBroker.PlanID,
							}
							fakeAsyncServiceBroker := &fakes.FakeAsyncServiceBroker{
								FakeServiceBroker:    *fakeServiceBroker,
								ShouldProvisionAsync: false,
							}
							brokerAPI = brokerapi.New(fakeAsyncServiceBroker, brokerLogger, credentials)
						})

						It("returns a 201", func() {
							response := makeInstanceProvisioningRequestWithAcceptsIncomplete(instanceID, provisionDetails, true)
							Expect(response.StatusCode).To(Equal(http.StatusCreated))
						})
					})
				})

				Context("when the accepts_incomplete flag is false", func() {
					It("returns a 201", func() {
						response := makeInstanceProvisioningRequestWithAcceptsIncomplete(instanceID, provisionDetails, false)
						Expect(response.StatusCode).To(Equal(http.StatusCreated))
					})

					Context("when broker can only respond asynchronously", func() {
						BeforeEach(func() {
							fakeServiceBroker = &fakes.FakeServiceBroker{
								InstanceLimit: 3,
								ServiceID:     fakeServiceBroker.ServiceID,
								PlanID:        fakeServiceBroker.PlanID,
							}
							fakeAsyncServiceBroker := &fakes.FakeAsyncOnlyServiceBroker{
								FakeServiceBroker: *fakeServiceBroker,
							}
							brokerAPI = brokerapi.New(fakeAsyncServiceBroker, brokerLogger, credentials)
						})

						It("returns a 422", func() {
							acceptsIncomplete := false
							response := makeInstanceProvisioningRequestWithAcceptsIncomplete(instanceID, provisionDetails, acceptsIncomplete)
							Expect(response.StatusCode).To(Equal(http.StatusUnprocessableEntity))
							Expect(response.Body).To(MatchJSON(fixture("async_required.json")))
						})
					})
				})

				Context("when the accepts_incomplete flag is missing", func() {
					It("returns a 201", func() {
						response := makeInstanceProvisioningRequest(instanceID, provisionDetails, "")
						Expect(response.StatusCode).To(Equal(http.StatusCreated))
					})

					Context("when broker can only respond asynchronously", func() {
						BeforeEach(func() {
							fakeServiceBroker = &fakes.FakeServiceBroker{
								InstanceLimit: 3,
								ServiceID:     fakeServiceBroker.ServiceID,
								PlanID:        fakeServiceBroker.PlanID,
							}
							fakeAsyncServiceBroker := &fakes.FakeAsyncOnlyServiceBroker{
								FakeServiceBroker: *fakeServiceBroker,
							}
							brokerAPI = brokerapi.New(fakeAsyncServiceBroker, brokerLogger, credentials)
						})

						It("returns a 422", func() {
							acceptsIncomplete := false
							response := makeInstanceProvisioningRequestWithAcceptsIncomplete(instanceID, provisionDetails, acceptsIncomplete)
							Expect(response.StatusCode).To(Equal(http.StatusUnprocessableEntity))
							Expect(response.Body).To(MatchJSON(fixture("async_required.json")))
						})
					})
				})
			})

			Context("the request is malformed", func() {
				It("missing header X-Broker-API-Version", func() {
					apiVersion = ""
					response := makeInstanceProvisioningRequest(instanceID, provisionDetails, "")
					Expect(response.StatusCode).To(Equal(412))
					Expect(lastLogLine().Message).To(ContainSubstring(".provision.broker-api-version-invalid"))
					Expect(lastLogLine().Data["error"]).To(ContainSubstring("X-Broker-API-Version Header not set"))
				})

				It("has wrong version of API", func() {
					apiVersion = "1.14"
					response := makeInstanceProvisioningRequest(instanceID, provisionDetails, "")
					Expect(response.StatusCode).To(Equal(412))
					Expect(lastLogLine().Message).To(ContainSubstring(".provision.broker-api-version-invalid"))
					Expect(lastLogLine().Data["error"]).To(ContainSubstring("X-Broker-API-Version Header must be 2.x"))
				})

				It("missing service_id", func() {
					delete(provisionDetails, "service_id")
					response := makeInstanceProvisioningRequest(instanceID, provisionDetails, "")
					Expect(response.StatusCode).To(Equal(400))
					Expect(lastLogLine().Message).To(ContainSubstring(".provision.service-id-missing"))
					Expect(lastLogLine().Data["error"]).To(ContainSubstring("service_id missing"))
				})

				It("missing plan_id", func() {
					delete(provisionDetails, "plan_id")
					response := makeInstanceProvisioningRequest(instanceID, provisionDetails, "")
					Expect(response.StatusCode).To(Equal(400))
					Expect(lastLogLine().Message).To(ContainSubstring(".provision.plan-id-missing"))
					Expect(lastLogLine().Data["error"]).To(ContainSubstring("plan_id missing"))
				})

				It("service_id not in the catalog", func() {
					provisionDetails["service_id"] = "not-in-the-catalogue"
					response := makeInstanceProvisioningRequest(instanceID, provisionDetails, "")
					Expect(response.StatusCode).To(Equal(400))
					Expect(lastLogLine().Message).To(ContainSubstring(".provision.invalid-service-id"))
					Expect(lastLogLine().Data["error"]).To(ContainSubstring("service-id not in the catalog"))
				})

				It("plan_id not in the catalog", func() {
					provisionDetails["plan_id"] = "not-in-the-catalogue"
					response := makeInstanceProvisioningRequest(instanceID, provisionDetails, "")
					Expect(response.StatusCode).To(Equal(400))
					Expect(lastLogLine().Message).To(ContainSubstring(".provision.invalid-plan-id"))
					Expect(lastLogLine().Data["error"]).To(ContainSubstring("plan-id not in the catalog"))
				})
			})
		})

		Describe("updating", func() {
			var (
				instanceID  string
				details     map[string]interface{}
				queryString string
				response    *testflight.Response
			)

			makeInstanceUpdateRequest := func(instanceID string, details map[string]interface{}, queryString string, apiVersion string) *testflight.Response {
				response := &testflight.Response{}

				testflight.WithServer(brokerAPI, func(r *testflight.Requester) {
					path := "/v2/service_instances/" + instanceID + queryString

					buffer := &bytes.Buffer{}
					json.NewEncoder(buffer).Encode(details)
					request, err := http.NewRequest("PATCH", path, buffer)
					Expect(err).NotTo(HaveOccurred())
					if apiVersion != "" {
						request.Header.Add("X-Broker-Api-Version", apiVersion)
					}
					request.Header.Add("Content-Type", "application/json")
					request.SetBasicAuth(credentials.Username, credentials.Password)

					response = r.Do(request)
				})
				return response
			}

			BeforeEach(func() {
				instanceID = uniqueInstanceID()
				details = map[string]interface{}{
					"service_id": "some-service-id",
					"plan_id":    "new-plan",
					"parameters": map[string]interface{}{
						"new-param": "new-param-value",
					},
					"previous_values": map[string]interface{}{
						"service_id":      "service-id",
						"plan_id":         "old-plan",
						"organization_id": "org-id",
						"space_id":        "space-id",
					},
					"context": map[string]interface{}{
						"new-context": "new-context-value",
					},
				}
				queryString = "?accept_incomplete=true"
			})

			JustBeforeEach(func() {
				response = makeInstanceUpdateRequest(instanceID, details, queryString, "2.14")
			})

			Context("the request is malformed", func() {
				It("missing header X-Broker-API-Version", func() {
					response := makeInstanceUpdateRequest("instance-id", details, queryString, "")
					Expect(response.StatusCode).To(Equal(412))
					Expect(lastLogLine().Message).To(ContainSubstring(".update.broker-api-version-invalid"))
					Expect(lastLogLine().Data["error"]).To(ContainSubstring("X-Broker-API-Version Header not set"))
				})

				It("has wrong version of API", func() {
					response := makeInstanceUpdateRequest("instance-id", details, queryString, "1.14")
					Expect(response.StatusCode).To(Equal(412))
					Expect(lastLogLine().Message).To(ContainSubstring(".update.broker-api-version-invalid"))
					Expect(lastLogLine().Data["error"]).To(ContainSubstring("X-Broker-API-Version Header must be 2.x"))
				})

				It("missing service-id", func() {
					delete(details, "service_id")
					response := makeInstanceUpdateRequest("instance-id", details, queryString, "2.14")
					Expect(response.StatusCode).To(Equal(400))
					Expect(lastLogLine().Message).To(ContainSubstring(".update.service-id-missing"))
					Expect(lastLogLine().Data["error"]).To(ContainSubstring("service_id missing"))
				})
			})

			Context("when the broker returns no error", func() {
				Context("when the broker responds synchronously", func() {
					It("returns HTTP 200", func() {
						Expect(response.StatusCode).To(Equal(http.StatusOK))
					})

					It("returns JSON content type", func() {
						Expect(response.RawResponse.Header.Get("Content-Type")).To(Equal("application/json"))
					})

					It("returns empty JSON body", func() {
						Expect(response.Body).To(Equal("{}\n"))
					})

					It("calls broker with instanceID and update details", func() {
						Expect(fakeServiceBroker.UpdatedInstanceIDs).To(ConsistOf(instanceID))
						Expect(fakeServiceBroker.UpdateDetails.ServiceID).To(Equal("some-service-id"))
						Expect(fakeServiceBroker.UpdateDetails.PlanID).To(Equal("new-plan"))
						Expect(fakeServiceBroker.UpdateDetails.PreviousValues).To(Equal(brokerapi.PreviousValues{
							PlanID:    "old-plan",
							ServiceID: "service-id",
							OrgID:     "org-id",
							SpaceID:   "space-id",
						},
						))
						Expect(fakeServiceBroker.UpdateDetails.RawParameters).To(Equal(json.RawMessage(`{"new-param":"new-param-value"}`)))
					})

					It("calls update with details with raw parameters", func() {
						detailsWithRawParameters := brokerapi.DetailsWithRawParameters(fakeServiceBroker.UpdateDetails)
						rawParameters := detailsWithRawParameters.GetRawParameters()
						Expect(rawParameters).To(Equal(json.RawMessage(`{"new-param":"new-param-value"}`)))
					})

					It("calls update with details with raw context", func() {
						Expect(fakeServiceBroker.UpdateDetails.RawContext).To(
							Equal(json.RawMessage(`{"new-context":"new-context-value"}`)),
						)
					})

					Context("when accepts_incomplete=true", func() {
						BeforeEach(func() {
							queryString = "?accepts_incomplete=true"
						})

						It("tells broker async is allowed", func() {
							Expect(fakeServiceBroker.AsyncAllowed).To(BeTrue())
						})
					})

					Context("when accepts_incomplete is not supplied", func() {
						BeforeEach(func() {
							queryString = ""
						})

						It("tells broker async not allowed", func() {
							Expect(fakeServiceBroker.AsyncAllowed).To(BeFalse())
						})
					})
				})

				Context("when the broker responds asynchronously", func() {
					BeforeEach(func() {
						fakeServiceBroker.ShouldReturnAsync = true
					})

					It("returns HTTP 202", func() {
						Expect(response.StatusCode).To(Equal(http.StatusAccepted))
					})

					Context("when the broker responds with operation data", func() {
						BeforeEach(func() {
							fakeServiceBroker.OperationDataToReturn = "some-operation-data"
						})

						It("returns the operation data to the cloud controller", func() {
							Expect(response.Body).To(MatchJSON(fixture("operation_data_response.json")))
						})
					})
				})
			})

			Context("when the broker indicates that it needs async support", func() {
				BeforeEach(func() {
					fakeServiceBroker.UpdateError = brokerapi.ErrAsyncRequired
				})

				It("returns HTTP 422", func() {
					Expect(response.StatusCode).To(Equal(http.StatusUnprocessableEntity))
				})

				It("returns a descriptive message", func() {
					var body map[string]string
					err := json.Unmarshal([]byte(response.Body), &body)
					Expect(err).ToNot(HaveOccurred())
					Expect(body["error"]).To(Equal("AsyncRequired"))
					Expect(body["description"]).To(Equal("This service plan requires client support for asynchronous service operations."))
				})
			})

			Context("when the broker indicates that the plan cannot be upgraded", func() {
				BeforeEach(func() {
					fakeServiceBroker.UpdateError = brokerapi.ErrPlanChangeNotSupported
				})

				It("returns HTTP 422", func() {
					Expect(response.StatusCode).To(Equal(http.StatusUnprocessableEntity))
				})

				It("returns a descriptive message", func() {
					var body map[string]string
					err := json.Unmarshal([]byte(response.Body), &body)
					Expect(err).ToNot(HaveOccurred())
					Expect(body["error"]).To(Equal("PlanChangeNotSupported"))
					Expect(body["description"]).To(Equal("The requested plan migration cannot be performed"))
				})
			})

			Context("when the broker errors in an unknown way", func() {
				BeforeEach(func() {
					fakeServiceBroker.UpdateError = errors.New("some horrible internal error")
				})

				It("returns HTTP 500", func() {
					Expect(response.StatusCode).To(Equal(500))
				})

				It("returns a descriptive message", func() {
					var body map[string]string
					err := json.Unmarshal([]byte(response.Body), &body)
					Expect(err).ToNot(HaveOccurred())
					Expect(body["description"]).To(Equal("some horrible internal error"))
				})
			})
		})

		Describe("deprovisioning", func() {
			It("calls Deprovision on the service broker with the instance id", func() {
				instanceID := uniqueInstanceID()
				makeInstanceDeprovisioningRequest(instanceID, "")
				Expect(fakeServiceBroker.DeprovisionedInstanceIDs).To(ContainElement(instanceID))
			})

			Context("when the instance exists", func() {
				var instanceID string
				var provisionDetails map[string]interface{}

				BeforeEach(func() {
					instanceID = uniqueInstanceID()

					provisionDetails = map[string]interface{}{
						"service_id":        fakeServiceBroker.ServiceID,
						"plan_id":           "plan-id",
						"organization_guid": "organization-guid",
						"space_guid":        "space-guid",
					}
					makeInstanceProvisioningRequest(instanceID, provisionDetails, "")
				})

				itReturnsStatus := func(expectedStatus int, queryString string) {
					It(fmt.Sprintf("returns HTTP %d", expectedStatus), func() {
						response := makeInstanceDeprovisioningRequest(instanceID, queryString)
						Expect(response.StatusCode).To(Equal(expectedStatus))
					})
				}

				itReturnsEmptyJsonObject := func(queryString string) {
					It("returns an empty JSON object", func() {
						response := makeInstanceDeprovisioningRequest(instanceID, queryString)
						Expect(response.Body).To(MatchJSON(`{}`))
					})
				}

				Context("when the broker can only operate synchronously", func() {
					Context("when the accepts_incomplete flag is not set", func() {
						itReturnsStatus(200, "")
						itReturnsEmptyJsonObject("")
					})

					Context("when the accepts_incomplete flag is set to true", func() {
						itReturnsStatus(200, "accepts_incomplete=true")
						itReturnsEmptyJsonObject("accepts_incomplete=true")
					})
				})

				Context("when the broker can only operate asynchronously", func() {
					BeforeEach(func() {
						fakeAsyncServiceBroker := &fakes.FakeAsyncOnlyServiceBroker{
							FakeServiceBroker: *fakeServiceBroker,
						}
						brokerAPI = brokerapi.New(fakeAsyncServiceBroker, brokerLogger, credentials)
					})

					Context("when the accepts_incomplete flag is not set", func() {
						itReturnsStatus(http.StatusUnprocessableEntity, "")

						It("returns a descriptive error", func() {
							response := makeInstanceDeprovisioningRequest(instanceID, "")
							Expect(response.Body).To(MatchJSON(fixture("async_required.json")))
						})
					})

					Context("when the accepts_incomplete flag is set to true", func() {
						itReturnsStatus(202, "accepts_incomplete=true")
						itReturnsEmptyJsonObject("accepts_incomplete=true")
					})

					Context("when the broker returns operation data", func() {
						BeforeEach(func() {
							fakeServiceBroker.OperationDataToReturn = "some-operation-data"
							fakeAsyncServiceBroker := &fakes.FakeAsyncOnlyServiceBroker{
								FakeServiceBroker: *fakeServiceBroker,
							}
							brokerAPI = brokerapi.New(fakeAsyncServiceBroker, brokerLogger, credentials)
						})

						itReturnsStatus(202, "accepts_incomplete=true")

						It("returns the operation data to the cloud controller", func() {
							response := makeInstanceDeprovisioningRequest(instanceID, "accepts_incomplete=true")
							Expect(response.Body).To(MatchJSON(fixture("operation_data_response.json")))
						})
					})
				})

				Context("when the broker can operate both synchronously and asynchronously", func() {
					BeforeEach(func() {
						fakeAsyncServiceBroker := &fakes.FakeAsyncServiceBroker{
							FakeServiceBroker: *fakeServiceBroker,
						}
						brokerAPI = brokerapi.New(fakeAsyncServiceBroker, brokerLogger, credentials)
					})

					Context("when the accepts_incomplete flag is not set", func() {
						itReturnsStatus(200, "")
						itReturnsEmptyJsonObject("")
					})

					Context("when the accepts_incomplete flag is set to true", func() {
						itReturnsStatus(202, "accepts_incomplete=true")
						itReturnsEmptyJsonObject("accepts_incomplete=true")
					})
				})

				It("contains plan_id", func() {
					makeInstanceDeprovisioningRequest(instanceID, "")
					Expect(fakeServiceBroker.DeprovisionDetails.PlanID).To(Equal("plan-id"))
				})

				It("contains service_id", func() {
					makeInstanceDeprovisioningRequest(instanceID, "")
					Expect(fakeServiceBroker.DeprovisionDetails.ServiceID).To(Equal("service-id"))
				})
			})

			Context("when the instance does not exist", func() {
				var instanceID string

				It("returns a 410", func() {
					response := makeInstanceDeprovisioningRequest(uniqueInstanceID(), "")
					Expect(response.StatusCode).To(Equal(410))
				})

				It("returns an empty JSON object", func() {
					response := makeInstanceDeprovisioningRequest(uniqueInstanceID(), "")
					Expect(response.Body).To(MatchJSON(`{}`))
				})

				It("logs an appropriate error", func() {
					instanceID = uniqueInstanceID()
					makeInstanceDeprovisioningRequest(instanceID, "")
					Expect(lastLogLine().Message).To(ContainSubstring(".deprovision.instance-missing"))
					Expect(lastLogLine().Data["error"]).To(ContainSubstring("instance does not exist"))
				})
			})

			Context("when instance deprovisioning fails", func() {
				var instanceID string
				var provisionDetails map[string]interface{}

				BeforeEach(func() {
					instanceID = uniqueInstanceID()
					provisionDetails = map[string]interface{}{
						"plan_id":           "plan-id",
						"organization_guid": "organization-guid",
						"space_guid":        "space-guid",
					}
					makeInstanceProvisioningRequest(instanceID, provisionDetails, "")
				})

				Context("when an unexpected error occurs", func() {
					BeforeEach(func() {
						fakeServiceBroker.DeprovisionError = errors.New("broker failed")
					})

					It("returns a 500", func() {
						response := makeInstanceDeprovisioningRequest(instanceID, "")
						Expect(response.StatusCode).To(Equal(500))
					})

					It("returns json with a description field and a useful error message", func() {
						response := makeInstanceDeprovisioningRequest(instanceID, "")
						Expect(response.Body).To(MatchJSON(`{"description":"broker failed"}`))
					})

					It("logs an appropriate error", func() {
						makeInstanceDeprovisioningRequest(instanceID, "")
						Expect(lastLogLine().Message).To(ContainSubstring(".deprovision.unknown-error"))
						Expect(lastLogLine().Data["error"]).To(ContainSubstring("broker failed"))
					})
				})

				Context("when a custom error occurs", func() {
					BeforeEach(func() {
						fakeServiceBroker.DeprovisionError = brokerapi.NewFailureResponse(
							errors.New("I failed in unique and interesting ways"),
							http.StatusTeapot,
							"interesting-failure",
						)
					})

					It("returns status teapot", func() {
						response := makeInstanceDeprovisioningRequest(instanceID, "")
						Expect(response.StatusCode).To(Equal(http.StatusTeapot))
					})

					It("returns json with a description field and a useful error message", func() {
						response := makeInstanceDeprovisioningRequest(instanceID, "")
						Expect(response.Body).To(MatchJSON(`{"description":"I failed in unique and interesting ways"}`))
					})

					It("logs an appropriate error", func() {
						makeInstanceDeprovisioningRequest(instanceID, "")
						Expect(lastLogLine().Message).To(ContainSubstring(".deprovision.interesting-failure"))
						Expect(lastLogLine().Data["error"]).To(ContainSubstring("I failed in unique and interesting ways"))
					})
				})
			})

			Context("the request is malformed", func() {
				It("missing header X-Broker-API-Version", func() {
					apiVersion = ""
					response := makeInstanceDeprovisioningRequestFull("instance-id", "service-id", "plan-id", "")
					Expect(response.StatusCode).To(Equal(412))
					Expect(lastLogLine().Message).To(ContainSubstring(".deprovision.broker-api-version-invalid"))
					Expect(lastLogLine().Data["error"]).To(ContainSubstring("X-Broker-API-Version Header not set"))
				})

				It("has wrong version of API", func() {
					apiVersion = "1.1"
					response := makeInstanceDeprovisioningRequestFull("instance-id", "service-id", "plan-id", "")
					Expect(response.StatusCode).To(Equal(412))
					Expect(lastLogLine().Message).To(ContainSubstring(".deprovision.broker-api-version-invalid"))
					Expect(lastLogLine().Data["error"]).To(ContainSubstring("X-Broker-API-Version Header must be 2.x"))
				})

				It("missing service-id", func() {
					response := makeInstanceDeprovisioningRequestFull("instance-id", "", "plan-id", "")
					Expect(response.StatusCode).To(Equal(400))
					Expect(lastLogLine().Message).To(ContainSubstring(".deprovision.service-id-missing"))
					Expect(lastLogLine().Data["error"]).To(ContainSubstring("service_id missing"))
				})

				It("missing plan-id", func() {
					response := makeInstanceDeprovisioningRequestFull("instance-id", "service-id", "", "")
					Expect(response.StatusCode).To(Equal(400))
					Expect(lastLogLine().Message).To(ContainSubstring(".deprovision.plan-id-missing"))
					Expect(lastLogLine().Data["error"]).To(ContainSubstring("plan_id missing"))
				})
			})
		})
	})

	Describe("binding lifecycle endpoint", func() {
		makeBindingRequestWithSpecificAPIVersion := func(instanceID, bindingID string, details map[string]interface{}, apiVersion string) *testflight.Response {
			response := &testflight.Response{}
			testflight.WithServer(brokerAPI, func(r *testflight.Requester) {
				path := fmt.Sprintf("/v2/service_instances/%s/service_bindings/%s",
					instanceID, bindingID)

				buffer := &bytes.Buffer{}

				if details != nil {
					json.NewEncoder(buffer).Encode(details)
				}

				request, err := http.NewRequest("PUT", path, buffer)

				Expect(err).NotTo(HaveOccurred())

				if apiVersion != "" {
					request.Header.Add("X-Broker-Api-Version", apiVersion)
				}
				request.Header.Add("Content-Type", "application/json")
				request.SetBasicAuth("username", "password")

				response = r.Do(request)
			})
			return response
		}

		makeBindingRequest := func(instanceID, bindingID string, details map[string]interface{}) *testflight.Response {
			return makeBindingRequestWithSpecificAPIVersion(instanceID, bindingID, details, "2.10")
		}

		Describe("binding", func() {
			var (
				instanceID string
				bindingID  string
				details    map[string]interface{}
			)

			BeforeEach(func() {
				instanceID = uniqueInstanceID()
				bindingID = uniqueBindingID()
				details = map[string]interface{}{
					"app_guid":   "app_guid",
					"plan_id":    "plan_id",
					"service_id": "service_id",
					"parameters": map[string]interface{}{
						"new-param": "new-param-value",
					},
				}
			})

			Context("the request is malformed", func() {
				BeforeEach(func() {
					bindingID = uniqueBindingID()
				})

				It("missing header X-Broker-API-Version", func() {
					response := makeBindingRequestWithSpecificAPIVersion(instanceID, bindingID, map[string]interface{}{}, "")
					Expect(response.StatusCode).To(Equal(412))
					Expect(lastLogLine().Message).To(ContainSubstring(".bind.broker-api-version-invalid"))
					Expect(lastLogLine().Data["error"]).To(ContainSubstring("X-Broker-API-Version Header not set"))
				})

				It("has wrong version of API", func() {
					response := makeBindingRequestWithSpecificAPIVersion(instanceID, bindingID, map[string]interface{}{}, "1.14")
					Expect(response.StatusCode).To(Equal(412))
					Expect(lastLogLine().Message).To(ContainSubstring(".bind.broker-api-version-invalid"))
					Expect(lastLogLine().Data["error"]).To(ContainSubstring("X-Broker-API-Version Header must be 2.x"))
				})

				It("missing service-id", func() {
					response := makeBindingRequestWithSpecificAPIVersion(instanceID, bindingID, map[string]interface{}{"plan_id": "123"}, "2.14")
					Expect(response.StatusCode).To(Equal(400))
					Expect(lastLogLine().Message).To(ContainSubstring(".bind.service-id-missing"))
					Expect(lastLogLine().Data["error"]).To(ContainSubstring("service_id missing"))
				})

				It("missing plan-id", func() {
					response := makeBindingRequestWithSpecificAPIVersion(instanceID, bindingID, map[string]interface{}{"service_id": "123"}, "2.14")
					Expect(response.StatusCode).To(Equal(400))
					Expect(lastLogLine().Message).To(ContainSubstring(".bind.plan-id-missing"))
					Expect(lastLogLine().Data["error"]).To(ContainSubstring("plan_id missing"))
				})
			})

			Context("when the associated instance exists", func() {
				It("calls Bind on the service broker with the instance and binding ids", func() {
					makeBindingRequest(instanceID, bindingID, details)
					Expect(fakeServiceBroker.BoundInstanceIDs).To(ContainElement(instanceID))
					Expect(fakeServiceBroker.BoundBindingIDs).To(ContainElement(bindingID))
					Expect(fakeServiceBroker.BoundBindingDetails).To(Equal(brokerapi.BindDetails{
						AppGUID:       "app_guid",
						PlanID:        "plan_id",
						ServiceID:     "service_id",
						RawParameters: json.RawMessage(`{"new-param":"new-param-value"}`),
					}))
				})

				It("calls bind with details with raw parameters", func() {
					makeBindingRequest(instanceID, bindingID, details)
					detailsWithRawParameters := brokerapi.DetailsWithRawParameters(fakeServiceBroker.BoundBindingDetails)
					rawParameters := detailsWithRawParameters.GetRawParameters()
					Expect(rawParameters).To(Equal(json.RawMessage(`{"new-param":"new-param-value"}`)))
				})

				It("returns the credentials returned by Bind", func() {
					response := makeBindingRequest(uniqueInstanceID(), uniqueBindingID(), details)
					Expect(response.Body).To(MatchJSON(fixture("binding.json")))
				})

				It("returns a 201", func() {
					response := makeBindingRequest(uniqueInstanceID(), uniqueBindingID(), details)
					Expect(response.StatusCode).To(Equal(201))
				})

				Context("when syslog_drain_url is being passed", func() {
					BeforeEach(func() {
						fakeServiceBroker.SyslogDrainURL = "some-drain-url"
					})

					It("responds with the syslog drain url", func() {
						response := makeBindingRequest(uniqueInstanceID(), uniqueBindingID(), details)
						Expect(response.Body).To(MatchJSON(fixture("binding_with_syslog.json")))
					})
				})

				Context("when route_service_url is being passed", func() {
					BeforeEach(func() {
						fakeServiceBroker.RouteServiceURL = "some-route-url"
					})

					It("responds with the route service url", func() {
						response := makeBindingRequest(uniqueInstanceID(), uniqueBindingID(), details)
						Expect(response.Body).To(MatchJSON(fixture("binding_with_route_service.json")))
					})
				})

				Context("when a volume mount is being passed", func() {
					BeforeEach(func() {
						fakeServiceBroker.VolumeMounts = []brokerapi.VolumeMount{{
							Driver:       "driver",
							ContainerDir: "/dev/null",
							Mode:         "rw",
							DeviceType:   "shared",
							Device: brokerapi.SharedDevice{
								VolumeId:    "some-guid",
								MountConfig: map[string]interface{}{"key": "value"},
							},
						}}
					})

					Context("when the broker API version is greater than 2.9", func() {
						It("responds with a volume mount", func() {
							response := makeBindingRequest(uniqueInstanceID(), uniqueBindingID(), details)
							Expect(response.Body).To(MatchJSON(fixture("binding_with_volume_mounts.json")))
						})
					})

					Context("when the broker API version is 2.9", func() {
						It("responds with an experimental volume mount", func() {
							response := makeBindingRequestWithSpecificAPIVersion(uniqueInstanceID(), uniqueBindingID(), details, "2.9")
							Expect(response.Body).To(MatchJSON(fixture("binding_with_experimental_volume_mounts.json")))
						})
					})

					Context("when the broker API version is 2.8", func() {
						It("responds with an experimental volume mount", func() {
							response := makeBindingRequestWithSpecificAPIVersion(uniqueInstanceID(), uniqueBindingID(), details, "2.8")
							Expect(response.Body).To(MatchJSON(fixture("binding_with_experimental_volume_mounts.json")))
						})
					})
				})

				Context("when no bind details are being passed", func() {
					It("returns a 422", func() {
						response := makeBindingRequest(uniqueInstanceID(), uniqueBindingID(), nil)
						Expect(response.StatusCode).To(Equal(http.StatusUnprocessableEntity))
					})
				})

				Context("when there are arbitrary params", func() {
					var (
						rawParams string
						rawCtx    string
					)

					BeforeEach(func() {
						details["parameters"] = map[string]interface{}{
							"string": "some-string",
							"number": 1,
							"object": struct{ Name string }{"some-name"},
							"array":  []interface{}{"a", "b", "c"},
						}

						details["context"] = map[string]interface{}{
							"platform":      "fake-platform",
							"serial-number": 12648430,
							"object":        struct{ Name string }{"parameter"},
							"array":         []interface{}{"1", "2", "3"},
						}

						rawParams = `{
							"string":"some-string",
							"number":1,
							"object": { "Name": "some-name" },
							"array": [ "a", "b", "c" ]
						}`
						rawCtx = `{
							"platform":"fake-platform",
							"serial-number":12648430,
							"object": {"Name":"parameter"},
							"array":[ "1", "2", "3" ]
						}`
					})

					It("calls Bind on the service broker with all params", func() {
						makeBindingRequest(instanceID, bindingID, details)
						Expect(string(fakeServiceBroker.BoundBindingDetails.RawParameters)).To(MatchJSON(rawParams))
					})

					It("calls Bind with details with raw parameters", func() {
						makeBindingRequest(instanceID, bindingID, details)
						detailsWithRawParameters := brokerapi.DetailsWithRawParameters(fakeServiceBroker.BoundBindingDetails)
						rawParameters := detailsWithRawParameters.GetRawParameters()
						Expect(string(rawParameters)).To(MatchJSON(rawParams))
					})

					It("calls Bind with details with raw context", func() {
						makeBindingRequest(instanceID, bindingID, details)
						detailsWithRawContext := brokerapi.DetailsWithRawContext(fakeServiceBroker.BoundBindingDetails)
						rawContext := detailsWithRawContext.GetRawContext()
						Expect(string(rawContext)).To(MatchJSON(rawCtx))
					})
				})

				When("there are details in the bind_resource", func() {

					It("calls Bind on the service broker with the bind_resource", func() {

						details["bind_resource"] = map[string]interface{}{
							"app_guid": "a-guid",
							"space_guid": "a-space-guid",
							"route": "route.cf-apps.com",
							"credential_client_id": "some-credentials",
						}

						makeBindingRequest(instanceID, bindingID, details)
						Expect(fakeServiceBroker.BoundBindingDetails.BindResource).NotTo(BeNil())
						Expect(fakeServiceBroker.BoundBindingDetails.BindResource.AppGuid).To(Equal("a-guid"))
						Expect(fakeServiceBroker.BoundBindingDetails.BindResource.SpaceGuid).To(Equal("a-space-guid"))
						Expect(fakeServiceBroker.BoundBindingDetails.BindResource.Route).To(Equal("route.cf-apps.com"))
						Expect(fakeServiceBroker.BoundBindingDetails.BindResource.CredentialClientID).To(Equal("some-credentials"))
					})
				})

				When("there are no details in the bind_resource", func() {

					It("calls Bind on the service broker with an empty bind_resource", func() {

						details["bind_resource"] = map[string]interface{}{}

						makeBindingRequest(instanceID, bindingID, details)
						Expect(fakeServiceBroker.BoundBindingDetails.BindResource).NotTo(BeNil())
						Expect(fakeServiceBroker.BoundBindingDetails.BindResource.AppGuid).To(BeEmpty())
						Expect(fakeServiceBroker.BoundBindingDetails.BindResource.SpaceGuid).To(BeEmpty())
						Expect(fakeServiceBroker.BoundBindingDetails.BindResource.Route).To(BeEmpty())
						Expect(fakeServiceBroker.BoundBindingDetails.BindResource.CredentialClientID).To(BeEmpty())
					})
				})

			})

			Context("when the associated instance does not exist", func() {
				var instanceID string

				BeforeEach(func() {
					fakeServiceBroker.BindError = brokerapi.ErrInstanceDoesNotExist
				})

				It("returns a 404", func() {
					response := makeBindingRequest(uniqueInstanceID(), uniqueBindingID(), details)
					Expect(response.StatusCode).To(Equal(404))
				})

				It("returns an error JSON object", func() {
					response := makeBindingRequest(uniqueInstanceID(), uniqueBindingID(), details)
					Expect(response.Body).To(MatchJSON(`{"description":"instance does not exist"}`))
				})

				It("logs an appropriate error", func() {
					instanceID = uniqueInstanceID()
					makeBindingRequest(instanceID, uniqueBindingID(), details)
					Expect(lastLogLine().Message).To(ContainSubstring(".bind.instance-missing"))
					Expect(lastLogLine().Data["error"]).To(ContainSubstring("instance does not exist"))
				})
			})

			Context("when the requested binding already exists", func() {
				var instanceID string

				BeforeEach(func() {
					fakeServiceBroker.BindError = brokerapi.ErrBindingAlreadyExists
				})

				It("returns a 409", func() {
					response := makeBindingRequest(uniqueInstanceID(), uniqueBindingID(), details)
					Expect(response.StatusCode).To(Equal(409))
				})

				It("returns an error JSON object", func() {
					response := makeBindingRequest(uniqueInstanceID(), uniqueBindingID(), details)
					Expect(response.Body).To(MatchJSON(`{"description":"binding already exists"}`))
				})

				It("logs an appropriate error", func() {
					instanceID = uniqueInstanceID()
					makeBindingRequest(instanceID, uniqueBindingID(), details)
					makeBindingRequest(instanceID, uniqueBindingID(), details)

					Expect(lastLogLine().Message).To(ContainSubstring(".bind.binding-already-exists"))
					Expect(lastLogLine().Data["error"]).To(ContainSubstring("binding already exists"))
				})
			})

			Context("when the binding returns an unknown error", func() {
				BeforeEach(func() {
					fakeServiceBroker.BindError = errors.New("unknown error")
				})

				It("returns a generic 500 error response", func() {
					response := makeBindingRequest(uniqueInstanceID(), uniqueBindingID(), details)
					Expect(response.StatusCode).To(Equal(500))
					Expect(response.Body).To(MatchJSON(`{"description":"unknown error"}`))
				})

				It("logs a detailed error message", func() {
					makeBindingRequest(uniqueInstanceID(), uniqueBindingID(), details)

					Expect(lastLogLine().Message).To(ContainSubstring(".bind.unknown-error"))
					Expect(lastLogLine().Data["error"]).To(ContainSubstring("unknown error"))
				})
			})

			Context("when the binding returns a custom error", func() {
				BeforeEach(func() {
					fakeServiceBroker.BindError = brokerapi.NewFailureResponse(
						errors.New("I failed in unique and interesting ways"),
						http.StatusTeapot,
						"interesting-failure",
					)
				})

				It("returns status teapot", func() {
					response := makeBindingRequest(uniqueInstanceID(), uniqueBindingID(), details)
					Expect(response.StatusCode).To(Equal(http.StatusTeapot))
				})

				It("returns json with a description field and a useful error message", func() {
					response := makeBindingRequest(uniqueInstanceID(), uniqueBindingID(), details)
					Expect(response.Body).To(MatchJSON(`{"description":"I failed in unique and interesting ways"}`))
				})

				It("logs an appropriate error", func() {
					makeBindingRequest(uniqueInstanceID(), uniqueBindingID(), details)
					Expect(lastLogLine().Message).To(ContainSubstring(".bind.interesting-failure"))
					Expect(lastLogLine().Data["error"]).To(ContainSubstring("I failed in unique and interesting ways"))
				})
			})
		})

		Describe("unbinding", func() {
			makeUnbindingRequestWithServiceIDPlanID := func(instanceID, bindingID, serviceID, planID, apiVersion string) *testflight.Response {
				response := &testflight.Response{}
				testflight.WithServer(brokerAPI, func(r *testflight.Requester) {
					path := fmt.Sprintf("/v2/service_instances/%s/service_bindings/%s?plan_id=%s&service_id=%s",
						instanceID, bindingID, planID, serviceID)
					request, _ := http.NewRequest("DELETE", path, strings.NewReader(""))
					request.Header.Add("Content-Type", "application/json")
					request.Header.Add("X-Broker-API-Version", apiVersion)
					request.SetBasicAuth("username", "password")

					response = r.Do(request)
				})
				return response
			}

			makeUnbindingRequest := func(instanceID string, bindingID string) *testflight.Response {
				return makeUnbindingRequestWithServiceIDPlanID(instanceID, bindingID, "service-id", "plan-id", "2.13")
			}

			Context("when the associated instance exists", func() {
				var instanceID string
				var provisionDetails map[string]interface{}

				BeforeEach(func() {
					instanceID = uniqueInstanceID()
					provisionDetails = map[string]interface{}{
						"service_id":        fakeServiceBroker.ServiceID,
						"plan_id":           "plan-id",
						"organization_guid": "organization-guid",
						"space_guid":        "space-guid",
					}
					makeInstanceProvisioningRequest(instanceID, provisionDetails, "")
				})

				Context("the request is malformed", func() {
					var bindingID string

					BeforeEach(func() {
						bindingID = uniqueBindingID()
						makeBindingRequest(instanceID, bindingID, map[string]interface{}{})
					})

					It("missing header X-Broker-API-Version", func() {
						response := makeUnbindingRequestWithServiceIDPlanID(instanceID, bindingID, "service-id", "plan-id", "")
						Expect(response.StatusCode).To(Equal(412))
						Expect(lastLogLine().Message).To(ContainSubstring(".unbind.broker-api-version-invalid"))
						Expect(lastLogLine().Data["error"]).To(ContainSubstring("X-Broker-API-Version Header not set"))
					})

					It("has wrong version of API", func() {
						response := makeUnbindingRequestWithServiceIDPlanID(instanceID, bindingID, "service-id", "plan-id", "1.1")
						Expect(response.StatusCode).To(Equal(412))
						Expect(lastLogLine().Message).To(ContainSubstring(".unbind.broker-api-version-invalid"))
						Expect(lastLogLine().Data["error"]).To(ContainSubstring("X-Broker-API-Version Header must be 2.x"))
					})

					It("missing service-id", func() {
						response := makeUnbindingRequestWithServiceIDPlanID(instanceID, bindingID, "", "plan-id", "2.13")
						Expect(response.StatusCode).To(Equal(400))
						Expect(lastLogLine().Message).To(ContainSubstring(".unbind.service-id-missing"))
						Expect(lastLogLine().Data["error"]).To(ContainSubstring("service_id missing"))
					})

					It("missing plan-id", func() {
						response := makeUnbindingRequestWithServiceIDPlanID(instanceID, bindingID, "service-id", "", "2.13")
						Expect(response.StatusCode).To(Equal(400))
						Expect(lastLogLine().Message).To(ContainSubstring(".unbind.plan-id-missing"))
						Expect(lastLogLine().Data["error"]).To(ContainSubstring("plan_id missing"))
					})
				})

				Context("and the binding exists", func() {
					var bindingID string

					BeforeEach(func() {
						bindingID = uniqueBindingID()
						makeBindingRequest(instanceID, bindingID, map[string]interface{}{
							"service_id": "service_id", "plan_id": "plan_id",
						})
					})

					It("returns a 200", func() {
						response := makeUnbindingRequest(instanceID, bindingID)
						Expect(response.StatusCode).To(Equal(200))
					})

					It("returns an empty JSON object", func() {
						response := makeUnbindingRequest(instanceID, bindingID)
						Expect(response.Body).To(MatchJSON(`{}`))
					})

					It("contains plan_id", func() {
						makeUnbindingRequest(instanceID, bindingID)
						Expect(fakeServiceBroker.UnbindingDetails.PlanID).To(Equal("plan-id"))
					})

					It("contains service_id", func() {
						makeUnbindingRequest(instanceID, bindingID)
						Expect(fakeServiceBroker.UnbindingDetails.ServiceID).To(Equal("service-id"))
					})
				})

				Context("but the binding does not exist", func() {
					It("returns a 410", func() {
						response := makeUnbindingRequest(instanceID, "does-not-exist")
						Expect(response.StatusCode).To(Equal(410))
					})

					It("logs an appropriate error message", func() {
						makeUnbindingRequest(instanceID, "does-not-exist")

						Expect(lastLogLine().Message).To(ContainSubstring(".unbind.binding-missing"))
						Expect(lastLogLine().Data["error"]).To(ContainSubstring("binding does not exist"))
					})

					It("returns an empty JSON object", func() {
						response := makeUnbindingRequest(instanceID, "does-not-exist")
						Expect(response.Body).To(MatchJSON(`{}`))
					})
				})
			})

			Context("when the associated instance does not exist", func() {
				var instanceID string

				It("returns a 410", func() {
					response := makeUnbindingRequest(uniqueInstanceID(), uniqueBindingID())
					Expect(response.StatusCode).To(Equal(http.StatusGone))
				})

				It("returns an empty JSON object", func() {
					response := makeUnbindingRequest(uniqueInstanceID(), uniqueBindingID())
					Expect(response.Body).To(MatchJSON(`{}`))
				})

				It("logs an appropriate error", func() {
					instanceID = uniqueInstanceID()
					makeUnbindingRequest(instanceID, uniqueBindingID())

					Expect(lastLogLine().Message).To(ContainSubstring(".unbind.instance-missing"))
					Expect(lastLogLine().Data["error"]).To(ContainSubstring("instance does not exist"))
				})
			})

			Context("when unbinding returns an unknown error", func() {
				BeforeEach(func() {
					fakeServiceBroker.UnbindError = errors.New("unknown error")
				})

				It("returns a generic 500 error response", func() {
					response := makeUnbindingRequest(uniqueInstanceID(), uniqueBindingID())
					Expect(response.StatusCode).To(Equal(500))
					Expect(response.Body).To(MatchJSON(`{"description":"unknown error"}`))
				})

				It("logs a detailed error message", func() {
					makeUnbindingRequest(uniqueInstanceID(), uniqueBindingID())

					Expect(lastLogLine().Message).To(ContainSubstring(".unbind.unknown-error"))
					Expect(lastLogLine().Data["error"]).To(ContainSubstring("unknown error"))
				})
			})

			Context("when unbinding returns a custom error", func() {
				BeforeEach(func() {
					fakeServiceBroker.UnbindError = brokerapi.NewFailureResponse(
						errors.New("I failed in unique and interesting ways"),
						http.StatusTeapot,
						"interesting-failure",
					)
				})

				It("returns status teapot", func() {
					response := makeUnbindingRequest(uniqueInstanceID(), uniqueBindingID())
					Expect(response.StatusCode).To(Equal(http.StatusTeapot))
				})

				It("returns json with a description field and a useful error message", func() {
					response := makeUnbindingRequest(uniqueInstanceID(), uniqueBindingID())
					Expect(response.Body).To(MatchJSON(`{"description":"I failed in unique and interesting ways"}`))
				})

				It("logs an appropriate error", func() {
					makeUnbindingRequest(uniqueInstanceID(), uniqueBindingID())
					Expect(lastLogLine().Message).To(ContainSubstring(".unbind.interesting-failure"))
					Expect(lastLogLine().Data["error"]).To(ContainSubstring("I failed in unique and interesting ways"))
				})
			})
		})

		Describe("last_operation", func() {
			makeLastOperationRequest := func(instanceID, operationData, apiVersion string) *testflight.Response {
				response := &testflight.Response{}
				testflight.WithServer(brokerAPI, func(r *testflight.Requester) {
					path := fmt.Sprintf("/v2/service_instances/%s/last_operation", instanceID)
					if operationData != "" {
						path = fmt.Sprintf("%s?operation=%s", path, url.QueryEscape(operationData))
					}

					request, _ := http.NewRequest("GET", path, strings.NewReader(""))
					if apiVersion != "" {
						request.Header.Add("X-Broker-API-Version", apiVersion)
					}
					request.Header.Add("Content-Type", "application/json")
					request.SetBasicAuth("username", "password")

					response = r.Do(request)
				})
				return response
			}

			It("calls the broker with the relevant instance ID", func() {
				instanceID := "instanceID"
				makeLastOperationRequest(instanceID, "", "2.14")
				Expect(fakeServiceBroker.LastOperationInstanceID).To(Equal(instanceID))
			})

			It("calls the broker with the URL decoded operation data", func() {
				instanceID := "an-instance"
				operationData := `{"foo":"bar"}`
				makeLastOperationRequest(instanceID, operationData, "2.14")
				Expect(fakeServiceBroker.LastOperationData).To(Equal(operationData))
			})

			It("should return succeeded if the operation completed successfully", func() {
				fakeServiceBroker.LastOperationState = "succeeded"
				fakeServiceBroker.LastOperationDescription = "some description"

				instanceID := "instanceID"
				response := makeLastOperationRequest(instanceID, "", "2.14")

				logs := brokerLogger.Logs()

				Expect(logs[0].Message).To(ContainSubstring(".lastOperation.starting-check-for-operation"))
				Expect(logs[0].Data["instance-id"]).To(ContainSubstring(instanceID))

				Expect(logs[1].Message).To(ContainSubstring(".lastOperation.done-check-for-operation"))
				Expect(logs[1].Data["instance-id"]).To(ContainSubstring(instanceID))
				Expect(logs[1].Data["state"]).To(ContainSubstring(string(fakeServiceBroker.LastOperationState)))

				Expect(response.StatusCode).To(Equal(200))
				Expect(response.Body).To(MatchJSON(fixture("last_operation_succeeded.json")))
			})

			It("should return a 410 and log in case the instance id is not found", func() {
				fakeServiceBroker.LastOperationError = brokerapi.ErrInstanceDoesNotExist
				instanceID := "non-existing"
				response := makeLastOperationRequest(instanceID, "", "2.14")

				Expect(lastLogLine().Message).To(ContainSubstring(".lastOperation.instance-missing"))
				Expect(lastLogLine().Data["error"]).To(ContainSubstring("instance does not exist"))

				Expect(response.StatusCode).To(Equal(410))
				Expect(response.Body).To(MatchJSON(`{}`))
			})

			Context("when last_operation returns an unknown error", func() {
				BeforeEach(func() {
					fakeServiceBroker.LastOperationError = errors.New("unknown error")
				})

				It("returns a generic 500 error response", func() {
					response := makeLastOperationRequest("instanceID", "", "2.14")

					Expect(response.StatusCode).To(Equal(500))
					Expect(response.Body).To(MatchJSON(`{"description": "unknown error"}`))
				})

				It("logs a detailed error message", func() {
					makeLastOperationRequest("instanceID", "", "2.14")

					Expect(lastLogLine().Message).To(ContainSubstring(".lastOperation.unknown-error"))
					Expect(lastLogLine().Data["error"]).To(ContainSubstring("unknown error"))
				})
			})

			Context("when last_operation returns a custom error", func() {
				BeforeEach(func() {
					fakeServiceBroker.LastOperationError = brokerapi.NewFailureResponse(
						errors.New("I failed in unique and interesting ways"),
						http.StatusTeapot,
						"interesting-failure",
					)
				})

				It("returns status teapot", func() {
					response := makeLastOperationRequest("instanceID", "", "2.14")
					Expect(response.StatusCode).To(Equal(http.StatusTeapot))
				})

				It("returns json with a description field and a useful error message", func() {
					response := makeLastOperationRequest("instanceID", "", "2.14")
					Expect(response.Body).To(MatchJSON(`{"description":"I failed in unique and interesting ways"}`))
				})

				It("logs an appropriate error", func() {
					makeLastOperationRequest("instanceID", "", "2.14")
					Expect(lastLogLine().Message).To(ContainSubstring(".lastOperation.interesting-failure"))
					Expect(lastLogLine().Data["error"]).To(ContainSubstring("I failed in unique and interesting ways"))
				})
			})

			Context("the request is malformed", func() {
				It("missing header X-Broker-API-Version", func() {
					response := makeLastOperationRequest("instance-id", "", "")
					Expect(response.StatusCode).To(Equal(412))
					Expect(lastLogLine().Message).To(ContainSubstring(".lastOperation.broker-api-version-invalid"))
					Expect(lastLogLine().Data["error"]).To(ContainSubstring("X-Broker-API-Version Header not set"))
				})

				It("has wrong version of API", func() {
					response := makeLastOperationRequest("instance-id", "", "1.2")
					Expect(response.StatusCode).To(Equal(412))
					Expect(lastLogLine().Message).To(ContainSubstring(".lastOperation.broker-api-version-invalid"))
					Expect(lastLogLine().Data["error"]).To(ContainSubstring("X-Broker-API-Version Header must be 2.x"))
				})
			})
		})
	})
})
