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

package base

import (
	"code.cloudfoundry.org/lager"
	"golang.org/x/oauth2/jwt"
)

// NewPeeredNetworkServiceBase creates a new PeeredNetworkServiceBase from the
// given settings.
func NewPeeredNetworkServiceBase(projectID string, auth *jwt.Config, logger lager.Logger) PeeredNetworkServiceBase {
	return PeeredNetworkServiceBase{
		HTTPConfig:       auth,
		DefaultProjectID: projectID,
		Logger:           logger,
	}
}

// PeeredNetworkServiceBase is a base for services that are attached to a
// project via peered network.
type PeeredNetworkServiceBase struct {
	MergedInstanceCredsMixin

	AccountManager   ServiceAccountManager
	HTTPConfig       *jwt.Config
	DefaultProjectID string
	Logger           lager.Logger
}
