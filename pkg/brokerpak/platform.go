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

package brokerpak

import (
	"fmt"
	"runtime"

	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/validation"
)

// Platform holds an os/architecture pair.
type Platform struct {
	Os   string `yaml:"os"`
	Arch string `yaml:"arch"`
}

var _ validation.Validatable = (*Platform)(nil)

// Validate implements validation.Validatable.
func (p Platform) Validate() (errs *validation.FieldError) {
	return errs.Also(
		validation.ErrIfBlank(p.Os, "os"),
		validation.ErrIfBlank(p.Arch, "arch"),
	)
}

// String formats the platform as an os/arch pair.
func (p Platform) String() string {
	return fmt.Sprintf("%s/%s", p.Os, p.Arch)
}

// Equals is an equality test between this platform and the other.
func (p Platform) Equals(other Platform) bool {
	return p.String() == other.String()
}

// MatchesCurrent returns true if the platform matches this binary's GOOS/GOARCH combination.
func (p Platform) MatchesCurrent() bool {
	return p.Equals(CurrentPlatform())
}

// CurrentPlatform returns the platform defined by GOOS and GOARCH.
func CurrentPlatform() Platform {
	return Platform{Os: runtime.GOOS, Arch: runtime.GOARCH}
}
