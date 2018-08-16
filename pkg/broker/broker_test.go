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
	"fmt"
	"testing"

	"github.com/spf13/viper"
)

func ExampleBrokerService_EnabledProperty() {
	service := BrokerService{
		Name: "left-handed-smoke-sifter",
	}

	fmt.Println(service.EnabledProperty())

	// Output: service.left-handed-smoke-sifter.enabled
}

func ExampleBrokerService_DefinitionProperty() {
	service := BrokerService{
		Name: "left-handed-smoke-sifter",
	}

	fmt.Println(service.DefinitionProperty())

	// Output: service.left-handed-smoke-sifter.definition
}

func ExampleBrokerService_UserDefinedPlansProperty() {
	service := BrokerService{
		Name: "left-handed-smoke-sifter",
	}

	fmt.Println(service.UserDefinedPlansProperty())

	// Output: service.left-handed-smoke-sifter.plans
}

func ExampleBrokerService_RoleWhitelistProperty() {
	service := BrokerService{
		Name: "left-handed-smoke-sifter",
	}

	fmt.Println(service.RoleWhitelistProperty())

	// Output: service.left-handed-smoke-sifter.whitelist
}

func ExampleBrokerService_IsEnabled() {
	service := BrokerService{
		Name: "left-handed-smoke-sifter",
	}

	viper.Set(service.EnabledProperty(), true)
	fmt.Println(service.IsEnabled())

	viper.Set(service.EnabledProperty(), false)
	fmt.Println(service.IsEnabled())

	// Output: true
	// false
}

func ExampleBrokerService_IsRoleWhitelistEnabled() {
	service := BrokerService{
		Name:                 "left-handed-smoke-sifter",
		DefaultRoleWhitelist: []string{"a", "b", "c"},
	}
	fmt.Println(service.IsRoleWhitelistEnabled())

	service.DefaultRoleWhitelist = nil
	fmt.Println(service.IsRoleWhitelistEnabled())

	// Output: true
	// false
}

func ExampleBrokerService_RoleWhitelist() {
	service := BrokerService{
		Name:                 "my-service",
		DefaultRoleWhitelist: []string{"a", "b", "c"},
	}
	viper.Set(service.RoleWhitelistProperty(), "")
	fmt.Println(service.RoleWhitelist())

	viper.Set(service.RoleWhitelistProperty(), "x,y,z")
	fmt.Println(service.RoleWhitelist())

	// Output: [a b c]
	// [x y z]
}

func ExampleBrokerService_TileUserDefinedPlansVariable() {
	service := BrokerService{
		Name: "google-spanner",
	}

	fmt.Println(service.TileUserDefinedPlansVariable())

	// Output: SPANNER_CUSTOM_PLANS
}

func ExampleBrokerService_ServiceDefinition() {
	service := BrokerService{
		Name: "left-handed-smoke-sifter",
		DefaultServiceDefinition: `{"id":"abcd-efgh-ijkl"}`,
	}

	// Default definition
	defn, err := service.ServiceDefinition()
	fmt.Printf("%q %v\n", defn.ID, err)

	// Override
	viper.Set(service.DefinitionProperty(), `{"id":"override-id"}`)
	defn, err = service.ServiceDefinition()
	fmt.Printf("%q %v\n", defn.ID, err)

	// Bad Value
	viper.Set(service.DefinitionProperty(), "nil")
	_, err = service.ServiceDefinition()
	fmt.Printf("%v\n", err)

	// Cleanup
	viper.Set(service.DefinitionProperty(), nil)

	// Output: "abcd-efgh-ijkl" <nil>
	// "override-id" <nil>
	// Error parsing service definition for "left-handed-smoke-sifter": invalid character 'i' in literal null (expecting 'u')
}

func ExampleBrokerService_GetPlanById() {
	service := BrokerService{
		Name: "left-handed-smoke-sifter",
		DefaultServiceDefinition: `{"id":"abcd-efgh-ijkl", "plans": [{"id": "builtin-plan", "name": "Builtin!"}]}`,
	}

	viper.Set(service.UserDefinedPlansProperty(), `[{"id":"custom-plan", "name": "Custom!"}]`)
	defer viper.Set(service.UserDefinedPlansProperty(), nil)

	plan, err := service.GetPlanById("builtin-plan")
	fmt.Printf("builtin-plan: %q %v\n", plan.Name, err)

	plan, err = service.GetPlanById("custom-plan")
	fmt.Printf("custom-plan: %q %v\n", plan.Name, err)

	_, err = service.GetPlanById("missing-plan")
	fmt.Printf("missing-plan: %s\n", err)

	// Output: builtin-plan: "Builtin!" <nil>
	// custom-plan: "Custom!" <nil>
	// missing-plan: Plan ID "missing-plan" could not be found
}

