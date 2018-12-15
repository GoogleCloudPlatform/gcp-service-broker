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

package brokerpak

import (
	"encoding/json"
	"fmt"

	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/broker"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/providers/tf"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/validation"
	"github.com/GoogleCloudPlatform/gcp-service-broker/utils"
	"github.com/spf13/viper"
)

// BrokerpakResource represents a single configuration of a brokerpak.
type BrokerpakResource struct {
	// BrokerpakUri holds the URI for loading the Brokerpak.
	BrokerpakUri string `json:"uri" validate:"required,uri"`
	// ServicePrefix holds an optional prefix that will be prepended to every service name.
	ServicePrefix string `json:"service_prefix" validate:"osbname"`
	// ExcludedPlans holds a newline delimited list of service plan UUIDs that will be excluded at registration time.
	ExcludedPlans string `json:"excluded_plans"`
	// Config holds the configuration options for the Brokerpak as a JSON object.
	Config string `json:"config" validate:"required,json"`
	// Notes holds user-defined notes about the Brokerpak and shouldn't be used programatically.
	Notes string `json:"notes"`
}

// Validate returns an error if the resource is invalid.
func (b *BrokerpakResource) Validate() error {
	return validation.ValidateStruct(b)
}

// ExcludedPlansSet gets the list of ExcludedPlans as a slice.
func (b *BrokerpakResource) ExcludedPlansSlice() []string {
	return utils.SplitNewlineDelimitedList(b.ExcludedPlans)
}

// RegisterPak fetches the brokerpak and registers it with the given registry.
func (b *BrokerpakResource) RegisterPak(pack string, registry broker.BrokerRegistry) error {
	brokerPak, err := OpenBrokerPak(pack)
	if err != nil {
		return fmt.Errorf("couldn't open brokerpak: %q: %v", pack, err)
	}
	defer brokerPak.Close()

	toIgnore := utils.NewStringSet(b.ExcludedPlansSlice()...)
	brokerPak.ServiceTransformer = func(t tf.TfServiceDefinitionV1) *tf.TfServiceDefinitionV1 {
		if toIgnore.Contains(t.Id) {
			return nil
		}

		t.Name = b.ServicePrefix + t.Name

		return &t
	}

	if err := brokerPak.Register(registry); err != nil {
		return fmt.Errorf("couldn't register brokerpak: %q: %v", pack, err)
	}

	return nil
}

// NewBrokerpakResourceFromPath creates a new BrokerpakResource from a path.
func NewBrokerpakResourceFromPath(path string) *BrokerpakResource {
	return &BrokerpakResource{
		BrokerpakUri: path,
		Config:       "{}",
	}
}

// LoaderConfiguration holds the global configuration object for the Brokerpak
// loader.
type LoaderConfiguration struct {
	// Config holds global configuration options for the Brokerpak as a JSON object.
	Config string `validate:"required,json"`

	// Brokerpaks holds list of brokerpaks to load.
	Brokerpaks map[string]BrokerpakResource `validate:"dive"`
}

// Validate returns an error if the configuration is invalid.
func (cfg *LoaderConfiguration) Validate() error {
	return validation.ValidateStruct(cfg)
}

// NewLoaderConfigurationFromEnv loads the global configuration from Viper.
func NewLoaderConfigurationFromEnv() (*LoaderConfiguration, error) {
	paks := map[string]BrokerpakResource{}
	sources := viper.GetString("brokerpak.sources")
	if err := json.Unmarshal([]byte(sources), &paks); err != nil {
		return nil, fmt.Errorf("couldn't deserialize brokerpak source config: %v", err)
	}

	cfg := LoaderConfiguration{
		Config:     viper.GetString("brokerpak.config"),
		Brokerpaks: paks,
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("brokerpak config was invalid: %v", err)
	}

	return &cfg, nil
}
