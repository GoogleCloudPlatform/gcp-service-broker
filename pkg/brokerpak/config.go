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
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/toggles"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/validation"
	"github.com/GoogleCloudPlatform/gcp-service-broker/utils"
	"github.com/spf13/viper"
)

const (
	// BuiltinPakLocation is the file-system location to load brokerpaks from to
	// make them look builtin.
	BuiltinPakLocation      = "/usr/share/gcp-service-broker/builtin-brokerpaks"
	brokerpakSourcesKey     = "brokerpak.sources"
	brokerpakConfigKey      = "brokerpak.config"
	brokerpakBuiltinPathKey = "brokerpak.builtin.path"
)

var loadBuiltinToggle = toggles.Features.Toggle("enable-builtin-brokerpaks", true, `Load brokerpaks that are built-in to the software.`)

func init() {
	viper.SetDefault(brokerpakSourcesKey, "{}")
	viper.SetDefault(brokerpakConfigKey, "{}")
	viper.SetDefault(brokerpakBuiltinPathKey, BuiltinPakLocation)
}

// BrokerpakSourceConfig represents a single configuration of a brokerpak.
type BrokerpakSourceConfig struct {
	// BrokerpakUri holds the URI for loading the Brokerpak.
	BrokerpakUri string `json:"uri"`
	// ServicePrefix holds an optional prefix that will be prepended to every service name.
	ServicePrefix string `json:"service_prefix"`
	// ExcludedServices holds a newline delimited list of service UUIDs that will be excluded at registration time.
	ExcludedServices string `json:"excluded_services"`
	// Config holds the configuration options for the Brokerpak as a JSON object.
	Config string `json:"config"`
	// Notes holds user-defined notes about the Brokerpak and shouldn't be used programatically.
	Notes string `json:"notes"`
}

var _ validation.Validatable = (*BrokerpakSourceConfig)(nil)

// Validate implements validation.Validatable.
func (b *BrokerpakSourceConfig) Validate() (errs *validation.FieldError) {

	errs = errs.Also(validation.ErrIfBlank(b.BrokerpakUri, "uri"))

	if b.ServicePrefix != "" {
		errs = errs.Also(validation.ErrIfNotOSBName(b.ServicePrefix, "service_prefix"))
	}

	errs = errs.Also(validation.ErrIfNotJSON(json.RawMessage(b.Config), "config"))

	return errs
}

// ExcludedServicesSlice gets the ExcludedServices as a slice of UUIDs.
func (b *BrokerpakSourceConfig) ExcludedServicesSlice() []string {
	return utils.SplitNewlineDelimitedList(b.ExcludedServices)
}

// SetExcludedServices sets the ExcludedServices from a slice of UUIDs.
func (b *BrokerpakSourceConfig) SetExcludedServices(services []string) {
	b.ExcludedServices = strings.Join(services, "\n")
}

// NewBrokerpakSourceConfigFromPath creates a new BrokerpakSourceConfig from a path.
func NewBrokerpakSourceConfigFromPath(path string) BrokerpakSourceConfig {
	return BrokerpakSourceConfig{
		BrokerpakUri: path,
		Config:       "{}",
	}
}

// ServerConfig holds the Brokerpak configuration for the server.
type ServerConfig struct {
	// Config holds global configuration options for the Brokerpak as a JSON object.
	Config string

	// Brokerpaks holds list of brokerpaks to load.
	Brokerpaks map[string]BrokerpakSourceConfig
}

var _ validation.Validatable = (*ServerConfig)(nil)

// Validate returns an error if the configuration is invalid.
func (cfg *ServerConfig) Validate() (errs *validation.FieldError) {
	errs = errs.Also(validation.ErrIfNotJSON(json.RawMessage(cfg.Config), "Config"))

	for k, v := range cfg.Brokerpaks {
		errs = errs.Also(validation.ErrIfNotOSBName(k, "").ViaFieldKey("Brokerpaks", k))
		errs = errs.Also(v.Validate().ViaFieldKey("Brokerpaks", k))
	}

	return errs
}

// NewServerConfigFromEnv loads the global Brokerpak config from Viper.
func NewServerConfigFromEnv() (*ServerConfig, error) {
	paks := map[string]BrokerpakSourceConfig{}
	sources := viper.GetString(brokerpakSourcesKey)
	if err := json.Unmarshal([]byte(sources), &paks); err != nil {
		return nil, fmt.Errorf("couldn't deserialize brokerpak source config: %v", err)
	}

	cfg := ServerConfig{
		Config:     viper.GetString(brokerpakConfigKey),
		Brokerpaks: paks,
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("brokerpak config was invalid: %v", err)
	}

	// Builtin paks fail validation because they reference the local filesystem
	// but do work.
	if loadBuiltinToggle.IsActive() {
		log.Println("loading builtin brokerpaks")
		paks, err := ListBrokerpaks(viper.GetString(brokerpakBuiltinPathKey))
		if err != nil {
			return nil, fmt.Errorf("couldn't load builtin brokerpaks: %v", err)
		}

		for i, path := range paks {
			key := fmt.Sprintf("builtin-%d", i)
			config := NewBrokerpakSourceConfigFromPath(path)
			config.Notes = fmt.Sprintf("This pak was automatically loaded because the toggle %s was enabled", loadBuiltinToggle.EnvironmentVariable())
			cfg.Brokerpaks[key] = config
		}
	}

	return &cfg, nil
}

// ListBrokerpaks gets all brokerpaks in a given directory.
func ListBrokerpaks(directory string) ([]string, error) {
	var paks []string
	err := filepath.Walk(filepath.FromSlash(directory), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if filepath.Ext(path) == ".brokerpak" {
			paks = append(paks, path)
		}

		return nil
	})

	sort.Strings(paks)

	if os.IsNotExist(err) {
		return paks, nil
	}

	return paks, err
}

func newLocalFileServerConfig(path string) *ServerConfig {
	return &ServerConfig{
		Config: viper.GetString(brokerpakConfigKey),
		Brokerpaks: map[string]BrokerpakSourceConfig{
			"local-brokerpak": NewBrokerpakSourceConfigFromPath(path),
		},
	}
}
