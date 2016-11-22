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

package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"time"
)

type ServiceBindingCredentials struct {
	gorm.Model

	OtherDetails string `sql:"type:text"`

	ServiceId         string
	ServiceInstanceId string
	BindingId         string
}

type ServiceInstanceDetails struct {
	ID        string `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time

	Name         string
	Location     string
	Url          string
	OtherDetails string `sql:"type:text"`

	ServiceId        string
	PlanId           string
	SpaceGuid        string
	OrganizationGuid string
}

type ProvisionRequestDetails struct {
	gorm.Model

	ServiceInstanceId string
	// is a json.Marshal of models.ProvisionDetails
	RequestDetails string
}

type PlanDetails struct {
	ID        string `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time

	ServiceId string
	Name      string
	Features  string `sql:"type:text"`
}

type Migration struct {
	gorm.Model

	MigrationId int
}

type CloudOperation struct {
	gorm.Model

	Name          string
	Status        string
	OperationType string
	ErrorMessage  string
	InsertTime    time.Time
	StartTime     time.Time
	TargetId      string
	TargetLink    string

	ServiceId         string
	ServiceInstanceId string
}
