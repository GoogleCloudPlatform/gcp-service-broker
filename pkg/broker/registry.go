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

package broker

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	osbapi "github.com/pivotal-cf/brokerapi"
	"github.com/spf13/viper"
)

var brokerRegistry = make(map[string]*BrokerService)

func Register(service *BrokerService) {
	brokerRegistry[service.Name] = service

	err := service.init()
	if err != nil {
		log.Fatalf("Error registering service %q, %s", service.Name, err)
	}
}

func GetServiceById(id string) {

}

func GetEnabledServices() []*BrokerService {
	var out []*BrokerService

	for _, svc := range brokerRegistry {
		if svc.IsEnabled() {
			out = append(out, svc)
		}
	}

	return out
}

func GetAllServices() []*BrokerService {
	var out []*BrokerService

	for _, svc := range brokerRegistry {
		out = append(out, svc)
	}

	return out
}

// Get service plans

type BrokerService struct {
	Name                     string
	DefaultServiceDefinition string
	ProvisionInputVariables  []BrokerVariable
	BindInputVariables       []BrokerVariable
	BindOutputVariables      []BrokerVariable
	Examples                 []ServiceExample

	// Not modifiable
	serviceDefinition osbapi.Service
	userDefinedPlans  []osbapi.ServicePlan

	enabledProperty            string
	userDefinedPlansProperty   string
	enableDefaultPlansProperty string
	definitionProperty         string
}

// CLOUDSQL_MYSQL_CUSTOM_PLANS -> GSB_SERVICE_CLOUDSQL_MYSQL(|_PLANS|_ENABLED)

func (svc *BrokerService) init() error {
	// create properties
	svc.definitionProperty = fmt.Sprintf("service.%s.definition", svc.Name)
	svc.enabledProperty = fmt.Sprintf("service.%s.enabled", svc.Name)
	svc.userDefinedPlansProperty = fmt.Sprintf("service.%s.plans", svc.Name)
	svc.enableDefaultPlansProperty = fmt.Sprintf("service.%s.enable_default_plans", svc.Name)

	// Set up environment variables to be compatible with legacy tile.yml configurations.
	// Bind a name of a service like google-datastore to an environment variable GOOGLE_DATASTORE
	replacer := strings.NewReplacer("-", "_")
	env := replacer.Replace(strings.ToUpper(svc.Name))
	viper.BindEnv(svc.definitionProperty, env)

	// set defaults
	viper.SetDefault(svc.definitionProperty, svc.DefaultServiceDefinition)
	viper.SetDefault(svc.enabledProperty, true)
	viper.SetDefault(svc.enableDefaultPlansProperty, true)
	viper.SetDefault(svc.userDefinedPlansProperty, "[]")

	// Parse the service definition from the properties
	rawDefinition := []byte(viper.GetString(svc.definitionProperty))

	var defn osbapi.Service
	if err := json.Unmarshal(rawDefinition, &defn); err != nil {
		return fmt.Errorf("Error getting catalog info for service %q: %v", svc.Name, err)
	}
	svc.serviceDefinition = defn

	// TODO Parse the user-defined plans

	return nil
}

func (svc *BrokerService) IsEnabled() bool {
	return viper.GetBool(svc.enabledProperty)
}

func (svc *BrokerService) AreDefaultPlansEnabled() bool {
	return viper.GetBool(svc.enableDefaultPlansProperty)
}

func (svc *BrokerService) CatalogEntry() osbapi.Service {
	metadata := svc.serviceDefinition
	// TODO User defined plans and schemas
	// If default plans are not enabled, remove them from the service
	// If user defined plans are present, add them to the service
	// Generate the schemas based on provision/bind input variables

	return metadata
}

func (svc *BrokerService) DefaultPlans() interface{} {
	return nil
}

// DefaultPlans() []osbapi.ServicePlan
// HelpUrl() string

// service.{name}.plans
// service.{name}.service

// ProvisionExample() string
// BindExample() string

// - name: service
// 	label: Service
// 	type: dropdown_select
// 	description: "The service this plan is associated with"
// 	options:
// 		- name: '51b3e27e-d323-49ce-8c5f-1211e6409e82'
// 			label: 'Spanner'
