package domain

import "time"

type WorkflowExecutionArtifacts struct {
	ProjectName    string `json:"project_name,omitempty"`
	Template       string `json:"template,omitempty"`
	PlanningSource string `json:"planning_source,omitempty"`
	URL            string `json:"url,omitempty"`
	Message        string `json:"message,omitempty"`
}

type WorkflowExecutionRecord struct {
	WorkflowID  string                      `json:"workflow_id"`
	Kind        string                      `json:"kind"`
	Name        string                      `json:"name"`
	Status      string                      `json:"status"`
	CurrentStep string                      `json:"current_step"`
	Logs        []string                    `json:"logs"`
	Artifacts   *WorkflowExecutionArtifacts `json:"artifacts,omitempty"`
	Error       string                      `json:"error,omitempty"`
	CreatedAt   time.Time                   `json:"created_at"`
	UpdatedAt   time.Time                   `json:"updated_at"`
}
