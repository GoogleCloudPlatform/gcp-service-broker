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

	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/validation"
	"github.com/GoogleCloudPlatform/gcp-service-broker/utils/stream"
	"github.com/GoogleCloudPlatform/gcp-service-broker/utils/ziputil"

	getter "github.com/hashicorp/go-getter"
)

const manifestName = "manifest.yml"

type Manifest struct {
	// Package metadata
	PackVersion int `yaml:"packversion" validate:"required,eq=1"`

	// User modifiable values
	Name               string              `yaml:"name" validate:"required"`
	Version            string              `yaml:"version" validate:"required"`
	Metadata           map[string]string   `yaml:"metadata"`
	Platforms          []Platform          `yaml:"platforms" validate:"required,dive"`
	TerraformResources []TerraformResource `yaml:"terraform_binaries" validate:"required,dive"`
	ServiceDefinitions []string            `yaml:"service_definitions" validate:"required"`
}

// Validate will run struct validation on the fields of this manifest.
func (m *Manifest) Validate() error {
	return validation.ValidateStruct(m)
}

// AppliesToCurrentPlatform returns true if the one of the platforms in the
// manifest match the current GOOS and GOARCH.
func (m *Manifest) AppliesToCurrentPlatform() bool {
	for _, platform := range m.Platforms {
		if platform.MatchesCurrent() {
			return true
		}
	}

	return false
}

// Pack creates a brokerpak from the manifest and definitions.
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
		if err := fetchArchive(resource.Source, destination); err != nil {
			return err
		}
	}
	return nil
}

func (m *Manifest) packBinaries(tmp string) error {
	for _, platform := range m.Platforms {
		platformPath := filepath.Join(tmp, "bin", platform.Os, platform.Arch)
		for _, resource := range m.TerraformResources {
			if err := getter.GetAny(platformPath, resource.Url(platform)); err != nil {
				return err
			}
		}
	}

	return nil
}

func (m *Manifest) packDefinitions(tmp, base string) error {
	// users can place definitions in any directory structure they like, even
	// above the current directory so we standardize their location and names
	// for the zip to avoid collisions
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
