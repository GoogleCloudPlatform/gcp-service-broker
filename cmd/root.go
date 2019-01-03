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
	"log"
	"os"

	"github.com/GoogleCloudPlatform/gcp-service-broker/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "gcp-service-broker",
	Short: "GCP Service Broker is an OSB compatible service broker",
	Long:  `An OSB compatible service broker for Google Cloud Platform.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("WARNING: In the future running the broker from the root")
		fmt.Println("WARNING: command will show help instead.")
		fmt.Println("WARNING: Update your scripts to run gcp-service-broker serve")

		serve()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "Configuration file to be read")
	viper.SetEnvPrefix(utils.EnvironmentVarPrefix)
	viper.SetEnvKeyReplacer(utils.PropertyToEnvReplacer)
	viper.AutomaticEnv()
}

func initConfig() {
	if cfgFile == "" {
		return
	}

	viper.SetConfigFile(cfgFile)

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Can't read config: %v\n", err)
	}
}
