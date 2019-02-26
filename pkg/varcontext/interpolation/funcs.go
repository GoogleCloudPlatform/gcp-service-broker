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
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"sync/atomic"
	"time"

	"github.com/hashicorp/hil"
	"github.com/hashicorp/hil/ast"
	"github.com/spf13/cast"
)

var hilStandardLibrary = createStandardLibrary()

// createStandardLibrary instantiates all the functions and associates them
// to their names in a lookup table for our standard library.
func createStandardLibrary() map[string]ast.Function {
	return map[string]ast.Function{
		"time.nano":       hilFuncTimeNano(),
		"str.truncate":    hilFuncStrTruncate(),
		"str.queryEscape": hilFuncStrQueryEscape(),
		"regexp.matches":  hilFuncRegexpMatches(),
		"counter.next":    hilFuncCounterNext(),
		"rand.base64":     hilFuncRandBase64(),
		"assert":          hilFuncAssert(),
		"json.marshal":    hilFuncJsonMarshal(),
		"map.flatten":     hilFuncMapFlatten(),
	}
}

// hilFuncTimeNano creates a function that returns the current UNIX timestamp
// in nanoseconds as a string. time.nano() -> "1538770941497"
func hilFuncTimeNano() ast.Function {
	return ast.Function{
		ArgTypes:   []ast.Type{},
		ReturnType: ast.TypeString,
		Callback: func(args []interface{}) (interface{}, error) {
			return fmt.Sprintf("%d", time.Now().UnixNano()), nil
		},
	}
}

// hilFuncStrTruncate creates a hil function that truncates a string to a given
// length. str.truncate(3, "hello") -> "hel"
func hilFuncStrTruncate() ast.Function {
	return ast.Function{
		ArgTypes:   []ast.Type{ast.TypeInt, ast.TypeString},
		ReturnType: ast.TypeString,
		Callback: func(args []interface{}) (interface{}, error) {
			maxLength := args[0].(int)
			str := args[1].(string)
			if len(str) > maxLength {
				return str[:maxLength], nil
			}

			return str, nil
		},
	}
}

// hilfuncRegexpMatches creates a hil function that checks if a string matches a given
// regular expression. regexp.matches("^d[0-9]+$", "d2)
func hilFuncRegexpMatches() ast.Function {
	return ast.Function{
		ArgTypes:   []ast.Type{ast.TypeString, ast.TypeString},
		ReturnType: ast.TypeBool,
		Callback: func(args []interface{}) (interface{}, error) {
			return regexp.MatchString(args[0].(string), args[1].(string))
		},
	}
}

// hilFuncCounterNext creates the hil function counter.next() which
// increments a counter and returns the incremented value.
// The counter is bound to the function definition, so multiple calls to
// this method will create different counters.
func hilFuncCounterNext() ast.Function {
	var counter int32

	return ast.Function{
		ArgTypes:   []ast.Type{},
		ReturnType: ast.TypeInt,
		Callback: func(args []interface{}) (interface{}, error) {
			return cast.ToIntE(atomic.AddInt32(&counter, 1))
		},
	}
}

// hilFuncRandBase64 creates n cryptographically-secure random bytes and
// converts them to Base64 rand.base64(10) -> "YWJjZGVmZ2hpag==".
func hilFuncRandBase64() ast.Function {
	return ast.Function{
		ArgTypes:   []ast.Type{ast.TypeInt},
		ReturnType: ast.TypeString,
		Callback: func(args []interface{}) (interface{}, error) {
			passwordLength := args[0].(int)
			rb := make([]byte, passwordLength)
			if _, err := rand.Read(rb); err != nil {
				return "", err
			}

			return base64.URLEncoding.EncodeToString(rb), nil
		},
	}
}

// hilFuncStrQueryEscape escapes a string suitable for embedding in a URL.
func hilFuncStrQueryEscape() ast.Function {
	return ast.Function{
		ArgTypes:   []ast.Type{ast.TypeString},
		ReturnType: ast.TypeString,
		Callback: func(args []interface{}) (interface{}, error) {
			return url.QueryEscape(args[0].(string)), nil
		},
	}
}

// hilFuncAssert throws an error with the second param if the first param is falsy.
func hilFuncAssert() ast.Function {
	return ast.Function{
		ArgTypes:   []ast.Type{ast.TypeBool, ast.TypeString},
		ReturnType: ast.TypeBool,
		Callback: func(args []interface{}) (interface{}, error) {
			condition := args[0].(bool)
			message := args[1].(string)

			if !condition {
				return false, fmt.Errorf("Assertion failed: %s", message)
			}

			return true, nil
		},
	}
}

// hilFuncJsonMarshal marshals a value as JSON.
func hilFuncJsonMarshal() ast.Function {
	return ast.Function{
		ArgTypes:   []ast.Type{ast.TypeAny},
		ReturnType: ast.TypeString,
		Callback: func(args []interface{}) (interface{}, error) {
			unwrapped, err := hilToInterface(args[0])
			if err != nil {
				return nil, err
			}

			bytes, err := json.Marshal(unwrapped)
			if err != nil {
				return nil, fmt.Errorf("couldn't convert: %v to JSON %s", args[0], err)
			}
			return string(bytes), nil
		},
	}
}

// hilFuncJsonMarshal flattens a map into a string of key/value pairs with
// given separators.
func hilFuncMapFlatten() ast.Function {
	return ast.Function{
		ArgTypes:   []ast.Type{ast.TypeString, ast.TypeString, ast.TypeMap},
		ReturnType: ast.TypeString,
		Callback: func(args []interface{}) (interface{}, error) {
			kvSep := args[0].(string)
			tupleSep := args[1].(string)
			unwrapped, err := hilToInterface(args[2])
			if err != nil {
				return nil, err
			}

			outArr := []string{}
			for k, v := range unwrapped.(map[string]interface{}) {
				outArr = append(outArr, fmt.Sprintf("%v%s%v", k, kvSep, v))
			}

			sort.Strings(outArr)

			return strings.Join(outArr, tupleSep), nil
		},
	}
}

func hilToInterface(arg interface{}) (interface{}, error) {
	// The types here cover what HIL supports.
	switch arg.(type) {
	case map[string]ast.Variable:
		out := make(map[string]interface{})
		for key, v := range arg.(map[string]ast.Variable) {
			val, verr := hilToInterface(v)
			if verr != nil {
				return nil, verr
			}

			out[key] = val
		}
		return out, nil

	case []ast.Variable:
		var out []interface{}
		for _, v := range arg.([]ast.Variable) {
			unwrapped, err := hil.VariableToInterface(v)
			if err != nil {
				return nil, err
			}
			out = append(out, unwrapped)
		}
		return out, nil

	case ast.Variable:
		return hil.VariableToInterface(arg.(ast.Variable))

	default:
		return arg, nil
	}
}
