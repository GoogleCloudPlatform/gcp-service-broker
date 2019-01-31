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

	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/toggles"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/validation"
	"github.com/GoogleCloudPlatform/gcp-service-broker/utils"
)

var (
	// The following flags enable and disable services based on their tags.
	// The guiding philosophy for defaults is optimistic about new technology and pessimistic about old.
	lifecycleTagToggles = map[string]toggles.Toggle{
		"preview":      toggles.Features.Toggle("enable-preview-services", true, `Enable services that are new to the broker this release.`),
		"unmaintained": toggles.Features.Toggle("enable-unmaintained-services", false, `Enable broker services that are unmaintained.`),
		"eol":          toggles.Features.Toggle("enable-eol-services", false, `Enable broker services that are end of life.`),
		"beta":         toggles.Features.Toggle("enable-gcp-beta-services", true, "Enable services that are in GCP Beta. These have no SLA or support policy."),
		"deprecated":   toggles.Features.Toggle("enable-gcp-deprecated-services", false, "Enable services that use deprecated GCP components."),
		"terraform":    toggles.Features.Toggle("enable-terraform-services", false, "Enable services that use the experimental, unstable, Terraform back-end."),
	}

	enableBuiltinServices = toggles.Features.Toggle("enable-builtin-services", true, `Enable services that are built in to the broker i.e. not brokerpaks.`)
)

// ServiceRegistry holds the list of ServiceDefinitions that can be provisioned
// by the GCP Service Broker.
type ServiceRegistry struct {
	services map[string]*ServiceDefinition
	config   ServiceConfigMap
}

// Register a ServiceDefinition with the service registry that various commands
// poll to create the catalog, documentation, etc.
func (registry *ServiceRegistry) Register(service *ServiceDefinition) {
	name := service.Name

	if _, ok := registry.services[name]; ok {
		log.Fatalf("Tried to register multiple instances of: %q", name)
	}

	cfg := registry.config[service.Id]
	if err := service.SetConfig(cfg); err != nil {
		log.Fatalf("Error registering service %q, %s", name, err)
	}

	if err := validation.ValidateStruct(service); err != nil {
		log.Fatalf("Error validating service %q, %s", name, err)
	}

	registry.services[name] = service
}

// GetEnabledServices returns a list of all registered brokers that the user
// has enabled the use of.
func (registry *ServiceRegistry) GetEnabledServices() []*ServiceDefinition {
	var out []*ServiceDefinition

	for _, svc := range registry.GetAllServices() {
		isEnabled := true

		if svc.IsBuiltin {
			isEnabled = enableBuiltinServices.IsActive()
		}

		tags := utils.NewStringSet(svc.Tags...)
		for tag, toggle := range lifecycleTagToggles {
			if !toggle.IsActive() && tags.Contains(tag) {
				isEnabled = false
				break
			}
		}

		if isEnabled {
			out = append(out, svc)
		}
	}

	return out
}

// GetAllServices returns a list of all registered brokers whether or not the
// user has enabled them. The brokers are sorted in lexocographic order based
// on name.
func (registry *ServiceRegistry) GetAllServices() []*ServiceDefinition {
	var out []*ServiceDefinition

	for _, svc := range registry.services {
		out = append(out, svc)
	}

	// Sort by name so there's a consistent order in the UI and tests.
	sort.Slice(out, func(i int, j int) bool { return out[i].Name < out[j].Name })

	return out
}

// GetServiceById returns the service with the given ID, if it does not exist
// or one of the services has a parse error then an error is returned.
func (registry *ServiceRegistry) GetServiceById(id string) (*ServiceDefinition, error) {
	for _, svc := range registry.services {
		if svc.Id == id {
			return svc, nil
		}
	}

	return nil, fmt.Errorf("Unknown service ID: %q", id)
}

// ServiceCount returns the total number of registered services.
func (registry *ServiceRegistry) ServiceCount() int {
	return len(registry.services)
}

// NewServiceRegistry creates a blank registry with the given config map.
func NewServiceRegistry(cfg ServiceConfigMap) *ServiceRegistry {
	return &ServiceRegistry{
		config:   cfg,
		services: make(map[string]*ServiceDefinition),
	}
}
