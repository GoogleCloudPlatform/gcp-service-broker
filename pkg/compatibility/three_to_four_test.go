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

package compatibility

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/models"
	"github.com/GoogleCloudPlatform/gcp-service-broker/db_service"
	"github.com/jinzhu/gorm"
	"github.com/pivotal-cf/brokerapi"

	// Import so the brokers register with the registry
	_ "github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers"
)

// These plan details are dumped from a 3.6 service broker to be realistic, the
// list also includes user-defined plans that were saved to the database.
const planDetails = `INSERT INTO plan_details (id, created_at, updated_at, deleted_at, service_id, name, features) VALUES
  ('000ad064-f354-4f2d-a615-672fc2b98551','2018-08-22 19:16:49','2018-08-22 19:16:49',NULL,'b8e19880-ac58-42ef-b033-f7cd9c94d1fe','bigtable-ssd-20','{\"num_nodes\":\"20\",\"storage_type\":\"SDD\"}'),
  ('11308f80-dc53-44f7-a49b-c3e65258b421','2018-08-22 19:16:49','2018-08-22 19:16:49',NULL,'83837945-1547-41e0-b661-ea31d76eed11','default','\"\"'),
  ('1228b43c-15ca-43e6-ae2d-70a33e47efa1','2018-08-22 19:16:49','2018-08-22 19:16:49',NULL,'b9e4332e-b42b-4680-bda5-ea1506797474','nearline','{\"storage_class\":\"NEARLINE\"}'),
  ('1588ad53-4b42-4fd4-af8e-e2935a9a1e61','2018-08-22 19:16:49','2018-08-22 19:16:49',NULL,'cbad6d78-a73c-432d-b8ff-b219a17a803a','postgres-n1-standard-2','{\"max_disk_size\":\"10000\",\"pricing_plan\":\"PER_USE\",\"tier\":\"db-n1-standard-2\"}'),('22b6c084-7977-4e85-a4bd-16c7df79d4df','2018-08-22 19:16:49','2018-08-22 19:16:49',NULL,'cbad6d78-a73c-432d-b8ff-b219a17a803a','postgres-n1-highmem-16','{\"max_disk_size\":\"10000\",\"pricing_plan\":\"PER_USE\",\"tier\":\"db-n1-highmem-16\"}'),
  ('289ed477-56db-4545-acc9-5e92b2842d11','2018-08-22 19:16:49','2018-08-22 19:16:49',NULL,'4bc59b9a-8520-409f-85da-1c7552315863','mysql-micro-dev','{\"max_disk_size\":\"100\",\"pricing_plan\":\"PER_USE\",\"tier\":\"db-f1-micro\"}'),('4d7e90b2-80d3-4df7-a9a1-3d400362e974','2018-08-22 19:16:49','2018-08-22 19:16:49',NULL,'4bc59b9a-8520-409f-85da-1c7552315863','mysql-n1-standard-2','{\"max_disk_size\":\"10000\",\"pricing_plan\":\"PACKAGE\",\"tier\":\"db-n1-standard-2\"}'),
  ('5554757a-4c1a-4435-a455-441bfc9617c4','2018-08-22 19:16:49','2018-08-22 19:16:49',NULL,'b8e19880-ac58-42ef-b033-f7cd9c94d1fe','bigtable-ssd-10','{\"num_nodes\":\"10\",\"storage_type\":\"SDD\"}'),
  ('5dd56dd9-cc53-4b71-aa4e-45d9c1266dbb','2018-08-22 19:16:49','2018-08-22 19:16:49',NULL,'5ad2dce0-51f7-4ede-8b46-293d6df1e8d4','default','\"\"'),
  ('5e1829b7-4f12-4a52-a7a0-39321a6a82e6','2018-08-22 19:16:49','2018-08-22 19:16:49',NULL,'51b3e27e-d323-49ce-8c5f-1211e6409e82','spanner-regional-10','{\"num_nodes\":\"10\"}'),
  ('6c36cc15-cb2e-494d-a0e7-f4a34e0ff501','2018-08-22 19:16:49','2018-08-22 19:16:49',NULL,'b9e4332e-b42b-4680-bda5-ea1506797474','reduced_availability','{\"storage_class\":\"DURABLE_REDUCED_AVAILABILITY\"}'),
  ('728b61e5-2b7b-4a2e-a2a1-88797c6db1c2','2018-08-22 19:16:49','2018-08-22 19:16:49',NULL,'628629e3-79f5-4255-b981-d14c6c7856be','default','\"\"'),
  ('898b92fe-9d78-4838-a5c2-a674cc6fc845','2018-08-22 19:16:49','2018-08-22 19:16:49',NULL,'c5ddfe15-24d9-47f8-8ffe-f6b7daa9cf4a','default','\"\"'),
  ('905f3ca5-e41f-4d88-a1d9-abbef29a3d57','2018-08-22 19:16:49','2018-08-22 19:16:49',NULL,'4bc59b9a-8520-409f-85da-1c7552315863','mysql-n1-highmem-16','{\"max_disk_size\":\"10000\",\"pricing_plan\":\"PACKAGE\",\"tier\":\"db-n1-highmem-16\"}'),
  ('9824e39b-dbc0-4f95-a96a-1b02d4820b61','2018-08-22 19:16:49','2018-08-22 19:16:49',NULL,'cbad6d78-a73c-432d-b8ff-b219a17a803a','postgres-micro-dev','{\"max_disk_size\":\"100\",\"pricing_plan\":\"PER_USE\",\"tier\":\"db-f1-micro\"}'),('9ddcfbf4-cde8-4598-aeab-18a31304dc15','2018-08-22 19:16:49','2018-08-22 19:16:49',NULL,'b8e19880-ac58-42ef-b033-f7cd9c94d1fe','bigtable-hdd-10','{\"num_nodes\":\"3\",\"storage_type\":\"HDD\"}'),
  ('aca09d2b-2e70-4da3-a94a-a49c5b9b2f1a','2018-08-22 19:16:49','2018-08-22 19:16:49',NULL,'76d4abb2-fee7-4c8f-aee1-bcea2837f02b','default','\"\"'),
  ('ae472c3f-ac49-4094-ab31-e59b80ced397','2018-08-22 19:16:49','2018-08-22 19:16:49',NULL,'f80c0a3e-bd4d-4809-a900-b4e33a6450f1','default','\"\"'),
  ('bdcef6e7-e546-486c-ad54-af732b5ba840','2018-08-22 19:16:49','2018-08-22 19:16:49',NULL,'b9e4332e-b42b-4680-bda5-ea1506797474','standard','{\"storage_class\":\"STANDARD\"}'),
  ('c07ab411-cf55-48a6-adeb-ff17c632a043','2018-08-22 19:16:49','2018-08-22 19:16:49',NULL,'51b3e27e-d323-49ce-8c5f-1211e6409e82','spanner-regional-3','{\"num_nodes\":\"3\"}'),('c2fa09d4-8b0e-4730-a515-03ab19ad5c60','2018-08-22 19:16:49','2018-08-22 19:16:49',NULL,'b8e19880-ac58-42ef-b033-f7cd9c94d1fe','bigtable-ssd-30','{\"num_nodes\":\"30\",\"storage_type\":\"SDD\"}'),
  ('ccb53708-2399-449d-ac38-e0718872c867','2018-08-22 19:16:49','2018-08-22 19:16:49',NULL,'51b3e27e-d323-49ce-8c5f-1211e6409e82','spanner-regional-micro-dev','{\"num_nodes\":\"1\"}'),
  ('d90782de-a504-4f4a-a9f8-8b20d71ca4ad','2018-08-22 19:16:49','2018-08-22 19:16:49',NULL,'cbad6d78-a73c-432d-b8ff-b219a17a803a','postgres-n1-standard-16','{\"max_disk_size\":\"10000\",\"pricing_plan\":\"PER_USE\",\"tier\":\"db-n1-standard-16\"}'),('f29b4519-229f-44dc-a2fb-593f451f8b7a','2018-08-22 19:16:49','2018-08-22 19:16:49',NULL,'4bc59b9a-8520-409f-85da-1c7552315863','mysql-n1-standard-16','{\"max_disk_size\":\"10000\",\"pricing_plan\":\"PACKAGE\",\"tier\":\"db-n1-standard-16\"}');`

