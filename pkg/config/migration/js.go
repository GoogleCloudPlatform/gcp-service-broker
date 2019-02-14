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

	multierror "github.com/hashicorp/go-multierror"
	"github.com/robertkrimen/otto"
)

const jsTransformContext = `
// env2prop converts an environment variable to property name
function env2prop(envVar) {return '.properties.' + envVar.toLowerCase();}

// lookupProp lookus up an environment variable as a prop and returns the value.
function lookupProp(envVar) {
	var prop = env2prop(envVar);
	return (prop in properties.properties) ? properties.properties[prop].value : null;
}

// deleteProp removes a property from the properties list.
function deleteProp(envVar) { delete properties.properties[env2prop(envVar)];}

// setProp sets a property on the properties list.
function setProp(envVar, value) {
	var prop = env2prop(envVar);

	if (!(prop in properties.properties)) {
		properties.properties[prop] = {'value': value, 'type': 'text'};
	} else {
		properties['properties'][prop].value = value;
	}
}
`

// JsTransform is a backing tech for building config migrations in JavaScript.
// It works for tile variables that have string values that can be replaced
// in-place.
type JsTransform struct {
	// EnvironmentVariables holds the environment vars this migratino operates on.
	EnvironmentVariables []string

	// MigrationJs holds the JavaScript migration function. It MUST follow the
	// signature func(string) -> string. The param is the value of the environment
	// variable, and the result is the new value.
	MigrationJs string

	// TransformFuncName is the name of the JavaScript function defined in
	// MigrationJs.
	TransformFuncName string
}

func (j *JsTransform) RunGo(env map[string]string) error {
	var runtimeErrors *multierror.Error

	vm := otto.New()

	vm.Set("lookupProp", func(call otto.FunctionCall) otto.Value {
		value, exists := env[call.Argument(0).String()]
		if !exists {
			return otto.NullValue()
		}

		jsOut, err := vm.ToValue(value)
		if err != nil {
			runtimeErrors = multierror.Append(runtimeErrors, err)
		}
		return jsOut
	})

	vm.Set("deleteProp", func(call otto.FunctionCall) otto.Value {
		delete(env, call.Argument(0).String())
		return otto.UndefinedValue()
	})

	vm.Set("setProp", func(call otto.FunctionCall) otto.Value {
		env[call.Argument(0).String()] = call.Argument(1).String()
		return otto.UndefinedValue()
	})

	if _, err := vm.Run(j.toJs(false)); err != nil {
		return fmt.Errorf("loading JS function: %v", err)
	}

	return runtimeErrors.ErrorOrNil()
}

// ToJs converts this migration to a JavaScript suitable for running in the
// tile.
func (j *JsTransform) ToJs() string {
	return j.toJs(true)
}

func (j *JsTransform) toJs(includeFunctions bool) string {
	buf := &bytes.Buffer{}
	fmt.Fprintln(buf, "{")
	if includeFunctions {
		fmt.Fprintln(buf, jsTransformContext)
	}
	fmt.Fprintln(buf, j.MigrationJs)
	fmt.Fprintln(buf)

	for _, v := range j.EnvironmentVariables {
		fmt.Fprintf(buf, "%s(%q);", j.TransformFuncName, v)
		fmt.Fprintln(buf)
	}

	fmt.Fprintln(buf, "}")
	return buf.String()
}

// ToMigration converts this transform into a migration with the given name.
func (j *JsTransform) ToMigration(name string) Migration {
	return Migration{
		Name:       name,
		TileScript: j.ToJs(),
		GoFunc:     j.RunGo,
	}
}
