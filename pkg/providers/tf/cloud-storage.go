package tf

import (
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/broker"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/varcontext"
	"github.com/pivotal-cf/brokerapi"
)

var cloudStorage = TfServiceDefinitionV1{
	Version:          1,
	Name:             "google-storage-2",
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
	ProvisionSettings: &TfServiceDefinitionV1Action{
		PlanInputs: []broker.BrokerVariable{},
		UserInputs: []broker.BrokerVariable{},
		Computed:   []varcontext.DefaultVariable{},
		Template:   ``,
		Outputs:    []broker.BrokerVariable{},
	},
	// BindSettings      TfServiceDefinitionV1Action `yaml:"bind" validate:"required,dive"`
	// Examples          []broker.ServiceExample     `yaml:"examples" validate:"required,dive"`

	// Internal SHOULD be set to true for Google maintained services.
	Internal: true,
}
