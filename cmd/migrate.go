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
	"github.com/GoogleCloudPlatform/gcp-service-broker/db_service"
	"github.com/GoogleCloudPlatform/gcp-service-broker/utils"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "migrate",
		Short: "Upgrade your database",
		Long:  `Upgrade your database to be compatible with this service broker.`,
		Run: func(cmd *cobra.Command, args []string) {
			logger := utils.NewLogger("migrations-cmd")

			logger.Debug("Setting up the database")
			db := db_service.SetupDb(logger)

			logger.Debug("Starting the migration")
			if err := db_service.RunMigrations(db); err != nil {
				logger.Fatal("Error running migrations", err)
			}

			logger.Debug("Finished migration")
		},
	})
}
