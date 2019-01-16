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

package toggles

import (
	"fmt"

	"github.com/spf13/viper"
)

func ExampleToggle_EnvironmentVariable() {
	ts := NewToggleSet("foo.")
	toggle := ts.Toggle("bar", true, "bar gets a default of true")

	fmt.Println(toggle.EnvironmentVariable())

	// Output: GSB_FOO_BAR
}

func ExampleToggle_IsActive() {
	ts := NewToggleSet("foo.")
	toggle := ts.Toggle("bar", true, "bar gets a default of true")

	fmt.Println(toggle.IsActive())
	viper.Set("foo.bar", "false")
	defer viper.Reset()

	fmt.Println(toggle.IsActive())

	// Output: true
	// false
}

func ExampleToggleSet_Toggles() {
	ts := NewToggleSet("foo.")

	//  add some toggles
	ts.Toggle("z", true, "a toggle")
	ts.Toggle("a", false, "another toggle")
	ts.Toggle("b", true, "a third toggle")

	for _, tgl := range ts.Toggles() {
		fmt.Printf("name: %s, var: %s, description: %q, default: %v\n", tgl.Name, tgl.EnvironmentVariable(), tgl.Description, tgl.Default)
	}

	// Output: name: a, var: GSB_FOO_A, description: "another toggle", default: false
	// name: b, var: GSB_FOO_B, description: "a third toggle", default: true
	// name: z, var: GSB_FOO_Z, description: "a toggle", default: true
}
