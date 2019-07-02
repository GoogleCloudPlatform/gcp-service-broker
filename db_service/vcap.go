package db_service

import (
	"code.cloudfoundry.org/lager"
	"encoding/json"
	"github.com/spf13/viper"
	"net/url"
	"os"
)

type VcapService struct {
	BindingName  string            `json:"binding_name"`  // The name assigned to the service binding by the user.
	InstanceName string            `json:"instance_name"` // The name assigned to the service instance by the user.
	Name         string            `json:"name"`          // The binding_name if it exists; otherwise the instance_name.
	Label        string            `json:"label"`         // The name of the service offering.
	Tags         []string          `json:"tags"`          // An array of strings an app can use to identify a service instance.
	Plan         string            `json:"plan"`          // The service plan selected when the service instance was created.
	Credentials  map[string]string `json:"credentials"`   // The service-specific credentials needed to access the service instance.
}

func useVcapServices(logger lager.Logger) {
	vcapData, vcapExists := os.LookupEnv("VCAP_SERVICES")
	if vcapExists {
		vcapService := parseVcapServices(vcapData, logger)

		u, err := url.Parse(vcapService.Credentials["uri"])
		if err != nil {
			panic(err)
		}

		viper.Set(dbPathProp, u.Path)
		viper.Set(dbTypeProp, DbTypeMysql)
		viper.Set(dbHostProp, vcapService.Credentials["host"])
		viper.Set(dbUserProp, vcapService.Credentials["Username"])
		viper.Set(dbPassProp, vcapService.Credentials["Password"])
		viper.Set(dbNameProp, vcapService.Credentials["database_name"])

		if contains(vcapService.Tags, "gcp") {
			viper.Set(caCertProp, vcapService.Credentials["CaCert"])
			viper.Set(clientCertProp, vcapService.Credentials["ClientCert"])
			viper.Set(clientKeyProp, vcapService.Credentials["ClientKey"])
		}
	}
}

// Parse VCAP_SERVICES environment variable
func parseVcapServices(vcapServicesEnv string, logger lager.Logger) VcapService {
	var vcapServiceMap map[string]*json.RawMessage
	err := json.Unmarshal([]byte(vcapServicesEnv), &vcapServiceMap)
	if err != nil {
		logger.Error("Error parsing VCAP_SERVICES environment variable", err)
	}
	var vcapServices []VcapService
	for _,v := range vcapServiceMap {
		err := json.Unmarshal(*v, &vcapServices)
		if err != nil {
			logger.Error("Error parsing VCAP_SERVICES environment variable", err)
		}
	}
	index := findMySqlTag(vcapServices, "mysql")
	if index == -1 {
		logger.Info("The VCAP_SERVICES environment variable may only contain one MySQL database.")
		os.Exit(1)
	}
	return vcapServices[index]
}

// contains tells whether a given string array arr contains string key
func contains(arr []string, key string) bool {
	for _, n := range arr {
		if key == n {
			return true
		}
	}
	return false
}

// We'll want to search the list for credentials that have a tag of "mysql", fail if we find more or fewer than 1, and use the credentials there.
func findMySqlTag(VcapServices []VcapService, key string) int {
	index := -1
	count := 0
	for i, vcapService := range VcapServices {
		if contains(vcapService.Tags, key) {
			count += 1
			index = i
		}
	}
	if count != 1 {
		return -1
	} else {
		return index
	}
}