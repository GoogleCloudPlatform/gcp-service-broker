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

package varcontext

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/broker"
)

func TestContextBuilder(t *testing.T) {
	cases := map[string]struct {
		Builder     *ContextBuilder
		Expected    map[string]interface{}
		ErrContains string
	}{
		"an empty context": {
			Builder:     Builder(),
			Expected:    map[string]interface{}{},
			ErrContains: "",
		},

		// MergeMap
		"MergeMap blank okay": {
			Builder:  Builder().MergeMap(map[string]interface{}{}),
			Expected: map[string]interface{}{},
		},

		"MergeMap multi-key": {
			Builder:  Builder().MergeMap(map[string]interface{}{"a": "a", "b": "b"}),
			Expected: map[string]interface{}{"a": "a", "b": "b"},
		},

		"MergeMap overwrite": {
			Builder:  Builder().MergeMap(map[string]interface{}{"a": "a"}).MergeMap(map[string]interface{}{"a": "aaa"}),
			Expected: map[string]interface{}{"a": "aaa"},
		},

		// nil default, non-string default, string default
		// func (builder *ContextBuilder) MergeDefaults(brokerVariables []broker.BrokerVariable) *ContextBuilder {
		"MergeDefaults no defaults": {
			Builder:  Builder().MergeDefaults([]broker.BrokerVariable{{FieldName: "foo"}}),
			Expected: map[string]interface{}{},
		},

		"MergeDefaults non-string": {
			Builder:  Builder().MergeDefaults([]broker.BrokerVariable{{FieldName: "h2g2", Default: 42}}),
			Expected: map[string]interface{}{"h2g2": 42},
		},

		"MergeDefaults basic-string": {
			Builder:  Builder().MergeDefaults([]broker.BrokerVariable{{FieldName: "a", Default: "no-template"}}),
			Expected: map[string]interface{}{"a": "no-template"},
		},

		"MergeDefaults template string": {
			Builder:  Builder().MergeDefaults([]broker.BrokerVariable{{FieldName: "a", Default: "a"}, {FieldName: "b", Default: "${a}"}}),
			Expected: map[string]interface{}{"a": "a", "b": "a"},
		},

		// MergeEvalResult
		"MergeEvalResult accumulates context": {
			Builder:  Builder().MergeEvalResult("a", "a").MergeEvalResult("b", "${a}"),
			Expected: map[string]interface{}{"a": "a", "b": "a"},
		},
		"MergeEvalResult errors": {
			Builder:     Builder().MergeEvalResult("a", "${dne}"),
			ErrContains: `couldn't compute the value for "a"`,
		},

		// MergeJsonObject
		"MergeJsonObject blank message": {
			Builder:  Builder().MergeJsonObject(json.RawMessage{}),
			Expected: map[string]interface{}{},
		},
		"MergeJsonObject valid message": {
			Builder:  Builder().MergeJsonObject(json.RawMessage(`{"a":"a"}`)),
			Expected: map[string]interface{}{"a": "a"},
		},
		"MergeJsonObject invalid message": {
			Builder:     Builder().MergeJsonObject(json.RawMessage(`{{{}}}`)),
			ErrContains: "invalid character '{'",
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {

			vc, err := tc.Builder.Build()

			if vc == nil && tc.Expected != nil {
				t.Fatalf("Expected: %v, got: %v", tc.Expected, vc)
			}

			if vc != nil && !reflect.DeepEqual(vc.ToMap(), tc.Expected) {
				t.Errorf("Expected: %v, got: %v", tc.Expected, vc.ToMap())
			}

			switch {
			case err == nil && tc.ErrContains == "":
				break
			case err == nil && tc.ErrContains != "":
				t.Errorf("Got no error when %q was expected", tc.ErrContains)
			case err != nil && tc.ErrContains == "":
				t.Errorf("Got error %v when none was expected", err)
			case !strings.Contains(err.Error(), tc.ErrContains):
				t.Errorf("Got error %v, but expected it to contain %q", err, tc.ErrContains)
			}
		})
	}
}
