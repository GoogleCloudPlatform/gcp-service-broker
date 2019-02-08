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

package migration

import "github.com/GoogleCloudPlatform/gcp-service-broker/utils"

// Diff holds the difference between two strings.
type Diff struct {
	Old string `yaml:"old,omitempty"`
	New string `yaml:"new"` // new is intentionally not omitempty to show change
}

// DiffStringMap creates a diff between the two maps.
func DiffStringMap(old, new map[string]string) map[string]Diff {
	allKeys := utils.NewStringSetFromStringMapKeys(old)
	newKeys := utils.NewStringSetFromStringMapKeys(new)
	allKeys.Add(newKeys.ToSlice()...)

	out := make(map[string]Diff)
	for _, k := range allKeys.ToSlice() {
		if old[k] != new[k] {
			out[k] = Diff{Old: old[k], New: new[k]}
		}
	}

	return out
}
