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

package brokerpak

import (
	"reflect"
	"testing"
)

func TestSplitBrokerpakList(t *testing.T) {
	cases := map[string]struct {
		Input    string
		Expected []string
	}{
		"none": {
			Input:    ``,
			Expected: nil,
		},
		"single": {
			Input:    `gs://foo/bar`,
			Expected: []string{"gs://foo/bar"},
		},
		"crlf": {
			Input:    "a://foo\r\nb://bar",
			Expected: []string{"a://foo", "b://bar"},
		},
		"trim": {
			Input:    "  a://foo  \n\tb://bar  ",
			Expected: []string{"a://foo", "b://bar"},
		},
		"blank": {
			Input:    "\n\r\r\n\n\t\t\v\n\n\n\n\n \n\n\n \n \n \n\n",
			Expected: nil,
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			actual := splitBrokerpakList(tc.Input)

			if !reflect.DeepEqual(tc.Expected, actual) {
				t.Errorf("Expected: %v actual: %v", tc.Expected, actual)
			}
		})
	}
}
