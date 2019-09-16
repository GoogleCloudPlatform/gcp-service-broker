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

package interpolation

import (
	"reflect"

	"github.com/hashicorp/hil"
	"github.com/hashicorp/hil/ast"
)

// Eval evaluates the tempate string using hil https://github.com/hashicorp/hil
// with the given variables that can be accessed form the string.
func Eval(templateString string, variables map[string]interface{}) (interface{}, error) {
	tree, err := hil.Parse(templateString)
	if err != nil {
		return nil, err
	}

	varMap := make(map[string]ast.Variable)
	for vn, vv := range variables {
		converted, err := hil.InterfaceToVariable(vv)
		if err != nil {
			return nil, err
		}
		varMap[vn] = converted
	}

	config := &hil.EvalConfig{
		GlobalScope: &ast.BasicScope{
			VarMap:  varMap,
			FuncMap: hilStandardLibrary,
		},
	}

	result, err := hil.Eval(tree, config)
	if err != nil {
		return nil, err
	}

	return result.Value, err
}

// IsHILExpression returns true if the template is a HIL expression and false
// otherwise.
func IsHILExpression(template string) bool {
	tree, err := hil.Parse(template)
	if err != nil {
		return false
	}

	// Eval will error if it can't resolve a reference so we know the template is
	// a HIL expression
	result, err := hil.Eval(tree, &hil.EvalConfig{GlobalScope: &ast.BasicScope{}})
	if err != nil {
		return true
	}

	// if the template doesn't match the result value then we know something was
	// evaluated
	return !reflect.DeepEqual(template, result.Value)
}
