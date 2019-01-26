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
	"reflect"
	"testing"
)

func TestDiffStringMap(t *testing.T) {
	cases := map[string]struct {
		Old      map[string]string
		New      map[string]string
		Expected map[string]Diff
	}{
		"same": {
			Old:      map[string]string{"a": "b"},
			New:      map[string]string{"a": "b"},
			Expected: map[string]Diff{},
		},
		"removed": {
			Old:      map[string]string{"a": "old"},
			New:      map[string]string{},
			Expected: map[string]Diff{"a": Diff{Old: "old", New: ""}},
		},
		"added": {
			Old:      map[string]string{},
			New:      map[string]string{"a": "new"},
			Expected: map[string]Diff{"a": Diff{Old: "", New: "new"}},
		},
		"changed": {
			Old:      map[string]string{"a": "old"},
			New:      map[string]string{"a": "new"},
			Expected: map[string]Diff{"a": Diff{Old: "old", New: "new"}},
		},
		"full-gambit": {
			Old: map[string]string{
				"removed": "removed",
				"changed": "orig",
			},
			New: map[string]string{
				"changed": "new",
				"added":   "added",
			},
			Expected: map[string]Diff{
				"removed": Diff{Old: "removed", New: ""},
				"changed": Diff{Old: "orig", New: "new"},
				"added":   Diff{Old: "", New: "added"},
			},
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			actual := DiffStringMap(tc.Old, tc.New)

			if !reflect.DeepEqual(tc.Expected, actual) {
				t.Errorf("Expected: %#v Actual: %#v", tc.Expected, actual)
			}
		})
	}
}
