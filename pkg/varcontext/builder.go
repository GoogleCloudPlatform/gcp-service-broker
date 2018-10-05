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

	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/broker"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/varcontext/interpolation"
)

// ContextBuilder is a builder for VariableContexts.
type ContextBuilder struct {
	ErrorCollector
	context map[string]interface{}
}

// Builder creates a new ContextBuilder for constructing VariableContexts.
func Builder() *ContextBuilder {
	return &ContextBuilder{
		context: make(map[string]interface{}),
	}
}

// MergeDefaults gets the default values from the given BrokerVariables and
// if they're a string, it tries to evaluet it in the built up context.
func (builder *ContextBuilder) MergeDefaults(brokerVariables []broker.BrokerVariable) *ContextBuilder {
	for _, v := range brokerVariables {
		if v.Default == nil {
			continue
		}

		if strVal, ok := v.Default.(string); ok {
			result, err := interpolation.Eval(strVal, builder.context)
			if err != nil {
				builder.AddError(fmt.Errorf("couldn't compute the default value for %q, template: %q, %v", v.FieldName, strVal, err))
				continue
			}

			builder.context[v.FieldName] = result
		} else {
			builder.context[v.FieldName] = v.Default
		}
	}

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
	builder.AddError(json.Unmarshal(data, &out))
	builder.MergeMap(out)

	return builder
}

// Build generates a finalized VarContext based on the state of the builder.
// Exactly one of VarContext and error will be nil.
func (builder *ContextBuilder) Build() (*VarContext, error) {
	if builder.HasErrors() {
		return nil, builder.Error()
	}

	return &VarContext{context: builder.context}, nil
}
