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

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/src-d/go-license-detector.v2/licensedb/filer"
)

const ExampleLicense = `Copyright 2018 the Service Broker Project Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.`

func TestShouldIncludeFileInOsdf(t *testing.T) {
	cases := map[string]struct {
		FileName string
		SpdxCode string
		Expected bool
	}{
		"apache2 notice":    {"NOTICE", "Apache-2.0", true},
		"MIT notice":        {"NOTICE", "MIT", false},
		"upper license":     {"LICENSE", "MIT", true},
		"lower license":     {"license", "MIT", true},
		"mixed license":     {"License", "MIT", true},
		"extension license": {"License.txt", "MIT", true},
		"non-license":       {"main.go", "MIT", false},
		"license directory": {"license/", "MIT", false},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			fd := filer.File{IsDir: strings.HasSuffix(tc.FileName, "/"), Name: tc.FileName}
			if shouldIncludeFileInOsdf(fd, tc.SpdxCode) != tc.Expected {
				t.Error("Expected", tc.Expected, " Got", !tc.Expected)
			}
		})
	}
}

func TestMostLikelyLicense(t *testing.T) {
	cases := map[string]struct {
		LicenseList map[string]float32
		Expected    string
	}{
		"single":   {map[string]float32{"mit": 0}, "mit"},
		"multiple": {map[string]float32{"mit": 1, "bsd": 0}, "mit"},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			actual := mostLikelyLicense(tc.LicenseList)
			if tc.Expected != actual {
				t.Error("Expected", tc.Expected, " Got", actual)
			}
		})
	}
}

func ExampleDetectLicenses() {
	dir, err := ioutil.TempDir("", "lic")
	if err != nil {
		panic(err)
	}

	if err := ioutil.WriteFile(filepath.Join(dir, "LICENSE"), []byte(ExampleLicense), 0666); err != nil {
		panic(err)
	}

	lic, err := detectLicenses(dir)
	if err != nil {
		panic(err)
	}

	fmt.Println(lic)

	// Output: map[Apache-2.0:1]
}

func ExampleGetLicenseText() {
	dir, err := ioutil.TempDir("", "lic")
	if err != nil {
		panic(err)
	}

	proj := &Dependency{
		ParentDirectory: dir,
		Project: Project{
			Name:     "foo",
			Revision: "xxx-my-revision-here-xxx",
		},
	}

	if err := os.MkdirAll(proj.Directory(), 0777); err != nil {
		panic(err)
	}

	if err := ioutil.WriteFile(filepath.Join(proj.Directory(), "LICENSE"), []byte(ExampleLicense), 0666); err != nil {
		panic(err)
	}

	lic, err := getLicenseText(proj, "Apache-2.0")
	if err != nil {
		panic(err)
	}

	fmt.Println(lic)

	// Output: Contents of: foo/LICENSE@xxx-my-revision-here-xxx
	//
	// Copyright 2018 the Service Broker Project Authors.
	//
	// Licensed under the Apache License, Version 2.0 (the "License");
	// you may not use this file except in compliance with the License.
	// You may obtain a copy of the License at
	//
	// http://www.apache.org/licenses/LICENSE-2.0
	//
	// Unless required by applicable law or agreed to in writing, software
	// distributed under the License is distributed on an "AS IS" BASIS,
	// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	// See the License for the specific language governing permissions and
	// limitations under the License.
}

func ExampleGetProjects() {
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		panic(err)
	}

	lockFile := `
[[projects]]
  digest = "1:82f6c9a55c0bd9744064418f049d5232bb8b8cc45eb32e72c0adefaf158a0f9b"
  name = "github.com/hashicorp/go-safetemp"
  packages = ["."]
  pruneopts = ""
  revision = "c9a55de4fe06c920a71964b53cfe3dd293a3c743"
  version = "v1.0.0"

[[projects]]
  digest = "1:8c7fb7f81c06add10a17362abc1ae569ff9765a26c061c6b6e67c909f4f414db"
  name = "github.com/hashicorp/go-version"
  packages = ["."]
  pruneopts = ""
  revision = "b5a281d3160aa11950a6182bd9a9dc2cb1e02d50"
  version = "v1.0.0"
`

	if err := ioutil.WriteFile(filepath.Join(dir, "Gopkg.lock"), []byte(lockFile), 0666); err != nil {
		panic(err)
	}

	for _, p := range getProjects(dir) {
		fmt.Println(p.Name)
	}

	// Output: github.com/hashicorp/go-safetemp
	// github.com/hashicorp/go-version
}
