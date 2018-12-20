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
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/models"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/varcontext"
	"github.com/pivotal-cf/brokerapi"
	"github.com/spf13/viper"
)

func ExampleServiceDefinition_EnabledProperty() {
	service := ServiceDefinition{
		Name: "left-handed-smoke-sifter",
	}

	fmt.Println(service.EnabledProperty())

	// Output: service.left-handed-smoke-sifter.enabled
}

func ExampleServiceDefinition_DefinitionProperty() {
	service := ServiceDefinition{
		Name: "left-handed-smoke-sifter",
	}

	fmt.Println(service.DefinitionProperty())

	// Output: service.left-handed-smoke-sifter.definition
}

func ExampleServiceDefinition_UserDefinedPlansProperty() {
	service := ServiceDefinition{
		Name: "left-handed-smoke-sifter",
	}

	fmt.Println(service.UserDefinedPlansProperty())

	// Output: service.left-handed-smoke-sifter.plans
}

func ExampleServiceDefinition_IsEnabled() {
	service := ServiceDefinition{
		Name: "left-handed-smoke-sifter",
	}

	viper.Set(service.EnabledProperty(), true)
	fmt.Println(service.IsEnabled())

	viper.Set(service.EnabledProperty(), false)
	fmt.Println(service.IsEnabled())

	// Output: true
	// false
}

func ExampleServiceDefinition_IsRoleWhitelistEnabled() {
	service := ServiceDefinition{
		Name:                 "left-handed-smoke-sifter",
		DefaultRoleWhitelist: []string{"a", "b", "c"},
	}
	fmt.Println(service.IsRoleWhitelistEnabled())

	service.DefaultRoleWhitelist = nil
	fmt.Println(service.IsRoleWhitelistEnabled())

	// Output: true
	// false
}

func ExampleServiceDefinition_TileUserDefinedPlansVariable() {
	service := ServiceDefinition{
		Name: "google-spanner",
	}

	fmt.Println(service.TileUserDefinedPlansVariable())

	// Output: SPANNER_CUSTOM_PLANS
}

