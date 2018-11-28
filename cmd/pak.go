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
		Long:  `interact with user-defined service definition bundles`,
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
		Use:   "build",
		Short: "bundle up the service definition files and Terraform resources",
		RunE: func(cmd *cobra.Command, args []string) error {
			return brokerpak.Pack("")
		},
	})

	pakCmd.AddCommand(&cobra.Command{
		Use:   "info pack",
		Short: "get info about a service definition pack",
		RunE: func(cmd *cobra.Command, args []string) error {
			return brokerpak.Info(args[0])
		},
	})

	pakCmd.AddCommand(&cobra.Command{
		Use:   "validate pack",
		Short: "validate a service definition pack",
		RunE: func(cmd *cobra.Command, args []string) error {
			return brokerpak.Validate(args[0])
		},
	})

}
