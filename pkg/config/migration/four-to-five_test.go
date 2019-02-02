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

import "testing"

const defaultEmptyGsbServiceConfig = `{
  "00b9ca4a-7cd6-406a-a5b7-2f43f41ade75": {
    "//": "Builtin STACKDRIVER_PROFILER",
    "bind_defaults": {},
    "custom_plans": [],
    "disabled": false,
    "provision_defaults": {}
  },
  "2bc0d9ed-3f68-4056-b842-4a85cfbc727f": {
    "//": "Builtin STACKDRIVER_MONITORING",
    "bind_defaults": {},
    "custom_plans": [],
    "disabled": false,
    "provision_defaults": {}
  },
  "3e897eb3-9062-4966-bd4f-85bda0f73b3d": {
    "//": "Builtin DATAFLOW",
    "bind_defaults": {},
    "custom_plans": [],
    "disabled": false,
    "provision_defaults": {}
  },
  "4bc59b9a-8520-409f-85da-1c7552315863": {
    "//": "Builtin CLOUDSQL_MYSQL",
    "bind_defaults": {},
    "custom_plans": [],
    "disabled": false,
    "provision_defaults": {}
  },
  "51b3e27e-d323-49ce-8c5f-1211e6409e82": {
    "//": "Builtin SPANNER",
    "bind_defaults": {},
    "custom_plans": [],
    "disabled": false,
    "provision_defaults": {}
  },
  "5ad2dce0-51f7-4ede-8b46-293d6df1e8d4": {
    "//": "Builtin ML_APIS",
    "bind_defaults": {},
    "custom_plans": [],
    "disabled": false,
    "provision_defaults": {}
  },
  "628629e3-79f5-4255-b981-d14c6c7856be": {
    "//": "Builtin PUBSUB",
    "bind_defaults": {},
    "custom_plans": [],
    "disabled": false,
    "provision_defaults": {}
  },
  "76d4abb2-fee7-4c8f-aee1-bcea2837f02b": {
    "//": "Builtin DATASTORE",
    "bind_defaults": {},
    "custom_plans": [],
    "disabled": false,
    "provision_defaults": {}
  },
  "83837945-1547-41e0-b661-ea31d76eed11": {
    "//": "Builtin STACKDRIVER_DEBUGGER",
    "bind_defaults": {},
    "custom_plans": [],
    "disabled": false,
    "provision_defaults": {}
  },
  "a2b7b873-1e34-4530-8a42-902ff7d66b43": {
    "//": "Builtin FIRESTORE",
    "bind_defaults": {},
    "custom_plans": [],
    "disabled": false,
    "provision_defaults": {}
  },
  "b8e19880-ac58-42ef-b033-f7cd9c94d1fe": {
    "//": "Builtin BIGTABLE",
    "bind_defaults": {},
    "custom_plans": [],
    "disabled": false,
    "provision_defaults": {}
  },
  "b9e4332e-b42b-4680-bda5-ea1506797474": {
    "//": "Builtin STORAGE",
    "bind_defaults": {},
    "custom_plans": [],
    "disabled": false,
    "provision_defaults": {}
  },
  "c5ddfe15-24d9-47f8-8ffe-f6b7daa9cf4a": {
    "//": "Builtin STACKDRIVER_TRACE",
    "bind_defaults": {},
    "custom_plans": [],
    "disabled": false,
    "provision_defaults": {}
  },
  "cbad6d78-a73c-432d-b8ff-b219a17a803a": {
    "//": "Builtin CLOUDSQL_POSTGRES",
    "bind_defaults": {},
    "custom_plans": [],
    "disabled": false,
    "provision_defaults": {}
  },
  "e84b69db-3de9-4688-8f5c-26b9d5b1f129": {
    "//": "Builtin DIALOGFLOW",
    "bind_defaults": {},
    "custom_plans": [],
    "disabled": false,
    "provision_defaults": {}
  },
  "f80c0a3e-bd4d-4809-a900-b4e33a6450f1": {
    "//": "Builtin BIGQUERY",
    "bind_defaults": {},
    "custom_plans": [],
    "disabled": false,
    "provision_defaults": {}
  }
}`

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
          ".properties.gsb_service_google_bigquery_whitelist": { "value": 1 }
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
          ".properties.gsb_service_google_storage_whitelist": { "value": 1 }
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

