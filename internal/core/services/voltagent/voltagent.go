package voltagent

import (
	"context"
	"fmt"
	"strings"

	"nexus-super-node-v3/internal/adapters/ai"
	"nexus-super-node-v3/internal/core/domain"
	"nexus-super-node-v3/internal/workflow"

	"go.temporal.io/sdk/client"
)

// StreamChat forwards chat messages to the AI adapter via VoltAgent
func (s *VoltAgentService) StreamChat(ctx context.Context, messages []ai.ChatMessage) (<-chan string, <-chan error) {
	if s.aiClient == nil {
		errChan := make(chan error, 1)
		errChan <- fmt.Errorf("AI client not configured")
		close(errChan)
		return nil, errChan
	}
	// In the future, VoltAgent might inspect messages, inject system prompts based on context,
	// or decide to call tools directly before responding.
	// For now, we pass through to the AI adapter.
	return s.aiClient.StreamChat(ctx, messages)
}

// GetManifest returns a manifest of all available tools for VoltAgent
func (s *VoltAgentService) GetManifest() (*VoltAgentManifest, error) {
	manifest := &VoltAgentManifest{
		Version: "1.0.0",
		Tools:   []VoltAgentTool{},
	}

	// 1. Add local tools
	for _, belt := range s.mcpSvc.ToolBelts {
		for _, tool := range belt.Tools {
			manifest.Tools = append(manifest.Tools, VoltAgentTool{
				Name:        fmt.Sprintf("%s__%s", belt.Name, tool.Name),
				Description: tool.Description,
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"args": map[string]interface{}{
							"type": "array",
							"items": map[string]interface{}{
								"type": "string",
							},
						},
					},
				},
			})
		}
	}

	// 2. Add dynamic MCP tools
	servers := s.mcpSvc.ListDynamicServers()
	for _, server := range servers {
		manifest.Tools = append(manifest.Tools, VoltAgentTool{
			Name:        fmt.Sprintf("mcp__%s__execute", server.ID),
			Description: fmt.Sprintf("Execute tools on the %s MCP server", server.Name),
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"tool_name": map[string]string{"type": "string"},
					"arguments": map[string]interface{}{"type": "object"},
				},
			},
		})
	}

	// 3. Add Website Deployment Tool (System Tool)
	manifest.Tools = append(manifest.Tools, VoltAgentTool{
		Name:        "system__deploy_website",
		Description: "Generate and deploy a website based on a description/prompt using AI and Temporal workflows.",
		Parameters: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"project_name": map[string]string{"type": "string"},
				"prompt":       map[string]string{"type": "string"},
				"theme":        map[string]string{"type": "string", "enum": "light,dark,modern,minimal"},
				"framework":    map[string]string{"type": "string", "enum": "svelte,react,vue"},
			},
			"required": []string{"project_name", "prompt"},
		},
	})

	// 4. Add Crypto Analysis Tool (System Tool)
	manifest.Tools = append(manifest.Tools, VoltAgentTool{
		Name:        "system__crypto_analysis",
		Description: "Perform deep market analysis and optionally execute trades with human approval.",
		Parameters: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"user_id":    map[string]string{"type": "string"},
				"time_frame": map[string]string{"type": "string", "enum": "daily,weekly,3week"},
			},
			"required": []string{"user_id"},
		},
	})

	// 5. Add Human Handoff Tool (System Tool)
	manifest.Tools = append(manifest.Tools, VoltAgentTool{
		Name:        "system__human_handoff",
		Description: "Initiate human handoff protocol when AI cannot resolve a request or a voice call is requested.",
		Parameters: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"room_id":        map[string]string{"type": "string"},
				"target_user_id": map[string]string{"type": "string"},
			},
			"required": []string{"room_id", "target_user_id"},
		},
	})

	return manifest, nil
}