func TestBrokerService_UserDefinedPlans(t *testing.T) {
	cases := map[string]struct {
		Value       interface{}
		PlanIds     map[string]bool
		ExpectError bool
	}{
		"default-no-plans": {
			Value:       nil,
			PlanIds:     map[string]bool{},
			ExpectError: false,
		},
		"single-plan": {
			Value:       `[{"id":"aaa","name":"aaa","instances":"3"}]`,
			PlanIds:     map[string]bool{"aaa": true},
			ExpectError: false,
		},
		"bad-json": {
			Value:       `42`,
			PlanIds:     map[string]bool{},
			ExpectError: true,
		},
		"multiple-plans": {
			Value:       `[{"id":"aaa","name":"aaa","instances":"3"},{"id":"bbb","name":"bbb","instances":"3"}]`,
			PlanIds:     map[string]bool{"aaa": true, "bbb": true},
			ExpectError: false,
		},
		"missing-name": {
			Value:       `[{"id":"aaa","instances":"3"}]`,
			PlanIds:     map[string]bool{},
			ExpectError: true,
		},
		"missing-id": {
			Value:       `[{"name":"aaa","instances":"3"}]`,
			PlanIds:     map[string]bool{},
			ExpectError: true,
		},
		"missing-instances": {
			Value:       `[{"name":"aaa","id":"aaa"}]`,
			PlanIds:     map[string]bool{},
			ExpectError: true,
		},
	}

	service := BrokerService{
		Name: "left-handed-smoke-sifter",
		DefaultServiceDefinition: `{"id":"abcd-efgh-ijkl", "name":"lhss"}`,
		PlanVariables: []BrokerVariable{
			BrokerVariable{
				Required:  true,
				FieldName: "instances",
				Type:      JsonTypeString,
			},
		},
	}

	for tn, tc := range cases {
		viper.Set(service.UserDefinedPlansProperty(), tc.Value)
		plans, err := service.UserDefinedPlans()

		// Check errors
		hasErr := err != nil
		if hasErr != tc.ExpectError {
			t.Errorf("%s) Expected Error? %v, got error: %v", tn, tc.ExpectError, err)
			continue
		}

		// Check IDs
		if len(plans) != len(tc.PlanIds) {
			t.Errorf("%s) Expected %d plans, but got %d (%v)", tn, len(tc.PlanIds), len(plans), plans)
		}

		for _, plan := range plans {
			if _, ok := tc.PlanIds[plan.ID]; !ok {
				t.Errorf("%s) Got unexpected plan id %s, expected %+v", tn, plan.ID, tc.PlanIds)
			}
		}

		// Reset Environment
		viper.Set(service.UserDefinedPlansProperty(), nil)
	}
}

func TestBrokerService_CatalogEntry(t *testing.T) {
	cases := map[string]struct {
		UserDefinition interface{}
		UserPlans      interface{}
		PlanIds        map[string]bool
		ExpectError    bool
	}{
		"no-customization": {
			UserDefinition: nil,
			UserPlans:      nil,
			PlanIds:        map[string]bool{},
			ExpectError:    false,
		},
		"custom-definition": {
			UserDefinition: `{"id":"abcd-efgh-ijkl", "plans":[{"id":"zzz","name":"zzz"}]}`,
			UserPlans:      nil,
			PlanIds:        map[string]bool{"zzz": true},
			ExpectError:    false,
		},
		"custom-plans": {
			UserDefinition: nil,
			UserPlans:      `[{"id":"aaa","name":"aaa"},{"id":"bbb","name":"bbb"}]`,
			PlanIds:        map[string]bool{"aaa": true, "bbb": true},
			ExpectError:    false,
		},
		"custom-plans-and-definition": {
			UserDefinition: `{"id":"abcd-efgh-ijkl", "plans":[{"id":"zzz","name":"zzz"}]}`,
			UserPlans:      `[{"id":"aaa","name":"aaa"},{"id":"bbb","name":"bbb"}]`,
			PlanIds:        map[string]bool{"aaa": true, "bbb": true, "zzz": true},
			ExpectError:    false,
		},
		"bad-definition-json": {
			UserDefinition: `333`,
			UserPlans:      nil,
			PlanIds:        map[string]bool{},
			ExpectError:    true,
		},
		"bad-plan-json": {
			UserDefinition: nil,
			UserPlans:      `333`,
			PlanIds:        map[string]bool{},
			ExpectError:    true,
		},
	}

	service := BrokerService{
		Name: "left-handed-smoke-sifter",
		DefaultServiceDefinition: `{"id":"abcd-efgh-ijkl"}`,
	}

	for tn, tc := range cases {
		viper.Set(service.DefinitionProperty(), tc.UserDefinition)
		viper.Set(service.UserDefinedPlansProperty(), tc.UserPlans)

		srvc, err := service.CatalogEntry()
		hasErr := err != nil
		if hasErr != tc.ExpectError {
			t.Errorf("%s) Expected Error? %v, got error: %v", tn, tc.ExpectError, err)
		}

		if err == nil && len(srvc.Plans) != len(tc.PlanIds) {
			t.Errorf("%s) Expected %d plans, but got %d (%+v)", tn, len(tc.PlanIds), len(srvc.Plans), srvc.Plans)

			for _, plan := range srvc.Plans {
				if _, ok := tc.PlanIds[plan.ID]; !ok {
					t.Errorf("%s) Got unexpected plan id %s, expected %+v", tn, plan.ID, tc.PlanIds)
				}
			}
		}
	}

	viper.Set(service.DefinitionProperty(), nil)
	viper.Set(service.UserDefinedPlansProperty(), nil)
}
