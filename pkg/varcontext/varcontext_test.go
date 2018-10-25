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

package varcontext

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

func TestVarContext_GetString(t *testing.T) {
	// The following tests operate on the following example map
	testContext := map[string]interface{}{
		"anInt":   42,
		"aString": "value",
	}

	tests := map[string]struct {
		Key      string
		Expected string
		Error    string
	}{
		"int":         {"anInt", "42", ""},
		"string":      {"aString", "value", ""},
		"missing key": {"DNE", "", `missing value for key "DNE"`},
	}

	for tn, tc := range tests {
		t.Run(tn, func(t *testing.T) {
			vc := &VarContext{context: testContext}

			result := vc.GetString(tc.Key)
			if result != tc.Expected {
				t.Errorf("Expected to get: %q actual: %q", tc.Expected, result)
			}

			expectedErrors := tc.Error != ""
			hasError := vc.Error() != nil
			if hasError != expectedErrors {
				t.Error("Got error when not expecting or missing error that was expected")
			}

			if tc.Error != "" && !strings.Contains(vc.Error().Error(), tc.Error) {
				t.Errorf("Expected error to contain %q, but got: %v", tc.Error, vc.Error())
			}
		})
	}
}

func TestVarContext_GetInt(t *testing.T) {
	// The following tests operate on the following example map
	testContext := map[string]interface{}{
		"anInt":   42,
		"aString": "value",
	}

	tests := map[string]struct {
		Key      string
		Expected int
		Error    string
	}{
		"int":         {"anInt", 42, ""},
		"string":      {"aString", 0, `value for "aString" must be a integer`},
		"missing key": {"DNE", 0, `missing value for key "DNE"`},
	}

	for tn, tc := range tests {
		t.Run(tn, func(t *testing.T) {
			vc := &VarContext{context: testContext}

			result := vc.GetInt(tc.Key)
			if result != tc.Expected {
				t.Errorf("Expected to get: %q actual: %q", tc.Expected, result)
			}

			expectedErrors := tc.Error != ""
			hasError := vc.Error() != nil
			if hasError != expectedErrors {
				t.Errorf("Got error when not expecting or missing error that was expected: %v", vc.Error())
			}

			if tc.Error != "" && !strings.Contains(vc.Error().Error(), tc.Error) {
				t.Errorf("Expected error to contain %q, but got: %v", tc.Error, vc.Error())
			}
		})
	}
}

func TestVarContext_GetBool(t *testing.T) {
	// The following tests operate on the following example map
	testContext := map[string]interface{}{
		"anInt":   42,
		"zero":    0,
		"tsBool":  "true",
		"fsBool":  "false",
		"tBool":   true,
		"fBool":   false,
		"aString": "value",
	}

	tests := map[string]struct {
		Key      string
		Expected bool
		Error    string
	}{
		"true bool":         {"tBool", true, ""},
		"false bool":        {"fBool", false, ""},
		"true string bool":  {"tsBool", true, ""},
		"false string bool": {"fsBool", false, ""},
		"int":               {"anInt", true, ""},
		"zero":              {"zero", false, ""},
		"string":            {"aString", false, `value for "aString" must be a boolean`},
		"missing key":       {"DNE", false, `missing value for key "DNE"`},
	}

	for tn, tc := range tests {
		t.Run(tn, func(t *testing.T) {
			vc := &VarContext{context: testContext}

			result := vc.GetBool(tc.Key)
			if result != tc.Expected {
				t.Errorf("Expected to get: %v actual: %v", tc.Expected, result)
			}

			expectedErrors := tc.Error != ""
			hasError := vc.Error() != nil
			if hasError != expectedErrors {
				t.Errorf("Got error when not expecting or missing error that was expected: %v", vc.Error())
			}

			if tc.Error != "" && !strings.Contains(vc.Error().Error(), tc.Error) {
				t.Errorf("Expected error to contain %q, but got: %v", tc.Error, vc.Error())
			}
		})
	}
}

func TestVarContext_GetStringMapString(t *testing.T) {
	// The following tests operate on the following example map
	testContext := map[string]interface{}{
		"single":  map[string]string{"foo": "bar"},
		"aString": "value",
		"json":    `{"foo":"bar"}`,
	}

	tests := map[string]struct {
		Key      string
		Expected map[string]string
		Error    string
	}{
		"single map": {"single", map[string]string{"foo": "bar"}, ""},
		"json map":   {"json", map[string]string{"foo": "bar"}, ""},
		"string":     {"aString", map[string]string{}, `value for "aString" must be a map[string]string`},
	}

	for tn, tc := range tests {
		t.Run(tn, func(t *testing.T) {
			vc := &VarContext{context: testContext}

			result := vc.GetStringMapString(tc.Key)
			if !reflect.DeepEqual(result, tc.Expected) {
				t.Errorf("Expected to get: %v actual: %v", tc.Expected, result)
			}

			expectedErrors := tc.Error != ""
			hasError := vc.Error() != nil
			if hasError != expectedErrors {
				t.Fatalf("Got error when not expecting or missing error that was expected: %v", vc.Error())
			}

			if tc.Error != "" && !strings.Contains(vc.Error().Error(), tc.Error) {
				t.Errorf("Expected error to contain %q, but got: %v", tc.Error, vc.Error())
			}
		})
	}
}

func TestVarContext_ToJson(t *testing.T) {
	vc := &VarContext{context: map[string]interface{}{
		"t": true,
		"f": false,
		"s": "a string",
		"a": []interface{}{"an", "array"},
		"F": 123.45,
	}}
	expected := vc.ToMap()

	serialized, _ := vc.ToJson()
	actual := make(map[string]interface{})
	json.Unmarshal(serialized, &actual)

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Expected: %#v, Got: %#v", expected, actual)
	}
}
