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

/*
  Package policy defines a way to create simple cascading rule systems similar
	to CSS.

	A policy is broken into two parts, conditions and declarations. Conditions
	are the test that is run when it'd determined if a rule should fire.
	Declarations are the values that are set by the rule.

	Rules are executed in a low to high precidence order, and values are merged
	with the values from higher precidence rules overwriting values with the same
	keys that were set earlier.

	Rules systems can be painfully difficult to debug and test. This is especially
	true because they're often built as a safe way for non-programmers to modify
	business process and don't have any way to programatically assert their logic
	without resorting to hand-testing.

	To combat this issue, this rules system introduces three separate concepts.

	1. Conditions are all a strict equality check.
	2. Rules are executed from top-to-bottom, eliminating the need for complex
	   state analysis and backtracking algorithms.
	3. There is a built-in system for assertion checking that's exposed to the
	   rule authors.
*/

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/validation"
	"github.com/GoogleCloudPlatform/gcp-service-broker/utils"
)

// Conditions are a set of values that can be compared with a base truth and
// return true if all of the facets of the condition match.
type Condition map[string]string

// AppliesTo returns true if all the facets of this condition match the given
// truth.
func (cond Condition) AppliesTo(truth Condition) bool {
	for k, v := range cond {
		truthValue, ok := truth[k]
		if !ok || v != truthValue {
			return false
		}
	}

	return true
}

// ValidateKeys ensures all of the keys of the condition exist in the set of
// allowed keys.
func (cond Condition) ValidateKeys(allowedKeys []string) error {
	allowedSet := utils.NewStringSet(allowedKeys...)
	condKeys := utils.NewStringSetFromStringMapKeys(cond)

	invalidKeys := condKeys.Minus(allowedSet)

	if invalidKeys.IsEmpty() {
		return nil
	}

	return fmt.Errorf("unknown condition keys: %v condition keys must one of: %v, check their capitalization and spelling", invalidKeys, allowedKeys)
}

// Policy combines a condition with several sets of values that are set if
// the condition holds true.
type Policy struct {
	Comment   string    `json:"//"`
	Condition Condition `json:"if"`

	Declarations map[string]interface{} `json:"then"`
}

// PolicyList contains the set of policies.
type PolicyList struct {
	// Policies are ordered from least to greatest precidence.
	Policies []Policy `json:"policy" validate:"dive"`

	// Assertions are used to validate the ordering of Policies.
	Assertions []Policy `json:"assert" validate:"dive"`
}

// Validate checks that the PolicyList struct is valid, that the keys for the
// conditions are valid and that all assertions hold.
func (pl *PolicyList) Validate(validConditionKeys []string) error {
	if err := validation.ValidateStruct(pl); err != nil {
		return fmt.Errorf("invalid PolicyList: %v", err)
	}

	for i, pol := range pl.Policies {
		if err := pol.Condition.ValidateKeys(validConditionKeys); err != nil {
			return fmt.Errorf("error in policy[%d], comment: %q, error: %v", i, pol.Comment, err)
		}
	}

	return pl.CheckAssertions()
}

// CheckAssertions tests each assertion in the Assertions list against the
// policies list. The condition is used as the ground truth and the
// actions are used as the expected output. If the actions don't match then
// an error is returned.
func (pl *PolicyList) CheckAssertions() error {
	for i, assertion := range pl.Assertions {
		expected := assertion.Declarations
		actual := pl.Apply(assertion.Condition)

		if !reflect.DeepEqual(actual, expected) {
			return fmt.Errorf("error in assertion[%d], comment: %q, expected: %v, actual: %v", i, assertion.Comment, expected, actual)
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

		for k, v := range policy.Declarations {
			out[k] = v
		}
	}

	return out
}

// NewPolicyListFromJson creates a PolicyList from the given JSON version.
// It will fail on invalid condition names and failed assertions.
//
// Exactly one of PolicyList or error will be returned.
func NewPolicyListFromJson(value json.RawMessage, validConditionKeys []string) (*PolicyList, error) {
	decoder := json.NewDecoder(bytes.NewBuffer(value))

	// be noisy if the structure is invalid
	decoder.DisallowUnknownFields()

	pl := &PolicyList{}
	if err := decoder.Decode(pl); err != nil {
		return nil, fmt.Errorf("couldn't decode PolicyList from JSON: %v", err)
	}

	if err := pl.Validate(validConditionKeys); err != nil {
		return nil, err
	}

	return pl, nil
}
