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
)

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
