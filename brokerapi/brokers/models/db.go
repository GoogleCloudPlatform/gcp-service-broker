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
	"encoding/json"
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

func (sbc ServiceBindingCredentials) GetOtherDetails() map[string]string {
	var creds map[string]string
	if err := json.Unmarshal([]byte(sbc.OtherDetails), &creds); err != nil {
		panic(err)
	}
	return creds
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

func (si ServiceInstanceDetails) GetOtherDetails() map[string]string {
	var instanceDetails map[string]string
	// if the instance has access details saved
	if si.OtherDetails != "" {
		if err := json.Unmarshal([]byte(si.OtherDetails), &instanceDetails); err != nil {
			panic(err)
		}
	}
	return instanceDetails

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
	InsertTime    string
	StartTime     string
	TargetId      string
	TargetLink    string

	ServiceId         string
	ServiceInstanceId string
}
