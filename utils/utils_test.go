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

package utils

import "fmt"

func ExampleSetParameter() {
	// Creates an object if none is input
	out, err := SetParameter(nil, "foo", 42)
	fmt.Printf("%s, %v\n", string(out), err)

	// Replaces existing values
	out, err = SetParameter([]byte(`{"replace": "old"}`), "replace", "new")
	fmt.Printf("%s, %v\n", string(out), err)

	// Output: {"foo":42}, <nil>
	// {"replace":"new"}, <nil>
}
