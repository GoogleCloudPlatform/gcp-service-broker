package db_service

import (
	"code.cloudfoundry.org/lager"
	"encoding/json"
	"fmt"
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

func useVcapServices(logger lager.Logger) error {
	vcapData, vcapExists := os.LookupEnv("VCAP_SERVICES")

	if !vcapExists {
		return nil
	}
	vcapService, err := parseVcapServices(vcapData)
	if err != nil {
		return fmt.Errorf("Error parsing VCAP_SERVICES: %s", err)
	}

	u, err := url.Parse(vcapService.Credentials["uri"])
	if err != nil {
		return fmt.Errorf("Error parsing credentials uri field: %s", err)
	}

	logger.Info("Using MySQL database injected via VCAP_SERVICES environment variable")
	viper.Set(dbTypeProp, DbTypeMysql)
	viper.Set(dbHostProp, vcapService.Credentials["host"])
	viper.Set(dbUserProp, coalesce(u.User.Username(), vcapService.Credentials["Username"], vcapService.Credentials["username"]))
	viper.Set(dbPassProp, coalesce(vcapService.Credentials["Password"], vcapService.Credentials["password"]))
	viper.Set(dbNameProp, coalesce(vcapService.Credentials["database_name"], vcapService.Credentials["name"]))

	//  if database is provided by gcp service broker, retrieve the client_cert, ca_cert and client_key fields
	if contains(vcapService.Tags, "gcp") {
		viper.Set(caCertProp, vcapService.Credentials["CaCert"])
		viper.Set(clientCertProp, vcapService.Credentials["ClientCert"])
		viper.Set(clientKeyProp, vcapService.Credentials["ClientKey"])
	}

	return nil
}

// Return first non-null string in list of arguments
func coalesce(credentials ...string) string {
	for _, credential := range credentials {
		if credential != "" {
			return credential
		}
	}
	return ""
}

func parseVcapServices(vcapServicesEnv string) (VcapService, error) {
	var vcapMap map[string][]VcapService
	err := json.Unmarshal([]byte(vcapServicesEnv), &vcapMap)
	if err != nil {
		return VcapService{}, fmt.Errorf("Error unmarshaling VCAP_SERVICES: %s", err)
	}

	for _, vcapArray := range vcapMap {
		vcapService, err := findMySqlTag(vcapArray, "mysql")
		if err != nil { break }
		return vcapService, nil
	}

	return VcapService{}, fmt.Errorf("Error finding MySQL tag: %s", err)
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
