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
	"github.com/pivotal-cf/brokerapi"

	"errors"

	"net/http"

	"code.cloudfoundry.org/lager"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

var _ = Describe("FailureResponse", func() {
	Describe("ErrorResponse", func() {
		It("returns a ErrorResponse containing the error message", func() {
			failureResponse := brokerapi.NewFailureResponse(errors.New("my error message"), http.StatusForbidden, "log-key")
			Expect(failureResponse.ErrorResponse()).To(Equal(brokerapi.ErrorResponse{
				Description: "my error message",
			}))
		})

		Context("when the error key is provided", func() {
			It("returns a ErrorResponse containing the error message and the error key", func() {
				failureResponse := brokerapi.NewFailureResponseBuilder(errors.New("my error message"), http.StatusForbidden, "log-key").WithErrorKey("error key").Build()
				Expect(failureResponse.ErrorResponse()).To(Equal(brokerapi.ErrorResponse{
					Description: "my error message",
					Error:       "error key",
				}))
			})
		})

		Context("when created with empty response", func() {
			It("returns an EmptyResponse", func() {
				failureResponse := brokerapi.NewFailureResponseBuilder(errors.New("my error message"), http.StatusForbidden, "log-key").WithEmptyResponse().Build()
				Expect(failureResponse.ErrorResponse()).To(Equal(brokerapi.EmptyResponse{}))
			})
		})
	})

	Describe("AppendErrorMessage", func() {
		It("returns the error with the additional error message included, with a non-empty body", func() {
			failureResponse := brokerapi.NewFailureResponseBuilder(errors.New("my error message"), http.StatusForbidden, "log-key").WithErrorKey("some-key").Build()
			Expect(failureResponse.Error()).To(Equal("my error message"))

			newError := failureResponse.AppendErrorMessage("and some more details")

			Expect(newError.Error()).To(Equal("my error message and some more details"))
			Expect(newError.ValidatedStatusCode(nil)).To(Equal(http.StatusForbidden))
			Expect(newError.LoggerAction()).To(Equal(failureResponse.LoggerAction()))
			
			errorResponse, typeCast := newError.ErrorResponse().(brokerapi.ErrorResponse)
			Expect(typeCast).To(BeTrue())
			Expect(errorResponse.Error).To(Equal("some-key"))
			Expect(errorResponse.Description).To(Equal("my error message and some more details"))
		})

		It("returns the error with the additional error message included, with an empty body", func() {
			failureResponse := brokerapi.NewFailureResponseBuilder(errors.New("my error message"), http.StatusForbidden, "log-key").WithEmptyResponse().Build()
			Expect(failureResponse.Error()).To(Equal("my error message"))

			newError := failureResponse.AppendErrorMessage("and some more details")

			Expect(newError.Error()).To(Equal("my error message and some more details"))
			Expect(newError.ValidatedStatusCode(nil)).To(Equal(http.StatusForbidden))
			Expect(newError.LoggerAction()).To(Equal(failureResponse.LoggerAction()))
			Expect(newError.ErrorResponse()).To(Equal(failureResponse.ErrorResponse()))
		})
	})

	Describe("ValidatedStatusCode", func() {
		It("returns the status code that was passed in", func() {
			failureResponse := brokerapi.NewFailureResponse(errors.New("my error message"), http.StatusForbidden, "log-key")
			Expect(failureResponse.ValidatedStatusCode(nil)).To(Equal(http.StatusForbidden))
		})

		It("when error key is provided it returns the status code that was passed in", func() {
			failureResponse := brokerapi.NewFailureResponseBuilder(errors.New("my error message"), http.StatusForbidden, "log-key").WithErrorKey("error key").Build()
			Expect(failureResponse.ValidatedStatusCode(nil)).To(Equal(http.StatusForbidden))
		})

		Context("when the status code is invalid", func() {
			It("returns 500", func() {
				failureResponse := brokerapi.NewFailureResponse(errors.New("my error message"), 600, "log-key")
				Expect(failureResponse.ValidatedStatusCode(nil)).To(Equal(http.StatusInternalServerError))
			})

			It("logs that the status has been changed", func() {
				log := gbytes.NewBuffer()
				logger := lager.NewLogger("test")
				logger.RegisterSink(lager.NewWriterSink(log, lager.DEBUG))
				failureResponse := brokerapi.NewFailureResponse(errors.New("my error message"), 600, "log-key")
				failureResponse.ValidatedStatusCode(logger)
				Expect(log).To(gbytes.Say("Invalid failure http response code: 600, expected 4xx or 5xx, returning internal server error: 500."))
			})
		})
	})

	Describe("LoggerAction", func() {
		It("returns the logger action that was passed in", func() {
			failureResponse := brokerapi.NewFailureResponseBuilder(errors.New("my error message"), http.StatusForbidden, "log-key").WithErrorKey("error key").Build()
			Expect(failureResponse.LoggerAction()).To(Equal("log-key"))
		})

		It("when error key is provided it returns the logger action that was passed in", func() {
			failureResponse := brokerapi.NewFailureResponse(errors.New("my error message"), http.StatusForbidden, "log-key")
			Expect(failureResponse.LoggerAction()).To(Equal("log-key"))
		})
	})
})
