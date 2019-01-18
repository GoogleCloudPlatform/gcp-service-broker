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
)

// Migration holds the information necessary to modify the values in the tile.
type Migration struct {
	Name       string
	TileScript string
}

// NoOp is an empty migration.
func NoOp() Migration {
	return Migration{
		Name:       "Noop",
		TileScript: ``,
	}
}

// DeleteWhitelistKeys removes the whitelist keys used prior to 4.0
func DeleteWhitelistKeys() Migration {
	return Migration{
		Name: "Delete unused keys from 4.x",
		TileScript: `
    var TO_REMOVE = [
      ".properties.gsb_service_google_bigquery_whitelist",
      ".properties.gsb_service_google_bigtable_whitelist",
      ".properties.gsb_service_google_cloudsql_mysql_whitelist",
      ".properties.gsb_service_google_cloudsql_postgres_whitelist",
      ".properties.gsb_service_google_ml_apis_whitelist",
      ".properties.gsb_service_google_pubsub_whitelist",
      ".properties.gsb_service_google_spanner_whitelist",
      ".properties.gsb_service_google_storage_whitelist",
    ];

    for (var i = 0; i < TO_REMOVE.length; i++) {
      delete properties.properties[TO_REMOVE[i]];
    }`,
	}
}

// FullMigration holds the complete migration path.
func FullMigration() Migration {
	migrations := []Migration{
		DeleteWhitelistKeys(),
	}

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
	}
}
