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

package generator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"text/template"

	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/broker"
)

const (
	formDocumentation = `
# Installation Customization

This file documents the various environment variables you can set to change the functionality of the service broker.
If you are using the PCF Tile deployment, then you can manage all of these options through the operator forms.
If you are running your own, then you can set them in the application manifest of a PCF deployment, or in your pod configuration for Kubernetes.

{{ range $i, $f := .Forms }}{{ template "normalform" $f }}{{ end }}

## Install Brokerpaks

You can install one or more brokerpaks using the <tt>GSB_BROKERPAK_SOURCES</tt>
environment variable.

The value should be a JSON array containing zero or more brokerpak configuration
objects with the following properties:

{{ with .BrokerpakForm  }}
| Property | Type | Description |
|----------|------|-------------|
{{ range .Properties -}}
| <tt>{{.Name}}</tt>{{ if not .Optional }} <b>*</b>{{end}} | {{ .Type }} | <p>{{ .Label }}. {{ .Description }}{{if .Default }} Default: <code>{{ js .Default }}</code>{{- end }}</p>|
{{ end }}

\* = Required
{{ end }}

### Example

Here is an example that loads three brokerpaks.

	[
		{
			"notes":"GA services for all users.",
			"uri":"https://link/to/artifact.brokerpak?checksum=md5:3063a2c62e82ef8614eee6745a7b6b59",
			"excluded_services":"00000000-0000-0000-0000-000000000000",
			"config":{}
		},
		{
			"notes":"Beta services for all users.",
			"uri":"gs://link/to/beta.brokerpak",
			"service_prefix":"beta-",
			"config":{}
		},
		{
			"notes":"Services for the marketing department. They use their own GCP Project.",
			"uri":"https://link/to/marketing.brokerpak",
			"service_prefix":"marketing-",
			"config":{"PROJECT_ID":"my-marketing-project"}
		},
	]

## Customizing Services

You can customize specific services by changing their defaults, disabling them, or creating custom plans.

The <tt>GSB_SERVICE_CONFIG</tt> environment variable holds all the customizations as a JSON map.
The keys of the map are the service's GUID, and the value is a configuration object.

Example:

	{
		"51b3e27e-d323-49ce-8c5f-1211e6409e82":{ /* Spanner Configuration Object */ },
		"628629e3-79f5-4255-b981-d14c6c7856be":{ /* Pub/Sub Configuration Object */ },
		...
	}

**Configuration Object**

| Property | Type | Description |
|----------|------|-------------|
| <tt>//</tt> | string | Space for your notes. |
| <tt>disabled</tt> | boolean | If set to true, this service will be hidden from the catalog. |
| <tt>provision_defaults</tt> | string:any map | A map of provision property/default pairs that are used to populate missing values in provision requests. |
| <tt>bind_defaults</tt> | string:any map | A map of bind property/default pairs that are used to populate missing values in bind requests. |
| <tt>custom_plans</tt> | array of custom plan objects | You can add custom service plans here. See below for the object structure. |

**Custom Plan Object**

| Property | Type | Description |
|----------|------|-------------|
| <tt>guid</tt> \* | string | A GUID for this plan, must be unique. Changing this value after services are using it WILL BREAK your instances. |
| <tt>name</tt> \* | string | A CLI friendly name for this plan. This can be changed without affecting existing instances, but may break scripts you build referencing it. |
| <tt>display_name</tt> \* | string | A human readable name for this plan, this can be changed. |
| <tt>description</tt> \* | string | A human readable description for this plan, this can be changed. |
| <tt>properties</tt> \* | string:string map | Properties used to configure the plan. Each service has its own set of properties used to customize it. |

\* = Required

{{ range .Services }}

### {{ .DisplayName }}<a id="{{.Name}}"></a>

{{ .Description }}

Configuration needs to be done under the GUID: <tt>{{.Id}}</tt>.

#### Example

{{ exampleServiceConfig . }}

_Note: the example includes the configuration and the GUID it should be nested under._

#### Provision Defaults

Setting a value for any of these in the <tt>provision_defaults</tt> map
will override the default value the provision call uses for the property.

{{ documentBrokerVariables .ProvisionInputVariables true }}

#### Bind Defaults

Setting a value for any of these in the <tt>bind_defaults</tt> map
will override the default value the provision call uses for the property.

{{ documentBrokerVariables .BindInputVariables true }}

#### Custom Plan Properties

{{ documentBrokerVariables .PlanVariables false }}

{{ end }}

---------------------------------------

_Note: **Do not edit this file**, it was auto-generated by running <code>gcp-service-broker generate customization</code>. If you find an error, change the source code in <tt>customization-md.go</tt> or file a bug._

{{/*=======================================================================*/}}
{{ define "normalform" }}
## {{ .Label }}

{{ .Description }}

You can configure the following environment variables:

| Environment Variable | Type | Description |
|----------------------|------|-------------|
{{ range .Properties -}}
| <tt>{{upper .Name}}</tt>{{ if not .Optional }} <b>*</b>{{end}} | {{ .Type }} | <p>{{ .Label }}. {{ .Description }}{{if .Default }} Default: <code>{{ js .Default }}</code>{{- end }}</p>|
{{ end }}

\* = Required

{{ end }}
`
)

