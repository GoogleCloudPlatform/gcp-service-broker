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

package validation

import (
	"reflect"
	"testing"
)

func TestConstraintBuilder(t *testing.T) {

	// NOTE: Keep keys strings rather than constants in Expected so we can also
	// validate that our constants don't get changed.
	cases := map[string]struct {
		Constraints ConstraintBuilder
		Expected    map[string]interface{}
	}{
		"empty": {
			Constraints: NewConstraintBuilder().Build(),
			Expected:    map[string]interface{}{},
		},
		"annotations": {
			Constraints: NewConstraintBuilder().
				Description("desc").
				Examples("exa", "exb").
				Type("string"),
			Expected: map[string]interface{}{
				"description": "desc",
				"examples":    []interface{}{"exa", "exb"},
				"type":        "string",
			},
		},

		"any type": {
			Constraints: NewConstraintBuilder().
				Enum("a", "b", "c").
				Const("exa"),
			Expected: map[string]interface{}{
				"enum":  []interface{}{"a", "b", "c"},
				"const": "exa",
			},
		},

		"numeric": {
			Constraints: NewConstraintBuilder().
				Maximum(3).
				Minimum(1).
				ExclusiveMaximum(4).
				ExclusiveMinimum(0).
				MultipleOf(1),
			Expected: map[string]interface{}{
				"maximum":          3,
				"minimum":          1,
				"exclusiveMaximum": 4,
				"exclusiveMinimum": 0,
				"multipleOf":       1,
			},
		},

		"strings": {
			Constraints: NewConstraintBuilder().
				MaxLength(30).
				MinLength(10).
				Pattern("^[A-Za-z]+[A-Za-z0-9]+$"),
			Expected: map[string]interface{}{
				"maxLength": 30,
				"minLength": 10,
				"pattern":   "^[A-Za-z]+[A-Za-z0-9]+$",
			},
		},

		"arrays": {
			Constraints: NewConstraintBuilder().
				MaxItems(30).
				MinItems(10),
			Expected: map[string]interface{}{
				"maxItems": 30,
				"minItems": 10,
			},
		},

		"objects": {
			Constraints: NewConstraintBuilder().
				MaxProperties(30).
				MinProperties(10).
				Required("a", "b", "c").
				PropertyNames(map[string]interface{}{"type": "string"}),
			Expected: map[string]interface{}{
				"maxProperties": 30,
				"minProperties": 10,
				"required":      []string{"a", "b", "c"},
				"propertyNames": map[string]interface{}{"type": "string"},
			},
		},

		"secondOverwritesFirst": {
			Constraints: NewConstraintBuilder().MaxLength(3).MaxLength(5).Build(),
			Expected: map[string]interface{}{
				KeyMaxLength: 5,
			},
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			actual := tc.Constraints.Build()

			if !reflect.DeepEqual(tc.Expected, actual) {
				t.Errorf("Error constructing constraints, expected: %#v got: %#v", tc.Expected, actual)
			}
		})
	}
}
