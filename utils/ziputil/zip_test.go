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

package ziputil

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/GoogleCloudPlatform/gcp-service-broker/utils/stream"
)

// A zip file with a single file with the contents "Hello, world!":
// MODE        SIZE  NAME
// -rw-rw-rw-  13    path/to/my/file.txt
var exampleFile = []byte{
	0x50, 0x4b, 0x3, 0x4, 0x14, 0x0, 0x8, 0x0, 0x8, 0x0, 0x0, 0x0, 0x0, 0x0,
	0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x13, 0x0,
	0x0, 0x0, 0x70, 0x61, 0x74, 0x68, 0x2f, 0x74, 0x6f, 0x2f, 0x6d, 0x79,
	0x2f, 0x66, 0x69, 0x6c, 0x65, 0x2e, 0x74, 0x78, 0x74, 0xf2, 0x48, 0xcd,
	0xc9, 0xc9, 0xd7, 0x51, 0x28, 0xcf, 0x2f, 0xca, 0x49, 0x51, 0x4, 0x4, 0x0,
	0x0, 0xff, 0xff, 0x50, 0x4b, 0x7, 0x8, 0xe6, 0xc6, 0xe6, 0xeb, 0x13, 0x0,
	0x0, 0x0, 0xd, 0x0, 0x0, 0x0, 0x50, 0x4b, 0x1, 0x2, 0x14, 0x0, 0x14, 0x0,
	0x8, 0x0, 0x8, 0x0, 0x0, 0x0, 0x0, 0x0, 0xe6, 0xc6, 0xe6, 0xeb, 0x13, 0x0,
	0x0, 0x0, 0xd, 0x0, 0x0, 0x0, 0x13, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
	0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x70, 0x61, 0x74, 0x68, 0x2f,
	0x74, 0x6f, 0x2f, 0x6d, 0x79, 0x2f, 0x66, 0x69, 0x6c, 0x65, 0x2e, 0x74, 0x78,
	0x74, 0x50, 0x4b, 0x5, 0x6, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0, 0x1, 0x0, 0x41,
	0x0, 0x0, 0x0, 0x54, 0x0, 0x0, 0x0, 0x0, 0x0}

func openZipfile() *zip.Reader {
	r, _ := zip.NewReader(bytes.NewReader(exampleFile), int64(len(exampleFile)))
	return r
}

func ExampleList() {
	zf := openZipfile()
	List(zf, os.Stdout)

	// Output: MODE        SIZE  NAME
	// -rw-rw-rw-  13    path/to/my/file.txt
}

func ExampleFind() {
	zf := openZipfile()
	fmt.Println(Find(zf, "does", "not", "exist"))
	fmt.Println(Find(zf, "path/to/my", "file.txt").Name)

	// Output: <nil>
	// path/to/my/file.txt
}

func ExampleOpen() {
	zf := openZipfile()
	source := stream.FromReadCloserError(Open(zf, "path/to/my/file.txt"))
	stream.Copy(source, stream.ToWriter(os.Stdout))

	// Output: Hello, world!
}

func ExampleExtract() {
	tmp, err := ioutil.TempDir("", "ziptest")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(tmp)

	zf := openZipfile()

	// extract all files/dirs under the path/to directory to tmp
	if err := Extract(zf, "path/to", tmp); err != nil {
		panic(err)
	}

	stream.Copy(stream.FromFile(tmp, "my", "file.txt"), stream.ToWriter(os.Stdout))

	// Output: Hello, world!
}

func TestClean(t *testing.T) {
	cases := map[string]struct {
		Case     []string
		Expected string
	}{
		"blank": {
			Case:     []string{},
			Expected: "",
		},
		"absolute": {
			Case:     []string{"/foo", "bar"},
			Expected: "foo/bar",
		},
		"relative-dot": {
			Case:     []string{"./foo", "bar"},
			Expected: "foo/bar",
		},
		"backslash": {
			// This case will ONLY test validity on Windows based machines.
			Case:     []string{filepath.Join("windows", "system32")},
			Expected: "windows/system32",
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			actual := Clean(tc.Case...)
			if actual != tc.Expected {
				t.Errorf("expected: %q, actual: %q", tc.Expected, actual)
			}
		})
	}
}

// This is a modified version of Golang's net/http/fs_test.go:TestServeFile_DotDot
func TestContainsDotDot(t *testing.T) {
	tests := []struct {
		req      string
		expected bool
	}{
		{"/testdata/file", false},
		{"/../file", true},
		{"/..", true},
		{"/../", true},
		{"/../foo", true},
		{"/..\\foo", true},
		{"/file/a", false},
		{"/file/a..", false},
		{"/file/a/..", true},
		{"/file/a\\..", true},
	}
	for _, tt := range tests {
		actual := containsDotDot(tt.req)
		if actual != tt.expected {
			t.Errorf("expected containsDotDot to be %t got %t", tt.expected, actual)
		}
	}
}
