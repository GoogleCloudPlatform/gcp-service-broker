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

package server

import (
	"encoding/json"
	"net/http"

	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/broker"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/client"
)

func NewExampleHandler(registry broker.BrokerRegistry) http.HandlerFunc {
	allExamples, err := client.GetAllCompleteServiceExamples(registry)
	if err != nil {
		return func(w http.ResponseWriter, rep *http.Request) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}

	json, err := json.Marshal(allExamples)

	if err != nil {
		return func(w http.ResponseWriter, rep *http.Request) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}

	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(json)
	}
}
