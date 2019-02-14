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

package cmd

import (
	"context"
	"net/http"

	"code.cloudfoundry.org/lager"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers"
	"github.com/GoogleCloudPlatform/gcp-service-broker/db_service"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/server"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/toggles"
	"github.com/GoogleCloudPlatform/gcp-service-broker/utils"
	"github.com/pivotal-cf/brokerapi"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	apiUserProp     = "api.user"
	apiPasswordProp = "api.password"
	apiPortProp     = "api.port"
)

var cfCompatibilityToggle = toggles.Features.Toggle("enable-cf-sharing", false, `Set all services to have the Sharable flag so they can be shared
	across spaces in PCF.`)

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
	logger := utils.NewLogger("gcp-service-broker")
	db_service.New(logger)

	// init broker
	cfg, err := brokers.NewBrokerConfigFromEnv()
	if err != nil {
		logger.Fatal("Error initializing service broker config: %s", err)
	}
	var serviceBroker brokerapi.ServiceBroker
	serviceBroker, err = brokers.New(cfg, logger)
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

	if cfCompatibilityToggle.IsActive() {
		logger.Info("Enabling Cloud Foundry service sharing")
		serviceBroker = server.NewCfSharingWrapper(serviceBroker)
	}

	services, err := serviceBroker.Services(context.Background())
	if err != nil {
		logger.Error("creating service catalog", err)
	}
	logger.Info("service catalog", lager.Data{"catalog": services})

	brokerAPI := brokerapi.New(serviceBroker, logger, credentials)
	http.Handle("/", brokerAPI)

	docsHandler, err := server.NewDocsHandler(cfg.Registry)
	if err != nil {
		logger.Error("creating docs endpoint", err)
	}
	http.Handle("/docs", docsHandler)

	configHandler, err := server.NewServiceConfigHandler(cfg.Registry)
	if err != nil {
		logger.Error("creating service config endpoint", err)
	}
	http.Handle("/service-config", configHandler)

	http.ListenAndServe(":"+port, nil)
}
