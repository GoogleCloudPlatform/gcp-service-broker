package tf

import (
	"log"

	accountmanagers "github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/account_managers"

	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/models"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/broker"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/validation"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/varcontext"
	"github.com/pivotal-cf/brokerapi"
)

func init() {
	roleWhitelist := []string{
		"storage.objectCreator",
		"storage.objectViewer",
		"storage.objectAdmin",
	}

	cloudStorage := TfServiceDefinitionV1{
		Version:          1,
		Name:             "google-storage-experimental",
		Id:               "68d094ae-e727-4c14-af07-ee34133c8dfb",
		Description:      "Unified object storage for developers and enterprises. Cloud Storage allows world-wide storage and retrieval of any amount of data at any time.",
		DisplayName:      "Google Cloud Storage 2",
		ImageUrl:         "",
		DocumentationUrl: "",
		SupportUrl:       "",
		Tags:             []string{"preview", "gcp", "terraform", "storage"},
		Plans: []broker.ServicePlan{
			{
				ServicePlan: brokerapi.ServicePlan{
					ID:   "e1d11f65-da66-46ad-977c-6d56513baf43",
					Name: "standard",
					Metadata: &brokerapi.ServicePlanMetadata{
						DisplayName: "Standard",
					},
					Description: "Standard storage class.",
				},
				ServiceProperties: map[string]string{
					"storage_class": "STANDARD",
				},
			},
		},
		ProvisionSettings: TfServiceDefinitionV1Action{
			PlanInputs: []broker.BrokerVariable{
				{
					FieldName: "storage_class",
					Type:      broker.JsonTypeString,
					Details:   "The storage class of the bucket. See: https://cloud.google.com/storage/docs/storage-classes.",
					Required:  true,
				},
			},
			UserInputs: []broker.BrokerVariable{
				{
					FieldName: "name",
					Type:      broker.JsonTypeString,
					Details:   "The name of the bucket. There is a single global namespace shared by all buckets so it MUST be unique.",
					Default:   "pcf_sb_${counter.next()}_${time.nano()}",
					Constraints: validation.NewConstraintBuilder(). // https://cloud.google.com/storage/docs/naming
											Pattern("^[A-Za-z0-9_\\.]+$").
											MinLength(3).
											MaxLength(222).
											Build(),
				},
				{
					FieldName: "location",
					Type:      broker.JsonTypeString,
					Default:   "US",
					Details:   `The location of the bucket. Object data for objects in the bucket resides in physical storage within this region. See: https://cloud.google.com/storage/docs/bucket-locations`,
					Constraints: validation.NewConstraintBuilder().
						Pattern("^[A-Za-z][-a-z0-9A-Z]+$").
						Examples("US", "EU", "southamerica-east1").
						Build(),
				},
			},
			Computed: []varcontext.DefaultVariable{
				{Name: "labels", Default: "${json.marshal(request.default_labels)}", Overwrite: true},
			},
			Template: `
	variable name {type = "string"}
	variable location {type = "string"}
	variable storage_class {type = "string"}

	resource "google_storage_bucket" "bucket" {
	  name     = "${var.name}"
	  location = "${var.location}"
	  storage_class = "${var.storage_class}"
	}

	output id {value = "${google_storage_bucket.bucket.id}"}
	output bucket_name {value = "${var.name}"}
	`,
			Outputs: []broker.BrokerVariable{
				{
					FieldName: "bucket_name",
					Type:      broker.JsonTypeString,
					Details:   "Name of the bucket this binding is for.",
					Required:  true,
					Constraints: validation.NewConstraintBuilder(). // https://cloud.google.com/storage/docs/naming
											Pattern("^[A-Za-z0-9_\\.]+$").
											MinLength(3).
											MaxLength(222).
											Build(),
				},
				{
					FieldName:   "id",
					Type:        broker.JsonTypeString,
					Details:     "The GCP ID of this bucket.",
					Required:    true,
					Constraints: validation.NewConstraintBuilder().Build(),
				},
			},
		},
		BindSettings: TfServiceDefinitionV1Action{
			PlanInputs: []broker.BrokerVariable{},
			UserInputs: accountmanagers.ServiceAccountBindInputVariables(models.StorageName, roleWhitelist),
			Computed: append(accountmanagers.ServiceAccountBindComputedVariables(),
				varcontext.DefaultVariable{
					Name:      "bucket",
					Default:   `${instance.details["bucket_name"]}`,
					Overwrite: true,
				}),
			Template: `
	variable role {type = "string"}
	variable service_account_name {type = "string"}
	variable service_account_display_name {type = "string"}
	variable bucket {type = "string"}

	resource "google_service_account" "account" {
	  account_id = "${var.service_account_name}"
	  display_name = "${var.service_account_display_name}"
	}

	resource "google_service_account_key" "key" {
	  service_account_id = "${google_service_account.account.name}"
	}

	resource "google_storage_bucket_iam_member" "member" {
	  bucket = "${var.bucket}"
	  role   = "roles/${var.role}"
	  member = "serviceAccount:${google_service_account.account.email}"
	}

	output "Name" {value = "${google_service_account.account.display_name}"}
	output "Email" {value = "${google_service_account.account.email}"}
	output "UniqueId" {value = "${google_service_account.account.unique_id}"}
	output "PrivateKeyData" {value = "${google_service_account_key.key.private_key}"}
	output "ProjectId" {value = "${google_service_account.account.project}"}
	`,
			Outputs: accountmanagers.ServiceAccountBindOutputVariables(),
		},

		Examples: []broker.ServiceExample{
			{
				Name:            "Basic Configuration",
				Description:     "Create a bucket with a service account that can create/read/delete the objects in it.",
				PlanId:          "e1d11f65-da66-46ad-977c-6d56513baf43",
				ProvisionParams: map[string]interface{}{"location": "us"},
				BindParams: map[string]interface{}{
					"role": "storage.objectAdmin",
				},
			},
		},

		// Internal SHOULD be set to true for Google maintained services.
		Internal: true,
	}

	service, err := cloudStorage.ToService()
	if err != nil {
		log.Fatal(err)
	}
	broker.Register(service)
}
