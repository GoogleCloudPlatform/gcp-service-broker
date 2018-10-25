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

package varcontext

import (
	"encoding/json"
	"fmt"

	"github.com/GoogleCloudPlatform/gcp-service-broker/utils"
	"github.com/hashicorp/go-multierror"
	"github.com/spf13/cast"
)

type VarContext struct {
	errors  *multierror.Error
	context map[string]interface{}
}

func (vc *VarContext) validate(key, typeName string, validator func(interface{}) error) {
	val, ok := vc.context[key]
	if !ok {
		vc.errors = multierror.Append(vc.errors, fmt.Errorf("missing value for key %q", key))
		return
	}

	if err := validator(val); err != nil {
		vc.errors = multierror.Append(vc.errors, fmt.Errorf("value for %q must be a %s", key, typeName))
	}
}

// GetString gets a string from the context, storing an error if the key doesn't
// exist or the variable couldn't be converted to a string.
func (vc *VarContext) GetString(key string) (res string) {
	vc.validate(key, "string", func(val interface{}) (err error) {
		res, err = cast.ToStringE(val)
		return err
	})

	return
}

// GetInt gets an integer from the context, storing an error if the key doesn't
// exist or the variable couldn't be converted to an int.
func (vc *VarContext) GetInt(key string) (res int) {
	vc.validate(key, "integer", func(val interface{}) (err error) {
		res, err = cast.ToIntE(val)
		return err
	})

	return
}

// GetBool gets a boolean from the context, storing an error if the key doesn't
// exist or the variable couldn't be converted to a bool.
// Integers can behave like bools in C style, 0 is false.
// The strings "true" and "false" are also cast to their bool values.
func (vc *VarContext) GetBool(key string) (res bool) {
	vc.validate(key, "boolean", func(val interface{}) (err error) {
		res, err = cast.ToBoolE(val)
		return err
	})

	return
}

// GetStringMapString gets map[string]string from the context,
// storing an error if the key doesn't exist or the variable couldn't be cast.
func (vc *VarContext) GetStringMapString(key string) (res map[string]string) {
	vc.validate(key, "map[string]string", func(val interface{}) (err error) {
		res, err = cast.ToStringMapStringE(val)
		return err
	})

	return
}

// ToMap gets the underlying map representaiton of the variable context.
func (vc *VarContext) ToMap() map[string]interface{} {
	output := make(map[string]interface{})

	for k, v := range vc.context {
		output[k] = v
	}

	return output
}

// ToJson gets the underlying JSON representaiton of the variable context.
func (vc *VarContext) ToJson() (json.RawMessage, error) {
	return json.Marshal(vc.ToMap())
}

// Error gets the accumulated error(s) that this VarContext holds.
func (vc *VarContext) Error() error {
	if vc.errors == nil {
		return nil
	}

	vc.errors.ErrorFormat = utils.SingleLineErrorFormatter
	return vc.errors
}