func setup3xDatabase(t *testing.T) {
	os.Remove("test.db")
	testDb, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		t.Errorf("Error setting up testing database %s", err)
	}

	db_service.RunMigrations(testDb)
	db_service.DbConnection = testDb

	testDb.Exec(planDetails)
}

type FakeBroker struct{}

func (t *FakeBroker) LastOperation(ctx context.Context, instanceID, operationData string) (brokerapi.LastOperation, error) {
	return brokerapi.LastOperation{}, nil
}

func (t *FakeBroker) Provision(ctx context.Context, instanceID string, details brokerapi.ProvisionDetails, asyncAllowed bool) (brokerapi.ProvisionedServiceSpec, error) {
	return brokerapi.ProvisionedServiceSpec{}, nil
}

func (t *FakeBroker) Deprovision(ctx context.Context, instanceID string, details brokerapi.DeprovisionDetails, asyncAllowed bool) (brokerapi.DeprovisionServiceSpec, error) {
	return brokerapi.DeprovisionServiceSpec{}, nil
}

func (t *FakeBroker) Bind(ctx context.Context, instanceID, bindingID string, details brokerapi.BindDetails) (brokerapi.Binding, error) {
	return brokerapi.Binding{}, nil
}

func (t *FakeBroker) Unbind(ctx context.Context, instanceID, bindingID string, details brokerapi.UnbindDetails) error {
	return nil
}

