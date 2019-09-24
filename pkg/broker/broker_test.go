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
	"errors"
	"fmt"
	"reflect"
	"testing"
	"os"

	"github.com/GoogleCloudPlatform/gcp-service-broker/db_service/models"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/validation"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/varcontext"
	"github.com/pivotal-cf/brokerapi"
	"github.com/spf13/viper"
)

func ExampleServiceDefinition_UserDefinedPlansProperty() {
	service := ServiceDefinition{
		Id:   "00000000-0000-0000-0000-000000000000",
		Name: "left-handed-smoke-sifter",
	}

	fmt.Println(service.UserDefinedPlansProperty())

	// Output: service.left-handed-smoke-sifter.plans
}

func ExampleServiceDefinition_IsRoleWhitelistEnabled() {
	service := ServiceDefinition{
		Id:                   "00000000-0000-0000-0000-000000000000",
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
		Id:   "00000000-0000-0000-0000-000000000000",
		Name: "google-spanner",
	}

	fmt.Println(service.TileUserDefinedPlansVariable())

	// Output: SPANNER_CUSTOM_PLANS
}

func ExampleServiceDefinition_GetPlanById() {
	service := ServiceDefinition{
		Id:   "00000000-0000-0000-0000-000000000000",
		Name: "left-handed-smoke-sifter",
		Plans: []ServicePlan{
			{ServicePlan: brokerapi.ServicePlan{ID: "builtin-plan", Name: "Builtin!"}},
		},
	}

	viper.Set(service.UserDefinedPlansProperty(), `[{"id":"custom-plan", "name": "Custom!"}]`)
	defer viper.Reset()

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
		TileValue   string
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
		"tile environment variable": {
			TileValue: `{
				"plan-100":{
					"description":"plan-100",
					"display_name":"plan-100",
					"guid":"495bf186-e1c2-4c7e-abc1-84b1a8634858",
					"instances":"100",
					"name":"plan-100",
					"service":"4bc59b9a-8520-409f-85da-1c7552315863"
				},
				"custom-plan2":{
					"description":"test",
					"display_name":"asdf",
					"guid":"938cfc91-bca3-4f9d-b384-1e4ad6f965ce",
					"instances":"10",
					"name":"custom-plan2",
					"service":"4bc59b9a-8520-409f-85da-1c7552315863"
				}
			}`,
			PlanIds: map[string]bool{
				"495bf186-e1c2-4c7e-abc1-84b1a8634858": true,
				"938cfc91-bca3-4f9d-b384-1e4ad6f965ce": true,
			},
			ExpectError: false,
		},
	}

	service := ServiceDefinition{
		Id:   "abcd-efgh-ijkl",
		Name: "left-handed-smoke-sifter",
		PlanVariables: []BrokerVariable{
			{
				Required:  true,
				FieldName: "instances",
				Type:      JsonTypeString,
			},
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			os.Setenv(service.TileUserDefinedPlansVariable(), tc.TileValue)
			defer os.Unsetenv(service.TileUserDefinedPlansVariable())

			viper.Set(service.UserDefinedPlansProperty(), tc.Value)
			defer viper.Reset()

			plans, err := service.UserDefinedPlans()

			// Check errors
			hasErr := err != nil
			if hasErr != tc.ExpectError {
				t.Fatalf("Expected Error? %v, got error: %v", tc.ExpectError, err)
			}

			// Check IDs
			if len(plans) != len(tc.PlanIds) {
				t.Errorf("Expected %d plans, but got %d (%v)", len(tc.PlanIds), len(plans), plans)
			}

			for _, plan := range plans {
				if _, ok := tc.PlanIds[plan.ID]; !ok {
					t.Errorf("Got unexpected plan id %s, expected %+v", plan.ID, tc.PlanIds)
				}
			}
		})
	}
}

