// Copyright the Service Broker Project Authors.
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
//
////////////////////////////////////////////////////////////////////////////////
//

package config

import (
	"encoding/json"
	"fmt"
	"gcp-service-broker/brokerapi/brokers/models"
	"gcp-service-broker/utils"
	"golang.org/x/oauth2/jwt"
	"gopkg.in/validator.v2"
	"os"
)

type BrokerConfig struct {
	Catalog    map[string]models.Service
	HttpConfig *jwt.Config
	ProjectId  string
}

func NewBrokerConfigFromEnv() (*BrokerConfig, error) {
	var err error
	bc := BrokerConfig{}
	creds, err := bc.GetCredentialsFromEnv()

	if err != nil {
		return &BrokerConfig{}, err
	}
	bc.ProjectId = creds.ProjectId
	conf, err := utils.GetAuthedConfig()
	if err != nil {
		return &BrokerConfig{}, err
	}
	bc.HttpConfig = conf
	cat, err := bc.InitCatalogFromEnv()
	if err != nil {
		return &BrokerConfig{}, err
	}
	bc.Catalog = cat
	return &bc, nil
}

// reads the service account json string from the environment variable ROOT_SERVICE_ACCOUNT_JSON, writes it to a file,
// and then exports the file location to the environment variable GOOGLE_APPLICATION_CREDENTIALS, making it visible to
// all google cloud apis
func (bc *BrokerConfig) GetCredentialsFromEnv() (models.GCPCredentials, error) {
	var err error
	g := models.GCPCredentials{}

	rootCreds := os.Getenv(models.RootSaEnvVar)
	if err = json.Unmarshal([]byte(rootCreds), &g); err != nil {
		return models.GCPCredentials{}, fmt.Errorf("Error unmarshalling service account json: %s", err)
	}

	return g, nil
}

// pulls SERVICES, PLANS, and PRECONFIGURED_PLANS environment variables to construct catalog
func (bc *BrokerConfig) InitCatalogFromEnv() (map[string]models.Service, error) {

	// set up services
	serviceMap := make(map[string]models.Service)

	for _, varname := range models.ServiceEnvVarNames {

		var svc models.Service
		if err := json.Unmarshal([]byte(os.Getenv(varname)), &svc); err != nil {
			return map[string]models.Service{}, err
		} else {
			if errs := validator.Validate(svc); errs != nil {
				return map[string]models.Service{}, errs
			} else {
				serviceMap[svc.ID] = svc
			}
		}

	}

	return serviceMap, nil
}