func (t *FakeBroker) Services(ctx context.Context) ([]brokerapi.Service, error) {
	return []brokerapi.Service{}, nil
}

func (t *FakeBroker) Update(ctx context.Context, instanceID string, details brokerapi.UpdateDetails, asyncAllowed bool) (brokerapi.UpdateServiceSpec, error) {
	return brokerapi.UpdateServiceSpec{}, nil
}

func TestThreeToFour_Update(t *testing.T) {
	setup3xDatabase(t)
	defer os.Remove("test.db")

	broker := NewLegacyPlanUpgrader(&FakeBroker{})

	cases := map[string]struct {
		ServiceId    string
		PlanId       string
		NewPlanId    string
		ExpectUpdate bool
		ErrContains  string
	}{
		"valid-stackdriver-debugger-default-upgrade": {"83837945-1547-41e0-b661-ea31d76eed11", "11308f80-dc53-44f7-a49b-c3e65258b421", "10866183-a775-49e8-96e3-4e7a901e4a79", true, ""},
		"valid-nearline-upgrade":                     {"b9e4332e-b42b-4680-bda5-ea1506797474", "1228b43c-15ca-43e6-ae2d-70a33e47efa1", "a42c1182-d1a0-4d40-82c1-28220518b360", true, ""},
		"valid-default-ml-upgrade":                   {"5ad2dce0-51f7-4ede-8b46-293d6df1e8d4", "5dd56dd9-cc53-4b71-aa4e-45d9c1266dbb", "be7954e1-ecfb-4936-a0b6-db35e6424c7a", true, ""},
		"valid-reduced_availability-upgrade":         {"b9e4332e-b42b-4680-bda5-ea1506797474", "6c36cc15-cb2e-494d-a0e7-f4a34e0ff501", "1a1f4fe6-1904-44d0-838c-4c87a9490a6b", true, ""},
		"valid-default-pubsub-upgrade":               {"628629e3-79f5-4255-b981-d14c6c7856be", "728b61e5-2b7b-4a2e-a2a1-88797c6db1c2", "622f4da3-8731-492a-af29-66a9146f8333", true, ""},
		"valid-default-stackdriver-trace-upgrade":    {"c5ddfe15-24d9-47f8-8ffe-f6b7daa9cf4a", "898b92fe-9d78-4838-a5c2-a674cc6fc845", "ab6c2287-b4bc-4ff4-a36a-0575e7910164", true, ""},
		"valid-datastore-default-datastore-upgrade":  {"76d4abb2-fee7-4c8f-aee1-bcea2837f02b", "aca09d2b-2e70-4da3-a94a-a49c5b9b2f1a", "05f1fb6b-b5f0-48a2-9c2b-a5f236507a97", true, ""},
		"valid-bq-default-upgrade":                   {"f80c0a3e-bd4d-4809-a900-b4e33a6450f1", "ae472c3f-ac49-4094-ab31-e59b80ced397", "10ff4e72-6e84-44eb-851f-bdb38a791914", true, ""},
		"valid-gcs-standard-upgrade":                 {"b9e4332e-b42b-4680-bda5-ea1506797474", "bdcef6e7-e546-486c-ad54-af732b5ba840", "e1d11f65-da66-46ad-977c-6d56513baf43", true, ""},

		"invalid-gcs-standard-to-nearline": {"b9e4332e-b42b-4680-bda5-ea1506797474", "bdcef6e7-e546-486c-ad54-af732b5ba840", "a42c1182-d1a0-4d40-82c1-28220518b360", false, "can only upgrade this legacy plan to \"standard\""},
		"invalid-legacy-to-legacy":         {"b9e4332e-b42b-4680-bda5-ea1506797474", "bdcef6e7-e546-486c-ad54-af732b5ba840", "6c36cc15-cb2e-494d-a0e7-f4a34e0ff501", false, "can only upgrade this legacy plan to \"standard\""},
		"not-legacy-plan":                  {"b9e4332e-b42b-4680-bda5-ea1506797474", "e1d11f65-da66-46ad-977c-6d56513baf43", "some-other-plan", false, ""},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			db_service.CreateServiceInstanceDetails(&models.ServiceInstanceDetails{
				ID:        tn,
				ServiceId: tc.ServiceId,
				PlanId:    tc.PlanId,
			})

			_, err := broker.Update(context.Background(), tn, brokerapi.UpdateDetails{PlanID: tc.NewPlanId}, true)
			checkErrorMatches(t, err, tc.ErrContains)

			details, err := db_service.GetServiceInstanceDetailsById(tn)
			if err != nil {
				t.Errorf("Error getting details: %s", err)
			}

			planWasUpdated := details.PlanId == tc.NewPlanId
			if planWasUpdated != tc.ExpectUpdate {
				t.Errorf("Unexpected plan state, expected update? %v, got: %v", tc.ExpectUpdate, planWasUpdated)
			}
		})
	}
}

