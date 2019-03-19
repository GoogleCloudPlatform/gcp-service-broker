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
package utils

import (
	"errors"
	"fmt"
	"time"
)

func ExampleRetry_full() {
	i := 0

	err := Retry(3, 1*time.Second, func() error {
		i++
		if i != 3 {
			fmt.Printf("failing: %d\n", i)
			return errors.New("error")
		}

		fmt.Printf("succeeding: %d\n", i)
		return nil
	})

	fmt.Println(err)

	// Output: failing: 1
	// failing: 2
	// succeeding: 3
	// <nil>
}

func ExampleRetry_fail() {
	err := Retry(15, 1*time.Millisecond, func() error {
		return errors.New("error")
	})

	fmt.Println(err)

	// Output: error
}
