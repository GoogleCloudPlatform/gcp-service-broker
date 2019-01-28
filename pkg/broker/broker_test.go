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

	"github.com/GoogleCloudPlatform/gcp-service-broker/db_service/models"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/validation"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/varcontext"
	"github.com/pivotal-cf/brokerapi"
	"github.com/spf13/viper"
)

func ExampleServiceDefinition_GetPlanById() {
	service := ServiceDefinition{
		Id:   "00000000-0000-0000-0000-000000000000",
		Name: "left-handed-smoke-sifter",
		Plans: []ServicePlan{
			{ServicePlan: brokerapi.ServicePlan{ID: "builtin-plan", Name: "Builtin!"}},
		},
		config: ServiceConfig{
			CustomPlans: []CustomPlan{
				{
					GUID: "custom-plan",
					Name: "Custom!",
				},
			},
		},
	}

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

func TestServiceDefinition_SetConfig(t *testing.T) {
	cases := map[string]struct {
		Value       []CustomPlan
		ExpectError error
	}{
		"default-no-plans": {
			Value:       nil,
			ExpectError: nil,
		},
		"single-plan": {
			Value: []CustomPlan{
				{GUID: "aaa", Name: "aaa", Properties: map[string]string{"instances": "3"}},
			},
			ExpectError: nil,
		},
		"multiple-plans": {
			Value: []CustomPlan{
				{GUID: "aaa", Name: "aaa", Properties: map[string]string{"instances": "3"}},
				{GUID: "bbb", Name: "bbb", Properties: map[string]string{"instances": "3"}},
			},
			ExpectError: nil,
		},
		"missing-name": {
			Value: []CustomPlan{
				{GUID: "aaa", Properties: map[string]string{"instances": "3"}},
			},
			ExpectError: errors.New("left-handed-smoke-sifter custom_plans[0] is missing a name"),
		},
		"missing-id": {
			Value: []CustomPlan{
				{Name: "aaa", Properties: map[string]string{"instances": "3"}},
			},
			ExpectError: errors.New("left-handed-smoke-sifter custom_plans[0] is missing an id"),
		},
		"missing-instances": {
			Value: []CustomPlan{
				{Name: "aaa", GUID: "aaa"},
			},
			ExpectError: errors.New("left-handed-smoke-sifter custom_plans[0] is missing required property instances"),
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
			{
				Required:  false,
				FieldName: "optional_field",
				Type:      JsonTypeString,
			},
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			config := ServiceConfig{CustomPlans: tc.Value}
			err := service.SetConfig(config)
			defer service.SetConfig(ServiceConfig{})

			if !reflect.DeepEqual(err, tc.ExpectError) {
				t.Fatalf("Expected Error? %v, got error: %v", tc.ExpectError, err)
			}
		})
	}
}

func TestServiceDefinition_CatalogEntry(t *testing.T) {
	cases := map[string]struct {
		UserPlans   []CustomPlan
		PlanIds     map[string]bool
		ExpectError bool
	}{
		"no-customization": {
			UserPlans:   nil,
			PlanIds:     map[string]bool{},
			ExpectError: false,
		},
		"custom-plans": {
			UserPlans: []CustomPlan{
				{
					GUID:       "aaa",
					Name:       "aaa",
					Properties: map[string]string{"instances": "3"},
				},
				{
					GUID:       "bbb",
					Name:       "bbb",
					Properties: map[string]string{"instances": "3"},
				},
			},
			PlanIds:     map[string]bool{"aaa": true, "bbb": true},
			ExpectError: false,
		},
	}

	service := ServiceDefinition{
		Id:   "00000000-0000-0000-0000-000000000000",
		Name: "left-handed-smoke-sifter",
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			err := service.SetConfig(ServiceConfig{
				CustomPlans: tc.UserPlans,
			})

			hasErr := err != nil
			if hasErr != tc.ExpectError {
				t.Errorf("Expected Error? %v, got error: %v", tc.ExpectError, err)
			}

			srvc := service.CatalogEntry()
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

	// Schemas should be nil by default
	srvc := service.CatalogEntry()
	fmt.Println("schemas with flag off:", srvc.ToPlain().Plans[0].Schemas)

	viper.Set("compatibility.enable-catalog-schemas", true)
	defer viper.Reset()

	srvc = service.CatalogEntry()
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
		UserParams        string
		ServiceProperties map[string]string
		DefaultOverride   map[string]interface{}
		ExpectedError     error
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
			DefaultOverride:   map[string]interface{}{"location": "eu"},
			ServiceProperties: map[string]string{},
			ExpectedContext: map[string]interface{}{
				"location":      "eu",
				"name":          "name-eu",
				"maybe-missing": "default",
			},
		},
		"user values override operator defaults": {
			UserParams:        `{"location":"nz"}`,
			DefaultOverride:   map[string]interface{}{"location": "eu"},
			ServiceProperties: map[string]string{},
			ExpectedContext: map[string]interface{}{
				"location":      "nz",
				"name":          "name-nz",
				"maybe-missing": "default",
			},
		},
		"operator defaults are not evaluated": {
			UserParams:        `{"location":"us"}`,
			DefaultOverride:   map[string]interface{}{"name": "foo-${location}"},
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
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			if err := service.SetConfig(ServiceConfig{ProvisionDefaults: tc.DefaultOverride}); err != nil {
				t.Fatal(err)
			}

			details := brokerapi.ProvisionDetails{RawParameters: json.RawMessage(tc.UserParams)}
			plan := ServicePlan{ServiceProperties: tc.ServiceProperties}
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
			{ServicePlan: brokerapi.ServicePlan{ID: "builtin-plan", Name: "Builtin!"}},
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
		},
	}

	cases := map[string]struct {
		UserParams      string
		DefaultOverride map[string]interface{}
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
			DefaultOverride: map[string]interface{}{"location": "eu"},
			ExpectedContext: map[string]interface{}{
				"location":     "eu",
				"name":         "name-eu",
				"instance-foo": "default",
			},
		},
		"user values override operator defaults": {
			UserParams:      `{"location":"nz"}`,
			InstanceVars:    `{"foo":"default"}`,
			DefaultOverride: map[string]interface{}{"location": "eu"},
			ExpectedContext: map[string]interface{}{
				"location":     "nz",
				"name":         "name-nz",
				"instance-foo": "default",
			},
		},
		"operator defaults are not evaluated": {
			UserParams:      `{"location":"us"}`,
			InstanceVars:    `{"foo":"default"}`,
			DefaultOverride: map[string]interface{}{"name": "foo-${location}"},
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
		"invalid-request": {
			UserParams:    `{"name":"some-name-that-is-longer-than-thirty-characters"}`,
			InstanceVars:  `{"foo":""}`,
			ExpectedError: errors.New("1 error(s) occurred: name: String length must be less than or equal to 30"),
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			if err := service.SetConfig(ServiceConfig{BindDefaults: tc.DefaultOverride}); err != nil {
				t.Fatal(err)
			}
			defer service.SetConfig(ServiceConfig{})

			details := brokerapi.BindDetails{RawParameters: json.RawMessage(tc.UserParams)}
			instance := models.ServiceInstanceDetails{OtherDetails: tc.InstanceVars}
			vars, err := service.BindVariables(instance, "binding-id-here", details)

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
