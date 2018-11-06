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

/*
Package wrapper provides a way to set up, tear down, and execute Terraform
programatically.
*/
package wrapper

func NewTerraformProxy(projectId, serviceAccountKey string) (*TerraformProxy, error) {
	return &TerraformProxy{}, nil
}

type TerraformProxy struct {
}

type TerraformBuilder struct {
	// GOOGLE_CREDENTIALS
	// GOOGLE_APPLICATION_CREDENTIALS
	// GOOGLE_PROJECT
}

type TerraformWorkspace struct {
	// directory
}

func (builder *TerraformBuilder) AddModule(name, definition string) {

}

func (builder *TerraformBuilder) AddInstance(module, name, string, vars interface{}) {

}

func (builder *TerraformBuilder) Build() (*TerraformWorkspace, error) {
	// create temp directory
	// create provider
	// create modules
	// re-hydrate state (if available)
	// create variables
	// register modules
	return nil, nil
}

// Initialize(project, serviceAccountKey, )
// .AddModule(Module)
// .AddInstance(Module, name, vars)
// .PackState()
// .UnpackState()
//
// // Shell out to TF
// .Plan()
// .Apply()
// .Destroy()
//
// Module(name, definition)
// .Inputs()  -> []string
// .Outputs()  -> []string
