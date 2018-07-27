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

	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/models"
	"github.com/spf13/viper"
)

var brokerRegistry = make(map[string]*BrokerService)

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

	for _, svc := range brokerRegistry {
		if svc.IsEnabled() {
			out = append(out, svc)
		}
	}

	return out
}

// GetAllServices returns a list of all registered brokers whether or not the
// user has enabled them.
func GetAllServices() []*BrokerService {
	var out []*BrokerService

	for _, svc := range brokerRegistry {
		out = append(out, svc)
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

	// Not modifiable
	serviceDefinition models.Service
	userDefinedPlans  []models.ServicePlan

	enabledProperty          string
	userDefinedPlansProperty string
	definitionProperty       string
}

func (svc *BrokerService) init() error {
	// create properties
	svc.definitionProperty = fmt.Sprintf("service.%s.definition", svc.Name)
	svc.enabledProperty = fmt.Sprintf("service.%s.enabled", svc.Name)
	svc.userDefinedPlansProperty = fmt.Sprintf("service.%s.plans", svc.Name)

	// Set up environment variables to be compatible with legacy tile.yml configurations.
	// Bind a name of a service like google-datastore to an environment variable GOOGLE_DATASTORE
	replacer := strings.NewReplacer("-", "_")
	env := replacer.Replace(strings.ToUpper(svc.Name))
	viper.BindEnv(svc.definitionProperty, env)

	// set defaults
	viper.SetDefault(svc.definitionProperty, svc.DefaultServiceDefinition)
	viper.SetDefault(svc.enabledProperty, true)
	viper.SetDefault(svc.userDefinedPlansProperty, "[]")

	// Parse the service definition from the properties
	rawDefinition := []byte(viper.GetString(svc.definitionProperty))

	var defn models.Service
	if err := json.Unmarshal(rawDefinition, &defn); err != nil {
		return err
	}
	svc.serviceDefinition = defn

	// TODO Parse any user-defined plans and include them

	return nil
}

func (svc *BrokerService) IsEnabled() bool {
	return viper.GetBool(svc.enabledProperty)
}

func (svc *BrokerService) CatalogEntry() models.Service {
	return svc.serviceDefinition
}

func (svc *BrokerService) GetPlanById(planId string) *models.ServicePlan {
	for _, plan := range svc.CatalogEntry().Plans {
		if plan.ID == planId {
			return &plan
		}
	}

	return nil
}
