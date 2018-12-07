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

package account_managers

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/broker"
	"github.com/spf13/viper"
)

func TestWhitelistAllows(t *testing.T) {
	cases := map[string]struct {
		Whitelist []string
		Role      string
		Expected  bool
	}{
		"Empty Whitelist": {
			Whitelist: []string{},
			Role:      "test",
			Expected:  false,
		},
		"Contained": {
			Whitelist: []string{"foo", "bar", "bazz"},
			Role:      "bar",
			Expected:  true,
		},
		"Not Contained": {
			Whitelist: []string{"foo", "bar", "bazz"},
			Role:      "bazzz",
			Expected:  false,
		},
	}

	for name, testcase := range cases {
		actual := whitelistAllows(testcase.Whitelist, testcase.Role)
		if actual != testcase.Expected {
			t.Errorf("%s) test failed expected? %v actual: %v, test: %#v", name, testcase.Expected, actual, testcase)
		}
	}
}

func ExampleRoleWhitelistProperty() {
	serviceName := "left-handed-smoke-sifter"

	fmt.Println(RoleWhitelistProperty(serviceName))

	// Output: service.left-handed-smoke-sifter.whitelist
}

func ExampleroleWhitelist() {
	serviceName := "my-service"
	defaultRoleWhitelist := []string{"a", "b", "c"}

	viper.Set(RoleWhitelistProperty(serviceName), "")
	fmt.Println(roleWhitelist(serviceName, defaultRoleWhitelist))

	viper.Set(RoleWhitelistProperty(serviceName), "x,y,z")
	fmt.Println(roleWhitelist(serviceName, defaultRoleWhitelist))

	// Output: [a b c]
	// [x y z]
}

func TestServiceAccountBindInputVariables(t *testing.T) {

	cases := map[string]struct {
		Whitelist   []string
		Override    string
		DefaultRole string
		Expected    broker.BrokerVariable
	}{
		"default in whitelist": {
			Whitelist:   []string{"foo"},
			DefaultRole: "foo",
			Expected: broker.BrokerVariable{
				FieldName: "role",
				Type:      broker.JsonTypeString,
				Details:   overridableBindMessage,

				Required: false,
				Default:  "foo",
				Enum:     map[interface{}]string{"foo": "roles/foo"},
			},
		},

		"default not in whitelist": {
			Whitelist:   []string{"foo"},
			DefaultRole: "test",
			Expected: broker.BrokerVariable{
				FieldName: "role",
				Type:      broker.JsonTypeString,
				Details:   overridableBindMessage,

				Required: true,
				Default:  nil,
				Enum:     map[interface{}]string{"foo": "roles/foo"},
			},
		},

		"default not in override whitelist": {
			Whitelist:   []string{"foo"},
			Override:    "bar,bazz",
			DefaultRole: "foo",
			Expected: broker.BrokerVariable{
				FieldName: "role",
				Type:      broker.JsonTypeString,
				Details:   overridableBindMessage,

				Required: true,
				Default:  nil,
				Enum:     map[interface{}]string{"bar": "roles/bar", "bazz": "roles/bazz"},
			},
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			viper.Set(RoleWhitelistProperty("my-service"), tc.Override)
			vars := ServiceAccountBindInputVariables("my-service", tc.Whitelist, tc.DefaultRole)
			if len(vars) != 1 {
				t.Fatalf("Expected 1 input variable, got %d", len(vars))
			}

			if !reflect.DeepEqual(vars[0], tc.Expected) {
				t.Fatalf("Expected %#v, got %#v", tc.Expected, vars[0])

			}
		})

	}
}
