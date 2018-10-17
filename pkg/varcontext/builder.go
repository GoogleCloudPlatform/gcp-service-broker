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

	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/varcontext/interpolation"
	multierror "github.com/hashicorp/go-multierror"
	"github.com/GoogleCloudPlatform/gcp-service-broker/utils"
)

// ContextBuilder is a builder for VariableContexts.
type ContextBuilder struct {
	errors  *multierror.Error
	context map[string]interface{}
}

// Builder creates a new ContextBuilder for constructing VariableContexts.
func Builder() *ContextBuilder {
	return &ContextBuilder{
		context: make(map[string]interface{}),
	}
}

// DefaultVariable holds a value that may or may not be evaluated.
// If the value is a string then it will be evaluated.
type DefaultVariable struct {
	Name      string      `json:"name"`
	Default   interface{} `json:"default"`
	Overwrite bool        `json:"overwrite"`
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
			builder.MergeEvalResult(v.Name, strVal)
		} else {
			builder.context[v.Name] = v.Default
		}
	}

	return builder
}

// MergeEvalResult evaluates the template against the templating engine and
// merges in the value if the result is not an error.
func (builder *ContextBuilder) MergeEvalResult(key, template string) *ContextBuilder {
	result, err := interpolation.Eval(template, builder.context)
	if err != nil {
		builder.errors = multierror.Append(fmt.Errorf("couldn't compute the value for %q, template: %q, %v", key, template, err))
		return builder
	}

	builder.context[key] = result

	return builder
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

// Build generates a finalized VarContext based on the state of the builder.
// Exactly one of VarContext and error will be nil.
func (builder *ContextBuilder) Build() (*VarContext, error) {
	if builder.errors != nil {
		builder.errors.ErrorFormat = utils.LineErrorFormatter
		return nil, builder.errors
	}

	return &VarContext{context: builder.context}, nil
}
