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

import (
	"fmt"
	"os"
)

func ExamplePropertyToEnv() {
	env := PropertyToEnv("my.property.key-value")
	fmt.Println(env)

	// Output: GSB_MY_PROPERTY_KEY_VALUE
}

func ExamplePropertyToEnvUnprefixed() {
	env := PropertyToEnvUnprefixed("my.property.key-value")
	fmt.Println(env)

	// Output: MY_PROPERTY_KEY_VALUE
}

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

func ExampleUnmarshalObjectRemainder() {
	var obj struct {
		A string `json:"a_str"`
		B int
	}

	remainder, err := UnmarshalObjectRemainder([]byte(`{"a_str":"hello", "B": 33, "C": 123}`), &obj)
	fmt.Printf("%s, %v\n", string(remainder), err)

	remainder, err = UnmarshalObjectRemainder([]byte(`{"a_str":"hello", "B": 33}`), &obj)
	fmt.Printf("%s, %v\n", string(remainder), err)

	// Output: {"C":123}, <nil>
	// {}, <nil>
}

func ExampleGetDefaultProjectId() {
	os.Setenv("ROOT_SERVICE_ACCOUNT_JSON", `{"project_id": "my-project-123"}`)
	defer os.Unsetenv("ROOT_SERVICE_ACCOUNT_JSON")

	projectId, err := GetDefaultProjectId()
	fmt.Printf("%s, %v\n", projectId, err)

	// Output: my-project-123, <nil>
}
