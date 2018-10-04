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

package broker

import (
	"fmt"
	"sort"

	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/validation"
)

const (
	JsonTypeString  JsonType = "string"
	JsonTypeNumeric JsonType = "number"
	JsonTypeInteger JsonType = "integer"
	JsonTypeBoolean JsonType = "boolean"
)

type JsonType string

type BrokerVariable struct {
	// Is this variable required?
	Required bool
	// The name of the JSON field this variable serializes/deserializes to
	FieldName string
	// The JSONSchema type of the field
	Type JsonType
	// Human readable info about the field.
	Details string
	// The default value of the field.
	Default interface{}
	// If there are a limited number of valid values for this field then
	// Enum will hold them in value:friendly name pairs
	Enum map[interface{}]string
	// Constraints holds JSON Schema validations defined for this variable.
	// Keys are valid JSON Schema validation keywords, and values are their
	// associated values.
	// http://json-schema.org/latest/json-schema-validation.html
	Constraints map[string]interface{}
}

// ToSchema converts the BrokerVariable into the value part of a JSON Schema.
func (bv *BrokerVariable) ToSchema() map[string]interface{} {
	schema := map[string]interface{}{}

	for k, v := range bv.Constraints {
		schema[k] = v
	}

	if len(bv.Enum) > 0 {
		enumeration := []interface{}{}
		for k, _ := range bv.Enum {
			enumeration = append(enumeration, k)
		}

		// Sort enumerations lexocographically for documentation consistency.
		sort.Slice(enumeration, func(i int, j int) bool {
			return fmt.Sprintf("%v", enumeration[i]) < fmt.Sprintf("%v", enumeration[j])
		})

		schema[validation.KeyEnum] = enumeration
	}

	if bv.Details != "" {
		schema[validation.KeyDescription] = bv.Details
	}

	if bv.Type != "" {
		schema[validation.KeyType] = bv.Type
	}

	if bv.Default != nil {
		schema[validation.KeyDefault] = bv.Default
	}

	return schema
}
