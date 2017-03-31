package integration_tests

import (
	. "gcp-service-broker/brokerapi/brokers"

	"golang.org/x/net/context"

	"fmt"
	"gcp-service-broker/brokerapi/brokers"
	"gcp-service-broker/brokerapi/brokers/models"
	"gcp-service-broker/brokerapi/brokers/name_generator"
	"gcp-service-broker/db_service"
	"gcp-service-broker/fakes"
	"hash/crc32"
	"net/http"
	"os"

	googlepubsub "cloud.google.com/go/pubsub"

	"encoding/json"

	googlebigtable "cloud.google.com/go/bigtable"
	googlestorage "cloud.google.com/go/storage"
	"code.cloudfoundry.org/lager"
	"github.com/jinzhu/gorm"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	googlebigquery "google.golang.org/api/bigquery/v2"
	"google.golang.org/api/iam/v1"
	"google.golang.org/api/option"
)

const timeout = 60

type genericService struct {
	serviceId              string
	planId                 string
	bindingId              string
	rawBindingParams       map[string]interface{}
	instanceId             string
	serviceExistsFn        func(bool) bool
	cleanupFn              func()
	serviceMetadataSavedFn func(string) bool
}

type iamService struct {
	bindingId string
	serviceId string
	planId    string
}

func getAndUnmarshalInstanceDetails(instanceId string) map[string]string {
	var instanceRecord models.ServiceInstanceDetails
	db_service.DbConnection.Find(&instanceRecord).Where("id = ?", instanceId)
	var instanceDetails map[string]string
	json.Unmarshal([]byte(instanceRecord.OtherDetails), &instanceDetails)
	return instanceDetails
}

func testGenericService(gcpBroker *GCPAsyncServiceBroker, params *genericService) {
	// If the service already exists (eg, failed previous test), clean it up before the run
	if params.serviceExistsFn != nil && params.serviceExistsFn(false) {
		params.cleanupFn()
	}
	//
	// Provision
	//
	provisionDetails := models.ProvisionDetails{
		ServiceID: params.serviceId,
		PlanID:    params.planId,
	}

	_, err := gcpBroker.Provision(params.instanceId, provisionDetails, true)
	Expect(err).ToNot(HaveOccurred())

	// Provision is registered in the database
	var count int
	db_service.DbConnection.Model(&models.ServiceInstanceDetails{}).Where("id = ?", params.instanceId).Count(&count)
	Expect(count).To(Equal(1))

	if params.serviceExistsFn != nil {
		Expect(params.serviceExistsFn(true)).To(BeTrue())
	}
	Expect(params.serviceMetadataSavedFn(params.instanceId)).To(BeTrue())

	//
	// Bind
	//
	bindDetails := models.BindDetails{
		ServiceID:  params.serviceId,
		PlanID:     params.planId,
		Parameters: params.rawBindingParams,
	}
	creds, err := gcpBroker.Bind(params.instanceId, params.bindingId, bindDetails)
	Expect(err).ToNot(HaveOccurred())

	db_service.DbConnection.Model(&models.ServiceBindingCredentials{}).Where("binding_id = ?", params.bindingId).Count(&count)
	Expect(count).To(Equal(1))

	iamService, err := iam.New(gcpBroker.GCPClient)
	Expect(err).ToNot(HaveOccurred())
	saService := iam.NewProjectsServiceAccountsService(iamService)
	resourceName := "projects/" + gcpBroker.RootGCPCredentials.ProjectId + "/serviceAccounts/" + creds.Credentials.(map[string]string)["UniqueId"]
	_, err = saService.Get(resourceName).Do()
	Expect(err).ToNot(HaveOccurred())

	//
	// Unbind
	//
	unbindDetails := models.UnbindDetails{
		ServiceID: params.serviceId,
		PlanID:    params.planId,
	}
	err = gcpBroker.Unbind(params.instanceId, params.bindingId, unbindDetails)
	Expect(err).ToNot(HaveOccurred())

	binding := models.ServiceBindingCredentials{}
	if err := db_service.DbConnection.Unscoped().Where("binding_id = ?", params.bindingId).First(&binding).Error; err != nil {
		panic("error checking for binding details: " + err.Error())
	}
	Expect(binding.DeletedAt).NotTo(Equal(nil))

	_, err = saService.Get(resourceName).Do()
	Expect(err).To(HaveOccurred())

	//
	// Deprovision

	deprovisionDetails := models.DeprovisionDetails{
		ServiceID: params.serviceId,
		PlanID:    params.planId,
	}
	_, err = gcpBroker.Deprovision(params.instanceId, deprovisionDetails, true)
	Expect(err).ToNot(HaveOccurred())
	instance := models.ServiceInstanceDetails{}
	if err := db_service.DbConnection.Unscoped().Where("ID = ?", params.instanceId).First(&instance).Error; err != nil {
		panic("error checking for service instance details: " + err.Error())
	}
	Expect(instance.DeletedAt).NotTo(Equal(nil))

	if params.serviceExistsFn != nil {
		Expect(params.serviceExistsFn(false)).To(BeFalse())
	}
}

