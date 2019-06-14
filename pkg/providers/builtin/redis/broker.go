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

package redis

import (
	googleredis "cloud.google.com/go/redis/apiv1beta1"
	"google.golang.org/genproto/googleapis/cloud/redis/v1beta1"
	"github.com/GoogleCloudPlatform/gcp-service-broker/db_service/models"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/providers/builtin/base"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/varcontext"
	"golang.org/x/net/context"
)

// RedisBroker is the service-broker back-end for creating and binding Redis services.
type RedisBroker struct {
	base.BrokerBase
}

// Provision creates a new Redis instance from the settings in the user-provided details and service plan.
func (b *RedisBroker) Provision(ctx context.Context, provisionContext *varcontext.VarContext) (models.ServiceInstanceDetails, error) {

	// Build Instance parameter
	instance := &redis.Instance{
		Name: "default",
		Tier: redis.Instance_BASIC,
		MemorySizeGb: 32,
	}

	// Fill in InstanceRequest
	ir := &redis.CreateInstanceRequest{
		Parent: "parent",
		InstanceId: "generic_id",
		Instance: instance,
	}

	// CreateInstance creates a Redis instance and executes asynchronously
	ctx := context.Background()
	c, err := googleredis.NewCloudRedisClient(ctx)
	if err != nil {
		// TODO (hsophia):
	}

	op, err := c.CreateInstance(ctx, ir)
	if err != nil {
		// TODO (hsophia): Handle error
	}

	resp, err := op.Wait(ctx)
	if err != nil {
		// TODO (hsophia): Handle error
	}

	// Get ID to return to serviceInstanceDetails
	id := models.ServiceInstanceDetails{
		Name: "default",
	}

	return id, nil
}