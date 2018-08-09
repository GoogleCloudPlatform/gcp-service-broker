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
	"strings"

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

	// Set up environment variables to be compatible with legacy tile.yml configurations.
	// Bind a name of a service like google-datastore to an environment variable GOOGLE_DATASTORE
	env := utils.PropertyToEnvUnprefixed(service.Name)
	viper.BindEnv(service.DefinitionProperty(), env)

	// set defaults
	viper.SetDefault(service.EnabledProperty(), true)

	// Test deserializing the user defined plans and service definition
	if _, err := service.CatalogEntry(); err != nil {
		log.Fatalf("Error registering service %q, %s", name, err)
	}

	brokerRegistry[name] = service
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

// GetServiceById returns the service with the given ID, if it does not exist
// or one of the services has a parse error then an error is returned.
func GetServiceById(id string) (*BrokerService, error) {
	for _, svc := range brokerRegistry {
		if entry, err := svc.CatalogEntry(); err != nil {
			return nil, err
		} else {
			if entry.ID == id {
				return svc, nil
			}
		}
	}

	return nil, fmt.Errorf("Unknown service ID: %q", id)
}

type BrokerService struct {
	Name                     string
	DefaultServiceDefinition string
	ProvisionInputVariables  []BrokerVariable
	BindInputVariables       []BrokerVariable
	BindOutputVariables      []BrokerVariable
	PlanVariables            []BrokerVariable
	Examples                 []ServiceExample
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

// TileUserDefinedPlansVariable returns the name of the user defined plans
// variable for the broker tile.
func (svc *BrokerService) TileUserDefinedPlansVariable() string {
	prefix := "GOOGLE_"

	v := utils.PropertyToEnvUnprefixed(svc.Name)
	if strings.HasPrefix(v, prefix) {
		v = v[len(prefix):]
	}

	return v + "_CUSTOM_PLANS"
}

// IsEnabled returns false if the operator has explicitly disabled this service
// or true otherwise.
func (svc *BrokerService) IsEnabled() bool {
	return viper.GetBool(svc.EnabledProperty())
}

// CatalogEntry returns the service broker catalog entry for this service, it
// has metadata about the service so operators and programmers know which
// service and plan will work best for their purposes.
func (svc *BrokerService) CatalogEntry() (*models.Service, error) {
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
func (svc *BrokerService) GetPlanById(planId string) (*models.ServicePlan, error) {
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
// the plans were not valid JSON.
func (svc *BrokerService) UserDefinedPlans() ([]models.ServicePlan, error) {
	plans := []models.ServicePlan{}

	userPlanJson := viper.GetString(svc.UserDefinedPlansProperty())
	if userPlanJson == "" {
		return plans, nil
	}

	err := json.Unmarshal([]byte(userPlanJson), &plans)
	return plans, err
}

// ServiceDefinition extracts service definition from the environment, failing
// if the definition was not valid JSON.
func (svc *BrokerService) ServiceDefinition() (*models.Service, error) {
	jsonDefinition := viper.GetString(svc.DefinitionProperty())
	if jsonDefinition == "" {
		jsonDefinition = svc.DefaultServiceDefinition
	}

	var defn models.Service
	err := json.Unmarshal([]byte(jsonDefinition), &defn)
	if err != nil {
		return nil, fmt.Errorf("Error parsing service definition for %q: %s", svc.Name, err)
	}
	return &defn, err
}
