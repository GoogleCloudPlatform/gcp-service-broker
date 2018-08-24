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

package models

import (
	"encoding/json"
)

type ServiceBindingCredentials ServiceBindingCredentialsV1

func (sbc ServiceBindingCredentials) GetOtherDetails() map[string]string {
	var creds map[string]string
	if err := json.Unmarshal([]byte(sbc.OtherDetails), &creds); err != nil {
		panic(err)
	}
	return creds
}

type ServiceInstanceDetails ServiceInstanceDetailsV1

func (si ServiceInstanceDetails) GetOtherDetails() map[string]string {
	var instanceDetails map[string]string
	// if the instance has access details saved
	if si.OtherDetails != "" {
		if err := json.Unmarshal([]byte(si.OtherDetails), &instanceDetails); err != nil {
			panic(err)
		}
	}
	return instanceDetails

}

type ProvisionRequestDetails ProvisionRequestDetailsV1
type Migration MigrationV1
type CloudOperation CloudOperationV1
