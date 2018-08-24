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

	"code.cloudfoundry.org/lager"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/models"
	"github.com/GoogleCloudPlatform/gcp-service-broker/db_service"
	"github.com/GoogleCloudPlatform/gcp-service-broker/utils"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "plan-info",
		Short: "Dump plan information from the database",
		Long:  `Dump plan information from the database.`,
		Run: func(cmd *cobra.Command, args []string) {
			logger := lager.NewLogger("get_plan_info_cmd")
			db := db_service.SetupDb(logger)

			var pds []*models.PlanDetailsV1
			if err := db.Find(&pds).Error; err != nil {
				fmt.Errorf("Could not retrieve plan details rows from db: %s", err)
			}

			utils.PrettyPrintOrExit(pds)
		},
	})
}
