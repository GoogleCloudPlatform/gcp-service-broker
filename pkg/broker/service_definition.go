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

package broker

import (
	"encoding/json"
	"fmt"
	"strings"

	"code.cloudfoundry.org/lager"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/models"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/varcontext"
	"github.com/GoogleCloudPlatform/gcp-service-broker/utils"
	"github.com/pivotal-cf/brokerapi"
	"github.com/spf13/viper"
	"golang.org/x/oauth2/jwt"
)

// ServiceDefinition holds the necessary details to describe an OSB service and
// provision it.
type ServiceDefinition struct {
	Name                       string                       `validate:"osbname"`
	DefaultServiceDefinition   string                       `validate:"json"`
	ProvisionInputVariables    []BrokerVariable             `validate:"dive"`
	ProvisionComputedVariables []varcontext.DefaultVariable `validate:"dive"`
	BindInputVariables         []BrokerVariable             `validate:"dive"`
	BindOutputVariables        []BrokerVariable             `validate:"dive"`
	BindComputedVariables      []varcontext.DefaultVariable `validate:"dive"`
	PlanVariables              []BrokerVariable             `validate:"dive"`
	Examples                   []ServiceExample             `validate:"dive"`
	DefaultRoleWhitelist       []string

	// ProviderBuilder creates a new provider given the project, auth, and logger.
	ProviderBuilder func(projectId string, auth *jwt.Config, logger lager.Logger) ServiceProvider
}

// EnabledProperty computes the Viper property name for the boolean the user
// can set to disable or enable this service.
func (svc *ServiceDefinition) EnabledProperty() string {
	return fmt.Sprintf("service.%s.enabled", svc.Name)
}

// DefinitionProperty computes the Viper property name for the JSON service
// definition.
func (svc *ServiceDefinition) DefinitionProperty() string {
	return fmt.Sprintf("service.%s.definition", svc.Name)
}

// UserDefinedPlansProperty computes the Viper property name for the JSON list
// of user-defined service plans.
func (svc *ServiceDefinition) UserDefinedPlansProperty() string {
	return fmt.Sprintf("service.%s.plans", svc.Name)
}

// ProvisionDefaultOverrideProperty returns the Viper property name for the
// object users can set to override the default values on provision.
func (svc *ServiceDefinition) ProvisionDefaultOverrideProperty() string {
	return fmt.Sprintf("service.%s.provision.defaults", svc.Name)
}

// ProvisionDefaultOverrides returns the deserialized JSON object for the
// operator-provided property overrides.
func (svc *ServiceDefinition) ProvisionDefaultOverrides() map[string]interface{} {
	return viper.GetStringMap(svc.ProvisionDefaultOverrideProperty())
}

// IsRoleWhitelistEnabled returns false if the service has no default whitelist
// meaning it does not allow any roles.
func (svc *ServiceDefinition) IsRoleWhitelistEnabled() bool {
	return len(svc.DefaultRoleWhitelist) > 0
}

// BindDefaultOverrideProperty returns the Viper property name for the
// object users can set to override the default values on bind.
func (svc *ServiceDefinition) BindDefaultOverrideProperty() string {
	return fmt.Sprintf("service.%s.bind.defaults", svc.Name)
}

// BindDefaultOverrides returns the deserialized JSON object for the
// operator-provided property overrides.
func (svc *ServiceDefinition) BindDefaultOverrides() map[string]interface{} {
	return viper.GetStringMap(svc.BindDefaultOverrideProperty())
}

// TileUserDefinedPlansVariable returns the name of the user defined plans
// variable for the broker tile.
func (svc *ServiceDefinition) TileUserDefinedPlansVariable() string {
	prefix := "GOOGLE_"

	v := utils.PropertyToEnvUnprefixed(svc.Name)
	if strings.HasPrefix(v, prefix) {
		v = v[len(prefix):]
	}

	return v + "_CUSTOM_PLANS"
}

// IsEnabled returns false if the operator has explicitly disabled this service
// or true otherwise.
func (svc *ServiceDefinition) IsEnabled() bool {
	return viper.GetBool(svc.EnabledProperty())
}

