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
	"archive/zip"
	"errors"
	"fmt"

	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/providers/tf"
	yaml "gopkg.in/yaml.v2"
)

// BrokerPakReader reads bundled together Terraform and service definitions.
type BrokerPakReader struct {
	contents *zip.ReadCloser
}

func (pak *BrokerPakReader) find(name string) *zip.File {
	for _, f := range pak.contents.File {
		if f.Name == name {
			return f
		}
	}

	return nil
}

func (pak *BrokerPakReader) readYaml(name string, v interface{}) error {
	fd := pak.find(name)
	if fd == nil {
		return fmt.Errorf("Couldn't find the file with the givne name %q", name)
	}

	rc, err := fd.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	decoder := yaml.NewDecoder(rc)
	if err := decoder.Decode(v); err != nil {
		return err
	}

	return nil
}

// Manifest fetches the manifest out of the package.
func (pak *BrokerPakReader) Manifest() (*Manifest, error) {
	manifest := &Manifest{}

	if err := pak.readYaml(manifestName, manifest); err != nil {
		return nil, err
	}

	return manifest, nil
}

// Services gets the list of services included in the pack.
func (pak *BrokerPakReader) Services() ([]tf.TfServiceDefinitionV1, error) {
	manifest, err := pak.Manifest()
	if err != nil {
		return nil, err
	}

	var services []tf.TfServiceDefinitionV1

	for _, serviceDefinition := range manifest.ServiceDefinitions {
		tmp := tf.TfServiceDefinitionV1{}
		if err := pak.readYaml(serviceDefinition, &tmp); err != nil {
			return nil, err
		}

		services = append(services, tmp)
	}

	return services, nil
}

func (pak *BrokerPakReader) Validate() error {
	manifest, err := pak.Manifest()
	if err != nil {
		return err
	}

	if err := manifest.Validate(); err != nil {
		return err
	}

	services, err := pak.Services()
	if err != nil {
		return err
	}

	for _, svc := range services {
		if err := svc.Validate(); err != nil {
			return err
		}
	}

	return nil
}

func (pak *BrokerPakReader) Close() {
	pak.contents.Close()
}

func (pak *BrokerPakReader) Register() error {
	mf, err := pak.Manifest()
	if err != nil {
		return err
	}

	if !mf.AppliesToCurrentPlatform() {
		return errors.New("The .brokerpack does not contain the binaries necessary to run on the current platform.")
	}

	// TODO extract binary contents to a tmp path
	// TODO register YAML files
	return nil
}

func OpenBrokerPak(pakPath string) (*BrokerPakReader, error) {
	rc, err := zip.OpenReader(pakPath)
	if err != nil {
		return nil, err
	}
	return &BrokerPakReader{contents: rc}, nil
}
