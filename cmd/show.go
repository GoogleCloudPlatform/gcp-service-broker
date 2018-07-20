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

package cmd

import (
	"encoding/json"
	"fmt"

	"code.cloudfoundry.org/lager"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/models"
	"github.com/GoogleCloudPlatform/gcp-service-broker/db_service"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(showCmd)
	showCmd.AddCommand(showMigrationsCmd)
	showCmd.AddCommand(showBindingsCmd)
	showCmd.AddCommand(showInstancesCmd)
	showCmd.AddCommand(showProvisionsCmd)
}

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show info about the provisioned resources",
	Long:  `Show info about the provisioned resources`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var showMigrationsCmd = &cobra.Command{
	Use:   "migrations",
	Short: "Show info about the migrations",
	RunE: func(cmd *cobra.Command, args []string) error {
		var results []models.Migration
		return tableToJson(&results)
	},
}

var showBindingsCmd = &cobra.Command{
	Use:   "bindings",
	Short: "Show info about the bindings",
	RunE: func(cmd *cobra.Command, args []string) error {
		var results []models.ServiceBindingCredentials
		return tableToJson(&results)
	},
}

var showInstancesCmd = &cobra.Command{
	Use:   "instances",
	Short: "Show info about the service instances",
	RunE: func(cmd *cobra.Command, args []string) error {
		var results []models.ServiceInstanceDetails
		return tableToJson(&results)
	},
}

var showProvisionsCmd = &cobra.Command{
	Use:   "provisions",
	Short: "Show info about the service provision requests",
	RunE: func(cmd *cobra.Command, args []string) error {
		var results []models.ProvisionRequestDetails
		return tableToJson(&results)
	},
}

func tableToJson(results interface{}) error {
	logger := lager.NewLogger("show-command")
	db := db_service.SetupDb(logger)

	err := db.Find(results).Error

	if err != nil {
		return err
	}

	res, err := json.MarshalIndent(results, "", "    ")
	if err != nil {
		return err
	}

	fmt.Println(string(res))
	return nil
}
