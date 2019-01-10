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

package policy

import (
	"fmt"
	"reflect"
)

// Conditions are a set of values that can be compared with a base truth and
// return true if all of the facets of the condition match.
//
// Blank facets are assumed to match everything.
type Condition struct {
	ServiceId   string `json:"service_id"`
	ServiceName string `json:"service_name"`
}

// AppliesTo returns true if all the facets of this condition match the given
// truth.
func (cond *Condition) AppliesTo(truth Condition) bool {
	if cond.ServiceId != "" && cond.ServiceId != truth.ServiceId {
		return false
	}

	if cond.ServiceName != "" && cond.ServiceName != truth.ServiceName {
		return false
	}

	return true
}

// Policy combines a condition with several sets of values that are set if
// the condition holds true.
type Policy struct {
	Comment   string    `json:"//"`
	Condition Condition `json:"if"`

	Actions map[string]interface{} `json:"then"`
}

type PolicyList struct {
	Policies   []Policy `json:"policy"`
	Assertions []Policy `json:"assert"`
}

// CheckAssertions tests each assertion in the Assertions list against the
// policies list. The conditoin is used as the ground truth and the
// actions are used as the expected output. If the actions don't match then
// an error is returned.
func (pl *PolicyList) CheckAssertions() error {
	for i, assertion := range pl.Assertions {
		expected := assertion.Actions
		actual := pl.Apply(assertion.Condition)

		if !reflect.DeepEqual(actual, expected) {
			return fmt.Errorf("Error in assertion %d, comment: %q, expected: %v, actual: %v", i+1, assertion.Comment, expected, actual)
		}
	}

	return nil
}

// Apply runs through the list of policies, first to last, and cascades the
// values of each if they match the given condition, returning the merged
// map at the end.
func (pl *PolicyList) Apply(groundTruth Condition) map[string]interface{} {
	out := make(map[string]interface{})

	for _, policy := range pl.Policies {
		if !policy.Condition.AppliesTo(groundTruth) {
			continue
		}

		for k, v := range policy.Actions {
			out[k] = v
		}
	}

	return out
}
