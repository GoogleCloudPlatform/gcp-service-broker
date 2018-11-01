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

package generator

import (
	"fmt"
	"log"
	"strings"

	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/account_managers"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/broker"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/toggles"
	"github.com/GoogleCloudPlatform/gcp-service-broker/utils"
	yaml "gopkg.in/yaml.v2"
)

// TileFormsSections holds the top level fields in tile.yml responsible for
// the forms.
// https://docs.pivotal.io/tiledev/2-2/tile-structure.html
type TileFormsSections struct {
	Forms            []Form `yaml:"forms"`
	ServicePlanForms []Form `yaml:"service_plan_forms,omitempty"`
}

// Form is a PCF Ops Manager compatible form definition used to generate forms.
// See https://docs.pivotal.io/tiledev/2-2/product-template-reference.html#form-properties
// for details about the fields.
type Form struct {
	Name        string         `yaml:"name"`
	Label       string         `yaml:"label"`
	Description string         `yaml:"description"`
	Optional    bool           `yaml:"optional,omitempty"` // optional, default false
	Properties  []FormProperty `yaml:"properties"`
}

// FormOption is an enumerated element for FormProperties that can be selected
// from. Name is the value and label is the human-readable display name.
type FormOption struct {
	Name  string `yaml:"name"`
	Label string `yaml:"label"`
}

// FormProperty holds a single form element in a PCF Ops manager form.
type FormProperty struct {
	Name         string       `yaml:"name"`
	Type         string       `yaml:"type"`
	Default      interface{}  `yaml:"default,omitempty"`
	Label        string       `yaml:"label,omitempty"`
	Description  string       `yaml:"description,omitempty"`
	Configurable bool         `yaml:"configurable,omitempty"` // optional, default false
	Options      []FormOption `yaml:"options,omitempty"`
	Optional     bool         `yaml:"optional,omitempty"` // optional, default false
}

// GenerateFormsString creates all the forms for the user to fill out in the PCF tile
// and returns it as a string.
func GenerateFormsString() string {
	response, err := yaml.Marshal(GenerateForms())
	if err != nil {
		log.Fatalf("Error marshaling YAML: %s", err)
	}

	return string(response)
}

// GenerateForms creates all the forms for the user to fill out in the PCF tile.
func GenerateForms() TileFormsSections {
	// Add new forms at the bottom of the list because the order is reflected
	// in the generated UI and we don't want to mix things up on users.
	return TileFormsSections{
		Forms: []Form{
			generateServiceAccountForm(),
			generateDatabaseForm(),
			generateEnableDisableForm(),
			generateRoleWhitelistForm(),
			generateCompatibilityForm(),
			generateDefaultOverrideForm(),
		},

		ServicePlanForms: generateServicePlanForms(),
	}
}

// generateEnableDisableForm generates the form to enable and disable services.
func generateEnableDisableForm() Form {
	enablers := []FormProperty{}
	for _, svc := range broker.GetAllServices() {
		entry, err := svc.CatalogEntry()
		if err != nil {
			log.Fatalf("Error getting catalog entry for service %s, %v", svc.Name, err)
		}

		enableForm := FormProperty{
			Name:         strings.ToLower(utils.PropertyToEnv(svc.EnabledProperty())),
			Label:        fmt.Sprintf("Let the broker create and bind %s instances.", entry.Metadata.DisplayName),
			Type:         "boolean",
			Default:      true,
			Configurable: true,
		}

		enablers = append(enablers, enableForm)
	}

	return Form{
		Name:        "enable_disable",
		Label:       "Enable Services",
		Description: "Enable or disable services.",
		Properties:  enablers,
	}
}

