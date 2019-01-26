// Copyright 2019 the Service Broker Project Authors.
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

package brokers

import (
	"context"

	"github.com/GoogleCloudPlatform/gcp-service-broker/db_service"
	"github.com/GoogleCloudPlatform/gcp-service-broker/db_service/models"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/broker"
	"github.com/jinzhu/gorm"
	"github.com/pivotal-cf/brokerapi"
)

// ServiceContext holds service definition and plan information.
type ServiceContext struct {
	Definition *broker.ServiceDefinition
	Plan       *broker.ServicePlan
}

// InstanceContext holds context around a service and a particular instance of it.
type InstanceContext struct {
	ServiceContext

	ServiceInstance *models.ServiceInstanceDetails
}

// BindingContext holds context around a binding, it's service instance, and the service definition.
type BindingContext struct {
	InstanceContext

	BindingInstance *models.ServiceBindingCredentials
}

// RequestContextService looks up the context surrounding a particular request
// from the database.
type RequestContextService struct {
	registry broker.BrokerRegistry
}

// ServiceContext gets the service context from the registry.
func (rc *RequestContextService) ServiceContext(ctx context.Context, serviceID, planID string) (*ServiceContext, error) {
	defn, err := rc.registry.GetServiceById(serviceID)
	if err != nil {
		return nil, err
	}

	plan, err := defn.GetPlanById(planID)
	if err != nil {
		return nil, err
	}

	return &ServiceContext{
		Definition: defn,
		Plan:       plan,
	}, nil
}

// InstanceContext gathers infromation about the instance from the database and
// service info from the registry.
func (rc *RequestContextService) InstanceContext(ctx context.Context, instanceID string) (*InstanceContext, error) {
	instance, err := rc.LookupInstanceDetails(ctx, instanceID)
	if err != nil {
		return nil, err
	}

	sc, err := rc.ServiceContext(ctx, instance.ServiceId, instance.PlanId)
	if err != nil {
		return nil, err
	}

	return &InstanceContext{
		ServiceContext:  *sc,
		ServiceInstance: instance,
	}, nil
}

// LookupInstanceDetails gets ServiceInstanceDetails given the instanceID.
// It returns an OSB ErrInstanceDoesNotExist if the instance doesn't exist.
func (rc *RequestContextService) LookupInstanceDetails(ctx context.Context, instanceID string) (*models.ServiceInstanceDetails, error) {
	instance, err := db_service.GetServiceInstanceDetailsById(ctx, instanceID)
	if err != nil && gorm.IsRecordNotFoundError(err) {
		return nil, brokerapi.ErrInstanceDoesNotExist
	}

	return instance, err
}

// BindingContext gathers infromation about the binding and instance from the
// database and service info from the registry.
func (rc *RequestContextService) BindingContext(ctx context.Context, instanceID, bindingID string) (*BindingContext, error) {
	ic, err := rc.InstanceContext(ctx, instanceID)
	if err != nil {
		return nil, err
	}

	binding, err := db_service.GetServiceBindingCredentialsByServiceInstanceIdAndBindingId(ctx, instanceID, bindingID)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, brokerapi.ErrBindingDoesNotExist
		}
		return nil, err
	}

	return &BindingContext{
		InstanceContext: *ic,
		BindingInstance: binding,
	}, nil
}

// NewRequestContextService creates a new RequestContextService for the given registry.
func NewRequestContextService(registry broker.BrokerRegistry) *RequestContextService {
	return &RequestContextService{
		registry: registry,
	}
}
