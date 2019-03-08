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

package server

import (
	"bytes"
	"html/template"
	"net/http"

	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/broker"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/generator"
	"github.com/gorilla/mux"
	blackfriday "gopkg.in/russross/blackfriday.v2"
)

var pageTemplate = template.Must(template.New("docs-page").Parse(`
<!DOCTYPE html>
<html lang="en">
	<head>
		<meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
		<title>{{.Title}}</title>
		<meta charset="utf-8" />
		<link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/4.0.0/css/bootstrap.min.css" crossorigin="anonymous" />
	</head>
	<body>
		<nav class="navbar navbar-expand navbar-dark sticky-top" style="background-color:#4285F4;">
			<a class="navbar-brand" href="#">
				<img src="https://cloud.google.com/_static/images/cloud/products/logos/svg/gcp-button-icon.svg" width="30" height="30" class="d-inline-block align-top" alt="">
				GCP Service Broker
			</a>
			<div>
				<ul class="navbar-nav">
					<li class="nav-item">
						<a class="nav-link active" href="/docs">Docs</a>
					</li>
					<li class="nav-item">
						<a class="nav-link active" href="/service-config">Service Configuration</a>
					</li>
				</ul>
			</div>
		</nav>
		<div class="container" id="maincontent">
			<br />
			{{ .Contents }}
		</div>
		<!-- Fixups for rendering markdown docs in the browser -->
		<script type="text/javascript">
			// add classes to the tables to style them nicely
			document.querySelectorAll("#maincontent > table").forEach( (node) => {node.classList = "table table-striped"});
		</script>
	</body>
</html>
`))

// AddDocsHandler creates a handler func that generates HTML documentation for
// the given registry and adds it to the /docs and / routes.
func AddDocsHandler(router *mux.Router, registry *broker.ServiceRegistry) error {
	docsPageMd := generator.CatalogDocumentation(registry)

	handler, err := renderAsPage("Service Broker Documents", docsPageMd)
	if err != nil {
		return err
	}

	router.Handle("/docs", handler)
	router.Handle("/", handler)

	return nil
}

// AddServiceConfigHandler creates a handler func that generates HTML
// documentation for service configurations on the given registry and
// adds it to the /service-config route.
func AddServiceConfigHandler(router *mux.Router, registry *broker.ServiceRegistry) error {
	docsPageMd, err := generator.GenerateServiceConfigMd(registry)
	if err != nil {
		return err
	}

	handler, err := renderAsPage("Service Broker Configuration", docsPageMd)
	if err != nil {
		return err
	}
	router.Handle("/service-config", handler)
	return nil
}

func renderAsPage(title, markdownContents string) (http.HandlerFunc, error) {
	renderer := blackfriday.NewHTMLRenderer(blackfriday.HTMLRendererParameters{Flags: blackfriday.HTMLFlagsNone})
	contents := blackfriday.Run([]byte(markdownContents), blackfriday.WithExtensions(blackfriday.CommonExtensions), blackfriday.WithRenderer(renderer))

	buf := &bytes.Buffer{}
	err := pageTemplate.Execute(buf, map[string]interface{}{
		"Title":    title,
		"Contents": template.HTML(contents),
	})

	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(200)
		w.Write(buf.Bytes())
	}, err
}
