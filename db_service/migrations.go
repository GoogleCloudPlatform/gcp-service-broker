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

package db_service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/models"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/broker"
	"github.com/GoogleCloudPlatform/gcp-service-broker/utils"
	"github.com/jinzhu/gorm"
	googlecloudsql "google.golang.org/api/sqladmin/v1beta4"
)

// runs schema migrations on the provided service broker database to get it up to date
func RunMigrations(db *gorm.DB) error {

	migrations := make([]func() error, 3)

	// initial migration - creates tables
	migrations[0] = func() error {
		return autoMigrateTables(db,
			&models.ServiceInstanceDetailsV1{},
			&models.ServiceBindingCredentialsV1{},
			&models.ProvisionRequestDetailsV1{},
			&models.PlanDetailsV1{},
			&models.MigrationV1{})
	}

	// adds CloudOperation table
	migrations[1] = func() error {

		if err := autoMigrateTables(db, &models.CloudOperationV1{}); err != nil {
			return err
		}

		// copy provision request details into service instance details
		cfg, err := utils.GetAuthedConfig()
		if err != nil {
			return fmt.Errorf("Error getting authorized http client: %s", err)
		}

		prs := []models.ProvisionRequestDetailsV1{}
		if err := db.Find(&prs).Error; err != nil {
			return err
		}

		if len(prs) == 0 {
			return nil
		}

		projectId, err := utils.GetDefaultProjectId()
		if err != nil {
			return fmt.Errorf("couldn't get Project ID for database upgrades %s", err)
		}

		for _, pr := range prs {
			var si models.ServiceInstanceDetailsV1
			if err := db.Where("id = ?", pr.ServiceInstanceId).First(&si).Error; err != nil {
				return err
			}
			od := make(map[string]string)
			if err := json.Unmarshal([]byte(pr.RequestDetails), &od); err != nil {
				return err
			}
			newOd := make(map[string]string)

			// cloudsql
			svc, err := broker.GetServiceById(si.ServiceId)
			if err != nil {
				return err
			}

			switch svc.Name {
			case models.CloudsqlMySQLName:
				newOd["instance_name"] = od["instance_name"]
				newOd["database_name"] = od["database_name"]

				sqlService, err := googlecloudsql.New(cfg.Client(context.Background()))
				if err != nil {
					return fmt.Errorf("Error creating new CloudSQL Client: %s", err)
				}
				dbService := googlecloudsql.NewInstancesService(sqlService)
				clouddb, err := dbService.Get(projectId, od["instance_name"]).Do()
				if err != nil {
					return fmt.Errorf("Error getting instance from api: %s", err)
				}
				newOd["host"] = clouddb.IpAddresses[0].IpAddress

			// bigquery
			case models.BigqueryName:
				newOd["dataset_id"] = od["name"]
			// ml apis
			case models.MlName:
				// n/a
			// storage
			case models.StorageName:
				newOd["bucket_name"] = od["name"]

			// pubsub
			case models.PubsubName:
				newOd["topic_name"] = od["topic_name"]
				newOd["subscription_name"] = od["subscription_name"]
			default:
				return fmt.Errorf("unrecognized service: %s", si.ServiceId)
			}

			odBytes, err := json.Marshal(&newOd)
			if err != nil {
				return err
			}
			si.OtherDetails = string(odBytes)
			if err := db.Save(&si).Error; err != nil {
				return err
			}
		}

		return nil
	}

	// drops plan details table
	migrations[2] = func() error {
		// NOOP migration, this was used to drop the plan_details table, but
		// there's more of a disincentive than incentive to do that because it could
		// leave operators wiping out plain details accidentally and not being able
		// to recover if they don't follow the upgrade path.
		return nil
	}

	var lastMigrationNumber = -1

	// if we've run any migrations before, we should have a migrations table, so find the last one we ran
	if db.HasTable("migrations") {
		var storedMigrations []models.Migration
		if err := db.Order("migration_id desc").Find(&storedMigrations).Error; err != nil {
			return fmt.Errorf("Error getting last migration id even though migration table exists: %s", err)
		}
		lastMigrationNumber = storedMigrations[0].MigrationId
	}

	// starting from the last migration we ran + 1, run migrations until we are current
	for i := lastMigrationNumber + 1; i < len(migrations); i++ {
		tx := db.Begin()
		err := migrations[i]()
		if err != nil {
			tx.Rollback()

			return err
		} else {
			newMigration := models.Migration{
				MigrationId: i,
			}
			if err := db.Save(&newMigration).Error; err != nil {
				tx.Rollback()
				return err
			} else {
				tx.Commit()
			}
		}
	}
	return nil
}

func autoMigrateTables(db *gorm.DB, tables ...interface{}) error {
	if db.Dialect().GetName() == "mysql" {
		return db.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8").AutoMigrate(tables...).Error
	} else {
		return db.AutoMigrate(tables...).Error
	}
}
