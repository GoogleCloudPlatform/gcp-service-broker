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
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"sync"

	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/varcontext"
)

var (
	FsInitializationErr = errors.New("Filesystem must first be initialized.")
)

func NewWorkspace(variableContext *varcontext.VarContext, terraformTemplate string) (*TerraformWorkspace, error) {
	tfModule := ModuleDefinition{
		Name:       "brokertemplate",
		Definition: terraformTemplate,
	}

	inputList, err := tfModule.Inputs()
	if err != nil {
		return nil, err
	}

	limitedConfig := make(map[string]interface{})
	config := variableContext.ToMap()
	for _, name := range inputList {
		limitedConfig[name] = config[name]
	}

	workspace := TerraformWorkspace{
		Modules: []ModuleDefinition{tfModule},
		Instances: []ModuleInstance{
			{
				ModuleName:    tfModule.Name,
				InstanceName:  "instance",
				Configuration: limitedConfig,
			},
		},
	}

	return &workspace, nil
}

func DeserializeWorkspace(definition string) (*TerraformWorkspace, error) {
	ws := TerraformWorkspace{}
	if err := json.Unmarshal([]byte(definition), &ws); err != nil {
		return nil, err
	}

	return &ws, nil
}

type TerraformWorkspace struct {
	Environment map[string]string  `json:"-"` // GOOGLE_CREDENTIALS needs to be set to the JSON key and GOOGLE_PROJECT needs to be set to the project
	Modules     []ModuleDefinition `json:"modules"`
	Instances   []ModuleInstance   `json:"instances"`
	State       []byte             `json:"tfstate"`

	mux         sync.Mutex
	initialized bool
	dir         string
}

// Serialize converts the TerraformWorkspace into a JSON string.
func (workspace *TerraformWorkspace) Serialize() (string, error) {
	ws, err := json.Marshal(workspace)
	if err != nil {
		return "", err
	}

	return string(ws), nil
}

func (workspace *TerraformWorkspace) InitializeFs() error {
	workspace.mux.Lock()
	defer workspace.mux.Unlock()
	if workspace.initialized {
		return nil
	}

	// create a temp directory
	if dir, err := ioutil.TempDir("", "gsb"); err != nil {
		return err
	} else {
		workspace.dir = dir
	}

	// write the modulesTerraformWorkspace
	for _, module := range workspace.Modules {
		parent := path.Join(workspace.dir, module.Name)
		if err := os.Mkdir(parent, 0600); err != nil {
			return err
		}

		if err := ioutil.WriteFile(path.Join(parent, "definition.tf"), []byte(module.Definition), 0600); err != nil {
			return err
		}
	}

	// write the instances
	for _, instance := range workspace.Instances {
		contents, err := instance.MarshalDefinition()
		if err != nil {
			return err
		}

		if err := ioutil.WriteFile(path.Join(workspace.dir, instance.InstanceName+".tf"), contents, 0600); err != nil {
			return err
		}
	}

	// write the state if it exists
	if len(workspace.State) > 0 {
		if err := ioutil.WriteFile(workspace.tfStatePath(), workspace.State, 0600); err != nil {
			return err
		}
	}

	workspace.initialized = true
	return nil
}

func (workspace *TerraformWorkspace) TeardownFs() error {
	workspace.mux.Lock()
	defer workspace.mux.Unlock()

	if err := os.RemoveAll(workspace.dir); err != nil {
		return err
	}

	workspace.initialized = false
	return nil
}

func (workspace *TerraformWorkspace) Validate() error {
	// run tf validate
	return nil
}

func (workspace *TerraformWorkspace) Plan() (string, error) {
	workspace.mux.Lock()
	defer workspace.mux.Unlock()
	if !workspace.initialized {
		return "", FsInitializationErr
	}

	return "", nil
}

func (workspace *TerraformWorkspace) Outputs(instance string) (map[string]interface{}, error) {
	workspace.mux.Lock()
	defer workspace.mux.Unlock()

	if err := workspace.updateState(); err != nil {
		return nil, err
	}

	// TODO parse the state from the file and return it here.

	state := tfState{}
	if err := json.Unmarshal(workspace.State, &state); err != nil {
		return nil, err
	}

	// All root project modules get put under the "root" namespace
	module := state.GetModule("root", instance)
	if module == nil {
		return nil, fmt.Errorf("no instance exists with name %q", instance)
	}

	return module.GetOutputs(), nil
}

func (workspace *TerraformWorkspace) updateState() error {
	if !workspace.initialized {
		return nil
	}

	bytes, err := ioutil.ReadFile(workspace.tfStatePath())
	if err != nil {
		return err
	}

	workspace.State = bytes
	return nil
}

func (workspace *TerraformWorkspace) Apply() error {
	workspace.mux.Lock()
	defer workspace.mux.Unlock()
	if !workspace.initialized {
		return FsInitializationErr
	}

	return nil
}

func (workspace *TerraformWorkspace) Destroy() error {
	workspace.mux.Lock()
	defer workspace.mux.Unlock()
	if !workspace.initialized {
		return nil
	}

	if err := os.RemoveAll(workspace.dir); err != nil {
		return err
	}

	workspace.initialized = false
	return nil
}

func (workspace *TerraformWorkspace) tfStatePath() string {
	return path.Join(workspace.dir, "terraform.tfstate")
}

// tfState is a struct that can help us deserialize the tfstate JSON file.
type tfState struct {
	Version int        `json:"version"`
	Modules []tfModule `json:"modules"`
}

// GetModule gets a module at a given path or nil if none exists for that path.
func (state *tfState) GetModule(path ...string) *tfModule {
	for _, module := range state.Modules {
		if reflect.DeepEqual(module.Path, path) {
			return &module
		}
	}

	return nil
}

type tfModule struct {
	Path    []string `json:"path"`
	Outputs map[string]struct {
		Type  string      `json:"type"`
		Value interface{} `json:"value"`
	} `json:"outputs"`
}

// GetOutputs gets the key/value outputs defined for a module.
func (module *tfModule) GetOutputs() map[string]interface{} {
	out := make(map[string]interface{})

	for outputName, tfOutput := range module.Outputs {
		out[outputName] = tfOutput.Value
	}

	return out
}
