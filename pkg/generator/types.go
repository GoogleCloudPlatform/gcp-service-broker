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

// CatalogDocumentation generates markdown documentation for an entire service
// catalog.
func CatalogDocumentation() string {
	out := ""

	services := broker.GetAllServices()
	for _, svc := range services {
		out += generateServiceDocumentation(svc)
		out += "\n"
	}

	return out
}

// generateServiceDocumentation creates documentation for a single catalog entry
func generateServiceDocumentation(svc *broker.BrokerService) string {
	catalog, err := svc.CatalogEntry()
	if err != nil {
		log.Fatalf("Error getting catalog entry for service %s, %v", svc.Name, err)
	}

	vars := map[string]interface{}{
		"catalog":            catalog,
		"metadata":           catalog.Metadata,
		"bindIn":             svc.BindInputVariables,
		"bindOut":            svc.BindOutputVariables,
		"provisionInputVars": svc.ProvisionInputVariables,
		"examples":           svc.Examples,
	}

	funcMap := template.FuncMap{
		"code":          mdCode,
		"join":          strings.Join,
		"varNotes":      varNotes,
		"jsonCodeBlock": jsonCodeBlock,
	}

	templateText := `
--------------------------------------------------------------------------------

# ![]({{ .metadata.ImageUrl }}) {{ .metadata.DisplayName }}

{{ .metadata.LongDescription }}

 * [Documentation]({{.metadata.DocumentationUrl }})
 * [Support]({{ .metadata.SupportUrl }})
 * Catalog Metadata ID: {{code .catalog.ID}}
 * Tags: {{ join .catalog.Tags ", " }}

## Provisioning

**Request Parameters**

{{ if eq (len .provisionInputVars) 0 }}_No parameters supported._{{ end }}
{{ range $i, $var := .provisionInputVars }} * {{ varNotes $var }}
{{ end }}

## Binding

**Request Parameters**

{{ if eq (len .bindIn) 0 }}_No parameters supported._{{ end }}
{{ range $i, $var := .bindIn }} * {{ varNotes $var }}
{{ end }}
**Response Parameters**

{{ range $i, $var := .bindOut }} * {{ varNotes $var }}
{{ end }}
## Plans

{{ if eq (len .catalog.Plans) 0 }}_No plans available_{{ end }}
{{ range $i, $plan := .catalog.Plans }}  * **{{ $plan.Name }}**: {{ $plan.Description }} - Plan ID: {{code $plan.ID}}
{{ end }}

## Examples

{{ if eq (len .examples) 0 }}_No examples_{{ end }}

{{ range $i, $example := .examples}}
### {{ $example.Name }}


{{ $example.Description }}
Uses plan: {{ code $example.PlanId }}

**Provision**

{{ jsonCodeBlock $example.ProvisionParams }}

**Bind**

{{ jsonCodeBlock $example.BindParams }}
{{ end }}
`

	tmpl, err := template.New("titleTest").Funcs(funcMap).Parse(templateText)
	if err != nil {
		log.Fatalf("parsing: %s", err)
	}

	// Run the template to verify the output.
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, vars)
	if err != nil {
		log.Fatalf("execution: %s", err)
	}

	return buf.String()

}

func mdCode(text string) string {
	return fmt.Sprintf("`%s`", text)
}

func varNotes(variable broker.BrokerVariable) string {
	out := fmt.Sprintf("`%s` _%s_ - ", variable.FieldName, variable.Type)

	if variable.Required {
		out += "**Required** "
	}

	out += cleanLines(variable.Details)

	if variable.Default != nil {
		out += fmt.Sprintf(" Default: `%v`", variable.Default)
	}

	return out
}

// cleanLines concatenates multiple lines of text, trimming any leading/trailing
// whitespace
func cleanLines(text string) string {
	lines := strings.Split(text, "\n")
	for i, l := range lines {
		lines[i] = strings.TrimSpace(l)
	}

	return strings.Join(lines, " ")
}

// jsonCodeBlock formats the value as pretty JSON and wraps it in a Github style
// hilighted block.
func jsonCodeBlock(value interface{}) string {
	block, _ := json.MarshalIndent(value, "", "    ")
	return fmt.Sprintf("```javascript\n%s\n```", block)
}
