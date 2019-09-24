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
	"os"
	"strings"

	"code.cloudfoundry.org/lager"
	"github.com/GoogleCloudPlatform/gcp-service-broker/db_service/models"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/toggles"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/validation"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/varcontext"
	"github.com/GoogleCloudPlatform/gcp-service-broker/utils"
	"github.com/pivotal-cf/brokerapi"
	"github.com/spf13/viper"
	"golang.org/x/oauth2/jwt"
)

var enableCatalogSchemas = toggles.Features.Toggle("enable-catalog-schemas", false, `Enable generating JSONSchema for the service catalog.`)

// ServiceDefinition holds the necessary details to describe an OSB service and
// provision it.
type ServiceDefinition struct {
	Id               string
	Name             string
	Description      string
	DisplayName      string
	ImageUrl         string
	DocumentationUrl string
	SupportUrl       string
	Tags             []string
	Bindable         bool
	PlanUpdateable   bool
	Plans            []ServicePlan

	ProvisionInputVariables    []BrokerVariable
	ProvisionComputedVariables []varcontext.DefaultVariable
	BindInputVariables         []BrokerVariable
	BindOutputVariables        []BrokerVariable
	BindComputedVariables      []varcontext.DefaultVariable
	PlanVariables              []BrokerVariable
	Examples                   []ServiceExample
	DefaultRoleWhitelist       []string

	// ProviderBuilder creates a new provider given the project, auth, and logger.
	ProviderBuilder func(projectId string, auth *jwt.Config, logger lager.Logger) ServiceProvider

	// IsBuiltin is true if the service is built-in to the platform.
	IsBuiltin bool
}

var _ validation.Validatable = (*ServiceDefinition)(nil)

// Validate implements validation.Validatable.
func (sd *ServiceDefinition) Validate() (errs *validation.FieldError) {
	errs = errs.Also(
		validation.ErrIfNotUUID(sd.Id, "Id"),
		validation.ErrIfNotOSBName(sd.Name, "Name"),
	)

	if sd.ImageUrl != "" {
		errs = errs.Also(validation.ErrIfNotURL(sd.ImageUrl, "ImageUrl"))
	}

	if sd.DocumentationUrl != "" {
		errs = errs.Also(validation.ErrIfNotURL(sd.DocumentationUrl, "DocumentationUrl"))
	}

	if sd.SupportUrl != "" {
		errs = errs.Also(validation.ErrIfNotURL(sd.SupportUrl, "SupportUrl"))
	}

	for i, v := range sd.ProvisionInputVariables {
		errs = errs.Also(v.Validate().ViaFieldIndex("ProvisionInputVariables", i))
	}

	for i, v := range sd.ProvisionComputedVariables {
		errs = errs.Also(v.Validate().ViaFieldIndex("ProvisionComputedVariables", i))
	}

	for i, v := range sd.BindInputVariables {
		errs = errs.Also(v.Validate().ViaFieldIndex("BindInputVariables", i))
	}

	for i, v := range sd.BindOutputVariables {
		errs = errs.Also(v.Validate().ViaFieldIndex("BindOutputVariables", i))
	}

	for i, v := range sd.BindComputedVariables {
		errs = errs.Also(v.Validate().ViaFieldIndex("BindComputedVariables", i))
	}

	for i, v := range sd.PlanVariables {
		errs = errs.Also(v.Validate().ViaFieldIndex("PlanVariables", i))
	}

	return errs
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
	v := utils.PropertyToEnvUnprefixed(svc.Name)
	v = strings.TrimPrefix(v, "GOOGLE_")

	return v + "_CUSTOM_PLANS"
}

// CatalogEntry returns the service broker catalog entry for this service, it
// has metadata about the service so operators and programmers know which
// service and plan will work best for their purposes.
func (svc *ServiceDefinition) CatalogEntry() (*Service, error) {
	userPlans, err := svc.UserDefinedPlans()
	if err != nil {
		return nil, err
	}

	sd := &Service{
		Service: brokerapi.Service{
			ID:          svc.Id,
			Name:        svc.Name,
			Description: svc.Description,
			Metadata: &brokerapi.ServiceMetadata{
				DisplayName:     svc.DisplayName,
				LongDescription: svc.Description,

				DocumentationUrl: svc.DocumentationUrl,
				ImageUrl:         svc.ImageUrl,
				SupportUrl:       svc.SupportUrl,
			},
			Tags:          svc.Tags,
			Bindable:      svc.Bindable,
			PlanUpdatable: svc.PlanUpdateable,
		},
		Plans: append(svc.Plans, userPlans...),
	}

	if enableCatalogSchemas.IsActive() {
		for i, _ := range sd.Plans {
			sd.Plans[i].Schemas = svc.createSchemas()
		}
	}

	return sd, nil
}

