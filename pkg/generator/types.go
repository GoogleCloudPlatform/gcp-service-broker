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

// Generate documentation
func CatalogDocumentation() string {
	out := ""

	for _, svc := range broker.GetAllServices() {
		out += generateServiceDocumentation(svc)
		out += "\n"
	}

	return out
}

func generateServiceDocumentation(svc *broker.BrokerService) string {
	catalog := svc.CatalogEntry()

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
		"image":         mdIcon,
		"link":          mdLink,
		"join":          strings.Join,
		"varNotes":      varNotes,
		"jsonCodeBlock": jsonCodeBlock,
	}

	templateText := `
# {{ image .metadata.ImageUrl }} {{ .metadata.DisplayName }}

{{ .metadata.LongDescription }}

 * {{ link "Documentation" .metadata.DocumentationUrl }}
 * {{ link "Support" .metadata.SupportUrl }}
 * Catalog Metadata ID: {{code .catalog.ID}}
 * Tags: {{ join .catalog.Tags ", " }}

## Provisioning

* Request Parameters
{{ range $i, $var := .provisionInputVars }}    * {{ varNotes $var }}
{{ end }}

## Binding

 * Request Parameters
{{ range $i, $var := .bindIn }}    * {{ varNotes $var }}
{{ end }}
 * Response Parameters
{{ range $i, $var := .bindOut }}    * {{ varNotes $var }}
{{ end }}

## Plans

{{ range $i, $plan := .catalog.Plans }}  * **{{ $plan.Name }}**: {{ $plan.Description }} - Plan ID: {{code $plan.ID}}
{{ end }}

## Examples
{{ range $i, $example := .examples}}
### {{ $example.Name }}


{{ $example.Description }}
Uses plan: {{ code $example.PlanId }}

**Provision**
{{ jsonCodeBlock $example.ProvisionParams }}

{{ if $example.BindParams }}
**Bind**

{{ jsonCodeBlock $example.BindParams }}
{{end}}

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

func mdIcon(url string) string {
	return fmt.Sprintf("![](%s)", url)
}

func mdCode(text string) string {
	return fmt.Sprintf("`%s`", text)
}

func mdLink(text, url string) string {
	return fmt.Sprintf("[%s](%s)", text, url)
}

func varNotes(variable broker.BrokerVariable) string {
	out := fmt.Sprintf("`%s` _%s_ - ", variable.FieldName, variable.Type)

	if variable.Required {
		out += "**Required** "
	}

	out += variable.Details

	if variable.Default != nil {
		out += fmt.Sprintf(" Default: `%s`", variable.Default)
	}

	return out
}

func jsonCodeBlock(value interface{}) string {
	block, _ := json.MarshalIndent(value, "", "    ")
	return fmt.Sprintf("```.json\n%s\n```", block)
}
