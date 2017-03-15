// Copyright the Service Broker Project Authors.
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
//
////////////////////////////////////////////////////////////////////////////////
//

package main

import (
	"net/http"
	"os"

	"code.cloudfoundry.org/lager"
	"gcp-service-broker/brokerapi"
	"gcp-service-broker/brokerapi/brokers"
	"gcp-service-broker/brokerapi/brokers/name_generator"
	"gcp-service-broker/db_service"
	"gcp-service-broker/brokerapi/brokers/models"
)

func main() {
	// init logger

	logger := lager.NewLogger("my-service-broker")
	logger.RegisterSink(lager.NewWriterSink(os.Stderr, lager.DEBUG))

	models.ProductionizeUserAgent()

	db_service.New(logger)
	name_generator.New()

	// init broker
	serviceBroker, err := brokers.New(logger)
	if err != nil {
		logger.Fatal("Error initializing service broker: %s", err)
	}

	username := os.Getenv("SECURITY_USER_NAME")
	password := os.Getenv("SECURITY_USER_PASSWORD")

	credentials := brokerapi.BrokerCredentials{
		Username: username,
		Password: password,
	}

	// init api
	brokerAPI := brokerapi.New(serviceBroker, logger, credentials)
	http.Handle("/", brokerAPI)
	portEnvVar := os.Getenv("PORT")
	logger.Debug("starting application on " + portEnvVar)
	http.ListenAndServe(":"+portEnvVar, nil)
}
