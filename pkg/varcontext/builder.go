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
	"reflect"

	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/validation"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/varcontext/interpolation"
	"github.com/GoogleCloudPlatform/gcp-service-broker/utils"
	multierror "github.com/hashicorp/go-multierror"
	"github.com/spf13/cast"
)

const (
	TypeObject  = "object"
	TypeBoolean = "boolean"
	TypeArray   = "array"
	TypeNumber  = "number"
	TypeString  = "string"
	TypeInteger = "integer"
)

// ContextBuilder is a builder for VariableContexts.
type ContextBuilder struct {
	errors    *multierror.Error
	context   map[string]interface{}
	constants map[string]interface{}
}

// Builder creates a new ContextBuilder for constructing VariableContexts.
func Builder() *ContextBuilder {
	return &ContextBuilder{
		context:   make(map[string]interface{}),
		constants: make(map[string]interface{}),
	}
}

// SetEvalConstants sets constants that will be available to evaluation contexts
// but not in the final output produced by the Build() call.
// These can be used to set values users can't overwrite mistakenly or maliciously.
func (builder *ContextBuilder) SetEvalConstants(constants map[string]interface{}) *ContextBuilder {
	builder.constants = constants

	return builder
}

// DefaultVariable holds a value that may or may not be evaluated.
// If the value is a string then it will be evaluated.
type DefaultVariable struct {
	Name      string      `json:"name" yaml:"name"`
	Default   interface{} `json:"default" yaml:"default"`
	Overwrite bool        `json:"overwrite" yaml:"overwrite"`
	Type      string      `json:"type" yaml:"type"`
}

var _ validation.Validatable = (*DefaultVariable)(nil)

// Validate implements validation.Validatable.
func (dv *DefaultVariable) Validate() (errs *validation.FieldError) {
	return errs.Also(
		validation.ErrIfBlank(dv.Name, "name"),
		validation.ErrIfNil(dv.Default, "default"),
		validation.ErrIfNotJSONSchemaType(dv.Type, "type"),
	)
}

// MergeDefaults gets the default values from the given BrokerVariables and
// if they're a string, it tries to evaluet it in the built up context.
func (builder *ContextBuilder) MergeDefaults(brokerVariables []DefaultVariable) *ContextBuilder {
	for _, v := range brokerVariables {
		if v.Default == nil {
			continue
		}

		if _, exists := builder.context[v.Name]; exists && !v.Overwrite {
			continue
		}

		if strVal, ok := v.Default.(string); ok {
			builder.MergeEvalResult(v.Name, strVal, v.Type)
		} else {
			builder.context[v.Name] = v.Default
		}

		if _, exists := builder.context[v.Name]; exists && !v.Overwrite {
			continue
		}
	}

	return builder
}

// MergeEvalResult evaluates the template against the templating engine and
// merges in the value if the result is not an error.
func (builder *ContextBuilder) MergeEvalResult(key, template, resultType string) *ContextBuilder {
	evaluationContext := make(map[string]interface{})
	for k, v := range builder.context {
		evaluationContext[k] = v
	}
	for k, v := range builder.constants {
		evaluationContext[k] = v
	}

	result, err := interpolation.Eval(template, evaluationContext)
	if err != nil {
		builder.errors = multierror.Append(fmt.Errorf("couldn't compute the value for %q, template: %q, %v", key, template, err))
		return builder
	}

	converted, err := castTo(result, resultType)
	if err != nil {
		builder.errors = multierror.Append(err)
		return builder
	}

	builder.context[key] = converted

	return builder
}

func toSliceE(value interface{}) ([]interface{}, error) {
	kind := reflect.TypeOf(value).Kind()
	switch kind {
	case reflect.String:
		out := []interface{}{}
		err := json.Unmarshal([]byte(value.(string)), &out)
		return out, err
	default:
		return cast.ToSliceE(value)
	}
}

func castTo(value interface{}, jsonType string) (interface{}, error) {
	switch jsonType {
	case TypeObject:
		return cast.ToStringMapE(value)
	case TypeBoolean:
		return cast.ToBoolE(value)
	case TypeArray:
		return toSliceE(value)
	case TypeNumber:
		return cast.ToFloat64E(value)
	case TypeString:
		return cast.ToStringE(value)
	case TypeInteger:
		return cast.ToIntE(value)
	case "": // for legacy compatibility
		return value, nil
	default:
		return nil, fmt.Errorf("couldn't cast %v to %s, unknown type", value, jsonType)
	}
}

// MergeMap inserts all the keys and values from the map into the context.
func (builder *ContextBuilder) MergeMap(data map[string]interface{}) *ContextBuilder {
	for k, v := range data {
		builder.context[k] = v
	}

	return builder
}

// MergeJsonObject converts the raw message to a map[string]interface{} and
// merges the values into the context. Blank RawMessages are treated like
// empty objects.
func (builder *ContextBuilder) MergeJsonObject(data json.RawMessage) *ContextBuilder {
	if len(data) == 0 {
		return builder
	}

	out := map[string]interface{}{}
	if err := json.Unmarshal(data, &out); err != nil {
		builder.errors = multierror.Append(builder.errors, err)
	}
	builder.MergeMap(out)

	return builder
}

// MergeStruct merges the given struct using its JSON field names.
func (builder *ContextBuilder) MergeStruct(data interface{}) *ContextBuilder {
	if jo, err := json.Marshal(data); err != nil {
		builder.errors = multierror.Append(builder.errors, err)
	} else {
		builder.MergeJsonObject(jo)
	}

	return builder
}

// Build generates a finalized VarContext based on the state of the builder.
// Exactly one of VarContext and error will be nil.
func (builder *ContextBuilder) Build() (*VarContext, error) {
	if builder.errors != nil {
		builder.errors.ErrorFormat = utils.SingleLineErrorFormatter
		return nil, builder.errors
	}

	return &VarContext{context: builder.context}, nil
}

// BuildMap is a shorthand of calling build then turning the returned varcontext
// into a map. Exactly one of map and error will be nil.
func (builder *ContextBuilder) BuildMap() (map[string]interface{}, error) {
	vc, err := builder.Build()
	if err != nil {
		return nil, err
	}

	return vc.ToMap(), nil
}
