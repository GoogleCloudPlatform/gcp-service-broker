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
	"testing"

	"github.com/pivotal-cf/brokerapi"
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

func TestExtractRole(t *testing.T) {
	cases := map[string]struct {
		Details      string
		ExpectedRole string
		ShouldError  bool
	}{
		"Positive":     {`{"role":"foo"}`, "foo", false},
		"Bad Json":     {`}}`, "", true},
		"Missing Role": {`{"ROLE":"bar"}`, "", true},
		"bad type":     {`{"role":1234}`, "", true},
	}

	for name, testcase := range cases {
		details := brokerapi.BindDetails{RawParameters: []byte(testcase.Details)}
		role, err := extractRole(details)

		if role != testcase.ExpectedRole {
			t.Errorf("%s) expected role: %q actual: %q", name, testcase.ExpectedRole, role)
		}

		wasError := err != nil
		if wasError != testcase.ShouldError {
			t.Errorf("%s) expected error? %v got: %s", name, testcase.ShouldError, err)
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
