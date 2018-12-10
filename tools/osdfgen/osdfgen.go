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

package main

import (
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"gopkg.in/src-d/go-license-detector.v2/licensedb"
	"gopkg.in/src-d/go-license-detector.v2/licensedb/filer"

	toml "github.com/pelletier/go-toml"
)

var (
	out  = flag.String("o", "-", "Sets the output location of the OSDF csv")
	proj = flag.String("p", ".", "the project root")
)

type Lockfile struct {
	Projects []Project
}

type Project struct {
	Name     string
	Revision string
}

func main() {
	flag.Parse()

	tree, err := toml.LoadFile(filepath.Join(*proj, "Gopkg.lock"))
	if err != nil {
		log.Fatalf("Error loading Gopkg.lock, %s", err)
	}

	deps := Lockfile{}
	if err := tree.Unmarshal(&deps); err != nil {
		log.Fatalf("Error unmarshaling lockfile %s", err)
	}

	var writer *csv.Writer
	if *out == "-" {
		writer = csv.NewWriter(os.Stdout)
	} else {
		fd, err := os.Create(*out)
		if err != nil {
			log.Fatalf("Error opening %q, %s", *out, err)
		}
		defer fd.Close()
		writer = csv.NewWriter(fd)
	}

	for _, project := range deps.Projects {
		dir, err := filer.FromDirectory(filepath.Join(*proj, "vendor", project.Name))
		if err != nil {
			log.Fatalf("Could not find dep %q in vendor %s", project.Name, err)
		}

		licenses, err := licensedb.Detect(dir)
		if err != nil {
			log.Fatalf("Could not detect licenses %s", err)
		}

		spdxCode := mostLikelyLicense(licenses)

		licenseText, err := getLicenseText(project, spdxCode)
		if err != nil {
			log.Printf("Could not get license text for %q: %s", project.Name, err)
		}

		// name, hash, spdx, full text
		writer.Write([]string{project.Name, project.Revision, spdxCode, licenseText})
	}

	writer.Flush()

}

func mostLikelyLicense(m map[string]float32) string {
	var maxVal float32
	maxKey := ""

	for k, v := range m {
		if v > maxVal {
			maxVal = v
			maxKey = k
		}
	}

	return maxKey
}

func getLicenseText(project Project, spdxCode string) (string, error) {
	projectRoot := filepath.Join(*proj, "vendor", project.Name)
	entries, err := ioutil.ReadDir(projectRoot)
	if err != nil {
		log.Fatalf("Could not find dep %q in vendor %s", project.Name, err)
	}

	licenses := ""

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		base := strings.ToLower(path.Base(entry.Name()))
		fullPath := path.Join(projectRoot, entry.Name())
		if strings.Contains(base, "license") || (base == "notice" && spdxCode == "Apache-2.0") {
			licenses += fmt.Sprintf("Contents of: %s@%s\n\n", fullPath, project.Revision)
			text, err := ioutil.ReadFile(fullPath)
			if err != nil {
				return licenses, err
			}

			licenses += string(text)
			licenses += "\n\n"
		}
	}

	if licenses == "" {
		return "", errors.New("Could not find license text")
	}

	return licenses, nil
}
