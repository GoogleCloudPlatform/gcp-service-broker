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
	"reflect"
	"testing"
	"errors"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/validation"
)

func TestBrokerVariable_ToSchema(t *testing.T) {
	cases := map[string]struct {
		BrokerVar BrokerVariable
		Expected  map[string]interface{}
	}{
		"blank": {
			BrokerVariable{}, map[string]interface{}{},
		},
		"enums get copied": {
			BrokerVariable{Enum: map[interface{}]string{"a": "description", "b": "description"}},
			map[string]interface{}{
				"enum": []interface{}{"a", "b"},
			},
		},
		"details are copied": {
			BrokerVariable{Details: "more information"},
			map[string]interface{}{
				"description": "more information",
			},
		},
		"type is copied": {
			BrokerVariable{Type: JsonTypeString},
			map[string]interface{}{
				"type": JsonTypeString,
			},
		},
		"default is copied": {
			BrokerVariable{Default: "some-value"},
			map[string]interface{}{
				"default": "some-value",
			},
		},
		"full test": {
			BrokerVariable{
				Default: "some-value",
				Type:    JsonTypeString,
				Details: "more information",
				Enum:    map[interface{}]string{"b": "description", "a": "description"},
				Constraints: map[string]interface{}{
					"examples": []string{"SAMPLEA", "SAMPLEB"},
				},
			},
			map[string]interface{}{
				"default":     "some-value",
				"type":        JsonTypeString,
				"description": "more information",
				"enum":        []interface{}{"a", "b"},
				"examples":    []string{"SAMPLEA", "SAMPLEB"},
			},
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			actual := tc.BrokerVar.ToSchema()
			if !reflect.DeepEqual(actual, tc.Expected) {
				t.Errorf("Expected ToSchema to be: %v, got: %v", tc.Expected, actual)
			}
		})
	}
}

func TestBrokerVariable_ValidateVariables(t *testing.T) {
	cases := map[string]struct {
		Parameters map[string]interface{}
		Variables []BrokerVariable
		Expected  error
	} {
		"nil check": {
			Parameters: nil,
			Variables: nil,
			Expected: nil,
		},
		"integer": {
			Parameters: map[string]interface{}{
				"test":12,
			},
			Variables: []BrokerVariable{
				{
					Required:true,
					FieldName:"test",
					Type:JsonTypeInteger,
				},
			},
			Expected: nil,
		},
		"unexpected type": {
			Parameters: map[string]interface{}{
				"test":"didn't see that coming",
			},
			Variables: []BrokerVariable{
				{
					Required:true,
					FieldName:"test",
					Type:JsonTypeInteger,
				},
			},
			Expected: errors.New("1 error(s) occurred: test: Invalid type. Expected: integer, given: string"),
		},
		"double trouble": {
			Parameters: map[string]interface{}{
				"test":"didn't see that coming",
				"test2":"I am no good",
			},
			Variables: []BrokerVariable{
				{
					Required:true,
					FieldName:"test",
					Type:JsonTypeInteger,
				},
				{
					Required:true,
					FieldName:"test2",
					Type:JsonTypeInteger,
				},
			},
			Expected: errors.New("2 error(s) occurred: test: Invalid type. Expected: integer, given: string; test2: Invalid type. Expected: integer, given: string"),
		},
		"test constraints": {
			Parameters: map[string]interface{}{
				"test":0,
			},
			Variables: []BrokerVariable{
				{
					Required:true,
					FieldName:"test",
					Type:JsonTypeInteger,
					Constraints: validation.NewConstraintBuilder().
						Minimum(10).
						Build(),
				},
			},
			Expected: errors.New("1 error(s) occurred: test: Must be greater than or equal to 10"),
		},
		"test enum": {
			Parameters: map[string]interface{}{
				"test":"not this one",
			},
			Variables: []BrokerVariable{
				{
					Required:true,
					FieldName:"test",
					Type:JsonTypeString,
					Enum: map[interface{}]string {
						"one": "it's either this one",
						"theother": "or this one",
					},
				},
			},
			Expected: errors.New("1 error(s) occurred: test: test must be one of the following: \"one\", \"theother\""),
		},
		"test missing": {
			Parameters: map[string]interface{}{},
			Variables: []BrokerVariable{
				{
					Required:true,
					FieldName:"test",
					Type:JsonTypeString,
					Enum: map[interface{}]string {
						"one": "it's either this one",
						"theother": "or this one",
					},
				},
			},
			Expected: errors.New("1 error(s) occurred: test: test is required"),
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			actual := ValidateVariables(tc.Parameters, tc.Variables)

			if !reflect.DeepEqual(actual, tc.Expected) {
				if actual == nil {
					t.Errorf("Expected ValidateVariables to be: %v, got: %v", tc.Expected, actual)
				} else if actual.Error() != tc.Expected.Error() {
					t.Errorf("Expected ValidateVariables to be: %v, got: %v", tc.Expected, actual.Error())
				}
			}
		})
	}
}

func TestBrokerVariable_ApplyDefaults(t *testing.T) {
	cases := map[string]struct {
		Parameters map[string]interface{}
		Variables []BrokerVariable
		Expected  map[string]interface{}
	} {
		"nil check": {
			Parameters: nil,
			Variables: nil,
			Expected: nil,
		},
		"simple": {
			Parameters: map[string]interface{}{},
			Variables: []BrokerVariable{
				{
					FieldName:"test",
					Type:JsonTypeInteger,
					Default:123,
				},
			},
			Expected: map[string]interface{}{
				"test":123,
			},
		},
		"do not replace": {
			Parameters: map[string]interface{}{
				"test":123,
			},
			Variables: []BrokerVariable{
				{
					FieldName:"test",
					Type:JsonTypeInteger,
					Default:456,
				},
			},
			Expected: map[string]interface{}{
				"test":123,
			},
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			ApplyDefaults(tc.Parameters, tc.Variables)

			if !reflect.DeepEqual(tc.Parameters, tc.Expected) {

				t.Errorf("Expected ValidateVariables to be: %v, got: %v", tc.Expected, tc.Parameters)
			}
		})
	}
}
