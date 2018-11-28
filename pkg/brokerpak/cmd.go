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
	"text/tabwriter"

	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/providers/tf"
	yaml "gopkg.in/yaml.v2"
)

// Init initializes a new brokerpak in the given directory with an example manifest and service definition.
func Init(directory string) error {
	manifestPath := filepath.Join(directory, manifestName)
	exampleManifest := NewExampleManifest()
	if err := writeYaml(manifestPath, exampleManifest); err != nil {
		return err
	}

	for _, path := range exampleManifest.ServiceDefinitions {
		defnPath := filepath.Join(directory, path)
		if err := writeYaml(defnPath, tf.NewExampleTfServiceDefinition()); err != nil {
			return err
		}
	}

	return nil
}

func writeYaml(path string, v interface{}) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755|os.ModeDir); err != nil {
		return err
	}

	bytes, err := yaml.Marshal(v)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, bytes, 0644)
}

func Pack(directory string) error {
	abs, err := filepath.Abs(directory)
	if err != nil {
		return err
	}
	manifestPath := filepath.Join(directory, manifestName)
	manifest, err := OpenManifest(manifestPath)
	if err != nil {
		return err
	}

	packname := filepath.Base(abs) + ".zip"
	fmt.Printf("Packing to %q\n", packname)
	return manifest.Pack(directory, packname)
}

func Info(pack string) error {
	brokerPak, err := OpenBrokerPak(pack)
	if err != nil {
		return err
	}

	mf, err := brokerPak.Manifest()
	if err != nil {
		return err
	}

	services, err := brokerPak.Services()
	if err != nil {
		return err
	}

	// Pack information
	fmt.Println("Information")
	{
		iw := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.StripEscape)
		fmt.Fprintf(iw, "format\t%d\n", mf.PackVersion)
		fmt.Fprintf(iw, "name\t%s\n", mf.Name)
		fmt.Fprintf(iw, "version\t%s\n", mf.Version)
		fmt.Fprintln(iw, "platforms")
		for _, arch := range mf.Platforms {
			fmt.Fprintf(iw, "\t%s\n", arch.String())
		}
		fmt.Fprintln(iw, "metadata")
		for k, v := range mf.Metadata {
			fmt.Fprintf(iw, "\t%s\t%s\n", k, v)
		}

		iw.Flush()
		fmt.Println()
	}

	{
		fmt.Println("Dependencies")
		rw := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.StripEscape)
		fmt.Fprintln(rw, "NAME\tVERSION\tSOURCE")
		for _, resource := range mf.TerraformResources {
			fmt.Fprintf(rw, "%s\t%s\t%s\n", resource.Name, resource.Version, resource.Source)
		}
		rw.Flush()
		fmt.Println()
	}

	{
		fmt.Println("Services")
		sw := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.StripEscape)
		fmt.Fprintln(sw, "ID\tNAME\tDESCRIPTION\tPLANS")
		for _, svc := range services {
			fmt.Fprintf(sw, "%s\t%s\t%s\t%d\n", svc.Id, svc.Name, svc.Description, len(svc.Plans))
		}
		sw.Flush()
		fmt.Println()
	}

	{
		fmt.Println("Contents")
		sw := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.StripEscape)
		fmt.Fprintln(sw, "MODE\tSIZE\tNAME")
		for _, fd := range brokerPak.contents.File {
			fmt.Fprintf(sw, "%s\t%d\t%s\n", fd.Mode().String(), fd.UncompressedSize, fd.Name)
		}
		sw.Flush()
		fmt.Println()
	}
	return nil
}

func Validate(pack string) error {
	brokerPak, err := OpenBrokerPak(pack)
	if err != nil {
		return err
	}

	return brokerPak.Validate()
}
