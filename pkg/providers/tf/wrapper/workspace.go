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
	"os/exec"
	"path"
	"strings"
	"sync"

	"code.cloudfoundry.org/lager"
)

// DefaultInstanceName is the default name of an instance of a particular module.
const (
	DefaultInstanceName = "instance"
)

var (
	FsInitializationErr = errors.New("Filesystem must first be initialized.")
)

// TerraformExecutor is the function that shells out to Terraform.
// It can intercept, modify or retry the given command.
type TerraformExecutor func(*exec.Cmd) error

// NewWorkspace creates a new TerraformWorkspace from a given template and variables to populate an instance of it.
// The created instance will have the name specified by the DefaultInstanceName constant.
func NewWorkspace(templateVars map[string]interface{}, terraformTemplate string) (*TerraformWorkspace, error) {
	tfModule := ModuleDefinition{
		Name:       "brokertemplate",
		Definition: terraformTemplate,
	}

	inputList, err := tfModule.Inputs()
	if err != nil {
		return nil, err
	}

	limitedConfig := make(map[string]interface{})
	for _, name := range inputList {
		limitedConfig[name] = templateVars[name]
	}

	workspace := TerraformWorkspace{
		Modules: []ModuleDefinition{tfModule},
		Instances: []ModuleInstance{
			{
				ModuleName:    tfModule.Name,
				InstanceName:  DefaultInstanceName,
				Configuration: limitedConfig,
			},
		},
	}

	return &workspace, nil
}

// DeserializeWorkspace creates a new TerraformWorkspace from a given JSON
// serialization of one.
func DeserializeWorkspace(definition string) (*TerraformWorkspace, error) {
	ws := TerraformWorkspace{}
	if err := json.Unmarshal([]byte(definition), &ws); err != nil {
		return nil, err
	}

	return &ws, nil
}

// TerraformWorkspace represents the directory layout of a Terraform execution.
// The structure is strict, consiting of several Terraform modules and instances
// of those modules. The strictness is artificial, but maintains a clear
// separation between data and code.
//
// It manages the directory structure needed for the commands, serializing and
// deserializing Terraform state, and all the flags necessary to call Terraform.
//
// All public functions that shell out to Terraform maintain the following invariants:
// - The function blocks if another terraform shell is running.
// - The function updates the tfstate once finished.
// - The function creates and destroys its own dir.
type TerraformWorkspace struct {
	Environment map[string]string  `json:"-"` // GOOGLE_CREDENTIALS needs to be set to the JSON key and GOOGLE_PROJECT needs to be set to the project
	Modules     []ModuleDefinition `json:"modules"`
	Instances   []ModuleInstance   `json:"instances"`
	State       []byte             `json:"tfstate"`

	// Executor is a function that gets invoked to shell out to Terraform.
	// If left nil, the default executor is used.
	Executor TerraformExecutor `json:"-"`

	dirLock sync.Mutex
	dir     string
}

// String returns a human-friendly representation of the workspace suitable for
// printing to the console.
func (workspace *TerraformWorkspace) String() string {
	var b strings.Builder

	b.WriteString("# Terraform Workspace\n")
	fmt.Fprintf(&b, "modules: %d\n", len(workspace.Modules))
	fmt.Fprintf(&b, "instances: %d\n", len(workspace.Instances))
	fmt.Fprintln(&b)

	for _, instance := range workspace.Instances {
		fmt.Fprintf(&b, "## Instance %q\n", instance.InstanceName)
		fmt.Fprintf(&b, "module = %q\n", instance.ModuleName)

		for k, v := range instance.Configuration {
			fmt.Fprintf(&b, "input.%s = %#v\n", k, v)
		}

		if outputs, err := workspace.Outputs(instance.InstanceName); err != nil {
			for k, v := range outputs {
				fmt.Fprintf(&b, "output.%s = %#v\n", k, v)
			}
		}
		fmt.Fprintln(&b)
	}

	return b.String()
}

// Serialize converts the TerraformWorkspace into a JSON string.
func (workspace *TerraformWorkspace) Serialize() (string, error) {
	ws, err := json.Marshal(workspace)
	if err != nil {
		return "", err
	}

	return string(ws), nil
}

