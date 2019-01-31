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
		MigrationJs:          `function format(input){return "[" + input + "]"}`,
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

func ExampleJsTransform_ToJs() {
	em := JsTransform{
		EnvironmentVariables: []string{"GSB_FOO", "GSB_BAR"},
		MigrationJs:          `function format(input){return "[" + input + "]"}`,
		TransformFuncName:    "format",
	}

	fmt.Println(em.ToJs())

	// Output: {
	// function format(input){return "[" + input + "]"}
	//
	// {
	//   // GSB_FOO
	//   var prop = '.properties.gsb_foo'
	//   var context = {}; // no additional context defined
	//   if (prop in properties['properties']) {
	//     properties['properties'][prop].value = format(properties['properties'][prop].value, context);
	//   }
	// }
	// {
	//   // GSB_BAR
	//   var prop = '.properties.gsb_bar'
	//   var context = {}; // no additional context defined
	//   if (prop in properties['properties']) {
	//     properties['properties'][prop].value = format(properties['properties'][prop].value, context);
	//   }
	// }
	// }
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
				MigrationJs:          `function set(input){return "new-value"}`,
				TransformFuncName:    "set",
			},
			Env:         map[string]string{"FOO": "old-value"},
			ExpectedEnv: map[string]string{"FOO": "new-value"},
		},
		"ignores unknown": {
			Migration: JsTransform{
				EnvironmentVariables: []string{"FOO"},
				MigrationJs:          `function set(input){return "new-value"}`,
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
			ExpectedErr: errors.New("processing FOO loading JS function: (anonymous): Line 1:5 Unexpected token function (and 3 more errors)"),
			ExpectedEnv: map[string]string{"FOO": "old-value"},
		},
		"bad func": {
			Migration: JsTransform{
				EnvironmentVariables: []string{"FOO"},
				MigrationJs:          `function set(input){return "new-value"}`,
				TransformFuncName:    "invalidFuncName",
			},
			Env:         map[string]string{"FOO": "old-value"},
			ExpectedErr: errors.New("processing FOO calling: invalidFuncName: ReferenceError: 'invalidFuncName' is not defined"),
			ExpectedEnv: map[string]string{"FOO": "old-value"},
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
