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

package wrapper

import (
	"fmt"
	"strings"
	"testing"
)

func ExampleTerraformModule_Inputs() {
	module := TerraformModule{
		Name: "cloud_storage",
		Definition: `
    variable name {type = "string"}
    variable storage_class {type = "string"}

    resource "google_storage_bucket" "bucket" {
      name     = "${var.name}"
      storage_class = "${var.storage_class}"
    }
`,
	}

	inputs, err := module.Inputs()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v\n", inputs)

	// Output: [name storage_class]
}

func ExampleTerraformModule_Outputs() {
	module := TerraformModule{
		Name: "cloud_storage",
		Definition: `
    resource "google_storage_bucket" "bucket" {
      name     = "my-bucket"
      storage_class = "STANDARD"
    }

    output id {value = "${google_storage_bucket.bucket.id}"}
    output bucket_name {value = "my-bucket"}
`,
	}

	outputs, err := module.Outputs()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v\n", outputs)

	// Output: [bucket_name id]
}

func TestTerraformModule_Validate(t *testing.T) {
	cases := map[string]struct {
		Module      TerraformModule
		ErrContains string
	}{
		"nominal": {
			Module: TerraformModule{
				Name: "my_module",
				Definition: `
          resource "google_storage_bucket" "bucket" {
            name     = "my-bucket"
            storage_class = "STANDARD"
          }`,
			},
			ErrContains: "",
		},
		"bad-name": {
			Module: TerraformModule{
				Name: "my module",
				Definition: `
          resource "google_storage_bucket" "bucket" {
            name     = "my-bucket"
            storage_class = "STANDARD"
          }`,
			},
			ErrContains: "Field validation for 'Name' failed ",
		},
		"bad-hcl": {
			Module: TerraformModule{
				Name: "my_module",
				Definition: `
          resource "bucket" {
            name     = "my-bucket"`,
			},
			ErrContains: "Field validation for 'Definition' failed ",
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			err := tc.Module.Validate()
			if tc.ErrContains == "" {
				if err != nil {
					t.Fatalf("Expected no error but got: %v", err)
				}
			} else {
				if err == nil {
					t.Fatalf("Expected error containing %q but got nil", tc.ErrContains)
				}
				if !strings.Contains(err.Error(), tc.ErrContains) {
					t.Fatalf("Expected error containing %q but got %v", tc.ErrContains, err)
				}
			}
		})
	}
}
