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

package migration

import (
	"bytes"
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/robertkrimen/otto"
)

const tileToEnvFunctionJs = `
function defnToEnv(defn) {
  var shouldRecurse = defn.type === 'collection' && typeof(defn.value) === 'object'
  if (!shouldRecurse) {
    return defn.value;
  }

  var collectionEntries = defn.value;
  var out = [];
  for (var i = 0; i < collectionEntries.length; i++) {
    var element = collectionEntries[i];
    var elementObj = {};

    for (var elementKey in element) {
      elementObj[elementKey] = defnToEnv(element[elementKey]);
    }

    out.push(elementObj);
  }

  return JSON.stringify(out, null, '  ')
}

function tile2env(properties) {
  var inner = properties["properties"];

  var env = {};
  for(var key in inner) {
    // skip non properties
    if(key.indexOf(".properties.") !== 0){
      continue;
    }

    // strip the .properties. prefix
    var envVar = key.substr(12).toUpperCase();
    var defn = inner[key];

    env[envVar] = '' + defnToEnv(defn);
  }

  return JSON.stringify(env);
}

properties = {{properties}};

{{migrationScript}}

tile2env(properties);
`

func GetMigrationResults(t *testing.T, tilePropertiesJSON, migrationScript string) map[string]string {
	t.Helper()
	vm := otto.New()

	script := tileToEnvFunctionJs
	replacer := strings.NewReplacer(
		"{{properties}}", tilePropertiesJSON,
		"{{migrationScript}}", migrationScript)
	script = replacer.Replace(script)

	t.Log(script)

	result, err := vm.Run(script)
	if err != nil {
		t.Fatal(err)
	}

	value, err := result.ToString()
	if err != nil {
		t.Fatal(err)
	}

	out := map[string]string{}
	if err := json.Unmarshal([]byte(value), &out); err != nil {
		t.Fatal(err)
	}

	return out
}

type MigrationTest struct {
	TileProperties string
	Migration      Migration
	ExpectedEnv    map[string]string
}

func (m *MigrationTest) Run(t *testing.T) {
	if !json.Valid([]byte(m.TileProperties)) {
		t.Fatalf("invalid tile properties")
	}

	t.Log("Running JS migration function")
	out := GetMigrationResults(t, m.TileProperties, m.Migration.TileScript)
	normalizeJson(t, out)
	normalizeJson(t, m.ExpectedEnv)

	if !reflect.DeepEqual(out, m.ExpectedEnv) {
		t.Logf("Diff %#v", DiffStringMap(out, m.ExpectedEnv))
		t.Fatalf("Expected: %#v Got: %#v", m.ExpectedEnv, out)
	}

	t.Log("Running go migration function")
	// runGo runs a no-op JS migration to get the environment variables as the
	// application would see them, then applies the go migration function to
	// ensure it produces the same result for people migrating by hand as the
	// JavaScript does for the tile users
	envFromTile := GetMigrationResults(t, m.TileProperties, "")
	if err := m.Migration.GoFunc(envFromTile); err != nil {
		t.Fatal(err)
	}

	normalizeJson(t, envFromTile)
	normalizeJson(t, m.ExpectedEnv)

	if !reflect.DeepEqual(envFromTile, m.ExpectedEnv) {
		t.Fatalf("Expected: %#v Got: %#v", m.ExpectedEnv, envFromTile)
	}
}

func normalizeJson(t *testing.T, env map[string]string) {
	for k, v := range env {
		if !json.Valid([]byte(v)) {
			continue
		}
		dst := &bytes.Buffer{}
		if err := json.Compact(dst, []byte(v)); err != nil {
			t.Errorf("couldn't compact %q: %v", k, err)
		}
		env[k] = dst.String()
	}
}

// TestNoOp is primarially to ensure the embedded JavaScript works.
func TestNoOp(t *testing.T) {
	cases := map[string]MigrationTest{
		"basic": {
			TileProperties: `
        {
          "properties": {
            ".properties.org": {
              "type": "string",
              "value": "system"
            },
            ".properties.space": {
              "type": "string",
              "value": "gcp-service-broker-space"
            }
          }
        }
        `,
			Migration: NoOp(),
			ExpectedEnv: map[string]string{
				"ORG":   "system",
				"SPACE": "gcp-service-broker-space",
			},
		},
		"empty": {
			TileProperties: `{"properties":{}}`,
			Migration:      NoOp(),
			ExpectedEnv:    map[string]string{},
		},
		"collection": {
			TileProperties: `
      {
        "properties": {
          ".properties.storage_custom_plans": {
            "type": "collection",
            "configurable": true,
            "credential": false,
            "value": [
              {
                "guid": {
                  "type": "uuid",
                  "value": "424b9afa-b31a-4517-a764-fbd137e58d3d"
                },
                "name": {
                  "type": "string",
                  "value": "custom-storage-plan"
                },
                "display_name": {
                  "type": "string",
                  "value": "custom cloud storage plan"
                },
                "description": {
                  "type": "string",
                  "value": "custom storage plan description"
                },
                "service": {
                  "type": "dropdown_select",
                  "value": "b9e4332e-b42b-4680-bda5-ea1506797474"
                },
                "storage_class": {
                  "type": "string",
                  "value": "standard"
                }
              }
            ],
            "optional": true
          }
        }
      }`,
			Migration: NoOp(),
			ExpectedEnv: map[string]string{
				"STORAGE_CUSTOM_PLANS": `[
  {
    "description": "custom storage plan description",
    "display_name": "custom cloud storage plan",
    "guid": "424b9afa-b31a-4517-a764-fbd137e58d3d",
    "name": "custom-storage-plan",
    "service": "b9e4332e-b42b-4680-bda5-ea1506797474",
    "storage_class": "standard"
  }
]`,
			},
		},
	}

	for tn, tc := range cases {
		t.Run(tn, tc.Run)
	}
}

func TestFullMigration(t *testing.T) {
	migration := FullMigration()

	t.Run("migration creates valid JS", func(t *testing.T) {
		vm := otto.New()
		fakeProps := map[string]interface{}{
			"properties": map[string]interface{}{},
		}
		if err := vm.Set("properties", fakeProps); err != nil {
			t.Fatal(err)
		}

		if _, err := vm.Run(migration.TileScript); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("go migration executes without failure", func(t *testing.T) {
		env := map[string]string{}
		if err := migration.GoFunc(env); err != nil {
			t.Fatal(err)
		}
	})
}
