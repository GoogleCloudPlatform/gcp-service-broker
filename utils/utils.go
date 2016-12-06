/**
# Copyright 2016 Google Inc. All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
**/

package utils

import (
	"encoding/json"
	"fmt"
	"gcp-service-broker/brokerapi/brokers/models"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"net/http"
	"os"
)

func SetGCPCredsFromEnv() error {
	rootCreds := os.Getenv(models.RootSaEnvVar)

	fo, err := os.Create(models.AppCredsFileName)
	if err != nil {
		return fmt.Errorf("Error creating file: %s", err)
	}
	_, err = fo.Write([]byte(rootCreds))
	if err != nil {
		return fmt.Errorf("Error writing to file: %s", err)
	}
	if err = fo.Close(); err != nil {
		return fmt.Errorf("Error closing file: %s", err)
	}

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("Error getting cwd: %s", err)
	}

	os.Setenv(models.AppCredsEnvVar, cwd+"/"+models.AppCredsFileName)
	return nil
}

func MapServiceIdToName() (map[string]string, error) {
	idToNameMap := make(map[string]string)
	var serviceList []models.Service
	serviceStr := os.Getenv("SERVICES")
	if err := json.Unmarshal([]byte(serviceStr), &serviceList); err != nil {
		return idToNameMap, fmt.Errorf("Error unmarshalling service list %s", err)
	}

	for _, service := range serviceList {
		idToNameMap[service.ID] = service.Name
	}
	return idToNameMap, nil
}

func GetAuthedClient() (*http.Client, error) {
	rootCreds := os.Getenv(models.RootSaEnvVar)
	conf, err := google.JWTConfigFromJSON([]byte(rootCreds), models.CloudPlatformScope)
	if err != nil {
		return nil, fmt.Errorf("Error initializing default client from credentials: %s", err)
	}
	return conf.Client(oauth2.NoContext), nil
}

func MergeStringMaps(map1 map[string]string, map2 map[string]string) map[string]string {
	combined := make(map[string]string)
	for key, val := range map1 {
		combined[key] = val
	}
	for key, val := range map2 {
		combined[key] = val
	}
	return combined
}
