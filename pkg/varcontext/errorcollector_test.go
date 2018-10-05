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
	"errors"
	"fmt"
)

func ExampleErrorCollector_HasErrors() {
	collector := ErrorCollector{}

	fmt.Printf("%v\n", collector.HasErrors())
	collector.AddError(errors.New("Test"))
	fmt.Printf("%v\n", collector.HasErrors())

	// Output: false
	// true
}

func ExampleErrorCollector_AddError() {
	collector := ErrorCollector{}

	collector.AddError(nil)
	fmt.Printf("%v\n", collector.HasErrors())

	collector.AddError(errors.New("Test"))
	fmt.Printf("%v\n", collector.HasErrors())

	// Output: false
	// true
}

func ExampleErrorCollector_Error() {
	collector := ErrorCollector{}

	fmt.Println(collector.Error())

	collector.AddError(errors.New("error1"))
	collector.AddError(errors.New("error2"))

	fmt.Println(collector.Error().Error())

	// Output: <nil>
	// error1, error2
}