var (
	customizationTemplateFuncs = template.FuncMap{
		"upper":                   strings.ToUpper,
		"exampleServiceConfig":    exampleServiceConfig,
		"documentBrokerVariables": documentBrokerVariables,
	}
	formDocumentationTemplate = template.Must(template.New("name").Funcs(customizationTemplateFuncs).Parse(formDocumentation))
)

func GenerateCustomizationMd(registry *broker.ServiceRegistry) string {
	tileForms := GenerateForms()

	env := map[string]interface{}{
		"Forms":         tileForms.Forms,
		"BrokerpakForm": brokerpakConfigurationForm(),
		"Services":      registry.GetAllServices(),
	}

	var buf bytes.Buffer
	if err := formDocumentationTemplate.Execute(&buf, env); err != nil {
		log.Fatalf("Error rendering template: %s", err)
	}

	return cleanMdOutput(buf.String())
}

// Remove trailing whitespace from the document and every line
func cleanMdOutput(text string) string {
	text = strings.TrimSpace(text)

	lines := strings.Split(text, "\n")
	for i, l := range lines {
		lines[i] = strings.TrimRight(l, " \t")
	}

	return strings.Join(lines, "\n")
}

func exampleServiceConfig(defn *broker.ServiceDefinition) string {
	out := map[string]interface{}{
		defn.Id: broker.ServiceConfig{
			BindDefaults:      createExampleDefaults(defn.BindInputVariables, "bind"),
			ProvisionDefaults: createExampleDefaults(defn.ProvisionInputVariables, "provision"),
			CustomPlans:       createExampleCustomPlan(defn),
		},
	}

	bytes, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return err.Error()
	}

	return "```json\n" + string(bytes) + "```\n"
}

func createExampleDefaults(vars []broker.BrokerVariable, action string) map[string]interface{} {
	example := make(map[string]interface{})

	if len(vars) == 0 {
		example["//"] = fmt.Sprintf("The %s action takes no params so it can't be overridden.", action)
	} else {
		example["//"] = fmt.Sprintf("See the '%s defaults' section below for defaults you can change.", action)
	}

	return example
}

func documentBrokerVariables(vars []broker.BrokerVariable, tableOnly bool) string {
	if len(vars) == 0 {
		return "_There are no configurable properties for this object._"
	}

	buf := &bytes.Buffer{}

	fmt.Fprintln(buf, "| Property | Type | Description |")
	fmt.Fprintln(buf, "|----------|------|-------------|")

	for _, v := range vars {
		required := ""
		if !tableOnly && v.Required {
			required = " \\*"
		}

		fmt.Fprintf(buf, "| `%s`%s | %s | %s |", v.FieldName, required, v.Type, singleLine(v.Details))
		fmt.Fprintln(buf)
	}

	if !tableOnly {
		fmt.Fprintln(buf, `\* = Required`)
	}

	fmt.Fprintln(buf)

	return buf.String()
}

func createExampleCustomPlan(service *broker.ServiceDefinition) []broker.CustomPlan {
	planVars := service.PlanVariables

	// if this service isn't configurable, don't show an example plan.
	if len(planVars) == 0 {
		return []broker.CustomPlan{}
	}

	return []broker.CustomPlan{
		{
			GUID:        "00000000-0000-0000-0000-000000000000",
			Name:        "a-cli-friendly-name",
			DisplayName: "A human-readable name",
			Description: "What makes this plan different?",
			Properties: map[string]string{
				"//": "See the custom plan properties section below for configurable properties.",
			},
		},
	}
}
