// Copyright 2018 the Service Broker Project Authors.
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

package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/pivotal-cf/brokerapi"
)

func ExamplePropertyToEnv() {
	env := PropertyToEnv("my.property.key-value")
	fmt.Println(env)

	// Output: GSB_MY_PROPERTY_KEY_VALUE
}

func ExamplePropertyToEnvUnprefixed() {
	env := PropertyToEnvUnprefixed("my.property.key-value")
	fmt.Println(env)

	// Output: MY_PROPERTY_KEY_VALUE
}

func ExampleSetParameter() {
	// Creates an object if none is input
	out, err := SetParameter(nil, "foo", 42)
	fmt.Printf("%s, %v\n", string(out), err)

	// Replaces existing values
	out, err = SetParameter([]byte(`{"replace": "old"}`), "replace", "new")
	fmt.Printf("%s, %v\n", string(out), err)

	// Output: {"foo":42}, <nil>
	// {"replace":"new"}, <nil>
}

func ExampleUnmarshalObjectRemainder() {
	var obj struct {
		A string `json:"a_str"`
		B int
	}

	remainder, err := UnmarshalObjectRemainder([]byte(`{"a_str":"hello", "B": 33, "C": 123}`), &obj)
	fmt.Printf("%s, %v\n", string(remainder), err)

	remainder, err = UnmarshalObjectRemainder([]byte(`{"a_str":"hello", "B": 33}`), &obj)
	fmt.Printf("%s, %v\n", string(remainder), err)

	// Output: {"C":123}, <nil>
	// {}, <nil>
}

func ExampleGetDefaultProjectId() {
	serviceAccountJson := `{
	  "//": "Dummy account from https://github.com/GoogleCloudPlatform/google-cloud-java/google-cloud-clients/google-cloud-core/src/test/java/com/google/cloud/ServiceOptionsTest.java",
	  "private_key_id": "somekeyid",
	  "private_key": "-----BEGIN PRIVATE KEY-----\nMIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQC+K2hSuFpAdrJI\nnCgcDz2M7t7bjdlsadsasad+fvRSW6TjNQZ3p5LLQY1kSZRqBqylRkzteMOyHgaR\n0Pmxh3ILCND5men43j3h4eDbrhQBuxfEMalkG92sL+PNQSETY2tnvXryOvmBRwa/\nQP/9dJfIkIDJ9Fw9N4Bhhhp6mCcRpdQjV38H7JsyJ7lih/oNjECgYAt\nknddadwkwewcVxHFhcZJO+XWf6ofLUXpRwiTZakGMn8EE1uVa2LgczOjwWHGi99MFjxSer5m9\n1tCa3/KEGKiS/YL71JvjwX3mb+cewlkcmweBKZHM2JPTk0ZednFSpVZMtycjkbLa\ndYOS8V85AgMBewECggEBAKksaldajfDZDV6nGqbFjMiizAKJolr/M3OQw16K6o3/\n0S31xIe3sSlgW0+UbYlF4U8KifhManD1apVSC3csafaspP4RZUHFhtBywLO9pR5c\nr6S5aLp+gPWFyIp1pfXbWGvc5VY/v9x7ya1VEa6rXvLsKupSeWAW4tMj3eo/64ge\nsdaceaLYw52KeBYiT6+vpsnYrEkAHO1fF/LavbLLOFJmFTMxmsNaG0tuiJHgjshB\n82DpMCbXG9YcCgI/DbzuIjsdj2JC1cascSP//3PmefWysucBQe7Jryb6NQtASmnv\nCdDw/0jmZTEjpe4S1lxfHplAhHFtdgYTvyYtaLZiVVkCgYEA8eVpof2rceecw/I6\n5ng1q3Hl2usdWV/4mZMvR0fOemacLLfocX6IYxT1zA1FFJlbXSRsJMf/Qq39mOR2\nSpW+hr4jCoHeRVYLgsbggtrevGmILAlNoqCMpGZ6vDmJpq6ECV9olliDvpPgWOP+\nmYPDreFBGxWvQrADNbRt2dmGsrsCgYEAyUHqB2wvJHFqdmeBsaacewzV8x9WgmeX\ngUIi9REwXlGDW0Mz50dxpxcKCAYn65+7TCnY5O/jmL0VRxU1J2mSWyWTo1C+17L0\n3fUqjxL1pkefwecxwecvC+gFFYdJ4CQ/MHHXU81Lwl1iWdFCd2UoGddYaOF+KNeM\nHC7cmqra+JsCgYEAlUNywzq8nUg7282E+uICfCB0LfwejuymR93CtsFgb7cRd6ak\nECR8FGfCpH8ruWJINllbQfcHVCX47ndLZwqv3oVFKh6pAS/vVI4dpOepP8++7y1u\ncoOvtreXCX6XqfrWDtKIvv0vjlHBhhhp6mCcRpdQjV38H7JsyJ7lih/oNjECgYAt\nkndj5uNl5SiuVxHFhcZJO+XWf6ofLUregtevZakGMn8EE1uVa2AY7eafmoU/nZPT\n00YB0TBATdCbn/nBSuKDESkhSg9s2GEKQZG5hBmL5uCMfo09z3SfxZIhJdlerreP\nJ7gSidI12N+EZxYd4xIJh/HFDgp7RRO87f+WJkofMQKBgGTnClK1VMaCRbJZPriw\nEfeFCoOX75MxKwXs6xgrw4W//AYGGUjDt83lD6AZP6tws7gJ2IwY/qP7+lyhjEqN\nHtfPZRGFkGZsdaksdlaksd323423d+15/UvrlRSFPNj1tWQmNKkXyRDW4IG1Oa2p\nrALStNBx5Y9t0/LQnFI4w3aG\n-----END PRIVATE KEY-----\n",
	  "client_email": "someclientid@developer.gserviceaccount.com",
	  "client_id": "someclientid.apps.googleusercontent.com",
	  "type": "service_account",
	  "project_id": "my-project-123"
	}`

	os.Setenv("ROOT_SERVICE_ACCOUNT_JSON", serviceAccountJson)
	defer os.Unsetenv("ROOT_SERVICE_ACCOUNT_JSON")

	projectId, err := GetDefaultProjectId()
	fmt.Printf("%s, %v\n", projectId, err)

	// Output: my-project-123, <nil>
}

