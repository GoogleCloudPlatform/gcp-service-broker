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
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/validation"
)

func TestContextBuilder(t *testing.T) {
	cases := map[string]struct {
		Builder     *ContextBuilder
		Expected    map[string]interface{}
		ErrContains string
	}{
		"an empty context": {
			Builder:     Builder(),
			Expected:    map[string]interface{}{},
			ErrContains: "",
		},

		// MergeMap
		"MergeMap blank okay": {
			Builder:  Builder().MergeMap(map[string]interface{}{}),
			Expected: map[string]interface{}{},
		},
		"MergeMap multi-key": {
			Builder:  Builder().MergeMap(map[string]interface{}{"a": "a", "b": "b"}),
			Expected: map[string]interface{}{"a": "a", "b": "b"},
		},
		"MergeMap overwrite": {
			Builder:  Builder().MergeMap(map[string]interface{}{"a": "a"}).MergeMap(map[string]interface{}{"a": "aaa"}),
			Expected: map[string]interface{}{"a": "aaa"},
		},

		// MergeDefaults
		"MergeDefaults no defaults": {
			Builder:  Builder().MergeDefaults([]DefaultVariable{{Name: "foo"}}),
			Expected: map[string]interface{}{},
		},
		"MergeDefaults non-string": {
			Builder:  Builder().MergeDefaults([]DefaultVariable{{Name: "h2g2", Default: 42}}),
			Expected: map[string]interface{}{"h2g2": 42},
		},
		"MergeDefaults basic-string": {
			Builder:  Builder().MergeDefaults([]DefaultVariable{{Name: "a", Default: "no-template"}}),
			Expected: map[string]interface{}{"a": "no-template"},
		},
		"MergeDefaults template string": {
			Builder:  Builder().MergeDefaults([]DefaultVariable{{Name: "a", Default: "a"}, {Name: "b", Default: "${a}"}}),
			Expected: map[string]interface{}{"a": "a", "b": "a"},
		},
		"MergeDefaults no-overwrite": {
			Builder:  Builder().MergeDefaults([]DefaultVariable{{Name: "a", Default: "a"}, {Name: "a", Default: "b", Overwrite: false}}),
			Expected: map[string]interface{}{"a": "a"},
		},
		"MergeDefaults overwrite": {
			Builder:  Builder().MergeDefaults([]DefaultVariable{{Name: "a", Default: "a"}, {Name: "a", Default: "b", Overwrite: true}}),
			Expected: map[string]interface{}{"a": "b"},
		},

		"MergeDefaults object": {
			Builder:  Builder().MergeDefaults([]DefaultVariable{{Name: "o", Default: `{"foo": "bar"}`, Type: "object"}}),
			Expected: map[string]interface{}{"o": map[string]interface{}{"foo": "bar"}},
		},

		"MergeDefaults boolean": {
			Builder:  Builder().MergeDefaults([]DefaultVariable{{Name: "b", Default: `true`, Type: "boolean"}}),
			Expected: map[string]interface{}{"b": true},
		},
		"MergeDefaults array": {
			Builder:  Builder().MergeDefaults([]DefaultVariable{{Name: "a", Default: `["a","b","c","d"]`, Type: "array"}}),
			Expected: map[string]interface{}{"a": []interface{}{"a", "b", "c", "d"}},
		},
		"MergeDefaults number": {
			Builder:  Builder().MergeDefaults([]DefaultVariable{{Name: "n", Default: `1.234`, Type: "number"}}),
			Expected: map[string]interface{}{"n": 1.234},
		},
		"MergeDefaults integer": {
			Builder:  Builder().MergeDefaults([]DefaultVariable{{Name: "i", Default: `1234`, Type: "integer"}}),
			Expected: map[string]interface{}{"i": 1234},
		},
		"MergeDefaults string": {
			Builder:  Builder().MergeDefaults([]DefaultVariable{{Name: "s", Default: `1234`, Type: "string"}}),
			Expected: map[string]interface{}{"s": "1234"},
		},
		"MergeDefaults blank type": {
			Builder:  Builder().MergeDefaults([]DefaultVariable{{Name: "s", Default: `1234`, Type: ""}}),
			Expected: map[string]interface{}{"s": "1234"},
		},
		"MergeDefaults bad type": {
			Builder:     Builder().MergeDefaults([]DefaultVariable{{Name: "s", Default: `1234`, Type: "class"}}),
			ErrContains: "couldn't cast 1234 to class, unknown type",
		},

		// MergeEvalResult
		"MergeEvalResult accumulates context": {
			Builder:  Builder().MergeEvalResult("a", "a", "string").MergeEvalResult("b", "${a}", "string"),
			Expected: map[string]interface{}{"a": "a", "b": "a"},
		},
		"MergeEvalResult errors": {
			Builder:     Builder().MergeEvalResult("a", "${dne}", "string"),
			ErrContains: `couldn't compute the value for "a"`,
		},

		// MergeJsonObject
		"MergeJsonObject blank message": {
			Builder:  Builder().MergeJsonObject(json.RawMessage{}),
			Expected: map[string]interface{}{},
		},
		"MergeJsonObject valid message": {
			Builder:  Builder().MergeJsonObject(json.RawMessage(`{"a":"a"}`)),
			Expected: map[string]interface{}{"a": "a"},
		},
		"MergeJsonObject invalid message": {
			Builder:     Builder().MergeJsonObject(json.RawMessage(`{{{}}}`)),
			ErrContains: "invalid character '{'",
		},

		// MergeStruct
		"MergeStruct without JSON Tags": {
			Builder:  Builder().MergeStruct(struct{ Name string }{Name: "Foo"}),
			Expected: map[string]interface{}{"Name": "Foo"},
		},
		"MergeStruct with JSON Tags": {
			Builder: Builder().MergeStruct(struct {
				Name string `json:"username"`
			}{Name: "Foo"}),
			Expected: map[string]interface{}{"username": "Foo"},
		},

		// constants
		"Basic constants": {
			Builder: Builder().
				SetEvalConstants(map[string]interface{}{"PI": 3.14}).
				MergeEvalResult("out", "${PI}", "string"),
			Expected: map[string]interface{}{"out": "3.14"},
		},
		"User overrides constant": {
			Builder: Builder().
				SetEvalConstants(map[string]interface{}{"PI": 3.14}).
				MergeMap(map[string]interface{}{"PI": 3.2}). // reassign incorrectly, https://en.wikipedia.org/wiki/Indiana_Pi_Bill
				MergeEvalResult("PI", "${PI}", "string"),    // test which PI gets referenced
			Expected: map[string]interface{}{"PI": "3.14"},
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {

			vc, err := tc.Builder.Build()

			switch {
			case err == nil && tc.ErrContains == "":
				break
			case err == nil && tc.ErrContains != "":
				t.Fatalf("Got no error when %q was expected", tc.ErrContains)
			case err != nil && tc.ErrContains == "":
				t.Fatalf("Got error %v when none was expected", err)
			case !strings.Contains(err.Error(), tc.ErrContains):
				t.Fatalf("Got error %v, but expected it to contain %q", err, tc.ErrContains)
			}
			if vc == nil && tc.Expected != nil {
				t.Fatalf("Expected: %v, got: %v", tc.Expected, vc)
			}

			if vc != nil && !reflect.DeepEqual(vc.ToMap(), tc.Expected) {
				t.Errorf("Expected: %#v, got: %#v", tc.Expected, vc.ToMap())
			}

		})
	}
}

