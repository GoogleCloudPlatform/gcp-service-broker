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

func serviceNames() []string {
	return []string{
		"BIGQUERY",
		"BIGTABLE",
		"CLOUDSQL_MYSQL",
		"CLOUDSQL_POSTGRES",
		"DATAFLOW",
		"DATASTORE",
		"DIALOGFLOW",
		"FIRESTORE",
		"ML_APIS",
		"PUBSUB",
		"SPANNER",
		"STACKDRIVER_DEBUGGER",
		"STACKDRIVER_MONITORING",
		"STACKDRIVER_PROFILER",
		"STACKDRIVER_TRACE",
		"STORAGE",
	}
}

func customPlanKeys() (out []string) {
	for _, s := range serviceNames() {
		out = append(out, s+"_CUSTOM_PLANS")
	}

	return
}

// MergeToServiceConfig merges the custom plans, enabled, bind defaults, and
// provision defaults for every 4.x service into a single configuration option.
func MergeToServiceConfig() Migration {
	jst := JsTransform{
		EnvironmentVariables: []string{"GSB_SERVICE_CONFIG"},
		MigrationJs: `

		function anyExist(envVars) {
			for(var i in envVars) {
				if(lookupProp(envVars[i]) !== null) {
					return true;
				}
			}

			return false;
		}

    function mergeService(svcName) {
      var customPlanVar = svcName+"_CUSTOM_PLANS";
      var enabledVar = "GSB_SERVICE_GOOGLE_"+svcName+"_ENABLED";
      var bindDefaultsVar = "GSB_SERVICE_GOOGLE_"+svcName+"_BIND_DEFAULTS";
      var provisionDefaultsVar = "GSB_SERVICE_GOOGLE_"+svcName+"_PROVISION_DEFAULTS";

			if (! anyExist([customPlanVar, enabledVar, bindDefaultsVar, provisionDefaultsVar])) {
				return null;
			}

			// explicitly check for boolean enabled false and turn it into a string.
			var enabled = lookupProp(enabledVar);
			if (enabled === false) {
				enabled = 'false';
			}

      var context = {
        "custom_plans": lookupProp(customPlanVar) || '[]',
        "enabled": enabled || 'true',
        "bind_defaults": lookupProp(bindDefaultsVar) || '{}',
        "provision_defaults": lookupProp(provisionDefaultsVar) || '{}',
      }

      deleteProp(customPlanVar);
      deleteProp(enabledVar);
      deleteProp(bindDefaultsVar);
      deleteProp(provisionDefaultsVar);

      // convert to plain object
      for (var key in context) {
        context[key] = JSON.parse(context[key]);
      }

			context["disabled"] = !context["enabled"];
			delete context["enabled"];

      context["//"] = "Builtin " + svcName;
      return context;
    }

    function mergeServices(envVar) {
      var out = {
        "f80c0a3e-bd4d-4809-a900-b4e33a6450f1": mergeService("BIGQUERY"),
        "b8e19880-ac58-42ef-b033-f7cd9c94d1fe": mergeService("BIGTABLE"),
        "4bc59b9a-8520-409f-85da-1c7552315863": mergeService("CLOUDSQL_MYSQL"),
        "cbad6d78-a73c-432d-b8ff-b219a17a803a": mergeService("CLOUDSQL_POSTGRES"),
        "3e897eb3-9062-4966-bd4f-85bda0f73b3d": mergeService("DATAFLOW"),
        "76d4abb2-fee7-4c8f-aee1-bcea2837f02b": mergeService("DATASTORE"),
        "e84b69db-3de9-4688-8f5c-26b9d5b1f129": mergeService("DIALOGFLOW"),
        "a2b7b873-1e34-4530-8a42-902ff7d66b43": mergeService("FIRESTORE"),
        "5ad2dce0-51f7-4ede-8b46-293d6df1e8d4": mergeService("ML_APIS"),
        "628629e3-79f5-4255-b981-d14c6c7856be": mergeService("PUBSUB"),
        "51b3e27e-d323-49ce-8c5f-1211e6409e82": mergeService("SPANNER"),
        "83837945-1547-41e0-b661-ea31d76eed11": mergeService("STACKDRIVER_DEBUGGER"),
        "2bc0d9ed-3f68-4056-b842-4a85cfbc727f": mergeService("STACKDRIVER_MONITORING"),
        "00b9ca4a-7cd6-406a-a5b7-2f43f41ade75": mergeService("STACKDRIVER_PROFILER"),
        "c5ddfe15-24d9-47f8-8ffe-f6b7daa9cf4a": mergeService("STACKDRIVER_TRACE"),
        "b9e4332e-b42b-4680-bda5-ea1506797474": mergeService("STORAGE")
      };

			Object.keys(out).forEach(function(key) {
			if (out[key] == null) delete out[key];
		});


			for (var key in out) {
				if(out[key] == null) {
					delete out[key];
				}
			}

      setProp(envVar, JSON.stringify(out, null, "  "));
    }`,
		TransformFuncName: "mergeServices",
	}

	return Migration{
		Name:       "Merge Service Config",
		TileScript: jst.ToJs(),
		GoFunc:     jst.RunGo,
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

// CollapseCustomPlans converts the Tile object based custom plan structure to
// a JSON flattened structure suitable for passing in an environment variable.
func CollapseCustomPlans() Migration {
	jst := JsTransform{
		EnvironmentVariables: customPlanKeys(),
		MigrationJs: `
    function replaceCustomPlans(envVar) {
      var planList = lookupProp(envVar);

      // do nothing if the planList doesn't exist or if it's already a string
      if (planList == null || typeof (planList) === 'string') {
        return;
      }

      var out = [];
      for (var i in planList) {
        var plan = planList[i];
        var converted = {};
        for (var k in plan) {
          converted[k] = plan[k].value;
        }

        out.push(converted);
      }

      setProp(envVar, JSON.stringify(out, null, "  "));
    }

`,
		TransformFuncName: "replaceCustomPlans",
	}

	return Migration{
		Name:       "Collapse custom plans",
		TileScript: jst.ToJs(),
		GoFunc:     NoOp().GoFunc,
	}
}

// FormatCustomPlans converts a flat variable map that the tile used to produce
// into a structure the CustomPlan can parse where plan variables get put
// under a properties map.
func FormatCustomPlans() Migration {
	jst := JsTransform{
		EnvironmentVariables: customPlanKeys(),
		MigrationJs: `
    function formatPlans(envVar) {
      var plansListJson = lookupProp(envVar);
      if (plansListJson == null) {
        return;
      }

      var planList = JSON.parse(plansListJson);

      var out = [];
      for (var i in planList) {
        var plan = planList[i];
        var converted = {};
        var properties = {};
        for (var k in plan) {
          var value = plan[k];
          if (k === 'service') { continue; }  // skip service key, it's no longer needed

          if (['guid', 'description', 'name', 'display_name'].indexOf(k) >= 0) {
            converted[k] = value;
          } else {
            properties[k] = value;
          }
        }
        converted.properties = properties;
        out.push(converted);
      }

      setProp(envVar, JSON.stringify(out));
    }
`,
		TransformFuncName: "formatPlans",
	}

	return Migration{
		Name:       "Format custom plans",
		TileScript: jst.ToJs(),
		GoFunc:     jst.RunGo,
	}
}
