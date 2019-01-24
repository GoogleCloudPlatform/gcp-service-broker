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

/*
Package toggles defines a standard way to define, list, and use feature
toggles in the service broker.

It mimics Go's `flags` package, but uses Viper as a backing store to abstract
out how a particular flag is set.

*/
package toggles

import (
	"sort"

	"github.com/GoogleCloudPlatform/gcp-service-broker/utils"
	"github.com/spf13/viper"
)

// Features is the default set of flags for enabling different features.
// For legacy compatibility reasons the flags are put under the "compatibility"
// namespace.
var Features = NewToggleSet("compatibility.")

// Toggle represents a single feature that the user can enable or disable.
type Toggle struct {
	Name           string
	Default        bool
	Description    string
	propertyPrefix string
}

// EnvironmentVariable gets the environment variable used to control the toggle.
func (toggle Toggle) EnvironmentVariable() string {
	return utils.PropertyToEnv(toggle.viperProperty())
}

// viperProperty gets the viper property of the variable to control the toggle.
func (toggle Toggle) viperProperty() string {
	return toggle.propertyPrefix + toggle.Name
}

// IsActive returns true if the toggle is enabled and false if it isn't.
func (toggle Toggle) IsActive() bool {
	return viper.GetBool(toggle.viperProperty())
}

// A ToggleSet represents a set of defined toggles. The zero value of a ToggleSet
// has no property prefix.
type ToggleSet struct {
	toggles        []Toggle
	propertyPrefix string
}

// Toggles returns a list of all registered toggles sorted lexicographically by
// their property name.
func (set *ToggleSet) Toggles() []Toggle {
	var copy []Toggle

	for _, tgl := range set.toggles {
		copy = append(copy, tgl)
	}

	sort.Slice(copy, func(i, j int) bool { return copy[i].Name < copy[j].Name })
	return copy
}

// Toggle creates a new toggle with the given name, default value, label and description.
// It also adds the toggle to an internal registry and initializes the default value in viper.
func (set *ToggleSet) Toggle(name string, value bool, description string) Toggle {
	toggle := Toggle{
		Name:           name,
		Default:        value,
		Description:    description,
		propertyPrefix: set.propertyPrefix,
	}

	set.toggles = append(set.toggles, toggle)
	viper.SetDefault(toggle.viperProperty(), value)

	return toggle
}

// NewFlagSet returns a new, empty toggle set with the specified property prefix.
// The property prefix will be prepended to any toggles exactly as-is. You MUST
// specify a trailing period if you want your properties to be namespaced.
func NewToggleSet(propertyPrefix string) *ToggleSet {
	return &ToggleSet{
		propertyPrefix: propertyPrefix,
	}
}