// For services that only create a service account and bind those credentials.
func testIamBasedService(gcpBroker *GCPAsyncServiceBroker, params *iamService) {
	genericServiceParams := &genericService{
		serviceId:        params.serviceId,
		planId:           params.planId,
		instanceId:       "iam-instance",
		bindingId:        "iam-instance",
		rawBindingParams: map[string]interface{}{},
		serviceMetadataSavedFn: func(instanceId string) bool {
			// Metadata should be empty, there is no additional information required
			instanceDetails := getAndUnmarshalInstanceDetails(instanceId)
			return len(instanceDetails) == 0
		},
	}

	testGenericService(gcpBroker, genericServiceParams)
}

// Instance Name is used to name every instance created in GCP (eg, a storage bucket)
// The name should be consistent between runs to ensure there's bounds to the resources it creates
// and to have some insurance that they are properly destroyed.
//
// Why:
// - If we allow it to generate a random instance name every time the test will
//   not fail if the resource existed before hand.
// - If we always use a static one, globally named resources (eg, a storage bucket)
//   would fail to create when two different projects run these tests.
func generateInstanceName(projectId string, sep string) string {
	hashed := crc32.ChecksumIEEE([]byte(projectId))
	if sep == "" {
		sep = "_"
	}
	return fmt.Sprintf("pcf%ssb%s1%s%d", sep, sep, sep, hashed)
}

