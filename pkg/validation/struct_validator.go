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

import (
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"regexp"

	"github.com/hashicorp/hcl"
	validator "gopkg.in/go-playground/validator.v9"
)

var (
	osbName                  = `^[a-zA-Z0-9-\.]+$`
	osbNameRegex             = regexp.MustCompile(osbName)
	terraformIdentifier      = `^[a-z_]*$`
	terraformIdentifierRegex = regexp.MustCompile(terraformIdentifier)
	jsonSchemaType           = `^(|object|boolean|array|number|string|integer)$`
	jsonSchemaTypeRegex      = regexp.MustCompile(jsonSchemaType)

	uuidRegex = regexp.MustCompile(`^[0-9a-fA-F]{8}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{12}$`)
)

var validate = validator.New()

func init() {
	validate.RegisterValidation("osbname", regexValidation(osbNameRegex))
	validate.RegisterValidation("json", jsonValidation)
	validate.RegisterValidation("hcl", hclValidation)
	validate.RegisterValidation("terraform_identifier", regexValidation(terraformIdentifierRegex))
	validate.RegisterValidation("jsonschema_type", regexValidation(jsonSchemaTypeRegex))
}

// ValidateStruct executes the validation tags on a struct and returns any
// failures.
func ValidateStruct(s interface{}) error {
	return validate.Struct(s)
}

func not(inner validator.Func) validator.Func {
	return func(field validator.FieldLevel) bool {
		return !inner(field)
	}
}

func regexValidation(matcher *regexp.Regexp) validator.Func {
	return func(field validator.FieldLevel) bool {
		return matcher.MatchString(field.Field().String())
	}
}

func jsonValidation(field validator.FieldLevel) bool {
	fl := field.Field()

	switch fl.Type().Kind() {
	case reflect.String:
		return IsJSONString(fl.String())
	default:
		return IsJSONBytes(fl.Bytes())
	}
}

func hclValidation(field validator.FieldLevel) bool {
	value := field.Field().String()
	return IsHCL(value)
}

// IsHCL validates that a value is valid HCL.
func IsHCL(value string) bool {
	_, err := hcl.Parse(value)
	return err == nil
}

// IsJSONString returns true if the string is valid JSON.
func IsJSONString(value string) bool {
	return IsJSONBytes([]byte(value))
}

// IsJSONBytes returns true if the byte array is valid JSON.
func IsJSONBytes(value []byte) bool {
	return json.Valid(value)
}

// ErrIfNotHCL returns an error if the value is not valid HCL.
func ErrIfNotHCL(value string, field string) *FieldError {
	if _, err := hcl.Parse(value); err == nil {
		return nil
	}

	return &FieldError{
		Message: "invalid HCL",
		Paths:   []string{field},
	}
}

// ErrIfNotJSON returns an error if the value is not valid JSON.
func ErrIfNotJSON(value json.RawMessage, field string) *FieldError {
	if json.Valid(value) {
		return nil
	}

	return &FieldError{
		Message: "invalid JSON",
		Paths:   []string{field},
	}
}

// ErrIfBlank returns an error if the value is a blank string.
func ErrIfBlank(value string, field string) *FieldError {
	if value == "" {
		return ErrMissingField(field)
	}

	return nil
}

// ErrIfNil returns an error if the value is nil.
func ErrIfNil(value interface{}, field string) *FieldError {
	if value == nil {
		return ErrMissingField(field)
	}

	return nil
}

// ErrIfNotOSBName returns an error if the value is not a valid OSB name.
func ErrIfNotOSBName(value string, field string) *FieldError {
	return ErrIfNotMatch(value, osbNameRegex, field)
}

// ErrIfNotJSONSchemaType returns an error if the value is not a valid JSON
// schema type.
func ErrIfNotJSONSchemaType(value string, field string) *FieldError {
	return ErrIfNotMatch(value, jsonSchemaTypeRegex, field)
}

// ErrIfNotTerraformIdentifier returns an error if the value is not a valid
// Terraform identifier.
func ErrIfNotTerraformIdentifier(value string, field string) *FieldError {
	return ErrIfNotMatch(value, terraformIdentifierRegex, field)
}

// ErrIfNotUUID returns an error if the value is not a valid UUID.
func ErrIfNotUUID(value string, field string) *FieldError {
	if uuidRegex.MatchString(value) {
		return nil
	}

	return &FieldError{
		Message: "field must be a UUID",
		Paths:   []string{field},
	}
}

// ErrIfNotURL returns an error if the value is not a valid URL.
func ErrIfNotURL(value string, field string) *FieldError {
	// Validaiton inspired by: github.com/go-playground/validator/baked_in.go
	url, err := url.ParseRequestURI(value)
	if err != nil || url.Scheme == "" {
		return &FieldError{
			Message: "field must be a URL",
			Paths:   []string{field},
		}
	}

	return nil
}

// ErrIfNotMatch returns an error if the value doesn't match the regex.
func ErrIfNotMatch(value string, regex *regexp.Regexp, field string) *FieldError {
	if regex.MatchString(value) {
		return nil
	}

	return ErrMustMatch(value, regex, field)
}

// ErrMustMatch notifies the user a field must match a regex.
func ErrMustMatch(value string, regex *regexp.Regexp, field string) *FieldError {
	return &FieldError{
		Message: fmt.Sprintf("field must match '%s'", regex.String()),
		Paths:   []string{field},
	}
}

// Validatable indicates that a particular type may have its fields validated.
type Validatable interface {
	// Validate checks the validity of this types fields.
	Validate() *FieldError
}

// ValidatableTest is a standard way of testing Validatable types.
type ValidatableTest struct {
	Object Validatable
	Expect error
}

// Testable is a type derived from testing.T
type Testable interface {
	Errorf(format string, a ...interface{})
}

// Assert runs the validatae function and fails Testable.
func (vt *ValidatableTest) Assert(t Testable) {
	actual := vt.Object.Validate()
	expect := vt.Expect

	switch {
	case expect == nil && actual == nil:
		// success
	case expect == nil && actual != nil:
		t.Errorf("expected: <nil> got: %s", actual.Error())
	case expect != nil && actual == nil:
		t.Errorf("expected: %s got: <nil>", expect.Error())
	case expect.Error() != actual.Error():
		t.Errorf("expected: %s got: %s", expect.Error(), actual.Error())
	}
}
