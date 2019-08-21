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
	"errors"
	"testing"

	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/validation"
)

func TestTerraformResource_Validate(t *testing.T) {
	cases := map[string]validation.ValidatableTest{
		"blank obj": {
			Object: &TerraformResource{},
			Expect: errors.New("missing field(s): name, source, version"),
		},
		"good obj": {
			Object: &TerraformResource{
				Name:    "foo",
				Version: "1.0",
				Source:  "github.com/myproject",
			},
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			tc.Assert(t)
		})
	}
}
