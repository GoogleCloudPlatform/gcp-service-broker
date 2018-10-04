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
				Enum:    map[interface{}]string{"a": "description", "b": "description"},
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
