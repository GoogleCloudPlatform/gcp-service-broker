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
	"fmt"

	"code.cloudfoundry.org/lager"
	"github.com/GoogleCloudPlatform/gcp-service-broker/db_service/models"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/toggles"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/varcontext"
	"github.com/GoogleCloudPlatform/gcp-service-broker/utils"
	"github.com/pivotal-cf/brokerapi"
	"golang.org/x/oauth2/jwt"
)

var enableCatalogSchemas = toggles.Features.Toggle("enable-catalog-schemas", false, `Enable generating JSONSchema for the service catalog.`)

// ServiceDefinition holds the necessary details to describe an OSB service and
// provision it.
type ServiceDefinition struct {
	Id               string `validate:"required,uuid"`
	Name             string `validate:"required,osbname"`
	Description      string
	DisplayName      string
	ImageUrl         string `validate:"omitempty,url"`
	DocumentationUrl string `validate:"omitempty,url"`
	SupportUrl       string `validate:"omitempty,url"`
	Tags             []string
	Bindable         bool
	PlanUpdateable   bool
	Plans            []ServicePlan

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

	// IsBuiltin is true if the service is built-in to the platform.
	IsBuiltin bool

	// config is set by the owning registry and is used to augment the data in the
	// definition.
	config ServiceConfig
}

// SetConfig sets user-defined overrides for the ServiceDefinition.
// It returns an error if the config is invalid for the service.
func (svc *ServiceDefinition) SetConfig(cfg ServiceConfig) error {
	// all plans should have valid values
	for idx, plan := range cfg.CustomPlans {
		if err := svc.validatePlan(idx, plan); err != nil {
			return err
		}
	}

	svc.config = cfg
	return nil
}

// CatalogEntry returns the service broker catalog entry for this service, it
// has metadata about the service so operators and programmers know which
// service and plan will work best for their purposes.
func (svc *ServiceDefinition) CatalogEntry() *Service {
	userPlans := []ServicePlan{}
	for _, customPlan := range svc.config.CustomPlans {
		userPlans = append(userPlans, customPlan.ToServicePlan())
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
		for i := range sd.Plans {
			sd.Plans[i].Schemas = svc.createSchemas()
		}
	}

	return sd
}

// createSchemas creates JSONSchemas compatible with the OSB spec for provision and bind.
// It leaves the instance update schema empty to indicate updates are not supported.
func (svc *ServiceDefinition) createSchemas() *brokerapi.ServiceSchemas {
	return &brokerapi.ServiceSchemas{
		Instance: brokerapi.ServiceInstanceSchema{
			Create: brokerapi.Schema{
				Parameters: createJsonSchema(svc.ProvisionInputVariables),
			},
		},
		Binding: brokerapi.ServiceBindingSchema{
			Create: brokerapi.Schema{
				Parameters: createJsonSchema(svc.BindInputVariables),
			},
		},
	}
}

// GetPlanById finds a plan in this service by its UUID.
func (svc *ServiceDefinition) GetPlanById(planId string) (*ServicePlan, error) {
	catalogEntry := svc.CatalogEntry()

	for _, plan := range catalogEntry.Plans {
		if plan.ID == planId {
			return &plan, nil
		}
	}

	return nil, fmt.Errorf("Plan ID %q could not be found", planId)
}

func (svc *ServiceDefinition) validatePlan(index int, plan CustomPlan) error {
	if plan.GUID == "" {
		return fmt.Errorf("%s custom_plans[%d] is missing an id", svc.Name, index)
	}

	if plan.Name == "" {
		return fmt.Errorf("%s custom_plans[%d] is missing a name", svc.Name, index)
	}

	if svc.PlanVariables == nil {
		return nil
	}

	for _, customVar := range svc.PlanVariables {
		if !customVar.Required {
			continue
		}

		if _, ok := plan.Properties[customVar.FieldName]; !ok {
			return fmt.Errorf("%s custom_plans[%d] is missing required property %s", svc.Name, index, customVar.FieldName)
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
	// The namespaces of these values roughly align with the OSB spec.
	constants := map[string]interface{}{
		"request.plan_id":        details.PlanID,
		"request.service_id":     details.ServiceID,
		"request.instance_id":    instanceId,
		"request.default_labels": utils.ExtractDefaultLabels(instanceId, details),
	}

	builder := varcontext.Builder().
		SetEvalConstants(constants).
		MergeMap(svc.config.ProvisionDefaults).
		MergeJsonObject(details.GetRawParameters()).
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
// 3. User defined variables (in `bind_input_variables`)
// 4. Operator default variables loaded from the environment.
// 5. Default variables (in `bind_input_variables`).
//
func (svc *ServiceDefinition) BindVariables(instance models.ServiceInstanceDetails, bindingID string, details brokerapi.BindDetails) (*varcontext.VarContext, error) {
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
		"request.plan_id":    instance.PlanId,
		"request.service_id": instance.ServiceId,
		"request.app_guid":   appGuid,

		// specified by the existing instance
		"instance.name":    instance.Name,
		"instance.details": otherDetails,
	}

	builder := varcontext.Builder().
		SetEvalConstants(constants).
		MergeMap(svc.config.BindDefaults).
		MergeJsonObject(details.GetRawParameters()).
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