// ExecuteTool forwards a request from VoltAgent to the MCP service or triggers internal workflows
func (s *VoltAgentService) ExecuteTool(toolID string, args map[string]interface{}) (interface{}, error) {
	if toolID == "system__deploy_website" {
		// ... existing logic ...
		projectName, _ := args["project_name"].(string)
		prompt, _ := args["prompt"].(string)
		theme, _ := args["theme"].(string)
		framework, _ := args["framework"].(string)

		pipelineDef := domain.PipelineDefinition{
			ID:      "website-deployment-v1",
			Name:    "Standard Website Deployment",
			Version: "1.0.0",
			Inputs:  []string{"project_name", "prompt", "theme", "framework"},
			Steps: []domain.PipelineStep{
				{
					ID:           "gen_ui",
					ActivityName: "GenerateUISchemaWrapper",
					Args: map[string]interface{}{
						"project_name": "{{project_name}}",
						"prompt":       "{{prompt}}",
						"theme":        "{{theme}}",
						"framework":    "{{framework}}",
					},
					ResultKey: "ui_schema",
				},
				{
					ID:           "gen_code",
					ActivityName: "GenerateSourceCodeWrapper",
					Args: map[string]interface{}{
						"schema": "{{ui_schema}}",
					},
					ResultKey: "source_code",
				},
				{
					ID:           "git_push",
					ActivityName: "PushToRepositoryWrapper",
					Args: map[string]interface{}{
						"project_name": "{{project_name}}",
						"prompt":       "{{prompt}}",
						"files":        "{{source_code}}",
					},
					ResultKey: "pr_url",
				},
				{
					ID:           "build_wasm",
					ActivityName: "BuildWebsiteBundleWrapper",
					Args: map[string]interface{}{
						"project_name": "{{project_name}}",
					},
					ResultKey: "bundle_path",
				},
				{
					ID:           "deploy_hosting",
					ActivityName: "DeployToHostingWrapper",
					Args: map[string]interface{}{
						"bundle_path": "{{bundle_path}}",
					},
					ResultKey: "deployment_result",
				},
			},
		}

		options := client.StartWorkflowOptions{
			ID:        fmt.Sprintf("deploy-site-%s", projectName),
			TaskQueue: "handoff-task-queue",
		}

		inputs := map[string]interface{}{
			"project_name": projectName,
			"prompt":       prompt,
			"theme":        theme,
			"framework":    framework,
		}

		we, err := s.temporalClient.ExecuteWorkflow(context.Background(), options, workflow.DynamicPipelineWorkflow, pipelineDef, inputs)
		if err != nil {
			return nil, err
		}

		return map[string]string{
			"status":      "started",
			"workflow_id": we.GetID(),
			"run_id":      we.GetRunID(),
			"message":     "Website deployment initiated via Dynamic Pipeline. You will be notified once it is live.",
		}, nil
	}

	if toolID == "system__crypto_analysis" {
		userID, _ := args["user_id"].(string)
		timeFrame, _ := args["time_frame"].(string)

		options := client.StartWorkflowOptions{
			ID:        domain.GetCryptoAnalysisPipelineID(userID),
			TaskQueue: "handoff-task-queue",
		}

		inputs := map[string]interface{}{
			"user_id":    userID,
			"time_frame": timeFrame,
		}

		we, err := s.temporalClient.ExecuteWorkflow(context.Background(), options, workflow.DynamicPipelineWorkflow, domain.CryptoAnalysisPipeline, inputs)
		if err != nil {
			return nil, err
		}

		return map[string]string{
			"status":      "started",
			"workflow_id": we.GetID(),
			"message":     "Crypto analysis pipeline started. Please approve the trade signal when notified.",
		}, nil
	}

	if toolID == "system__human_handoff" {
		roomID, _ := args["room_id"].(string)
		targetUserID, _ := args["target_user_id"].(string)

		options := client.StartWorkflowOptions{
			ID:        domain.GetHandoffPipelineID(roomID),
			TaskQueue: "handoff-task-queue",
		}

		inputs := map[string]interface{}{
			"room_id":        roomID,
			"target_user_id": targetUserID,
		}

		we, err := s.temporalClient.ExecuteWorkflow(context.Background(), options, workflow.DynamicPipelineWorkflow, domain.HumanHandoffPipeline, inputs)
		if err != nil {
			return nil, err
		}

		return map[string]string{
			"status":      "started",
			"workflow_id": we.GetID(),
			"message":     "Human handoff initiated. Waiting for an operator to join.",
		}, nil
	}

	if strings.HasPrefix(toolID, "mcp__") && strings.HasSuffix(toolID, "__execute") {
		// Format: mcp__{serverID}__execute
		parts := strings.Split(toolID, "__")
		if len(parts) >= 3 {
			serverID := parts[1]
			toolName, ok1 := args["tool_name"].(string)
			toolArgs, ok2 := args["arguments"]
			if ok1 && ok2 {
				return s.mcpSvc.ExecuteTool(serverID, toolName, toolArgs)
			}
			return nil, fmt.Errorf("invalid arguments for dynamic tool execution: missing tool_name or arguments")
		}
	}

	// Handle local tools (format: belt__tool)
	if strings.Contains(toolID, "__") {
		parts := strings.SplitN(toolID, "__", 2)
		beltName := parts[0]
		toolName := parts[1]
		// Local tools expect "args" in the arguments
		if toolArgs, ok := args["args"]; ok {
			return s.mcpSvc.ExecuteTool(beltName, toolName, toolArgs)
		}
		// Fallback: try passing all args if "args" key is missing (though manifest defines it)
		return s.mcpSvc.ExecuteTool(beltName, toolName, args)
	}

	return nil, fmt.Errorf("unknown tool ID format: %s", toolID)
}
