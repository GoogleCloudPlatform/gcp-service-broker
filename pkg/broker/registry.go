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
	"log"
	"sort"

	"github.com/GoogleCloudPlatform/gcp-service-broker/utils"
	"github.com/spf13/viper"
)

var brokerRegistry = make(map[string]*ServiceDefinition)

// Registers a ServiceDefinition with the service registry that various commands
// poll to create the catalog, documentation, etc.
func Register(service *ServiceDefinition) {
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
func GetEnabledServices() []*ServiceDefinition {
	var out []*ServiceDefinition

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
func GetAllServices() []*ServiceDefinition {
	var out []*ServiceDefinition

	for _, svc := range brokerRegistry {
		out = append(out, svc)
	}

	// Sort by name so there's a consistent order in the UI and tests.
	sort.Slice(out, func(i int, j int) bool { return out[i].Name < out[j].Name })

	return out
}

// GetServiceById returns the service with the given ID, if it does not exist
// or one of the services has a parse error then an error is returned.
func GetServiceById(id string) (*ServiceDefinition, error) {
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
