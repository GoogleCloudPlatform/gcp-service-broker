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
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/broker"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/validation"
)

// CatalogDocumentation generates markdown documentation for the service catalog
// of the given registry. Returns all content in one string.
func CatalogDocumentation(registry broker.BrokerRegistry) string {
	const contentSeparator = "\n--------------------------------------------------------------------------------\n"
	out := ""

	services := registry.GetAllServices()
	for _, svc := range services {
		out += generateServiceDocumentation(svc)
		out += contentSeparator
	}

	return cleanMdOutput(out)
}

// CatalogDocumentationToDir generates markdown documentation for the service catalog
// of the given registry. Generated documentation is saved under given `dstDir` with such layout:
//
// ./{dstDir}
// ├── classes
// │   ├── {class_name}.md # contains all information about given service class
// │   └── ...
// └── use.md # contains table of contents (ToC) about all generated service class documentation
//
func CatalogDocumentationToDir(registry broker.BrokerRegistry, dstDir string) {
	const (
		mdExt    = ".md"
		classDir = "classes"
		perms    = 0644
	)

	classesDir := filepath.Join(dstDir, classDir)
	if err := ensureDir(classesDir); err != nil {
		log.Fatalf("Cannot create directory for storing broker classes documentation: %v", err)
	}

	var tocServiceClasses []ToCServiceClass
	for _, svc := range registry.GetAllServices() {
		out := generateServiceDocumentation(svc)
		out = cleanMdOutput(out)

		saveClassPath := filepath.Join(classesDir, svc.Name+mdExt)
		err := ioutil.WriteFile(saveClassPath, []byte(out), perms)
		if err != nil {
			log.Fatalf("Cannot save %s documentation class into %s: %v", svc.Name, saveClassPath, err)
		}

		tocClassPath, err := filepath.Rel(dstDir, saveClassPath)
		if err != nil {
			log.Fatalf("Cannot resolve relative path to %s file: %v", saveClassPath, err)
		}

		tocServiceClasses = append(tocServiceClasses, ToCServiceClass{
			DisplayName: svc.DisplayName,
			FilePath:    tocClassPath,
		})
	}

	toc := generateServiceClassesToC(tocServiceClasses)
	tocFileName := filepath.Join(dstDir, "use.md")
	if err := ioutil.WriteFile(tocFileName, []byte(toc), perms); err != nil {
		log.Fatalf("Cannot save documentation from %s class into %s", "", "")
	}
}

type ToCServiceClass struct {
	DisplayName string
	FilePath    string
}

func generateServiceClassesToC(tocServiceClasses []ToCServiceClass) string {
	tocTemplateText := `# GCP Service Broker usage

# Overview

The GCP service broker is a tool that developers use to provision access Google Cloud resources. The service broker currently supports a list of built-in services.
 
Read about creating and binding specific GCP services:
{{range .serviceClasses}}
-   [{{ .DisplayName }}]({{ .FilePath }})
{{ end -}}
`

	vars := map[string]interface{}{
		"serviceClasses": tocServiceClasses,
	}

	return render(tocTemplateText, vars, nil)
}

func ensureDir(dirPath string) error {
	err := os.Mkdir(dirPath, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		return err
	}
	return nil
}

