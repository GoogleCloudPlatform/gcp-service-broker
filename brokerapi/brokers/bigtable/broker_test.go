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

package bigtable

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/pivotal-cf/brokerapi"
)

func TestBigTableBroker_ProvisionVariables(t *testing.T) {
	service := serviceDefinition()

	hddPlan := "65a49268-2c73-481e-80f3-9fde5bd5a654"
	ssdPlan := "38aa0e65-624b-4998-9c06-f9194b56d252"

	cases := map[string]struct {
		UserParams      string
		PlanId          string
		ExpectedContext map[string]interface{}
	}{
		"hdd": {
			UserParams: `{"name":"my-bt-instance"}`,
			PlanId:     hddPlan,
			ExpectedContext: map[string]interface{}{
				"num_nodes":    "3",
				"name":         "my-bt-instance",
				"cluster_id":   "my-bt-instance-cluster",
				"display_name": "my-bt-instance",
				"zone":         "us-east1-b",
				"storage_type": "HDD",
			},
		},
		"ssd": {
			UserParams: `{"name":"my-bt-instance"}`,
			PlanId:     ssdPlan,
			ExpectedContext: map[string]interface{}{
				"num_nodes":    "3",
				"name":         "my-bt-instance",
				"cluster_id":   "my-bt-instance-cluster",
				"display_name": "my-bt-instance",
				"zone":         "us-east1-b",
				"storage_type": "SSD",
			},
		},
		"cluster truncates": {
			UserParams: `{"name":"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"}`,
			PlanId:     ssdPlan,
			ExpectedContext: map[string]interface{}{
				"num_nodes":    "3",
				"name":         "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
				"cluster_id":   "aaaaaaaaaaaaaaaaaaaa-cluster",
				"display_name": "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
				"zone":         "us-east1-b",
				"storage_type": "SSD",
			},
		},
		"no defaults": {
			UserParams: `{"name":"test", "cluster_id": "testcluster", "display_name":"test display"}`,
			PlanId:     ssdPlan,
			ExpectedContext: map[string]interface{}{
				"num_nodes":    "3",
				"name":         "test",
				"cluster_id":   "testcluster",
				"display_name": "test display",
				"zone":         "us-east1-b",
				"storage_type": "SSD",
			},
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			details := brokerapi.ProvisionDetails{RawParameters: json.RawMessage(tc.UserParams)}
			plan, err := service.GetPlanById(tc.PlanId)
			if err != nil {
				t.Errorf("got error trying to find plan %s %v", tc.PlanId, err)
			}
			vars, err := service.ProvisionVariables("instance-id-here", details, *plan)

			if err != nil {
				t.Errorf("got error while creating provision variables: %v", err)
			}

			if !reflect.DeepEqual(vars.ToMap(), tc.ExpectedContext) {
				t.Errorf("Expected context: %#v got %#v", tc.ExpectedContext, vars.ToMap())
			}
		})
	}
}
