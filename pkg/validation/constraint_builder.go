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

package validation

const (
	KeyDefault          = "default"
	KeyExamples         = "examples"
	KeyDescription      = "description"
	KeyTitle            = "title"
	KeyType             = "type"
	KeyConst            = "const"
	KeyEnum             = "enum"
	KeyMultipleOf       = "multipleOf"
	KeyMaximum          = "maximum"
	KeyMinimum          = "minimum"
	KeyExclusiveMaximum = "exclusiveMaximum"
	KeyExclusiveMinimum = "exclusiveMinimum"
	KeyMaxLength        = "maxLength"
	KeyMinLength        = "minLength"
	KeyPattern          = "pattern"
	KeyMaxItems         = "maxItems"
	KeyMinItems         = "minItems"
	KeyMaxProperties    = "maxProperties"
	KeyMinProperties    = "minProperties"
	KeyRequired         = "required"
	KeyPropertyNames    = "propertyNames"
)

//  NewConstraintBuilder creates a builder for JSON Schema compliant constraint
// lists. See http://json-schema.org/latest/json-schema-validation.html
// for types of validation available.
func NewConstraintBuilder() ConstraintBuilder {
	return ConstraintBuilder{}
}

// A builder for JSON Schema compliant constraint lists
type ConstraintBuilder map[string]interface{}

// Type adds a type constrinat.
func (cb ConstraintBuilder) Type(t string) ConstraintBuilder {
	cb[KeyType] = t
	return cb
}

// Description adds a human-readable description
func (cb ConstraintBuilder) Description(desc string) ConstraintBuilder {
	cb[KeyDescription] = desc

	return cb
}

// Title adds a human-readable label suitable for labeling a UI element.
func (cb ConstraintBuilder) Title(title string) ConstraintBuilder {
	cb[KeyTitle] = title

	return cb
}

// Examples adds one or more examples
func (cb ConstraintBuilder) Examples(ex ...interface{}) ConstraintBuilder {
	cb[KeyExamples] = ex

	return cb
}

// Const adds a constraint that the field must equal this value.
func (cb ConstraintBuilder) Const(value interface{}) ConstraintBuilder {
	cb[KeyConst] = value

	return cb
}

// Enum adds a constraint that the field must be one of these values.
func (cb ConstraintBuilder) Enum(value ...interface{}) ConstraintBuilder {
	cb[KeyEnum] = value

	return cb
}

// MultipleOf adds a constraint that the field must be a multiple of this integer.
func (cb ConstraintBuilder) MultipleOf(value int) ConstraintBuilder {
	cb[KeyMultipleOf] = value

	return cb
}

// Minimum adds a constraint that the field must be greater than or equal to
// this number.
func (cb ConstraintBuilder) Minimum(value int) ConstraintBuilder {
	cb[KeyMinimum] = value

	return cb
}

// Maximum adds a constraint that the field must be less than or equal to
// this number.
func (cb ConstraintBuilder) Maximum(value int) ConstraintBuilder {
	cb[KeyMaximum] = value

	return cb
}

// ExclusiveMaximum adds a constraint that the field must be less than this number.
func (cb ConstraintBuilder) ExclusiveMaximum(value int) ConstraintBuilder {
	cb[KeyExclusiveMaximum] = value

	return cb
}

// ExclusiveMinimum adds a constraint that the field must be greater than this number.
func (cb ConstraintBuilder) ExclusiveMinimum(value int) ConstraintBuilder {
	cb[KeyExclusiveMinimum] = value

	return cb
}

// MaxLength adds a constraint that the string field must have at most this many characters.
func (cb ConstraintBuilder) MaxLength(value int) ConstraintBuilder {
	cb[KeyMaxLength] = value

	return cb
}

// MinLength adds a constraint that the string field must have at least this many characters.
func (cb ConstraintBuilder) MinLength(value int) ConstraintBuilder {
	cb[KeyMinLength] = value

	return cb
}

// Pattern adds a constraint that the string must match the given pattern.
func (cb ConstraintBuilder) Pattern(value string) ConstraintBuilder {
	cb[KeyPattern] = value

	return cb
}

// MaxItems adds a constraint that the array must have at most this many items.
func (cb ConstraintBuilder) MaxItems(value int) ConstraintBuilder {
	cb[KeyMaxItems] = value

	return cb
}

// MinItems adds a constraint that the array must have at least this many items.
func (cb ConstraintBuilder) MinItems(value int) ConstraintBuilder {
	cb[KeyMinItems] = value

	return cb
}

// KeyMaxProperties adds a constraint that the object must have at most this many keys.
func (cb ConstraintBuilder) MaxProperties(value int) ConstraintBuilder {
	cb[KeyMaxProperties] = value

	return cb
}

// MinProperties adds a constraint that the object must have at least this many keys.
func (cb ConstraintBuilder) MinProperties(value int) ConstraintBuilder {
	cb[KeyMinProperties] = value

	return cb
}

// Required adds a constraint that the object must have at least these keys.
func (cb ConstraintBuilder) Required(properties ...string) ConstraintBuilder {
	cb[KeyRequired] = properties

	return cb
}

// PropertyNames adds a constraint that the object property names must match the given schema.
func (cb ConstraintBuilder) PropertyNames(properties map[string]interface{}) ConstraintBuilder {
	cb[KeyPropertyNames] = properties

	return cb
}

func (cb ConstraintBuilder) Build() map[string]interface{} {
	return cb
}
