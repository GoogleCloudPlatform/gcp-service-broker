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

	"code.cloudfoundry.org/lager"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/broker"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/client"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/generator"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/providers/tf"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/providers/tf/wrapper"
	"github.com/GoogleCloudPlatform/gcp-service-broker/utils"
	"github.com/GoogleCloudPlatform/gcp-service-broker/utils/stream"
	"github.com/GoogleCloudPlatform/gcp-service-broker/utils/ziputil"
)

const BrokerbakListConfigVar = "brokerpak.packs"

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
	manifest := &Manifest{}
	if err := stream.Copy(stream.FromFile(manifestPath), stream.ToYaml(manifest)); err != nil {
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
	registerLogger := utils.NewLogger("brokerpak-registration")
	// create a temp directory to hold all the paks
	pakDir, err := ioutil.TempDir("", "brokerpaks")
	if err != nil {
		return fmt.Errorf("couldn't create brokerpak staging area: %v", err)
	}
	defer os.RemoveAll(pakDir)

	pakConfig, err := NewServerConfigFromEnv()
	if err != nil {
		return err
	}

	// XXX(josephlewis42): this could be parallelized to increase performance
	// if we find people are pulling lots of data from the network.
	for name, pak := range pakConfig.Brokerpaks {
		registerLogger.Info("registering", lager.Data{
			"name":           name,
			"excluded-plans": pak.ExcludedPlansSlice(),
			"prefix":         pak.ServicePrefix,
		})

		if err := registerPak(pak, registry); err != nil {
			registerLogger.Error("registering", err, lager.Data{
				"name":           name,
				"excluded-plans": pak.ExcludedPlansSlice(),
				"prefix":         pak.ServicePrefix,
			})

			return err
		}
	}

	return nil
}

// RegisterPak fetches the brokerpak and registers it with the given registry.
func registerPak(config BrokerpakSourceConfig, registry broker.BrokerRegistry) error {
	brokerPak, err := DownloadAndOpenBrokerpak(config.BrokerpakUri)
	if err != nil {
		return fmt.Errorf("couldn't open brokerpak: %q: %v", config.BrokerpakUri, err)
	}
	defer brokerPak.Close()

	dir, err := ioutil.TempDir("", "brokerpak")
	if err != nil {
		return err
	}

	// extract the Terraform directory
	if err := brokerPak.ExtractPlatformBins(dir); err != nil {
		return err
	}

	binPath := filepath.Join(dir, "terraform")
	executor := wrapper.CustomTerraformExecutor(binPath, dir, wrapper.DefaultExecutor)

	// register the services
	services, err := brokerPak.Services()
	if err != nil {
		return err
	}

	toIgnore := utils.NewStringSet(config.ExcludedPlansSlice()...)
	for _, svc := range services {
		if toIgnore.Contains(svc.Id) {
			continue
		}

		svc.Name = config.ServicePrefix + svc.Name

		bs, err := svc.ToService(executor)
		if err != nil {
			return err
		}

		registry.Register(bs)
	}

	return nil
}

// RunExamples executes the examples from a brokerpak.
func RunExamples(pack string) error {
	registry := broker.BrokerRegistry{}
	if err := registerPak(NewBrokerpakSourceConfigFromPath(pack), registry); err != nil {
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
	if err := registerPak(NewBrokerpakSourceConfigFromPath(pack), registry); err != nil {
		return err
	}

	fmt.Println(generator.CatalogDocumentation(registry))
	return nil
}
