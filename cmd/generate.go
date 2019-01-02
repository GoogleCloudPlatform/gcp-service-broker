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
	"fmt"

	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/generator"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/providers/builtin"
	"github.com/spf13/cobra"
)

func init() {
	generateCmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate documentation and tiles",
		Long:  `Generate documentation and tiles`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
	rootCmd.AddCommand(generateCmd)

	generateCmd.AddCommand(&cobra.Command{
		Use:   "use",
		Short: "Generate use markdown file",
		Long: `Generates the use.md file with:

	 * details about what each service is
	 * available parameters

	`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(generator.CatalogDocumentation(builtin.BuiltinBrokerRegistry()))
		},
	})

	generateCmd.AddCommand(&cobra.Command{
		Use:   "forms",
		Short: "Generate PCF Tile Forms",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(generator.GenerateFormsString())
		},
	})

	generateCmd.AddCommand(&cobra.Command{
		Use:   "customization",
		Short: "Generate customization documentation",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(generator.GenerateCustomizationMd())
		},
	})

	generateCmd.AddCommand(&cobra.Command{
		Use:   "tile",
		Short: "Generate tile.yml file",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Print(generator.GenerateTile())
		},
	})

	generateCmd.AddCommand(&cobra.Command{
		Use:   "manifest",
		Short: "Generate manifest.yml file",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(generator.GenerateManifest())
		},
	})
}
