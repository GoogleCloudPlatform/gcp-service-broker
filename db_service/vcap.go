package db_service

import (
	"code.cloudfoundry.org/lager"
	"encoding/json"
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

// Parse VCAP_SERVICES environment variable
func parseVcapServices(logger lager.Logger) []VcapService {
	var vcapServiceMap map[string]*json.RawMessage
	err := json.Unmarshal([]byte(os.Getenv("VCAP_SERVICES")), &vcapServiceMap)
	if err != nil {
		logger.Error("Error parsing VCAP_SERVICES environment variable", err)
	}
	var vcapServices []VcapService
	for _,v := range vcapServiceMap {
		// Debug print statement where k is the key in vcapServiceMap
		//	logger.Info("hey" + k)
		//	fmt.Printf("%s\n", *v)
		err := json.Unmarshal(*v, &vcapServices)
		if err != nil {
			logger.Error("Error parsing VCAP_SERVICES environment variable", err)
		}
	}
	return vcapServices
}