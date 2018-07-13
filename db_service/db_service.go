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
	"fmt"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/models"
	"github.com/jinzhu/gorm"
	"sync"
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

// gets the totaly number of service instances that are currently provisioned
func GetServiceInstanceTotal() (int, error) {
	var provisionedInstancesCount int
	err := DbConnection.Model(&models.ServiceInstanceDetails{}).Count(&provisionedInstancesCount).Error
	return provisionedInstancesCount, err
}

// gets the count of service instances by instance id (i.e. 0 or 1)
func GetServiceInstanceCount(instanceID string) (int, error) {
	var count int
	err := DbConnection.Model(&models.ServiceInstanceDetails{}).Where("id = ?", instanceID).Count(&count).Error
	return count, err
}

// soft deletes an instance from the database by instance id
func SoftDeleteInstanceDetails(instanceID string) error {
	// TODO(cbriant): how do I know if this is a connection error or a does not exist error
	instance := models.ServiceInstanceDetails{}
	if err := DbConnection.Where("ID = ?", instanceID).First(&instance).Error; err != nil {
		return models.ErrInstanceDoesNotExist
	}
	return DbConnection.Delete(&instance).Error
}

func GetLastOperation(instanceId string) (models.CloudOperation, error) {
	var op models.CloudOperation

	if err := DbConnection.Where("service_instance_id = ?", instanceId).Order("created_at desc").First(&op).Error; err != nil {
		return models.CloudOperation{}, err
	}
	return op, nil
}
