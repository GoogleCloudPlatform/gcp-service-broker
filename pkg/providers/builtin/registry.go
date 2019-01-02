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
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/api_service"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/bigquery"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/bigtable"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/cloudsql"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/dataflow"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/datastore"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/dialogflow"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/firestore"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/pubsub"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/spanner"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/stackdriver"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/storage"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/broker"
)

// BuiltinBrokerRegistry creates a new registry with all the built-in brokers
// added to it.
func BuiltinBrokerRegistry() broker.BrokerRegistry {
	out := broker.BrokerRegistry{}
	RegisterBuiltinBrokers(out)
	return out
}

// RegisterBuiltinBrokers adds the built-in brokers to the given registry.
func RegisterBuiltinBrokers(registry broker.BrokerRegistry) {
	registry.Register(api_service.ServiceDefinition())
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
