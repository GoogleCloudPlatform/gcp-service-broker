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
	"reflect"
	"regexp"

	"gopkg.in/go-playground/validator.v9"
)

var validate = validator.New()

func init() {
	validate.RegisterValidation("osbname", regexValidation(`^[a-zA-Z0-9-\.]+$`))
	validate.RegisterValidation("json", jsonValidation)
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

func regexValidation(matches string) validator.Func {
	matcher := regexp.MustCompile(matches)

	return func(field validator.FieldLevel) bool {
		return matcher.MatchString(field.Field().String())
	}
}

func jsonValidation(field validator.FieldLevel) bool {
	fl := field.Field()

	switch fl.Type().Kind() {
	case reflect.String:
		return json.Valid([]byte(fl.String()))
	default:
		return json.Valid(fl.Bytes())
	}
}
