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
	"fmt"
	"github.com/spf13/viper"
	"os"
)

func ExampleUseVcapServices() {
	os.Setenv("VCAP_SERVICES", `{
  "p.mysql": [
    {
      "label": "p.mysql",
      "name": "my-instance",
      "plan": "db-medium",
      "provider": null,
      "syslog_drain_url": null,
      "tags": [
        "mysql"
      ],
      "credentials": {
        "hostname": "10.0.0.20",
        "jdbcUrl": "jdbc:mysql://10.0.0.20:3306/service_instance_db?user=fefcbe8360854a18a7994b870e7b0bf5\u0026password=z9z6eskdbs1rhtxt",
        "name": "service_instance_db",
        "password": "z9z6eskdbs1rhtxt",
        "port": 3306,
        "uri": "mysql://fefcbe8360854a18a7994b870e7b0bf5:z9z6eskdbs1rhtxt@10.0.0.20:3306/service_instance_db?reconnect=true",
        "username": "fefcbe8360854a18a7994b870e7b0bf5"
      },
      "volume_mounts": []
    }
  ]
}`)
	UseVcapServices()
	fmt.Println(viper.Get(dbHostProp))
	fmt.Println(viper.Get(dbUserProp))
	fmt.Println(viper.Get(dbPassProp))
	fmt.Println(viper.Get(dbNameProp))

	// Output:
	// 10.0.0.20
	// fefcbe8360854a18a7994b870e7b0bf5
	// z9z6eskdbs1rhtxt
	// service_instance_db
}
