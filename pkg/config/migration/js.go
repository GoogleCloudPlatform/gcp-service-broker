// Copyright 2019 the Service Broker Project Authors.
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

package migration

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/GoogleCloudPlatform/gcp-service-broker/utils"
	"github.com/robertkrimen/otto"
)

// JsTransform is a backing tech for building config migrations in JavaScript.
// It works for tile variables that have string values that can be replaced
// in-place.
type JsTransform struct {
	// EnvironmentVariables holds the environment vars this migratino operates on.
	EnvironmentVariables []string

	// CreateIfNotExists creates the environment varible or property
	// if it doesn't exist yet.
	CreateIfNotExists bool

	// Context returns a map of name->env-vars to look up and call the function
	// with. It's given the name of the environment variable in case the
	// other variables are based on that name.
	Context func(environmentVariable string) map[string]string

	// MigrationJs holds the JavaScript migration function. It MUST follow the
	// signature func(string) -> string. The param is the value of the environment
	// variable, and the result is the new value.
	MigrationJs string

	// TransformFuncName is the name of the JavaScript function defined in
	// MigrationJs.
	TransformFuncName string
}

func (j *JsTransform) lookupContext(variable string, env map[string]string) map[string]string {
	context := make(map[string]string)
	if j.Context == nil {
		return context
	}

	for k, v := range j.Context(variable) {
		variable, ok := env[v]
		if ok {
			context[k] = variable
		}
	}

	return context
}

func (j *JsTransform) RunGo(env map[string]string) error {
	for _, varname := range j.EnvironmentVariables {
		value, exists := env[varname]
		if !exists && !j.CreateIfNotExists {
			continue
		}

		vm := otto.New()
		if _, err := vm.Run(j.MigrationJs); err != nil {
			return fmt.Errorf("processing %s loading JS function: %v", varname, err)
		}

		jsValue, err := vm.ToValue(value)
		if err != nil {
			return fmt.Errorf("processing %s converting value to JS value: %v", varname, err)
		}

		additionalContext, err := vm.ToValue(j.lookupContext(varname, env))
		if err != nil {
			return fmt.Errorf("processing %s converting context to JS value: %v", varname, err)
		}

		result, err := vm.Call(j.TransformFuncName, nil, jsValue, additionalContext)
		if err != nil {
			return fmt.Errorf("processing %s calling: %s: %v", varname, j.TransformFuncName, err)
		}

		str, err := result.ToString()
		if err != nil {
			return fmt.Errorf("processing %s couldn't convert result to string: %v", varname, err)
		}

		env[varname] = str
	}

	return nil
}

func (*JsTransform) propName(prop string) string {
	return fmt.Sprintf("'.properties.%s'", strings.ToLower(prop))
}

func (j *JsTransform) propValue(prop string) string {
	propName := j.propName(prop)
	return fmt.Sprintf("((%s in properties['properties'])? properties['properties'][%s].value : null)", propName, propName)
}

func (j *JsTransform) jsContext(envVar string) string {
	if j.Context == nil {
		return "var context = {}; // no additional context defined"
	}

	buf := &bytes.Buffer{}
	fmt.Fprintln(buf, "var context = {")

	context := j.Context(envVar)
	sortedKeys := utils.NewStringSetFromStringMapKeys(context).ToSlice()

	for _, k := range sortedKeys {
		fmt.Fprintf(buf, " %q: %s,\n", k, j.propValue(context[k]))
	}
	fmt.Fprintln(buf, "};")
	return buf.String()
}

// ToJs converts this migration to a JavaScript suitable for running in the
// tile.
func (j *JsTransform) ToJs() string {
	buf := &bytes.Buffer{}
	fmt.Fprintln(buf, "{")
	fmt.Fprintln(buf, j.MigrationJs)
	fmt.Fprintln(buf)

	for _, v := range j.EnvironmentVariables {
		fmt.Fprintln(buf, "{")
		fmt.Fprintf(buf, "  // %s\n", v)
		fmt.Fprintf(buf, "  var prop = %s\n", j.propName(v))
		fmt.Fprintln(buf, utils.Indent(j.jsContext(v), "  "))

		if j.CreateIfNotExists {
			fmt.Fprintln(buf, "  if (!(prop in properties['properties'])) {")
			fmt.Fprintln(buf, "    properties['properties'][prop] = {'value': ''};")
			fmt.Fprintln(buf, "  }")
		}

		fmt.Fprintln(buf, "  if (prop in properties['properties']) {")
		fmt.Fprintf(buf, "    properties['properties'][prop].value = %s(properties['properties'][prop].value, context);\n", j.TransformFuncName)
		fmt.Fprintln(buf, "  }")
		fmt.Fprintln(buf, "}")
	}

	fmt.Fprintln(buf, "}")
	return buf.String()
}
