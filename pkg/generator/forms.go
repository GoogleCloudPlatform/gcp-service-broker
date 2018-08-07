// Copyright the Service Broker Project Authors.
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

	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/broker"
	"github.com/GoogleCloudPlatform/gcp-service-broker/utils"
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
	Optional    bool           `yaml:"optional"`
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
	Type         string       `yaml:"type,omitempty"`
	Default      interface{}  `yaml:"default,omitempty"`
	Label        string       `yaml:"label,omitempty"`
	Description  string       `yaml:"description,omitempty"`
	Configurable bool         `yaml:"configurable"`
	Options      []FormOption `yaml:"options,omitempty"`
	Optional     bool         `yaml:"optional"`
}

// GenerateForms creates all the forms for the user to fill out in the PCF tile.
func GenerateForms() TileFormsSections {
	// Add new forms at the bottom of the list because the order is reflected
	// in the generated UI and we don't want to mix things up on users.
	return TileFormsSections{
		Forms: []Form{
			GenerateServiceAccountForm(),
			GenerateDatabaseForm(),
			GenerateEnableDisableForm(),
		},

		ServicePlanForms: GenerateServicePlanForms(),
	}
}

// GenerateEnableDisableForm generates the form to enable and disable services.
func GenerateEnableDisableForm() Form {
	enablers := []FormProperty{}
	for _, svc := range broker.GetAllServices() {
		entry, err := svc.CatalogEntry()
		if err != nil {
			log.Fatalf("Error getting catalog entry for service %s, %v", svc.Name, err)
		}

		enableForm := FormProperty{
			Name:         strings.ToLower(utils.PropertyToEnv(svc.EnabledProperty())),
			Label:        fmt.Sprintf("Let the broker create and bind %s instances", entry.Metadata.DisplayName),
			Type:         "boolean",
			Default:      true,
			Configurable: true,
			Optional:     true,
		}

		enablers = append(enablers, enableForm)
	}

	return Form{
		Name:        "enable_disable",
		Label:       "Enable Services",
		Description: "Enable or disable services",
		Properties:  enablers,
	}
}

// GenerateDatabaseForm generates the form for configuring database settings.
func GenerateDatabaseForm() Form {
	return Form{
		Name:        "database_properties",
		Label:       "Database Properties",
		Description: "Connection details for the backing database for the service broker",
		Properties: []FormProperty{
			FormProperty{Name: "db_host", Type: "string", Label: "Database host"},
			FormProperty{Name: "db_username", Type: "string", Label: "Database username", Optional: true},
			FormProperty{Name: "db_password", Type: "secret", Label: "Database password", Optional: true},
			FormProperty{Name: "db_port", Type: "string", Label: "Database port (defaults to 3306)", Default: "3306"},
			FormProperty{Name: "ca_cert", Type: "text", Label: "Server CA cert", Optional: true},
			FormProperty{Name: "client_cert", Type: "text", Label: "Client cert", Optional: true},
			FormProperty{Name: "client_key", Type: "text", Label: "Client key", Optional: true},
		},
	}
}

// GenerateServiceAccountForm generates the form for configuring the service
// account.
func GenerateServiceAccountForm() Form {
	return Form{
		Name:        "root_service_account",
		Label:       "Root Service Account",
		Description: "Please paste in the contents of the json keyfile (un-encoded) for your service account with owner credentials",
		Properties: []FormProperty{
			FormProperty{Name: "root_service_account_json", Type: "text", Label: "Root Service Account JSON"},
		},
	}
}

func GenerateServicePlanForms() []Form {
	out := []Form{}

	for _, svc := range broker.GetAllServices() {
		planVars := svc.PlanVariables

		if planVars == nil || len(planVars) == 0 {
			continue
		}

		form, err := GenerateServicePlanForm(svc)
		if err != nil {
			log.Fatalf("Error generating form for %+v, %s", form, err)
		}

		out = append(out, form)
	}

	return out
}

func GenerateServicePlanForm(svc *broker.BrokerService) (Form, error) {
	entry, err := svc.CatalogEntry()
	if err != nil {
		return Form{}, err
	}

	displayName := entry.Metadata.DisplayName
	planForm := Form{
		Name:        strings.ToLower(svc.TileUserDefinedPlansVariable()),
		Description: fmt.Sprintf("Generate custom plans for %s", displayName),
		Label:       fmt.Sprintf("%s Custom Plans", displayName),
		Optional:    true,
		Properties: []FormProperty{
			FormProperty{
				Name:         "display_name",
				Label:        "Display Name",
				Type:         "string",
				Description:  "Display name",
				Configurable: true,
			},
			FormProperty{
				Name:         "description",
				Label:        "Plan description",
				Type:         "string",
				Description:  "Plan description",
				Configurable: true,
			},
			FormProperty{
				Name:        "service",
				Label:       "Service",
				Type:        "dropdown_select",
				Description: "The service this plan is associated with",
				Options: []FormOption{
					FormOption{
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

func propertyToLabel(property string) string {
	return strings.Title(strings.NewReplacer("_", " ").Replace(property))
}
