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

package name_generator

import (
	"fmt"
	"time"

	"crypto/rand"
	"encoding/base64"
	"sync"
)

var (
	Basic BasicInstance
	Sql   SqlInstance
	once  sync.Once
)

type SqlInstance interface {
	BasicInstance
	DatabaseName() string
	GenerateUsername(instanceID, bindingID string) (string, error)
	GeneratePassword() (string, error)
}

type BasicInstance interface {
	InstanceName() string
	InstanceNameWithSeparator(sep string) string
}

func New() (BasicInstance, SqlInstance) {
	once.Do(func() {
		Basic = &BasicNameGenerator{}
		Sql = &SqlNameGenerator{}
	})
	return Basic, Sql
}

type BasicNameGenerator struct {
	count int
}

type SqlNameGenerator struct {
	BasicNameGenerator
}

func (bng *BasicNameGenerator) newNameWithSeperator(sep string) string {
	bng.count += 1
	return fmt.Sprintf("pcf%ssb%s%d%s%d", sep, sep, bng.count, sep, time.Now().UnixNano())
}

func (bng *BasicNameGenerator) InstanceName() string {
	return bng.newNameWithSeperator("_")
}

func (bng *BasicNameGenerator) InstanceNameWithSeparator(sep string) string {
	return bng.newNameWithSeperator(sep)
}

func (sng *SqlNameGenerator) InstanceName() string {
	return sng.newNameWithSeperator("-")
}

func (sng *SqlNameGenerator) DatabaseName() string {
	return sng.InstanceName()
}

const (
	maxUsernameLength       = 16 // Limit from http://dev.mysql.com/doc/refman/5.7/en/user-names.html
	generatedPasswordLength = 32
)

func (*SqlNameGenerator) GenerateUsername(instanceID, bindingID string) (string, error) {
	if len(instanceID)+len(bindingID) == 0 {
		return "", fmt.Errorf("empty instanceID and bindingID")
	}

	username := bindingID + instanceID
	if len(username) > maxUsernameLength {
		username = username[:maxUsernameLength]
	}

	return username, nil
}

func (*SqlNameGenerator) GeneratePassword() (string, error) {
	rb := make([]byte, generatedPasswordLength)
	_, err := rand.Read(rb)

	if err != nil {
		return "", err
	}

	rs := base64.URLEncoding.EncodeToString(rb)

	return rs, nil
}
