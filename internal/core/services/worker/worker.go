package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"nexus-super-node-v3/internal/core/services/mcp"
	"nexus-super-node-v3/internal/core/services/wazero"
	"nexus-super-node-v3/internal/ports"
)

type WorkerService struct {
	mcpSvc      *mcp.MCPService
	wazeroSvc   *wazero.WazeroService
	rivetEngine ports.RivetEngine
}

func NewWorkerService(mcpSvc *mcp.MCPService, wazeroSvc *wazero.WazeroService, rivetEngine ports.RivetEngine) *WorkerService {
	return &WorkerService{
		mcpSvc:      mcpSvc,
		wazeroSvc:   wazeroSvc,
		rivetEngine: rivetEngine,
	}
}

// HasuraPayload represents the structure of a Hasura Action payload
type HasuraPayload struct {
	Action struct {
		Name string `json:"name"`
	} `json:"action"`
	Input json.RawMessage `json:"input"`
}

// WorkflowTask represents a task to run a Rivet workflow
type WorkflowTask struct {
	Type    string                 `json:"type"`
	GraphID string                 `json:"graph_id"`
	Inputs  map[string]interface{} `json:"inputs"`
}

// ProcessTask handles the incoming task from Redpanda
func (w *WorkerService) ProcessTask(ctx context.Context, payload []byte) error {
	// 1. Try to parse as Workflow Task
	var workflowTask WorkflowTask
	if err := json.Unmarshal(payload, &workflowTask); err == nil && workflowTask.Type == "workflow" {
		return w.handleWorkflowTask(ctx, workflowTask)
	}

	// 2. Try to parse as Hasura Payload
	var hasuraEvent HasuraPayload
	if err := json.Unmarshal(payload, &hasuraEvent); err == nil && hasuraEvent.Action.Name != "" {
		return w.handleHasuraAction(ctx, hasuraEvent)
	}

	// 3. Fallback to generic task (to be defined)
	log.Printf("Received generic task: %s", string(payload))
	return nil
}

func (w *WorkerService) handleWorkflowTask(ctx context.Context, task WorkflowTask) error {
	log.Printf("Processing Workflow Task: %s", task.GraphID)

	outputs, err := w.rivetEngine.ExecuteGraph(ctx, task.GraphID, task.Inputs)
	if err != nil {
		log.Printf("Error executing workflow %s: %v", task.GraphID, err)
		return err
	}

	log.Printf("Workflow %s executed successfully. Outputs: %v", task.GraphID, outputs)
	return nil
}

func (w *WorkerService) handleHasuraAction(ctx context.Context, event HasuraPayload) error {
	actionName := event.Action.Name
	log.Printf("Processing Hasura Action: %s", actionName)

	// Convention: Map action name to MCP tool
	// We'll look for a tool with the same name in the "default" or "sample" toolbelt for now
	// In a real app, you might want a more sophisticated mapping or look up via a registry

	// For demonstration, let's assume the action name maps to a tool in the "sample" toolbelt
	// If the action is "hello", we call "hello" tool in "sample" belt

	// Parse input to pass as args
	var inputMap map[string]interface{}
	if err := json.Unmarshal(event.Input, &inputMap); err != nil {
		return fmt.Errorf("failed to parse action input: %w", err)
	}

	// Convert input map values to args slice
	// This is a simplification; real MCP tools might expect structured input
	var args []interface{}
	for _, v := range inputMap {
		args = append(args, v)
	}

	result, err := w.mcpSvc.ExecuteTool("sample", actionName, args...)
	if err != nil {
		log.Printf("Error executing tool %s: %v", actionName, err)
		// We might want to push this error back to a "dead letter queue" or Hasura
		return err
	}

	log.Printf("Action %s executed successfully. Result: %v", actionName, result)

	// TODO: If this is a synchronous action, we might need to write the result back to Hasura or a response topic
	return nil
}
