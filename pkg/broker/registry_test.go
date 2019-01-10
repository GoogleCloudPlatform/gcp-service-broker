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
	"testing"

	"github.com/pivotal-cf/brokerapi"
	"github.com/spf13/viper"
)

func TestRegistry_GetEnabledServices(t *testing.T) {
	cases := map[string]struct {
		Tag      string
		Property string
	}{
		"preview": {
			Tag:      "preview",
			Property: "compatibility.enable-preview-services",
		},
		"unmaintained": {
			Tag:      "unmaintained",
			Property: "compatibility.enable-unmaintained-services",
		},
		"eol": {
			Tag:      "eol",
			Property: "compatibility.enable-eol-services",
		},
		"beta": {
			Tag:      "beta",
			Property: "compatibility.enable-gcp-beta-services",
		},
		"deprecated": {
			Tag:      "deprecated",
			Property: "compatibility.enable-gcp-deprecated-services",
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			defer viper.Reset()

			sd := ServiceDefinition{
				Id:   "b9e4332e-b42b-4680-bda5-ea1506797474",
				Name: "test-service",
				Tags: []string{"gcp", tc.Tag},
				Plans: []ServicePlan{
					{
						ServicePlan: brokerapi.ServicePlan{
							ID:          "e1d11f65-da66-46ad-977c-6d56513baf43",
							Name:        "Builtin!",
							Description: "Standard storage class",
						},
					},
				},
				IsBuiltin: true,
			}

			registry := BrokerRegistry{}
			registry.Register(&sd)

			// shouldn't show up when property is false even if builtins are enabled
			viper.Set("compatibility.enable-builtin-services", true)
			viper.Set(tc.Property, false)
			if defns, err := registry.GetEnabledServices(); err != nil {
				t.Fatal(err)
			} else if len(defns) != 0 {
				t.Fatalf("Expected 0 definitions with %s disabled, but got %d", tc.Property, len(defns))
			}

			// should show up when property is true
			viper.Set(tc.Property, true)
			if defns, err := registry.GetEnabledServices(); err != nil {
				t.Fatal(err)
			} else if len(defns) != 1 {
				t.Fatalf("Expected 1 definition with %s enabled, but got %d", tc.Property, len(defns))
			}

			// should not show up if the service is explicitly disabled
			viper.Set("compatibility.enable-builtin-services", false)
			if defns, err := registry.GetEnabledServices(); err != nil {
				t.Fatal(err)
			} else if len(defns) != 0 {
				t.Fatalf("Expected no definition with builtins disabled, but got %d", len(defns))
			}
		})
	}
}
