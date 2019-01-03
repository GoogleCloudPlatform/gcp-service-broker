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
	"reflect"
	"testing"

	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/broker"
)

func TestServiceAccountWhitelistWithDefault(t *testing.T) {
	details := `The role for the account without the "roles/" prefix. See: https://cloud.google.com/iam/docs/understanding-roles for more details.`

	cases := map[string]struct {
		Whitelist   []string
		DefaultRole string
		Expected    broker.BrokerVariable
	}{
		"default in whitelist": {
			Whitelist:   []string{"foo"},
			DefaultRole: "foo",
			Expected: broker.BrokerVariable{
				FieldName: "role",
				Type:      broker.JsonTypeString,
				Details:   details,

				Required: false,
				Default:  "foo",
				Enum:     map[interface{}]string{"foo": "roles/foo"},
			},
		},

		"default not in whitelist": {
			Whitelist:   []string{"foo"},
			DefaultRole: "test",
			Expected: broker.BrokerVariable{
				FieldName: "role",
				Type:      broker.JsonTypeString,
				Details:   details,

				Required: false,
				Default:  "test",
				Enum:     map[interface{}]string{"foo": "roles/foo"},
			},
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			vars := ServiceAccountWhitelistWithDefault(tc.Whitelist, tc.DefaultRole)
			if len(vars) != 1 {
				t.Fatalf("Expected 1 input variable, got %d", len(vars))
			}

			if !reflect.DeepEqual(vars[0], tc.Expected) {
				t.Fatalf("Expected %#v, got %#v", tc.Expected, vars[0])
			}
		})
	}
}
