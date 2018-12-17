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

// // BrokerpakSourceConfig represents a single configuration of a brokerpak.
// type BrokerpakSourceConfig struct {
// 	// BrokerpakUri holds the URI for loading the Brokerpak.
// 	BrokerpakUri string `json:"uri" validate:"required,uri"`
// 	// ServicePrefix holds an optional prefix that will be prepended to every service name.
// 	ServicePrefix string `json:"service_prefix" validate:"omitempty,osbname"`
// 	// ExcludedPlans holds a newline delimited list of service plan UUIDs that will be excluded at registration time.
// 	ExcludedPlans string `json:"excluded_plans"`
// 	// Config holds the configuration options for the Brokerpak as a JSON object.
// 	Config string `json:"config" validate:"required,json"`
// 	// Notes holds user-defined notes about the Brokerpak and shouldn't be used programatically.
// 	Notes string `json:"notes"`
// }
//
// // ExcludedPlansSlice gets the ExcludedPlans as a slice of UUIDs.
// func (b *BrokerpakSourceConfig) ExcludedPlansSlice() []string {
// 	return utils.SplitNewlineDelimitedList(b.ExcludedPlans)
// }
//
// // SetExcludedPlans sets the ExcludedPlans from a slice of UUIDs.
// func (b *BrokerpakSourceConfig) SetExcludedPlans(plans []string) {
// 	b.ExcludedPlans = strings.Join(plans, "\n")
// }

// // NewBrokerpakSourceConfigFromPath creates a new BrokerpakSourceConfig from a path.
// func NewBrokerpakSourceConfigFromPath(path string) BrokerpakSourceConfig {
// 	return BrokerpakSourceConfig{
// 		BrokerpakUri: path,
// 		Config:       "{}",
// 	}
// }

// // ServerConfig holds the Brokerpak configuration for the server.
// type ServerConfig struct {
// 	// Config holds global configuration options for the Brokerpak as a JSON object.
// 	Config string `validate:"required,json"`
//
// 	// Brokerpaks holds list of brokerpaks to load.
// 	Brokerpaks map[string]BrokerpakSourceConfig `validate:"dive,keys,osbname,endkeys"`
// }

// // Validate returns an error if the configuration is invalid.
// func (cfg *ServerConfig) Validate() error {
// 	return validation.ValidateStruct(cfg)
// }

// // NewServerConfigFromEnv loads the global Brokerpak config from Viper.
// func NewServerConfigFromEnv() (*ServerConfig, error) {
// 	paks := map[string]BrokerpakSourceConfig{}
// 	sources := viper.GetString("brokerpak.sources")
// 	if err := json.Unmarshal([]byte(sources), &paks); err != nil {
// 		return nil, fmt.Errorf("couldn't deserialize brokerpak source config: %v", err)
// 	}
//
// 	cfg := ServerConfig{
// 		Config:     viper.GetString("brokerpak.config"),
// 		Brokerpaks: paks,
// 	}
//
// 	if err := cfg.Validate(); err != nil {
// 		return nil, fmt.Errorf("brokerpak config was invalid: %v", err)
// 	}
//
// 	return &cfg, nil
// }
