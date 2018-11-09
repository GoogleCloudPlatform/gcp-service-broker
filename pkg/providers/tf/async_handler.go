package tf

import (
	"context"
	"fmt"

	"github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/models"
	"github.com/GoogleCloudPlatform/gcp-service-broker/db_service"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/providers/tf/wrapper"
	"github.com/pivotal-cf/brokerapi"
)

func Apply(id string, workspace *wrapper.TerraformWorkspace) (string, error) {
	if err := workspace.InitializeFs(); err != nil {
		return "", fmt.Errorf("couldn't initialize temporary workspace: %s", err.Error())
	}

	// Validate that TF is happy
	if err := workspace.Validate(); err != nil {
		return "", err
	}

	workspaceString, err := workspace.Serialize()
	if err != nil {
		return "", err
	}

	deployment := models.TerraformDeployment{
		ID:                   id,
		Workspace:            workspaceString,
		LastOperationType:    models.ProvisionOperationType,
		LastOperationState:   string(brokerapi.InProgress),
		LastOperationMessage: "Terraform Job Starting",
	}

	// create DB instance
	// go start job
	// return ID

	return deployment.ID, nil
}

func Destroy() {

}

func Cleanup() {
	// update DB with result status
}

func Status(ctx context.Context, id string) (*brokerapi.LastOperation, error) {
	deployment, err := db_service.GetTerraformDeploymentById(ctx, id)
	if err != nil {
		return nil, err
	}

	return &brokerapi.LastOperation{
		State:       brokerapi.LastOperationState(deployment.LastOperationState),
		Description: deployment.LastOperationMessage,
	}, nil
}

func Outputs(ctx context.Context, id, instanceName string) (map[string]interface{}, error) {
	// poll the TF DB, extract the outputs, match with
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
