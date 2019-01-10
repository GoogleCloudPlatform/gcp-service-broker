// Copyright 2019 the Service Broker Project Authors.
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

package policy

import (
	"errors"
	"reflect"
	"testing"
)

func TestCondition_AppliesTo(t *testing.T) {
	cases := map[string]struct {
		Condition Condition
		Truth     Condition
		Expected  bool
	}{
		"blank-condition": {
			Condition: Condition{},
			Truth:     Condition{"service_id": "my-service-id", "service_name": "service-name"},
			Expected:  true,
		},
		"partial-condition": {
			Condition: Condition{"service_id": "my-service-id"},
			Truth:     Condition{"service_id": "my-service-id", "service_name": "service-name"},
			Expected:  true,
		},
		"mismatching-condition": {
			Condition: Condition{"service_id": "abc"},
			Truth:     Condition{"service_id": "my-service-id", "service_name": "service-name"},
			Expected:  false,
		},
		"key-not-in-truth": {
			Condition: Condition{"service_id": ""},
			Truth:     Condition{},
			Expected:  false,
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			actual := tc.Condition.AppliesTo(tc.Truth)

			if tc.Expected != actual {
				t.Errorf("Expected condition to apply? %t but was: %t", tc.Expected, actual)
			}
		})
	}
}

func TestCondition_ValidateKeys(t *testing.T) {
	cases := map[string]struct {
		Condition   Condition
		AllowedKeys []string
		Expected    error
	}{
		"blank-condition": {
			Condition:   Condition{},
			AllowedKeys: []string{"service_id"},
			Expected:    nil,
		},
		"good-condition": {
			Condition:   Condition{"service_id": "my-service-id"},
			AllowedKeys: []string{"service_id"},
			Expected:    nil,
		},
		"key-mismatch": {
			Condition:   Condition{"service_name": "abc"},
			AllowedKeys: []string{"service_id"},
			Expected:    errors.New("unknown condition keys: [service_name] condition keys must one of: [service_id], check their capitalization and spelling"),
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			actual := tc.Condition.ValidateKeys(tc.AllowedKeys)

			if !reflect.DeepEqual(tc.Expected, actual) {
				t.Errorf("Expected error: %v got %v", tc.Expected, actual)
			}
		})
	}
}

func TestPolicyList_Validate(t *testing.T) {
	cases := map[string]struct {
		Policy      PolicyList
		AllowedKeys []string
		Expected    error
	}{
		"blank-policy": {
			Policy:      PolicyList{},
			AllowedKeys: []string{"a", "b"},
			Expected:    nil,
		},
		"good-policy": {
			Policy: PolicyList{
				Policies: []Policy{
					{Condition: Condition{"a": "a-value"}, Declarations: map[string]interface{}{"a-fired": true}},
					{Condition: Condition{"b": "b-value"}, Declarations: map[string]interface{}{"b-fired": true}},
				},
				Assertions: []Policy{
					{Condition: Condition{"a": "a-value", "b": "b-value"}, Declarations: map[string]interface{}{"a-fired": true, "b-fired": true}},
				},
			},
			AllowedKeys: []string{"a", "b"},
			Expected:    nil,
		},

		"bad-keys": {
			Policy: PolicyList{
				Policies: []Policy{
					{Condition: Condition{"unknown": "a-value"}, Comment: "some-user-comment"},
				},
			},
			AllowedKeys: []string{"a", "b"},
			Expected:    errors.New(`error in policy[0], comment: "some-user-comment", error: unknown condition keys: [unknown] condition keys must one of: [a b], check their capitalization and spelling`),
		},
		"bad-assertion": {
			Policy: PolicyList{
				Policies: []Policy{
					{Condition: Condition{"a": "a-value"}, Declarations: map[string]interface{}{"out": false}},
				},

				Assertions: []Policy{
					{Condition: Condition{"a": "a-value"}, Declarations: map[string]interface{}{"out": true}, Comment: "some-assertion"},
				},
			},
			AllowedKeys: []string{"a", "b"},
			Expected:    errors.New(`error in assertion[0], comment: "some-assertion", expected: map[out:true], actual: map[out:false]`),
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			actual := tc.Policy.Validate(tc.AllowedKeys)

			if !reflect.DeepEqual(tc.Expected, actual) {
				t.Errorf("Expected error: %v got %v", tc.Expected, actual)
			}
		})
	}
}

