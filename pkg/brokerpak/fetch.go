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
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func fetch(url, dest string) error {
	fmt.Printf("downloading %q to %q\n", url, dest)

	_, err := os.Stat(dest)
	exists := !os.IsNotExist(err)
	if exists {
		fmt.Println("file already exists, skipping")
		return nil
	}

	// Setup local files first because it's cheaper to make these errors before
	// ones involving the network.
	if err := os.MkdirAll(filepath.Dir(dest), os.ModePerm); err != nil {
		return fmt.Errorf("error creating local directory %q: %v", filepath.Dir(dest), err)
	}
	out, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("error opening local file %q: %v", dest, err)
	}
	defer out.Close()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("error creating HTTP request: %v", err)
	}
	req.Header.Set("User-Agent", "gcp-service-broker")

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error getting HTTP resource: %v", err)
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("got unexpected HTTP response code: %d", response.StatusCode)
	}

	if _, err := io.Copy(out, response.Body); err != nil {
		return fmt.Errorf("error copying output: %v", err)
	}

	return nil
}
