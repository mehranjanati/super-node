package workflow

import (
	"fmt"
	"nexus-super-node-v3/internal/core/domain"
	"strings"
	"time"

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
// This is a simplified implementation. Real-world would use text/template.
func resolveArgs(args map[string]interface{}, context map[string]interface{}) ([]interface{}, error) {
	// In Temporal ExecuteActivity, args are passed as variadic interface{}.
	// However, our activities usually take specific structs (e.g. WebsiteDeploymentParams).
	// This generic workflow implies that Activities need to accept a map or we need strict mapping.
	//
	// Strategy A: All dynamic activities take (ctx, map[string]interface{}) -> (map[string]interface{}, error)
	// Strategy B: We map args to a struct (complex reflection).
	//
	// For this implementation, we will assume Strategy A for truly dynamic activities,
	// OR we assume 'args' contains the exact single argument struct if the activity expects one.

	// Let's assume for now that we pass the whole 'args' map as the first argument,
	// and let the activity handle parsing. This requires wrapper activities.

	resolved := make(map[string]interface{})
	for k, v := range args {
		strVal, ok := v.(string)
		if ok && strings.HasPrefix(strVal, "{{") && strings.HasSuffix(strVal, "}}") {
			key := strings.TrimSuffix(strings.TrimPrefix(strVal, "{{"), "}}")
			key = strings.TrimSpace(key)
			if val, exists := context[key]; exists {
				resolved[k] = val
			} else {
				// Keep original if missing, or error out?
				resolved[k] = v
			}
		} else {
			resolved[k] = v
		}
	}

	// Return as a single argument (the map)
	return []interface{}{resolved}, nil
}
