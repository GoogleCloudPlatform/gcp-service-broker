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
			if vc.HasErrors() != expectedErrors {
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
			if vc.HasErrors() != expectedErrors {
				t.Errorf("Got error when not expecting or missing error that was expected: %v", vc.Error())
			}

			if tc.Error != "" && !strings.Contains(vc.Error().Error(), tc.Error) {
				t.Errorf("Expected error to contain %q, but got: %v", tc.Error, vc.Error())
			}
		})
	}
}