func ExampleServiceDefinition_ServiceDefinition() {
	service := ServiceDefinition{
		Name:                     "left-handed-smoke-sifter",
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

func ExampleServiceDefinition_GetPlanById() {
	service := ServiceDefinition{
		Name:                     "left-handed-smoke-sifter",
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

func TestServiceDefinition_UserDefinedPlans(t *testing.T) {
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

	service := ServiceDefinition{
		Name:                     "left-handed-smoke-sifter",
		DefaultServiceDefinition: `{"id":"abcd-efgh-ijkl", "name":"lhss"}`,
		PlanVariables: []BrokerVariable{
			{
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

func TestServiceDefinition_CatalogEntry(t *testing.T) {
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

	service := ServiceDefinition{
		Name:                     "left-handed-smoke-sifter",
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

func ExampleServiceDefinition_CatalogEntrySchema() {
	service := ServiceDefinition{
		Name:                     "left-handed-smoke-sifter",
		DefaultServiceDefinition: `{"id":"abcd-efgh-ijkl", "plans": [{"id": "builtin-plan", "name": "Builtin!"}]}`,
		ProvisionInputVariables: []BrokerVariable{
			{FieldName: "location", Type: JsonTypeString, Default: "us"},
		},
		BindInputVariables: []BrokerVariable{
			{FieldName: "name", Type: JsonTypeString, Default: "name"},
		},
	}

	srvc, err := service.CatalogEntry()
	if err != nil {
		panic(err)
	}

	// Schemas should be nil by default
	fmt.Println("schemas with flag off:", srvc.ToPlain().Plans[0].Schemas)

	viper.Set("compatibility.enable-catalog-schemas", true)
	defer viper.Set("compatibility.enable-catalog-schemas", false)

	srvc, err = service.CatalogEntry()
	if err != nil {
		panic(err)
	}

	eq := reflect.DeepEqual(srvc.ToPlain().Plans[0].Schemas, service.createSchemas())

	fmt.Println("schema was generated?", eq)

	// Output: schemas with flag off: <nil>
	// schema was generated? true
}

func TestServiceDefinition_ProvisionVariables(t *testing.T) {
	service := ServiceDefinition{
		Name:                     "left-handed-smoke-sifter",
		DefaultServiceDefinition: `{"id":"abcd-efgh-ijkl", "plans": [{"id": "builtin-plan", "name": "Builtin!"}]}`,
		ProvisionInputVariables: []BrokerVariable{
			{
				FieldName: "location",
				Type:      JsonTypeString,
				Default:   "us",
			},
			{
				FieldName: "name",
				Type:      JsonTypeString,
				Default:   "name-${location}",
			},
		},
		ProvisionComputedVariables: []varcontext.DefaultVariable{
			{
				Name:      "location",
				Default:   "${str.truncate(10, location)}",
				Overwrite: true,
			},
			{
				Name:      "maybe-missing",
				Default:   "default",
				Overwrite: false,
			},
		},
	}

	cases := map[string]struct {
		UserParams        string
		ServiceProperties map[string]string
		DefaultOverride   string
		ExpectedContext   map[string]interface{}
	}{
		"empty": {
			UserParams:        "",
			ServiceProperties: map[string]string{},
			ExpectedContext: map[string]interface{}{
				"location":      "us",
				"name":          "name-us",
				"maybe-missing": "default",
			},
		},
		"service has missing param": {
			UserParams:        "",
			ServiceProperties: map[string]string{"maybe-missing": "custom"},
			ExpectedContext: map[string]interface{}{
				"location":      "us",
				"name":          "name-us",
				"maybe-missing": "custom",
			},
		},
		"location gets truncated": {
			UserParams:        `{"location": "averylonglocation"}`,
			ServiceProperties: map[string]string{},
			ExpectedContext: map[string]interface{}{
				"location":      "averylongl",
				"name":          "name-averylonglocation",
				"maybe-missing": "default",
			},
		},
		"user location and name": {
			UserParams:        `{"location": "eu", "name":"foo"}`,
			ServiceProperties: map[string]string{},
			ExpectedContext: map[string]interface{}{
				"location":      "eu",
				"name":          "foo",
				"maybe-missing": "default",
			},
		},
		"user tries to overwrite service var": {
			UserParams:        `{"location": "eu", "name":"foo", "service-provided":"test"}`,
			ServiceProperties: map[string]string{"service-provided": "custom"},
			ExpectedContext: map[string]interface{}{
				"location":         "eu",
				"name":             "foo",
				"maybe-missing":    "default",
				"service-provided": "custom",
			},
		},
		"operator defaults override computed defaults": {
			UserParams:        "",
			DefaultOverride:   `{"location":"eu"}`,
			ServiceProperties: map[string]string{},
			ExpectedContext: map[string]interface{}{
				"location":      "eu",
				"name":          "name-eu",
				"maybe-missing": "default",
			},
		},
		"user values override operator defaults": {
			UserParams:        `{"location":"nz"}`,
			DefaultOverride:   `{"location":"eu"}`,
			ServiceProperties: map[string]string{},
			ExpectedContext: map[string]interface{}{
				"location":      "nz",
				"name":          "name-nz",
				"maybe-missing": "default",
			},
		},
		"operator defaults are not evaluated": {
			UserParams:        `{"location":"us"}`,
			DefaultOverride:   `{"name":"foo-${location}"}`,
			ServiceProperties: map[string]string{},
			ExpectedContext: map[string]interface{}{
				"location":      "us",
				"name":          "foo-${location}",
				"maybe-missing": "default",
			},
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			viper.Set(service.ProvisionDefaultOverrideProperty(), tc.DefaultOverride)
			details := brokerapi.ProvisionDetails{RawParameters: json.RawMessage(tc.UserParams)}
			plan := ServicePlan{ServiceProperties: tc.ServiceProperties}
			vars, err := service.ProvisionVariables("instance-id-here", details, plan)

			if err != nil {
				t.Errorf("got error while creating provision variables: %v", err)
			}

			if !reflect.DeepEqual(vars.ToMap(), tc.ExpectedContext) {
				t.Errorf("Expected context: %v got %v", tc.ExpectedContext, vars.ToMap())
			}
		})
	}
}

func TestServiceDefinition_BindVariables(t *testing.T) {
	service := ServiceDefinition{
		Name:                     "left-handed-smoke-sifter",
		DefaultServiceDefinition: `{"id":"abcd-efgh-ijkl", "plans": [{"id": "builtin-plan", "name": "Builtin!"}]}`,
		BindInputVariables: []BrokerVariable{
			{
				FieldName: "location",
				Type:      JsonTypeString,
				Default:   "us",
			},
			{
				FieldName: "name",
				Type:      JsonTypeString,
				Default:   "name-${location}",
			},
		},
		BindComputedVariables: []varcontext.DefaultVariable{
			{
				Name:      "location",
				Default:   "${str.truncate(10, location)}",
				Overwrite: true,
			},
			{
				Name:      "instance-foo",
				Default:   `${instance.details["foo"]}`,
				Overwrite: true,
			},
		},
	}

	cases := map[string]struct {
		UserParams      string
		DefaultOverride string
		ExpectedContext map[string]interface{}
		InstanceVars    string
	}{
		"empty": {
			UserParams:   "",
			InstanceVars: `{"foo":""}`,
			ExpectedContext: map[string]interface{}{
				"location":     "us",
				"name":         "name-us",
				"instance-foo": "",
			},
		},
		"location gets truncated": {
			UserParams:   `{"location": "averylonglocation"}`,
			InstanceVars: `{"foo":"default"}`,
			ExpectedContext: map[string]interface{}{
				"location":     "averylongl",
				"name":         "name-averylonglocation",
				"instance-foo": "default",
			},
		},
		"user location and name": {
			UserParams:   `{"location": "eu", "name":"foo"}`,
			InstanceVars: `{"foo":"default"}`,
			ExpectedContext: map[string]interface{}{
				"location":     "eu",
				"name":         "foo",
				"instance-foo": "default",
			},
		},
		"operator defaults override computed defaults": {
			UserParams:      "",
			InstanceVars:    `{"foo":"default"}`,
			DefaultOverride: `{"location":"eu"}`,
			ExpectedContext: map[string]interface{}{
				"location":     "eu",
				"name":         "name-eu",
				"instance-foo": "default",
			},
		},
		"user values override operator defaults": {
			UserParams:      `{"location":"nz"}`,
			InstanceVars:    `{"foo":"default"}`,
			DefaultOverride: `{"location":"eu"}`,
			ExpectedContext: map[string]interface{}{
				"location":     "nz",
				"name":         "name-nz",
				"instance-foo": "default",
			},
		},
		"operator defaults are not evaluated": {
			UserParams:      `{"location":"us"}`,
			InstanceVars:    `{"foo":"default"}`,
			DefaultOverride: `{"name":"foo-${location}"}`,
			ExpectedContext: map[string]interface{}{
				"location":     "us",
				"name":         "foo-${location}",
				"instance-foo": "default",
			},
		},
		"instance info can get parsed": {
			UserParams:   `{"location":"us"}`,
			InstanceVars: `{"foo":"bar"}`,
			ExpectedContext: map[string]interface{}{
				"location":     "us",
				"name":         "name-us",
				"instance-foo": "bar",
			},
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			viper.Set(service.BindDefaultOverrideProperty(), tc.DefaultOverride)
			details := brokerapi.BindDetails{RawParameters: json.RawMessage(tc.UserParams)}
			instance := models.ServiceInstanceDetails{OtherDetails: tc.InstanceVars}
			vars, err := service.BindVariables(instance, "binding-id-here", details)

			if err != nil {
				t.Fatalf("got error while creating provision variables: %v", err)
			}

			if !reflect.DeepEqual(vars.ToMap(), tc.ExpectedContext) {
				t.Errorf("Expected context: %v got %v", tc.ExpectedContext, vars.ToMap())
			}
		})
	}
}

func TestServiceDefinition_createSchemas(t *testing.T) {
	service := ServiceDefinition{
		Name:                     "left-handed-smoke-sifter",
		DefaultServiceDefinition: `{"id":"abcd-efgh-ijkl", "plans": [{"id": "builtin-plan", "name": "Builtin!"}]}`,
		ProvisionInputVariables: []BrokerVariable{
			{FieldName: "location", Type: JsonTypeString, Default: "us"},
		},
		BindInputVariables: []BrokerVariable{
			{FieldName: "name", Type: JsonTypeString, Default: "name"},
		},
	}

	schemas := service.createSchemas()
	if schemas == nil {
		t.Fatal("Schemas was nil, expected non-nil value")
	}

	// it populates the instance create schema with the fields in ProvisionInputVariables
	instanceCreate := schemas.Instance.Create
	if instanceCreate.Parameters == nil {
		t.Error("instance create params were nil, expected a schema")
	}

	expectedCreateParams := createJsonSchema(service.ProvisionInputVariables)
	if !reflect.DeepEqual(instanceCreate.Parameters, expectedCreateParams) {
		t.Errorf("expected create params to be: %v got %v", expectedCreateParams, instanceCreate.Parameters)
	}

	// It leaves the instance update schema blank.
	instanceUpdate := schemas.Instance.Update
	if instanceUpdate.Parameters != nil {
		t.Error("instance update params were not nil, expected nil")
	}

	// it populates the binding create schema with the fields in BindInputVariables.
	bindCreate := schemas.Binding.Create
	if bindCreate.Parameters == nil {
		t.Error("bind create params were not nil, expected a schema")
	}

	expectedBindCreateParams := createJsonSchema(service.BindInputVariables)
	if !reflect.DeepEqual(bindCreate.Parameters, expectedBindCreateParams) {
		t.Errorf("expected create params to be: %v got %v", expectedBindCreateParams, bindCreate.Parameters)
	}
}