func TestExtractDefaultLabels(t *testing.T) {
	tests := map[string]struct {
		instanceId string
		details    brokerapi.ProvisionDetails
		expected   map[string]string
	}{
		"empty everything": {
			instanceId: "",
			details:    brokerapi.ProvisionDetails{},
			expected: map[string]string{
				"pcf-organization-guid": "",
				"pcf-space-guid":        "",
				"pcf-instance-id":       "",
			},
		},
		"osb 2.13": {
			instanceId: "my-instance",
			details:    brokerapi.ProvisionDetails{OrganizationGUID: "org-guid", SpaceGUID: "space-guid"},
			expected: map[string]string{
				"pcf-organization-guid": "org-guid",
				"pcf-space-guid":        "space-guid",
				"pcf-instance-id":       "my-instance",
			},
		},
		"osb future": {
			instanceId: "my-instance",
			details: brokerapi.ProvisionDetails{
				OrganizationGUID: "org-guid",
				SpaceGUID:        "space-guid",
				RawContext:       json.RawMessage(`{"organization_guid":"org-override", "space_guid":"space-override"}`),
			},
			expected: map[string]string{
				"pcf-organization-guid": "org-override",
				"pcf-space-guid":        "space-override",
				"pcf-instance-id":       "my-instance",
			},
		},
		"osb special characters": {
			instanceId: "my~instance.",
			details:    brokerapi.ProvisionDetails{},
			expected: map[string]string{
				"pcf-organization-guid": "",
				"pcf-space-guid":        "",
				"pcf-instance-id":       "my_instance_",
			},
		},
	}

	for tn, tc := range tests {
		labels := ExtractDefaultLabels(tc.instanceId, tc.details)

		if !reflect.DeepEqual(labels, tc.expected) {
			t.Errorf("Error runniung case %q, expected: %v got: %v", tn, tc.expected, labels)
		}
	}
}

func TestSplitNewlineDelimitedList(t *testing.T) {
	cases := map[string]struct {
		Input    string
		Expected []string
	}{
		"none": {
			Input:    ``,
			Expected: nil,
		},
		"single": {
			Input:    `gs://foo/bar`,
			Expected: []string{"gs://foo/bar"},
		},
		"crlf": {
			Input:    "a://foo\r\nb://bar",
			Expected: []string{"a://foo", "b://bar"},
		},
		"trim": {
			Input:    "  a://foo  \n\tb://bar  ",
			Expected: []string{"a://foo", "b://bar"},
		},
		"blank": {
			Input:    "\n\r\r\n\n\t\t\v\n\n\n\n\n \n\n\n \n \n \n\n",
			Expected: nil,
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			actual := SplitNewlineDelimitedList(tc.Input)

			if !reflect.DeepEqual(tc.Expected, actual) {
				t.Errorf("Expected: %v actual: %v", tc.Expected, actual)
			}
		})
	}
}

func ExampleIndent() {
	weirdText := "First\n\tSecond"
	out := Indent(weirdText, "  ")
	fmt.Println(out == "  First\n  \tSecond")

	// Output: true
}

func ExampleCopyStringMap() {
	m := map[string]string{"a": "one"}
	copy := CopyStringMap(m)
	m["a"] = "two"

	fmt.Println(m["a"])
	fmt.Println(copy["a"])

	// Output: two
	// one
}
