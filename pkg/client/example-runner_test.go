package client

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"

	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/broker"
)

func ExampleGetAllCompleteServiceExamples_jsonSpec() {

	allExamples := []CompleteServiceExample{
		{
			ServiceExample: broker.ServiceExample{
				Name:            "Basic Configuration",
				Description:     "Creates an account with the permission `clouddebugger.agent`.",
				PlanId:          "10866183-a775-49e8-96e3-4e7a901e4a79",
				ProvisionParams: map[string]interface{}{},
				BindParams:      map[string]interface{}{}},
			ServiceName: "google-stackdriver-debugger",
			ServiceId:   "83837945-1547-41e0-b661-ea31d76eed11",
			ExpectedOutput: broker.CreateJsonSchema([]broker.BrokerVariable{
				{
					Required:  true,
					FieldName: "Email",
					Type:      "string",
					Details:   "Email address of the service account.",
				},
			}),
		},
	}

	b, err := json.MarshalIndent(allExamples, "", "\t")

	if err != nil {
		panic(err)
	}

	os.Stdout.Write(b)
}

func TestGetExamplesForAService(t *testing.T) {
	cases := map[string]struct {
		ServiceDefinition *broker.ServiceDefinition
		ExpectedResponse  []CompleteServiceExample
		ExpectedError     error
	}{
		"service with no examples": {
			ServiceDefinition: &broker.ServiceDefinition{
				Id: "TestService",
			},
			ExpectedResponse: []CompleteServiceExample(nil),
			ExpectedError:    nil,
		},
		"google-stackdriver-debugger": {
			ServiceDefinition: &broker.ServiceDefinition{
				Id:   "83837945-1547-41e0-b661-ea31d76eed11",
				Name: "google-stackdriver-debugger",
				BindOutputVariables: []broker.BrokerVariable{
					{
						Required:  true,
						FieldName: "Email",
						Type:      "string",
						Details:   "Email address of the service account.",
					},
				},
				Examples: []broker.ServiceExample{
					{
						Name:            "Basic Configuration",
						Description:     "Creates an account with the permission `clouddebugger.agent`.",
						PlanId:          "10866183-a775-49e8-96e3-4e7a901e4a79",
						ProvisionParams: map[string]interface{}{},
						BindParams:      map[string]interface{}{},
					},
				},
			},
			ExpectedResponse: []CompleteServiceExample{
				{
					ServiceExample: broker.ServiceExample{
						Name:            "Basic Configuration",
						Description:     "Creates an account with the permission `clouddebugger.agent`.",
						PlanId:          "10866183-a775-49e8-96e3-4e7a901e4a79",
						ProvisionParams: map[string]interface{}{},
						BindParams:      map[string]interface{}{}},
					ServiceName: "google-stackdriver-debugger",
					ServiceId:   "83837945-1547-41e0-b661-ea31d76eed11",
					ExpectedOutput: broker.CreateJsonSchema([]broker.BrokerVariable{
						{
							Required:  true,
							FieldName: "Email",
							Type:      "string",
							Details:   "Email address of the service account.",
						},
					}),
				},
			},
			ExpectedError: nil,
		},
		"google-dataflow": {
			ServiceDefinition: &broker.ServiceDefinition{
				Id:   "3e897eb3-9062-4966-bd4f-85bda0f73b3d",
				Name: "google-dataflow",
				BindOutputVariables: []broker.BrokerVariable{
					{
						Required:  true,
						FieldName: "Email",
						Type:      "string",
						Details:   "Email address of the service account.",
					},
					{
						Required:  true,
						FieldName: "Name",
						Type:      "string",
						Details:   "The name of the service account.",
					},
				},
				Examples: []broker.ServiceExample{
					{
						Name:            "Developer",
						Description:     "Creates a Dataflow user and grants it permission to create, drain and cancel jobs.",
						PlanId:          "8e956dd6-8c0f-470c-9a11-065537d81872",
						ProvisionParams: map[string]interface{}{},
						BindParams:      map[string]interface{}{},
					},
					{
						Name:            "Viewer",
						Description:     "Creates a Dataflow user and grants it permission to create, drain and cancel jobs.",
						PlanId:          "8e956dd6-8c0f-470c-9a11-065537d81872",
						ProvisionParams: map[string]interface{}{},
						BindParams:      map[string]interface{}{"role": "dataflow.viewer"},
					},
				},
			},
			ExpectedResponse: []CompleteServiceExample{
				{
					ServiceExample: broker.ServiceExample{
						Name:            "Developer",
						Description:     "Creates a Dataflow user and grants it permission to create, drain and cancel jobs.",
						PlanId:          "8e956dd6-8c0f-470c-9a11-065537d81872",
						ProvisionParams: map[string]interface{}{},
						BindParams:      map[string]interface{}{},
					},
					ServiceName: "google-dataflow",
					ServiceId:   "3e897eb3-9062-4966-bd4f-85bda0f73b3d",
					ExpectedOutput: broker.CreateJsonSchema([]broker.BrokerVariable{
						{
							Required:  true,
							FieldName: "Email",
							Type:      "string",
							Details:   "Email address of the service account.",
						},
						{
							Required:  true,
							FieldName: "Name",
							Type:      "string",
							Details:   "The name of the service account.",
						},
					}),
				},
				{
					ServiceExample: broker.ServiceExample{
						Name:            "Viewer",
						Description:     "Creates a Dataflow user and grants it permission to create, drain and cancel jobs.",
						PlanId:          "8e956dd6-8c0f-470c-9a11-065537d81872",
						ProvisionParams: map[string]interface{}{},
						BindParams:      map[string]interface{}{"role": "dataflow.viewer"},
					},
					ServiceName: "google-dataflow",
					ServiceId:   "3e897eb3-9062-4966-bd4f-85bda0f73b3d",
					ExpectedOutput: broker.CreateJsonSchema([]broker.BrokerVariable{
						{
							Required:  true,
							FieldName: "Email",
							Type:      "string",
							Details:   "Email address of the service account.",
						},
						{
							Required:  true,
							FieldName: "Name",
							Type:      "string",
							Details:   "The name of the service account.",
						},
					}),
				},
			},
			ExpectedError: nil,
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			actual, err := GetExamplesForAService(tc.ServiceDefinition)
			expectError(t, tc.ExpectedError, err)
			if !reflect.DeepEqual(tc.ExpectedResponse, actual) {
				t.Errorf("Expected: %v got %v", tc.ExpectedResponse, actual)
			}
		})
	}
}

