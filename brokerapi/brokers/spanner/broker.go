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

package spanner

import (
	"context"
	"encoding/json"
	"fmt"

	googlespanner "cloud.google.com/go/spanner/admin/instance/apiv1"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/broker_base"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/models"
	"github.com/GoogleCloudPlatform/gcp-service-broker/db_service"
	"github.com/GoogleCloudPlatform/gcp-service-broker/utils"
	"github.com/pivotal-cf/brokerapi"
	"google.golang.org/api/option"
	instancepb "google.golang.org/genproto/googleapis/spanner/admin/instance/v1"
)

// SpannerBroker is the service-broker back-end for creating Spanner databases
// and accounts.
type SpannerBroker struct {
	broker_base.BrokerBase
}

// InstanceInformation holds the details needed to connect to a Spanner instance
// after it has been provisioned.
type InstanceInformation struct {
	InstanceId string `json:"instance_id"`
}

// Provision creates a new Spanner instance from the settings in the user-provided details and service plan.
func (s *SpannerBroker) Provision(ctx context.Context, instanceId string, details brokerapi.ProvisionDetails, plan models.ServicePlan) (models.ServiceInstanceDetails, error) {
	variableContext, err := serviceDefinition().ProvisionVariables(instanceId, details, plan)
	if err != nil {
		return models.ServiceInstanceDetails{}, err
	}

	// create instance provision request
	instanceName := variableContext.GetString("name")
	instanceLocation := fmt.Sprintf("projects/%s/instanceConfigs/%s", s.ProjectId, variableContext.GetString("location"))

	creationRequest := instancepb.CreateInstanceRequest{
		Parent:     "projects/" + s.ProjectId,
		InstanceId: instanceName,
		Instance: &instancepb.Instance{
			Name:        s.qualifiedInstanceName(instanceName),
			DisplayName: variableContext.GetString("display_name"),
			NodeCount:   int32(variableContext.GetInt("num_nodes")),
			Config:      instanceLocation,
			Labels:      utils.ExtractDefaultLabels(instanceId, details),
		},
	}

	if err := variableContext.Error(); err != nil {
		return models.ServiceInstanceDetails{}, err
	}

	// Make request
	client, err := s.createAdminClient(ctx)
	if err != nil {
		return models.ServiceInstanceDetails{}, err
	}
	op, err := client.CreateInstance(ctx, &creationRequest)
	if err != nil {
		return models.ServiceInstanceDetails{}, fmt.Errorf("Error creating instance: %s", err)
	}

	// save off instance information
	ii := InstanceInformation{
		InstanceId: instanceName,
	}

	otherDetails, err := json.Marshal(ii)
	if err != nil {
		return models.ServiceInstanceDetails{}, fmt.Errorf("Error marshalling other details: %s", err)
	}

	return models.ServiceInstanceDetails{
		Name:          instanceName,
		Url:           "",
		Location:      instanceLocation,
		OtherDetails:  string(otherDetails),
		OperationType: models.ProvisionOperationType,
		OperationId:   op.Name(),
	}, nil
}

// PollInstance gets the last operation for this instance and polls its status.
func (s *SpannerBroker) PollInstance(ctx context.Context, instanceId string) (bool, error) {
	instance, err := db_service.GetServiceInstanceDetailsById(ctx, instanceId)
	if err != nil {
		return false, brokerapi.ErrInstanceDoesNotExist
	}

	if instance.OperationType == models.ClearOperationType {
		return false, fmt.Errorf("No pending operations could be found for this Spanner instance.")
	}

	if instance.OperationType != models.ProvisionOperationType {
		return false, fmt.Errorf("Couldn't poll Spanner instance, unknown operation type: %s", instance.OperationType)
	}

	client, err := s.createAdminClient(ctx)
	if err != nil {
		return false, err
	}

	// From https://godoc.org/cloud.google.com/go/spanner/admin/instance/apiv1#CreateInstanceOperation.Poll
	spannerOp := client.CreateInstanceOperation(instance.OperationId)
	_, err = spannerOp.Poll(ctx)
	done := spannerOp.Done()

	switch {
	case err != nil && !done: // There was a failure polling
		return false, fmt.Errorf("Error checking operation status: %s", err)

	case err != nil && done: // The operation completed in error
		return true, fmt.Errorf("Error provisioning instance: %v", err)

	case err == nil && done: // The operation was successful
		instance.OperationId = ""
		instance.OperationType = models.ClearOperationType
		if err := db_service.SaveServiceInstanceDetails(ctx, instance); err != nil {
			s.Logger.Error("updating instance after provision", err)
		}
		return true, nil

	default: // The operation hasn't completed yet
		return false, nil
	}
}

// Deprovision deletes the Spanner instance associated with the given instance.
func (s *SpannerBroker) Deprovision(ctx context.Context, instance models.ServiceInstanceDetails, details brokerapi.DeprovisionDetails) error {
	client, err := s.createAdminClient(ctx)
	if err != nil {
		return err
	}

	// delete instance
	err = client.DeleteInstance(ctx, &instancepb.DeleteInstanceRequest{
		Name: s.qualifiedInstanceName(instance.Name),
	})

	if err != nil {
		return fmt.Errorf("Error deleting instance: %s", err)
	}

	return nil
}

// ProvisionsAsync indicates that Spanner uses asynchronous provisioning
func (s *SpannerBroker) ProvisionsAsync() bool {
	return true
}

// qualifiedInstanceName gets the fully qualified instance name with
// regards to the project id.
func (s *SpannerBroker) qualifiedInstanceName(instanceName string) string {
	return fmt.Sprintf("projects/%s/instances/%s", s.ProjectId, instanceName)
}

// LastOperationWasDelete is used during polling of async operations to check
// if the workflow is a provision or deprovision flow based off the type of the
// most recent operation
// since spanner deprovisions synchronously, the last operation will never have
// been delete
func (s *SpannerBroker) LastOperationWasDelete(ctx context.Context, instanceId string) (bool, error) {
	return false, nil
}

func (s *SpannerBroker) createAdminClient(ctx context.Context) (*googlespanner.InstanceAdminClient, error) {
	co := option.WithUserAgent(models.CustomUserAgent)
	ct := option.WithTokenSource(s.HttpConfig.TokenSource(ctx))
	client, err := googlespanner.NewInstanceAdminClient(ctx, co, ct)
	if err != nil {
		return nil, fmt.Errorf("Couldn't instantiate Spanner API client: %s", err)
	}

	return client, nil
}
