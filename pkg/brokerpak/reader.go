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
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/broker"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/providers/tf"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/providers/tf/wrapper"
	"github.com/GoogleCloudPlatform/gcp-service-broker/utils/stream"
	"github.com/GoogleCloudPlatform/gcp-service-broker/utils/ziputil"
)

// BrokerPakReader reads bundled together Terraform and service definitions.
type BrokerPakReader struct {
	contents *zip.ReadCloser
}

func (pak *BrokerPakReader) readYaml(name string, v interface{}) error {
	fd := ziputil.Find(&pak.contents.Reader, name)
	if fd == nil {
		return fmt.Errorf("couldn't find the file with the name %q", name)
	}

	return stream.Copy(stream.FromReadCloserError(fd.Open()), stream.ToYaml(v))
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

// Validate checks the manifest and service definitions for syntactic and
// limited semantic errors.
func (pak *BrokerPakReader) Validate() error {
	manifest, err := pak.Manifest()
	if err != nil {
		return fmt.Errorf("couldn't open brokerpak manifest: %v", err)
	}

	if err := manifest.Validate(); err != nil {
		return fmt.Errorf("couldn't validate brokerpak manifest: %v", err)
	}

	services, err := pak.Services()
	if err != nil {
		return fmt.Errorf("couldn't list services: %v", err)
	}

	for _, svc := range services {
		if err := svc.Validate(); err != nil {
			return fmt.Errorf("service %q failed validation: %v", svc.Name, err)
		}
	}

	return nil
}

// Close closes the underlying reader for the BrokerPakReader.
func (pak *BrokerPakReader) Close() error {
	return pak.contents.Close()
}

// Register extracts the binaries used by the services provided by this pak,
// then registers those services with the given registry.
func (pak *BrokerPakReader) Register(registry broker.BrokerRegistry) error {
	dir, err := ioutil.TempDir("", "brokerpak")
	if err != nil {
		return err
	}

	// extract the Terraform directory
	if err := pak.extractPlatformBins(dir); err != nil {
		return err
	}

	binPath := filepath.Join(dir, "terraform")
	executor := wrapper.CustomTerraformExecutor(binPath, dir, wrapper.DefaultExecutor)

	// register the services
	services, err := pak.Services()
	if err != nil {
		return err
	}
	for _, svc := range services {
		bs, err := svc.ToService(executor)
		if err != nil {
			return err
		}

		registry.Register(bs)
	}

	return nil
}

func (pak *BrokerPakReader) extractPlatformBins(destination string) error {
	mf, err := pak.Manifest()
	if err != nil {
		return err
	}

	curr := CurrentPlatform()
	if !mf.AppliesToCurrentPlatform() {
		return fmt.Errorf("the package %q doesn't contain binaries compatible with the current platform %q", mf.Name, curr.String())
	}

	bindir := ziputil.Join("bin", curr.Os, curr.Arch)
	return ziputil.Extract(&pak.contents.Reader, bindir, destination)
}

// Opens the file at the given path as a BrokerPakReader.
func OpenBrokerPak(pakPath string) (*BrokerPakReader, error) {
	rc, err := zip.OpenReader(pakPath)
	if err != nil {
		return nil, err
	}
	return &BrokerPakReader{contents: rc}, nil
}
