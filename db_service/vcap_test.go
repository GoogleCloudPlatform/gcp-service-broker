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
	"os"
	"testing"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
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

func TestParseVcapServices(t *testing.T) {
	cases := map[string]struct {
		VcapServiceData string
		ExpectedError   error
	}{
		"empty vcap service": {
			VcapServiceData: "",
			ExpectedError:   errors.New("Error unmarshalling VCAP_SERVICES: unexpected end of JSON input"),
		},
		"google-cloud vcap-service": {
			VcapServiceData: `{
  "google-cloudsql-mysql": [
    {
      "binding_name": "testbinding",
      "instance_name": "testinstance",
      "name": "kf-binding-tt2-mystorage",
      "label": "google-storage",
      "tags": [
        "gcp",
        "cloudsql",
        "mysql"
      ],
      "plan": "nearline",
      "credentials": {
	"CaCert": "-truncated-",
	"ClientCert": "-truncated-",
	"ClientKey": "-truncated-",
	"Email": "pcf-binding-testbind@test-gsb.iam.gserviceaccount.com",
        "Name": "pcf-binding-testbind",
        "Password": "PASSWORD",
        "PrivateKeyData": "PRIVATEKEY",
        "ProjectId": "test-gsb",
        "Sha1Fingerprint": "aa3bade266136f733642ebdb4992b89eb05f83c4",
        "UniqueId": "108868434450972082663",
        "UriPrefix": "",
        "Username": "newuseraccount",
        "database_name": "service_broker",
        "host": "127.0.0.1",
        "instance_name": "pcf-sb-1-1561406852899716453",
        "last_master_operation_id": "",
        "region": "",
        "uri": "mysql://newuseraccount:PASSWORD@127.0.0.1/service_broker?ssl_mode=required"
      }
    }
  ]
}`,
			ExpectedError: nil,
		},
		"pivotal vcap service": {
			VcapServiceData: `{
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
}
`,
			ExpectedError: nil,
		},
		"invalid vcap service - more than one mysql tag": {
			VcapServiceData: `{
  "google-cloudsql-mysql": [
    {
      "binding_name": "testbinding",
      "instance_name": "testinstance",
      "name": "kf-binding-tt2-mystorage",
      "label": "google-storage",
      "tags": [
        "gcp",
        "cloudsql",
        "mysql"
      ],
      "plan": "nearline",
      "credentials": {
	"CaCert": "-truncated-",
	"ClientCert": "-truncated-",
	"ClientKey": "-truncated-",
	"Email": "pcf-binding-testbind@test-gsb.iam.gserviceaccount.com",
        "Name": "pcf-binding-testbind",
        "Password": "PASSWORD",
        "PrivateKeyData": "PRIVATEKEY",
        "ProjectId": "test-gsb",
        "Sha1Fingerprint": "aa3bade266136f733642ebdb4992b89eb05f83c4",
        "UniqueId": "108868434450972082663",
        "UriPrefix": "",
        "Username": "newuseraccount",
        "database_name": "service_broker",
        "host": "127.0.0.1",
        "instance_name": "pcf-sb-1-1561406852899716453",
        "last_master_operation_id": "",
        "region": "",
        "uri": "mysql://newuseraccount:PASSWORD@127.0.0.1/service_broker?ssl_mode=required"
      }
    },
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
}`,
			ExpectedError: errors.New("Error finding MySQL tag: The variable VCAP_SERVICES must have one VCAP service with a tag of 'mysql'. There are currently 2 VCAP services with the tag 'mysql'."),
		},
		"invalid vcap service - zero mysql tags": {
			VcapServiceData: `{
  "p.mysql": [
    {
      "label": "p.mysql",
      "name": "my-instance",
      "plan": "db-medium",
      "provider": null,
      "syslog_drain_url": null,
      "tags": [
        "notmysql"
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
}
`,
			ExpectedError: errors.New("Error finding MySQL tag: The variable VCAP_SERVICES must have one VCAP service with a tag of 'mysql'. There are currently 0 VCAP services with the tag 'mysql'."),
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			vcapService, err := ParseVcapServices(tc.VcapServiceData)
			if err == nil {
				fmt.Printf("\n%#v\n", vcapService)
			}
			expectError(t, tc.ExpectedError, err)
		})
	}

}