// generateRoleWhitelistForm generates a form for users to enable/disable the
// whitelist validation for new service accounts bound to the service.
// They are opt-out and on by default for safety.
func generateRoleWhitelistForm() Form {
	enablers := []FormProperty{}
	for _, svc := range broker.GetAllServices() {
		entry, err := svc.CatalogEntry()
		if err != nil {
			log.Fatalf("Error getting catalog entry for service %s, %v", svc.Name, err)
		}

		if !svc.IsRoleWhitelistEnabled() {
			continue
		}

		enableForm := FormProperty{
			Name:         strings.ToLower(utils.PropertyToEnv(account_managers.RoleWhitelistProperty(svc.Name))),
			Label:        fmt.Sprintf("Role whitelist for %s instances.", entry.Metadata.DisplayName),
			Description:  "A comma delimited list of roles (minus the role/ prefix) that can be used when creating bound users for this service.",
			Type:         "string",
			Default:      strings.Join(svc.DefaultRoleWhitelist, ","),
			Configurable: true,
		}

		enablers = append(enablers, enableForm)
	}

	return Form{
		Name:        "role_whitelists",
		Label:       "Role Whitelisting",
		Description: "Enable or disable role whitelisting.",
		Properties:  enablers,
	}
}

// generateDefaultOverrideForm generates a form for users to override the
// defaults in a plan.
func generateDefaultOverrideForm() Form {
	formElements := []FormProperty{}
	for _, svc := range broker.GetAllServices() {
		entry, err := svc.CatalogEntry()
		if err != nil {
			log.Fatalf("Error getting catalog entry for service %s, %v", svc.Name, err)
		}

		if !svc.IsRoleWhitelistEnabled() {
			continue
		}

		provisionForm := FormProperty{
			Name:         strings.ToLower(utils.PropertyToEnv(svc.ProvisionDefaultOverrideProperty())),
			Label:        fmt.Sprintf("Provision default override %s instances.", entry.Metadata.DisplayName),
			Description:  "A JSON object with key/value pairs. Keys MUST be the name of a user-defined provision property and values are the alternative default.",
			Type:         "text",
			Default:      "{}",
			Configurable: true,
		}
		formElements = append(formElements, provisionForm)

		bindForm := FormProperty{
			Name:         strings.ToLower(utils.PropertyToEnv(svc.BindDefaultOverrideProperty())),
			Label:        fmt.Sprintf("Bind default override %s instances.", entry.Metadata.DisplayName),
			Description:  "A JSON object with key/value pairs. Keys MUST be the name of a user-defined bind property and values are the alternative default.",
			Type:         "text",
			Default:      "{}",
			Configurable: true,
		}
		formElements = append(formElements, bindForm)
	}

	return Form{
		Name:        "default_override",
		Label:       "Default Overrides",
		Description: "Override the default values your users get when provisioning.",
		Properties:  formElements,
	}
}

// generateDatabaseForm generates the form for configuring database settings.
func generateDatabaseForm() Form {
	return Form{
		Name:        "database_properties",
		Label:       "Database Properties",
		Description: "Connection details for the backing database for the service broker.",
		Properties: []FormProperty{
			{Name: "db_host", Type: "string", Label: "Database host", Configurable: true},
			{Name: "db_username", Type: "string", Label: "Database username", Optional: true, Configurable: true},
			{Name: "db_password", Type: "secret", Label: "Database password", Optional: true, Configurable: true},
			{Name: "db_port", Type: "string", Label: "Database port (defaults to 3306)", Default: "3306", Configurable: true},
			{Name: "db_name", Type: "string", Label: "Database name", Default: "servicebroker", Configurable: true},
			{Name: "ca_cert", Type: "text", Label: "Server CA cert", Optional: true, Configurable: true},
			{Name: "client_cert", Type: "text", Label: "Client cert", Optional: true, Configurable: true},
			{Name: "client_key", Type: "text", Label: "Client key", Optional: true, Configurable: true},
		},
	}
}

// generateServiceAccountForm generates the form for configuring the service
// account.
func generateServiceAccountForm() Form {
	return Form{
		Name:        "root_service_account",
		Label:       "Root Service Account",
		Description: "Please paste in the contents of the json keyfile (un-encoded) for your service account with owner credentials.",
		Properties: []FormProperty{
			{Name: "root_service_account_json", Type: "text", Label: "Root Service Account JSON", Configurable: true},
		},
	}
}