var _ = Describe("LiveIntegrationTests", func() {
	var (
		gcpBroker           *GCPAsyncServiceBroker
		err                 error
		logger              lager.Logger
		serviceNameToId     map[string]string = make(map[string]string)
		serviceNameToPlanId map[string]string = make(map[string]string)
		instance_name       string
	)

	BeforeEach(func() {
		logger = lager.NewLogger("brokers_test")
		logger.RegisterSink(lager.NewWriterSink(GinkgoWriter, lager.DEBUG))

		testDb, _ := gorm.Open("sqlite3", "test.db")
		testDb.CreateTable(models.ServiceInstanceDetails{})
		testDb.CreateTable(models.ServiceBindingCredentials{})
		testDb.CreateTable(models.PlanDetails{})
		testDb.CreateTable(models.ProvisionRequestDetails{})

		db_service.DbConnection = testDb

		os.Setenv("SECURITY_USER_NAME", "username")
		os.Setenv("SECURITY_USER_PASSWORD", "password")
		os.Setenv("SERVICES", fakes.Services)
		os.Setenv("PRECONFIGURED_PLANS", fakes.PreconfiguredPlans)

		os.Setenv("CLOUDSQL_CUSTOM_PLANS", fakes.TestCloudSQLPlan)
		os.Setenv("BIGTABLE_CUSTOM_PLANS", fakes.TestBigtablePlan)
		os.Setenv("SPANNER_CUSTOM_PLANS", fakes.TestSpannerPlan)

		var creds models.GCPCredentials
		creds, err = brokers.GetCredentialsFromEnv()
		if err != nil {
			logger.Error("error", err)
		}
		instance_name = generateInstanceName(creds.ProjectId, "")
		name_generator.Basic = &fakes.StaticNameGenerator{Val: instance_name}

		gcpBroker, err = brokers.New(logger)
		if err != nil {
			logger.Error("error", err)
		}

		for _, service := range *gcpBroker.Catalog {
			serviceNameToId[service.Name] = service.ID
			serviceNameToPlanId[service.Name] = service.Plans[0].ID
		}
	})

	Describe("Broker init", func() {
		It("should have 8 services in sevices map", func() {
			Expect(len(gcpBroker.ServiceBrokerMap)).To(Equal(8))
		})

		It("should have a default client", func() {
			Expect(gcpBroker.GCPClient).NotTo(Equal(&http.Client{}))
		})

		It("should have loaded credentials correctly and have a project id", func() {
			Expect(gcpBroker.RootGCPCredentials.ProjectId).ToNot(BeEmpty())
		})
	})

	Describe("getting broker catalog", func() {
		It("should have 8 services available", func() {
			Expect(len(gcpBroker.Services())).To(Equal(8))
		})

		It("should have 3 storage plans available", func() {
			serviceList := gcpBroker.Services()
			for _, s := range serviceList {
				if s.ID == serviceNameToId[models.StorageName] {
					Expect(len(s.Plans)).To(Equal(3))
				}
			}

		})
	})

	Describe("bigquery", func() {
		It("can provision/bind/unbind/deprovision", func() {
			service, err := googlebigquery.New(gcpBroker.GCPClient)
			Expect(err).NotTo(HaveOccurred())

			params := &genericService{
				serviceId:  serviceNameToId[models.BigqueryName],
				planId:     serviceNameToPlanId[models.BigqueryName],
				bindingId:  "integration_test_bind",
				instanceId: "integration_test_dataset",
				rawBindingParams: map[string]interface{}{
					"role": "bigquery.admin",
				},
				serviceExistsFn: func(expected bool) bool {
					_, err = service.Datasets.Get(gcpBroker.RootGCPCredentials.ProjectId, instance_name).Do()

					return err == nil
				},
				serviceMetadataSavedFn: func(instanceId string) bool {
					instanceDetails := getAndUnmarshalInstanceDetails(instanceId)
					return instanceDetails["dataset_id"] != ""
				},
				cleanupFn: func() {
					err := service.Datasets.Delete(gcpBroker.RootGCPCredentials.ProjectId, instance_name).Do()
					Expect(err).NotTo(HaveOccurred())
				},
			}
			testGenericService(gcpBroker, params)
		}, timeout)
	})

	Describe("bigtable", func() {
		var bigtableInstanceName string
		BeforeEach(func() {
			bigtableInstanceName = generateInstanceName(gcpBroker.RootGCPCredentials.ProjectId, "-")
			name_generator.Basic = &fakes.StaticNameGenerator{Val: bigtableInstanceName}
		})

		AfterEach(func() {
			name_generator.Basic = &fakes.StaticNameGenerator{Val: instance_name}
		})

		It("can provision/bind/unbind/deprovision", func() {

			ctx := context.Background()
			service, err := googlebigtable.NewInstanceAdminClient(ctx, gcpBroker.RootGCPCredentials.ProjectId)
			Expect(err).NotTo(HaveOccurred())

			params := &genericService{
				serviceId:  serviceNameToId[models.BigtableName],
				planId:     serviceNameToPlanId[models.BigtableName],
				bindingId:  "integration_test_bind",
				instanceId: "integration_test_instance",
				rawBindingParams: map[string]interface{}{
					"role": "editor",
				},
				serviceExistsFn: func(expected bool) bool {
					instances, err := service.Instances(ctx)

					return err == nil && len(instances) == 1 && instances[0].Name == bigtableInstanceName
				},
				serviceMetadataSavedFn: func(instanceId string) bool {
					instanceDetails := getAndUnmarshalInstanceDetails(instanceId)
					return instanceDetails["instance_id"] != ""
				},
				cleanupFn: func() {
					err := service.DeleteInstance(ctx, bigtableInstanceName)
					Expect(err).NotTo(HaveOccurred())
				},
			}
			testGenericService(gcpBroker, params)
		}, timeout)
	})

	Describe("cloud storage", func() {
		It("can provision/bind/unbind/deprovision", func() {
			service, err := googlestorage.NewClient(context.Background(), option.WithUserAgent(models.CustomUserAgent))
			Expect(err).NotTo(HaveOccurred())

			params := &genericService{
				serviceId:  serviceNameToId[models.StorageName],
				planId:     serviceNameToPlanId[models.StorageName],
				instanceId: "integration_test_bucket",
				bindingId:  "integration_test_bucket_binding",
				rawBindingParams: map[string]interface{}{
					"role": "storage.admin",
				},
				serviceExistsFn: func(bool) bool {
					bucket := service.Bucket(instance_name)
					_, err = bucket.Attrs(context.Background())

					return err == nil
				},
				serviceMetadataSavedFn: func(instanceId string) bool {
					instanceDetails := getAndUnmarshalInstanceDetails(instanceId)
					return instanceDetails["bucket_name"] != ""
				},
				cleanupFn: func() {
					bucket := service.Bucket(instance_name)
					bucket.Delete(context.Background())
				},
			}

			testGenericService(gcpBroker, params)
		}, timeout)
	})

	Describe("pub sub", func() {
		It("can provision/bind/unbind/deprovision", func() {
			service, err := googlepubsub.NewClient(context.Background(), gcpBroker.RootGCPCredentials.ProjectId, option.WithUserAgent(models.CustomUserAgent))
			Expect(err).NotTo(HaveOccurred())

			topic := service.Topic(instance_name)

			params := &genericService{
				serviceId:  serviceNameToId[models.PubsubName],
				planId:     serviceNameToPlanId[models.PubsubName],
				instanceId: "integration_test_topic",
				bindingId:  "integration_test_topic_bindingId",
				rawBindingParams: map[string]interface{}{
					"role": "pubsub.admin",
				},
				serviceExistsFn: func(bool) bool {
					exists, err := topic.Exists(context.Background())
					return exists && err == nil
				},
				serviceMetadataSavedFn: func(instanceId string) bool {
					instanceDetails := getAndUnmarshalInstanceDetails(instanceId)
					return instanceDetails["topic_name"] != ""
				},
				cleanupFn: func() {
					err := topic.Delete(context.Background())
					Expect(err).NotTo(HaveOccurred())
				},
			}

			testGenericService(gcpBroker, params)
		}, timeout)
	})

	Describe("stadkdriver debugger", func() {
		It("can provision/bind/unbind/deprovision", func() {
			params := &iamService{
				serviceId: serviceNameToId[models.StackdriverDebuggerName],
				planId:    serviceNameToPlanId[models.StackdriverDebuggerName],
			}
			testIamBasedService(gcpBroker, params)
		}, timeout)
	})

	AfterEach(func() {
		os.Remove(models.AppCredsFileName)
		os.Remove("test.db")
	})
})
