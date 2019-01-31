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
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/GoogleCloudPlatform/gcp-service-broker/utils"
)

// Migration holds the information necessary to modify the values in the tile.
type Migration struct {
	Name       string
	TileScript string
	GoFunc     func(env map[string]string)
}

func deleteMigration(name string, env []string) Migration {
	js := ``
	for _, v := range env {
		js += "delete properties.properties['.properties." + strings.ToLower(v) + "'];\n"
	}

	return Migration{
		Name:       name,
		TileScript: js,
		GoFunc: func(envMap map[string]string) {
			for _, v := range env {
				delete(envMap, v)
			}
		},
	}
}

// NoOp is an empty migration.
func NoOp() Migration {
	return Migration{
		Name:       "Noop",
		TileScript: ``,
		GoFunc:     func(envMap map[string]string) {},
	}
}

// DeleteWhitelistKeys removes the whitelist keys used prior to 4.0
func DeleteWhitelistKeys() Migration {
	whitelistKeys := []string{
		"GSB_SERVICE_GOOGLE_BIGQUERY_WHITELIST",
		"GSB_SERVICE_GOOGLE_BIGTABLE_WHITELIST",
		"GSB_SERVICE_GOOGLE_CLOUDSQL_MYSQL_WHITELIST",
		"GSB_SERVICE_GOOGLE_CLOUDSQL_POSTGRES_WHITELIST",
		"GSB_SERVICE_GOOGLE_ML_APIS_WHITELIST",
		"GSB_SERVICE_GOOGLE_PUBSUB_WHITELIST",
		"GSB_SERVICE_GOOGLE_SPANNER_WHITELIST",
		"GSB_SERVICE_GOOGLE_STORAGE_WHITELIST",
	}

	return deleteMigration("Delete whitelist keys from 4.x", whitelistKeys)
}

// DeleteEnabledKeys removes the enable/disable flag keys.
func DeleteEnabledKeys() Migration {
	keys := []string{
		"GSB_SERVICE_GOOGLE_BIGQUERY_ENABLED",
		"GSB_SERVICE_GOOGLE_BIGTABLE_ENABLED",
		"GSB_SERVICE_GOOGLE_CLOUDSQL_MYSQL_ENABLED",
		"GSB_SERVICE_GOOGLE_CLOUDSQL_POSTGRES_ENABLED",
		"GSB_SERVICE_GOOGLE_ML_APIS_ENABLED",
		"GSB_SERVICE_GOOGLE_PUBSUB_ENABLED",
		"GSB_SERVICE_GOOGLE_SPANNER_ENABLED",
		"GSB_SERVICE_GOOGLE_STORAGE_ENABLED",
	}

	return deleteMigration("Delete enabled keys from 4.x", keys)
}

// DeleteProvisionDefaultKeys removes the provision default keys.
func DeleteProvisionDefaultKeys() Migration {
	keys := []string{
		"GSB_SERVICE_GOOGLE_BIGQUERY_PROVISION_DEFAULTS",
		"GSB_SERVICE_GOOGLE_BIGTABLE_PROVISION_DEFAULTS",
		"GSB_SERVICE_GOOGLE_CLOUDSQL_MYSQL_PROVISION_DEFAULTS",
		"GSB_SERVICE_GOOGLE_CLOUDSQL_POSTGRES_PROVISION_DEFAULTS",
		"GSB_SERVICE_GOOGLE_ML_APIS_PROVISION_DEFAULTS",
		"GSB_SERVICE_GOOGLE_PUBSUB_PROVISION_DEFAULTS",
		"GSB_SERVICE_GOOGLE_SPANNER_PROVISION_DEFAULTS",
		"GSB_SERVICE_GOOGLE_STORAGE_PROVISION_DEFAULTS",
	}

	return deleteMigration("Delete provision defaults from 4.x", keys)
}

// DeleteBindDefaultKeys removes the bind default keys.
func DeleteBindDefaultKeys() Migration {
	keys := []string{
		"GSB_SERVICE_GOOGLE_BIGQUERY_BIND_DEFAULTS",
		"GSB_SERVICE_GOOGLE_BIGTABLE_BIND_DEFAULTS",
		"GSB_SERVICE_GOOGLE_CLOUDSQL_MYSQL_BIND_DEFAULTS",
		"GSB_SERVICE_GOOGLE_CLOUDSQL_POSTGRES_BIND_DEFAULTS",
		"GSB_SERVICE_GOOGLE_ML_APIS_BIND_DEFAULTS",
		"GSB_SERVICE_GOOGLE_PUBSUB_BIND_DEFAULTS",
		"GSB_SERVICE_GOOGLE_SPANNER_BIND_DEFAULTS",
		"GSB_SERVICE_GOOGLE_STORAGE_BIND_DEFAULTS",
	}

	return deleteMigration("Delete bind defaults from 4.x", keys)
}

