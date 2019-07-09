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
			ExpectedError: errors.New("Error unmarshalling VCAP_SERVICES: unexpected end of JSON input"),
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
        "CaCert": "-----BEGIN CERTIFICATE-----\nMIIDfzCCAmegAwIBAgIBADANBgkqhkiG9w0BAQsFADB3MS0wKwYDVQQuEyRlNTJi\nMmY2MS1lYjJhLTQ2ZTQtYWVhOS1iMTgxNGRkMDk3ZTExIzAhBgNVBAMTGkdvb2ds\nZSBDbG91ZCBTUUwgU2VydmVyIENBMRQwEgYDVQQKEwtHb29nbGUsIEluYzELMAkG\nA1UEBhMCVVMwHhcNMTkwNjI0MjAxMTE4WhcNMjkwNjIxMjAxMjE4WjB3MS0wKwYD\nVQQuEyRlNTJiMmY2MS1lYjJhLTQ2ZTQtYWVhOS1iMTgxNGRkMDk3ZTExIzAhBgNV\nBAMTGkdvb2dsZSBDbG91ZCBTUUwgU2VydmVyIENBMRQwEgYDVQQKEwtHb29nbGUs\nIEluYzELMAkGA1UEBhMCVVMwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIB\nAQCxXcxUlYWr/dez+t3Azxzp75ejO1aI7Jdbp5kndAK7Myw3YxkSZbkjgHHtaVzd\n3547fCt3O3oaWDu+J/cIA7v5MHWynbxWir/ulxjPSOqd8iuETTv6bwx92nT1XQQb\nR9lbynVdx4oqP/xLsdi5Da8eknV9cc+ck/0NjOwhJ1dS3mFpj6PcKAapQmGsVM0M\n9wXiU5877tSZsY5L0pjVEKVIEFOnqGuQtB/p87LKKglnXBD7MDeV0FdzZ59spW0i\ndhwP+DVokcHJ9llDSrYDLvE5fNRXo7isxnTVpt0EsmxBW6xhucK/tPn2KCKMqXhp\n7687Q9gx+qztudtNhutxw23DAgMBAAGjFjAUMBIGA1UdEwEB/wQIMAYBAf8CAQAw\nDQYJKoZIhvcNAQELBQADggEBAKWk4jqBY8rwcHU7D4g6e4RDw3MZ/t3dcGdSEHpk\n2N3/GEZpZ1SunbxzMps+QiyggbbdmVsgV7E5E1QcE0gNatcaAKDomfrjW//2CoZC\nxQa29lOFBDP50fvYs53G2ySqB1KoU2oKDUxYnE4xumHbumrYnbWfUAgCUq3kqAsM\nyoxpLB2Z1FMMXQc3Fjci6tJdyx3S/eJ7vmiiQpM0/ruNjfZmeQAr7qwIkaW6R1FM\naQpOzBpUmDyc1eETlfgNO7Pp7DETKW2U2hcwAPkF7i25zTFBvMXRlzE1b6ysZ0Mn\nr1btp0wGz8Ivf01ehAwmecaC3UFN0zsIImul3godUoxTqIk=\n-----END CERTIFICATE-----",
        "ClientCert": "-----BEGIN CERTIFICATE-----\nMIIDfzCCAmegAwIBAgIEQIxUfzANBgkqhkiG9w0BAQsFADCBhjEtMCsGA1UELhMk\nMDQ5OGNlMTYtNGY3ZC00OTc1LWI2NTQtOGZjYWI4ZTU4NzRkMTIwMAYDVQQDEylH\nb29nbGUgQ2xvdWQgU1FMIENsaWVudCBDQSB0ZXN0YmluZGluY2VydDEUMBIGA1UE\nChMLR29vZ2xlLCBJbmMxCzAJBgNVBAYTAlVTMB4XDTE5MDYyNDIwMjUxOFoXDTI5\nMDYyMTIwMjYxOFowPDEXMBUGA1UEAxMOdGVzdGJpbmRpbmNlcnQxFDASBgNVBAoT\nC0dvb2dsZSwgSW5jMQswCQYDVQQGEwJVUzCCASIwDQYJKoZIhvcNAQEBBQADggEP\nADCCAQoCggEBAKfhU5zGair8x/GZbRay/AtAGdE0ibcKzmTV503nobOAbUV25qvD\ntHYMHBL+VUvoooor2MZOy6svqz8Ogl+8mm/bELOiZFD+9kMpE27M07+vlompddcN\n7HQol0xy10Nf7ctpnMB2GY1gVRj3BZw7u7ks/0kn5pjPqaC1aM0erOJlQGK6Q3TY\nk1YbY/xm1UYa88K7DNVlvR2pX7iWSmAIV7HRd1ojGfoLIXU3pl5LK9UU29LyNsOU\ndbg0R0nG/8O3XIyces98h4MuVqArFFcS/+JpzzPEVbylyBOa1iPD2POrQ9C/4OLH\nbhEvl87quWQl/sWNQH49/SfBk0j+KwhSNmsCAwEAAaM+MDwwOgYDVR0RBDMwMYEv\ndGVycmFmb3JtQGpsZXdpc2lpaS1nc2IuaWFtLmdzZXJ2aWNlYWNjb3VudC5jb20w\nDQYJKoZIhvcNAQELBQADggEBAEvHkoa7y2RhzaTgHM3kSw695/y0YKyo7LgFbz2O\nLeUj35sWJ4uA2xRNRoUf0Y7z2Ucgk8O4ZqdcvhQW5z/lzTtOlsBttdnibXcHB2wm\neToDu2XvqLd/AmjQNRBdTVr05TAhIYMimGR4MR5o/HrbcECuD4Q6fSkQQAhmi59x\n+NgREruyZHj+NgwKjPPls+hClkR0RWjDQ8x8pqJ1B9fXKhzABGMlOR5pgI71Ilxx\ntBY6KEyilMWR6xqOZXHkGCiVKnkVafa3nI00jIAm9o/K3tLlpTN7ROApci3/474L\nBZifVX48qBaRxpxbyK7Xl/+isIr7LiH0JmyRJuxj2cYITtY=\n-----END CERTIFICATE-----",
        "ClientKey": "-----BEGIN RSA PRIVATE KEY-----\nMIIEowIBAAKCAQEAp+FTnMZqKvzH8ZltFrL8C0AZ0TSJtwrOZNXnTeehs4BtRXbm\nq8O0dgwcEv5VS+iiiivYxk7Lqy+rPw6CX7yab9sQs6JkUP72QykTbszTv6+Wial1\n1w3sdCiXTHLXQ1/ty2mcwHYZjWBVGPcFnDu7uSz/SSfmmM+poLVozR6s4mVAYrpD\ndNiTVhtj/GbVRhrzwrsM1WW9HalfuJZKYAhXsdF3WiMZ+gshdTemXksr1RTb0vI2\nw5R1uDRHScb/w7dcjJx6z3yHgy5WoCsUVxL/4mnPM8RVvKXIE5rWI8PY86tD0L/g\n4sduES+Xzuq5ZCX+xY1Afj39J8GTSP4rCFI2awIDAQABAoIBAFtgCpl/aZQCSHXY\n84ZyXztkZWj4Nqj5acN6pc5CcEH6ef9gK0d8WwIRr0orQpPxiF66ZN/zTWncpVHJ\n/O5NAqY1T07m6cEoNTPy7I/XTr27va0qHmiyPGwxF8DVlRMn6I9Z6abb4SaRM2BG\nO7iAzrmIo17XJ+0uwn4ln2hd9O24GARgYeleWsWQ9TuuTwRZ6tS/NCfSZabNasVo\n9LkFwPpuHibuhlDBMczU2oJ8MbI0dSQwu65KpSycebQQYJ4vviq7kMCUl/Ccjz7q\n80PrqTkFFx9ZpO/Um6NYWz0nwKrUCI5cSmhXZCD6SQ2ttjvOGdQoExJIEouY8wMh\nkBINnukCgYEA87IisfK64YXP+gQ3yQTBRfmMUhPU/YTGR0FZ3x2cuhRT38l4gf/+\nHDoyI5RU9iuPwsy5JiaHuN6b1maqiYpHrDuu/+W7QgxzqFmzVvmw5e2RRVS8+zhZ\noFrF/KodFwzC1RpzQ9UOgjWoz4jmrdZMLJfQPkIoveXbP/jaW/8K2tUCgYEAsFs+\ncFxi3NL836Ql3vyyz+gDbQ+df35XCca7+DA9Bkk4mv7GCvIJusH9LvhGbfBhq8wo\nOrKMfnJeTN7gWEshpZONBZZfpJqSwCnRI3YfL1EK6v9ar/WKt9+dRrGWLjc5eeXT\nEDwttibSScNO65Xbj3ZU0NeUWQCqylSxZXhC7D8CgYEA8rXgSErwRdz4HpJE3TiX\nJhI85yJJZ5XtxNoZoFXl2o6UWrZWB2PmukZb2YPKesM4E1PCs3R8iGtt2kO2ZfYL\nHEb1LHip4EZ2ip2MOHvG67mIjfyvm6Wr3kGKHvNutZ5IDeaiFlUEdjrrPoei+FAO\n3fr1tIw/96IOk9BN6oJBVWECgYBbiQB/kXQ+6cQW0DxX4RFumB4vHUvCQPEsQdqO\nl3sVKCwZRuPECpzCMq4XEwZ7SaloYi7/SG1jtDj97TDEozpzloI7xDEgXpqM4yeK\nIGVPSeFA2AlaCzhU99vKNaKdmkxa2M8UPif7w2qinpz36nBrph+fxkVZbN845Xyu\nDh2uQQKBgBXKc+QVL1gNEMyYJyt7yjL/139M02bCrEKwCi3xTDVekGBIDuBCJgPJ\nUcAUJfwAX1ZAt1vYCTGOaAy83JiQ7BdNNYOftJONpeU0h/wzKV6jx4Bl/a8KdrTC\n7o+2nQvkm65Vr+ATAm//2mgpgKBxiPISSn/VkTqRNbWo+BLFy5WR\n-----END RSA PRIVATE KEY-----",
        "Email": "pcf-binding-testbind@jlewisiii-gsb.iam.gserviceaccount.com",
        "Name": "pcf-binding-testbind",
        "Password": "Ukd7QEmrfC7xMRqNmTzHCbNnmBtNceys1olOzLoSm4k",
        "PrivateKeyData": "ewogICJ0eXBlIjogInNlcnZpY2VfYWNjb3VudCIsCiAgInByb2plY3RfaWQiOiAiamxld2lzaWlpLWdzYiIsCiAgInByaXZhdGVfa2V5X2lkIjogImRkM2ZhMDU2MTY1YTFmOTFmNzJmMDk0NTJhNDA3OTUxN2IwZjNlYTEiLAogICJwcml2YXRlX2tleSI6ICItLS0tLUJFR0lOIFBSSVZBVEUgS0VZLS0tLS1cbk1JSUV2Z0lCQURBTkJna3Foa2lHOXcwQkFRRUZBQVNDQktnd2dnU2tBZ0VBQW9JQkFRRFM5dG4xd20wVUNyMldcbmRXL05HSlMveUw3Vy8xM2ROdWxYZ3hCMS9CK0Q2Q1hMZW5oSGMrYWs0d1ZmbVhGc2pnbUJLV2tLSnBZQjk3eW1cbklrRFJJNmVsZXlFbGZlYVFoV1pjVnROWEdjaHAyNHRvTkg2cENLRVdveFdtTVcwQys5OHdqR0FpZm5lOG5kdVdcbjJOT0E3K3FLVFpJaDlEd2ViKzgremN4SUZXZmpDTGN0bDFSUTRod0NDTnFHdHlQdW1pV3lBc0NhTmlFV0VKVGlcblduNDRva01qUUlZNGRQdGcvaHpjUit5cWxwWGg4a3p2czVITlI5MXg3eDJJbmlTeUpnK2lEQnlKSlZBdGowd3JcbjFKZlgzVTdpOFJLdThGTENGSkV4OEZZREd3eWRwNlZUYjkyL3RjR3NQSDl0NXVGSVc3TjUvTk5zZXVYaFBxSzNcbkgrdUQ4TkRGQWdNQkFBRUNnZ0VBSlVzNUxiNWdyUTNYQlIyZWxZbTJaZzcxV2FtTUxOcVR0b0kzYXp3V1VDbStcbllMRzJTSjlmRXlBRTU2a0hDWk0wYis1am9NVkFlSG1Va21QMHhHUUNzM2pJVzhuZGRBZjVGL0xMYXBibXZIdndcbnNZdXlKbXlkbVpSYjgrVEI2aWlmaElRVVRKVEIwd2l1OUlSQkk0YUdGa3Z2UE94aG9sblVWK3htcEFtUXMyd1pcblYxcEZSOTBtYVFWbEFLMVpGNlZWUmZTMnl4RzVtc1N4dUxkekpidjdadURYN2ZGL1NtMXlKbURjMFVzalVNWEJcbnp0WGU5dXFXZkE5RCsyamtGb1d3MGZjVnZPSE4zbDV3TFNpTTI3cXNNWksxL0YwemUwKzkwVUpyd0xnNGlxQ2dcbmxXVHFnR3FZaktOUHp6TFdJU0ZNMnpaek5yNk1Ha2dwMEIyYmhwWWVBd0tCZ1FEcEFnVUpHUGQwc0ZkYXhGVWVcbk1PbkVRZjNrYW5PRUY1eGR4UXNxdUpGNWhBMjE3YjFVQllCRjM1cGJKOURkMGsrUjdVV3dqL0Q2OUtiV3NscS9cbkFEamltY2NlalZ4V1ZUM1BTWnhySVpYRE1NbjR5eHFIaW1EaDdUSDZPa3RINFJZbU9zN2JqcEcxdzEvYnFvTlhcbjJCM1pOOEI2THRsZjdqSHdjbUREVHRqR2R3S0JnUURueC8yZ0F0ZVcxNkxaTnp1aWU3WFJjeWd2NGRZUXdoUlpcbmR6SHJLc3EvMjdjeTAwUE9wWk9BK0FtVjFwdVc1V1Y4OENtVmJsQzQwckR2cnlUaHhPODJWN0toZ3BMRWVsaUdcblMzQTQ1ell3eGluT2dYMTZKKzVaSmFYL0xLMTlyV1V5L3I0b0FMOUR4cEl2dVFpSkwyRHVZM0NzeUp6YWF0bjFcbjQ1UW5heTdsb3dLQmdFSndaLzByR0V3MmlBSUNuMzZuVmRDM1BHem9DWjR0bVZHSGdPS2lsQ0NCRGVQRk1Vb0dcbjg0ZDQ5YXR1VS9rY0ljSXJWTWErbEdrS1g1UXljUHVyVlkwUGFoNkZFa0l2dGhzb0V5amMvN1lUY0ZPM25nM3RcbjRDZ3JtU2VQZmEyMk9ibVc1U3Jub1JhaDZmQlowMisxMlBUNkY3RC9NTTVRdmY2Z3JvU2lNOStMQW9HQkFNVndcbnZHRlk2bk9KWHlTd0F6SEhOanVrWUNCaHZhdHEyRkRaMDRFalk3RUpwa1k2WnpHYUpFdWhmdkRQN3B3YzcxWDlcbmN6N2l5UXFZRjdjbE9FTEdNb3ZWS3NxZ1l3dlJ1S1Uxai9RNUtSVmxTT21ycnNxblIwZFRaZE00S05XOUprN0pcbmFBekZqaWhhOTk2RlBYczNDOWdtaHkzNGVuMG90bURhcXpMay8vOEhBb0dCQUtSQW5xeXpud08rVUZBQjBkSWFcbk1VYUdza0hJUmNvNmo4d2FqNkRlSzJqdThYRnZlZGxuOVExV2U3djRXRHBxc1BCWWhYMFdoc0ViM2pheGtRT2xcbjB4ZHprRDMrYUpBSjNFSlBpcEhFeDZscVM1aFRGdlBVM0RMTmFSbDhmQVZCdGREWVJNRmw2Y0xlMDF5b1h3UkRcbkJyU2drbHVTSU1yb1JxOGMwZVhseU53ZFxuLS0tLS1FTkQgUFJJVkFURSBLRVktLS0tLVxuIiwKICAiY2xpZW50X2VtYWlsIjogInBjZi1iaW5kaW5nLXRlc3RiaW5kQGpsZXdpc2lpaS1nc2IuaWFtLmdzZXJ2aWNlYWNjb3VudC5jb20iLAogICJjbGllbnRfaWQiOiAiMTA4ODY4NDM0NDUwOTcyMDgyNjYzIiwKICAiYXV0aF91cmkiOiAiaHR0cHM6Ly9hY2NvdW50cy5nb29nbGUuY29tL28vb2F1dGgyL2F1dGgiLAogICJ0b2tlbl91cmkiOiAiaHR0cHM6Ly9vYXV0aDIuZ29vZ2xlYXBpcy5jb20vdG9rZW4iLAogICJhdXRoX3Byb3ZpZGVyX3g1MDlfY2VydF91cmwiOiAiaHR0cHM6Ly93d3cuZ29vZ2xlYXBpcy5jb20vb2F1dGgyL3YxL2NlcnRzIiwKICAiY2xpZW50X3g1MDlfY2VydF91cmwiOiAiaHR0cHM6Ly93d3cuZ29vZ2xlYXBpcy5jb20vcm9ib3QvdjEvbWV0YWRhdGEveDUwOS9wY2YtYmluZGluZy10ZXN0YmluZCU0MGpsZXdpc2lpaS1nc2IuaWFtLmdzZXJ2aWNlYWNjb3VudC5jb20iCn0K",
        "ProjectId": "jlewisiii-gsb",
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
        "CaCert": "-----BEGIN CERTIFICATE-----\nMIIDfzCCAmegAwIBAgIBADANBgkqhkiG9w0BAQsFADB3MS0wKwYDVQQuEyRlNTJi\nMmY2MS1lYjJhLTQ2ZTQtYWVhOS1iMTgxNGRkMDk3ZTExIzAhBgNVBAMTGkdvb2ds\nZSBDbG91ZCBTUUwgU2VydmVyIENBMRQwEgYDVQQKEwtHb29nbGUsIEluYzELMAkG\nA1UEBhMCVVMwHhcNMTkwNjI0MjAxMTE4WhcNMjkwNjIxMjAxMjE4WjB3MS0wKwYD\nVQQuEyRlNTJiMmY2MS1lYjJhLTQ2ZTQtYWVhOS1iMTgxNGRkMDk3ZTExIzAhBgNV\nBAMTGkdvb2dsZSBDbG91ZCBTUUwgU2VydmVyIENBMRQwEgYDVQQKEwtHb29nbGUs\nIEluYzELMAkGA1UEBhMCVVMwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIB\nAQCxXcxUlYWr/dez+t3Azxzp75ejO1aI7Jdbp5kndAK7Myw3YxkSZbkjgHHtaVzd\n3547fCt3O3oaWDu+J/cIA7v5MHWynbxWir/ulxjPSOqd8iuETTv6bwx92nT1XQQb\nR9lbynVdx4oqP/xLsdi5Da8eknV9cc+ck/0NjOwhJ1dS3mFpj6PcKAapQmGsVM0M\n9wXiU5877tSZsY5L0pjVEKVIEFOnqGuQtB/p87LKKglnXBD7MDeV0FdzZ59spW0i\ndhwP+DVokcHJ9llDSrYDLvE5fNRXo7isxnTVpt0EsmxBW6xhucK/tPn2KCKMqXhp\n7687Q9gx+qztudtNhutxw23DAgMBAAGjFjAUMBIGA1UdEwEB/wQIMAYBAf8CAQAw\nDQYJKoZIhvcNAQELBQADggEBAKWk4jqBY8rwcHU7D4g6e4RDw3MZ/t3dcGdSEHpk\n2N3/GEZpZ1SunbxzMps+QiyggbbdmVsgV7E5E1QcE0gNatcaAKDomfrjW//2CoZC\nxQa29lOFBDP50fvYs53G2ySqB1KoU2oKDUxYnE4xumHbumrYnbWfUAgCUq3kqAsM\nyoxpLB2Z1FMMXQc3Fjci6tJdyx3S/eJ7vmiiQpM0/ruNjfZmeQAr7qwIkaW6R1FM\naQpOzBpUmDyc1eETlfgNO7Pp7DETKW2U2hcwAPkF7i25zTFBvMXRlzE1b6ysZ0Mn\nr1btp0wGz8Ivf01ehAwmecaC3UFN0zsIImul3godUoxTqIk=\n-----END CERTIFICATE-----",
        "ClientCert": "-----BEGIN CERTIFICATE-----\nMIIDfzCCAmegAwIBAgIEQIxUfzANBgkqhkiG9w0BAQsFADCBhjEtMCsGA1UELhMk\nMDQ5OGNlMTYtNGY3ZC00OTc1LWI2NTQtOGZjYWI4ZTU4NzRkMTIwMAYDVQQDEylH\nb29nbGUgQ2xvdWQgU1FMIENsaWVudCBDQSB0ZXN0YmluZGluY2VydDEUMBIGA1UE\nChMLR29vZ2xlLCBJbmMxCzAJBgNVBAYTAlVTMB4XDTE5MDYyNDIwMjUxOFoXDTI5\nMDYyMTIwMjYxOFowPDEXMBUGA1UEAxMOdGVzdGJpbmRpbmNlcnQxFDASBgNVBAoT\nC0dvb2dsZSwgSW5jMQswCQYDVQQGEwJVUzCCASIwDQYJKoZIhvcNAQEBBQADggEP\nADCCAQoCggEBAKfhU5zGair8x/GZbRay/AtAGdE0ibcKzmTV503nobOAbUV25qvD\ntHYMHBL+VUvoooor2MZOy6svqz8Ogl+8mm/bELOiZFD+9kMpE27M07+vlompddcN\n7HQol0xy10Nf7ctpnMB2GY1gVRj3BZw7u7ks/0kn5pjPqaC1aM0erOJlQGK6Q3TY\nk1YbY/xm1UYa88K7DNVlvR2pX7iWSmAIV7HRd1ojGfoLIXU3pl5LK9UU29LyNsOU\ndbg0R0nG/8O3XIyces98h4MuVqArFFcS/+JpzzPEVbylyBOa1iPD2POrQ9C/4OLH\nbhEvl87quWQl/sWNQH49/SfBk0j+KwhSNmsCAwEAAaM+MDwwOgYDVR0RBDMwMYEv\ndGVycmFmb3JtQGpsZXdpc2lpaS1nc2IuaWFtLmdzZXJ2aWNlYWNjb3VudC5jb20w\nDQYJKoZIhvcNAQELBQADggEBAEvHkoa7y2RhzaTgHM3kSw695/y0YKyo7LgFbz2O\nLeUj35sWJ4uA2xRNRoUf0Y7z2Ucgk8O4ZqdcvhQW5z/lzTtOlsBttdnibXcHB2wm\neToDu2XvqLd/AmjQNRBdTVr05TAhIYMimGR4MR5o/HrbcECuD4Q6fSkQQAhmi59x\n+NgREruyZHj+NgwKjPPls+hClkR0RWjDQ8x8pqJ1B9fXKhzABGMlOR5pgI71Ilxx\ntBY6KEyilMWR6xqOZXHkGCiVKnkVafa3nI00jIAm9o/K3tLlpTN7ROApci3/474L\nBZifVX48qBaRxpxbyK7Xl/+isIr7LiH0JmyRJuxj2cYITtY=\n-----END CERTIFICATE-----",
        "ClientKey": "-----BEGIN RSA PRIVATE KEY-----\nMIIEowIBAAKCAQEAp+FTnMZqKvzH8ZltFrL8C0AZ0TSJtwrOZNXnTeehs4BtRXbm\nq8O0dgwcEv5VS+iiiivYxk7Lqy+rPw6CX7yab9sQs6JkUP72QykTbszTv6+Wial1\n1w3sdCiXTHLXQ1/ty2mcwHYZjWBVGPcFnDu7uSz/SSfmmM+poLVozR6s4mVAYrpD\ndNiTVhtj/GbVRhrzwrsM1WW9HalfuJZKYAhXsdF3WiMZ+gshdTemXksr1RTb0vI2\nw5R1uDRHScb/w7dcjJx6z3yHgy5WoCsUVxL/4mnPM8RVvKXIE5rWI8PY86tD0L/g\n4sduES+Xzuq5ZCX+xY1Afj39J8GTSP4rCFI2awIDAQABAoIBAFtgCpl/aZQCSHXY\n84ZyXztkZWj4Nqj5acN6pc5CcEH6ef9gK0d8WwIRr0orQpPxiF66ZN/zTWncpVHJ\n/O5NAqY1T07m6cEoNTPy7I/XTr27va0qHmiyPGwxF8DVlRMn6I9Z6abb4SaRM2BG\nO7iAzrmIo17XJ+0uwn4ln2hd9O24GARgYeleWsWQ9TuuTwRZ6tS/NCfSZabNasVo\n9LkFwPpuHibuhlDBMczU2oJ8MbI0dSQwu65KpSycebQQYJ4vviq7kMCUl/Ccjz7q\n80PrqTkFFx9ZpO/Um6NYWz0nwKrUCI5cSmhXZCD6SQ2ttjvOGdQoExJIEouY8wMh\nkBINnukCgYEA87IisfK64YXP+gQ3yQTBRfmMUhPU/YTGR0FZ3x2cuhRT38l4gf/+\nHDoyI5RU9iuPwsy5JiaHuN6b1maqiYpHrDuu/+W7QgxzqFmzVvmw5e2RRVS8+zhZ\noFrF/KodFwzC1RpzQ9UOgjWoz4jmrdZMLJfQPkIoveXbP/jaW/8K2tUCgYEAsFs+\ncFxi3NL836Ql3vyyz+gDbQ+df35XCca7+DA9Bkk4mv7GCvIJusH9LvhGbfBhq8wo\nOrKMfnJeTN7gWEshpZONBZZfpJqSwCnRI3YfL1EK6v9ar/WKt9+dRrGWLjc5eeXT\nEDwttibSScNO65Xbj3ZU0NeUWQCqylSxZXhC7D8CgYEA8rXgSErwRdz4HpJE3TiX\nJhI85yJJZ5XtxNoZoFXl2o6UWrZWB2PmukZb2YPKesM4E1PCs3R8iGtt2kO2ZfYL\nHEb1LHip4EZ2ip2MOHvG67mIjfyvm6Wr3kGKHvNutZ5IDeaiFlUEdjrrPoei+FAO\n3fr1tIw/96IOk9BN6oJBVWECgYBbiQB/kXQ+6cQW0DxX4RFumB4vHUvCQPEsQdqO\nl3sVKCwZRuPECpzCMq4XEwZ7SaloYi7/SG1jtDj97TDEozpzloI7xDEgXpqM4yeK\nIGVPSeFA2AlaCzhU99vKNaKdmkxa2M8UPif7w2qinpz36nBrph+fxkVZbN845Xyu\nDh2uQQKBgBXKc+QVL1gNEMyYJyt7yjL/139M02bCrEKwCi3xTDVekGBIDuBCJgPJ\nUcAUJfwAX1ZAt1vYCTGOaAy83JiQ7BdNNYOftJONpeU0h/wzKV6jx4Bl/a8KdrTC\n7o+2nQvkm65Vr+ATAm//2mgpgKBxiPISSn/VkTqRNbWo+BLFy5WR\n-----END RSA PRIVATE KEY-----",
        "Email": "pcf-binding-testbind@jlewisiii-gsb.iam.gserviceaccount.com",
        "Name": "pcf-binding-testbind",
        "Password": "Ukd7QEmrfC7xMRqNmTzHCbNnmBtNceys1olOzLoSm4k",
        "PrivateKeyData": "ewogICJ0eXBlIjogInNlcnZpY2VfYWNjb3VudCIsCiAgInByb2plY3RfaWQiOiAiamxld2lzaWlpLWdzYiIsCiAgInByaXZhdGVfa2V5X2lkIjogImRkM2ZhMDU2MTY1YTFmOTFmNzJmMDk0NTJhNDA3OTUxN2IwZjNlYTEiLAogICJwcml2YXRlX2tleSI6ICItLS0tLUJFR0lOIFBSSVZBVEUgS0VZLS0tLS1cbk1JSUV2Z0lCQURBTkJna3Foa2lHOXcwQkFRRUZBQVNDQktnd2dnU2tBZ0VBQW9JQkFRRFM5dG4xd20wVUNyMldcbmRXL05HSlMveUw3Vy8xM2ROdWxYZ3hCMS9CK0Q2Q1hMZW5oSGMrYWs0d1ZmbVhGc2pnbUJLV2tLSnBZQjk3eW1cbklrRFJJNmVsZXlFbGZlYVFoV1pjVnROWEdjaHAyNHRvTkg2cENLRVdveFdtTVcwQys5OHdqR0FpZm5lOG5kdVdcbjJOT0E3K3FLVFpJaDlEd2ViKzgremN4SUZXZmpDTGN0bDFSUTRod0NDTnFHdHlQdW1pV3lBc0NhTmlFV0VKVGlcblduNDRva01qUUlZNGRQdGcvaHpjUit5cWxwWGg4a3p2czVITlI5MXg3eDJJbmlTeUpnK2lEQnlKSlZBdGowd3JcbjFKZlgzVTdpOFJLdThGTENGSkV4OEZZREd3eWRwNlZUYjkyL3RjR3NQSDl0NXVGSVc3TjUvTk5zZXVYaFBxSzNcbkgrdUQ4TkRGQWdNQkFBRUNnZ0VBSlVzNUxiNWdyUTNYQlIyZWxZbTJaZzcxV2FtTUxOcVR0b0kzYXp3V1VDbStcbllMRzJTSjlmRXlBRTU2a0hDWk0wYis1am9NVkFlSG1Va21QMHhHUUNzM2pJVzhuZGRBZjVGL0xMYXBibXZIdndcbnNZdXlKbXlkbVpSYjgrVEI2aWlmaElRVVRKVEIwd2l1OUlSQkk0YUdGa3Z2UE94aG9sblVWK3htcEFtUXMyd1pcblYxcEZSOTBtYVFWbEFLMVpGNlZWUmZTMnl4RzVtc1N4dUxkekpidjdadURYN2ZGL1NtMXlKbURjMFVzalVNWEJcbnp0WGU5dXFXZkE5RCsyamtGb1d3MGZjVnZPSE4zbDV3TFNpTTI3cXNNWksxL0YwemUwKzkwVUpyd0xnNGlxQ2dcbmxXVHFnR3FZaktOUHp6TFdJU0ZNMnpaek5yNk1Ha2dwMEIyYmhwWWVBd0tCZ1FEcEFnVUpHUGQwc0ZkYXhGVWVcbk1PbkVRZjNrYW5PRUY1eGR4UXNxdUpGNWhBMjE3YjFVQllCRjM1cGJKOURkMGsrUjdVV3dqL0Q2OUtiV3NscS9cbkFEamltY2NlalZ4V1ZUM1BTWnhySVpYRE1NbjR5eHFIaW1EaDdUSDZPa3RINFJZbU9zN2JqcEcxdzEvYnFvTlhcbjJCM1pOOEI2THRsZjdqSHdjbUREVHRqR2R3S0JnUURueC8yZ0F0ZVcxNkxaTnp1aWU3WFJjeWd2NGRZUXdoUlpcbmR6SHJLc3EvMjdjeTAwUE9wWk9BK0FtVjFwdVc1V1Y4OENtVmJsQzQwckR2cnlUaHhPODJWN0toZ3BMRWVsaUdcblMzQTQ1ell3eGluT2dYMTZKKzVaSmFYL0xLMTlyV1V5L3I0b0FMOUR4cEl2dVFpSkwyRHVZM0NzeUp6YWF0bjFcbjQ1UW5heTdsb3dLQmdFSndaLzByR0V3MmlBSUNuMzZuVmRDM1BHem9DWjR0bVZHSGdPS2lsQ0NCRGVQRk1Vb0dcbjg0ZDQ5YXR1VS9rY0ljSXJWTWErbEdrS1g1UXljUHVyVlkwUGFoNkZFa0l2dGhzb0V5amMvN1lUY0ZPM25nM3RcbjRDZ3JtU2VQZmEyMk9ibVc1U3Jub1JhaDZmQlowMisxMlBUNkY3RC9NTTVRdmY2Z3JvU2lNOStMQW9HQkFNVndcbnZHRlk2bk9KWHlTd0F6SEhOanVrWUNCaHZhdHEyRkRaMDRFalk3RUpwa1k2WnpHYUpFdWhmdkRQN3B3YzcxWDlcbmN6N2l5UXFZRjdjbE9FTEdNb3ZWS3NxZ1l3dlJ1S1Uxai9RNUtSVmxTT21ycnNxblIwZFRaZE00S05XOUprN0pcbmFBekZqaWhhOTk2RlBYczNDOWdtaHkzNGVuMG90bURhcXpMay8vOEhBb0dCQUtSQW5xeXpud08rVUZBQjBkSWFcbk1VYUdza0hJUmNvNmo4d2FqNkRlSzJqdThYRnZlZGxuOVExV2U3djRXRHBxc1BCWWhYMFdoc0ViM2pheGtRT2xcbjB4ZHprRDMrYUpBSjNFSlBpcEhFeDZscVM1aFRGdlBVM0RMTmFSbDhmQVZCdGREWVJNRmw2Y0xlMDF5b1h3UkRcbkJyU2drbHVTSU1yb1JxOGMwZVhseU53ZFxuLS0tLS1FTkQgUFJJVkFURSBLRVktLS0tLVxuIiwKICAiY2xpZW50X2VtYWlsIjogInBjZi1iaW5kaW5nLXRlc3RiaW5kQGpsZXdpc2lpaS1nc2IuaWFtLmdzZXJ2aWNlYWNjb3VudC5jb20iLAogICJjbGllbnRfaWQiOiAiMTA4ODY4NDM0NDUwOTcyMDgyNjYzIiwKICAiYXV0aF91cmkiOiAiaHR0cHM6Ly9hY2NvdW50cy5nb29nbGUuY29tL28vb2F1dGgyL2F1dGgiLAogICJ0b2tlbl91cmkiOiAiaHR0cHM6Ly9vYXV0aDIuZ29vZ2xlYXBpcy5jb20vdG9rZW4iLAogICJhdXRoX3Byb3ZpZGVyX3g1MDlfY2VydF91cmwiOiAiaHR0cHM6Ly93d3cuZ29vZ2xlYXBpcy5jb20vb2F1dGgyL3YxL2NlcnRzIiwKICAiY2xpZW50X3g1MDlfY2VydF91cmwiOiAiaHR0cHM6Ly93d3cuZ29vZ2xlYXBpcy5jb20vcm9ib3QvdjEvbWV0YWRhdGEveDUwOS9wY2YtYmluZGluZy10ZXN0YmluZCU0MGpsZXdpc2lpaS1nc2IuaWFtLmdzZXJ2aWNlYWNjb3VudC5jb20iCn0K",
        "ProjectId": "jlewisiii-gsb",
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
			_, err := ParseVcapServices(tc.VcapServiceData)
			//fmt.Println(vcapService)
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
