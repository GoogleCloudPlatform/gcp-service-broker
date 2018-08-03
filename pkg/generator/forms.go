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
	"strings"

	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/broker"
	"github.com/GoogleCloudPlatform/gcp-service-broker/utils"
)

// https://docs.pivotal.io/tiledev/2-2/tile-structure.html

type TileFormsSections struct {
	Forms            []Form `yaml:"forms"`
	ServicePlanForms []Form `yaml:"service_plan_forms,omitempty"`
}

type Form struct {
	Name        string         `yaml:"name"`
	Label       string         `yaml:"label"`
	Description string         `yaml:"description"`
	Properties  []FormProperty `yaml:"properties"`
}

type FormOption struct {
	Name  string `yaml:"name"`
	Label string `yaml:"label"`
}

type FormProperty struct {
	Name         string       `yaml:"name"`
	Label        string       `yaml:"label,omitempty"`
	Type         string       `yaml:"type,omitempty"`
	Default      interface{}  `yaml:"default,omitempty"`
	Configurable bool         `yaml:"configurable,omitempty"`
	Options      []FormOption `yaml:"options,omitempty"`
	Optional     bool         `yaml:"optional,omitempty"`
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
	}
}

// GenerateEnableDisableForm generates the form to enable and disable services.
func GenerateEnableDisableForm() Form {
	enablers := []FormProperty{}
	for _, svc := range broker.GetAllServices() {
		enableForm := FormProperty{
			Name:         strings.ToLower(utils.PropertyToEnv(svc.EnabledProperty())),
			Label:        fmt.Sprintf("Let the broker create and bind %s instances", svc.CatalogEntry().Metadata.DisplayName),
			Type:         "boolean",
			Default:      true,
			Configurable: true,
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
