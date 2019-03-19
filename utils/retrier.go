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

import "math/rand"
import "time"

// Retry retries the call repeatedly with a given delay and offset between 0 and jitter.
func Retry(times int, delay, jitter time.Duration, callback func() error) error {
	var err error

	for i := 0; i < times; i++ {
		err = callback()
		if err == nil {
			return nil
		}

		maxJitter := int64(jitter)
		jitterAmt := time.Duration(rand.Int63n(maxJitter))
		time.Sleep(delay + jitterAmt)
	}

	return err
}
