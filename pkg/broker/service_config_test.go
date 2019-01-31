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

package broker

import (
	"errors"
	"reflect"
	"testing"

	"github.com/pivotal-cf/brokerapi"
	"github.com/spf13/viper"
)

func TestServiceConfig_ToServicePlan(t *testing.T) {
	cases := map[string]struct {
		Plan     CustomPlan
		Expected ServicePlan
	}{
		"nominal": {
			Plan: CustomPlan{
				GUID:        "00000000-0000-0000-0000-000000000000",
				Name:        "my-normal-plan",
				DisplayName: "My Normal Plan",
				Description: "Some normal plan",
				Properties:  map[string]string{"a": "b"},
			},
			Expected: ServicePlan{
				ServicePlan: brokerapi.ServicePlan{
					Description: "Some normal plan",
					Name:        "my-normal-plan",
					ID:          "00000000-0000-0000-0000-000000000000",
					Metadata: &brokerapi.ServicePlanMetadata{
						DisplayName: "My Normal Plan",
					},
				},
				ServiceProperties: map[string]string{"a": "b"},
			},
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			actual := tc.Plan.ToServicePlan()

			if !reflect.DeepEqual(tc.Expected, actual) {
				t.Errorf("Expected: %v, Actual: %v", tc.Expected, actual)
			}
		})
	}
}

func TestNewServiceConfigMapFromEnv(t *testing.T) {
	cases := map[string]struct {
		Env         string
		ExpectedErr error
		Expected    ServiceConfigMap
	}{
		"blank": {
			Expected: ServiceConfigMap{},
		},
		"invalid-json": {
			Env:         "{",
			ExpectedErr: errors.New("couldn't deserialize ServiceConfigMap: unexpected end of JSON input"),
		},
		"multiple-objects": {
			Env: `{"00000000-0000-0000-0000-000000000000":{}, "00000000-0000-0000-0000-000000000001":{}}`,
			Expected: ServiceConfigMap{
				"00000000-0000-0000-0000-000000000000": ServiceConfig{},
				"00000000-0000-0000-0000-000000000001": ServiceConfig{},
			},
		},
		"populated-object": {
			Env: `{"00000000-0000-0000-0000-000000000000":{
        "provision_defaults":{"pd1":"pdv1", "pd2":"pdv2"},
        "bind_defaults":{"bd1":"bdv1", "bd2":"bdv2"},
        "custom_plans":[
          {
            "guid":"00000000-0000-0000-0000-000000000001",
            "name":"my-service-name",
            "display_name":"my-service-display-name",
            "description":"my-service-description",
            "properties":{"pk1":"pv1"}
          }
        ]
        }
      }`,
			Expected: ServiceConfigMap{
				"00000000-0000-0000-0000-000000000000": ServiceConfig{
					ProvisionDefaults: map[string]interface{}{"pd1": "pdv1", "pd2": "pdv2"},
					BindDefaults:      map[string]interface{}{"bd1": "bdv1", "bd2": "bdv2"},
					CustomPlans: []CustomPlan{
						{
							GUID:        "00000000-0000-0000-0000-000000000001",
							Name:        "my-service-name",
							DisplayName: "my-service-display-name",
							Description: "my-service-description",
							Properties:  map[string]string{"pk1": "pv1"},
						},
					},
				},
			},
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			viper.Set(ServiceConfigProperty, tc.Env)
			defer viper.Reset()

			actual, actualErr := NewServiceConfigMapFromEnv()
			if !reflect.DeepEqual(tc.ExpectedErr, actualErr) {
				t.Errorf("Expected: %v, Actual: %v", tc.ExpectedErr, actualErr)
			}

			if !reflect.DeepEqual(tc.Expected, actual) {
				t.Errorf("Expected: %v, Actual: %v", tc.Expected, actual)
			}
		})
	}
}