func TestCollapseCustomPlans(t *testing.T) {
	cases := map[string]MigrationTest{
		"no-overlap": {
			TileProperties: `
        {
          "properties": {
            ".properties.org": {"type": "string","value": "system"}
          }
        }`,
			Migration:   CollapseCustomPlans(),
			ExpectedEnv: map[string]string{"ORG": "system"},
		},
		"multiple-plans": {
			TileProperties: `
      {
        "properties": {
          ".properties.cloudsql_mysql_custom_plans": {
            "type": "collection",
            "value": [{"name": {"type": "string", "value": "plan1"}},{"name": {"type": "string", "value": "plan2"}}]
          }
        }
      }
`,
			Migration: CollapseCustomPlans(),
			ExpectedEnv: map[string]string{
				"CLOUDSQL_MYSQL_CUSTOM_PLANS": `[{"name": "plan1"},{"name": "plan2"}]`,
			},
		},
		"real-example": {
			TileProperties: `{
      	"properties": {
      		".properties.cloudsql_mysql_custom_plans": {
      			"type": "collection",
      			"value": [
      				{
      					"guid": {"type": "uuid", "value": "69b4b7c9-2175-4d00-a298-81ebf11de59f"},
      					"name": {"type": "string", "value": "custom-mysql-plan"},
      					"display_name": {"type": "string", "value": "custom mysql plan"},
      					"description": {"type": "string", "value": "custom mysql plan description"},
      					"service": {"type": "dropdown_select", "value": "4bc59b9a-8520-409f-85da-1c7552315863"},
      					"tier": {"type": "string", "value": "db-n1-standard-1"},
      					"pricing_plan": {"type": "dropdown_select", "value": "PER_USE"},
      					"max_disk_size": {"type": "string", "value": "10"}
      				}
      			]
      		}
      	}
      }`,
			Migration: CollapseCustomPlans(),
			ExpectedEnv: map[string]string{
				"CLOUDSQL_MYSQL_CUSTOM_PLANS": `[
          {
            "description": "custom mysql plan description",
            "display_name": "custom mysql plan",
            "guid": "69b4b7c9-2175-4d00-a298-81ebf11de59f",
            "max_disk_size": "10",
            "name": "custom-mysql-plan",
            "pricing_plan": "PER_USE",
            "service": "4bc59b9a-8520-409f-85da-1c7552315863",
            "tier": "db-n1-standard-1"
            }
          ]`,
			},
		},
		"all-4x-services": {
			TileProperties: `
      {
        "properties": {
          ".properties.bigquery_custom_plans": { "type": "collection", "value": [{"name": {"type": "string", "value": "plan1"}}]},
          ".properties.bigtable_custom_plans": { "type": "collection", "value": [{"name": {"type": "string", "value": "plan1"}}]},
          ".properties.cloudsql_mysql_custom_plans": { "type": "collection", "value": [{"name": {"type": "string", "value": "plan1"}}]},
          ".properties.cloudsql_postgres_custom_plans": { "type": "collection", "value": [{"name": {"type": "string", "value": "plan1"}}]},
          ".properties.dataflow_custom_plans": { "type": "collection", "value": [{"name": {"type": "string", "value": "plan1"}}]},
          ".properties.datastore_custom_plans": { "type": "collection", "value": [{"name": {"type": "string", "value": "plan1"}}]},
          ".properties.dialogflow_custom_plans": { "type": "collection", "value": [{"name": {"type": "string", "value": "plan1"}}]},
          ".properties.firestore_custom_plans": { "type": "collection", "value": [{"name": {"type": "string", "value": "plan1"}}]},
          ".properties.ml_apis_custom_plans": { "type": "collection", "value": [{"name": {"type": "string", "value": "plan1"}}]},
          ".properties.pubsub_custom_plans": { "type": "collection", "value": [{"name": {"type": "string", "value": "plan1"}}]},
          ".properties.spanner_custom_plans": { "type": "collection", "value": [{"name": {"type": "string", "value": "plan1"}}]},
          ".properties.stackdriver_debugger_custom_plans": { "type": "collection", "value": [{"name": {"type": "string", "value": "plan1"}}]},
          ".properties.stackdriver_monitoring_custom_plans": { "type": "collection", "value": [{"name": {"type": "string", "value": "plan1"}}]},
          ".properties.stackdriver_profiler_custom_plans": { "type": "collection", "value": [{"name": {"type": "string", "value": "plan1"}}]},
          ".properties.stackdriver_trace_custom_plans": { "type": "collection", "value": [{"name": {"type": "string", "value": "plan1"}}]},
          ".properties.storage_custom_plans": { "type": "collection", "value": [{"name": {"type": "string", "value": "plan1"}}]}
        }
      }
`,
			Migration: CollapseCustomPlans(),
			ExpectedEnv: map[string]string{
				"BIGQUERY_CUSTOM_PLANS":               `[{"name":"plan1"}]`,
				"BIGTABLE_CUSTOM_PLANS":               `[{"name":"plan1"}]`,
				"CLOUDSQL_MYSQL_CUSTOM_PLANS":         `[{"name":"plan1"}]`,
				"CLOUDSQL_POSTGRES_CUSTOM_PLANS":      `[{"name":"plan1"}]`,
				"DATAFLOW_CUSTOM_PLANS":               `[{"name":"plan1"}]`,
				"DATASTORE_CUSTOM_PLANS":              `[{"name":"plan1"}]`,
				"DIALOGFLOW_CUSTOM_PLANS":             `[{"name":"plan1"}]`,
				"FIRESTORE_CUSTOM_PLANS":              `[{"name":"plan1"}]`,
				"ML_APIS_CUSTOM_PLANS":                `[{"name":"plan1"}]`,
				"PUBSUB_CUSTOM_PLANS":                 `[{"name":"plan1"}]`,
				"SPANNER_CUSTOM_PLANS":                `[{"name":"plan1"}]`,
				"STACKDRIVER_DEBUGGER_CUSTOM_PLANS":   `[{"name":"plan1"}]`,
				"STACKDRIVER_MONITORING_CUSTOM_PLANS": `[{"name":"plan1"}]`,
				"STACKDRIVER_PROFILER_CUSTOM_PLANS":   `[{"name":"plan1"}]`,
				"STACKDRIVER_TRACE_CUSTOM_PLANS":      `[{"name":"plan1"}]`,
				"STORAGE_CUSTOM_PLANS":                `[{"name":"plan1"}]`,
			},
		},
	}

	for tn, tc := range cases {
		t.Run(tn, tc.Run)
	}
}

