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

package interpolation

import (
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/hil"
	"github.com/spf13/cast"
)

func TestEval(t *testing.T) {
	tests := map[string]struct {
		Template      string
		Variables     map[string]interface{}
		Expected      interface{}
		ErrorContains string
	}{
		"Non-Templated String":  {Template: "foo", Expected: "foo"},
		"Basic Evaluation":      {Template: "${33}", Expected: "33"},
		"Escaped Evaluation":    {Template: "$${33}", Expected: "${33}"},
		"Missing Variable":      {Template: "${a}", ErrorContains: "unknown variable accessed: a"},
		"Variable Substitution": {Template: "${foo}", Variables: map[string]interface{}{"foo": 33}, Expected: "33"},
		"Bad Template":          {Template: "${", ErrorContains: "expected expression"},
		"Truncate Required":     {Template: `${str.truncate(2, "expression")}`, Expected: "ex"},
		"Truncate Not Required": {Template: `${str.truncate(200, "expression")}`, Expected: "expression"},
		"Counter":               {Template: "${counter.next()},${counter.next()},${counter.next()}", Expected: "1,2,3"},
		"Query Escape":          {Template: `${str.queryEscape("hello world")}`, Expected: "hello+world"},
		"Query Amp":             {Template: `${str.queryEscape("hello&world")}`, Expected: "hello%26world"},
		"Regex":                 {Template: `${regexp.matches("^(D|d)[0-9]+$", "d12345")}`, Expected: "true"},
		"Bad Regex":             {Template: `${regexp.matches("^($", "d12345")}`, ErrorContains: "error parsing regexp"},
		"Conditionals True":     {Template: `${true ? "foo" : "bar"}`, Expected: "foo"},
		"Conditionals False":    {Template: `${false ? "foo" : "bar"}`, Expected: "bar"},
		"No Short Circuit":      {Template: `${false ? counter.next() : counter.next()}`, Expected: "2"},
		"assert success":        {Template: `${assert(true, "nothing should happen")}`, Expected: "true"},
		"assert failure":        {Template: `${assert(false, "failure message")}`, ErrorContains: "failure message"},
		"assert message":        {Template: `${assert(false, "failure message ${1+1}")}`, ErrorContains: "failure message 2"},
		"json marshal":          {Template: "${json.marshal(mapval)}", Variables: map[string]interface{}{"mapval": map[string]string{"hello": "world"}}, Expected: `{"hello":"world"}`},
		"json marshal array":    {Template: "${json.marshal(list)}", Variables: map[string]interface{}{"list": []string{"a", "b", "c"}}, Expected: `["a","b","c"]`},
		"json marshal numeric":  {Template: "${json.marshal(42)}", Expected: `42`},
		"json marshal string":   {Template: `${json.marshal("str")}`, Expected: `"str"`},
		"json marshal true":     {Template: "${json.marshal(true)}", Expected: `true`},
		"json marshal false":    {Template: "${json.marshal(false)}", Expected: `false`},
		"map flatten blank":     {Template: `${map.flatten(":", ";", mapval)}`, Variables: map[string]interface{}{"mapval": map[string]string{}}, Expected: ``},
		"map flatten one":       {Template: `${map.flatten(":", ";", mapval)}`, Variables: map[string]interface{}{"mapval": map[string]string{"key1": "val1"}}, Expected: `key1:val1`},
		"map flatten":           {Template: `${map.flatten(":", ";", mapval)}`, Variables: map[string]interface{}{"mapval": map[string]string{"key1": "val1", "key2": "val2"}}, Expected: `key1:val1;key2:val2`},
	}

	for tn, tc := range tests {
		hilStandardLibrary = createStandardLibrary()

		t.Run(tn, func(t *testing.T) {
			res, err := Eval(tc.Template, tc.Variables)
			expectingErr := tc.ErrorContains != ""
			hasErr := err != nil
			if expectingErr != hasErr {
				t.Errorf("Expecting error? %v, got: %v", expectingErr, err)
			}

			if expectingErr && !strings.Contains(err.Error(), tc.ErrorContains) {
				t.Errorf("Expected error: %v to contain %q", err, tc.ErrorContains)
			}

			if !reflect.DeepEqual(tc.Expected, res) {
				t.Errorf("Expected result: %+v, got %+v", tc.Expected, res)
			}
		})
	}
}

func TestHilFuncTimeNano(t *testing.T) {
	before := time.Now().UnixNano()
	result, _ := Eval("${time.nano()}", nil)
	value := cast.ToInt64(result)
	after := time.Now().UnixNano()

	if before >= value || value >= after {
		t.Errorf("Expected %d < %d < %d", before, value, after)
	}
}

func TestHilFuncRandBase64(t *testing.T) {
	result, _ := Eval("${rand.base64(32)}", nil)
	length := len(result.(string))
	if length != 44 {
		t.Errorf("Expected length to be %d got %d", 44, length)
	}

	result, _ = Eval("${rand.base64(16)}", nil)
	length = len(result.(string))
	if length != 24 {
		t.Errorf("Expected length to be %d got %d", 44, length)
	}
}

func TestHilToInterface(t *testing.T) {
	// This function tests hilToInterface operates correctly with regards to
	// taking valid user inputs (i.e. only JSON values), converting them to HIL
	// values then converting them back.
	tests := map[string]struct {
		UserInput interface{}
		Expected  interface{}
	}{
		"string":      {UserInput: "foo", Expected: "foo"},
		"numeric":     {UserInput: 42, Expected: "42"},
		"bool-true":   {UserInput: true, Expected: "1"},
		"bool-false":  {UserInput: false, Expected: "0"},
		"str-array":   {UserInput: []interface{}{"a", "b"}, Expected: []interface{}{"a", "b"}},
		"mixed-array": {UserInput: []interface{}{"a", 2}, Expected: []interface{}{"a", "2"}},
		"object": {
			UserInput: map[string]interface{}{"s": "str", "n": 42.0, "a": []interface{}{"a", 2}},
			Expected:  map[string]interface{}{"s": "str", "n": "42", "a": []interface{}{"a", "2"}},
		},
	}

	for tn, tc := range tests {
		t.Run(tn, func(t *testing.T) {
			converted, err := hil.InterfaceToVariable(tc.UserInput)
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			res, err := hilToInterface(converted)
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			if !reflect.DeepEqual(tc.Expected, res) {
				t.Errorf("Expected result: %+v (%t), got %+v (%t)", tc.Expected, tc.Expected, res, res)
			}
		})
	}

}