func customPlanKeys() []string {
	return []string{
		"BIGQUERY_CUSTOM_PLANS",
		"BIGTABLE_CUSTOM_PLANS",
		"CLOUDSQL_MYSQL_CUSTOM_PLANS",
		"CLOUDSQL_POSTGRES_CUSTOM_PLANS",
		"ML_APIS_CUSTOM_PLANS",
		"PUBSUB_CUSTOM_PLANS",
		"SPANNER_CUSTOM_PLANS",
		"STORAGE_CUSTOM_PLANS",
	}
}

// DeleteCustomPlanKeys removes the custom plan keys.
func DeleteCustomPlanKeys() Migration {
	return deleteMigration("Delete custom plans from 4.x", customPlanKeys())
}

func CollapseCustomPlans() Migration {
	jst := JsTransform{
		EnvironmentVariables: customPlanKeys(),
		MigrationJs: `
		function replaceCustomPlans(planList) {
			// if this is already JSON rather than a tile custom format
			// just return it
			if (typeof(planList) === "string") {
				return planList;
			}

		  var out = [];
		  for (var i in planList) {
		    var plan = planList[i];
		    var converted = {}
		    for (var k in plan) {
		      converted[k] = plan[k].value
		    }

		    out.push(converted);
		  }

			return JSON.stringify(out);
		}`,
		TransformFuncName: "replaceCustomPlans",
	}

	return Migration{
		Name:       "Collapse custom plans",
		TileScript: jst.ToJs(),
		GoFunc: func(env map[string]string) {
			if err := jst.RunGo(env); err != nil {
				log.Fatal(err)
			}
		},
	}
}

func FormatCustomPlans() Migration {
	jst := JsTransform{
		EnvironmentVariables: customPlanKeys(),
		MigrationJs: `
		function formatPlans(plansListJson) {
			var planList = JSON.parse(plansListJson);

			var out = [];
			for (var i in planList) {
				var plan = planList[i];
				var converted = {};
				var properties = {};
				for (var k in plan) {
					var value = plan[k];
					if (k === 'service') { continue; } // skip service key, it's no longer needed

					if(["guid", "description", "name", "display_name"].indexOf(k) >= 0) {
						converted[k] = value
					} else {
						properties[k] = value;
					}
				}
				converted.properties = properties;

				out.push(converted);
			}

			return JSON.stringify(out);
		}`,
		TransformFuncName: "formatPlans",
	}

	return Migration{
		Name:       "Format custom plans",
		TileScript: jst.ToJs(),
		GoFunc: func(env map[string]string) {
			if err := jst.RunGo(env); err != nil {
				log.Fatal(err)
			}
		},
	}
}

func MergeServices() Migration {
	jst := JsTransform{
		EnvironmentVariables: customPlanKeys(),
		MigrationJs: `
		function mergeServices(v, context) {
			var out = {};
			for (var key in context) {
				var value = context[key];
				if (value != null) {
					out[key] = JSON.parse(context[key]);
				}
			}

			return JSON.stringify(out, null, "  ");
		}`,
		TransformFuncName: "mergeServices",
		Context: func(envVar string) map[string]string {
			svcName := strings.TrimSuffix(envVar, "_CUSTOM_PLANS")

			return map[string]string{
				"custom_plans":       envVar,
				"enabled":            fmt.Sprintf("GSB_SERVICE_GOOGLE_%s_ENABLED", svcName),
				"bind_defaults":      fmt.Sprintf("GSB_SERVICE_GOOGLE_%s_BIND_DEFAULTS", svcName),
				"provision_defaults": fmt.Sprintf("GSB_SERVICE_GOOGLE_%s_PROVISION_DEFAULTS", svcName),
			}
		},
	}

	return Migration{
		Name:       "Merge Services",
		TileScript: jst.ToJs(),
		GoFunc: func(env map[string]string) {
			if err := jst.RunGo(env); err != nil {
				log.Fatal(err)
			}
		},
	}
}