func generateCompatibilityForm() Form {
	var formEntries []FormProperty

	for _, toggle := range toggles.Compatibility.Toggles() {
		toggleEntry := FormProperty{
			Name:         strings.ToLower(toggle.EnvironmentVariable()),
			Type:         "boolean",
			Label:        toggle.Name,
			Configurable: true,
			Optional:     true,
			Default:      fmt.Sprintf("%v", toggle.Default), // the tile deals with all values as strings so a default string is acceptable.
			Description:  singleLine(toggle.Description),
		}

		formEntries = append(formEntries, toggleEntry)
	}

	return Form{
		Name:        "compatibility",
		Label:       "Compatibility",
		Description: "Legacy Compatibility Options",
		Properties:  formEntries,
	}
}

// generateServicePlanForms generates customized service plan forms for all
// registered services that have the ability to customize their variables.
func generateServicePlanForms() []Form {
	out := []Form{}

	for _, svc := range broker.GetAllServices() {
		planVars := svc.PlanVariables

		if planVars == nil || len(planVars) == 0 {
			continue
		}

		form, err := generateServicePlanForm(svc)
		if err != nil {
			log.Fatalf("Error generating form for %+v, %s", form, err)
		}

		out = append(out, form)
	}

	return out
}

// generateServicePlanForm creates a form for adding additional service plans
// to the broker for an existing service.
func generateServicePlanForm(svc *broker.ServiceDefinition) (Form, error) {
	entry, err := svc.CatalogEntry()
	if err != nil {
		return Form{}, err
	}

	displayName := entry.Metadata.DisplayName
	planForm := Form{
		Name:        strings.ToLower(svc.TileUserDefinedPlansVariable()),
		Description: fmt.Sprintf("Generate custom plans for %s.", displayName),
		Label:       fmt.Sprintf("%s Custom Plans", displayName),
		Optional:    true,
		Properties: []FormProperty{
			{
				Name:         "display_name",
				Label:        "Display Name",
				Type:         "string",
				Description:  "Name of the plan to be displayed to users.",
				Configurable: true,
			},
			{
				Name:         "description",
				Label:        "Plan description",
				Type:         "string",
				Description:  "The description of the plan shown to users.",
				Configurable: true,
			},
			{
				Name:        "service",
				Label:       "Service",
				Type:        "dropdown_select",
				Description: "The service this plan is associated with.",
				Default:     entry.ID,
				Options: []FormOption{
					{
						Name:  entry.ID,
						Label: displayName,
					},
				},
			},
		},
	}

	// Along with the above three fixed properties, each plan has optional
	// additional properties.

	for _, v := range svc.PlanVariables {
		prop := brokerVariableToFormProperty(v)
		planForm.Properties = append(planForm.Properties, prop)
	}

	return planForm, nil
}

func brokerVariableToFormProperty(v broker.BrokerVariable) FormProperty {
	formInput := FormProperty{
		Name:         v.FieldName,
		Label:        propertyToLabel(v.FieldName),
		Type:         string(v.Type),
		Description:  v.Details,
		Configurable: true,
		Optional:     !v.Required,
		Default:      v.Default,
	}

	if v.Enum != nil {
		formInput.Type = "dropdown_select"

		opts := []FormOption{}
		for name, label := range v.Enum {
			opts = append(opts, FormOption{Name: fmt.Sprintf("%v", name), Label: label})
		}

		formInput.Options = opts

		if len(opts) == 1 {
			formInput.Configurable = false
			formInput.Default = opts[0].Name
		}
	}

	return formInput
}

// propertyToLabel converts a JSON snake-case property into a title case
// human-readable alternative.
func propertyToLabel(property string) string {
	return strings.Title(strings.NewReplacer("_", " ").Replace(property))
}

func singleLine(text string) string {
	lines := strings.Split(text, "\n")

	var out []string
	for _, line := range lines {
		out = append(out, strings.TrimSpace(line))
	}

	return strings.Join(out, " ")
}
