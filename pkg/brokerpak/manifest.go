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
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/validation"
	"github.com/GoogleCloudPlatform/gcp-service-broker/utils/stream"
	"github.com/GoogleCloudPlatform/gcp-service-broker/utils/ziputil"
)

const manifestName = "manifest.yml"
const HashicorpUrlTemplate = "https://releases.hashicorp.com/${name}/${version}/${name}_${version}_${os}_${arch}.zip"

type TerraformResource struct {
	Name    string `yaml:"name" validate:"required"`
	Version string `yaml:"version" validate:"required"`
	Source  string `yaml:"source" validate:"required"`

	// UrlTemplate holds a custom URL template to get the release of the given tool.
	// Paramaters available are ${name}, ${version}, ${os}, and ${arch}.
	// If non is specified HashicorpUrlTemplate is used.
	UrlTemplate string `yaml:"url_template,omitempty"`
}

func (tr *TerraformResource) Url(platform Platform) string {
	replacer := strings.NewReplacer("${name}", tr.Name, "${version}", tr.Version, "${os}", platform.Os, "${arch}", platform.Arch)
	url := tr.UrlTemplate
	if url == "" {
		url = HashicorpUrlTemplate
	}

	return replacer.Replace(url)
}

type Platform struct {
	Os   string `yaml:"os" validate:"required"`
	Arch string `yaml:"arch" validate:"required"`
}

func (p *Platform) String() string {
	return fmt.Sprintf("%s/%s", p.Os, p.Arch)
}

func (p *Platform) Equals(other Platform) bool {
	return p.String() == other.String()
}

func CurrentPlatform() Platform {
	return Platform{Os: runtime.GOOS, Arch: runtime.GOARCH}
}

// MatchesCurrent returns true if the platform matches this binary's GOOS/GOARCH combination.
func (p *Platform) MatchesCurrent() bool {
	return p.Equals(CurrentPlatform())
}

type Manifest struct {
	PackVersion int `yaml:"packversion" validate:"required,eq=1"`

	Name               string              `yaml:"name" validate:"required"`
	Version            string              `yaml:"version" validate:"required"`
	Metadata           map[string]string   `yaml:"metadata"`
	Platforms          []Platform          `yaml:"platforms" validate:"required,dive"`
	TerraformResources []TerraformResource `yaml:"terraform_binaries" validate:"required,dive"`
	ServiceDefinitions []string            `yaml:"service_definitions" validate:"required"`
}

func (m *Manifest) Validate() error {
	return validation.ValidateStruct(m)
}

func (m *Manifest) AppliesToCurrentPlatform() bool {
	for _, platform := range m.Platforms {
		if platform.MatchesCurrent() {
			return true
		}
	}

	return false
}

func (m *Manifest) Pack(base, dest string) error {
	dir, err := ioutil.TempDir("", "brokerpak")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir) // clean up

	if err := m.packSources(dir); err != nil {
		return err
	}

	if err := m.packBinaries(dir); err != nil {
		return err
	}

	if err := m.packDefinitions(dir, base); err != nil {
		return err
	}

	return ziputil.Archive(dir, dest)
}

func (m *Manifest) packSources(tmp string) error {
	for _, resource := range m.TerraformResources {
		destination := filepath.Join(tmp, "src", resource.Name+".zip")
		if err := fetch(resource.Source, destination); err != nil {
			return err
		}
	}
	return nil
}

func (m *Manifest) packBinaries(tmp string) error {
	for _, platform := range m.Platforms {
		platformPath := filepath.Join(tmp, "bin", platform.Os, platform.Arch)
		for _, resource := range m.TerraformResources {
			destination := filepath.Join(platformPath, resource.Name+".zip")
			defer os.Remove(destination)
			if err := fetch(resource.Url(platform), destination); err != nil {
				return err
			}

			if err := ziputil.Unarchive(destination, platformPath); err != nil {
				return fmt.Errorf("problem extracting %q, %v", destination, err)
			}
		}
	}

	return nil
}

func (m *Manifest) packDefinitions(tmp, base string) error {
	manifestCopy := *m

	var servicePaths []string
	for i, sd := range m.ServiceDefinitions {
		packedName := fmt.Sprintf("service-%d.yml", i)

		if err := stream.Copy(stream.FromFile(base, sd), stream.ToFile(tmp, "definitions", packedName)); err != nil {
			return err
		}

		servicePaths = append(servicePaths, "definitions/"+packedName)
	}

	manifestCopy.ServiceDefinitions = servicePaths

	return stream.Copy(stream.FromYaml(manifestCopy), stream.ToFile(tmp, manifestName))
}

// OpenManifest reads a manifest from the given file, failing if the manifest
// couldn't be decoded or read.
func OpenManifest(filename string) (*Manifest, error) {
	out := &Manifest{}
	return out, stream.Copy(stream.FromFile(filename), stream.ToYaml(out))
}

// NewExampleManifest creates a new manifest with sample values for the service broker suitable for giving a user a template to manually edit.
func NewExampleManifest() Manifest {
	return Manifest{
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
		TerraformResources: []TerraformResource{
			{
				Name:    "terraform",
				Version: "0.11.9",
				Source:  "https://github.com/hashicorp/terraform/archive/v0.11.9.zip",
			},
			{
				Name:    "terraform-provider-google-beta",
				Version: "1.19.0",
				Source:  "https://github.com/terraform-providers/terraform-provider-google/archive/v1.19.0.zip",
			},
		},
		ServiceDefinitions: []string{"example-service-definition.yml"},
	}
}
