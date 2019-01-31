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

package builtin

import (
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/broker"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/providers/builtin/bigquery"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/providers/builtin/bigtable"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/providers/builtin/cloudsql"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/providers/builtin/dataflow"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/providers/builtin/datastore"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/providers/builtin/dialogflow"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/providers/builtin/firestore"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/providers/builtin/ml"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/providers/builtin/pubsub"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/providers/builtin/spanner"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/providers/builtin/stackdriver"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/providers/builtin/storage"
)

// NOTE(josephlewis42) unless there are extenuating circumstances, as of 2019
// no new builtin providers should be added. Instead, providers should be
// added using downloadable brokerpaks.

// BuiltinBrokerRegistry creates a new registry with all the built-in brokers
// added to it.
func BuiltinBrokerRegistry(cfg broker.ServiceConfigMap) *broker.ServiceRegistry {
	out := broker.NewServiceRegistry(cfg)
	RegisterBuiltinBrokers(out)
	return out
}

// RegisterBuiltinBrokers adds the built-in brokers to the given registry.
func RegisterBuiltinBrokers(registry *broker.ServiceRegistry) {
	registry.Register(ml.ServiceDefinition())
	registry.Register(bigquery.ServiceDefinition())
	registry.Register(bigtable.ServiceDefinition())
	registry.Register(cloudsql.MysqlServiceDefinition())
	registry.Register(cloudsql.PostgresServiceDefinition())
	registry.Register(dataflow.ServiceDefinition())
	registry.Register(datastore.ServiceDefinition())
	registry.Register(dialogflow.ServiceDefinition())
	registry.Register(firestore.ServiceDefinition())
	registry.Register(pubsub.ServiceDefinition())
	registry.Register(spanner.ServiceDefinition())
	registry.Register(stackdriver.StackdriverDebuggerServiceDefinition())
	registry.Register(stackdriver.StackdriverMonitoringServiceDefinition())
	registry.Register(stackdriver.StackdriverProfilerServiceDefinition())
	registry.Register(stackdriver.StackdriverTraceServiceDefinition())
	registry.Register(storage.ServiceDefinition())
}
