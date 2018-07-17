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

	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/models"
	"github.com/spf13/viper"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
)

func MapServiceIdToName() (map[string]string, error) {
	idToNameMap := make(map[string]string)

	for _, varname := range models.ServiceNameList {

		var svc models.Service
		if err := json.Unmarshal([]byte(viper.GetString(varname)), &svc); err != nil {
			return map[string]string{}, err
		} else {
			idToNameMap[svc.ID] = svc.Name
		}
	}

	return idToNameMap, nil
}

func GetAuthedConfig() (*jwt.Config, error) {
	rootCreds := models.GetServiceAccountJson()
	conf, err := google.JWTConfigFromJSON([]byte(rootCreds), models.CloudPlatformScope)
	if err != nil {
		return nil, fmt.Errorf("Error initializing config from credentials: %s", err)
	}
	return conf, nil
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
