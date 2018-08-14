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

package fakes

type StaticNameGenerator struct {
	Val string
}

func (sg *StaticNameGenerator) InstanceName() string {
	return sg.Val
}

func (sg *StaticNameGenerator) InstanceNameWithSeparator(sep string) string {
	return sg.Val
}

func (sg *StaticNameGenerator) DatabaseName() string {
	return sg.Val
}

type StaticSQLNameGenerator struct {
	StaticNameGenerator
}

func (sng *StaticSQLNameGenerator) InstanceName() string {
	return sng.Val
}

func (sng *StaticSQLNameGenerator) DatabaseName() string {
	return sng.Val
}

func (sng *StaticSQLNameGenerator) GenerateUsername(instanceID, bindingID string) (string, error) {
	return sng.Val[:16], nil
}

func (sng *StaticSQLNameGenerator) GeneratePassword() (string, error) {
	return sng.Val, nil
}