func TestPolicyList_Apply(t *testing.T) {
	cases := map[string]struct {
		Policy   PolicyList
		Truth    map[string]string
		Expected map[string]interface{}
	}{
		"cascading-overwrite": {
			Policy: PolicyList{
				Policies: []Policy{
					{Condition: Condition{}, Declarations: map[string]interface{}{"last-fired": 0}},
					{Condition: Condition{}, Declarations: map[string]interface{}{"last-fired": 1}},
				},
			},
			Truth: map[string]string{},
			Expected: map[string]interface{}{
				"last-fired": 1,
			},
		},
		"cascading-merge": {
			Policy: PolicyList{
				Policies: []Policy{
					{Condition: Condition{}, Declarations: map[string]interface{}{"first": 0}},
					{Condition: Condition{}, Declarations: map[string]interface{}{"second": 1}},
				},
			},
			Truth: map[string]string{},
			Expected: map[string]interface{}{
				"first":  0,
				"second": 1,
			},
		},
		"no-conditions-match": {
			Policy: PolicyList{
				Policies: []Policy{
					{Condition: Condition{"a": "true"}, Declarations: map[string]interface{}{"first": 0}},
					{Condition: Condition{"a": "true"}, Declarations: map[string]interface{}{"second": 1}},
				},
			},
			Truth:    map[string]string{},
			Expected: map[string]interface{}{},
		},
		"partial-conditions-match": {
			Policy: PolicyList{
				Policies: []Policy{
					{Condition: Condition{"a": "true"}, Declarations: map[string]interface{}{"last-fired": 0}},
					{Condition: Condition{"a": "false"}, Declarations: map[string]interface{}{"last-fired": 1}},
				},
			},
			Truth: map[string]string{"a": "true"},
			Expected: map[string]interface{}{
				"last-fired": 0,
			},
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			actual := tc.Policy.Apply(tc.Truth)

			if !reflect.DeepEqual(tc.Expected, actual) {
				t.Errorf("Expected: %v got %v", tc.Expected, actual)
			}
		})
	}
}

func TestNewPolicyListFromJson(t *testing.T) {
	cases := map[string]struct {
		Json        string
		AllowedKeys []string
		Expected    error
	}{
		"invalid-json": {
			Json:        `invalid-json`,
			AllowedKeys: []string{},
			Expected:    errors.New("couldn't decode PolicyList from JSON: invalid character 'i' looking for beginning of value"),
		},
		"unknown-field": {
			Json:        `{"unknown-field":[]}`,
			AllowedKeys: []string{},
			Expected:    errors.New(`couldn't decode PolicyList from JSON: json: unknown field "unknown-field"`),
		},
		"bad-key": {
			Json: `{"policy":[
				{"//":"user-comment", "if":{"unknown-condition":""}}
			]}`,
			AllowedKeys: []string{},
			Expected:    errors.New(`error in policy[0], comment: "user-comment", error: unknown condition keys: [unknown-condition] condition keys must one of: [], check their capitalization and spelling`),
		},
		"bad-assertion": {
			Json: `{
			"policy":[
				{"if":{}, "then":{"foo":"bar"}},
				{"if":{}, "then":{"foo":"bazz"}}
			],
			"assert":[{"//":"check bad-value", "if":{}, "then":{"foo":"bad-value"}}]
			}`,
			AllowedKeys: []string{},
			Expected:    errors.New(`error in assertion[0], comment: "check bad-value", expected: map[foo:bad-value], actual: map[foo:bazz]`),
		},
		"good-fizzbuzz": {
			Json: `{
			"policy": [
				{"if": {}, "then": {"print":"{{number}}"}},
				{"if": {"multiple-of-3":"true"}, "then": {"print":"fizz"}},
				{"if": {"multiple-of-5":"true"}, "then": {"print":"buzz"}},
				{"if": {"multiple-of-3":"true", "multiple-of-5":"true"}, "then": {"print":"fizzbuzz"}}
			],
			"assert": [{"if":{"multiple-of-3":"true", "multiple-of-5":"true"}, "then":{"print":"fizzbuzz"}}]
			}`,
			AllowedKeys: []string{"multiple-of-3", "multiple-of-5"},
			Expected:    nil,
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			pl, err := NewPolicyListFromJson([]byte(tc.Json), tc.AllowedKeys)
			if pl == nil && err == nil || pl != nil && err != nil {
				t.Fatalf("Expected exactly one of PolicyList and err to be nil PolicyList: %v, Error: %v", pl, err)
			}

			if !reflect.DeepEqual(err, tc.Expected) {
				t.Errorf("Expected error: %v got: %v", tc.Expected, err)
			}
		})
	}
}
