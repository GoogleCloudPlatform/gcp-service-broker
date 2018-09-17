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

// ServiceBindingCredentials holds credentials returned to the users after
// binding to a service.
type ServiceBindingCredentials ServiceBindingCredentialsV1

// GetOtherDetails returns an unmarshaled version of the OtherDetails field
// or errors.
func (sbc ServiceBindingCredentials) GetOtherDetails() (map[string]string, error) {
	var creds map[string]string
	err := json.Unmarshal([]byte(sbc.OtherDetails), &creds)
	return creds, err
}

// ServiceInstanceDetails holds information about provisioned services 
type ServiceInstanceDetails ServiceInstanceDetailsV1

// GetOtherDetails returns an unmarshaled version of the OtherDetails field
// or errors.
func (si ServiceInstanceDetails) GetOtherDetails() (map[string]string, error) {
	var instanceDetails map[string]string
	if si.OtherDetails == "" {
		return instanceDetails, nil
	}

	err := json.Unmarshal([]byte(si.OtherDetails), &instanceDetails)
	return instanceDetails, err
}

// ProvisionRequestDetails holds user-defined properties passed to a call
// to provision a service.
type ProvisionRequestDetails ProvisionRequestDetailsV1

// Migration represents the mgirations table. It holds a monotonically
// increasing number that gets incremented with every database schema revision.
type Migration MigrationV1

// CloudOperation holds information about the status of Google Cloud
// long-running operations.
type CloudOperation CloudOperationV1
