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

import "reflect"

// NewStringSet creates a new string set with the given contents.
func NewStringSet(contents ...string) StringSet {
	set := StringSet{}
	set.Add(contents...)
	return set
}

// StringSet is a set data structure for strings
type StringSet map[string]bool

// Add puts one or more elements into the set.
func (set StringSet) Add(str ...string) {
	for _, s := range str {
		set[s] = true
	}
}

// ToSlice converts the set to a slice with undefined contents order.
func (set StringSet) ToSlice() []string {
	out := []string{}
	for k, _ := range set {
		out = append(out, k)
	}

	return out
}

// IsEmpty determines if the set has zero elements.
func (set StringSet) IsEmpty() bool {
	return len(set) == 0
}

// Equals compares the contents of the two sets and returns true if they are
// the same.
func (set StringSet) Equals(other StringSet) bool {
	return reflect.DeepEqual(set, other)
}

// Contains performs a set membership check.
func (set StringSet) Contains(other string) bool {
	_, ok := set[other]
	return ok
}
