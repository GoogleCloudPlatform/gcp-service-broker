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

package config

import (
	"encoding/json"
	"fmt"

	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/models"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/broker"
	"github.com/GoogleCloudPlatform/gcp-service-broker/utils"
	"golang.org/x/oauth2/jwt"
)

type BrokerConfig struct {
	Catalog    map[string]models.Service
	HttpConfig *jwt.Config
	ProjectId  string
}

func NewBrokerConfigFromEnv() (*BrokerConfig, error) {
	var err error
	bc := BrokerConfig{}
	creds, err := bc.getCredentialsFromEnv()

	if err != nil {
		return &BrokerConfig{}, err
	}
	bc.ProjectId = creds.ProjectId
	conf, err := utils.GetAuthedConfig()
	if err != nil {
		return &BrokerConfig{}, err
	}
	bc.HttpConfig = conf
	bc.Catalog, err = bc.initCatalogFromEnv()
	if err != nil {
		return &BrokerConfig{}, err
	}

	return &bc, nil
}

// reads the service account json string from the environment variable ROOT_SERVICE_ACCOUNT_JSON, writes it to a file,
// and then exports the file location to the environment variable GOOGLE_APPLICATION_CREDENTIALS, making it visible to
// all google cloud apis
func (bc *BrokerConfig) getCredentialsFromEnv() (models.GCPCredentials, error) {
	var err error
	g := models.GCPCredentials{}

	rootCreds := models.GetServiceAccountJson()
	if err = json.Unmarshal([]byte(rootCreds), &g); err != nil {
		return models.GCPCredentials{}, fmt.Errorf("Error unmarshalling service account json: %s", err)
	}

	return g, nil
}

// pulls SERVICES, PLANS, and environment variables to construct catalog
func (bc *BrokerConfig) initCatalogFromEnv() (map[string]models.Service, error) {
	serviceMap := make(map[string]models.Service)

	for _, service := range broker.GetEnabledServices() {
		entry, err := service.CatalogEntry()
		if err != nil {
			return serviceMap, err
		}
		serviceMap[entry.ID] = *entry
	}

	return serviceMap, nil
}