func TestServiceDefinition_CatalogEntry(t *testing.T) {
	cases := map[string]struct {
		UserPlans   interface{}
		PlanIds     map[string]bool
		ExpectError bool
	}{
		"no-customization": {
			UserPlans:   nil,
			PlanIds:     map[string]bool{},
			ExpectError: false,
		},
		"custom-plans": {
			UserPlans:   `[{"id":"aaa","name":"aaa"},{"id":"bbb","name":"bbb"}]`,
			PlanIds:     map[string]bool{"aaa": true, "bbb": true},
			ExpectError: false,
		},
		"bad-plan-json": {
			UserPlans:   `333`,
			PlanIds:     map[string]bool{},
			ExpectError: true,
		},
	}

	service := ServiceDefinition{
		Id:   "00000000-0000-0000-0000-000000000000",
		Name: "left-handed-smoke-sifter",
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			viper.Set(service.UserDefinedPlansProperty(), tc.UserPlans)
			defer viper.Reset()

			srvc, err := service.CatalogEntry()
			hasErr := err != nil
			if hasErr != tc.ExpectError {
				t.Errorf("Expected Error? %v, got error: %v", tc.ExpectError, err)
			}

			if err == nil && len(srvc.Plans) != len(tc.PlanIds) {
				t.Errorf("Expected %d plans, but got %d (%+v)", len(tc.PlanIds), len(srvc.Plans), srvc.Plans)

				for _, plan := range srvc.Plans {
					if _, ok := tc.PlanIds[plan.ID]; !ok {
						t.Errorf("Got unexpected plan id %s, expected %+v", plan.ID, tc.PlanIds)
					}
				}
			}
		})
	}
}

func ExampleServiceDefinition_CatalogEntry() {
	service := ServiceDefinition{
		Id:   "00000000-0000-0000-0000-000000000000",
		Name: "left-handed-smoke-sifter",
		Plans: []ServicePlan{
			{ServicePlan: brokerapi.ServicePlan{ID: "builtin-plan", Name: "Builtin!"}},
		},
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
	defer viper.Reset()

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
		Id:   "00000000-0000-0000-0000-000000000000",
		Name: "left-handed-smoke-sifter",
		Plans: []ServicePlan{
			{ServicePlan: brokerapi.ServicePlan{ID: "builtin-plan", Name: "Builtin!"}},
		},
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
				Constraints: validation.NewConstraintBuilder().
					MaxLength(30).
					Build(),
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
		UserParams         string
		ServiceProperties  map[string]string
		DefaultOverride    string
		ProvisionOverrides map[string]interface{}
		ExpectedError      error
		ExpectedContext    map[string]interface{}
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
		"invalid-request": {
			UserParams:    `{"name":"some-name-that-is-longer-than-thirty-characters"}`,
			ExpectedError: errors.New("1 error(s) occurred: name: String length must be less than or equal to 30"),
		},
		"provision_overrides override user params but not computed defaults": {
			UserParams:         `{"location":"us"}`,
			DefaultOverride:    "{}",
			ServiceProperties:  map[string]string{},
			ProvisionOverrides: map[string]interface{}{"location": "eu"},
			ExpectedContext: map[string]interface{}{
				"location":      "eu",
				"name":          "name-eu",
				"maybe-missing": "default",
			},
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			viper.Set(service.ProvisionDefaultOverrideProperty(), tc.DefaultOverride)
			defer viper.Reset()

			details := brokerapi.ProvisionDetails{RawParameters: json.RawMessage(tc.UserParams)}
			plan := ServicePlan{ServiceProperties: tc.ServiceProperties, ProvisionOverrides: tc.ProvisionOverrides}
			vars, err := service.ProvisionVariables("instance-id-here", details, plan)

			expectError(t, tc.ExpectedError, err)

			if tc.ExpectedError == nil && !reflect.DeepEqual(vars.ToMap(), tc.ExpectedContext) {
				t.Errorf("Expected context: %v got %v", tc.ExpectedContext, vars.ToMap())
			}
		})
	}
}

