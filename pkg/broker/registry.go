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
	"sort"

	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/models"
	"github.com/GoogleCloudPlatform/gcp-service-broker/utils"
	"github.com/spf13/viper"
)

var brokerRegistry = make(map[string]*BrokerService)

// Registers a BrokerService with the service registry that various commands
// poll to create the catalog, documentation, etc.
func Register(service *BrokerService) {
	name := service.Name

	if _, ok := brokerRegistry[name]; ok {
		log.Fatalf("Tried to register multiple instances of: %q", name)
	}

	brokerRegistry[name] = service

	if err := service.init(); err != nil {
		log.Fatalf("Error registering service %q, %s", name, err)
	}
}

// GetEnabledServices returns a list of all registered brokers that the user
// has enabled the use of.
func GetEnabledServices() []*BrokerService {
	var out []*BrokerService

	for _, svc := range GetAllServices() {
		if svc.IsEnabled() {
			out = append(out, svc)
		}
	}

	return out
}

// GetAllServices returns a list of all registered brokers whether or not the
// user has enabled them. The brokers are sorted in lexocographic order based
// on name.
func GetAllServices() []*BrokerService {
	var out []*BrokerService

	for _, svc := range brokerRegistry {
		out = append(out, svc)
	}

	// Sort by name so there's a consistent order in the UI and tests.
	sort.Slice(out, func(i int, j int) bool { return out[i].Name < out[j].Name })

	return out
}

func MapServiceIdToName() map[string]string {
	out := map[string]string{}

	for _, svc := range brokerRegistry {
		out[svc.CatalogEntry().ID] = svc.Name
	}

	return out
}

type BrokerService struct {
	Name                     string
	DefaultServiceDefinition string
	ProvisionInputVariables  []BrokerVariable
	BindInputVariables       []BrokerVariable
	BindOutputVariables      []BrokerVariable
	Examples                 []ServiceExample

	// Not modifiable.
	serviceDefinition models.Service
	userDefinedPlans  []models.ServicePlan
}

func (svc *BrokerService) init() error {

	definitionProperty := svc.DefinitionProperty()

	// Set up environment variables to be compatible with legacy tile.yml configurations.
	// Bind a name of a service like google-datastore to an environment variable GOOGLE_DATASTORE
	env := utils.PropertyToEnvUnprefixed(svc.Name)
	viper.BindEnv(definitionProperty, env)

	// set defaults
	viper.SetDefault(definitionProperty, svc.DefaultServiceDefinition)
	viper.SetDefault(svc.EnabledProperty(), true)
	viper.SetDefault(svc.UserDefinedPlansProperty(), "[]")

	// Parse the service definition from the properties
	rawDefinition := []byte(viper.GetString(definitionProperty))

	var defn models.Service
	if err := json.Unmarshal(rawDefinition, &defn); err != nil {
		return err
	}
	svc.serviceDefinition = defn

	// TODO Parse any user-defined plans and include them

	return nil
}

// EnabledProperty computes the Viper property name for the boolean the user
// can set to disable or enable this service.
func (svc *BrokerService) EnabledProperty() string {
	return fmt.Sprintf("service.%s.enabled", svc.Name)
}

// DefinitionProperty computes the Viper property name for the JSON service
// definition.
func (svc *BrokerService) DefinitionProperty() string {
	return fmt.Sprintf("service.%s.definition", svc.Name)
}

// UserDefinedPlansProperty computes the Viper property name for the JSON list
// of user-defined service plans.
func (svc *BrokerService) UserDefinedPlansProperty() string {
	return fmt.Sprintf("service.%s.plans", svc.Name)
}

// IsEnabled returns false if the operator has explicitly disabled this service
// or true otherwise.
func (svc *BrokerService) IsEnabled() bool {
	return viper.GetBool(svc.EnabledProperty())
}

// CatalogEntry returns the service broker catalog entry for this service, it
// has metadata about the service so operators and programmers know which
// service and plan will work best for their purposes.
func (svc *BrokerService) CatalogEntry() models.Service {
	return svc.serviceDefinition
}

// GetPlanById finds a plan in this service by its UUID.
func (svc *BrokerService) GetPlanById(planId string) *models.ServicePlan {
	for _, plan := range svc.CatalogEntry().Plans {
		if plan.ID == planId {
			return &plan
		}
	}

	return nil
}
