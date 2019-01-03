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

package builtin

import (
	"reflect"
	"sort"
	"testing"
)

func TestBuiltinBrokerRegistry(t *testing.T) {
	builtinServiceNames := []string{
		"google-bigquery",
		"google-bigtable",
		"google-cloudsql-mysql",
		"google-cloudsql-postgres",
		"google-dataflow",
		"google-datastore",
		"google-dialogflow",
		"google-firestore",
		"google-ml-apis",
		"google-pubsub",
		"google-spanner",
		"google-stackdriver-debugger",
		"google-stackdriver-monitoring",
		"google-stackdriver-profiler",
		"google-stackdriver-trace",
		"google-storage",
	}

	sort.Strings(builtinServiceNames)

	t.Run("service-count", func(t *testing.T) {
		expectedServiceCount := len(builtinServiceNames)
		actualServiceCount := len(BuiltinBrokerRegistry())
		if actualServiceCount != expectedServiceCount {
			t.Errorf("Expected %d services, registered: %d", expectedServiceCount, actualServiceCount)
		}
	})

	t.Run("service-names", func(t *testing.T) {
		var actual []string
		registry := BuiltinBrokerRegistry()
		for name, _ := range registry {
			actual = append(actual, name)
		}

		sort.Strings(actual)
		if !reflect.DeepEqual(builtinServiceNames, actual) {
			t.Errorf("Expected service names: %v, got: %v", builtinServiceNames, actual)
		}
	})

	t.Run("has-builtin-flag", func(t *testing.T) {
		registry := BuiltinBrokerRegistry()
		for _, svc := range registry {
			if !svc.IsBuiltin {
				t.Errorf("Expected flag 'builtin' to be set for %s, but it was: %t", svc.Name, svc.IsBuiltin)
			}
		}
	})
}