func TestServiceDefinition_BindVariables(t *testing.T) {
	service := ServiceDefinition{
		Id:   "00000000-0000-0000-0000-000000000000",
		Name: "left-handed-smoke-sifter",
		Plans: []ServicePlan{
			{
				ServicePlan: brokerapi.ServicePlan{
					ID:   "builtin-plan",
					Name: "Builtin!",
				},
				ServiceProperties: map[string]string{
					"service-property": "operator-set",
				},
			},
		},
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
				Constraints: validation.NewConstraintBuilder().
					MaxLength(30).
					Build(),
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
			{
				Name:      "service-prop",
				Default:   `${request.plan_properties["service-property"]}`,
				Overwrite: true,
			},
		},
	}

	cases := map[string]struct {
		UserParams      string
		DefaultOverride string
		BindOverrides   map[string]interface{}
		ExpectedError   error
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
				"service-prop": "operator-set",
			},
		},
		"location gets truncated": {
			UserParams:   `{"location": "averylonglocation"}`,
			InstanceVars: `{"foo":"default"}`,
			ExpectedContext: map[string]interface{}{
				"location":     "averylongl",
				"name":         "name-averylonglocation",
				"instance-foo": "default",
				"service-prop": "operator-set",
			},
		},
		"user location and name": {
			UserParams:   `{"location": "eu", "name":"foo"}`,
			InstanceVars: `{"foo":"default"}`,
			ExpectedContext: map[string]interface{}{
				"location":     "eu",
				"name":         "foo",
				"instance-foo": "default",
				"service-prop": "operator-set",
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
				"service-prop": "operator-set",
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
				"service-prop": "operator-set",
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
				"service-prop": "operator-set",
			},
		},
		"instance info can get parsed": {
			UserParams:   `{"location":"us"}`,
			InstanceVars: `{"foo":"bar"}`,
			ExpectedContext: map[string]interface{}{
				"location":     "us",
				"name":         "name-us",
				"instance-foo": "bar",
				"service-prop": "operator-set",
			},
		},
		"invalid-request": {
			UserParams:    `{"name":"some-name-that-is-longer-than-thirty-characters"}`,
			InstanceVars:  `{"foo":""}`,
			ExpectedError: errors.New("1 error(s) occurred: name: String length must be less than or equal to 30"),
		},
		"bind_overrides override user params but not computed defaults": {
			UserParams:    `{"location":"us"}`,
			InstanceVars:  `{"foo":"default"}`,
			BindOverrides: map[string]interface{}{"location": "eu"},
			ExpectedContext: map[string]interface{}{
				"location":     "eu",
				"name":         "name-eu",
				"instance-foo": "default",
				"service-prop": "operator-set",
			},
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			viper.Set(service.BindDefaultOverrideProperty(), tc.DefaultOverride)
			defer viper.Reset()

			details := brokerapi.BindDetails{RawParameters: json.RawMessage(tc.UserParams)}
			instance := models.ServiceInstanceDetails{OtherDetails: tc.InstanceVars}
			service.Plans[0].BindOverrides = tc.BindOverrides
			vars, err := service.BindVariables(instance, "binding-id-here", details, &service.Plans[0])

			expectError(t, tc.ExpectedError, err)

			if tc.ExpectedError == nil && !reflect.DeepEqual(vars.ToMap(), tc.ExpectedContext) {
				t.Errorf("Expected context: %v got %v", tc.ExpectedContext, vars.ToMap())
			}
		})
	}
}

func TestServiceDefinition_createSchemas(t *testing.T) {
	service := ServiceDefinition{
		Id:   "00000000-0000-0000-0000-000000000000",
		Name: "left-handed-smoke-sifter",
		Plans: []ServicePlan{
			{ServicePlan: brokerapi.ServicePlan{ID: "builtin-plan", Name: "Builtin!"}},
		},
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

	expectedCreateParams := CreateJsonSchema(service.ProvisionInputVariables)
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

	expectedBindCreateParams := CreateJsonSchema(service.BindInputVariables)
	if !reflect.DeepEqual(bindCreate.Parameters, expectedBindCreateParams) {
		t.Errorf("expected create params to be: %v got %v", expectedBindCreateParams, bindCreate.Parameters)
	}
}

func expectError(t *testing.T, expected, actual error) {
	t.Helper()
	expectedErr := expected != nil
	gotErr := actual != nil

	switch {
	case expectedErr && gotErr:
		if expected.Error() != actual.Error() {
			t.Fatalf("Expected: %v, got: %v", expected, actual)
		}
	case expectedErr && !gotErr:
		t.Fatalf("Expected: %v, got: %v", expected, actual)
	case !expectedErr && gotErr:
		t.Fatalf("Expected no error but got: %v", actual)
	}
}
