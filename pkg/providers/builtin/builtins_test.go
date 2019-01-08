// Copyright 2010 the Service Broker Project Authors.
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
	"context"
	"testing"

	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/api_service"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/bigquery"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/bigtable"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/broker_base"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/broker_base/broker_basefakes"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/cloudsql"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/dataflow"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/datastore"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/dialogflow"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/firestore"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/models"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/pubsub"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/spanner"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/stackdriver"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/storage"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/broker"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/varcontext"
)

func TestServiceProviderAsync(t *testing.T) {
	cases := map[string]struct {
		AsyncProvisionExpected   bool
		AsyncDeprovisionExpected bool
		Provider                 broker.ServiceProvider
	}{
		"ml": {
			Provider: &api_service.ApiServiceBroker{},
		},
		"bigquery": {
			Provider: &bigquery.BigQueryBroker{},
		},
		"bigtable": {
			Provider: &bigtable.BigTableBroker{},
		},
		"cloudsql": {
			AsyncProvisionExpected:   true,
			AsyncDeprovisionExpected: true,
			Provider:                 &cloudsql.CloudSQLBroker{},
		},
		"dataflow": {
			Provider: &dataflow.DataflowBroker{},
		},
		"datastore": {
			Provider: &datastore.DatastoreBroker{},
		},
		"dialogflow": {
			Provider: &dialogflow.DialogflowBroker{},
		},
		"firestore": {
			Provider: &firestore.FirestoreBroker{},
		},
		"pubsub": {
			Provider: &pubsub.PubSubBroker{},
		},
		"spanner": {
			AsyncProvisionExpected: true,
			Provider:               &spanner.SpannerBroker{},
		},
		"storage": {
			Provider: &storage.StorageBroker{},
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			actualProvisionAsync := tc.Provider.ProvisionsAsync()
			if tc.AsyncProvisionExpected != actualProvisionAsync {
				t.Errorf("Expected async provision to match. Expected: %t, Actual: %t", tc.AsyncProvisionExpected, actualProvisionAsync)
			}

			actualDeprovisionAsync := tc.Provider.DeprovisionsAsync()
			if tc.AsyncDeprovisionExpected != actualDeprovisionAsync {
				t.Errorf("Expected async deprovision to match. Expected: %t, Actual: %t", tc.AsyncDeprovisionExpected, actualDeprovisionAsync)
			}
		})
	}
}

func TestThinWrapperServiceProviders(t *testing.T) {
	cases := map[string]func(broker_base.BrokerBase) broker.ServiceProvider{
		"pubsub": func(brokerBase broker_base.BrokerBase) broker.ServiceProvider {
			return &pubsub.PubSubBroker{BrokerBase: brokerBase}
		},
		"stackdriver": func(brokerBase broker_base.BrokerBase) broker.ServiceProvider {
			return &stackdriver.StackdriverAccountProvider{BrokerBase: brokerBase}
		},
		"ml": func(brokerBase broker_base.BrokerBase) broker.ServiceProvider {
			return &api_service.ApiServiceBroker{BrokerBase: brokerBase}
		},
		"bigquery": func(brokerBase broker_base.BrokerBase) broker.ServiceProvider {
			return &bigquery.BigQueryBroker{BrokerBase: brokerBase}
		},
		"dataflow": func(brokerBase broker_base.BrokerBase) broker.ServiceProvider {
			return &dataflow.DataflowBroker{BrokerBase: brokerBase}
		},
		"datastore": func(brokerBase broker_base.BrokerBase) broker.ServiceProvider {
			return &datastore.DatastoreBroker{BrokerBase: brokerBase}
		},
		"dialogflow": func(brokerBase broker_base.BrokerBase) broker.ServiceProvider {
			return &dialogflow.DialogflowBroker{BrokerBase: brokerBase}
		},
		"firestore": func(brokerBase broker_base.BrokerBase) broker.ServiceProvider {
			return &firestore.FirestoreBroker{BrokerBase: brokerBase}
		},
		"spanner": func(brokerBase broker_base.BrokerBase) broker.ServiceProvider {
			return &spanner.SpannerBroker{BrokerBase: brokerBase}
		},
		"storage": func(brokerBase broker_base.BrokerBase) broker.ServiceProvider {
			return &storage.StorageBroker{BrokerBase: brokerBase}
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			accountManager := broker_basefakes.FakeServiceAccountManager{}
			brokerBase := broker_base.BrokerBase{
				AccountManager: &accountManager,
			}
			serviceProvider := tc(brokerBase)

			if _, err := serviceProvider.Bind(context.Background(), &varcontext.VarContext{}); err != nil {
				t.Fatal(err)
			}
			if accountManager.CreateCredentialsCallCount() != 1 {
				t.Errorf("Expected CreateCredentials to be called once. Expected: %d, Actual: %d", 1, accountManager.CreateCredentialsCallCount())
			}

			if err := serviceProvider.Unbind(context.Background(), models.ServiceInstanceDetails{}, models.ServiceBindingCredentials{}); err != nil {
				t.Fatal(err)
			}
			if accountManager.DeleteCredentialsCallCount() != 1 {
				t.Errorf("Expected DeleteCredentials to be called once. Expected: %d, Actual: %d", 1, accountManager.DeleteCredentialsCallCount())
			}
		})
	}
}
