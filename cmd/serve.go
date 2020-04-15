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
	"database/sql"
	"net/http"

	"code.cloudfoundry.org/lager"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers"
	"github.com/GoogleCloudPlatform/gcp-service-broker/db_service"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/broker"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/brokerpak"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/providers/builtin"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/server"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/toggles"
	"github.com/GoogleCloudPlatform/gcp-service-broker/utils"
	"github.com/gorilla/mux"
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
	across spaces in Tanzu.`)

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

	rootCmd.AddCommand(&cobra.Command{
		Use:   "serve-docs",
		Short: "Just serve the docs normally available on the broker",
		Run: func(cmd *cobra.Command, args []string) {
			serveDocs()
		},
	})

	viper.BindEnv(apiUserProp, "SECURITY_USER_NAME")
	viper.BindEnv(apiPasswordProp, "SECURITY_USER_PASSWORD")
	viper.BindEnv(apiPortProp, "PORT")
}

func serve() {
	logger := utils.NewLogger("gcp-service-broker")
	db := db_service.New(logger)

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

	credentials := brokerapi.BrokerCredentials{
		Username: viper.GetString(apiUserProp),
		Password: viper.GetString(apiPasswordProp),
	}

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

	startServer(cfg.Registry, db.DB(), brokerAPI)
}

func serveDocs() {
	logger := utils.NewLogger("gcp-service-broker")
	// init broker
	registry := builtin.BuiltinBrokerRegistry()
	if err := brokerpak.RegisterAll(registry); err != nil {
		logger.Error("loading brokerpaks", err)
	}

	startServer(registry, nil, nil)
}

func startServer(registry broker.BrokerRegistry, db *sql.DB, brokerapi http.Handler) {
	logger := utils.NewLogger("gcp-service-broker")

	router := mux.NewRouter()

	// match paths going to the brokerapi first
	if brokerapi != nil {
		router.PathPrefix("/v2").Handler(brokerapi)
	}

	server.AddDocsHandler(router, registry)
	router.HandleFunc("/examples", server.NewExampleHandler(registry))
	server.AddHealthHandler(router, db)

	port := viper.GetString(apiPortProp)
	logger.Info("Serving", lager.Data{"port": port})
	http.ListenAndServe(":"+port, router)
}
