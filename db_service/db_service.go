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

//go:generate go run dao_generator.go

package db_service

import (
	"fmt"
	"sync"

	"code.cloudfoundry.org/lager"
	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var DbConnection *gorm.DB
var once sync.Once

// Instantiates the db connection and runs migrations
func New(logger lager.Logger) *gorm.DB {
	once.Do(func() {
		DbConnection = SetupDb(logger)
		if err := RunMigrations(DbConnection); err != nil {
			panic(fmt.Sprintf("Error migrating database: %s", err.Error()))
		}
	})
	return DbConnection
}

// defaultDatastore gets the default datastore for the given default database
// instantiated in New(). In the future, all accesses of DbConnection will be
// done through SqlDatastore and it will become the globally shared instance.
func defaultDatastore() *SqlDatastore {
	return &SqlDatastore{db: DbConnection}
}

type SqlDatastore struct {
	db *gorm.DB
}
