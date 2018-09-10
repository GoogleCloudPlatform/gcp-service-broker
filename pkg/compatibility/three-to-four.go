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

package compatibility

import (
	"context"
	"errors"
	"fmt"

	"github.com/GoogleCloudPlatform/gcp-service-broker/db_service"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/broker"
	"github.com/pivotal-cf/brokerapi"
)

var (
	errOperationUnsupported = errors.New("this operation is unsupported in 3 to 4 upgrade mode")
)

type ThreeToFour struct {
}

func (t *ThreeToFour) LastOperation(ctx context.Context, instanceID, operationData string) (brokerapi.LastOperation, error) {
	// pass on
	return brokerapi.LastOperation{}, errOperationUnsupported
}

func (t *ThreeToFour) Provision(ctx context.Context, instanceID string, details brokerapi.ProvisionDetails, asyncAllowed bool) (brokerapi.ProvisionedServiceSpec, error) {
	// if legacy, reject, alert to new
	return brokerapi.ProvisionedServiceSpec{}, errOperationUnsupported
}

func (t *ThreeToFour) Deprovision(ctx context.Context, instanceID string, details brokerapi.DeprovisionDetails, asyncAllowed bool) (brokerapi.DeprovisionServiceSpec, error) {
	// if legacy, reject, alert to new
	return brokerapi.DeprovisionServiceSpec{}, errOperationUnsupported
}

func (t *ThreeToFour) Bind(ctx context.Context, instanceID, bindingID string, details brokerapi.BindDetails) (brokerapi.Binding, error) {
	// if legacy, reject, alert to new
	return brokerapi.Binding{}, errOperationUnsupported
}

func (t *ThreeToFour) Unbind(ctx context.Context, instanceID, bindingID string, details brokerapi.UnbindDetails) error {
	// if legacy, reject, alert to new
	return errOperationUnsupported
}

func (t *ThreeToFour) Services(ctx context.Context) ([]brokerapi.Service, error) {
	// inject legacy definitions + prefix with "legacy-"

	services := broker.GetEnabledServices()

	var marketplace []brokerapi.Service

	for _, svc := range services {
		catalog, err := svc.CatalogEntry()
		if err != nil {
			return nil, err
		}
		serviceDefn := catalog.ToPlain()
		var additionalPlans []brokerapi.ServicePlan

		for _, plan := range serviceDefn.Plans {
			planDetails, err := db_service.GetPlanDetailsV1ByServiceIdAndName(serviceDefn.ID, plan.Name)
			if err == nil {
				plan.Name = "upgrademe-" + plan.Name
				plan.ID = planDetails.ID
				additionalPlans = append(additionalPlans, plan)
			}
		}

		serviceDefn.Plans = append(serviceDefn.Plans, additionalPlans...)
		marketplace = append(marketplace, serviceDefn)
	}

	return marketplace, nil
}

func (t *ThreeToFour) Update(ctx context.Context, instanceID string, details brokerapi.UpdateDetails, asyncAllowed bool) (brokerapi.UpdateServiceSpec, error) {
	// if not legacy, pass on

	// if legacy, allow upgrading to new

	// TODO validate that the new plan id is an acceptable upgrade from the old one.
	// does the name of the new plan equal that of the old one

	instanceDetails, err := db_service.GetServiceInstanceDetailsById(instanceID)
	if err != nil {
		return brokerapi.UpdateServiceSpec{}, fmt.Errorf("error updating %q, %s", instanceID, err)
	}

	instanceDetails.PlanId = details.PlanID
	if err := db_service.SaveServiceInstanceDetails(instanceDetails); err != nil {
		return brokerapi.UpdateServiceSpec{}, fmt.Errorf("updating the existing database record %s", err)
	}

	return brokerapi.UpdateServiceSpec{}, nil
}
