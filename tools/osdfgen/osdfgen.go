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

// +build !service_broker

// `osdfgen` can be used to build a CSV suitable for uploading to Pivotal's
// [OSDF Generator](http://osdf-generator.cfapps.io/static/index.html).
// It determines licenses by sniffing the dependencies listed in `Gopkg.lock`.
// Example: go run osdfgen.go -p ../../ -o test.csv
// The `-p` flag points at the project root and the `-o` flag is the place to put the output (stdout by default).
package main

import (
	"bytes"
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/src-d/go-license-detector.v2/licensedb"
	"gopkg.in/src-d/go-license-detector.v2/licensedb/filer"

	toml "github.com/pelletier/go-toml"
)

type Lockfile struct {
	Projects []Project
}

type Project struct {
	Name     string
	Revision string
}

type Dependency struct {
	ParentDirectory string
	Project         Project
}

func (d *Dependency) Directory() string {
	return filepath.Join(d.ParentDirectory, "vendor", d.Project.Name)
}

func main() {
	out := flag.String("o", "-", "Sets the output location of the OSDF csv")
	proj := flag.String("p", ".", "The project root")

	flag.Parse()

	buf := &bytes.Buffer{}
	writer := csv.NewWriter(buf)
	for _, project := range getProjects(*proj) {
		dep := &Dependency{
			ParentDirectory: *proj,
			Project:         project,
		}

		licenses, err := detectLicenses(dep.Directory())
		if err != nil {
			log.Fatal(err)
		}

		spdxCode := mostLikelyLicense(licenses)
		licenseText, err := getLicenseText(dep, spdxCode)
		if err != nil {
			log.Fatalf("Could not get license text for %q: %s", project.Name, err)
		}

		// name, hash, spdx, full text
		writer.Write([]string{project.Name, project.Revision, spdxCode, licenseText})
	}

	writer.Flush()
	writeOutput(*out, buf)
}

// detectLicenses returns a map of the SPDX codes for the licenses in a given
// directory.
func detectLicenses(directory string) (map[string]float32, error) {
	dir, err := filer.FromDirectory(directory)
	if err != nil {
		return nil, fmt.Errorf("Could not find dep %q in vendor %s", dir, err)
	}

	return licensedb.Detect(dir)
}

// writeOutput writes the contents of src into the file denoted by fileName.
// if fileName is "-" the UNIX convention is followed and the contents are
// written to stdout.
func writeOutput(fileName string, src io.Reader) {
	var dest io.Writer
	if fileName == "-" {
		dest = os.Stdout
	} else {
		var err error
		fd, err := os.Create(fileName)
		if err != nil {
			log.Fatalf("Error opening %q, %s", fileName, err)
		}
		defer fd.Close()
		dest = fd
	}

	if _, err := io.Copy(dest, src); err != nil {
		log.Fatalf("Error trying to write output: %v", err)
	}
}

// getProjects reads a Gopkg.lock file and returns a list of the projects it
// contains.
func getProjects(projectRoot string) []Project {
	tree, err := toml.LoadFile(filepath.Join(projectRoot, "Gopkg.lock"))
	if err != nil {
		log.Fatalf("Error loading Gopkg.lock, %s", err)
	}

	deps := Lockfile{}
	if err := tree.Unmarshal(&deps); err != nil {
		log.Fatalf("Error unmarshaling lockfile %s", err)
	}

	return deps.Projects
}

// mostLikelyLicense finds the key in a map with the greatest value.
func mostLikelyLicense(m map[string]float32) string {
	var maxVal float32 = -math.MaxFloat32
	maxKey := ""

	for k, v := range m {
		if v > maxVal {
			maxVal = v
			maxKey = k
		}
	}

	return maxKey
}

// getLicenseText creates a copy of the licence text(s) and notice(s) for a
// given project.
// Returns an error if no license could be found.
func getLicenseText(project *Dependency, spdxCode string) (string, error) {
	dir, err := filer.FromDirectory(project.Directory())
	if err != nil {
		log.Fatalf("Could not find dep %q in vendor %s", project.Project.Name, err)
	}

	entries, err := dir.ReadDir(".")
	if err != nil {
		log.Fatalf("Could not find dep %q in vendor %s", project.Project.Name, err)
	}

	licenses := ""
	for _, entry := range entries {
		if !shouldIncludeFileInOsdf(entry, spdxCode) {
			continue
		}

		licenses += fmt.Sprintf("Contents of: %s/%s@%s\n\n", project.Project.Name, entry.Name, project.Project.Revision)

		text, err := dir.ReadFile(entry.Name)
		if err != nil {
			return licenses, err
		}

		licenses += string(text)
		licenses += "\n\n"
	}

	if licenses == "" {
		return "", errors.New("Could not find license text")
	}

	return licenses, nil
}

// A file should be included in the OSDF if it's a license, or the license is
// Apache and it's a notice.
func shouldIncludeFileInOsdf(file filer.File, spdxCode string) bool {
	if file.IsDir {
		return false
	}

	lowerName := strings.ToLower(file.Name)
	isLicense := strings.Contains(lowerName, "license")

	// We're required to include NOTICE files for Apache 2 licensed products.
	isNotice := lowerName == "notice" && spdxCode == "Apache-2.0"
	return isLicense || isNotice
}
