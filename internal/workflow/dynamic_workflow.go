package workflow

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"nexus-super-node-v3/internal/core/domain"

	"go.temporal.io/sdk/workflow"
)

// DynamicPipelineWorkflow is a generic workflow that executes a pipeline based on its definition.
func DynamicPipelineWorkflow(ctx workflow.Context, pipelineDef domain.PipelineDefinition, inputs map[string]interface{}) (map[string]interface{}, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("Starting Dynamic Pipeline", "PipelineID", pipelineDef.ID, "Name", pipelineDef.Name)

	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute * 10, // Default timeout, could be configurable per step
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	// Execution Context (stores inputs and step results)
	execContext := make(map[string]interface{})
	for k, v := range inputs {
		execContext[k] = v
	}

	// Iterate through steps
	for _, step := range pipelineDef.Steps {
		logger.Info("Executing Step", "StepID", step.ID, "Activity", step.ActivityName)

		// 1. Resolve Arguments
		resolvedArgs, err := resolveArgs(step.Args, execContext)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve args for step %s: %w", step.ID, err)
		}

		// 2. Configure Activity Options for this step
		stepAO := ao
		if step.Timeout != "" {
			timeout, err := time.ParseDuration(step.Timeout)
			if err == nil {
				stepAO.StartToCloseTimeout = timeout
			}
		}
		if step.TaskQueue != "" {
			stepAO.TaskQueue = step.TaskQueue
		}
		stepCtx := workflow.WithActivityOptions(ctx, stepAO)

		// 3. Execute Activity (if provided)
		if step.ActivityName != "" {
			var result interface{}
			err = workflow.ExecuteActivity(stepCtx, step.ActivityName, resolvedArgs...).Get(stepCtx, &result)
			if err != nil {
				logger.Error("Step Failed", "StepID", step.ID, "Error", err)
				return nil, err
			}

			// Store Result
			if step.ResultKey != "" {
				execContext[step.ResultKey] = result
			}
		}

		// 4. Wait for Signal (if provided)
		if step.WaitSignal != "" {
			logger.Info("Pausing for signal", "Signal", step.WaitSignal)
			var signalData interface{}
			signalChan := workflow.GetSignalChannel(ctx, step.WaitSignal)

			// Wait for signal
			signalChan.Receive(ctx, &signalData)
			logger.Info("Received signal", "Signal", step.WaitSignal, "Data", signalData)

			// Store signal data if result key is provided
			if step.ResultKey != "" {
				execContext[step.ResultKey] = signalData
			}
		}
	}

	logger.Info("Pipeline Completed Successfully", "PipelineID", pipelineDef.ID)
	return execContext, nil
}

// resolveArgs replaces placeholders like "{{.project_name}}" with actual values from context.
func resolveArgs(args map[string]interface{}, context map[string]interface{}) ([]interface{}, error) {
	resolved := make(map[string]interface{})
	for k, v := range args {
		strVal, ok := v.(string)
		if ok && strings.Contains(strVal, "{{") && strings.Contains(strVal, "}}") {
			// Handle simple exact match "{{key}}"
			if strings.HasPrefix(strVal, "{{") && strings.HasSuffix(strVal, "}}") {
				key := strings.TrimSuffix(strings.TrimPrefix(strVal, "{{"), "}}")
				key = strings.TrimSpace(key)
				if val, found := lookupValue(context, key); found {
					resolved[k] = val
				} else {
					resolved[k] = v // Keep original if not found
				}
			} else {
				// Handle string interpolation "Hello {{name}}"
				// Simple implementation: replace only exact {{key}} occurrences
				// A regex based approach would be better for multiple replacements
				// For now, let's keep it simple or use strings.Replace
				// TODO: Implement full template engine support
				resolved[k] = v
			}
		} else {
			resolved[k] = v
		}
	}

	// Return as a single argument (the map)
	return []interface{}{resolved}, nil
}

// lookupValue resolves dot notation keys like "analysis.TopPick"
func lookupValue(data map[string]interface{}, path string) (interface{}, bool) {
	parts := strings.Split(path, ".")
	if len(parts) == 0 {
		return nil, false
	}

	current := interface{}(data)
	for _, part := range parts {
		// If current is nil, stop
		if current == nil {
			return nil, false
		}

		// Check if current is a map
		if m, ok := current.(map[string]interface{}); ok {
			val, exists := m[part]
			if !exists {
				return nil, false
			}
			current = val
			continue
		}

		// Check if current is a struct (convert to map via JSON for simplicity)
		// This is expensive but handles arbitrary structs without complex reflection
		jsonBytes, err := json.Marshal(current)
		if err != nil {
			return nil, false
		}
		var m map[string]interface{}
		if err := json.Unmarshal(jsonBytes, &m); err == nil {
			if val, exists := m[part]; exists {
				current = val
				continue
			}
		}

		// If neither map nor convertible struct, fail
		return nil, false
	}

	return current, true
}
