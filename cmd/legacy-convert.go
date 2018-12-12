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

	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/broker"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/providers/tf"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/varcontext"
	"github.com/GoogleCloudPlatform/gcp-service-broker/utils"
	"github.com/GoogleCloudPlatform/gcp-service-broker/utils/stream"
	"github.com/spf13/cobra"
)

const header = `# Copyright 2018 the Service Broker Project Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http:#www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
---`

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:    "dump-legacy",
		Short:  "Dump a legacy plan in a Terraform format",
		Args:   cobra.ExactArgs(1),
		Hidden: true,
		Run: func(cmd *cobra.Command, args []string) {
			reg := broker.DefaultRegistry
			desiredService := args[0]

			for _, svc := range reg.GetAllServices() {
				if svc.Name == desiredService {
					defn := legacyToTerraformDefinition(svc)
					fmt.Println(header)
					stream.Copy(stream.FromYaml(defn), stream.ToWriter(os.Stdout))
				}
			}
		},
	})
}

func legacyToTerraformDefinition(sd *broker.ServiceDefinition) tf.TfServiceDefinitionV1 {
	ce, err := sd.CatalogEntry()
	if err != nil {
		log.Fatal(err)
	}

	return tf.TfServiceDefinitionV1{
		Version:          1,
		Name:             sd.Name,
		Id:               ce.ID,
		Description:      ce.Description,
		DisplayName:      ce.Metadata.DisplayName,
		ImageUrl:         ce.Metadata.ImageUrl,
		DocumentationUrl: ce.Metadata.DocumentationUrl,
		SupportUrl:       ce.Metadata.SupportUrl,
		Tags:             ce.Tags,
		Plans:            legacyPlansToTerraformDef(ce.Plans),

		ProvisionSettings: createAction(sd.PlanVariables, sd.ProvisionInputVariables, nil, sd.ProvisionComputedVariables),
		BindSettings:      createAction(nil, sd.BindInputVariables, sd.BindOutputVariables, sd.BindComputedVariables),

		Examples: sd.Examples,
	}
}

func legacyPlansToTerraformDef(plans []broker.ServicePlan) []tf.TfServiceDefinitionV1Plan {
	var out []tf.TfServiceDefinitionV1Plan

	for _, p := range plans {
		converted := tf.TfServiceDefinitionV1Plan{
			Description: p.Description,
			Free:        false,
			Id:          p.ID,
			Name:        p.Name,
			Properties:  p.ServiceProperties,
		}

		if p.Metadata != nil {
			converted.DisplayName = p.Metadata.DisplayName
			converted.Bullets = p.Metadata.Bullets
		}

		out = append(out, converted)
	}

	return out
}

func createAction(plan, user, outputs []broker.BrokerVariable, computed []varcontext.DefaultVariable) tf.TfServiceDefinitionV1Action {
	template := ``

	varnames := utils.NewStringSet()
	for _, v := range plan {
		varnames.Add(v.FieldName)
	}
	for _, v := range user {
		varnames.Add(v.FieldName)
	}
	for _, v := range computed {
		varnames.Add(v.Name)
	}

	for _, v := range varnames.ToSlice() {
		template += fmt.Sprintf("variable %s {type = \"string\"}\n", v)
	}

	template += "\n\n"

	for _, v := range outputs {
		template += fmt.Sprintf("output %s {value = \"${}\"}\n", v.FieldName)
	}

	return tf.TfServiceDefinitionV1Action{
		PlanInputs: plan,
		UserInputs: user,
		Outputs:    outputs,
		Computed:   computed,
		Template:   template,
	}
}