func TestThreeToFour_migrationErrorMessage(t *testing.T) {
	setup3xDatabase(t)
	defer os.Remove("test.db")

	broker := NewLegacyPlanUpgrader(&FakeBroker{})

	cases := map[string]struct {
		ServiceId   string
		PlanId      string
		ErrContains string
	}{
		"upgrade-needed":    {"83837945-1547-41e0-b661-ea31d76eed11", "11308f80-dc53-44f7-a49b-c3e65258b421", "cf update-service my-service -p default"},
		"no-upgrade-needed": {"b9e4332e-b42b-4680-bda5-ea1506797474", "a42c1182-d1a0-4d40-82c1-28220518b360", ""},
	}

	for tn, tc := range cases {
		db_service.CreateServiceInstanceDetails(&models.ServiceInstanceDetails{
			ID:        tn,
			Name:      "my-service",
			ServiceId: tc.ServiceId,
			PlanId:    tc.PlanId,
		})

		t.Run(tn+"-deprovision", func(t *testing.T) {
			_, err := broker.Deprovision(context.Background(), tn, brokerapi.DeprovisionDetails{PlanID: tc.PlanId, ServiceID: tc.ServiceId}, true)
			checkErrorMatches(t, err, tc.ErrContains)
		})

		t.Run(tn+"-bind", func(t *testing.T) {
			_, err := broker.Bind(context.Background(), tn, tn, brokerapi.BindDetails{PlanID: tc.PlanId, ServiceID: tc.ServiceId})
			checkErrorMatches(t, err, tc.ErrContains)
		})

		t.Run(tn+"-unbind", func(t *testing.T) {
			err := broker.Unbind(context.Background(), tn, tn, brokerapi.UnbindDetails{PlanID: tc.PlanId, ServiceID: tc.ServiceId})
			checkErrorMatches(t, err, tc.ErrContains)
		})
	}
}

