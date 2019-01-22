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
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/robertkrimen/otto"
)

const tileToEnvFunctionJs = `
function defnToEnv(defn) {
  if (defn.type !== 'collection') {
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
	t.Log("Running JS migration function")
	m.runJs(t)

	t.Log("Running go migration function")
	m.runGo(t)
}

func (m *MigrationTest) runJs(t *testing.T) {
	out := GetMigrationResults(t, m.TileProperties, m.Migration.TileScript)

	if !reflect.DeepEqual(out, m.ExpectedEnv) {
		t.Errorf("Expected: %v Got: %v", m.ExpectedEnv, out)
	}
}

func (m *MigrationTest) runGo(t *testing.T) {
	envFromTile := GetMigrationResults(t, m.TileProperties, "")
	m.Migration.GoFunc(envFromTile)

	if !reflect.DeepEqual(envFromTile, m.ExpectedEnv) {
		t.Errorf("Expected: %v Got: %v", m.ExpectedEnv, envFromTile)
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
              "configurable": true,
              "credential": false,
              "value": "system",
              "optional": false
            },
            ".properties.space": {
              "type": "string",
              "configurable": true,
              "credential": false,
              "value": "gcp-service-broker-space",
              "optional": false
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
                  "configurable": false,
                  "credential": false,
                  "value": "424b9afa-b31a-4517-a764-fbd137e58d3d",
                  "optional": false
                },
                "name": {
                  "type": "string",
                  "configurable": true,
                  "credential": false,
                  "value": "custom-storage-plan",
                  "optional": false
                },
                "display_name": {
                  "type": "string",
                  "configurable": true,
                  "credential": false,
                  "value": "custom cloud storage plan",
                  "optional": false
                },
                "description": {
                  "type": "string",
                  "configurable": true,
                  "credential": false,
                  "value": "custom storage plan description",
                  "optional": false
                },
                "service": {
                  "type": "dropdown_select",
                  "configurable": false,
                  "credential": false,
                  "value": "b9e4332e-b42b-4680-bda5-ea1506797474",
                  "optional": false,
                  "options": [
                    {
                      "label": "Google Cloud Storage",
                      "value": "b9e4332e-b42b-4680-bda5-ea1506797474"
                    }
                  ]
                },
                "storage_class": {
                  "type": "string",
                  "configurable": true,
                  "credential": false,
                  "value": "standard",
                  "optional": false
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

func TestDeleteWhitelistKeys(t *testing.T) {
	cases := map[string]MigrationTest{
		"no-overlap": {
			TileProperties: `
        {
          "properties": {
            ".properties.org": {
              "type": "string",
              "configurable": true,
              "credential": false,
              "value": "system",
              "optional": false
            }
          }
        }`,
			Migration:   DeleteWhitelistKeys(),
			ExpectedEnv: map[string]string{"ORG": "system"},
		},
		"one-key": {
			TileProperties: `
      {
        "properties": {
          ".properties.org": {"type": "string", "value": "system"},
          ".properties.gsb_service_google_bigquery_whitelist": { "value": 1 },
        }
      }`,
			Migration:   DeleteWhitelistKeys(),
			ExpectedEnv: map[string]string{"ORG": "system"},
		},
		"all-keys": {
			TileProperties: `
      {
        "properties": {
          ".properties.gsb_service_google_bigquery_whitelist": { "value": 1 },
          ".properties.gsb_service_google_bigtable_whitelist": { "value": 1 },
          ".properties.gsb_service_google_cloudsql_mysql_whitelist": { "value": 1 },
          ".properties.gsb_service_google_cloudsql_postgres_whitelist": { "value": 1 },
          ".properties.gsb_service_google_ml_apis_whitelist": { "value": 1 },
          ".properties.gsb_service_google_pubsub_whitelist": { "value": 1 },
          ".properties.gsb_service_google_spanner_whitelist": { "value": 1 },
          ".properties.gsb_service_google_storage_whitelist": { "value": 1 },
        }
      }`,
			Migration:   DeleteWhitelistKeys(),
			ExpectedEnv: map[string]string{},
		},
	}

	for tn, tc := range cases {
		t.Run(tn, tc.Run)
	}
}

func TestFullMigration(t *testing.T) {
	cases := map[string]MigrationTest{
		"no-fail-on-empty-properties": {
			TileProperties: `
        {
          "properties": {}
        }`,
			Migration:   FullMigration(),
			ExpectedEnv: map[string]string{},
		},
	}

	for tn, tc := range cases {
		t.Run(tn, tc.Run)
	}
}
