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

package utils

import "fmt"

func ExampleStringSet_Add() {
	set := NewStringSet()
	set.Add("a")
	set.Add("b")

	fmt.Println(set)
	set.Add("a")
	fmt.Println(set)

	// Output: [a b]
	// [a b]
}

func ExampleNewStringSet() {
	a := NewStringSet()
	a.Add("a")
	a.Add("b")

	b := NewStringSet("b", "a")

	fmt.Println(a.Equals(b))

	// Output: true
}

func ExampleNewStringSetFromStringMapKeys() {
	m := map[string]string{
		"a": "some a value",
		"b": "some b value",
	}

	set := NewStringSetFromStringMapKeys(m)
	fmt.Println(set)

	// Output: [a b]
}

func ExampleStringSet_ToSlice() {
	a := NewStringSet()
	a.Add("z")
	a.Add("b")

	fmt.Println(a.ToSlice())

	// Output: [b z]
}

func ExampleStringSet_IsEmpty() {
	a := NewStringSet()

	fmt.Println(a.IsEmpty())
	a.Add("a")
	fmt.Println(a.IsEmpty())

	// Output: true
	// false
}

func ExampleStringSet_Equals() {
	a := NewStringSet("a", "b")
	b := NewStringSet("a", "b", "c")
	fmt.Println(a.Equals(b))

	a.Add("c")
	fmt.Println(a.Equals(b))

	// Output: false
	// true
}

func ExampleStringSet_Contains() {
	a := NewStringSet("a", "b")
	fmt.Println(a.Contains("z"))
	fmt.Println(a.Contains("a"))

	// Output: false
	// true
}

func ExampleStringSet_Minus() {
	a := NewStringSet("a", "b")
	b := NewStringSet("b", "c")
	delta := a.Minus(b)

	fmt.Println(delta)
	// Output: [a]
}
