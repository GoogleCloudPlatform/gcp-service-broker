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

// migrations returns a list of all migrations to be performed in order.
func migrations() []Migration {
	return []Migration{
		DeleteWhitelistKeys(),
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
