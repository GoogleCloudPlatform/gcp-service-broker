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
	"os"
	"path/filepath"
	"text/tabwriter"

	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/broker"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/client"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/generator"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/providers/tf"
	"github.com/GoogleCloudPlatform/gcp-service-broker/utils/stream"
	"github.com/GoogleCloudPlatform/gcp-service-broker/utils/ziputil"
	"github.com/spf13/viper"
)

// Init initializes a new brokerpak in the given directory with an example manifest and service definition.
func Init(directory string) error {
	exampleManifest := NewExampleManifest()
	if err := stream.Copy(stream.FromYaml(exampleManifest), stream.ToFile(directory, manifestName)); err != nil {
		return err
	}

	for _, path := range exampleManifest.ServiceDefinitions {
		if err := stream.Copy(stream.FromYaml(tf.NewExampleTfServiceDefinition()), stream.ToFile(directory, path)); err != nil {
			return err
		}
	}

	return nil
}

// Pack creates a new brokerpak from the given directory which MUST contain a
// manifest.yml file. If the pack was successful, the returned string will be
// the path to the created brokerpak.
func Pack(directory string) (string, error) {
	abs, err := filepath.Abs(directory)
	if err != nil {
		return "", err
	}
	manifestPath := filepath.Join(directory, manifestName)
	manifest, err := OpenManifest(manifestPath)
	if err != nil {
		return "", err
	}

	packname := filepath.Base(abs) + ".brokerpak"
	return packname, manifest.Pack(directory, packname)
}

// Info writes out human-readable information about the brokerpak.
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
		w := cmdTabWriter()
		fmt.Fprintf(w, "format\t%d\n", mf.PackVersion)
		fmt.Fprintf(w, "name\t%s\n", mf.Name)
		fmt.Fprintf(w, "version\t%s\n", mf.Version)
		fmt.Fprintln(w, "platforms")
		for _, arch := range mf.Platforms {
			fmt.Fprintf(w, "\t%s\n", arch.String())
		}
		fmt.Fprintln(w, "metadata")
		for k, v := range mf.Metadata {
			fmt.Fprintf(w, "\t%s\t%s\n", k, v)
		}

		w.Flush()
		fmt.Println()
	}

	{
		fmt.Println("Dependencies")
		w := cmdTabWriter()
		fmt.Fprintln(w, "NAME\tVERSION\tSOURCE")
		for _, resource := range mf.TerraformResources {
			fmt.Fprintf(w, "%s\t%s\t%s\n", resource.Name, resource.Version, resource.Source)
		}
		w.Flush()
		fmt.Println()
	}

	{
		fmt.Println("Services")
		w := cmdTabWriter()
		fmt.Fprintln(w, "ID\tNAME\tDESCRIPTION\tPLANS")
		for _, svc := range services {
			fmt.Fprintf(w, "%s\t%s\t%s\t%d\n", svc.Id, svc.Name, svc.Description, len(svc.Plans))
		}
		w.Flush()
		fmt.Println()
	}

	fmt.Println("Contents")
	ziputil.List(&brokerPak.contents.Reader, os.Stdout)
	fmt.Println()

	return nil
}

func cmdTabWriter() *tabwriter.Writer {
	// args: output, minwidth, tabwidth, padding, padchar, flags
	return tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.StripEscape)
}

// Validate checks the brokerpak for syntactic and limited semantic errors.
func Validate(pack string) error {
	brokerPak, err := OpenBrokerPak(pack)
	if err != nil {
		return err
	}
	defer brokerPak.Close()

	return brokerPak.Validate()
}

// RegisterAll fetches all brokerpaks from the settings file and registers them
// with the given registry.
func RegisterAll(registry broker.BrokerRegistry) error {
	localPacks := viper.GetStringSlice("brokerpak.local_packs")

	for _, pack := range localPacks {
		if err := registerPak(pack, registry); err != nil {
			return err
		}
	}

	return nil
}

func registerPak(pack string, registry broker.BrokerRegistry) error {
	brokerPak, err := OpenBrokerPak(pack)
	if err != nil {
		return fmt.Errorf("couldn't open brokerpak: %q: %v", pack, err)
	}
	defer brokerPak.Close()

	if err := brokerPak.Register(registry); err != nil {
		return fmt.Errorf("couldn't register brokerpak: %q: %v", pack, err)
	}

	return nil
}

// RunExamples executes the examples from a brokerpak.
func RunExamples(pack string) error {
	registry := broker.BrokerRegistry{}
	if err := registerPak(pack, registry); err != nil {
		return err
	}

	apiClient, err := client.NewClientFromEnv()
	if err != nil {
		return err
	}

	return client.RunExamplesForService(registry, apiClient, "")
}

// Docs generates the markdown usage docs for the given pack and writes them to stdout.
func Docs(pack string) error {
	registry := broker.BrokerRegistry{}
	if err := registerPak(pack, registry); err != nil {
		return err
	}

	fmt.Println(generator.CatalogDocumentation(registry))
	return nil
}
