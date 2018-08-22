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
	"database/sql"
	"fmt"
	"os"

	"code.cloudfoundry.org/lager"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/models"
	"github.com/GoogleCloudPlatform/gcp-service-broker/fakes"
	"github.com/jinzhu/gorm"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func getLocalTestConnectionStr(dbName string) string {

	username := os.Getenv("TEST_DB_USERNAME")
	password := os.Getenv("TEST_DB_PASSWORD")
	host := os.Getenv("TEST_DB_HOST")

	return fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?parseTime=true", username, password, host, dbName)
}

func createTestDatabase() {

	db, err := sql.Open("mysql", getLocalTestConnectionStr(""))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	res, err := db.Query("SHOW DATABASES LIKE 'servicebrokertest'")
	if err != nil {
		panic(err)
	}
	if res.Next() {
		dropTestDatabase()
	}

	_, err = db.Exec("CREATE DATABASE servicebrokertest")
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("USE servicebrokertest")
	if err != nil {
		panic(err)
	}
}

func dropTestDatabase() {
	db, err := sql.Open("mysql", getLocalTestConnectionStr(""))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	_, err = db.Exec("DROP DATABASE servicebrokertest")
	if err != nil {
		panic(err)
	}
}

var _ = Describe("DbService", func() {
	var (
		err    error
		logger lager.Logger
	)

	BeforeEach(func() {
		logger = lager.NewLogger("brokers_test")
		logger.RegisterSink(lager.NewWriterSink(GinkgoWriter, lager.DEBUG))

		fakes.SetUpTestServices()

		createTestDatabase()
		testDb, _ := gorm.Open("mysql", getLocalTestConnectionStr("servicebrokertest"))

		DbConnection = testDb
	})

	Describe("Migrations", func() {
		It("should create a migrations table", func() {
			err = RunMigrations(DbConnection)
			Expect(err).NotTo(HaveOccurred())
			Expect(DbConnection.HasTable("migrations")).To(BeTrue())
		})

		It("should apply all migrations when run", func() {
			err = RunMigrations(DbConnection)
			Expect(err).NotTo(HaveOccurred())
			var storedMigrations []models.Migration
			err = DbConnection.Order("id desc").Find(&storedMigrations).Error
			Expect(err).NotTo(HaveOccurred())
			lastMigrationNumber := storedMigrations[0].MigrationId
			Expect(lastMigrationNumber).To(Equal(2))
		})

		It("should be able to run migrations multiple times", func() {
			err = RunMigrations(DbConnection)
			Expect(err).NotTo(HaveOccurred())
			err = RunMigrations(DbConnection)
			Expect(err).NotTo(HaveOccurred())
		})

	})

	AfterEach(func() {
		dropTestDatabase()
	})
})
