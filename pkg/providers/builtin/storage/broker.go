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
	"fmt"

	googlestorage "cloud.google.com/go/storage"
	"github.com/GoogleCloudPlatform/gcp-service-broker/db_service/models"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/providers/builtin/base"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/varcontext"
	"github.com/GoogleCloudPlatform/gcp-service-broker/utils"
	"github.com/pivotal-cf/brokerapi"
	"golang.org/x/net/context"
	"google.golang.org/api/option"
)

// StorageBroker is the service-broker back-end for creating and binding to
// Google Cloud Storage buckets.
type StorageBroker struct {
	base.BrokerBase
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

	if provisionContext.GetBool("force_delete") {
		attrs.Labels["sb-force-delete"] = "true"
	} else {
		attrs.Labels["sb-force-delete"] = "false"
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

	id := models.ServiceInstanceDetails{
		Name:     attrs.Name,
		Location: attrs.Location,
	}

	if err := id.SetOtherDetails(ii); err != nil {
		return models.ServiceInstanceDetails{}, fmt.Errorf("Error marshalling json: %s", err)
	}

	return id, nil
}

// Deprovision deletes the bucket associated with the given instance.
// Note that all objects within the bucket must be deleted first.
func (b *StorageBroker) Deprovision(ctx context.Context, bucket models.ServiceInstanceDetails, details brokerapi.DeprovisionDetails) (*string, error) {
	storageService, err := b.createClient(ctx)
	if err != nil {
		return nil, err
	}

	attrs, err := storageService.Bucket(bucket.Name).Attrs(ctx)
	if err != nil {
		return nil, err
	}

	if attrs.Labels["sb-force-delete"] == "true" {
		objects := storageService.Bucket(bucket.Name).Objects(ctx, nil)

		for {
			obj, err := objects.Next()
			if err != nil || obj == nil {
				break
			}

			storageService.Bucket(bucket.Name).Object(obj.Name).Delete(ctx)
		}
	}

	if err = storageService.Bucket(bucket.Name).Delete(ctx); err != nil {
		return nil, fmt.Errorf("error deleting bucket: %s (to delete a non-empty bucket, set the label sb-force-delete=true on it)", err)
	}

	return nil, nil
}

func (b *StorageBroker) createClient(ctx context.Context) (*googlestorage.Client, error) {
	co := option.WithUserAgent(utils.CustomUserAgent)
	ct := option.WithTokenSource(b.HttpConfig.TokenSource(ctx))
	storageService, err := googlestorage.NewClient(ctx, co, ct)
	if err != nil {
		return nil, fmt.Errorf("Couldn't instantiate Cloud Storage API client: %s", err)
	}

	return storageService, nil
}
