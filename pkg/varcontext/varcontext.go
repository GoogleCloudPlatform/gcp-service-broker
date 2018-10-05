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
	"fmt"

	"github.com/spf13/cast"
)

type VarContext struct {
	ErrorCollector

	context map[string]interface{}
}

func (vc *VarContext) validate(key, typeName string, validator func(interface{}) error) {
	val, ok := vc.context[key]
	if !ok {
		vc.AddError(fmt.Errorf("missing value for key %q", key))
		return
	}

	if err := validator(val); err != nil {
		vc.AddError(fmt.Errorf("value for %q must be a %s", key, typeName))
	}
}

// GetString gets a string from the context, storing an error if the key doesn't
// exist or the variable couldn't be converted to a string.
func (vc *VarContext) GetString(key string) string {
	vc.validate(key, "string", func(val interface{}) error {
		_, err := cast.ToStringE(val)
		return err
	})

	return cast.ToString(vc.context[key])
}

// GetInt gets an integer from the context, storing an error if the key doesn't
// exist or the variable couldn't be converted to an int.
func (vc *VarContext) GetInt(key string) int {
	vc.validate(key, "integer", func(val interface{}) error {
		_, err := cast.ToIntE(val)
		return err
	})

	return cast.ToInt(vc.context[key])
}

// ToMap gets the underlying map representaiton of the variable context.
func (vc *VarContext) ToMap() map[string]interface{} {
	output := make(map[string]interface{})

	for k, v := range vc.context {
		output[k] = v
	}

	return output
}
