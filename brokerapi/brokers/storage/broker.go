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

package storage

import (
	"net/http"

	googlestorage "cloud.google.com/go/storage"
	"code.cloudfoundry.org/lager"
	"encoding/json"
	"fmt"
	"gcp-service-broker/brokerapi/brokers/broker_base"
	"gcp-service-broker/brokerapi/brokers/models"
	"gcp-service-broker/db_service"
	"golang.org/x/net/context"
	"google.golang.org/api/option"
)

type StorageBroker struct {
	Client    *http.Client
	ProjectId string
	Logger    lager.Logger

	broker_base.BrokerBase
}

// creates a new bucket with the name given in provision details and optional location
// (defaults to "US", for acceptable location values see: https://cloud.google.com/storage/docs/bucket-locations)
func (b *StorageBroker) Provision(instanceId string, details models.ProvisionDetails, plan models.PlanDetails) (models.ServiceInstanceDetails, error) {

	var err error

	var planDetails map[string]string
	if err = json.Unmarshal([]byte(plan.Features), &planDetails); err != nil {
		return models.ServiceInstanceDetails{}, fmt.Errorf("Error unmarshalling plan features: %s", err)
	}
	storageClass := planDetails["storage_class"]

	var params map[string]string
	if err = json.Unmarshal(details.RawParameters, &params); err != nil {
		return models.ServiceInstanceDetails{}, fmt.Errorf("Error unmarshalling parameters: %s", err)
	}

	// make a new bucket
	ctx := context.Background()
	co := option.WithUserAgent(models.CustomUserAgent)
	storageService, err := googlestorage.NewClient(ctx, co)
	if err != nil {
		return models.ServiceInstanceDetails{}, fmt.Errorf("Error creating new storage client: %s", err)
	}

	bucket := storageService.Bucket(params["name"])

	loc := "US"
	userLoc, locOk := params["location"]
	if locOk {
		loc = userLoc
	}
	attrs := googlestorage.BucketAttrs{
		Name:         params["name"],
		StorageClass: storageClass,
		Location:     loc,
	}

	// create the bucket. Nil uses default bucket attributes
	err = bucket.Create(ctx, b.ProjectId, &attrs)
	if err != nil {
		return models.ServiceInstanceDetails{}, fmt.Errorf("Error creating new bucket: %s", err)

	}

	otherDetails, err := json.Marshal(map[string]string{
		"StorageClass": attrs.StorageClass,
	})
	if err != nil {
		return models.ServiceInstanceDetails{}, fmt.Errorf("Error marshalling json: %s", err)
	}

	i := models.ServiceInstanceDetails{
		Name:         attrs.Name,
		Url:          "",
		Location:     attrs.Location,
		OtherDetails: string(otherDetails),
	}

	return i, nil
}

// Deletes the bucket associated with the given instance id
// Note that all objects within the bucket must be deleted first
func (b *StorageBroker) Deprovision(instanceID string, details models.DeprovisionDetails) error {
	bucket := models.ServiceInstanceDetails{}
	if err := db_service.DbConnection.Where("ID = ?", instanceID).First(&bucket).Error; err != nil {
		return models.ErrInstanceDoesNotExist
	}

	ctx := context.Background()
	storageService, err := googlestorage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("Error creating storage client: %s", err)
	}

	if err = storageService.Bucket(bucket.Name).Delete(ctx); err != nil {
		return fmt.Errorf("Error deleting bucket: %s", err)
	}

	return nil
}