func TestThreeToFour_Provision(t *testing.T) {
	setup3xDatabase(t)
	defer os.Remove("test.db")

	broker := NewLegacyPlanUpgrader(&FakeBroker{})

	cases := map[string]struct {
		ServiceId   string
		PlanId      string
		ErrContains string
	}{
		"legacy-plan":  {"83837945-1547-41e0-b661-ea31d76eed11", "11308f80-dc53-44f7-a49b-c3e65258b421", "The plan \"legacy3-default\" is only availble for compatibility purposes"},
		"current-plan": {"b9e4332e-b42b-4680-bda5-ea1506797474", "a42c1182-d1a0-4d40-82c1-28220518b360", ""},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			_, err := broker.Provision(context.Background(), tn, brokerapi.ProvisionDetails{PlanID: tc.PlanId, ServiceID: tc.ServiceId}, true)
			checkErrorMatches(t, err, tc.ErrContains)
		})
	}
}

func TestThreeToFour_Services(t *testing.T) {
	setup3xDatabase(t)
	defer os.Remove("test.db")

	broker := NewLegacyPlanUpgrader(&FakeBroker{})
	services, _ := broker.Services(context.Background())

	cases := map[string]struct {
		ServiceId string
		PlanId    string
	}{
		"stackdriver-debugger":        {"83837945-1547-41e0-b661-ea31d76eed11", "11308f80-dc53-44f7-a49b-c3e65258b421"},
		"nearline":                    {"b9e4332e-b42b-4680-bda5-ea1506797474", "1228b43c-15ca-43e6-ae2d-70a33e47efa1"},
		"default-ml":                  {"5ad2dce0-51f7-4ede-8b46-293d6df1e8d4", "5dd56dd9-cc53-4b71-aa4e-45d9c1266dbb"},
		"reduced_availability":        {"b9e4332e-b42b-4680-bda5-ea1506797474", "6c36cc15-cb2e-494d-a0e7-f4a34e0ff501"},
		"default-pubsub-upgrade":      {"628629e3-79f5-4255-b981-d14c6c7856be", "728b61e5-2b7b-4a2e-a2a1-88797c6db1c2"},
		"default-stackdriver-trace":   {"c5ddfe15-24d9-47f8-8ffe-f6b7daa9cf4a", "898b92fe-9d78-4838-a5c2-a674cc6fc845"},
		"datastore-default-datastore": {"76d4abb2-fee7-4c8f-aee1-bcea2837f02b", "aca09d2b-2e70-4da3-a94a-a49c5b9b2f1a"},
		"bq-default":                  {"f80c0a3e-bd4d-4809-a900-b4e33a6450f1", "ae472c3f-ac49-4094-ab31-e59b80ced397"},
		"gcs-standard":                {"b9e4332e-b42b-4680-bda5-ea1506797474", "bdcef6e7-e546-486c-ad54-af732b5ba840"},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			for _, service := range services {
				if service.ID != tc.ServiceId {
					continue
				}
				for _, plan := range service.Plans {
					if plan.ID == tc.PlanId {
						return
					}
				}
			}

			t.Errorf("Service/Plan pair not found in list %s/%s", tc.ServiceId, tc.PlanId)
		})
	}
}

func checkErrorMatches(t *testing.T, err error, matches string) {
	hasErr := err != nil
	expectingErr := matches != ""

	if hasErr != expectingErr {
		t.Fatalf("Expecting err? %v, got: %v", expectingErr, err)
	}

	if expectingErr && !strings.Contains(err.Error(), matches) {
		t.Fatalf("Wrong error, expected to contain %q, got: %v", matches, err)
	}
}
