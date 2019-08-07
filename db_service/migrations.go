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
	"errors"
	"fmt"

	"github.com/GoogleCloudPlatform/gcp-service-broker/db_service/models"
	"github.com/jinzhu/gorm"
)

const numMigrations = 7

// runs schema migrations on the provided service broker database to get it up to date
func RunMigrations(db *gorm.DB) error {
	migrations := make([]func() error, numMigrations)

	// initial migration - creates tables
	migrations[0] = func() error { // v1.0
		return autoMigrateTables(db,
			&models.ServiceInstanceDetailsV1{},
			&models.ServiceBindingCredentialsV1{},
			&models.ProvisionRequestDetailsV1{},
			&models.PlanDetailsV1{},
			&models.MigrationV1{})
	}

	// adds CloudOperation table
	migrations[1] = func() error { // v2.x
		// NOTE: this migration used to have lots of custom logic, however it has
		// been removed because brokers starting at v4 no longer support the
		// functionality the migration required.
		//
		// It is acceptable to pass through this migration step on the way to
		// intiailize a _new_ databse, but it is not acceptable to use this step
		// in a path through the upgrade.
		return autoMigrateTables(db, &models.CloudOperationV1{})
	}

	// drops plan details table
	migrations[2] = func() error { // 4.0.0
		// NOOP migration, this was used to drop the plan_details table, but
		// there's more of a disincentive than incentive to do that because it could
		// leave operators wiping out plain details accidentally and not being able
		// to recover if they don't follow the upgrade path.
		return nil
	}

	migrations[3] = func() error { // v4.1.0
		return autoMigrateTables(db, &models.ServiceInstanceDetailsV2{})
	}

	migrations[4] = func() error { // v4.2.0
		return autoMigrateTables(db, &models.TerraformDeploymentV1{})
	}

	migrations[5] = func() error { // v4.2.3
		return autoMigrateTables(db, &models.ProvisionRequestDetailsV2{})
	}

	migrations[6] = func() error { // v4.2.4
		if db.Dialect().GetName() == "sqlite3" {
			// sqlite does not support changing column data types
			return nil
		} else {
			return db.Model(&models.ProvisionRequestDetailsV2{}).ModifyColumn("request_details", "text").Error
		}
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

	if err := ValidateLastMigration(lastMigrationNumber); err != nil {
		return err
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

// ValidateLastMigration returns an error if the database version is newer than
// this tool supports or is too old to be updated.
func ValidateLastMigration(lastMigration int) error {
	switch {
	case lastMigration >= numMigrations:
		return errors.New("The database you're connected to is newer than this tool supports.")

	case lastMigration == 0:
		return errors.New("Migration from broker versions <= 2.0 is no longer supported, upgrade using a v3.x broker then try again.")

	default:
		return nil
	}
}

func autoMigrateTables(db *gorm.DB, tables ...interface{}) error {
	if db.Dialect().GetName() == "mysql" {
		return db.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8").AutoMigrate(tables...).Error
	} else {
		return db.AutoMigrate(tables...).Error
	}
}
