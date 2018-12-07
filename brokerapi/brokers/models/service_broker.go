// Copyright 2018 the Service Broker Project Authors.
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

package models

import "github.com/GoogleCloudPlatform/gcp-service-broker/utils"

// This custom user agent string is added to provision calls so that Google can track the aggregated use of this tool
// We can better advocate for devoting resources to supporting cloud foundry and this service broker if we can show
// good usage statistics for it, so if you feel the need to fork this repo, please leave this string in place!
var CustomUserAgent = "cf-gcp-service-broker-test " + utils.Version

func ProductionizeUserAgent() {
	CustomUserAgent = "cf-gcp-service-broker " + utils.Version
}

const StorageName = "google-storage"
const BigqueryName = "google-bigquery"
const BigtableName = "google-bigtable"
const CloudsqlMySQLName = "google-cloudsql-mysql"
const CloudsqlPostgresName = "google-cloudsql-postgres"
const PubsubName = "google-pubsub"
const MlName = "google-ml-apis"
const SpannerName = "google-spanner"
const StackdriverTraceName = "google-stackdriver-trace"
const StackdriverDebuggerName = "google-stackdriver-debugger"
const StackdriverProfilerName = "google-stackdriver-profiler"
const DatastoreName = "google-datastore"
