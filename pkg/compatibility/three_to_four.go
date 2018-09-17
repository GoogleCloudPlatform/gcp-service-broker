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
	"fmt"

	"github.com/GoogleCloudPlatform/gcp-service-broker/db_service"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/broker"
	"github.com/pivotal-cf/brokerapi"
)

type upgradePath struct {
	ServiceId      string
	LegacyPlanId   string
	LegacyPlanName string
	NewPlanId      string
	NewPlanName    string
}

func (u *upgradePath) ToServicePlan() brokerapi.ServicePlan {
	return brokerapi.ServicePlan{
		ID:          u.LegacyPlanId,
		Name:        fmt.Sprintf("legacy3-%s", u.LegacyPlanName),
		Description: fmt.Sprintf("Legacy plan, must be upgraded to %q", u.NewPlanName),
	}
}

// NewLegacyPlanUpgrader wraps a service broker with an interface that requires
// provisioned services which are instances of legacy plans (which had GUIDs
// generated at runtime rather than fixed and are therefore different per-install)
// to upgrade to their fixed counterpart before any operations can be done on
// them.
func NewLegacyPlanUpgrader(wrapped brokerapi.ServiceBroker) *ThreeToFour {
	legacyPlanUpgrades := []upgradePath{
		{"83837945-1547-41e0-b661-ea31d76eed11", "", "default", "10866183-a775-49e8-96e3-4e7a901e4a79", "default"},                           // Stackdriver Debugger
		{"5ad2dce0-51f7-4ede-8b46-293d6df1e8d4", "", "default", "be7954e1-ecfb-4936-a0b6-db35e6424c7a", "default"},                           // Cloud ML APIs
		{"628629e3-79f5-4255-b981-d14c6c7856be", "", "default", "622f4da3-8731-492a-af29-66a9146f8333", "default"},                           // Pub/Sub
		{"c5ddfe15-24d9-47f8-8ffe-f6b7daa9cf4a", "", "default", "ab6c2287-b4bc-4ff4-a36a-0575e7910164", "default"},                           // Stackdriver Trace
		{"76d4abb2-fee7-4c8f-aee1-bcea2837f02b", "", "default", "05f1fb6b-b5f0-48a2-9c2b-a5f236507a97", "default"},                           // Datastore
		{"f80c0a3e-bd4d-4809-a900-b4e33a6450f1", "", "default", "10ff4e72-6e84-44eb-851f-bdb38a791914", "default"},                           // BigQuery
		{"b9e4332e-b42b-4680-bda5-ea1506797474", "", "nearline", "a42c1182-d1a0-4d40-82c1-28220518b360", "nearline"},                         // Cloud Storage
		{"b9e4332e-b42b-4680-bda5-ea1506797474", "", "reduced_availability", "1a1f4fe6-1904-44d0-838c-4c87a9490a6b", "reduced-availability"}, // Cloud Storage
		{"b9e4332e-b42b-4680-bda5-ea1506797474", "", "standard", "e1d11f65-da66-46ad-977c-6d56513baf43", "standard"},                         // Cloud Storage
	}

	var allowedUpgrades []upgradePath
	for _, upgrade := range legacyPlanUpgrades {
		// Check each upgrade in the DB, if it does not exist there then the user
		// didn't get it from one of their earlier migrations.
		legacyDefinition, err := db_service.GetPlanDetailsV1ByServiceIdAndName(upgrade.ServiceId, upgrade.LegacyPlanName)
		if err != nil {
			continue // continue on missing, users may not have all legacy plans
		}

		upgrade.LegacyPlanId = legacyDefinition.ID
		allowedUpgrades = append(allowedUpgrades, upgrade)
	}

	return &ThreeToFour{Wrapped: wrapped, allowedUpgrades: allowedUpgrades}
}

// ThreetoFour is a brokerapi.ServiceBroker wrapper that tells users when
// they are using legacy plan IDs and how to upgrade them.
type ThreeToFour struct {
	Wrapped         brokerapi.ServiceBroker
	allowedUpgrades []upgradePath
}

// LastOperation calls the wrapped service broker
func (t *ThreeToFour) LastOperation(ctx context.Context, instanceID, operationData string) (brokerapi.LastOperation, error) {
	return t.Wrapped.LastOperation(ctx, instanceID, operationData)
}

// Provision calls the wrapped service broker unless the plan is legacy, in which case
// the user is told which plan to use instead.
func (t *ThreeToFour) Provision(ctx context.Context, instanceID string, details brokerapi.ProvisionDetails, asyncAllowed bool) (brokerapi.ProvisionedServiceSpec, error) {
	if up, ok := t.getUpgradePath(details.ServiceID, details.PlanID); ok {
		return brokerapi.ProvisionedServiceSpec{}, fmt.Errorf("The plan %q is only availble for compatibility purposes, use %q instead.", "legacy3-"+up.LegacyPlanName, up.NewPlanName)
	}

	return t.Wrapped.Provision(ctx, instanceID, details, asyncAllowed)
}

