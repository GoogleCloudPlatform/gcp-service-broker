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
	"code.cloudfoundry.org/lager"
	"github.com/jinzhu/gorm"
	"gcp-service-broker/brokerapi/brokers/models"
	"sync"
)

var DbConnection *gorm.DB
var once sync.Once

func New(logger lager.Logger) *gorm.DB {
	once.Do(func() {
		DbConnection = SetupDb(logger)
		RunMigrations(DbConnection)
	})
	return DbConnection
}

func GetServiceInstanceTotal() (int, error) {
	var provisionedInstancesCount int
	err := DbConnection.Model(&models.ServiceInstanceDetails{}).Count(&provisionedInstancesCount).Error
	return provisionedInstancesCount, err
}

func GetServiceInstanceCount(instanceID string) (int, error) {
	var count int
	err := DbConnection.Model(&models.ServiceInstanceDetails{}).Where("id = ?", instanceID).Count(&count).Error
	return count, err
}

func SoftDeleteInstanceDetails(instanceID string) error {
	// TODO(cbriant): how do I know if this is a connection error or a does not exist error
	instance := models.ServiceInstanceDetails{}
	if err := DbConnection.Where("ID = ?", instanceID).First(&instance).Error; err != nil {
		return models.ErrInstanceDoesNotExist
	}
	return DbConnection.Delete(&instance).Error
}
