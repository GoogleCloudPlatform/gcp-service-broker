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
	"context"
	"fmt"
	"log"
	"os"

	"code.cloudfoundry.org/lager"
	"github.com/GoogleCloudPlatform/gcp-service-broker/db_service"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/providers/tf"
	"github.com/spf13/cobra"
)

func init() {
	tfCmd := &cobra.Command{
		Use:   "tf",
		Short: "Interact with the Terraform backend",
		Long:  `Interact with the Terraform backend`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	rootCmd.AddCommand(tfCmd)

	tfCmd.AddCommand(&cobra.Command{
		Use:   "dump",
		Short: "dump a Terraform workspace",
		Run: func(cmd *cobra.Command, args []string) {

			logger := lager.NewLogger("dump-command")
			db_service.New(logger)

			fmt.Printf("Dumping %q\n", args[0])

			tfjr := tf.TfJobRunner{}
			ws, err := tfjr.Dump(context.Background(), args[0])
			if err != nil {
				fmt.Printf("Error: %s\n", err.Error())
				log.Fatal(err)
			}

			fmt.Println(ws)
		},
	})

	tfCmd.AddCommand(&cobra.Command{
		Use:   "apply",
		Short: "apply a Terraform workspace",
		Run: func(cmd *cobra.Command, args []string) {

			logger := lager.NewLogger("apply-command")
			logger.RegisterSink(lager.NewWriterSink(os.Stderr, lager.ERROR))
			logger.RegisterSink(lager.NewWriterSink(os.Stdout, lager.DEBUG))

			db_service.New(logger)

			fmt.Printf("Applying %q\n", args[0])

			tfjr, err := tf.NewTfJobRunerFromEnv()
			if err != nil {
				log.Fatal(err.Error())
			}
			if err := tfjr.Create(context.Background(), args[0]); err != nil {
				log.Fatal(err.Error())
			}

			for {
				isDone, err := tfjr.Status(context.Background(), args[0])

				if err != nil {
					log.Printf("Got error: %s", err.Error())
				}

				if isDone {
					return
				}
			}
		},
	})

	tfCmd.AddCommand(&cobra.Command{
		Use:   "destroy",
		Short: "destroy a Terraform workspace",
		Run: func(cmd *cobra.Command, args []string) {

			logger := lager.NewLogger("apply-command")
			logger.RegisterSink(lager.NewWriterSink(os.Stderr, lager.ERROR))
			logger.RegisterSink(lager.NewWriterSink(os.Stdout, lager.DEBUG))

			db_service.New(logger)

			fmt.Printf("Applying %q\n", args[0])

			tfjr, err := tf.NewTfJobRunerFromEnv()
			if err != nil {
				log.Fatal(err.Error())
			}
			if err := tfjr.Destroy(context.Background(), args[0]); err != nil {
				log.Fatal(err.Error())
			}

			for {
				isDone, err := tfjr.Status(context.Background(), args[0])

				if err != nil {
					log.Printf("Got error: %s", err.Error())
				}

				if isDone {
					return
				}
			}
		},
	})

	tfCmd.AddCommand(&cobra.Command{
		Use:   "outputs",
		Short: "outputs of a Terraform workspace",
		Run: func(cmd *cobra.Command, args []string) {

			logger := lager.NewLogger("outputs-command")
			logger.RegisterSink(lager.NewWriterSink(os.Stderr, lager.ERROR))
			logger.RegisterSink(lager.NewWriterSink(os.Stdout, lager.DEBUG))

			db_service.New(logger)

			fmt.Printf("Applying %q\n", args[0])

			tfjr, err := tf.NewTfJobRunerFromEnv()
			if err != nil {
				log.Fatal(err.Error())
			}

			outputs, err := tfjr.Outputs(context.Background(), args[0], "instance")
			if err != nil {
				log.Fatal(err.Error())
			}

			logger.Info("outputs", outputs)
		},
	})

}
