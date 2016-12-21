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

package db_service

import (
	"encoding/json"
	"fmt"
	"gcp-service-broker/brokerapi/brokers/models"
	"gcp-service-broker/utils"
	"github.com/jinzhu/gorm"
	googlecloudsql "google.golang.org/api/sqladmin/v1beta4"
	"os"
)

// runs schema migrations on the provided service broker database to get it up to date
func RunMigrations(db *gorm.DB) error {

	migrations := make([]func() error, 2)

	// initial migration - creates tables
	migrations[0] = func() error {
		if err := db.Exec(`CREATE TABLE service_instance_details (
			  id varchar(255) NOT NULL DEFAULT '',
			  created_at timestamp NULL DEFAULT NULL,
			  updated_at timestamp NULL DEFAULT NULL,
			  deleted_at timestamp NULL DEFAULT NULL,
			  name varchar(255) DEFAULT NULL,
			  location varchar(255) DEFAULT NULL,
			  url varchar(255) DEFAULT NULL,
			  other_details text,
			  service_id varchar(255) DEFAULT NULL,
			  plan_id varchar(255) DEFAULT NULL,
			  space_guid varchar(255) DEFAULT NULL,
			  organization_guid varchar(255) DEFAULT NULL,
			  PRIMARY KEY (id)
			) ENGINE=InnoDB DEFAULT CHARSET=utf8`).Error; err != nil {
			return err
		}
		if err := db.Exec(`CREATE TABLE service_binding_credentials (
			  id int(10) unsigned NOT NULL AUTO_INCREMENT,
			  created_at timestamp NULL DEFAULT NULL,
			  updated_at timestamp NULL DEFAULT NULL,
			  deleted_at timestamp NULL DEFAULT NULL,
			  other_details text,
			  service_id varchar(255) DEFAULT NULL,
			  service_instance_id varchar(255) DEFAULT NULL,
			  binding_id varchar(255) DEFAULT NULL,
			  PRIMARY KEY (id),
			  KEY idx_service_binding_credentials_deleted_at (deleted_at)
			) ENGINE=InnoDB DEFAULT CHARSET=utf8`).Error; err != nil {
			return err
		}
		if err := db.Exec(`CREATE TABLE provision_request_details (
			  id int(10) unsigned NOT NULL AUTO_INCREMENT,
			  created_at timestamp NULL DEFAULT NULL,
			  updated_at timestamp NULL DEFAULT NULL,
			  deleted_at timestamp NULL DEFAULT NULL,
			  service_instance_id varchar(255) DEFAULT NULL,
			  request_details varchar(255) DEFAULT NULL,
			  PRIMARY KEY (id),
			  KEY idx_provision_request_details_deleted_at (deleted_at)
			) ENGINE=InnoDB DEFAULT CHARSET=utf8`).Error; err != nil {
			return err
		}
		if err := db.Exec(`CREATE TABLE plan_details (
			  id varchar(255) NOT NULL DEFAULT '',
			  created_at timestamp NULL DEFAULT NULL,
			  updated_at timestamp NULL DEFAULT NULL,
			  deleted_at timestamp NULL DEFAULT NULL,
			  service_id varchar(255) DEFAULT NULL,
			  name varchar(255) DEFAULT NULL,
			  features text,
			  PRIMARY KEY (id)
			) ENGINE=InnoDB DEFAULT CHARSET=utf8`).Error; err != nil {
			return err
		}
		if err := db.Exec(`CREATE TABLE migrations (
			  id int(10) unsigned NOT NULL AUTO_INCREMENT,
			  created_at timestamp NULL DEFAULT NULL,
			  updated_at timestamp NULL DEFAULT NULL,
			  deleted_at timestamp NULL DEFAULT NULL,
			  migration_id int(10) DEFAULT NULL,
			  PRIMARY KEY (id)
			) ENGINE=InnoDB DEFAULT CHARSET=utf8`).Error; err != nil {
			return err
		}
		return nil
	}

	// adds CloudOperation table
	migrations[1] = func() error {
		if err := db.Exec(`CREATE TABLE cloud_operations (
			  id int(10) unsigned NOT NULL AUTO_INCREMENT,
			  created_at timestamp NULL DEFAULT NULL,
			  updated_at timestamp NULL DEFAULT NULL,
			  deleted_at timestamp NULL DEFAULT NULL,
			  name varchar(255) DEFAULT NULL,
			  status varchar(255) DEFAULT NULL,
			  operation_type varchar(255) DEFAULT NULL,
			  error_message text,
			  insert_time varchar(255) DEFAULT NULL,
			  start_time varchar(255) DEFAULT NULL,
			  target_id varchar(255) DEFAULT NULL,
			  target_link varchar(255) DEFAULT NULL,
			  service_id varchar(255) DEFAULT NULL,
			  service_instance_id varchar(255) DEFAULT NULL,
			  PRIMARY KEY (id)
			) ENGINE=InnoDB DEFAULT CHARSET=utf8`).Error; err != nil {
			return err
		}

		// copy provision request details into service instance details

		serviceAccount := make(map[string]string)
		if err := json.Unmarshal([]byte(os.Getenv("ROOT_SERVICE_ACCOUNT_JSON")), &serviceAccount); err != nil {
			return err
		}

		client, err := utils.GetAuthedClient()
		if err != nil {
			return fmt.Errorf("Error getting authorized http client: %s", err)
		}

		idToNameMap, err := utils.MapServiceIdToName()
		if err != nil {
			return err
		}

		var prs []models.ProvisionRequestDetails
		if err := DbConnection.Find(&prs).Error; err != nil {
			return err
		}

		for _, pr := range prs {
			var si models.ServiceInstanceDetails
			if err := DbConnection.Where("id = ?", pr.ServiceInstanceId).First(&si).Error; err != nil {
				return err
			}
			od := make(map[string]string)
			if err := json.Unmarshal([]byte(pr.RequestDetails), &od); err != nil {
				return err
			}
			newOd := make(map[string]string)

			// cloudsql
			switch serviceName := idToNameMap[si.ServiceId]; serviceName {
			case models.CloudsqlName:
				newOd["instance_name"] = od["instance_name"]
				newOd["database_name"] = od["database_name"]

				sqlService, err := googlecloudsql.New(client)
				if err != nil {
					return fmt.Errorf("Error creating new CloudSQL Client: %s", err)
				}
				dbService := googlecloudsql.NewInstancesService(sqlService)
				clouddb, err := dbService.Get(serviceAccount["project_id"], od["instance_name"]).Do()
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
				println(fmt.Sprintf("%v", idToNameMap))
				return fmt.Errorf("unrecognized service: %s", si.ServiceId)
			}

			odBytes, err := json.Marshal(&newOd)
			if err != nil {
				return err
			}
			si.OtherDetails = string(odBytes)
			if err := DbConnection.Save(&si).Error; err != nil {
				return err
			}
		}

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