// CatalogEntry returns the service broker catalog entry for this service, it
// has metadata about the service so operators and programmers know which
// service and plan will work best for their purposes.
func (svc *ServiceDefinition) CatalogEntry() (*Service, error) {
	sd, err := svc.ServiceDefinition()
	if err != nil {
		return nil, err
	}

	plans, err := svc.UserDefinedPlans()
	if err != nil {
		return nil, err
	}

	sd.Plans = append(sd.Plans, plans...)

	return sd, nil
}

// GetPlanById finds a plan in this service by its UUID.
func (svc *ServiceDefinition) GetPlanById(planId string) (*ServicePlan, error) {
	catalogEntry, err := svc.CatalogEntry()
	if err != nil {
		return nil, err
	}

	for _, plan := range catalogEntry.Plans {
		if plan.ID == planId {
			return &plan, nil
		}
	}

	return nil, fmt.Errorf("Plan ID %q could not be found", planId)
}

// UserDefinedPlans extracts user defined plans from the environment, failing if
// the plans were not valid JSON or were missing required properties/variables.
func (svc *ServiceDefinition) UserDefinedPlans() ([]ServicePlan, error) {
	plans := []ServicePlan{}

	userPlanJson := viper.GetString(svc.UserDefinedPlansProperty())
	if userPlanJson == "" {
		return plans, nil
	}

	// There's a mismatch between how plans are used internally and defined by
	// the user and the tile. In the environment variables we parse an array of
	// flat maps, but internally extra variables need to be put into a sub-map.
	// e.g. they come in as [{"id":"1234", "name":"foo", "A": 1, "B": 2}]
	// but we need [{"id":"1234", "name":"foo", "service_properties":{"A": 1, "B": 2}}]
	// Go doesn't support this natively so we do it manually here.
	rawPlans := []json.RawMessage{}
	if err := json.Unmarshal([]byte(userPlanJson), &rawPlans); err != nil {
		return plans, err
	}

	for _, rawPlan := range rawPlans {
		plan := ServicePlan{}
		remainder, err := utils.UnmarshalObjectRemainder(rawPlan, &plan)
		if err != nil {
			return []ServicePlan{}, err
		}

		plan.ServiceProperties = make(map[string]string)
		if err := json.Unmarshal(remainder, &plan.ServiceProperties); err != nil {
			return []ServicePlan{}, err
		}

		if err := svc.validatePlan(plan); err != nil {
			return []ServicePlan{}, err
		}

		plans = append(plans, plan)
	}

	return plans, nil
}

func (svc *ServiceDefinition) validatePlan(plan ServicePlan) error {
	if plan.ID == "" {
		return fmt.Errorf("%s custom plan %+v is missing an id", svc.Name, plan)
	}

	if plan.Name == "" {
		return fmt.Errorf("%s custom plan %+v is missing a name", svc.Name, plan)
	}

	if svc.PlanVariables == nil {
		return nil
	}

	for _, customVar := range svc.PlanVariables {
		if !customVar.Required {
			continue
		}

		if _, ok := plan.ServiceProperties[customVar.FieldName]; !ok {
			return fmt.Errorf("%s custom plan %+v is missing required property %s", svc.Name, plan, customVar.FieldName)
		}
	}

	return nil
}

// ServiceDefinition extracts service definition from the environment, failing
// if the definition was not valid JSON.
func (svc *ServiceDefinition) ServiceDefinition() (*Service, error) {
	jsonDefinition := viper.GetString(svc.DefinitionProperty())
	if jsonDefinition == "" {
		jsonDefinition = svc.DefaultServiceDefinition
	}

	var defn Service
	err := json.Unmarshal([]byte(jsonDefinition), &defn)
	if err != nil {
		return nil, fmt.Errorf("Error parsing service definition for %q: %s", svc.Name, err)
	}
	return &defn, err
}

func (svc *ServiceDefinition) provisionDefaults() []varcontext.DefaultVariable {
	var out []varcontext.DefaultVariable
	for _, provisionVar := range svc.ProvisionInputVariables {
		out = append(out, varcontext.DefaultVariable{Name: provisionVar.FieldName, Default: provisionVar.Default, Overwrite: false, Type: string(provisionVar.Type)})
	}
	return out
}

