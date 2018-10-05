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
	"strings"
)

// ErrorCollector is a utility class that concatenates errors.
// It is useful for builder-style objects that want to keep errors tracked
// internally until a finalization call.
type ErrorCollector struct {
	errs []error
}

// HasErrors returns true if the collector has one or more errors.
func (ec *ErrorCollector) HasErrors() bool {
	return len(ec.errs) > 0
}

// Error gets a combined error message for all logged errors or nil if no errors
// were logged.
func (ec *ErrorCollector) Error() error {
	if len(ec.errs) == 0 {
		return nil
	}

	var out []string
	for _, err := range ec.errs {
		out = append(out, err.Error())
	}

	return errors.New(strings.Join(out, ", "))
}

// AddError adds an error to the internal error queue if err is not nil.
func (ec *ErrorCollector) AddError(err error) {
	if err != nil {
		ec.errs = append(ec.errs, err)
	}
}
