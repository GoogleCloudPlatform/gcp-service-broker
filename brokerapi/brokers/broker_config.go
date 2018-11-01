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

package brokers

import (
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/broker"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/toggles"
	"github.com/GoogleCloudPlatform/gcp-service-broker/utils"
	"golang.org/x/oauth2/jwt"
)

var enableInputValidation = toggles.Compatibility.Toggle("enable-input-validation", true, "Enables validating user input variables against JSON Schema definitions.")

type BrokerConfig struct {
	HttpConfig            *jwt.Config
	ProjectId             string
	EnableInputValidation bool
	Registry              broker.BrokerRegistry
}

func NewBrokerConfigFromEnv() (*BrokerConfig, error) {
	projectId, err := utils.GetDefaultProjectId()
	if err != nil {
		return nil, err
	}

	conf, err := utils.GetAuthedConfig()
	if err != nil {
		return nil, err
	}

	return &BrokerConfig{
		ProjectId:             projectId,
		HttpConfig:            conf,
		EnableInputValidation: enableInputValidation.IsActive(),
		Registry:              broker.DefaultRegistry,
	}, nil
}
