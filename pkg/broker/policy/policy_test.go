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

import "testing"

func TestCondition_AppliesTo(t *testing.T) {
	cases := map[string]struct {
		Condition Condition
		Truth     Condition
		Expected  bool
	}{
		"matches everything": {
			Condition: Condition{},
			Truth:     Condition{},
			Expected:  true,
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			actual := tc.Condition.AppliesTo(tc.Truth)

			if tc.Expected != actual {
				t.Errorf("Expected condition to apply? %t but was: %t", tc.Expected, actual)
			}
		})
	}
}

const examplePolicy = `
{
  "policies":[
    {
      "//":"always applies",
      "if":{},
      "then":{
        "cascade-true": true,
        "cascade-false": true

      }
    },
    {
      "//":"some comment here",
      "if":{
        "service_name": "cloud-storage"
      },
      "then":{
        "cascade-false": false
      }
    }
  ],
  "assert":[
  {
    "//":"cascading works correctly",
    "if":{
      "service_name": "cloud-storage"
    },
    "then":{
      "cascade-true": true,
      "cascade-false": false
    }
  }


  ],
}
`
