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
	"log"
	"strings"

	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/models"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
)

const (
	EnvironmentVarPrefix = "gsb"
)

var (
	PropertyToEnvReplacer = strings.NewReplacer(".", "_", "-", "_")
)

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

// PrettyPrintOrExit writes a JSON serialized version of the content to stdout.
// If a failure occurs during marshaling, the error is logged along with a
// formatted version of the object and the program exits with a failure status.
func PrettyPrintOrExit(content interface{}) {
	err := prettyPrint(content)

	if err != nil {
		log.Fatalf("Could not format results: %s, results were: %+v", err, content)
	}
}

// PrettyPrintOrErr writes a JSON serialized version of the content to stdout.
// If a failure occurs during marshaling, the error is logged along with a
// formatted version of the object and the function will return the error.
func PrettyPrintOrErr(content interface{}) error {
	err := prettyPrint(content)

	if err != nil {
		log.Printf("Could not format results: %s, results were: %+v", err, content)
	}

	return err
}

func prettyPrint(content interface{}) error {
	prettyResults, err := json.MarshalIndent(content, "", "    ")
	if err == nil {
		fmt.Println(string(prettyResults))
	}

	return err
}

// PropertyToEnv converts a Viper configuration property name into an
// environment variable prefixed with EnvironmentVarPrefix
func PropertyToEnv(propertyName string) string {
	return PropertyToEnvUnprefixed(EnvironmentVarPrefix + "." + propertyName)
}

// PropertyToEnvUnprefixed converts a Viper configuration property name into an
// environment variable using PropertyToEnvReplacer
func PropertyToEnvUnprefixed(propertyName string) string {
	return PropertyToEnvReplacer.Replace(strings.ToUpper(propertyName))
}

// SetParameter sets a value on a JSON raw message and returns a modified
// version with the value set
func SetParameter(input json.RawMessage, key string, value interface{}) (json.RawMessage, error) {
	params := make(map[string]interface{})

	if input != nil && len(input) != 0 {
		err := json.Unmarshal(input, &params)
		if err != nil {
			return nil, err
		}
	}

	params[key] = value

	return json.Marshal(params)
}

// UnmarshalObjectRemaidner unmarshals an object into v and returns the
// remaining key/value pairs as a JSON string by doing a set difference.
func UnmarshalObjectRemainder(data []byte, v interface{}) ([]byte, error) {
	if err := json.Unmarshal(data, v); err != nil {
		return nil, err
	}

	encoded, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	return jsonDiff(data, encoded)
}

func jsonDiff(superset, subset json.RawMessage) ([]byte, error) {
	usedKeys := make(map[string]json.RawMessage)
	if err := json.Unmarshal(subset, &usedKeys); err != nil {
		return nil, err
	}

	allKeys := make(map[string]json.RawMessage)
	if err := json.Unmarshal(superset, &allKeys); err != nil {
		return nil, err
	}

	remainder := make(map[string]json.RawMessage)
	for key, value := range allKeys {
		if _, ok := usedKeys[key]; !ok {
			remainder[key] = value
		}
	}

	return json.Marshal(remainder)
}
