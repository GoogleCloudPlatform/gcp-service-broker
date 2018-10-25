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

package storage

import (
	"encoding/json"
	"fmt"

	googlestorage "cloud.google.com/go/storage"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/broker_base"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/models"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/varcontext"
	"github.com/pivotal-cf/brokerapi"
	"golang.org/x/net/context"
	"google.golang.org/api/option"
)

// StorageBroker is the service-broker back-end for creating and binding to
// Google Cloud Storage buckets.
type StorageBroker struct {
	broker_base.BrokerBase
}

// InstanceInformation holds the details needed to connect to a GCS instance
// after it has been provisioned.
type InstanceInformation struct {
	BucketName string `json:"bucket_name"`
}

// Provision creates a new GCS bucket from the settings in the user-provided details and service plan.
func (b *StorageBroker) Provision(ctx context.Context, provisionContext *varcontext.VarContext) (models.ServiceInstanceDetails, error) {
	attrs := googlestorage.BucketAttrs{
		Name:         provisionContext.GetString("name"),
		StorageClass: provisionContext.GetString("storage_class"),
		Location:     provisionContext.GetString("location"),
		Labels:       provisionContext.GetStringMapString("labels"),
	}

	if err := provisionContext.Error(); err != nil {
		return models.ServiceInstanceDetails{}, err
	}

	// make a new bucket
	storageService, err := b.createClient(ctx)
	if err != nil {
		return models.ServiceInstanceDetails{}, err
	}

	// create the bucket. Nil uses default bucket attributes
	if err := storageService.Bucket(attrs.Name).Create(ctx, b.ProjectId, &attrs); err != nil {
		return models.ServiceInstanceDetails{}, fmt.Errorf("Error creating new bucket: %s", err)
	}

	ii := InstanceInformation{
		BucketName: attrs.Name,
	}

	otherDetails, err := json.Marshal(ii)
	if err != nil {
		return models.ServiceInstanceDetails{}, fmt.Errorf("Error marshalling json: %s", err)
	}

	return models.ServiceInstanceDetails{
		Name:         attrs.Name,
		Url:          "",
		Location:     attrs.Location,
		OtherDetails: string(otherDetails),
	}, nil
}

// Deprovision deletes the bucket associated with the given instance.
// Note that all objects within the bucket must be deleted first.
func (b *StorageBroker) Deprovision(ctx context.Context, bucket models.ServiceInstanceDetails, details brokerapi.DeprovisionDetails) (*string, error) {
	storageService, err := b.createClient(ctx)
	if err != nil {
		return nil, err
	}

	if err = storageService.Bucket(bucket.Name).Delete(ctx); err != nil {
		return nil, fmt.Errorf("Error deleting bucket: %s", err)
	}

	return nil, nil
}

func (b *StorageBroker) createClient(ctx context.Context) (*googlestorage.Client, error) {
	co := option.WithUserAgent(models.CustomUserAgent)
	ct := option.WithTokenSource(b.HttpConfig.TokenSource(ctx))
	storageService, err := googlestorage.NewClient(ctx, co, ct)
	if err != nil {
		return nil, fmt.Errorf("Couldn't instantiate Cloud Storage API client: %s", err)
	}

	return storageService, nil
}
