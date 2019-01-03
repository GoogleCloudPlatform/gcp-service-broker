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
	"fmt"

	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/broker"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/brokerpak"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/providers/builtin"
	"github.com/GoogleCloudPlatform/gcp-service-broker/utils"
	"golang.org/x/oauth2/jwt"
)

type BrokerConfig struct {
	HttpConfig *jwt.Config
	ProjectId  string
	Registry   broker.BrokerRegistry
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

	registry := builtin.BuiltinBrokerRegistry()
	if err := brokerpak.RegisterAll(registry); err != nil {
		return nil, fmt.Errorf("Error loading brokerpaks: %v", err)
	}

	return &BrokerConfig{
		ProjectId:  projectId,
		HttpConfig: conf,
		Registry:   registry,
	}, nil
}
