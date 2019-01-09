// Copyright 2019 the Service Broker Project Authors.
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

package brokers

import (
	"os"
	"reflect"
	"testing"

	"golang.org/x/oauth2/jwt"
)

const testServiceAccountJson = `{
  "type": "service_account",
  "project_id": "foo",
  "private_key_id": "something",
  "private_key": "foobar",
  "client_email": "example@gmail.com",
  "client_id": "1",
  "auth_uri": "somelink",
  "token_uri": "somelink",
  "auth_provider_x509_cert_url": "somelink",
  "client_x509_cert_url": "somelink"
}`

func TestNewBrokerConfigFromEnv(t *testing.T) {
	os.Setenv("ROOT_SERVICE_ACCOUNT_JSON", testServiceAccountJson)
	defer os.Unsetenv("ROOT_SERVICE_ACCOUNT_JSON")

	cfg, err := NewBrokerConfigFromEnv()
	if err != nil {
		t.Fatal(err)
	}

	t.Run("has-default-client", func(t *testing.T) {
		if cfg.HttpConfig == nil {
			t.Fatal("Expected HttpCofnig to be non-nil, got: <nil>")
		}

		if reflect.DeepEqual(cfg.HttpConfig, &jwt.Config{}) {
			t.Errorf("Expected HttpConfig to not be an empty JWT config, got: %#v", cfg.HttpConfig)
		}
	})

	t.Run("parsed-projectid-from-config", func(t *testing.T) {
		if !reflect.DeepEqual(cfg.ProjectId, "foo") {
			t.Errorf("Expected ProjectId to be %v, got: %v", "foo", cfg.ProjectId)
		}
	})
}
