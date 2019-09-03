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
	"github.com/GoogleCloudPlatform/gcp-service-broker/utils"
	cloudresourcemanager "google.golang.org/api/cloudresourcemanager/v1"
)

// Merge multiple Bindings such that Bindings with the same Role result in
// a single Binding with combined Members
func mergeBindings(bindings []*cloudresourcemanager.Binding) []*cloudresourcemanager.Binding {
	rb := []*cloudresourcemanager.Binding{}

	for role, members := range rolesToMembersMap(bindings) {
		if members.IsEmpty() {
			continue
		}

		rb = append(rb, &cloudresourcemanager.Binding{
			Role:    role,
			Members: members.ToSlice(),
		})
	}

	return rb
}

// Map a role to a map of members, allowing easy merging of multiple bindings.
func rolesToMembersMap(bindings []*cloudresourcemanager.Binding) map[string]utils.StringSet {
	bm := make(map[string]utils.StringSet)
	for _, b := range bindings {
		if set, ok := bm[b.Role]; ok {
			set.Add(b.Members...)
		} else {
			bm[b.Role] = utils.NewStringSet(b.Members...)
		}
	}

	return bm
}

// Remove a member from all role bindings to clean up the iam policy before deleting a service account.
func removeMemberFromBindings(bindings []*cloudresourcemanager.Binding, memberToRemove string) []*cloudresourcemanager.Binding {
	for _, binding := range bindings {
		for i, member := range binding.Members {
			if member == memberToRemove {
				binding.Members = append(binding.Members[:i], binding.Members[i+1:]...)
				break
			}
		}
	}

	return bindings
}
