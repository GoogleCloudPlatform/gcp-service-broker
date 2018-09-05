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
	"strconv"

	googlespanner "cloud.google.com/go/spanner/admin/instance/apiv1"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/broker_base"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/models"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/name_generator"
	"github.com/GoogleCloudPlatform/gcp-service-broker/db_service"
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
func (s *SpannerBroker) Provision(instanceId string, details brokerapi.ProvisionDetails, plan models.ServicePlan) (models.ServiceInstanceDetails, error) {
	var err error
	var params map[string]string

	if len(details.RawParameters) == 0 {
		params = map[string]string{}
	} else if err = json.Unmarshal(details.RawParameters, &params); err != nil {
		return models.ServiceInstanceDetails{}, fmt.Errorf("Error unmarshalling parameters: %s", err)
	}

	// Ensure there is a name for this instance
	if _, ok := params["name"]; !ok {
		params["name"] = name_generator.Basic.InstanceNameWithSeparator("-")
	}

	// set up client

	co := option.WithUserAgent(models.CustomUserAgent)
	ct := option.WithTokenSource(s.HttpConfig.TokenSource(context.Background()))
	client, err := googlespanner.NewInstanceAdminClient(context.Background(), co, ct)
	if err != nil {
		return models.ServiceInstanceDetails{}, fmt.Errorf("Error creating client: %s", err)
	}

	// set up params
	numNodes, err := strconv.Atoi(plan.ServiceProperties["num_nodes"])
	if err != nil {
		return models.ServiceInstanceDetails{}, fmt.Errorf("Error getting number of nodes: %s", err)
	}

	displayName := params["name"]
	if customerDisplayName, ok := params["display_name"]; ok {
		displayName = customerDisplayName
	}

	loc, ok := params["location"]
	if !ok {
		loc = "projects/" + s.ProjectId + "/instanceConfigs/regional-us-central1"
	} else {
		loc = "projects/" + s.ProjectId + "/instanceConfigs/" + loc
	}

	// create instance
	op, err := client.CreateInstance(context.Background(), &instancepb.CreateInstanceRequest{
		Parent:     "projects/" + s.ProjectId,
		InstanceId: params["name"],
		Instance: &instancepb.Instance{
			Name:        "projects/" + s.ProjectId + "/instances/" + params["name"],
			DisplayName: displayName,
			NodeCount:   int32(numNodes),
			Config:      loc,
		},
	})
	if err != nil {
		return models.ServiceInstanceDetails{}, fmt.Errorf("Error creating instance: %s", err)
	}

	// save off instance information
	ii := InstanceInformation{
		InstanceId: params["name"],
	}

	otherDetails, err := json.Marshal(ii)
	if err != nil {
		return models.ServiceInstanceDetails{}, fmt.Errorf("Error marshalling other details: %s", err)
	}

	i := models.ServiceInstanceDetails{
		Name:         params["name"],
		Url:          "",
		Location:     loc,
		OtherDetails: string(otherDetails),
	}

	err = createCloudOperation(op, instanceId, details.ServiceID)
	if err != nil {
		return models.ServiceInstanceDetails{}, fmt.Errorf("Error saving operation to database: %s", err)
	}

	return i, nil
}

