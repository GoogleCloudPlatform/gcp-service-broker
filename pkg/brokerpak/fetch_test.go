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

package brokerpak

import (
	"testing"

	getter "github.com/hashicorp/go-getter"
)

func TestDefaultGetters(t *testing.T) {
	getters := defaultGetters()

	// gcs SHOULD NOT be in there
	for _, prefix := range []string{"gs", "gcs"} {
		if getters[prefix] != nil {
			t.Errorf("expected default getters not to contain %q", prefix)
		}
	}

	for _, prefix := range []string{"http", "https", "git", "hg"} {
		if getters[prefix] == nil {
			t.Errorf("expected default getters not to contain %q", prefix)
		}
	}
}

func TestNewFileGetterClient(t *testing.T) {
	source := "http://www.example.com/foo/bar/bazz"
	dest := "/tmp/path/to/dest"
	client := newFileGetterClient(source, dest)

	if client.Src != source {
		t.Errorf("Expected Src to be %q got %q", source, client.Src)
	}

	if client.Dst != dest {
		t.Errorf("Expected Dst to be %q got %q", dest, client.Dst)
	}

	if client.Mode != getter.ClientModeFile {
		t.Errorf("Expected Dst to be %q got %q", getter.ClientModeFile, client.Mode)
	}

	if client.Getters == nil {
		t.Errorf("Expected getters to be set")
	}

	if client.Decompressors == nil {
		t.Errorf("Expected decompressors to be set")
	}

	if len(client.Decompressors) != 0 {
		t.Errorf("Expected decompressors to be empty")
	}
}
