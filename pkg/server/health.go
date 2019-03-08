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

package server

import (
	"database/sql"
	"time"

	"github.com/gorilla/mux"
	"github.com/heptiolabs/healthcheck"
)

// AddHealthHandler creates a new handler for health and liveness checks and
// adds it to the /live and /ready endpoints.
func AddHealthHandler(router *mux.Router, db *sql.DB) healthcheck.Handler {
	health := healthcheck.NewHandler()

	health.AddReadinessCheck("database", healthcheck.DatabasePingCheck(db, 2*time.Second))

	router.HandleFunc("/live", health.LiveEndpoint)
	router.HandleFunc("/ready", health.ReadyEndpoint)

	return health
}
