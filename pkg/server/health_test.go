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
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"

	// Needed to open the sqlite3 database
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func TestNewHealthHandler(t *testing.T) {
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		t.Fatalf("couldn't create database: %v", err)
	}
	defer os.Remove("test.db")

	cases := map[string]struct {
		Endpoint       string
		ExpectedStatus int
		ExpectedBody   string
		LiveErr        error
		ReadyErr       error
	}{
		"live endpoint full": {
			Endpoint:       "/live?full=1",
			ExpectedStatus: 200,
			ExpectedBody:   `{"test-live":"OK"}`,
		},
		"live endpoint minimal": {
			Endpoint:       "/live",
			ExpectedStatus: 200,
			ExpectedBody:   `{}`,
		},
		"live endpoint bad": {
			Endpoint:       "/live?full=1",
			ExpectedStatus: 503,
			ExpectedBody:   `{"test-live":"bad-value"}`,
			LiveErr:        errors.New("bad-value"),
		},
		"ready endpoint minimal": {
			Endpoint:       "/ready",
			ExpectedStatus: 200,
			ExpectedBody:   `{}`,
		},
		"ready endpoint full": {
			Endpoint:       "/ready?full=1",
			ExpectedStatus: 200,
			ExpectedBody:   `{"database":"OK","test-live":"OK","test-ready":"OK"}`,
		},
		"ready endpoint bad liveness": {
			Endpoint:       "/ready?full=1",
			ExpectedStatus: 503,
			ExpectedBody:   `{"database":"OK","test-live":"bad-value","test-ready":"OK"}`,
			LiveErr:        errors.New("bad-value"),
		},
		"ready endpoint bad": {
			Endpoint:       "/ready?full=1",
			ExpectedStatus: 503,
			ExpectedBody:   `{"database":"OK","test-live":"OK","test-ready":"bad-value"}`,
			ReadyErr:       errors.New("bad-value"),
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			router := mux.NewRouter()
			handler := AddHealthHandler(router, db.DB())
			handler.AddLivenessCheck("test-live", func() error {
				return tc.LiveErr
			})
			handler.AddReadinessCheck("test-ready", func() error {
				return tc.ReadyErr
			})

			request := httptest.NewRequest(http.MethodGet, tc.Endpoint, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, request)

			if w.Code != tc.ExpectedStatus {
				t.Fatalf("Expected response code: %v got: %v", tc.ExpectedStatus, w.Code)
			}

			compacted := &bytes.Buffer{}
			if err := json.Compact(compacted, w.Body.Bytes()); err != nil {
				t.Fatal(err)
			}

			if compacted.String() != tc.ExpectedBody {
				t.Fatalf("Expected response: %v got: %v", tc.ExpectedBody, compacted.String())
			}
		})
	}
}