func (svc *ServiceDefinition) bindDefaults() []varcontext.DefaultVariable {
	var out []varcontext.DefaultVariable
	for _, v := range svc.BindInputVariables {
		out = append(out, varcontext.DefaultVariable{Name: v.FieldName, Default: v.Default, Overwrite: false, Type: string(v.Type)})
	}
	return out
}

// ProvisionVariables gets the variable resolution context for a provision request.
// Variables have a very specific resolution order, and this function populates the context to preserve that.
// The variable resolution order is the following:
//
// 1. Variables defined in your `computed_variables` JSON list.
// 2. Variables defined by the selected service plan in its `service_properties` map.
// 3. User defined variables (in `provision_input_variables` or `bind_input_variables`)
// 4. Operator default variables loaded from the environment.
// 5. Default variables (in `provision_input_variables` or `bind_input_variables`).
//
// Loading into the map occurs slightly differently.
// Default variables and computed_variables get executed by interpolation.
// User defined varaibles are not to prevent side-channel attacks.
// Default variables may reference user provided variables.
// For example, to create a default database name based on a user-provided instance name.
// Therefore, they get executed conditionally if a user-provided variable does not exist.
// Computed variables get executed either unconditionally or conditionally for greater flexibility.
func (svc *ServiceDefinition) ProvisionVariables(instanceId string, details brokerapi.ProvisionDetails, plan ServicePlan) (*varcontext.VarContext, error) {
	defaults := svc.provisionDefaults()

	// The namespaces of these values roughly align with the OSB spec.
	constants := map[string]interface{}{
		"request.plan_id":        details.PlanID,
		"request.service_id":     details.ServiceID,
		"request.instance_id":    instanceId,
		"request.default_labels": utils.ExtractDefaultLabels(instanceId, details),
	}

	return varcontext.Builder().
		SetEvalConstants(constants).
		MergeMap(svc.ProvisionDefaultOverrides()).
		MergeJsonObject(details.GetRawParameters()).
		MergeDefaults(defaults).
		MergeMap(plan.GetServiceProperties()).
		MergeDefaults(svc.ProvisionComputedVariables).
		Build()
}

// BindVariables gets the variable resolution context for a bind request.
// Variables have a very specific resolution order, and this function populates the context to preserve that.
// The variable resolution order is the following:
//
// 1. Variables defined in your `computed_variables` JSON list.
// 3. User defined variables (in `bind_input_variables`)
// 4. Operator default variables loaded from the environment.
// 5. Default variables (in `bind_input_variables`).
//
func (svc *ServiceDefinition) BindVariables(instance models.ServiceInstanceDetails, bindingID string, details brokerapi.BindDetails) (*varcontext.VarContext, error) {
	defaults := svc.bindDefaults()

	otherDetails := make(map[string]interface{})
	if instance.OtherDetails != "" {
		if err := json.Unmarshal(json.RawMessage(instance.OtherDetails), &otherDetails); err != nil {
			return nil, err
		}
	}

	appGuid := ""
	if details.BindResource != nil {
		appGuid = details.BindResource.AppGuid
	}

	// The namespaces of these values roughly align with the OSB spec.
	constants := map[string]interface{}{
		// specified in the URL
		"request.binding_id":  bindingID,
		"request.instance_id": instance.ID,

		// specified in the request body
		// Note: the value in instance is considered the official record so values
		// are pulled from there rather than the request. In a future version of OSB
		// the duplicate sending of fields is likely to be removed.
		"request.plan_id":    instance.PlanId,
		"request.service_id": instance.ServiceId,
		"request.app_guid":   appGuid,

		// specified by the existing instance
		"instance.name":    instance.Name,
		"instance.details": otherDetails,
	}

	return varcontext.Builder().
		SetEvalConstants(constants).
		MergeMap(svc.BindDefaultOverrides()).
		MergeJsonObject(details.GetRawParameters()).
		MergeDefaults(defaults).
		MergeDefaults(svc.BindComputedVariables).
		Build()
}
