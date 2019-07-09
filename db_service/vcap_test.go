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
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"os"
	"testing"
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
        "Password": "Ukd7QEmrfC7xMRqNmTzHCbNnmBtNceys1olOzLoSm4k",
        "PrivateKeyData": "ewogICJ0eXBlIjogInNlcnZpY2VfYWNjb3VudCIsCiAgInByb2plY3RfaWQiOiAiamxld2lzaWlpLWdzYiIsCiAgInByaXZhdGVfa2V5X2lkIjogImRkM2ZhMDU2MTY1YTFmOTFmNzJmMDk0NTJhNDA3OTUxN2IwZjNlYTEiLAogICJwcml2YXRlX2tleSI6ICItLS0tLUJFR0lOIFBSSVZBVEUgS0VZLS0tLS1cbk1JSUV2Z0lCQURBTkJna3Foa2lHOXcwQkFRRUZBQVNDQktnd2dnU2tBZ0VBQW9JQkFRRFM5dG4xd20wVUNyMldcbmRXL05HSlMveUw3Vy8xM2ROdWxYZ3hCMS9CK0Q2Q1hMZW5oSGMrYWs0d1ZmbVhGc2pnbUJLV2tLSnBZQjk3eW1cbklrRFJJNmVsZXlFbGZlYVFoV1pjVnROWEdjaHAyNHRvTkg2cENLRVdveFdtTVcwQys5OHdqR0FpZm5lOG5kdVdcbjJOT0E3K3FLVFpJaDlEd2ViKzgremN4SUZXZmpDTGN0bDFSUTRod0NDTnFHdHlQdW1pV3lBc0NhTmlFV0VKVGlcblduNDRva01qUUlZNGRQdGcvaHpjUit5cWxwWGg4a3p2czVITlI5MXg3eDJJbmlTeUpnK2lEQnlKSlZBdGowd3JcbjFKZlgzVTdpOFJLdThGTENGSkV4OEZZREd3eWRwNlZUYjkyL3RjR3NQSDl0NXVGSVc3TjUvTk5zZXVYaFBxSzNcbkgrdUQ4TkRGQWdNQkFBRUNnZ0VBSlVzNUxiNWdyUTNYQlIyZWxZbTJaZzcxV2FtTUxOcVR0b0kzYXp3V1VDbStcbllMRzJTSjlmRXlBRTU2a0hDWk0wYis1am9NVkFlSG1Va21QMHhHUUNzM2pJVzhuZGRBZjVGL0xMYXBibXZIdndcbnNZdXlKbXlkbVpSYjgrVEI2aWlmaElRVVRKVEIwd2l1OUlSQkk0YUdGa3Z2UE94aG9sblVWK3htcEFtUXMyd1pcblYxcEZSOTBtYVFWbEFLMVpGNlZWUmZTMnl4RzVtc1N4dUxkekpidjdadURYN2ZGL1NtMXlKbURjMFVzalVNWEJcbnp0WGU5dXFXZkE5RCsyamtGb1d3MGZjVnZPSE4zbDV3TFNpTTI3cXNNWksxL0YwemUwKzkwVUpyd0xnNGlxQ2dcbmxXVHFnR3FZaktOUHp6TFdJU0ZNMnpaek5yNk1Ha2dwMEIyYmhwWWVBd0tCZ1FEcEFnVUpHUGQwc0ZkYXhGVWVcbk1PbkVRZjNrYW5PRUY1eGR4UXNxdUpGNWhBMjE3YjFVQllCRjM1cGJKOURkMGsrUjdVV3dqL0Q2OUtiV3NscS9cbkFEamltY2NlalZ4V1ZUM1BTWnhySVpYRE1NbjR5eHFIaW1EaDdUSDZPa3RINFJZbU9zN2JqcEcxdzEvYnFvTlhcbjJCM1pOOEI2THRsZjdqSHdjbUREVHRqR2R3S0JnUURueC8yZ0F0ZVcxNkxaTnp1aWU3WFJjeWd2NGRZUXdoUlpcbmR6SHJLc3EvMjdjeTAwUE9wWk9BK0FtVjFwdVc1V1Y4OENtVmJsQzQwckR2cnlUaHhPODJWN0toZ3BMRWVsaUdcblMzQTQ1ell3eGluT2dYMTZKKzVaSmFYL0xLMTlyV1V5L3I0b0FMOUR4cEl2dVFpSkwyRHVZM0NzeUp6YWF0bjFcbjQ1UW5heTdsb3dLQmdFSndaLzByR0V3MmlBSUNuMzZuVmRDM1BHem9DWjR0bVZHSGdPS2lsQ0NCRGVQRk1Vb0dcbjg0ZDQ5YXR1VS9rY0ljSXJWTWErbEdrS1g1UXljUHVyVlkwUGFoNkZFa0l2dGhzb0V5amMvN1lUY0ZPM25nM3RcbjRDZ3JtU2VQZmEyMk9ibVc1U3Jub1JhaDZmQlowMisxMlBUNkY3RC9NTTVRdmY2Z3JvU2lNOStMQW9HQkFNVndcbnZHRlk2bk9KWHlTd0F6SEhOanVrWUNCaHZhdHEyRkRaMDRFalk3RUpwa1k2WnpHYUpFdWhmdkRQN3B3YzcxWDlcbmN6N2l5UXFZRjdjbE9FTEdNb3ZWS3NxZ1l3dlJ1S1Uxai9RNUtSVmxTT21ycnNxblIwZFRaZE00S05XOUprN0pcbmFBekZqaWhhOTk2RlBYczNDOWdtaHkzNGVuMG90bURhcXpMay8vOEhBb0dCQUtSQW5xeXpud08rVUZBQjBkSWFcbk1VYUdza0hJUmNvNmo4d2FqNkRlSzJqdThYRnZlZGxuOVExV2U3djRXRHBxc1BCWWhYMFdoc0ViM2pheGtRT2xcbjB4ZHprRDMrYUpBSjNFSlBpcEhFeDZscVM1aFRGdlBVM0RMTmFSbDhmQVZCdGREWVJNRmw2Y0xlMDF5b1h3UkRcbkJyU2drbHVTSU1yb1JxOGMwZVhseU53ZFxuLS0tLS1FTkQgUFJJVkFURSBLRVktLS0tLVxuIiwKICAiY2xpZW50X2VtYWlsIjogInBjZi1iaW5kaW5nLXRlc3RiaW5kQGpsZXdpc2lpaS1nc2IuaWFtLmdzZXJ2aWNlYWNjb3VudC5jb20iLAogICJjbGllbnRfaWQiOiAiMTA4ODY4NDM0NDUwOTcyMDgyNjYzIiwKICAiYXV0aF91cmkiOiAiaHR0cHM6Ly9hY2NvdW50cy5nb29nbGUuY29tL28vb2F1dGgyL2F1dGgiLAogICJ0b2tlbl91cmkiOiAiaHR0cHM6Ly9vYXV0aDIuZ29vZ2xlYXBpcy5jb20vdG9rZW4iLAogICJhdXRoX3Byb3ZpZGVyX3g1MDlfY2VydF91cmwiOiAiaHR0cHM6Ly93d3cuZ29vZ2xlYXBpcy5jb20vb2F1dGgyL3YxL2NlcnRzIiwKICAiY2xpZW50X3g1MDlfY2VydF91cmwiOiAiaHR0cHM6Ly93d3cuZ29vZ2xlYXBpcy5jb20vcm9ib3QvdjEvbWV0YWRhdGEveDUwOS9wY2YtYmluZGluZy10ZXN0YmluZCU0MGpsZXdpc2lpaS1nc2IuaWFtLmdzZXJ2aWNlYWNjb3VudC5jb20iCn0K",
        "ProjectId": "test-gsb",
        "Sha1Fingerprint": "aa3bade266136f733642ebdb4992b89eb05f83c4",
        "UniqueId": "108868434450972082663",
        "UriPrefix": "",
        "Username": "newuseraccount",
        "database_name": "service_broker",
        "host": "104.154.90.3",
        "instance_name": "pcf-sb-1-1561406852899716453",
        "last_master_operation_id": "",
        "region": "",
        "uri": "mysql://newuseraccount:Ukd7QEmrfC7xMRqNmTzHCbNnmBtNceys1olOzLoSm4k@104.154.90.3/service_broker?ssl_mode=required"
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
        "Password": "Ukd7QEmrfC7xMRqNmTzHCbNnmBtNceys1olOzLoSm4k",
        "PrivateKeyData": "ewogICJ0eXBlIjogInNlcnZpY2VfYWNjb3VudCIsCiAgInByb2plY3RfaWQiOiAiamxld2lzaWlpLWdzYiIsCiAgInByaXZhdGVfa2V5X2lkIjogImRkM2ZhMDU2MTY1YTFmOTFmNzJmMDk0NTJhNDA3OTUxN2IwZjNlYTEiLAogICJwcml2YXRlX2tleSI6ICItLS0tLUJFR0lOIFBSSVZBVEUgS0VZLS0tLS1cbk1JSUV2Z0lCQURBTkJna3Foa2lHOXcwQkFRRUZBQVNDQktnd2dnU2tBZ0VBQW9JQkFRRFM5dG4xd20wVUNyMldcbmRXL05HSlMveUw3Vy8xM2ROdWxYZ3hCMS9CK0Q2Q1hMZW5oSGMrYWs0d1ZmbVhGc2pnbUJLV2tLSnBZQjk3eW1cbklrRFJJNmVsZXlFbGZlYVFoV1pjVnROWEdjaHAyNHRvTkg2cENLRVdveFdtTVcwQys5OHdqR0FpZm5lOG5kdVdcbjJOT0E3K3FLVFpJaDlEd2ViKzgremN4SUZXZmpDTGN0bDFSUTRod0NDTnFHdHlQdW1pV3lBc0NhTmlFV0VKVGlcblduNDRva01qUUlZNGRQdGcvaHpjUit5cWxwWGg4a3p2czVITlI5MXg3eDJJbmlTeUpnK2lEQnlKSlZBdGowd3JcbjFKZlgzVTdpOFJLdThGTENGSkV4OEZZREd3eWRwNlZUYjkyL3RjR3NQSDl0NXVGSVc3TjUvTk5zZXVYaFBxSzNcbkgrdUQ4TkRGQWdNQkFBRUNnZ0VBSlVzNUxiNWdyUTNYQlIyZWxZbTJaZzcxV2FtTUxOcVR0b0kzYXp3V1VDbStcbllMRzJTSjlmRXlBRTU2a0hDWk0wYis1am9NVkFlSG1Va21QMHhHUUNzM2pJVzhuZGRBZjVGL0xMYXBibXZIdndcbnNZdXlKbXlkbVpSYjgrVEI2aWlmaElRVVRKVEIwd2l1OUlSQkk0YUdGa3Z2UE94aG9sblVWK3htcEFtUXMyd1pcblYxcEZSOTBtYVFWbEFLMVpGNlZWUmZTMnl4RzVtc1N4dUxkekpidjdadURYN2ZGL1NtMXlKbURjMFVzalVNWEJcbnp0WGU5dXFXZkE5RCsyamtGb1d3MGZjVnZPSE4zbDV3TFNpTTI3cXNNWksxL0YwemUwKzkwVUpyd0xnNGlxQ2dcbmxXVHFnR3FZaktOUHp6TFdJU0ZNMnpaek5yNk1Ha2dwMEIyYmhwWWVBd0tCZ1FEcEFnVUpHUGQwc0ZkYXhGVWVcbk1PbkVRZjNrYW5PRUY1eGR4UXNxdUpGNWhBMjE3YjFVQllCRjM1cGJKOURkMGsrUjdVV3dqL0Q2OUtiV3NscS9cbkFEamltY2NlalZ4V1ZUM1BTWnhySVpYRE1NbjR5eHFIaW1EaDdUSDZPa3RINFJZbU9zN2JqcEcxdzEvYnFvTlhcbjJCM1pOOEI2THRsZjdqSHdjbUREVHRqR2R3S0JnUURueC8yZ0F0ZVcxNkxaTnp1aWU3WFJjeWd2NGRZUXdoUlpcbmR6SHJLc3EvMjdjeTAwUE9wWk9BK0FtVjFwdVc1V1Y4OENtVmJsQzQwckR2cnlUaHhPODJWN0toZ3BMRWVsaUdcblMzQTQ1ell3eGluT2dYMTZKKzVaSmFYL0xLMTlyV1V5L3I0b0FMOUR4cEl2dVFpSkwyRHVZM0NzeUp6YWF0bjFcbjQ1UW5heTdsb3dLQmdFSndaLzByR0V3MmlBSUNuMzZuVmRDM1BHem9DWjR0bVZHSGdPS2lsQ0NCRGVQRk1Vb0dcbjg0ZDQ5YXR1VS9rY0ljSXJWTWErbEdrS1g1UXljUHVyVlkwUGFoNkZFa0l2dGhzb0V5amMvN1lUY0ZPM25nM3RcbjRDZ3JtU2VQZmEyMk9ibVc1U3Jub1JhaDZmQlowMisxMlBUNkY3RC9NTTVRdmY2Z3JvU2lNOStMQW9HQkFNVndcbnZHRlk2bk9KWHlTd0F6SEhOanVrWUNCaHZhdHEyRkRaMDRFalk3RUpwa1k2WnpHYUpFdWhmdkRQN3B3YzcxWDlcbmN6N2l5UXFZRjdjbE9FTEdNb3ZWS3NxZ1l3dlJ1S1Uxai9RNUtSVmxTT21ycnNxblIwZFRaZE00S05XOUprN0pcbmFBekZqaWhhOTk2RlBYczNDOWdtaHkzNGVuMG90bURhcXpMay8vOEhBb0dCQUtSQW5xeXpud08rVUZBQjBkSWFcbk1VYUdza0hJUmNvNmo4d2FqNkRlSzJqdThYRnZlZGxuOVExV2U3djRXRHBxc1BCWWhYMFdoc0ViM2pheGtRT2xcbjB4ZHprRDMrYUpBSjNFSlBpcEhFeDZscVM1aFRGdlBVM0RMTmFSbDhmQVZCdGREWVJNRmw2Y0xlMDF5b1h3UkRcbkJyU2drbHVTSU1yb1JxOGMwZVhseU53ZFxuLS0tLS1FTkQgUFJJVkFURSBLRVktLS0tLVxuIiwKICAiY2xpZW50X2VtYWlsIjogInBjZi1iaW5kaW5nLXRlc3RiaW5kQGpsZXdpc2lpaS1nc2IuaWFtLmdzZXJ2aWNlYWNjb3VudC5jb20iLAogICJjbGllbnRfaWQiOiAiMTA4ODY4NDM0NDUwOTcyMDgyNjYzIiwKICAiYXV0aF91cmkiOiAiaHR0cHM6Ly9hY2NvdW50cy5nb29nbGUuY29tL28vb2F1dGgyL2F1dGgiLAogICJ0b2tlbl91cmkiOiAiaHR0cHM6Ly9vYXV0aDIuZ29vZ2xlYXBpcy5jb20vdG9rZW4iLAogICJhdXRoX3Byb3ZpZGVyX3g1MDlfY2VydF91cmwiOiAiaHR0cHM6Ly93d3cuZ29vZ2xlYXBpcy5jb20vb2F1dGgyL3YxL2NlcnRzIiwKICAiY2xpZW50X3g1MDlfY2VydF91cmwiOiAiaHR0cHM6Ly93d3cuZ29vZ2xlYXBpcy5jb20vcm9ib3QvdjEvbWV0YWRhdGEveDUwOS9wY2YtYmluZGluZy10ZXN0YmluZCU0MGpsZXdpc2lpaS1nc2IuaWFtLmdzZXJ2aWNlYWNjb3VudC5jb20iCn0K",
        "ProjectId": "test-gsb",
        "Sha1Fingerprint": "aa3bade266136f733642ebdb4992b89eb05f83c4",
        "UniqueId": "108868434450972082663",
        "UriPrefix": "",
        "Username": "newuseraccount",
        "database_name": "service_broker",
        "host": "104.154.90.3",
        "instance_name": "pcf-sb-1-1561406852899716453",
        "last_master_operation_id": "",
        "region": "",
        "uri": "mysql://newuseraccount:Ukd7QEmrfC7xMRqNmTzHCbNnmBtNceys1olOzLoSm4k@104.154.90.3/service_broker?ssl_mode=required"
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
			ExpectedError: errors.New("Error parsing credentials uri field: parse mys@!ql://fefcbe8360854a18a7994b870e7b0bf5:z9z6eskdbs1rhtxt@10.0.0.20:3306/service_instance_db?reconnect=true: first path segment in URL cannot contain colon"),
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
