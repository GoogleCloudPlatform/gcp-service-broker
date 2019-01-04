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
	"fmt"
	"testing"

	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/validation"
	"github.com/spf13/viper"
)

func TestNewBrokerpakSourceConfigFromPath(t *testing.T) {
	t.Run("is-valid-by-default", func(t *testing.T) {
		cfg := NewBrokerpakSourceConfigFromPath("/path/to/my/pak.brokerpak")
		if err := validation.ValidateStruct(cfg); err != nil {
			t.Fatalf("Expected no error got %v", err)
		}
	})

	t.Run("has-empty-config-by-default", func(t *testing.T) {
		cfg := NewBrokerpakSourceConfigFromPath("/path/to/my/pak.brokerpak")
		if cfg.Config != "{}" {
			t.Fatalf("Expected empty config '{}' got %q", cfg.Config)
		}
	})

	t.Run("has-no-excluded-plans-by-default", func(t *testing.T) {
		cfg := NewBrokerpakSourceConfigFromPath("/path/to/my/pak.brokerpak")
		if cfg.ExcludedPlans != "" {
			t.Fatalf("Expected no excluded plans, got: %v", cfg.ExcludedPlans)
		}
	})
}

func ExampleBrokerpakSourceConfig_ExcludedPlansSlice() {
	cfg := BrokerpakSourceConfig{ExcludedPlans: "FOO\nBAR"}

	fmt.Println(cfg.ExcludedPlansSlice())

	// Output: [FOO BAR]
}

func ExampleBrokerpakSourceConfig_SetExcludedPlans() {
	cfg := BrokerpakSourceConfig{}
	cfg.SetExcludedPlans([]string{"plan1", "plan2"})

	fmt.Println("slice:", cfg.ExcludedPlansSlice())
	fmt.Println("text:", cfg.ExcludedPlans)

	// Output: slice: [plan1 plan2]
	// text: plan1
	// plan2
}

func TestServiceConfig_Validate(t *testing.T) {
	cases := map[string]struct {
		Cfg ServerConfig
		Err string
	}{
		"missing config": {
			Cfg: ServerConfig{
				Config:     "",
				Brokerpaks: nil,
			},
			Err: "Key: 'ServerConfig.Config' Error:Field validation for 'Config' failed on the 'required' tag",
		},
		"bad config": {
			Cfg: ServerConfig{
				Config:     "{}aaa",
				Brokerpaks: nil,
			},
			Err: "Key: 'ServerConfig.Config' Error:Field validation for 'Config' failed on the 'json' tag",
		},
		"bad brokerpak keys": {
			Cfg: ServerConfig{
				Config: "{}",
				Brokerpaks: map[string]BrokerpakSourceConfig{
					"bad key": NewBrokerpakSourceConfigFromPath("file:///some/path"),
				},
			},
			Err: "Key: 'ServerConfig.Brokerpaks[bad key]' Error:Field validation for 'Brokerpaks[bad key]' failed on the 'osbname' tag",
		},
		"bad brokerpak values": {
			Cfg: ServerConfig{
				Config: "{}",
				Brokerpaks: map[string]BrokerpakSourceConfig{
					"good-key": BrokerpakSourceConfig{
						BrokerpakUri: "file:///some/path",
						Config:       "{}aaa",
					},
				},
			},
			Err: "Key: 'ServerConfig.Brokerpaks[good-key].Config' Error:Field validation for 'Config' failed on the 'json' tag",
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			err := tc.Cfg.Validate()
			if err == nil {
				t.Fatal("expected error, got nil")
			}

			if err.Error() != tc.Err {
				t.Fatalf("Expected %q got %q", tc.Err, err.Error())
			}
		})
	}
}

func ExampleNewServerConfigFromEnv() {
	viper.Set("brokerpak.sources", `{"good-key":{"uri":"file://path/to/brokerpak", "config":"{}"}}`)
	viper.Set("brokerpak.config", `{}`)
	defer viper.Reset() // cleanup

	cfg, err := NewServerConfigFromEnv()
	if err != nil {
		panic(err)
	}

	fmt.Println("global config:", cfg.Config)
	fmt.Println("num services:", len(cfg.Brokerpaks))

	// Output: global config: {}
	// num services: 1
}

func TestNewServerConfigFromEnv(t *testing.T) {
	cases := map[string]struct {
		Config  string
		Sources string
		Err     string
	}{
		"missing config": {
			Config:  ``,
			Sources: `{}`,
			Err:     `brokerpak config was invalid: Key: 'ServerConfig.Config' Error:Field validation for 'Config' failed on the 'required' tag`,
		},
		"bad config": {
			Config:  `{}aaa`,
			Sources: `{}`,
			Err:     `brokerpak config was invalid: Key: 'ServerConfig.Config' Error:Field validation for 'Config' failed on the 'json' tag`,
		},
		"bad brokerpak keys": {
			Config:  `{}`,
			Sources: `{"bad key":{"uri":"file://path/to/brokerpak", "config":"{}"}}`,
			Err:     `brokerpak config was invalid: Key: 'ServerConfig.Brokerpaks[bad key]' Error:Field validation for 'Brokerpaks[bad key]' failed on the 'osbname' tag`,
		},
		"bad brokerpak values": {
			Config:  `{}`,
			Sources: `{"good-key":{"uri":"file://path/to/brokerpak", "config":"aaa{}"}}`,
			Err:     `brokerpak config was invalid: Key: 'ServerConfig.Brokerpaks[good-key].Config' Error:Field validation for 'Config' failed on the 'json' tag`,
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			viper.Set("brokerpak.sources", tc.Sources)
			viper.Set("brokerpak.config", tc.Config)
			defer viper.Reset()

			cfg, err := NewServerConfigFromEnv()
			if err == nil {
				t.Fatal("expected error, got nil")
			}

			if cfg != nil {
				t.Error("Exactly one of cfg and err should be nil")
			}

			if err.Error() != tc.Err {
				t.Fatalf("Expected %q got %q", tc.Err, err.Error())
			}
		})
	}
}