// PollInstance gets the last operation for this instance and polls its status.
func (s *SpannerBroker) PollInstance(instanceId string) (bool, error) {

	op, err := db_service.GetCloudOperationByServiceInstanceId(instanceId)
	if err != nil {
		return false, fmt.Errorf("Could not locate CloudOperation in database: %s", err)
	}

	if _, err := db_service.GetServiceInstanceDetailsById(instanceId); err != nil {
		return false, brokerapi.ErrInstanceDoesNotExist
	}

	// we're polling on instance deletion, which is synchronous, unlike creation. Exit early if the instance has been deleted
	wasDelete, err := s.LastOperationWasDelete(instanceId)

	if err != nil {
		return false, fmt.Errorf("Can't check last operation type: %s", err)
	}
	if wasDelete {
		return true, nil
	}

	ct := option.WithTokenSource(s.HttpConfig.TokenSource(context.Background()))
	client, err := googlespanner.NewInstanceAdminClient(context.Background(), ct)
	if err != nil {
		return false, fmt.Errorf("Error creating client: %s", err)
	}

	spannerOp := client.CreateInstanceOperation(op.Name)

	spannerInstance, err := spannerOp.Poll(context.Background())
	done := spannerOp.Done()

	// from https://godoc.org/cloud.google.com/go/spanner/admin/instance/apiv1#InstanceOperation.Poll
	if spannerInstance == nil && err != nil && !done {
		return false, fmt.Errorf("Error checking operation status: %s", err)
	} else if spannerInstance == nil && err != nil && done {
		op.Status = "FAILED"
		op.ErrorMessage = err.Error()

		if dberr := db_service.SaveCloudOperation(op); dberr != nil {
			return false, fmt.Errorf(`Error saving operation details to database: %s.`, dberr)
		}

		return true, fmt.Errorf("Error provisioning instance: %v", err)
	} else if spannerInstance == nil && err == nil && !done {
		op.Status = string(instancepb.Instance_STATE_UNSPECIFIED)

		if err := db_service.SaveCloudOperation(op); err != nil {
			return false, fmt.Errorf(`Error saving operation details to database: %s.`, err)
		}

		return false, nil
	} else if spannerInstance != nil && err == nil && done {
		op.Status = spannerInstance.State.String()

		if err := db_service.SaveCloudOperation(op); err != nil {
			return false, fmt.Errorf(`Error saving operation details to database: %s.`, err)
		}

		return true, nil
	}

	return false, fmt.Errorf("unknown error")
}

func createCloudOperation(op *googlespanner.CreateInstanceOperation, instanceId string, serviceId string) error {
	errorStr := ""
	if _, err := op.Poll(context.Background()); err != nil {
		errorStr = err.Error()
	}

	metadata, err := op.Metadata()
	if err != nil {
		return fmt.Errorf("Error getting operation metadata: %s", err)
	}

	startTime := ""
	if metadata.StartTime != nil {
		startTime = metadata.StartTime.String()
	}

	currentState := models.CloudOperation{
		Name:              op.Name(),
		ErrorMessage:      errorStr,
		InsertTime:        startTime,
		OperationType:     "SPANNER_OPERATION",
		StartTime:         startTime,
		Status:            metadata.Instance.State.String(),
		ServiceId:         serviceId,
		ServiceInstanceId: instanceId,
	}

	if err = db_service.CreateCloudOperation(&currentState); err != nil {
		return fmt.Errorf("Error saving operation details to database: %s. Services relying on async deprovisioning will not be able to complete deprovisioning", err)
	}
	return nil
}

// Deprovision deletes the Spanner instance associated with the given instance.
func (s *SpannerBroker) Deprovision(ctx context.Context, instance models.ServiceInstanceDetails, details brokerapi.DeprovisionDetails) error {
	// set up client
	co := option.WithUserAgent(models.CustomUserAgent)
	ct := option.WithTokenSource(s.HttpConfig.TokenSource(ctx))
	client, err := googlespanner.NewInstanceAdminClient(ctx, co, ct)
	if err != nil {
		return fmt.Errorf("Error creating client: %s", err)
	}

	// delete instance
	err = client.DeleteInstance(context.Background(), &instancepb.DeleteInstanceRequest{
		Name: "projects/" + s.ProjectId + "/instances/" + instance.Name,
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

// LastOperationWasDelete is used during polling of async operations to check
// if the workflow is a provision or deprovision flow based off the type of the
// most recent operation
// since spanner deprovisions synchronously, the last operation will never have
// been delete
func (s *SpannerBroker) LastOperationWasDelete(instanceId string) (bool, error) {
	return false, nil
}
