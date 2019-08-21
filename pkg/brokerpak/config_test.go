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
	"reflect"
	"testing"

	"github.com/spf13/viper"
)

func TestNewBrokerpakSourceConfigFromPath(t *testing.T) {
	t.Run("is-valid-by-default", func(t *testing.T) {
		cfg := NewBrokerpakSourceConfigFromPath("/path/to/my/pak.brokerpak")

		if err := cfg.Validate(); err != nil {
			t.Fatalf("Expected no error got %v", err)
		}
	})

	t.Run("has-empty-config-by-default", func(t *testing.T) {
		cfg := NewBrokerpakSourceConfigFromPath("/path/to/my/pak.brokerpak")
		if cfg.Config != "{}" {
			t.Fatalf("Expected empty config '{}' got %q", cfg.Config)
		}
	})

	t.Run("has-no-excluded-services-by-default", func(t *testing.T) {
		cfg := NewBrokerpakSourceConfigFromPath("/path/to/my/pak.brokerpak")
		if cfg.ExcludedServices != "" {
			t.Fatalf("Expected no excluded services, got: %v", cfg.ExcludedServices)
		}
	})
}

func ExampleBrokerpakSourceConfig_ExcludedServicesSlice() {
	cfg := BrokerpakSourceConfig{ExcludedServices: "FOO\nBAR"}

	fmt.Println(cfg.ExcludedServicesSlice())

	// Output: [FOO BAR]
}

func ExampleBrokerpakSourceConfig_SetExcludedServices() {
	cfg := BrokerpakSourceConfig{}
	cfg.SetExcludedServices([]string{"plan1", "plan2"})

	fmt.Println("slice:", cfg.ExcludedServicesSlice())
	fmt.Println("text:", cfg.ExcludedServices)

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
			Err: "invalid JSON: Config",
		},
		"bad config": {
			Cfg: ServerConfig{
				Config:     "{}aaa",
				Brokerpaks: nil,
			},
			Err: "invalid JSON: Config",
		},
		"bad brokerpak keys": {
			Cfg: ServerConfig{
				Config: "{}",
				Brokerpaks: map[string]BrokerpakSourceConfig{
					"bad key": NewBrokerpakSourceConfigFromPath("file:///some/path"),
				},
			},
			Err: "field must match '^[a-zA-Z0-9-\\.]+$': Brokerpaks[bad key]",
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
			Err: "invalid JSON: Brokerpaks[good-key].config",
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

func ExampleNewServerConfigFromEnv_customBuiltin() {
	viper.Set("brokerpak.sources", `{}`)
	viper.Set("brokerpak.config", `{}`)
	viper.Set(brokerpakBuiltinPathKey, "testdata/dummy-brokerpaks")
	viper.Set("compatibility.enable-builtin-brokerpaks", "true")
	defer viper.Reset() // cleanup

	cfg, err := NewServerConfigFromEnv()
	if err != nil {
		panic(err)
	}

	fmt.Println("num services:", len(cfg.Brokerpaks))

	// Output: num services: 2
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
			Err:     `brokerpak config was invalid: invalid JSON: Config`,
		},
		"bad config": {
			Config:  `{}aaa`,
			Sources: `{}`,
			Err:     `brokerpak config was invalid: invalid JSON: Config`,
		},
		"bad brokerpak keys": {
			Config:  `{}`,
			Sources: `{"bad key":{"uri":"file://path/to/brokerpak", "config":"{}"}}`,
			Err:     `brokerpak config was invalid: field must match '^[a-zA-Z0-9-\.]+$': Brokerpaks[bad key]`,
		},
		"bad brokerpak values": {
			Config:  `{}`,
			Sources: `{"good-key":{"uri":"file://path/to/brokerpak", "config":"aaa{}"}}`,
			Err:     `brokerpak config was invalid: invalid JSON: Brokerpaks[good-key].config`,
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

func TestListBrokerpaks(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		path         string
		expectedErr  error
		expectedPaks []string
	}{
		"directory does not exist": {
			path:         "testdata/dne",
			expectedErr:  nil,
			expectedPaks: nil,
		},
		"directory contains no brokerpaks": {
			path:         "testdata/no-brokerpaks",
			expectedErr:  nil,
			expectedPaks: nil,
		},
		"directory contains brokerpaks": {
			path:        "testdata/dummy-brokerpaks",
			expectedErr: nil,
			expectedPaks: []string{
				"testdata/dummy-brokerpaks/first.brokerpak",
				"testdata/dummy-brokerpaks/second.brokerpak",
			},
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			paks, err := ListBrokerpaks(tc.path)

			if err != nil || tc.expectedErr != nil {
				if fmt.Sprint(err) != fmt.Sprint(tc.expectedErr) {
					t.Fatalf("expected err: %v got: %v", err, tc.expectedErr)
				}

				return
			}

			if !reflect.DeepEqual(tc.expectedPaks, paks) {
				t.Fatalf("expected paks: %v got: %v", tc.expectedPaks, paks)
			}
		})
	}
}
