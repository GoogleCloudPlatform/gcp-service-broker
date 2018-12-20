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

package brokerpak

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/providers/tf"
	"github.com/GoogleCloudPlatform/gcp-service-broker/utils/stream"
)

func fakeBrokerpak() (string, error) {
	dir, err := ioutil.TempDir("", "fakepak")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(dir)

	tfSrc := filepath.Join(dir, "terraform")
	if err := stream.Copy(stream.FromString("dummy-file"), stream.ToFile(tfSrc)); err != nil {
		return "", err
	}

	exampleManifest := &Manifest{
		PackVersion: 1,
		Name:        "my-services-pack",
		Version:     "1.0.0",
		Metadata: map[string]string{
			"author": "me@example.com",
		},
		Platforms: []Platform{
			{Os: "linux", Arch: "386"},
			{Os: "linux", Arch: "amd64"},
		},
		// These resources are stubbed with a local dummy file
		TerraformResources: []TerraformResource{
			{
				Name:        "terraform",
				Version:     "0.11.9",
				Source:      tfSrc,
				UrlTemplate: tfSrc,
			},
			{
				Name:        "terraform-provider-google-beta",
				Version:     "1.19.0",
				Source:      tfSrc,
				UrlTemplate: tfSrc,
			},
		},
		ServiceDefinitions: []string{"example-service-definition.yml"},
		Parameters: []ManifestParameter{
			{Name: "TEST_PARAM", Description: "An example paramater that will be injected into Terraform's environment variables."},
		},
	}

	if err := stream.Copy(stream.FromYaml(exampleManifest), stream.ToFile(dir, manifestName)); err != nil {
		return "", err
	}

	for _, path := range exampleManifest.ServiceDefinitions {
		if err := stream.Copy(stream.FromYaml(tf.NewExampleTfServiceDefinition()), stream.ToFile(dir, path)); err != nil {
			return "", err
		}
	}

	return Pack(dir)
}

func ExampleValidate() {
	pk, err := fakeBrokerpak()
	defer os.Remove(pk)

	if err != nil {
		panic(err)
	}

	if err := Validate(pk); err != nil {
		panic(err)
	} else {
		fmt.Println("ok!")
	}

	// Output: ok!
}

func TestFinfo(t *testing.T) {
	pk, err := fakeBrokerpak()
	defer os.Remove(pk)

	if err != nil {
		t.Fatal(err)
	}

	buf := &bytes.Buffer{}
	if err := finfo(pk, buf); err != nil {
		t.Fatal(err)
	}

	// Check for "important strings" which MUST exist for this to be a valid
	// output
	importantStrings := []string{
		"Information",      // heading
		"my-services-pack", // name
		"1.0.0",            // version

		"Parameters", // heading
		"TEST_PARAM", // value

		"Dependencies",                   // heading
		"terraform",                      // dependency
		"terraform-provider-google-beta", // dependency

		"Services",                             // heading
		"00000000-0000-0000-0000-000000000000", // guid
		"example-service",                      // name

		"Contents",                               // heading
		"bin/",                                   // directory
		"definitions/",                           // directory
		"manifest.yml",                           // manifest
		"src/terraform-provider-google-beta.zip", // file
		"src/terraform.zip",                      // file

	}
	actual := string(buf.Bytes())
	for _, str := range importantStrings {
		if !strings.Contains(actual, str) {
			fmt.Errorf("Expected output to contain %s but it didn't", str)
		}
	}
}

func TestRegistryFromLocalBrokerpak(t *testing.T) {
	pk, err := fakeBrokerpak()
	defer os.Remove(pk)

	if err != nil {
		t.Fatal(err)
	}

	abs, err := filepath.Abs(pk)
	if err != nil {
		t.Fatal(err)
	}

	registry, err := registryFromLocalBrokerpak(abs)
	if err != nil {
		t.Fatal(err)
	}

	if len(registry) != 1 {
		t.Fatalf("Expected %d services but got %d", 1, len(registry))
	}

	svc, err := registry.GetServiceById("00000000-0000-0000-0000-000000000000")
	if err != nil {
		t.Fatal(err)
	}

	if svc.Name != "example-service" {
		t.Errorf("Expected exapmle-service, got %q", svc.Name)
	}
}
