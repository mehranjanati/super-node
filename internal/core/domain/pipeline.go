package domain

// PipelineDefinition defines the structure of a dynamic pipeline.
type PipelineDefinition struct {
	ID          string         `json:"id" yaml:"id"`
	Name        string         `json:"name" yaml:"name"`
	Description string         `json:"description" yaml:"description"`
	Version     string         `json:"version" yaml:"version"`
	Inputs      []string       `json:"inputs" yaml:"inputs"` // Required input keys
	Steps       []PipelineStep `json:"steps" yaml:"steps"`
}

// PipelineStep defines a single step in the pipeline.
type PipelineStep struct {
	ID           string                 `json:"id" yaml:"id"`
	ActivityName string                 `json:"activity_name" yaml:"activity_name"` // e.g., "GenerateSourceCode"
	Args         map[string]interface{} `json:"args" yaml:"args"`                   // Static args or templates like "{{.Inputs.project_name}}"
	ResultKey    string                 `json:"result_key" yaml:"result_key"`       // Where to store the result in the context
	WaitSignal   string                 `json:"wait_signal" yaml:"wait_signal"`     // Optional: If set, workflow pauses until this signal is received
	Timeout      string                 `json:"timeout" yaml:"timeout"`             // Optional: Timeout for this step (e.g., "1h")
	TaskQueue    string                 `json:"task_queue" yaml:"task_queue"`       // Optional: Route this activity to a specific task queue (External Worker)
}

// PipelineExecutionRequest represents a request to execute a pipeline.
type PipelineExecutionRequest struct {
	PipelineID string                 `json:"pipeline_id"`
	Inputs     map[string]interface{} `json:"inputs"`
}
