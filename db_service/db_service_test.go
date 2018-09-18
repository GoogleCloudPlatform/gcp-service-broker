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

		os.Setenv("ROOT_SERVICE_ACCOUNT_JSON", `{
			  "//": "Dummy account from https://github.com/GoogleCloudPlatform/google-cloud-java/google-cloud-clients/google-cloud-core/src/test/java/com/google/cloud/ServiceOptionsTest.java",
			  "private_key_id": "somekeyid",
			  "private_key": "-----BEGIN PRIVATE KEY-----\nMIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQC+K2hSuFpAdrJI\nnCgcDz2M7t7bjdlsadsasad+fvRSW6TjNQZ3p5LLQY1kSZRqBqylRkzteMOyHgaR\n0Pmxh3ILCND5men43j3h4eDbrhQBuxfEMalkG92sL+PNQSETY2tnvXryOvmBRwa/\nQP/9dJfIkIDJ9Fw9N4Bhhhp6mCcRpdQjV38H7JsyJ7lih/oNjECgYAt\nknddadwkwewcVxHFhcZJO+XWf6ofLUXpRwiTZakGMn8EE1uVa2LgczOjwWHGi99MFjxSer5m9\n1tCa3/KEGKiS/YL71JvjwX3mb+cewlkcmweBKZHM2JPTk0ZednFSpVZMtycjkbLa\ndYOS8V85AgMBewECggEBAKksaldajfDZDV6nGqbFjMiizAKJolr/M3OQw16K6o3/\n0S31xIe3sSlgW0+UbYlF4U8KifhManD1apVSC3csafaspP4RZUHFhtBywLO9pR5c\nr6S5aLp+gPWFyIp1pfXbWGvc5VY/v9x7ya1VEa6rXvLsKupSeWAW4tMj3eo/64ge\nsdaceaLYw52KeBYiT6+vpsnYrEkAHO1fF/LavbLLOFJmFTMxmsNaG0tuiJHgjshB\n82DpMCbXG9YcCgI/DbzuIjsdj2JC1cascSP//3PmefWysucBQe7Jryb6NQtASmnv\nCdDw/0jmZTEjpe4S1lxfHplAhHFtdgYTvyYtaLZiVVkCgYEA8eVpof2rceecw/I6\n5ng1q3Hl2usdWV/4mZMvR0fOemacLLfocX6IYxT1zA1FFJlbXSRsJMf/Qq39mOR2\nSpW+hr4jCoHeRVYLgsbggtrevGmILAlNoqCMpGZ6vDmJpq6ECV9olliDvpPgWOP+\nmYPDreFBGxWvQrADNbRt2dmGsrsCgYEAyUHqB2wvJHFqdmeBsaacewzV8x9WgmeX\ngUIi9REwXlGDW0Mz50dxpxcKCAYn65+7TCnY5O/jmL0VRxU1J2mSWyWTo1C+17L0\n3fUqjxL1pkefwecxwecvC+gFFYdJ4CQ/MHHXU81Lwl1iWdFCd2UoGddYaOF+KNeM\nHC7cmqra+JsCgYEAlUNywzq8nUg7282E+uICfCB0LfwejuymR93CtsFgb7cRd6ak\nECR8FGfCpH8ruWJINllbQfcHVCX47ndLZwqv3oVFKh6pAS/vVI4dpOepP8++7y1u\ncoOvtreXCX6XqfrWDtKIvv0vjlHBhhhp6mCcRpdQjV38H7JsyJ7lih/oNjECgYAt\nkndj5uNl5SiuVxHFhcZJO+XWf6ofLUregtevZakGMn8EE1uVa2AY7eafmoU/nZPT\n00YB0TBATdCbn/nBSuKDESkhSg9s2GEKQZG5hBmL5uCMfo09z3SfxZIhJdlerreP\nJ7gSidI12N+EZxYd4xIJh/HFDgp7RRO87f+WJkofMQKBgGTnClK1VMaCRbJZPriw\nEfeFCoOX75MxKwXs6xgrw4W//AYGGUjDt83lD6AZP6tws7gJ2IwY/qP7+lyhjEqN\nHtfPZRGFkGZsdaksdlaksd323423d+15/UvrlRSFPNj1tWQmNKkXyRDW4IG1Oa2p\nrALStNBx5Y9t0/LQnFI4w3aG\n-----END PRIVATE KEY-----\n",
			  "client_email": "someclientid@developer.gserviceaccount.com",
			  "client_id": "someclientid.apps.googleusercontent.com",
			  "type": "service_account",
			  "project_id": "my-project-123"
			}`)

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
