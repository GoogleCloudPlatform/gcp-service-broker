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
	"os"

	"code.cloudfoundry.org/lager"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/models"
	"github.com/GoogleCloudPlatform/gcp-service-broker/fakes"
	"github.com/jinzhu/gorm"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("DbService", func() {
	var (
		err    error
		logger lager.Logger
		testDb *gorm.DB
	)

	BeforeEach(func() {
		logger = lager.NewLogger("brokers_test")
		logger.RegisterSink(lager.NewWriterSink(GinkgoWriter, lager.DEBUG))

		fakes.SetUpTestServices()

		testDb, err = gorm.Open("sqlite3", "test.sqlite3")
		Expect(err).NotTo(HaveOccurred())
	})

	Describe("Migrations", func() {
		It("should create a migrations table", func() {
			err = RunMigrations(testDb)
			Expect(err).NotTo(HaveOccurred())
			Expect(testDb.HasTable("migrations")).To(BeTrue())
		})

		It("should apply all migrations when run", func() {
			err = RunMigrations(testDb)
			Expect(err).NotTo(HaveOccurred())
			var storedMigrations []models.Migration
			err = testDb.Order("id desc").Find(&storedMigrations).Error
			Expect(err).NotTo(HaveOccurred())
			lastMigrationNumber := storedMigrations[0].MigrationId
			Expect(lastMigrationNumber).To(Equal(2))
		})

		It("should be able to run migrations multiple times", func() {
			err = RunMigrations(testDb)
			Expect(err).NotTo(HaveOccurred())
			err = RunMigrations(testDb)
			Expect(err).NotTo(HaveOccurred())
		})

	})

	AfterEach(func() {
		os.Remove("test.sqlite3")
	})
})
