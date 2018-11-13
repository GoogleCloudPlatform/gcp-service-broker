package tf

import (
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/broker"
)

func init() {
	// if !viper.GetBool("terraform.enable") {
	// 	return
	// }

	// TODO load definitions from a directory and instantiate them

	service, err := cloudStorage.ToService()
	if err != nil {
		panic(err)
	}
	broker.Register(service)
}