func MergeToServiceConfig() Migration {
	jst := JsTransform{
		EnvironmentVariables: []string{"GSB_SERVICE_CONFIG"},
		CreateIfNotExists:    true,
		MigrationJs: `
		function mergeServices(v, context) {
			var out = {};
			for (var key in context) {
				var value = context[key];
				if (value != null) {
					out[key] = JSON.parse(context[key]);
				}
			}

			return JSON.stringify(out, null, "  ");
		}`,
		TransformFuncName: "mergeServices",
		Context: func(envVar string) map[string]string {
			return map[string]string{
				"f80c0a3e-bd4d-4809-a900-b4e33a6450f1": "BIGQUERY_CUSTOM_PLANS",
				"b8e19880-ac58-42ef-b033-f7cd9c94d1fe": "BIGTABLE_CUSTOM_PLANS",
				"4bc59b9a-8520-409f-85da-1c7552315863": "CLOUDSQL_MYSQL_CUSTOM_PLANS",
				"cbad6d78-a73c-432d-b8ff-b219a17a803a": "CLOUDSQL_POSTGRES_CUSTOM_PLANS",
				"3e897eb3-9062-4966-bd4f-85bda0f73b3d": "DATAFLOW_CUSTOM_PLANS",
				"76d4abb2-fee7-4c8f-aee1-bcea2837f02b": "DATASTORE_CUSTOM_PLANS",
				"e84b69db-3de9-4688-8f5c-26b9d5b1f129": "DIALOGFLOW_CUSTOM_PLANS",
				"a2b7b873-1e34-4530-8a42-902ff7d66b43": "FIRESTORE_CUSTOM_PLANS",
				"5ad2dce0-51f7-4ede-8b46-293d6df1e8d4": "ML_APIS_CUSTOM_PLANS",
				"628629e3-79f5-4255-b981-d14c6c7856be": "PUBSUB_CUSTOM_PLANS",
				"51b3e27e-d323-49ce-8c5f-1211e6409e82": "SPANNER_CUSTOM_PLANS",
				"83837945-1547-41e0-b661-ea31d76eed11": "STACKDRIVER_DEBUGGER_CUSTOM_PLANS",
				"2bc0d9ed-3f68-4056-b842-4a85cfbc727f": "STACKDRIVER_MONITORING_CUSTOM_PLANS",
				"00b9ca4a-7cd6-406a-a5b7-2f43f41ade75": "STACKDRIVER_PROFILER_CUSTOM_PLANS",
				"c5ddfe15-24d9-47f8-8ffe-f6b7daa9cf4a": "STACKDRIVER_TRACE_CUSTOM_PLANS",
				"b9e4332e-b42b-4680-bda5-ea1506797474": "STORAGE_CUSTOM_PLANS",
			}
		},
	}

	return Migration{
		Name:       "Merge Service Config",
		TileScript: jst.ToJs(),
		GoFunc: func(env map[string]string) {
			if _, ok := env["GSB_SERVICE_CONFIG"]; !ok {
				env["GSB_SERVICE_CONFIG"] = ""
			}

			if err := jst.RunGo(env); err != nil {
				log.Fatal(err)
			}
		},
	}
}

// migrations returns a list of all migrations to be performed in order.
func migrations() []Migration {
	return []Migration{
		DeleteWhitelistKeys(),
		CollapseCustomPlans(),
		FormatCustomPlans(),
		MergeServices(),
		MergeToServiceConfig(),
		DeleteProvisionDefaultKeys(),
		DeleteBindDefaultKeys(),
		DeleteCustomPlanKeys(),
		DeleteEnabledKeys(),
	}
}

// FullMigration holds the complete migration path.
func FullMigration() Migration {
	migrations := migrations()

	var buf bytes.Buffer
	for i, migration := range migrations {
		fmt.Fprintln(&buf, "{")
		fmt.Fprintf(&buf, "// migration: %d, %s\n", i, migration.Name)
		fmt.Fprintln(&buf, migration.TileScript)
		fmt.Fprintln(&buf, "}")
	}

	return Migration{
		Name:       "Full Migration",
		TileScript: string(buf.Bytes()),
		GoFunc: func(env map[string]string) {
			for _, migration := range migrations {
				migration.GoFunc(env)
			}
		},
	}
}

// MigrateEnv migrates the currently set environment variables and returns
// the diff.
func MigrateEnv() map[string]Diff {
	env := splitEnv()
	orig := utils.CopyStringMap(env)

	FullMigration().GoFunc(env)

	return DiffStringMap(orig, env)
}

// splitEnv splits os.Environ style environment variables that
// are in an array of "key=value" strings.
func splitEnv() map[string]string {
	envMap := make(map[string]string)

	for _, e := range os.Environ() {
		split := strings.Split(e, "=")
		envMap[split[0]] = strings.Join(split[1:], "=")
	}

	return envMap
}