func TestSetDatabaseCredentials(t *testing.T) {
	cases := map[string]struct {
		VcapService   VcapService
		ExpectedError error
	}{
		"empty vcap service": {
			VcapService:   VcapService{},
			ExpectedError: nil,
		},
		"valid gcp vcap service": {
			VcapService: VcapService{
				BindingName:  "testbinding",
				InstanceName: "testinstance",
				Name:         "kf-binding-tt2-mystorage",
				Label:        "google-storage",
				Tags:         []string{"gcp", "cloudsql", "mysql"},
				Plan:         "nearline",
				Credentials: map[string]interface{}{
					"CaCert":                   "-truncated-",
					"ClientCert":               "-truncated-",
					"ClientKey":                "-truncated-",
					"Email":                    "pcf-binding-testbind@test-gsb.iam.gserviceaccount.com",
					"Name":                     "pcf-binding-testbind",
					"Password":                 "Ukd7QEmrfC7xMRqNmTzHCbNnmBtNceys1olOzLoSm4k",
					"PrivateKeyData":           "-truncated-",
					"ProjectId":                "test-gsb",
					"Sha1Fingerprint":          "aa3bade266136f733642ebdb4992b89eb05f83c4",
					"UniqueId":                 "108868434450972082663",
					"UriPrefix":                "",
					"Username":                 "newuseraccount",
					"database_name":            "service_broker",
					"host":                     "104.154.90.3",
					"instance_name":            "pcf-sb-1-1561406852899716453",
					"last_master_operation_id": "",
					"region":                   "",
					"uri":                      "mysql://newuseraccount:Ukd7QEmrfC7xMRqNmTzHCbNnmBtNceys1olOzLoSm4k@104.154.90.3/service_broker?ssl_mode=required",
				},
			},
			ExpectedError: nil,
		},
		"valid pivotal vcap service": {
			VcapService: VcapService{
				BindingName:  "",
				InstanceName: "",
				Name:         "my-instance",
				Label:        "p.mysql",
				Tags:         []string{"mysql"},
				Plan:         "db-medium",
				Credentials: map[string]interface{}{
					"hostname": "10.0.0.20",
					"jdbcUrl":  "jdbc:mysql://10.0.0.20:3306/service_instance_db?user=fefcbe8360854a18a7994b870e7b0bf5&password=z9z6eskdbs1rhtxt",
					"name":     "service_instance_db",
					"password": "z9z6eskdbs1rhtxt",
					"port":     3306,
					"uri":      "mysql://fefcbe8360854a18a7994b870e7b0bf5:z9z6eskdbs1rhtxt@10.0.0.20:3306/service_instance_db?reconnect=true",
					"username": "fefcbe8360854a18a7994b870e7b0bf5",
				},
			},
			ExpectedError: nil,
		},
		"invalid vcap service - malformed uri": {
			VcapService: VcapService{
				BindingName:  "",
				InstanceName: "",
				Name:         "my-instance",
				Label:        "p.mysql",
				Tags:         []string{"mysql"},
				Plan:         "db-medium",
				Credentials: map[string]interface{}{
					"hostname": "10.0.0.20",
					"jdbcUrl":  "jdbc:mysql://10.0.0.20:3306/service_instance_db?user=fefcbe8360854a18a7994b870e7b0bf5&password=z9z6eskdbs1rhtxt",
					"name":     "service_instance_db",
					"password": "z9z6eskdbs1rhtxt",
					"port":     3306,
					"uri":      "mys@!ql://fefcbe8360854a18a7994b870e7b0bf5:z9z6eskdbs1rhtxt@10.0.0.20:3306/service_instance_db?reconnect=true",
					"username": "fefcbe8360854a18a7994b870e7b0bf5",
				},
			},
			ExpectedError: errors.New("Error parsing credentials uri field: parse \"mys@!ql://fefcbe8360854a18a7994b870e7b0bf5:z9z6eskdbs1rhtxt@10.0.0.20:3306/service_instance_db?reconnect=true\": first path segment in URL cannot contain colon"),
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			err := SetDatabaseCredentials(tc.VcapService)
			expectError(t, tc.ExpectedError, err)
		})
	}

}

func expectError(t *testing.T, expected, actual error) {
	t.Helper()
	expectedErr := expected != nil
	gotErr := actual != nil

	switch {
	case expectedErr && gotErr:
		if expected.Error() != actual.Error() {
			t.Fatalf("Expected: %v, got: %v", expected, actual)
		}
	case expectedErr && !gotErr:
		t.Fatalf("Expected: %v, got: %v", expected, actual)
	case !expectedErr && gotErr:
		t.Fatalf("Expected no error but got: %v", actual)
	}
}
