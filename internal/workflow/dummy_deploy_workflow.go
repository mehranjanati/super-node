package workflow

import (
	"fmt"
	"time"

	"go.temporal.io/sdk/workflow"
)

const DummyDeploymentStatusQuery = "deployment_status"

type DummyDeploymentInput struct {
	ProjectName string `json:"project_name"`
	Template    string `json:"template"`
}

type DummyDeploymentArtifacts struct {
	ProjectName string `json:"project_name"`
	Template    string `json:"template"`
	URL         string `json:"url"`
	Message     string `json:"message"`
}

type DummyDeploymentStatus struct {
	WorkflowID  string                    `json:"workflow_id,omitempty"`
	Status      string                    `json:"status"`
	CurrentStep string                    `json:"current_step"`
	Logs        []string                  `json:"logs"`
	Artifacts   *DummyDeploymentArtifacts `json:"artifacts,omitempty"`
	Error       string                    `json:"error,omitempty"`
}

// DummyDeploymentWorkflow is a short deterministic workflow for day 5.
func DummyDeploymentWorkflow(ctx workflow.Context, input DummyDeploymentInput) (DummyDeploymentStatus, error) {
	logger := workflow.GetLogger(ctx)
	state := DummyDeploymentStatus{
		Status:      "RUNNING",
		CurrentStep: "INIT",
		Logs: []string{
			fmt.Sprintf("Accepted deployment request for %s.", input.ProjectName),
			fmt.Sprintf("Using template %s.", input.Template),
		},
	}

	if err := workflow.SetQueryHandler(ctx, DummyDeploymentStatusQuery, func() (DummyDeploymentStatus, error) {
		return state, nil
	}); err != nil {
		return DummyDeploymentStatus{}, err
	}

	steps := []struct {
		step string
		log  string
		wait time.Duration
	}{
		{step: "GEN_SCHEMA", log: "Generating deployment schema.", wait: time.Second},
		{step: "GEN_CODE", log: "Preparing source bundle.", wait: time.Second},
		{step: "BUILDING", log: "Building deployment artifact.", wait: time.Second},
		{step: "DEPLOYING", log: "Publishing preview environment.", wait: time.Second},
	}

	for _, item := range steps {
		if err := workflow.Sleep(ctx, item.wait); err != nil {
			state.Status = "FAILED"
			state.Error = err.Error()
			state.Logs = append(state.Logs, "Workflow interrupted before completion.")
			return state, err
		}

		state.CurrentStep = item.step
		state.Logs = append(state.Logs, item.log)
		logger.Info("Dummy deployment step completed", "step", item.step, "project", input.ProjectName)
	}

	state.Status = "COMPLETED"
	state.CurrentStep = "DONE"
	state.Logs = append(state.Logs, "Deployment finished successfully.")
	state.Artifacts = &DummyDeploymentArtifacts{
		ProjectName: input.ProjectName,
		Template:    input.Template,
		URL:         fmt.Sprintf("https://%s.nexus.local", input.ProjectName),
		Message:     "Dummy deployment completed via Temporal workflow",
	}

	return state, nil
}