func ExampleContextBuilder_BuildMap() {
	_, e := Builder().MergeEvalResult("a", "${assert(false, \"failure!\")}", "string").BuildMap()
	fmt.Printf("Error: %v\n", e)

	m, _ := Builder().MergeEvalResult("a", "${1+1}", "string").BuildMap()
	fmt.Printf("Map: %v\n", m)

	//Output: Error: 1 error(s) occurred: couldn't compute the value for "a", template: "${assert(false, \"failure!\")}", assert: Assertion failed: failure!
	// Map: map[a:2]
}

func TestDefaultVariable_Validate(t *testing.T) {
	cases := map[string]validation.ValidatableTest{
		"empty": validation.ValidatableTest{
			Object: &DefaultVariable{},
			Expect: errors.New("missing field(s): default, name"),
		},
		"bad type": validation.ValidatableTest{
			Object: &DefaultVariable{
				Name:    "my-name",
				Default: 123,
				Type:    "stringss",
			},
			Expect: errors.New("field must match '^(|object|boolean|array|number|string|integer)$': type"),
		},
		"good": validation.ValidatableTest{
			Object: &DefaultVariable{
				Name:    "my-name",
				Default: 123,
				Type:    "string",
			},
			Expect: nil,
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			tc.Assert(t)
		})
	}
}
