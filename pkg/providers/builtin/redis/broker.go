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
	"fmt"
	"github.com/GoogleCloudPlatform/gcp-service-broker/db_service/models"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/providers/builtin/base"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/varcontext"
	"github.com/pivotal-cf/brokerapi"
	"golang.org/x/net/context"
	"google.golang.org/genproto/googleapis/cloud/redis/v1beta1"
)

// RedisBroker is the service-broker back-end for creating and binding Redis services.
type RedisBroker struct {
	base.BrokerBase
}

// serviceTiers holds the valid value mapping for string service tiers to their REST call equivalent
var serviceTiers = map[string]redis.Instance_Tier{
	"basic": redis.Instance_BASIC,
	"standard_ha": redis.Instance_STANDARD_HA,
}

// Provision creates a new Redis instance from the settings in the user-provided details and service plan.
func (b *RedisBroker) Provision(ctx context.Context, provisionContext *varcontext.VarContext) (models.ServiceInstanceDetails, error) {

	authorizedNetwork := provisionContext.GetString("authorized_network")
	capacityTier := int32(provisionContext.GetInt("capacity_tier"))
	displayName := provisionContext.GetString("display_name")
	instanceId := provisionContext.GetString("instance_id")
	locationId := provisionContext.GetString("location_id")
	serviceTier := serviceTiers[provisionContext.GetString("service_tier")]
	parent := fmt.Sprintf("projects/%s/locations/%s", b.ProjectId, locationId)
	name := fmt.Sprintf("%s/instances/%s", parent, instanceId)

	// Build Redis Instance
	instance := &redis.Instance{
		Name: name,
		DisplayName: displayName,
		Tier: serviceTier,
		MemorySizeGb: capacityTier,
		AuthorizedNetwork: authorizedNetwork,
	}

	ir := &redis.CreateInstanceRequest{
		Parent: parent,
		InstanceId: instanceId,
		Instance: instance,
	}

	c, err := googleredis.NewCloudRedisClient(ctx)
	if err != nil {
		return models.ServiceInstanceDetails{}, err
	}

	op, err := c.CreateInstance(ctx, ir)
	if err != nil {
		return models.ServiceInstanceDetails{}, fmt.Errorf("Error creating new Redis instance: %s", err)
	}

	resp, err := op.Wait(ctx)
	if err != nil {
		return models.ServiceInstanceDetails{}, err
	}

	id := models.ServiceInstanceDetails{
		Name: resp.Name,
	}

	return id, nil
}

// Deprovision deletes the Redis instance with the given instance ID
func (b *RedisBroker) Deprovision(ctx context.Context, instance models.ServiceInstanceDetails, details brokerapi.DeprovisionDetails) (*string, error) {
	c, err := googleredis.NewCloudRedisClient(ctx)
	if err != nil {
		return nil, err
	}

	req := &redis.DeleteInstanceRequest{
		Name: instance.Name,
	}

	op, err := c.DeleteInstance(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("Error deleting Redis instance: %s", err)
	}

	err = op.Wait(ctx)

	return nil, nil
}