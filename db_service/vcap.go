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

package db_service

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"

	"github.com/spf13/viper"
)

type VcapService struct {
	BindingName  string                 `json:"binding_name"`  // The name assigned to the service binding by the user.
	InstanceName string                 `json:"instance_name"` // The name assigned to the service instance by the user.
	Name         string                 `json:"name"`          // The binding_name if it exists; otherwise the instance_name.
	Label        string                 `json:"label"`         // The name of the service offering.
	Tags         []string               `json:"tags"`          // An array of strings an app can use to identify a service instance.
	Plan         string                 `json:"plan"`          // The service plan selected when the service instance was created.
	Credentials  map[string]interface{} `json:"credentials"`   // The service-specific credentials needed to access the service instance.
}

func UseVcapServices() error {
	vcapData, vcapExists := os.LookupEnv("VCAP_SERVICES")

	// CloudFoundry provides an empty VCAP_SERVICES hash by default
	if !vcapExists || vcapData == "{}" {
		return nil
	}

	vcapService, err := ParseVcapServices(vcapData)
	if err != nil {
		return fmt.Errorf("Error parsing VCAP_SERVICES: %s", err)
	}

	return SetDatabaseCredentials(vcapService)
}

func SetDatabaseCredentials(vcapService VcapService) error {

	u, err := url.Parse(coalesce(vcapService.Credentials["uri"]))
	if err != nil {
		return fmt.Errorf("Error parsing credentials uri field: %s", err)
	}

	// Set up database credentials using environment variables
	viper.Set(dbTypeProp, DbTypeMysql)
	viper.Set(dbHostProp, coalesce(vcapService.Credentials["host"], vcapService.Credentials["hostname"]))
	viper.Set(dbUserProp, coalesce(u.User.Username(), vcapService.Credentials["Username"], vcapService.Credentials["username"]))
	if pw, _ := u.User.Password(); pw != "" {
		viper.Set(dbPassProp, coalesce(pw, vcapService.Credentials["Password"], vcapService.Credentials["password"]))
	} else {
		viper.Set(dbPassProp, coalesce(vcapService.Credentials["Password"], vcapService.Credentials["password"]))
	}
	viper.Set(dbNameProp, coalesce(vcapService.Credentials["database_name"], vcapService.Credentials["name"]))

	// If database is provided by gcp service broker, retrieve the client_cert, ca_cert and client_key fields
	if contains(vcapService.Tags, "gcp") {
		viper.Set(caCertProp, vcapService.Credentials["CaCert"])
		viper.Set(clientCertProp, vcapService.Credentials["ClientCert"])
		viper.Set(clientKeyProp, vcapService.Credentials["ClientKey"])
	}

	return nil
}

// Return first non-null string in list of arguments
func coalesce(credentials ...interface{}) string {
	for _, credential := range credentials {
		if credential != nil {
			switch credential.(type) {
			case int:
				return fmt.Sprintf("%d", credential)
			default:
				return fmt.Sprintf("%v", credential)
			}
		}
	}
	return ""
}

func ParseVcapServices(vcapServicesData string) (VcapService, error) {
	var vcapMap map[string][]VcapService
	err := json.Unmarshal([]byte(vcapServicesData), &vcapMap)
	if err != nil {
		return VcapService{}, fmt.Errorf("Error unmarshalling VCAP_SERVICES: %s", err)
	}

	for _, vcapArray := range vcapMap {
		vcapService, err := findMySqlTag(vcapArray, "mysql")
		if err != nil {
			return VcapService{}, fmt.Errorf("Error finding MySQL tag: %s", err)
		}
		return vcapService, nil
	}

	return VcapService{}, fmt.Errorf("Error parsing VCAP_SERVICES")
}

// whether a given string array arr contains string key
func contains(arr []string, key string) bool {
	for _, n := range arr {
		if key == n {
			return true
		}
	}
	return false
}

// return the index of the VcapService with a tag of "mysql" in the list of VcapServices, fail if we find more or fewer than 1
func findMySqlTag(vcapServices []VcapService, key string) (VcapService, error) {
	index := -1
	count := 0
	for i, vcapService := range vcapServices {
		if contains(vcapService.Tags, key) {
			count += 1
			index = i
		}
	}
	if count != 1 {
		return VcapService{}, fmt.Errorf("The variable VCAP_SERVICES must have one VCAP service with a tag of %s. There are currently %d VCAP services with the tag %s.", "'mysql'", count, "'mysql'")
	}
	return vcapServices[index], nil
}
