// Copyright the Service Broker Project Authors.
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
}
