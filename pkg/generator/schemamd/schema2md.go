package schemamd

import (
	"encoding/json"
	"fmt"
)

type Config struct {
	HeadingOffset int
	AnchorPrefix  string
	FieldName     string
}

type JsonSchema struct {
	Type    string          `json:"type"`
	Default json.RawMessage `json:"default"`

	// documentation
	Title       string        `json:"title"`
	Description string        `json:"description"`
	Examples    []interface{} `json:"examples"`

	// object settings
	Properties    map[string]JsonSchema `json:"parameters"`
	Required      []string              `json:"required"`
	PropertyNames []string              `json:"propertyNames"`
	MaxProperties interface{}           `json:"maxProperties"`
	MinProperties interface{}           `json:"minProperties"`

	// array settings
	MaxItems    interface{} `json:"maxItems"`
	MinItems    interface{} `json:"minItems"`
	UniqueItems bool        `json:"uniqueItems"`

	// primitive settings
	Const            interface{}   `json:"const"`
	Enum             []interface{} `json:"enum"`
	MultipleOf       interface{}   `json:"multipleOf"`
	Maximum          interface{}   `json:"maximum"`
	Minimum          interface{}   `json:"minimum"`
	ExclusiveMaximum interface{}   `json:"exclusiveMaximum"`
	ExclusiveMinimum interface{}   `json:"exclusiveMinimum"`
	MaxLength        interface{}   `json:"maxLength"`
	MinLength        interface{}   `json:"minLength"`
	Pattern          interface{}   `json:"pattern"`
}

func (schema *JsonSchema) Markdown(cfg Config) (string, error) {
	// heading
	// description
	// if object
	//  table
	//  for each item, list and add constraint
	// if array
	//  notes
	// if item
	//  notes, constraint

	switch schema.Type {
	case "object":
		return schema.objectMarkdown(cfg)
	case "array":
		return "", fmt.Errorf("Documentation for arrays is not supported.")
	case "string", "number", "integer", "boolean", "null":
		// render basic type
		return schema.primitiveMarkdown(cfg)
	default:
		return "", fmt.Errorf("Unknown type: %s", schema.Type)
	}
}

func (schema *JsonSchema) objectMarkdown(cfg Config) (string, error) {

}

func (schema *JsonSchema) primitiveMarkdown(cfg Config) (string, error) {

}

// "parameters": {
//   "$schema": "http://json-schema.org/draft-04/schema#",
//   "properties": {
//     "location": {
//       "default": "US",
//       "description": "The location of the bucket. Object data for objects in the bucket resides in physical storage within this region. See: https://cloud.google.com/storage/docs/bucket-locations",
//       "examples": [
//         "US",
//         "EU",
//         "southamerica-east1"
//       ],
//       "pattern": "^[A-Za-z][-a-z0-9A-Z]+$",
//       "title": "Location",
//       "type": "string"
//     },
//     "name": {
//       "default": "pcf_sb_${counter.next()}_${time.nano()}",
//       "description": "The name of the bucket. There is a single global namespace shared by all buckets so it MUST be unique.",
//       "maxLength": 222,
//       "minLength": 3,
//       "pattern": "^[A-Za-z0-9_\\.]+$",
//       "title": "Name",
//       "type": "string"
//     }
//   },
//   "type": "object"
// }
// },
