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
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/validation"
	"github.com/GoogleCloudPlatform/gcp-service-broker/utils"
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

// ExcludedPlansSet gets the list of ExcludedPlans as a set.
func (b *BrokerpakResource) ExcludedPlansSet() utils.StringSet {
	excludedList := utils.SplitNewlineDelimitedList(b.ExcludedPlans)
	return utils.NewStringSet(excludedList...)
}

// LoaderConfiguration holds the global configuration object for the Brokerpak
// loader.
type LoaderConfiguration struct {
	// Config holds global configuration options for the Brokerpak as a JSON object.
	Config string `json:"config" yaml:"config" validate:"required,json"`

	// Brokerpaks holds list of brokerpaks to load.
	Brokerpaks []BrokerpakResource `json:"sources" yaml:"sources" validate:"dive"`
}

// Validate returns an error if the configuration is invalid.
func (cfg *LoaderConfiguration) Validate() error {
	return validation.ValidateStruct(cfg)
}

// NewLoaderConfigurationFromEnv loads the global LoaderConfiguration from
// Viper.
func NewLoaderConfigurationFromEnv() (*LoaderConfiguration, error) {
	// load from viper
	// validate
	return nil, nil
}
