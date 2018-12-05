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

package wrapper

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"reflect"
	"testing"
)

func TestTerraformWorkspace_Invariants(t *testing.T) {

	// This function tests the following two invariants of the workspace:
	// - The function updates the tfstate once finished.
	// - The function creates and destroys its own dir.

	cases := map[string]struct {
		Exec func(ws *TerraformWorkspace)
	}{
		"validate": {Exec: func(ws *TerraformWorkspace) {
			ws.Validate()
		}},
		"apply": {Exec: func(ws *TerraformWorkspace) {
			ws.Apply()
		}},
		"destroy": {Exec: func(ws *TerraformWorkspace) {
			ws.Destroy()
		}},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			// construct workspace
			ws, err := NewWorkspace(map[string]interface{}{}, ``)
			if err != nil {
				t.Fatal(err)
			}

			// substitute the executor so we can validate the state at the time of
			// "running" tf
			executorRan := false
			cmdDir := ""
			ws.Executor = func(cmd *exec.Cmd) error {
				executorRan = true
				cmdDir = cmd.Dir

				// validate that the directory exists
				_, err := os.Stat(cmd.Dir)
				if err != nil {
					t.Fatalf("couldn't stat the cmd execution dir %v", err)
				}

				// write dummy state file
				if err := ioutil.WriteFile(path.Join(cmdDir, "terraform.tfstate"), []byte(tn), 0755); err != nil {
					t.Fatal(err)
				}

				return nil
			}

			// run function
			tc.Exec(ws)

			// check validator got ran
			if !executorRan {
				t.Fatal("Executor did not get run as part of the function")
			}

			// check workspace destroyed
			if _, err := os.Stat(cmdDir); !os.IsNotExist(err) {
				t.Fatalf("command directory didn't %q get torn down %v", cmdDir, err)
			}

			// check tfstate updated
			if !reflect.DeepEqual(ws.State, []byte(tn)) {
				t.Fatalf("Expected state %v got %v", []byte(tn), ws.State)
			}
		})
	}
}

func TestCustomTerraformExecutor(t *testing.T) {
	customBinary := "/path/to/terraform"
	customPlugins := "/path/to/terraform-plugins"
	pluginsFlag := "-plugin-dir=" + customPlugins

	cases := map[string]struct {
		Input    *exec.Cmd
		Expected *exec.Cmd
	}{
		"destroy": {
			Input:    exec.Command("terraform", "destroy", "-auto-approve", "-no-color"),
			Expected: exec.Command(customBinary, "destroy", "-auto-approve", "-no-color"),
		},
		"apply": {
			Input:    exec.Command("terraform", "apply", "-auto-approve", "-no-color"),
			Expected: exec.Command(customBinary, "apply", "-auto-approve", "-no-color"),
		},
		"validate": {
			Input:    exec.Command("terraform", "validate", "-no-color"),
			Expected: exec.Command(customBinary, "validate", "-no-color"),
		},
		"init": {
			Input:    exec.Command("terraform", "init", "-no-color"),
			Expected: exec.Command(customBinary, "init", "-get-plugins=false", pluginsFlag, "-no-color"),
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			actual := exec.Command("!actual-never-got-called!")

			executor := CustomTerraformExecutor(customBinary, customPlugins, func(c *exec.Cmd) error {
				actual = c
				return nil
			})

			executor(tc.Input)

			if actual.Path != tc.Expected.Path {
				t.Errorf("path wasn't updated, expected: %q, actual: %q", tc.Expected.Path, actual.Path)
			}

			if !reflect.DeepEqual(actual.Args, tc.Expected.Args) {
				t.Errorf("args weren't updated correctly, expected: %#v, actual: %#v", tc.Expected.Args, actual.Args)
			}
		})
	}
}
