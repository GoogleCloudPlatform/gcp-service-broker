// Copyright the Service Broker Project Authors.
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
	"encoding/json"
	"fmt"

	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	serviceId      string
	planId         string
	instanceId     string
	bindingId      string
	parametersJson string
)

func init() {
	rootCmd.AddCommand(clientCmd)
	clientCmd.AddCommand(clientCatalogCmd, provisionCmd, deprovisionCmd, bindCmd, unbindCmd, lastCmd)

	resourceSubcommands := []*cobra.Command{provisionCmd, deprovisionCmd, bindCmd, unbindCmd}
	creationSubcommands := []*cobra.Command{provisionCmd, bindCmd}
	bindingSubcommands := []*cobra.Command{bindCmd, unbindCmd}

	for _, sc := range resourceSubcommands {
		sc.Flags().StringVarP(&instanceId, "instanceid", "", "", "id of the service instance to operate on (user defined)")
		sc.MarkFlagRequired("instanceid")
		sc.Flags().StringVarP(&serviceId, "serviceid", "", "", "GUID of the service instanceid references (see catalog)")
		sc.MarkFlagRequired("serviceid")
		sc.Flags().StringVarP(&planId, "planid", "", "", "GUID of the service instanceid references (see catalog entry for the associated serviceid)")
		sc.MarkFlagRequired("planid")
	}

	for _, sc := range creationSubcommands {
		sc.Flags().StringVarP(&parametersJson, "params", "", "{}", "JSON string of user-defined paramaters to pass to the request")
	}

	for _, sc := range bindingSubcommands {
		sc.Flags().StringVarP(&bindingId, "bindingid", "", "", "GUID of the binding to work on (user defined)")
		sc.MarkFlagRequired("bindingid")
	}

	lastCmd.Flags().StringVarP(&instanceId, "instanceid", "", "", "id of the service instance to operate on (user defined)")
	lastCmd.MarkFlagRequired("instanceid")
}

var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "A CLI client for the service broker",
	Long: `A CLI client for the service broker.

The client commands use the same configuration values as the server and operate
on localhost using the HTTP protocol.

Configuration Params:

 - api.user
 - api.password
 - api.port

Environment Variables:

 - GSB_API_USER
 - GSB_API_PASSWORD
 - GSB_API_PORT

The client commands return formatted JSON when run if the exit code is 0:

	{
	    "url": "http://user:pass@localhost:8000/v2/catalog",
	    "http_method": "GET",
	    "status_code": 200,
	    "response": // Response Body as JSON
	}

Exit codes DO NOT correspond with status_code, if a request was made and the
response could be parsed then the exit code will be 0.
Non-zero exit codes indicate a failure in the executable.

Because of the format, you can use the client to do automated testing of your
user-defined plans.
`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var clientCatalogCmd = newClientCommand("catalog", "Show the service catalog", func(client *client.Client) *client.BrokerResponse {
	return client.Catalog()
})

var provisionCmd = newClientCommand("provision", "Provision a service", func(client *client.Client) *client.BrokerResponse {
	return client.Provision(instanceId, serviceId, planId, json.RawMessage(parametersJson))
})

var deprovisionCmd = newClientCommand("deprovision", "Derovision a service", func(client *client.Client) *client.BrokerResponse {
	return client.Deprovision(instanceId, serviceId, planId)
})

var bindCmd = newClientCommand("bind", "Bind to a service", func(client *client.Client) *client.BrokerResponse {
	return client.Bind(instanceId, bindingId, serviceId, planId, json.RawMessage(parametersJson))
})

var unbindCmd = newClientCommand("unbind", "Unbind a service", func(client *client.Client) *client.BrokerResponse {
	return client.Unbind(instanceId, bindingId, serviceId, planId)
})

var lastCmd = newClientCommand("last", "Get the status of the last operation", func(client *client.Client) *client.BrokerResponse {
	return client.LastOperation(instanceId)
})

func newClientCommand(use, short string, run func(*client.Client) *client.BrokerResponse) *cobra.Command {
	return &cobra.Command{
		Use:   use,
		Short: short,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := createClient()
			if err != nil {
				return err
			}

			return printJsonResults(run(client))
		},
	}
}

func createClient() (*client.Client, error) {
	user := viper.GetString("api.user")
	pass := viper.GetString("api.password")
	port := viper.GetInt("api.port")

	return client.New(user, pass, "localhost", port)
}

func printJsonResults(results interface{}) error {
	prettyResults, err := json.MarshalIndent(results, "", "    ")
	if err != nil {
		return err
	}

	fmt.Println(string(prettyResults))
	return nil
}
