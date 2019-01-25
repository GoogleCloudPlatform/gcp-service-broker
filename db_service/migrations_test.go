// Copyright 2019 the Service Broker Project Authors.
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
	"math"
	"os"
	"reflect"
	"testing"

	"github.com/GoogleCloudPlatform/gcp-service-broker/db_service/models"
	"github.com/jinzhu/gorm"
)

func TestValidateLastMigration(t *testing.T) {
	cases := map[string]struct {
		LastMigration int
		Expected      error
	}{
		"new-db": {
			LastMigration: -1,
			Expected:      nil,
		},
		"before-v2": {
			LastMigration: 0,
			Expected:      errors.New("Migration from broker versions <= 2.0 is no longer supported, upgrade using a v3.x broker then try again."),
		},
		"v3-to-v4": {
			LastMigration: 1,
			Expected:      nil,
		},
		"v4-to-v4.1": {
			LastMigration: 2,
			Expected:      nil,
		},
		"v4.1-to-v4.2": {
			LastMigration: 3,
			Expected:      nil,
		},
		"up-to-date": {
			LastMigration: numMigrations - 1,
			Expected:      nil,
		},
		"future": {
			LastMigration: numMigrations,
			Expected:      errors.New("The database you're connected to is newer than this tool supports."),
		},
		"far-future": {
			LastMigration: math.MaxInt32,
			Expected:      errors.New("The database you're connected to is newer than this tool supports."),
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			actual := ValidateLastMigration(tc.LastMigration)

			if !reflect.DeepEqual(actual, tc.Expected) {
				t.Errorf("Expected error %v, got %v", tc.Expected, actual)
			}
		})
	}
}

func TestRunMigrations_Failures(t *testing.T) {
	cases := map[string]struct {
		LastMigration int
		Expected      error
	}{
		"before-v2": {
			LastMigration: 0,
			Expected:      errors.New("Migration from broker versions <= 2.0 is no longer supported, upgrade using a v3.x broker then try again."),
		},
		"future": {
			LastMigration: numMigrations,
			Expected:      errors.New("The database you're connected to is newer than this tool supports."),
		},
		"far-future": {
			LastMigration: math.MaxInt32,
			Expected:      errors.New("The database you're connected to is newer than this tool supports."),
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			db, err := gorm.Open("sqlite3", "test.sqlite3")
			defer os.Remove("test.sqlite3")
			if err != nil {
				t.Fatal(err)
			}

			if err := autoMigrateTables(db, &models.MigrationV1{}); err != nil {
				t.Fatal(err)
			}

			if err := db.Save(&models.Migration{MigrationId: tc.LastMigration}).Error; err != nil {
				t.Fatal(err)
			}

			actual := RunMigrations(db)
			if !reflect.DeepEqual(actual, tc.Expected) {
				t.Errorf("Expected error %v, got %v", tc.Expected, actual)
			}
		})
	}
}

func TestRunMigrations(t *testing.T) {
	cases := map[string]func(t *testing.T, db *gorm.DB){
		"creates-migrations-table": func(t *testing.T, db *gorm.DB) {
			if err := RunMigrations(db); err != nil {
				t.Fatal(err)
			}

			if !db.HasTable("migrations") {
				t.Error("Expected db to have migrations table")
			}
		},

		"applies-all-migrations-when-run": func(t *testing.T, db *gorm.DB) {
			if err := RunMigrations(db); err != nil {
				t.Fatal(err)
			}

			var storedMigrations []models.Migration
			if err := db.Order("id desc").Find(&storedMigrations).Error; err != nil {
				t.Fatal(err)
			}

			lastMigrationNumber := storedMigrations[0].MigrationId
			if lastMigrationNumber != numMigrations-1 {
				t.Errorf("expected lastMigrationNumber to be %d, got %d", numMigrations-1, lastMigrationNumber)
			}
		},

		"can-run-migrations-multiple-times": func(t *testing.T, db *gorm.DB) {
			for i := 0; i < 10; i++ {
				if err := RunMigrations(db); err != nil {
					t.Fatal(err)
				}
			}
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			db, err := gorm.Open("sqlite3", "test.sqlite3")
			defer os.Remove("test.sqlite3")
			if err != nil {
				t.Fatal(err)
			}

			tc(t, db)
		})
	}
}
