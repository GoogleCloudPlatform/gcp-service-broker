package tf

import (
	"context"
	"errors"
	"fmt"
	"time"

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

func NewTfJobRunerFromEnv() (*TfJobRunner, error) {
	projectId, err := utils.GetDefaultProjectId()
	if err != nil {
		return nil, err
	}

	return &TfJobRunner{
		ProjectId:      projectId,
		ServiceAccount: utils.GetServiceAccountJson(),
	}, nil
}

type TfJobRunner struct {
	ProjectId      string
	ServiceAccount string
}

func (runner *TfJobRunner) StageJob(ctx context.Context, jobId string, workspace *wrapper.TerraformWorkspace) error {
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
		LastOperationType: models.ProvisionOperationType,
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

	return ws, nil
}

func (runner *TfJobRunner) Dump(ctx context.Context, id string) (*wrapper.TerraformWorkspace, error) {
	deployment, err := db_service.GetTerraformDeploymentById(ctx, id)
	if err != nil {
		return nil, err
	}

	return runner.hydrateWorkspace(ctx, deployment)
}

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

func (runner *TfJobRunner) Status(ctx context.Context, id string) (bool, error) {
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
