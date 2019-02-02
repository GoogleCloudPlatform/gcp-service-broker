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

package migration

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

func ExampleJsTransform_RunGo() {
	em := JsTransform{
		EnvironmentVariables: []string{"GSB_FOO", "GSB_BAR"},
		MigrationJs:          `function format(input){setProp(input, "[" + lookupProp(input) + "]")}`,
		TransformFuncName:    "format",
	}

	env := map[string]string{
		"GSB_FOO": "to,encapsulate",
	}

	err := em.RunGo(env)
	if err != nil {
		panic(err)
	}

	fmt.Println(env["GSB_FOO"])

	// Output: [to,encapsulate]
}

func TestJsTransform_RunGo(t *testing.T) {
	cases := map[string]struct {
		Migration   JsTransform
		Env         map[string]string
		ExpectedErr error
		ExpectedEnv map[string]string
	}{
		"set value": {
			Migration: JsTransform{
				EnvironmentVariables: []string{"FOO"},
				MigrationJs:          `function set(envVar){setProp(envVar, "new-value")}`,
				TransformFuncName:    "set",
			},
			Env:         map[string]string{"FOO": "old-value"},
			ExpectedEnv: map[string]string{"FOO": "new-value"},
		},
		"delete value": {
			Migration: JsTransform{
				EnvironmentVariables: []string{"FOO"},
				MigrationJs:          `function process(envVar){deleteProp(envVar)}`,
				TransformFuncName:    "process",
			},
			Env:         map[string]string{"FOO": "old-value"},
			ExpectedEnv: map[string]string{},
		},
		"ignores unknown": {
			Migration: JsTransform{
				EnvironmentVariables: []string{"FOO"},
				MigrationJs:          `function set(input){return (lookupProp(input)) ? "new-value": null; }`,
				TransformFuncName:    "set",
			},
			Env:         map[string]string{"BAR": "old-value"},
			ExpectedEnv: map[string]string{"BAR": "old-value"},
		},
		"bad js": {
			Migration: JsTransform{
				EnvironmentVariables: []string{"FOO"},
				MigrationJs:          `set function(input) return{return "new-value"}`,
				TransformFuncName:    "set",
			},
			Env:         map[string]string{"FOO": "old-value"},
			ExpectedErr: errors.New("loading JS function: (anonymous): Line 2:5 Unexpected token function (and 4 more errors)"),
			ExpectedEnv: map[string]string{"FOO": "old-value"},
		},
		"bad func": {
			Migration: JsTransform{
				EnvironmentVariables: []string{"FOO"},
				MigrationJs:          `function set(input){return "new-value"}`,
				TransformFuncName:    "invalidFuncName",
			},
			Env:         map[string]string{"FOO": "old-value"},
			ExpectedErr: errors.New("loading JS function: ReferenceError: 'invalidFuncName' is not defined"),
			ExpectedEnv: map[string]string{"FOO": "old-value"},
		},
		"access props": {
			Migration: JsTransform{
				EnvironmentVariables: []string{"FOO"},
				MigrationJs:          `function getRelated(envVar){setProp(envVar, lookupProp(envVar+"_PLUS"))}`,
				TransformFuncName:    "getRelated",
			},
			Env:         map[string]string{"FOO": "old-value", "FOO_PLUS": "new-vlaue"},
			ExpectedEnv: map[string]string{"FOO": "new-vlaue", "FOO_PLUS": "new-vlaue"},
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			err := tc.Migration.RunGo(tc.Env)
			if !reflect.DeepEqual(err, tc.ExpectedErr) {
				t.Errorf("Expected error: %v Got: %v", tc.ExpectedErr, err)
			}

			if !reflect.DeepEqual(tc.Env, tc.ExpectedEnv) {
				t.Errorf("Expected: %v Got: %v", tc.ExpectedEnv, tc.Env)
			}
		})
	}
}
