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

package cmd

import (
	"net/http"
	"os"

	"code.cloudfoundry.org/lager"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/config"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/models"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/name_generator"
	"github.com/GoogleCloudPlatform/gcp-service-broker/db_service"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	apiUserProp     = "api.user"
	apiPasswordProp = "api.password"
	apiPortProp     = "api.port"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "serve",
		Short: "Start the service broker",
		Long: `Starts the service broker listening on a port defined by the
	PORT environment variable.`,
		Run: func(cmd *cobra.Command, args []string) {
			serve()
		},
	})

	viper.BindEnv(apiUserProp, "SECURITY_USER_NAME")
	viper.BindEnv(apiPasswordProp, "SECURITY_USER_PASSWORD")
	viper.BindEnv(apiPortProp, "PORT")
}

func serve() {

	logger := lager.NewLogger("gcp-service-broker")
	logger.RegisterSink(lager.NewWriterSink(os.Stderr, lager.DEBUG))

	models.ProductionizeUserAgent()

	db_service.New(logger)
	name_generator.New()

	// init broker
	cfg, err := config.NewBrokerConfigFromEnv()
	if err != nil {
		logger.Fatal("Error initializing service broker config: %s", err)
	}
	serviceBroker, err := brokers.New(cfg, logger)
	if err != nil {
		logger.Fatal("Error initializing service broker: %s", err)
	}

	username := viper.GetString(apiUserProp)
	password := viper.GetString(apiPasswordProp)
	port := viper.GetString(apiPortProp)

	credentials := brokerapi.BrokerCredentials{
		Username: username,
		Password: password,
	}

	// init api
	logger.Info("Serving", lager.Data{
		"port":     port,
		"username": username,
	})

	brokerAPI := brokerapi.New(serviceBroker, logger, credentials)
	http.Handle("/", brokerAPI)
	http.ListenAndServe(":"+port, nil)
}