func TestFormatCustomPlans(t *testing.T) {
	cases := map[string]MigrationTest{
		"no-overlap": {
			TileProperties: `
        {
          "properties": {
            ".properties.org": {"type": "string","value": "system"}
          }
        }`,
			Migration:   FormatCustomPlans(),
			ExpectedEnv: map[string]string{"ORG": "system"},
		},
		"service-gets-deleted": {
			TileProperties: `
        {
          "properties": {
            ".properties.cloudsql_mysql_custom_plans": {"value":"[{\"service\":\"4bc59b9a-8520-409f-85da-1c7552315863\"}]"}
          }
        }`,
			Migration: FormatCustomPlans(),
			ExpectedEnv: map[string]string{
				"CLOUDSQL_MYSQL_CUSTOM_PLANS": `[{"properties":{}}]`,
			},
		},
		"guid-gets-set": {
			TileProperties: `
        {
          "properties": {
            ".properties.cloudsql_mysql_custom_plans": {"value":"[{\"guid\":\"4bc59b9a-8520-409f-85da-1c7552315863\"}]"}
          }
        }`,
			Migration: FormatCustomPlans(),
			ExpectedEnv: map[string]string{
				"CLOUDSQL_MYSQL_CUSTOM_PLANS": `[{"guid":"4bc59b9a-8520-409f-85da-1c7552315863","properties":{}}]`,
			},
		},
		"description-gets-set": {
			TileProperties: `
        {
          "properties": {
            ".properties.cloudsql_mysql_custom_plans": {"value":"[{\"description\":\"some-value\"}]"}
          }
        }`,
			Migration: FormatCustomPlans(),
			ExpectedEnv: map[string]string{
				"CLOUDSQL_MYSQL_CUSTOM_PLANS": `[{"description":"some-value","properties":{}}]`,
			},
		},
		"name-gets-set": {
			TileProperties: `
        {
          "properties": {
            ".properties.cloudsql_mysql_custom_plans": {"value":"[{\"name\":\"some-value\"}]"}
          }
        }`,
			Migration: FormatCustomPlans(),
			ExpectedEnv: map[string]string{
				"CLOUDSQL_MYSQL_CUSTOM_PLANS": `[{"name":"some-value","properties":{}}]`,
			},
		},
		"display-name-gets-set": {
			TileProperties: `
        {
          "properties": {
            ".properties.cloudsql_mysql_custom_plans": {"value":"[{\"display_name\":\"some-value\"}]"}
          }
        }`,
			Migration: FormatCustomPlans(),
			ExpectedEnv: map[string]string{
				"CLOUDSQL_MYSQL_CUSTOM_PLANS": `[{"display_name":"some-value","properties":{}}]`,
			},
		},
		"multiple-plans": {
			TileProperties: `
        {
          "properties": {
            ".properties.cloudsql_mysql_custom_plans": {"value":"[{\"name\":\"plan1\"},{\"name\":\"plan2\"}]"}
          }
        }`,
			Migration: FormatCustomPlans(),
			ExpectedEnv: map[string]string{
				"CLOUDSQL_MYSQL_CUSTOM_PLANS": `[{"name":"plan1","properties":{}},{"name":"plan2","properties":{}}]`,
			},
		},
		"unknown-props-go-to-properties": {
			TileProperties: `
        {
          "properties": {
            ".properties.cloudsql_mysql_custom_plans": {"value":"[{\"prop1\":\"some-value1\",\"prop2\":\"some-value2\"}]"}
          }
        }`,
			Migration: FormatCustomPlans(),
			ExpectedEnv: map[string]string{
				"CLOUDSQL_MYSQL_CUSTOM_PLANS": `[{"properties":{"prop1":"some-value1","prop2":"some-value2"}}]`,
			},
		},
		// The v4 tile only supports MySQL, Postgres, BigTable, Spanner, and Storage custom plans
		"applies-to-v4-tile-values": {
			TileProperties: `
        {
          "properties": {
            ".properties.cloudsql_mysql_custom_plans": {"value":"[{\"name\":\"mysql\"}]"},
            ".properties.bigtable_custom_plans": {"value":"[{\"name\":\"bigtable\"}]"},
            ".properties.cloudsql_postgres_custom_plans": {"value":"[{\"name\":\"postgres\"}]"},
            ".properties.spanner_custom_plans": {"value":"[{\"name\":\"spanner\"}]"},
            ".properties.storage_custom_plans": {"value":"[{\"name\":\"storage\"}]"}
          }
        }`,
			Migration: FormatCustomPlans(),
			ExpectedEnv: map[string]string{
				"CLOUDSQL_MYSQL_CUSTOM_PLANS":    `[{"name":"mysql","properties":{}}]`,
				"BIGTABLE_CUSTOM_PLANS":          `[{"name":"bigtable","properties":{}}]`,
				"CLOUDSQL_POSTGRES_CUSTOM_PLANS": `[{"name":"postgres","properties":{}}]`,
				"SPANNER_CUSTOM_PLANS":           `[{"name":"spanner","properties":{}}]`,
				"STORAGE_CUSTOM_PLANS":           `[{"name":"storage","properties":{}}]`,
			},
		},
	}

	for tn, tc := range cases {
		t.Run(tn, tc.Run)
	}
}

func TestMergeToServiceConfig(t *testing.T) {
	cases := map[string]MigrationTest{
		"no-overlap": {
			TileProperties: `
        {
          "properties": {
            ".properties.org": {"type": "string","value": "system"}
          }
        }`,
			Migration:   MergeToServiceConfig(),
			ExpectedEnv: map[string]string{"ORG": "system", "GSB_SERVICE_CONFIG": defaultEmptyGsbServiceConfig},
		},
		"all-keys": {
			TileProperties: `
        {
          "properties": {
            ".properties.org": {"type": "string","value": "system"}
          }
        }`,
			Migration:   MergeToServiceConfig(),
			ExpectedEnv: map[string]string{"ORG": "system", "GSB_SERVICE_CONFIG": defaultEmptyGsbServiceConfig},
		},
	}

	for tn, tc := range cases {
		t.Run(tn, tc.Run)
	}
}
