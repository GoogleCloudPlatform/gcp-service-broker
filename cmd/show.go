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
	"fmt"
	"log"

	"code.cloudfoundry.org/lager"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/models"
	"github.com/GoogleCloudPlatform/gcp-service-broker/db_service"
	"github.com/GoogleCloudPlatform/gcp-service-broker/utils"
	"github.com/spf13/cobra"
)

func init() {
	showCmd := &cobra.Command{
		Use:   "show",
		Short: "Show info about the provisioned resources",
		Long:  `Show info about the provisioned resources`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	rootCmd.AddCommand(showCmd)

	addDumpTableCommand(showCmd, "bindings", &[]models.ServiceBindingCredentials{})
	addDumpTableCommand(showCmd, "instances", &[]models.ServiceInstanceDetails{})
	addDumpTableCommand(showCmd, "migrations", &[]models.Migration{})
	addDumpTableCommand(showCmd, "operations", &[]models.CloudOperation{})
	addDumpTableCommand(showCmd, "provisions", &[]models.ProvisionRequestDetails{})
}

func addDumpTableCommand(parent *cobra.Command, name string, value interface{}) {
	tmp := &cobra.Command{
		Use:   name,
		Short: fmt.Sprintf("Show the %s table as JSON", name),
		Run: func(cmd *cobra.Command, args []string) {
			tableToJson(value)
		},
	}

	parent.AddCommand(tmp)
}

func tableToJson(results interface{}) {
	logger := lager.NewLogger("show-command")
	db := db_service.SetupDb(logger)

	if err := db.Find(results).Error; err != nil {
		log.Fatal(err)
	}

	utils.PrettyPrintOrExit(results)
}
