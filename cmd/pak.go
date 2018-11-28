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

package cmd

import (
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/brokerpak"
	"github.com/spf13/cobra"
)

func init() {
	pakCmd := &cobra.Command{
		Use:   "pak",
		Short: "interact with user-defined service definition bundles",
		Long: `Lets you create, validate, and view service definition bundles.

A service definition bundle is a zip file containing all the elements needed
to define and run a custom service.

Bundles include source code (for legal compliance), service definitions, and
Terraform/provider binaries for multiple platforms. They give you a contained
way to deploy new services to existing brokers or augment the broker to fit
your needs.

To start building a pack, create a new directory and within it run init:

	mkdir my-pak && cd my-pak
	gcp-service-broker pak init

You'll get a new pack with a manifest and example service definition.
Define the architectures and Terraform plugins you need in your manifest along
with any metadata you want, and include the names of all service definition
files.

When you're done, you can build the bundle which will download the sources,
Terraform resources, and pack them together.

	cd .. && gcp-service-broker pak build my-pak

This will produce a pack:

	my-pak.zip

You can validate the pack:

	gcp-service-broker pak validate my-pak.zip

You can also list information about the pack which includes metadata,
dependencies, services it provides, and the contents.

	gcp-service-broker pak info my-pak.zip

`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	rootCmd.AddCommand(pakCmd)

	pakCmd.AddCommand(&cobra.Command{
		Use:   "init",
		Short: "initialize a pak manifest and example service in the current directory",
		RunE: func(cmd *cobra.Command, args []string) error {
			return brokerpak.Init("")
		},
	})

	pakCmd.AddCommand(&cobra.Command{
		Use:   "build [path/to/pack/directory]",
		Short: "bundle up the service definition files and Terraform resources",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			directory := ""
			if len(args) == 1 {
				directory = args[0]
			}

			return brokerpak.Pack(directory)
		},
	})

	pakCmd.AddCommand(&cobra.Command{
		Use:   "info [pack.zip]",
		Short: "get info about a service definition pack",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return brokerpak.Info(args[0])
		},
	})

	pakCmd.AddCommand(&cobra.Command{
		Use:   "validate [pack.zip]",
		Short: "validate a service definition pack",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return brokerpak.Validate(args[0])
		},
	})

}
