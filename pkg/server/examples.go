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
	"sort"
	"strings"

	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/broker"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/client"
)

func GetAllCompleteServiceExamples(registry broker.BrokerRegistry) ([]client.CompleteServiceExample, error) {

	var allExamples []client.CompleteServiceExample

	services := registry.GetAllServices()

	for _, service := range services {

		serviceExamples, err := client.GetExamplesForAService(service)

		if err != nil {
			return nil, err
		}

		allExamples = append(allExamples, serviceExamples...)
	}

	// Sort by ServiceName and ExampleName so there's a consistent order in the UI and tests.
	sort.Slice(allExamples, func(i int, j int) bool {
		if strings.Compare(allExamples[i].ServiceName, allExamples[j].ServiceName) != 0 {
			return allExamples[i].ServiceName < allExamples[j].ServiceName
		} else {
			return allExamples[i].ServiceExample.Name < allExamples[j].ServiceExample.Name
		}
	})

	return allExamples, nil
}

func NewExampleHandler(registry broker.BrokerRegistry) http.HandlerFunc {
	allExamples, err := GetAllCompleteServiceExamples(registry)

	if err != nil {
		return func(w http.ResponseWriter, rep *http.Request) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}

	exampleJSON, err := json.Marshal(allExamples)

	if err != nil {
		return func(w http.ResponseWriter, rep *http.Request) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}

	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(exampleJSON)
	}
}