// createSchemas creates JSONSchemas compatible with the OSB spec for provision and bind.
// It leaves the instance update schema empty to indicate updates are not supported.
func (svc *ServiceDefinition) createSchemas() *brokerapi.ServiceSchemas {
	return &brokerapi.ServiceSchemas{
		Instance: brokerapi.ServiceInstanceSchema{
			Create: brokerapi.Schema{
				Parameters: CreateJsonSchema(svc.ProvisionInputVariables),
			},
		},
		Binding: brokerapi.ServiceBindingSchema{
			Create: brokerapi.Schema{
				Parameters: CreateJsonSchema(svc.BindInputVariables),
			},
		},
	}
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

	// There's a mismatch between how plans are used internally and defined by
	// the user and the tile. In the environment variables we parse an array of
	// flat maps, but internally extra variables need to be put into a sub-map.
	// e.g. they come in as [{"id":"1234", "name":"foo", "A": 1, "B": 2}]
	// but we need [{"id":"1234", "name":"foo", "service_properties":{"A": 1, "B": 2}}]
	// Go doesn't support this natively so we do it manually here.
	rawPlans := []json.RawMessage{}

	// Unmarshal the plans from the viper configuration which is just a JSON list
	// of plans
	if userPlanJSON := viper.GetString(svc.UserDefinedPlansProperty()); userPlanJSON != "" {
		if err := json.Unmarshal([]byte(userPlanJSON), &rawPlans); err != nil {
			return []ServicePlan{}, err
		}
	}

	// Unmarshal tile plans if they're included, which are a JSON object where
	// keys are
	if tilePlans := os.Getenv(svc.TileUserDefinedPlansVariable()); tilePlans != "" {
		var rawTilePlans map[string]json.RawMessage
		if err := json.Unmarshal([]byte(tilePlans), &rawTilePlans); err != nil {
			return []ServicePlan{}, err
		}

		for _, v := range rawTilePlans {
			rawPlans = append(rawPlans, v)
		}
	}

	plans := []ServicePlan{}
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

		// reading from a tile we need to move their GUID to an ID field
		if plan.ID == "" {
			plan.ID = plan.ServiceProperties["guid"]
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
// 3. Variables overridden in the plan's `provision_overrides` map.
// 4. User defined variables (in `provision_input_variables` or `bind_input_variables`)
// 5. Operator default variables loaded from the environment.
// 6. Default variables (in `provision_input_variables` or `bind_input_variables`).
//
// Loading into the map occurs slightly differently.
// Default variables and computed_variables get executed by interpolation.
// User defined varaibles are not to prevent side-channel attacks.
// Default variables may reference user provided variables.
// For example, to create a default database name based on a user-provided instance name.
// Therefore, they get executed conditionally if a user-provided variable does not exist.
// Computed variables get executed either unconditionally or conditionally for greater flexibility.
func (svc *ServiceDefinition) ProvisionVariables(instanceId string, details brokerapi.ProvisionDetails, plan ServicePlan) (*varcontext.VarContext, error) {
	// The namespaces of these values roughly align with the OSB spec.
	constants := map[string]interface{}{
		"request.plan_id":        details.PlanID,
		"request.service_id":     details.ServiceID,
		"request.instance_id":    instanceId,
		"request.default_labels": utils.ExtractDefaultLabels(instanceId, details),
	}

	builder := varcontext.Builder().
		SetEvalConstants(constants).
		MergeMap(svc.ProvisionDefaultOverrides()).
		MergeJsonObject(details.GetRawParameters()).
		MergeMap(plan.ProvisionOverrides).
		MergeDefaults(svc.provisionDefaults()).
		MergeMap(plan.GetServiceProperties()).
		MergeDefaults(svc.ProvisionComputedVariables)

	return buildAndValidate(builder, svc.ProvisionInputVariables)
}

// BindVariables gets the variable resolution context for a bind request.
// Variables have a very specific resolution order, and this function populates the context to preserve that.
// The variable resolution order is the following:
//
// 1. Variables defined in your `computed_variables` JSON list.
// 2. Variables overridden in the plan's `bind_overrides` map.
// 3. User defined variables (in `bind_input_variables`)
// 4. Operator default variables loaded from the environment.
// 5. Default variables (in `bind_input_variables`).
//
func (svc *ServiceDefinition) BindVariables(instance models.ServiceInstanceDetails, bindingID string, details brokerapi.BindDetails, plan *ServicePlan) (*varcontext.VarContext, error) {
	otherDetails := make(map[string]interface{})
	if err := instance.GetOtherDetails(&otherDetails); err != nil {
		return nil, err
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
		"request.plan_id":         instance.PlanId,
		"request.service_id":      instance.ServiceId,
		"request.app_guid":        appGuid,
		"request.plan_properties": plan.GetServiceProperties(),

		// specified by the existing instance
		"instance.name":    instance.Name,
		"instance.details": otherDetails,
	}

	builder := varcontext.Builder().
		SetEvalConstants(constants).
		MergeMap(svc.BindDefaultOverrides()).
		MergeJsonObject(details.GetRawParameters()).
		MergeMap(plan.BindOverrides).
		MergeDefaults(svc.bindDefaults()).
		MergeDefaults(svc.BindComputedVariables)

	return buildAndValidate(builder, svc.BindInputVariables)
}

// buildAndValidate builds the varcontext and if it's valid validates the
// resulting context against the JSONSchema defined by the BrokerVariables
// exactly one of VarContext and error will be nil upon return.
func buildAndValidate(builder *varcontext.ContextBuilder, vars []BrokerVariable) (*varcontext.VarContext, error) {
	vc, err := builder.Build()
	if err != nil {
		return nil, err
	}

	if err := ValidateVariables(vc.ToMap(), vars); err != nil {
		return nil, err
	}

	return vc, nil
}
