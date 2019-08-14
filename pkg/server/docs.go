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
	"github.com/russross/blackfriday"
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
func AddDocsHandler(router *mux.Router, registry broker.BrokerRegistry) {
	docsPageMd := generator.CatalogDocumentation(registry)
	handler := renderAsPage("Service Broker Documents", docsPageMd)

	router.Handle("/docs", handler)
	router.Handle("/", handler)
}

func renderAsPage(title, markdownContents string) http.HandlerFunc {
	renderer := blackfriday.HtmlRenderer(
		blackfriday.EXTENSION_FENCED_CODE|
			blackfriday.EXTENSION_AUTOLINK,
		title,
		"https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/css/bootstrap.min.css",
	)
	page := blackfriday.Markdown([]byte(markdownContents), renderer, 0)

	buf := &bytes.Buffer{}
	err := pageTemplate.Execute(buf, map[string]interface{}{
		"Title":    title,
		"Contents": template.HTML(page),
	})

	if err != nil {
		return func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
	}

	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(200)
		w.Write(buf.Bytes())
	}
}
