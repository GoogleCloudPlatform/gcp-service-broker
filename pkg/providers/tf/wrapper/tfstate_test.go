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

import "fmt"

func ExampleNewTfstate_Good() {
	state := `{
    "version": 3,
    "terraform_version": "0.11.10",
    "serial": 2,
    "modules": [
        {
            "path": ["root"],
            "outputs": {},
            "resources": {},
            "depends_on": []
        },
        {
            "path": ["root", "instance"],
            "outputs": {},
            "resources": {},
            "depends_on": []
        }
    ]
  }`

	_, err := NewTfstate([]byte(state))
	fmt.Printf("%v", err)

	// Output: <nil>
}

func ExampleNewTfstate_BadVersion() {
	state := `{
    "version": 4,
    "terraform_version": "0.11.10",
    "serial": 2,
    "modules": [
        {
            "path": ["root"],
            "outputs": {},
            "resources": {},
            "depends_on": []
        }
    ]
  }`

	_, err := NewTfstate([]byte(state))
	fmt.Printf("%v", err)

	// Output: unsupported tfstate version: 4
}

func ExampleTfstate_GetModule() {
	state := `{
    "version": 3,
    "terraform_version": "0.11.10",
    "serial": 2,
    "modules": [
        {
            "path": ["root", "instance"],
            "outputs": {
                "Name": {
                    "sensitive": false,
                    "type": "string",
                    "value": "cf-binding-ex351277"
                }
            },
            "resources": {},
            "depends_on": []
        }
    ]
  }`

	tfstate, _ := NewTfstate([]byte(state))
	fmt.Printf("%v\n", tfstate.GetModule("does-not-exist"))
	fmt.Printf("%v\n", tfstate.GetModule("root", "instance"))

	// Output: <nil>
	// [module: root/instance with 1 outputs]
}
