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
	"os"
	"path/filepath"
)

func makeParents(dest string) error {
	if err := os.MkdirAll(filepath.Dir(dest), os.ModePerm); err != nil {
		return fmt.Errorf("error creating local directory %q: %v", filepath.Dir(dest), err)
	}

	return nil
}

func cp(from, to string) error {
	in, err := os.Open(from)
	if err != nil {
		return err
	}
	defer in.Close()

	return cpReader(in, to)
}

func cpReader(from io.Reader, to string) error {
	if err := makeParents(to); err != nil {
		return err
	}

	out, err := os.Create(to)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err = io.Copy(out, from); err != nil {
		return err
	}
	return out.Close()
}