// generateServiceDocumentation creates documentation for a single catalog entry
func generateServiceDocumentation(svc *broker.ServiceDefinition) string {
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
		"exampleCommands": func(example broker.ServiceExample) string {
			planName := "unknown-plan"
			for _, plan := range catalog.Plans {
				if plan.ID == example.PlanId {
					planName = plan.Name
				}
			}

			params, err := json.Marshal(example.ProvisionParams)
			if err != nil {
				return err.Error()
			}
			provision := fmt.Sprintf("$ cf create-service %s %s my-%s-example -c `%s`", catalog.Name, planName, catalog.Name, params)

			params, err = json.Marshal(example.BindParams)
			if err != nil {
				return err.Error()
			}
			bind := fmt.Sprintf("$ cf bind-service my-app my-%s-example -c `%s`", catalog.Name, params)
			return provision + "\n" + bind
		},
	}

	templateText := `
# <a name="{{ .catalog.Name }}"></a> ![]({{ .metadata.ImageUrl }}) {{ .metadata.DisplayName }}
{{ .metadata.LongDescription }}

 * [Documentation]({{.metadata.DocumentationUrl }})
 * [Support]({{ .metadata.SupportUrl }})
 * Catalog Metadata ID: {{code .catalog.ID}}
 * Tags: {{ join .catalog.Tags ", " }}
 * Service Name: {{ code .catalog.Name }}

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

The following plans are built-in to the GCP Service Broker and may be overridden
or disabled by the broker administrator.

{{ if eq (len .catalog.Plans) 0 }}_No plans available_{{ end }}
{{ range $i, $plan := .catalog.Plans -}}
* **{{code $plan.Name }}**
  * Plan ID: {{code $plan.ID}}.
  * Description: {{ $plan.Description }}
  * This plan {{ if eq (len $plan.ProvisionOverrides) 0 -}} doesn't override {{- else -}} overrides the following {{- end}} user variables on provision.
{{ range $k, $v := $plan.ProvisionOverrides}}    * {{ code $k }} = {{code $v}}
{{ end }}  * This plan {{ if eq (len $plan.BindOverrides) 0 -}} doesn't override {{- else -}} overrides the following {{- end}} user variables on bind.
{{ range $k, $v := $plan.BindOverrides}}    * {{ code $k }} = {{code $v}}
{{ end }}
{{- end }}

## Examples

{{ if eq (len .examples) 0 }}_No examples._{{ end }}

{{ range $i, $example := .examples}}
### {{ $example.Name }}


{{ $example.Description }}
Uses plan: {{ code $example.PlanId }}.

**Provision**

{{ jsonCodeBlock $example.ProvisionParams }}

**Bind**

{{ jsonCodeBlock $example.BindParams }}

**Cloud Foundry Example**

<pre>
{{exampleCommands $example}}
</pre>

{{ end }}
`

	return render(templateText, vars, funcMap)
}

func render(tmplText string, vars interface{}, funcMap template.FuncMap) string {
	tmpl, err := template.New("rendered").Funcs(funcMap).Parse(tmplText)
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

func mdCode(text interface{}) string {
	return fmt.Sprintf("`%v`", text)
}

func varNotes(variable broker.BrokerVariable) string {
	out := fmt.Sprintf("`%s` _%s_ - ", variable.FieldName, variable.Type)

	if variable.Required {
		out += "**Required** "
	}

	out += cleanLines(variable.Details)

	if variable.Default != nil {
		out += fmt.Sprintf(" Default: `%v`.", variable.Default)
	}

	bullets := constraintsToDoc(variable.ToSchema())
	if len(bullets) > 0 {
		out += "\n    * "
		out += strings.Join(bullets, "\n    * ")
	}

	return out
}

// constraintsToDoc converts a map of JSON Schema validation key/values to human-readable bullet points.
func constraintsToDoc(schema map[string]interface{}) []string {
	// We use an anonymous struct rather than a map to get a strict ordering of
	// constraints so they are generated consistently in documentation.
	// Not all JSON Schema constraints can be cleanly expressed in this format,
	// nor do we use them all so some are missing.
	constraintFormatters := []struct {
		SchemaKey string
		DocString string
	}{
		// Schema Annotations
		{validation.KeyExamples, "Examples: %+v."},

		// Validation for any instance type
		{validation.KeyEnum, "The value must be one of: %+v."},
		{validation.KeyConst, "The value must be: `%v`."},

		// Validation keywords for numeric instances
		{validation.KeyMultipleOf, "The value must be a multiple of %v."},
		{validation.KeyMaximum, "The value must be less than or equal to %v."},
		{validation.KeyExclusiveMaximum, "The value must be strictly less than %v."},
		{validation.KeyMinimum, "The value must be greater than or equal to %v."},
		{validation.KeyExclusiveMaximum, "The value must be strictly greater than %v."},

		// Validation keywords for strings
		{validation.KeyMaxLength, "The string must have at most %v characters."},
		{validation.KeyMinLength, "The string must have at least %v characters."},
		{validation.KeyPattern, "The string must match the regular expression `%v`."},

		// Validation keywords for arrays
		{validation.KeyMaxItems, "The array must have at most %v items."},
		{validation.KeyMinItems, "The array must have at least %v items."},

		// Validation keywords for objects
		{validation.KeyMaxProperties, "The object must have at most %v properties."},
		{validation.KeyMinProperties, "The object must have at least %v properties."},
		{validation.KeyRequired, "The following properties are required: %v."},
		{validation.KeyMaxProperties, "Property names must match the JSON Schema: `%+v`."},
	}

	var bullets []string
	for _, formatter := range constraintFormatters {
		if v, ok := schema[formatter.SchemaKey]; ok {
			bullets = append(bullets, fmt.Sprintf(formatter.DocString, v))
		}
	}

	return bullets
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
