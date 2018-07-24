package broker

var bs = &BrokerService{
	Name: "cloud-storage",
	DefaultServiceDefinition: `{
        "id": "b9e4332e-b42b-4680-bda5-ea1506797474",
        "description": "A Powerful, Simple and Cost Effective Object Storage Service",
        "name": "google-storage",
        "bindable": true,
        "plan_updateable": false,
        "metadata": {
          "displayName": "Google Cloud Storage",
          "longDescription": "A Powerful, Simple and Cost Effective Object Storage Service",
          "documentationUrl": "https://cloud.google.com/storage/docs/overview",
          "supportUrl": "https://cloud.google.com/support/",
          "imageUrl": "https://cloud.google.com/_static/images/cloud/products/logos/svg/storage.svg"
        },
        "tags": ["gcp", "storage"],
        "plans": [
          {
            "id": "e1d11f65-da66-46ad-977c-6d56513baf43",
            "service_id": "b9e4332e-b42b-4680-bda5-ea1506797474",
            "name": "standard",
            "display_name": "Standard",
            "description": "Standard storage class",
            "service_properties": {"storage_class": "STANDARD"}
          },
          {
            "id": "a42c1182-d1a0-4d40-82c1-28220518b360",
            "service_id": "b9e4332e-b42b-4680-bda5-ea1506797474",
            "name": "nearline",
            "display_name": "Nearline",
            "description": "Nearline storage class",
            "service_properties": {"storage_class": "NEARLINE"}
          },
          {
            "id": "1a1f4fe6-1904-44d0-838c-4c87a9490a6b",
            "service_id": "b9e4332e-b42b-4680-bda5-ea1506797474",
            "name": "reduced_availability",
            "display_name": "Durable Reduced Availability",
            "description": "Durable Reduced Availability storage class",
            "service_properties": {"storage_class": "DURABLE_REDUCED_AVAILABILITY"}
          }
        ]
      }`,
	ProvisionInputVariables: []BrokerVariable{
		BrokerVariable{
			FieldName: "name",
			Type:      JsonTypeString,
			Details:   "The name of the bucket. There is a single global namespace shared by all buckets so it MUST be unique.",
			Default:   "a generated value",
		},
		BrokerVariable{
			FieldName: "location",
			Type:      JsonTypeString,
			Default:   "us",
			Details:   `The location of the bucket. Object data for objects in the bucket resides in physical storage within this region. See https://cloud.google.com/storage/docs/bucket-locations`,
		},
	},
	BindInputVariables: []BrokerVariable{
		BrokerVariable{
			Required:  true,
			FieldName: "role",
			Type:      JsonTypeString,
			Details:   `The role for the account without the "roles/" prefix. See https://cloud.google.com/iam/docs/understanding-roles for available roles.`,
		},
	},
	BindOutputVariables: []BrokerVariable{
		BrokerVariable{
			FieldName: "Email",
			Type:      JsonTypeString,
			Details:   "Email address of the service account",
		},
		BrokerVariable{
			FieldName: "Name",
			Type:      JsonTypeString,
			Details:   "The name of the service account",
		},
		BrokerVariable{
			FieldName: "PrivateKeyData",
			Type:      JsonTypeString,
			Details:   "Service account private key data. Base-64 encoded JSON.",
		},
		BrokerVariable{
			FieldName: "ProjectId",
			Type:      JsonTypeString,
			Details:   "ID of the project that owns the service account",
		},
		BrokerVariable{
			FieldName: "UniqueId",
			Type:      JsonTypeString,
			Details:   "Unique and stable id of the service account",
		},
		BrokerVariable{
			FieldName: "bucket_name",
			Type:      JsonTypeString,
			Details:   "Name of the bucket this binding is for",
		},
	},

	Examples: []ServiceExample{
		ServiceExample{
			Name:            "Basic Configuration",
			Description:     "Create a nearline bucket with a service account that can create/read/delete the objects in it.",
			PlanId:          "a42c1182-d1a0-4d40-82c1-28220518b360",
			ProvisionParams: map[string]interface{}{"location": "us"},
			BindParams: map[string]interface{}{
				"role": "storage.objectAdmin",
			},
		},
	},
}

//
// type ServiceExample struct {
// 	// Name is a human-readable name of the example
// 	Name string
// 	// Descrpition is a long-form description of what this example is about
// 	Description string
// 	// PlanId is the plan this example will run against.
// 	PlanId string
//
// 	// ProvisionParams is the JSON object that will be passed to provision
// 	ProvisionParams map[string]interface{}
//
// 	// BindParams is the JSON object that will be passed to bind. If nil,
// 	// this example DOES NOT include a bind portion.
// 	BindParams map[string]interface{}
// }

func init() {
	Register(bs)
}
