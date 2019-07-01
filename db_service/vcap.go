package db_service

import (
	"code.cloudfoundry.org/lager"
	"encoding/json"
	"github.com/spf13/viper"
	"net/url"
	"os"
)

type VcapServiceMap struct {
	VcapServiceMap map[string][]VcapService
}

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
		viper.Set(caCertProp, vcapService.Credentials["CaCert"])
		viper.Set(clientCertProp, vcapService.Credentials["ClientCert"])
		viper.Set(clientKeyProp, vcapService.Credentials["ClientKey"])
		viper.Set(dbHostProp, vcapService.Credentials["host"])
		viper.Set(dbUserProp, vcapService.Credentials["Username"])
		viper.Set(dbPassProp, vcapService.Credentials["Password"])
		viper.Set(dbNameProp, vcapService.Credentials["database_name"])
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
	if len(vcapServices) > 1 {
		// TODO (hsophia): Change to logger.Error
		logger.Info("The VCAP_SERVICES environment variable may only contain one database.")
		os.Exit(1)
	}
	return vcapServices[0]
}