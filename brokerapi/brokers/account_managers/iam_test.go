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

package account_managers

import (
	"testing"

	"github.com/GoogleCloudPlatform/gcp-service-broker/utils"
	cloudresourcemanager "google.golang.org/api/cloudresourcemanager/v1"
)

func TestIamMergeBindings(t *testing.T) {
	table := map[string]struct {
		input  []*cloudresourcemanager.Binding
		expect map[string][]string
	}{
		"combine matching roles": {
			input: []*cloudresourcemanager.Binding{
				{Role: "role-1", Members: []string{"m1"}},
				{Role: "role-1", Members: []string{"m2"}},
			},
			expect: map[string][]string{
				"role-1": {"m1", "m2"},
			},
		},
		"deduplicates members": {
			input: []*cloudresourcemanager.Binding{
				{Role: "role-1", Members: []string{"m1"}},
				{Role: "role-1", Members: []string{"m1"}},
			},
			expect: map[string][]string{
				"role-1": {"m1"},
			},
		},
		"removes roles with no members": {
			input: []*cloudresourcemanager.Binding{
				{Role: "role-1", Members: []string{}},
			},
			expect: map[string][]string{},
		},
		"does not merge different roles": {
			input: []*cloudresourcemanager.Binding{
				{Role: "role-1", Members: []string{"m1", "m2"}},
				{Role: "role-2", Members: []string{"m1", "m3"}},
			},
			expect: map[string][]string{
				"role-1": {"m1", "m2"},
				"role-2": {"m1", "m3"},
			},
		},
	}

	for tn, tc := range table {
		merged := mergeBindings(tc.input)

		if len(merged) != len(tc.expect) {
			t.Errorf("%s) expected %d merged bindings, got: %d", tn, len(tc.expect), len(merged))
		}

		for _, binding := range merged {
			expset := utils.NewStringSet(tc.expect[binding.Role]...)
			gotset := utils.NewStringSet(binding.Members...)

			if !expset.Equals(gotset) {
				t.Errorf("%s) expected %v members in %s role, got %v", tn, expset.ToSlice(), binding.Role, gotset.ToSlice())
			}
		}
	}
}