func expectError(t *testing.T, expected, actual error) {
	t.Helper()
	expectedErr := expected != nil
	gotErr := actual != nil

	switch {
	case expectedErr && gotErr:
		if expected.Error() != actual.Error() {
			t.Fatalf("Expected: %v, got: %v", expected, actual)
		}
	case expectedErr && !gotErr:
		t.Fatalf("Expected: %v, got: %v", expected, actual)
	case !expectedErr && gotErr:
		t.Fatalf("Expected no error but got: %v", actual)
	}
}

func TestFilterMatchingServiceExamples(t *testing.T) {
	cases := map[string]struct {
		ServiceExamples []CompleteServiceExample
		ServiceName     string
		ExampleName     string
		Response        []CompleteServiceExample
	}{
		"No ServiceName or ExampleName specified.": {
			ServiceExamples: []CompleteServiceExample{
				{
					ServiceExample: broker.ServiceExample{
						Name: "Example-Foo",
					},
					ServiceName: "Service-Foo",
				},
				{
					ServiceExample: broker.ServiceExample{
						Name: "Example-Bar",
					},
					ServiceName: "Service-Bar",
				},
			},
			ServiceName: "",
			ExampleName: "",
			Response: []CompleteServiceExample{
				{
					ServiceExample: broker.ServiceExample{
						Name: "Example-Foo",
					},
					ServiceName: "Service-Foo",
				},
				{
					ServiceExample: broker.ServiceExample{
						Name: "Example-Bar",
					},
					ServiceName: "Service-Bar",
				},
			},
		},
		"ServiceName is specified, no ExampleName is specified. No matching CompleteServiceExample exists": {
			ServiceExamples: []CompleteServiceExample{
				{
					ServiceExample: broker.ServiceExample{
						Name: "Example-Foo",
					},
					ServiceName: "Service-Foo",
				},
				{
					ServiceExample: broker.ServiceExample{
						Name: "Example-Bar",
					},
					ServiceName: "Service-Bar",
				},
			},
			ServiceName: "Service-Unicorn",
			ExampleName: "",
			Response:    []CompleteServiceExample(nil),
		},
		"ServiceName is specified, no ExampleName is specified. Two matching CompleteServiceExample exist": {
			ServiceExamples: []CompleteServiceExample{
				{
					ServiceExample: broker.ServiceExample{
						Name: "Example-Foo",
					},
					ServiceName: "Service-Foo",
				},
				{
					ServiceExample: broker.ServiceExample{
						Name: "Example-Bar",
					},
					ServiceName: "Service-Bar",
				},
				{
					ServiceExample: broker.ServiceExample{
						Name: "Example-FooBar",
					},
					ServiceName: "Service-Bar",
				},
			},
			ServiceName: "Service-Bar",
			ExampleName: "",
			Response: []CompleteServiceExample{
				{
					ServiceExample: broker.ServiceExample{
						Name: "Example-Bar",
					},
					ServiceName: "Service-Bar",
				},
				{
					ServiceExample: broker.ServiceExample{
						Name: "Example-FooBar",
					},
					ServiceName: "Service-Bar",
				},
			},
		},
		"Both ServiceName and ExampleName are provided, but no matching CompleteServiceExample exists.": {
			ServiceExamples: []CompleteServiceExample{
				{
					ServiceExample: broker.ServiceExample{
						Name: "Example-Foo",
					},
					ServiceName: "Service-Foo",
				},
				{
					ServiceExample: broker.ServiceExample{
						Name: "Example-Bar",
					},
					ServiceName: "Service-Bar",
				},
			},
			ServiceName: "Service-Bar",
			ExampleName: "Example-Hello",
			Response:    []CompleteServiceExample(nil),
		},
		"Both ServiceName and ExampleName are provided, one matching CompleteServiceExample exists.": {
			ServiceExamples: []CompleteServiceExample{
				{
					ServiceExample: broker.ServiceExample{
						Name: "Example-Foo",
					},
					ServiceName: "Service-Foo",
				},
				{
					ServiceExample: broker.ServiceExample{
						Name: "Example-Bar",
					},
					ServiceName: "Service-Bar",
				},
			},
			ServiceName: "Service-Bar",
			ExampleName: "Example-Bar",
			Response: []CompleteServiceExample{
				{
					ServiceExample: broker.ServiceExample{
						Name: "Example-Bar",
					},
					ServiceName: "Service-Bar",
				},
			},
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			actual := FilterMatchingServiceExamples(tc.ServiceExamples, tc.ServiceName, tc.ExampleName)
			if !reflect.DeepEqual(tc.Response, actual) {
				t.Errorf("Expected: %v got %v", tc.Response, actual)
			}
		})
	}
}
