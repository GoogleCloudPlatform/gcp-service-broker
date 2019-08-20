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

package validation

import (
	"encoding/json"
	"fmt"
	"testing"
)

type osbNameStruct struct {
	Name string `validate:"osbname"`
}

type jsonStruct struct {
	Json []byte `validate:"json"`
}

type jsonStringStruct struct {
	Json string `validate:"json"`
}

type hclStruct struct {
	Hcl string `validate:"hcl"`
}

type terraformIdentifierStruct struct {
	Id string `validate:"terraform_identifier"`
}

type jsonschemaTypeStruct struct {
	Type string `validate:"jsonschema_type"`
}

func TestValidateStruct(t *testing.T) {
	cases := map[string]struct {
		Validate  interface{}
		ExpectErr bool
	}{
		"osbname missing": {
			Validate:  osbNameStruct{},
			ExpectErr: true,
		},
		"osbname valid": {
			Validate:  osbNameStruct{Name: "google-storage"},
			ExpectErr: false,
		},
		"osbname spaces": {
			Validate:  osbNameStruct{Name: " google-storage  "},
			ExpectErr: true,
		},
		"osbname dots": {
			Validate:  osbNameStruct{Name: "google.storage"},
			ExpectErr: false,
		},
		"osbname alpha": {
			Validate:  osbNameStruct{Name: "googlestorage"},
			ExpectErr: false,
		},
		"osbname upper": {
			Validate:  osbNameStruct{Name: "GOOGLESTORAGE"},
			ExpectErr: false,
		},
		"osbname numeric": {
			Validate:  osbNameStruct{Name: "12345"},
			ExpectErr: false,
		},
		"json bytes blank": {
			Validate:  jsonStruct{},
			ExpectErr: true,
		},
		"json bytes empty object": {
			Validate:  jsonStruct{Json: []byte("{}")},
			ExpectErr: false,
		},
		"json bytes bad object": {
			Validate:  jsonStruct{Json: []byte("")},
			ExpectErr: true,
		},
		"json bytes full object": {
			Validate:  jsonStruct{Json: []byte(`{"a":42, "s":"foo"}`)},
			ExpectErr: false,
		},
		"json string blank": {
			Validate:  jsonStringStruct{},
			ExpectErr: true,
		},
		"json string empty object": {
			Validate:  jsonStringStruct{Json: "{}"},
			ExpectErr: false,
		},
		"json string bad object": {
			Validate:  jsonStringStruct{Json: ""},
			ExpectErr: true,
		},
		"json string full object": {
			Validate:  jsonStringStruct{Json: `{"a":42, "s":"foo"}`},
			ExpectErr: false,
		},
		"hcl blank": {
			Validate:  hclStruct{Hcl: ""},
			ExpectErr: false,
		},
		"hcl bad": {
			Validate:  hclStruct{Hcl: "asfd"},
			ExpectErr: true,
		},
		"hcl json": {
			Validate:  hclStruct{Hcl: `{"a":42, "s":"foo"}`},
			ExpectErr: false,
		},
		"hcl terraform provider": {
			Validate: hclStruct{Hcl: `
				provider "google" {
				  credentials = "${file("account.json")}"
				  project     = "my-project-id"
				  region      = "us-central1"
				}
				provider "google-beta" {
				  credentials = "${file("account.json")}"
				  project     = "my-project-id"
				  region      = "us-central1"
				}
				`},
			ExpectErr: false,
		},
		"terraform identifier good": {
			Validate:  terraformIdentifierStruct{Id: "good_value_here"},
			ExpectErr: false,
		},
		"terraform identifier bad": {
			Validate:  terraformIdentifierStruct{Id: "bad.value.here"},
			ExpectErr: true,
		},
		"terraform identifier blank": {
			Validate:  terraformIdentifierStruct{Id: ""},
			ExpectErr: false,
		},
		"jsonschema type blank": {
			Validate:  jsonschemaTypeStruct{Type: ""},
			ExpectErr: false,
		},
		"jsonschema type object": {
			Validate:  jsonschemaTypeStruct{Type: "object"},
			ExpectErr: false,
		},
		"jsonschema type boolean": {
			Validate:  jsonschemaTypeStruct{Type: "boolean"},
			ExpectErr: false,
		},
		"jsonschema type number": {
			Validate:  jsonschemaTypeStruct{Type: "number"},
			ExpectErr: false,
		},
		"jsonschema type string": {
			Validate:  jsonschemaTypeStruct{Type: "string"},
			ExpectErr: false,
		},
		"jsonschema type integer": {
			Validate:  jsonschemaTypeStruct{Type: "integer"},
			ExpectErr: false,
		},
		"jsonschema type invalid": {
			Validate:  jsonschemaTypeStruct{Type: "invalid"},
			ExpectErr: true,
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			err := ValidateStruct(tc.Validate)
			gotErr := err != nil
			if gotErr != tc.ExpectErr {
				t.Errorf("expected error? %v got %v", tc.ExpectErr, err)
			}
		})
	}
}

func ExampleErrIfNotOSBName() {
	fmt.Println("Good is nil:", ErrIfNotOSBName("google-storage", "my-field") == nil)
	fmt.Println("Bad:", ErrIfNotOSBName("google storage", "my-field"))

	// Output: Good is nil: true
	// Bad: field must match '^[a-zA-Z0-9-\.]+$': my-field
}

func ExampleErrIfNotJSONSchemaType() {
	fmt.Println("Good is nil:", ErrIfNotJSONSchemaType("string", "my-field") == nil)
	fmt.Println("Bad:", ErrIfNotJSONSchemaType("str", "my-field"))

	// Output: Good is nil: true
	// Bad: field must match '^(|object|boolean|array|number|string|integer)$': my-field
}

func ExampleErrIfNotHCL() {
	fmt.Println("Good HCL is nil:", ErrIfNotHCL(`provider "google" {
		credentials = "${file("account.json")}"
		project     = "my-project-id"
		region      = "us-central1"
	}`, "my-field") == nil)

	fmt.Println("Good JSON is nil:", ErrIfNotHCL(`{"a":42, "s":"foo"}`, "my-field") == nil)

	fmt.Println("Bad:", ErrIfNotHCL("google storage", "my-field"))

	// Output: Good HCL is nil: true
	// Good JSON is nil: true
	// Bad: invalid HCL: my-field
}

func ExampleErrIfNotTerraformIdentifier() {
	fmt.Println("Good is nil:", ErrIfNotTerraformIdentifier("good_id", "my-field") == nil)
	fmt.Println("Bad:", ErrIfNotTerraformIdentifier("bad id", "my-field"))

	// Output: Good is nil: true
	// Bad: field must match '^[a-z_]*$': my-field
}

func ExampleErrIfNotJSON() {
	fmt.Println("Good is nil:", ErrIfNotJSON(json.RawMessage("{}"), "my-field") == nil)
	fmt.Println("Bad:", ErrIfNotJSON(json.RawMessage(""), "my-field"))

	// Output: Good is nil: true
	// Bad: invalid JSON: my-field
}