// Deprovision calls the wrapped service broker unless the plan is legacy, in which case
// the user is told how to upgrade first.
func (t *ThreeToFour) Deprovision(ctx context.Context, instanceID string, details brokerapi.DeprovisionDetails, asyncAllowed bool) (brokerapi.DeprovisionServiceSpec, error) {
	if err := t.migrationErrorMessage("deprovision", instanceID); err != nil {
		return brokerapi.DeprovisionServiceSpec{}, err
	}

	return t.Wrapped.Deprovision(ctx, instanceID, details, asyncAllowed)
}

// Bind calls the wrapped service broker unless the plan is legacy, in which case
// the user is told how to upgrade first.
func (t *ThreeToFour) Bind(ctx context.Context, instanceID, bindingID string, details brokerapi.BindDetails) (brokerapi.Binding, error) {
	if err := t.migrationErrorMessage("bind", instanceID); err != nil {
		return brokerapi.Binding{}, err
	}

	return t.Wrapped.Bind(ctx, instanceID, bindingID, details)
}

// Unbind calls the wrapped service broker unless the plan is legacy, in which case
// the user is told how to upgrade first.
func (t *ThreeToFour) Unbind(ctx context.Context, instanceID, bindingID string, details brokerapi.UnbindDetails) error {
	if err := t.migrationErrorMessage("unbind", instanceID); err != nil {
		return err
	}

	return t.Wrapped.Unbind(ctx, instanceID, bindingID, details)
}

// Services returns the list of enabled services, with dummy services injected
// for legacy compatibility.
func (t *ThreeToFour) Services(ctx context.Context) ([]brokerapi.Service, error) {
	services := broker.GetEnabledServices()

	var marketplace []brokerapi.Service

	for _, svc := range services {
		catalog, err := svc.CatalogEntry()
		if err != nil {
			return nil, err
		}
		serviceDefn := t.augmentServiceCatalog(catalog.ToPlain())
		marketplace = append(marketplace, serviceDefn)
	}

	return marketplace, nil
}

func (t *ThreeToFour) augmentServiceCatalog(entry brokerapi.Service) brokerapi.Service {
	serviceGuid := entry.ID

	var compatPlans []brokerapi.ServicePlan
	for _, upgradePath := range t.allowedUpgrades {
		if upgradePath.ServiceId == serviceGuid {
			compatPlans = append(compatPlans, upgradePath.ToServicePlan())
		}
	}

	if len(compatPlans) > 0 {
		entry.PlanUpdatable = true
		entry.Plans = append(entry.Plans, compatPlans...)
	}

	return entry
}

// Update checks if the update is for a plan change from legacy to the defined
// acceptable upgrade plan and modifies the database if so.
func (t *ThreeToFour) Update(ctx context.Context, instanceID string, details brokerapi.UpdateDetails, asyncAllowed bool) (brokerapi.UpdateServiceSpec, error) {
	instanceDetails, err := db_service.GetServiceInstanceDetailsById(instanceID)
	if err != nil {
		return brokerapi.UpdateServiceSpec{}, brokerapi.ErrInstanceDoesNotExist
	}

	path, ok := t.getUpgradePath(instanceDetails.ServiceId, instanceDetails.PlanId)
	if !ok {
		return t.Wrapped.Update(ctx, instanceID, details, asyncAllowed)
	}

	if path.NewPlanId != details.PlanID {
		return brokerapi.UpdateServiceSpec{}, fmt.Errorf("you can only upgrade this legacy plan to %q", path.NewPlanName)
	}

	return brokerapi.UpdateServiceSpec{}, t.updatePlanId(instanceID, details.PlanID)
}

func (t *ThreeToFour) getUpgradePath(serviceId, planId string) (upgradePath, bool) {
	for _, up := range t.allowedUpgrades {
		if up.ServiceId == serviceId && planId == up.LegacyPlanId {
			return up, true
		}
	}

	return upgradePath{}, false
}

func (t *ThreeToFour) migrationErrorMessage(verb, instanceId string) error {
	service, err := db_service.GetServiceInstanceDetailsById(instanceId)
	if err != nil {
		return brokerapi.ErrInstanceDoesNotExist
	}

	path, ok := t.getUpgradePath(service.ServiceId, service.PlanId)
	if !ok {
		return nil
	}

	command := fmt.Sprintf("cf update-service SERVICE_NAME -p %s", path.NewPlanName)
	return fmt.Errorf("The instance you're trying to %s is using an unsupported plan. You must update it first by running `%s`", verb, command)
}

func (ThreeToFour) updatePlanId(instanceID, newPlanId string) error {
	instanceDetails, err := db_service.GetServiceInstanceDetailsById(instanceID)
	if err != nil {
		return fmt.Errorf("couldn't find instance %q: %s", instanceID, err)
	}

	instanceDetails.PlanId = newPlanId

	if err := db_service.SaveServiceInstanceDetails(instanceDetails); err != nil {
		return fmt.Errorf("updating the existing database record %s", err)
	}

	return nil
}
