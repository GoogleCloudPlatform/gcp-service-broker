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

package tf

import (
	"context"
	"errors"
	"fmt"
	"time"

	"code.cloudfoundry.org/lager"
	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/models"
	"github.com/GoogleCloudPlatform/gcp-service-broker/db_service"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/providers/tf/wrapper"
	"github.com/GoogleCloudPlatform/gcp-service-broker/utils"
)

const (
	InProgress = "in progress"
	Succeeded  = "succeeded"
	Failed     = "failed"
)

// NewTfJobRunerFromEnv creates a new TfJobRunner with default configuration values.
func NewTfJobRunerFromEnv() (*TfJobRunner, error) {
	projectId, err := utils.GetDefaultProjectId()
	if err != nil {
		return nil, err
	}

	return NewTfJobRunnerForProject(projectId), nil
}

// Construct a new JobRunner for the given project.
func NewTfJobRunnerForProject(projectId string) *TfJobRunner {
	return &TfJobRunner{
		ProjectId:      projectId,
		ServiceAccount: utils.GetServiceAccountJson(),
	}
}

// TfJobRunner is responsible for executing terraform jobs in the background and
// providing a way to log and access the state of those background tasks.
//
// Jobs are given an ID and a workspace to operate in, and then the TfJobRunner
// is told which Terraform commands to execute on the given job.
// The TfJobRunner executes those commands in the background and keeps track of
// their state in a database table which gets updated once the task is completed.
//
// The TfJobRunner keeps track of the workspace and the Terraform state file so
// subsequent commands will operate on the same structure.
type TfJobRunner struct {
	ProjectId      string
	ServiceAccount string

	// Executor holds a custom executor that will be called when commands are run.
	Executor wrapper.TerraformExecutor
}

// StageJob stages a job to be executed. Before the workspace is saved to the
// database, the modules and inputs are validated by Terraform.
func (runner *TfJobRunner) StageJob(ctx context.Context, jobId string, workspace *wrapper.TerraformWorkspace) error {
	workspace.Executor = runner.Executor

	// Validate that TF is happy with the workspace
	if err := workspace.Validate(); err != nil {
		return err
	}

	workspaceString, err := workspace.Serialize()
	if err != nil {
		return err
	}

	deployment := &models.TerraformDeployment{
		ID:                jobId,
		Workspace:         workspaceString,
		LastOperationType: "validation",
	}

	return runner.operationFinished(nil, workspace, deployment)
}

func (runner *TfJobRunner) markJobStarted(ctx context.Context, deployment *models.TerraformDeployment, operationType string) error {
	// update the deployment info
	deployment.LastOperationType = operationType
	deployment.LastOperationState = InProgress
	deployment.LastOperationMessage = ""

	if err := db_service.SaveTerraformDeployment(ctx, deployment); err != nil {
		return err
	}

	return nil
}

func (runner *TfJobRunner) hydrateWorkspace(ctx context.Context, deployment *models.TerraformDeployment) (*wrapper.TerraformWorkspace, error) {
	ws, err := wrapper.DeserializeWorkspace(deployment.Workspace)
	if err != nil {
		return nil, err
	}

	// set environment variables
	ws.Environment = map[string]string{
		"GOOGLE_CREDENTIALS": runner.ServiceAccount,
		"GOOGLE_PROJECT":     runner.ProjectId,
	}

	ws.Executor = runner.Executor

	logger := utils.NewLogger("job-runner")
	logger.Info("wrapping", lager.Data{
		"wrapper": ws,
	})

	return ws, nil
}

// Create runs `terraform apply` on the given workspace in the background.
// The status of the job can be found by polling the Status function.
func (runner *TfJobRunner) Create(ctx context.Context, id string) error {
	deployment, err := db_service.GetTerraformDeploymentById(ctx, id)
	if err != nil {
		return err
	}

	workspace, err := runner.hydrateWorkspace(ctx, deployment)
	if err != nil {
		return err
	}

	if err := runner.markJobStarted(ctx, deployment, models.ProvisionOperationType); err != nil {
		return err
	}

	go func() {
		err := workspace.Apply()
		runner.operationFinished(err, workspace, deployment)
	}()

	return nil
}

// Destroy runs `terraform destroy` on the given workspace in the background.
// The status of the job can be found by polling the Status function.
func (runner *TfJobRunner) Destroy(ctx context.Context, id string) error {
	deployment, err := db_service.GetTerraformDeploymentById(ctx, id)
	if err != nil {
		return err
	}

	workspace, err := runner.hydrateWorkspace(ctx, deployment)
	if err != nil {
		return err
	}

	if err := runner.markJobStarted(ctx, deployment, models.DeprovisionOperationType); err != nil {
		return err
	}

	go func() {
		err := workspace.Destroy()
		runner.operationFinished(err, workspace, deployment)
	}()

	return nil
}

// operationFinished closes out the state of the background job so clients that
// are polling can get the results.
func (runner *TfJobRunner) operationFinished(err error, workspace *wrapper.TerraformWorkspace, deployment *models.TerraformDeployment) error {
	if err == nil {
		deployment.LastOperationState = Succeeded
		deployment.LastOperationMessage = ""
	} else {
		deployment.LastOperationState = Failed
		deployment.LastOperationMessage = err.Error()
	}

	workspaceString, err := workspace.Serialize()
	if err != nil {
		deployment.LastOperationState = Failed
		deployment.LastOperationMessage = fmt.Sprintf("couldn't serialize workspace, contact your operator for cleanup: %s", err.Error())
	}

	deployment.Workspace = workspaceString

	return db_service.SaveTerraformDeployment(context.Background(), deployment)
}

// Status gets the status of the most recent job on the workspace.
// If isDone is true, then the status of the operation will not change again.
// if isDone is false, then the operation is ongoing.
func (runner *TfJobRunner) Status(ctx context.Context, id string) (isDone bool, err error) {
	deployment, err := db_service.GetTerraformDeploymentById(ctx, id)
	if err != nil {
		return true, err
	}

	switch deployment.LastOperationState {
	case Succeeded:
		return true, nil
	case Failed:
		return true, errors.New(deployment.LastOperationMessage)
	default:
		return false, nil
	}
}

// Outputs gets the output variables for the given module instance in the workspace.
func (runner *TfJobRunner) Outputs(ctx context.Context, id, instanceName string) (map[string]interface{}, error) {
	deployment, err := db_service.GetTerraformDeploymentById(ctx, id)
	if err != nil {
		return nil, err
	}

	ws, err := wrapper.DeserializeWorkspace(deployment.Workspace)
	if err != nil {
		return nil, err
	}

	return ws.Outputs(instanceName)
}

// Wait waits for an operation to complete, polling its status once per second.
func (runner *TfJobRunner) Wait(ctx context.Context, id string) error {
	for {
		select {
		case <-ctx.Done():
			return nil

		case <-time.After(1 * time.Second):
			isDone, err := runner.Status(ctx, id)
			if isDone {
				return err
			}
		}
	}
}