// initializeFs initializes the filesystem directory necessary to run Terraform.
func (workspace *TerraformWorkspace) initializeFs() error {
	workspace.dirLock.Lock()
	// create a temp directory
	if dir, err := ioutil.TempDir("", "gsb"); err != nil {
		return err
	} else {
		workspace.dir = dir
	}

	// write the modulesTerraformWorkspace
	for _, module := range workspace.Modules {
		parent := path.Join(workspace.dir, module.Name)
		if err := os.Mkdir(parent, 0755); err != nil {
			return err
		}

		if err := ioutil.WriteFile(path.Join(parent, "definition.tf"), []byte(module.Definition), 0755); err != nil {
			return err
		}
	}

	// write the instances
	for _, instance := range workspace.Instances {
		contents, err := instance.MarshalDefinition()
		if err != nil {
			return err
		}

		if err := ioutil.WriteFile(path.Join(workspace.dir, instance.InstanceName+".tf"), contents, 0755); err != nil {
			return err
		}
	}

	// write the state if it exists
	if len(workspace.State) > 0 {
		if err := ioutil.WriteFile(workspace.tfStatePath(), workspace.State, 0755); err != nil {
			return err
		}
	}

	// run "terraform init"
	if err := workspace.runTf("init", "-no-color"); err != nil {
		return err
	}

	return nil
}

// TeardownFs removes the directory we executed Terraform in and updates the
// state from it.
func (workspace *TerraformWorkspace) teardownFs() error {
	bytes, err := ioutil.ReadFile(workspace.tfStatePath())
	if err != nil {
		return err
	}

	workspace.State = bytes

	if err := os.RemoveAll(workspace.dir); err != nil {
		return err
	}

	workspace.dir = ""
	workspace.dirLock.Unlock()
	return nil
}

// Outputs gets the Terraform outputs from the state for the instance with the
// given name. This function DOES NOT invoke Terraform and instead uses the stored state.
func (workspace *TerraformWorkspace) Outputs(instance string) (map[string]interface{}, error) {
	state, err := NewTfstate(workspace.State)
	if err != nil {
		return nil, err
	}

	// All root project modules get put under the "root" namespace
	module := state.GetModule("root", instance)
	if module == nil {
		return nil, fmt.Errorf("no instance exists with name %q", instance)
	}

	return module.GetOutputs(), nil
}

// Validate runs `terraform Validate` on this workspace.
// This funciton blocks if another Terraform command is running on this workspace.
func (workspace *TerraformWorkspace) Validate() error {
	err := workspace.initializeFs()
	defer workspace.teardownFs()
	if err != nil {
		return err
	}

	return workspace.runTf("validate", "-no-color")
}

// Apply runs `terraform apply` on this workspace.
// This funciton blocks if another Terraform command is running on this workspace.
func (workspace *TerraformWorkspace) Apply() error {
	err := workspace.initializeFs()
	defer workspace.teardownFs()
	if err != nil {
		return err
	}

	return workspace.runTf("apply", "-auto-approve", "-no-color")
}

// Destroy runs `terraform destroy` on this workspace.
// This funciton blocks if another Terraform command is running on this workspace.
func (workspace *TerraformWorkspace) Destroy() error {
	err := workspace.initializeFs()
	defer workspace.teardownFs()
	if err != nil {
		return err
	}

	return workspace.runTf("destroy", "-auto-approve", "-no-color")
}

func (workspace *TerraformWorkspace) tfStatePath() string {
	return path.Join(workspace.dir, "terraform.tfstate")
}

func (workspace *TerraformWorkspace) runTf(subCommand string, args ...string) error {
	sub := []string{subCommand}
	sub = append(sub, args...)

	env := os.Environ()
	for k, v := range workspace.Environment {
		env = append(env, fmt.Sprintf("%s=%s", k, v))
	}

	c := exec.Command("terraform", sub...)
	c.Env = env
	c.Dir = workspace.dir

	executor := defaultExecutor
	if workspace.Executor != nil {
		executor = workspace.Executor
	}

	return executor(c)
}

// CustomTerraformExecutor executes a custom Terraform binary that uses plugins
// from a given plugin directory rather than the Terraform that's on the PATH
// and downloading the binaries from the web.
func CustomTerraformExecutor(tfBinaryPath, tfPluginDir string) TerraformExecutor {
	return func(c *exec.Cmd) error {
		c.Path = tfBinaryPath

		// Add the -get-plugins=false and -plugin-dir={tfPluginDir} after the
		// sub-command to force Terraform to use a particular plugin.
		subCommand := c.Args[0]
		oldFlags := c.Args[1:]
		newArgs := []string{subCommand, "-get-plugins=false", fmt.Sprintf("-plugin-dir=%s", tfPluginDir)}
		c.Args = append(newArgs, oldFlags...)

		return defaultExecutor(c)
	}
}

func defaultExecutor(c *exec.Cmd) error {
	logger := lager.NewLogger("terraform@" + c.Dir)
	logger.RegisterSink(lager.NewWriterSink(os.Stderr, lager.ERROR))
	logger.RegisterSink(lager.NewWriterSink(os.Stdout, lager.DEBUG))

	logger.Info("starting process", lager.Data{
		"path": c.Path,
		"args": c.Args,
		"dir":  c.Dir,
	})
	output, err := c.CombinedOutput()
	logger.Info("results", lager.Data{
		"output": string(output),
		"error":  err,
	})

	return err
}
