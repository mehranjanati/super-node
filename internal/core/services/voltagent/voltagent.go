package voltagent

import (
	"context"
	"fmt"
	"strings"
	"time"

	"nexus-super-node-v3/internal/adapters/ai"
	"nexus-super-node-v3/internal/adapters/voltagentclient"
	"nexus-super-node-v3/internal/core/domain"
	"nexus-super-node-v3/internal/workflow"

	enumspb "go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/client"
)

const (
	planningSourceRemote   = "remote_voltagent"
	planningSourceEmbedded = "embedded_fallback"
)

type deploymentStartResult struct {
	Status         string
	WorkflowID     string
	RunID          string
	CurrentStep    string
	Message        string
	PlanningSource string
	Data           map[string]string
}

func (r deploymentStartResult) asStringMap() map[string]string {
	return map[string]string{
		"status":          r.Status,
		"workflow_id":     r.WorkflowID,
		"run_id":          r.RunID,
		"current_step":    r.CurrentStep,
		"message":         r.Message,
		"planning_source": r.PlanningSource,
	}
}

func (r deploymentStartResult) asAnyMap() map[string]interface{} {
	return map[string]interface{}{
		"status":          r.Status,
		"workflow_id":     r.WorkflowID,
		"run_id":          r.RunID,
		"current_step":    r.CurrentStep,
		"message":         r.Message,
		"planning_source": r.PlanningSource,
		"data":            r.Data,
	}
}

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

// ExecuteTool forwards a request from VoltAgent to the MCP service or triggers internal workflows.
func (s *VoltAgentService) ExecuteTool(ctx context.Context, toolID string, args map[string]interface{}, meta RequestMetadata) (interface{}, error) {
	if toolID == "system__deploy_website" {
		input, err := websiteDeploymentInputFromArgs(args)
		if err != nil {
			return nil, err
		}
		result, err := s.startWebsiteDeployment(ctx, input, meta)
		if err != nil {
			return nil, err
		}
		return result.asStringMap(), nil
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

func (s *VoltAgentService) StartWebsiteDeployment(ctx context.Context, input WebsiteDeploymentInput, meta RequestMetadata) (map[string]interface{}, error) {
	result, err := s.startWebsiteDeployment(ctx, input, meta)
	if err != nil {
		return nil, err
	}
	return result.asAnyMap(), nil
}

func (s *VoltAgentService) GetWorkflowStatus(workflowID string) (*workflow.DummyDeploymentStatus, error) {
	if s.temporalClient == nil {
		return nil, fmt.Errorf("temporal client not configured")
	}

	queryResult, err := s.temporalClient.QueryWorkflow(context.Background(), workflowID, "", workflow.DummyDeploymentStatusQuery)
	if err == nil {
		var status workflow.DummyDeploymentStatus
		if getErr := queryResult.Get(&status); getErr == nil {
			status.WorkflowID = workflowID
			return &status, nil
		}
	}

	description, describeErr := s.temporalClient.DescribeWorkflowExecution(context.Background(), workflowID, "")
	if describeErr != nil {
		if err != nil {
			return nil, err
		}
		return nil, describeErr
	}

	info := description.GetWorkflowExecutionInfo()
	if info != nil && info.Status == enumspb.WORKFLOW_EXECUTION_STATUS_COMPLETED {
		var status workflow.DummyDeploymentStatus
		if getErr := s.temporalClient.GetWorkflow(context.Background(), workflowID, "").Get(context.Background(), &status); getErr == nil {
			status.WorkflowID = workflowID
			return &status, nil
		}
	}

	status := &workflow.DummyDeploymentStatus{
		WorkflowID:  workflowID,
		Status:      mapTemporalStatus(info.GetStatus()),
		CurrentStep: "INIT",
		Logs:        []string{"Workflow status resolved from Temporal metadata."},
	}
	if info != nil && info.CloseTime != nil && info.Status == enumspb.WORKFLOW_EXECUTION_STATUS_COMPLETED {
		status.CurrentStep = "DONE"
	}

	return status, nil
}

func sanitizeWorkflowToken(value string) string {
	token := strings.ToLower(strings.TrimSpace(value))
	token = strings.ReplaceAll(token, " ", "-")
	token = strings.ReplaceAll(token, "_", "-")
	token = strings.ReplaceAll(token, "/", "-")
	token = strings.ReplaceAll(token, ".", "-")
	if token == "" {
		return "deploy"
	}
	return token
}

func mapTemporalStatus(status enumspb.WorkflowExecutionStatus) string {
	switch status {
	case enumspb.WORKFLOW_EXECUTION_STATUS_COMPLETED:
		return "COMPLETED"
	case enumspb.WORKFLOW_EXECUTION_STATUS_FAILED, enumspb.WORKFLOW_EXECUTION_STATUS_TERMINATED, enumspb.WORKFLOW_EXECUTION_STATUS_TIMED_OUT, enumspb.WORKFLOW_EXECUTION_STATUS_CANCELED:
		return "FAILED"
	default:
		return "RUNNING"
	}
}

func (s *VoltAgentService) startWebsiteDeployment(ctx context.Context, input WebsiteDeploymentInput, meta RequestMetadata) (deploymentStartResult, error) {
	normalized, err := normalizeWebsiteDeploymentInput(input)
	if err != nil {
		return deploymentStartResult{}, err
	}

	if s.remotePlanningEnabled() && s.remoteClient != nil {
		result, err := s.startRemotePlannedWebsiteDeployment(ctx, normalized, meta)
		if err == nil {
			return result, nil
		}
		if !s.embeddedFallbackEnabled() {
			return deploymentStartResult{}, err
		}
		return s.startEmbeddedWebsiteDeployment(ctx, normalized, fmt.Sprintf("Website deployment initiated via embedded fallback after remote planning failed: %v", err))
	}

	if s.remotePlanningEnabled() && s.remoteClient == nil && !s.embeddedFallbackEnabled() {
		return deploymentStartResult{}, fmt.Errorf("voltagent remote planning is enabled but client is not configured")
	}

	if !s.remotePlanningEnabled() && !s.embeddedFallbackEnabled() {
		return deploymentStartResult{}, fmt.Errorf("voltagent remote planning is disabled and embedded fallback is not allowed")
	}

	return s.startEmbeddedWebsiteDeployment(ctx, normalized, "Website deployment initiated via embedded fallback.")
}

func (s *VoltAgentService) startRemotePlannedWebsiteDeployment(ctx context.Context, input WebsiteDeploymentInput, meta RequestMetadata) (deploymentStartResult, error) {
	planReq := &voltagentclient.PlanRequest{
		ContractVersion: s.contractVersion(),
		Intent:          "deploy_website",
		Input: map[string]any{
			"project_name": input.ProjectName,
			"prompt":       input.Prompt,
			"framework":    input.Framework,
			"theme":        input.Theme,
		},
		Context: map[string]any{
			"source": firstNonEmptyString(meta.Source, "go-gateway"),
		},
	}
	if meta.RequestID != "" {
		planReq.Context["request_id"] = meta.RequestID
	}
	if meta.CorrelationID != "" {
		planReq.Context["correlation_id"] = meta.CorrelationID
	}

	planResp, err := s.remoteClient.Plan(ctx, planReq, voltagentclient.RequestOptions{
		RequestID:     meta.RequestID,
		CorrelationID: meta.CorrelationID,
	})
	if err != nil {
		return deploymentStartResult{}, err
	}

	return s.executeRemoteWebsitePlan(ctx, input, planResp.Plan)
}

func (s *VoltAgentService) executeRemoteWebsitePlan(ctx context.Context, input WebsiteDeploymentInput, plan *voltagentclient.ExecutionPlan) (deploymentStartResult, error) {
	if plan == nil {
		return deploymentStartResult{}, fmt.Errorf("remote voltagent plan is missing")
	}
	if plan.Intent != "deploy_website" {
		return deploymentStartResult{}, fmt.Errorf("unsupported remote plan intent: %s", plan.Intent)
	}
	if plan.Kind != "workflow" {
		return deploymentStartResult{}, fmt.Errorf("unsupported remote plan kind: %s", plan.Kind)
	}
	if plan.ExecutionTarget != "go-temporal" {
		return deploymentStartResult{}, fmt.Errorf("unsupported remote execution target: %s", plan.ExecutionTarget)
	}
	if plan.Workflow.Action != "start_dynamic_pipeline" {
		return deploymentStartResult{}, fmt.Errorf("unsupported remote workflow action: %s", plan.Workflow.Action)
	}

	plannedInput := mergeWebsiteDeploymentInput(input, plan.Workflow.Input)
	result, err := s.executeWebsiteDeploymentWorkflow(ctx, plannedInput, planningSourceRemote)
	if err != nil {
		return deploymentStartResult{}, err
	}

	if len(plan.Warnings) > 0 {
		result.Message = fmt.Sprintf("%s Warnings: %s", result.Message, strings.Join(plan.Warnings, "; "))
	}
	return result, nil
}

func (s *VoltAgentService) startEmbeddedWebsiteDeployment(ctx context.Context, input WebsiteDeploymentInput, message string) (deploymentStartResult, error) {
	result, err := s.executeWebsiteDeploymentWorkflow(ctx, input, planningSourceEmbedded)
	if err != nil {
		return deploymentStartResult{}, err
	}
	result.Message = message
	return result, nil
}

func (s *VoltAgentService) executeWebsiteDeploymentWorkflow(ctx context.Context, input WebsiteDeploymentInput, planningSource string) (deploymentStartResult, error) {
	if s.temporalClient == nil {
		return deploymentStartResult{}, fmt.Errorf("temporal client not configured")
	}

	workflowID := fmt.Sprintf("deploy-site-%s-%d", sanitizeWorkflowToken(input.ProjectName), time.Now().UnixMilli())
	options := client.StartWorkflowOptions{
		ID:        workflowID,
		TaskQueue: "handoff-task-queue",
	}

	workflowInputs := map[string]interface{}{
		"project_name": input.ProjectName,
		"prompt":       input.Prompt,
		"theme":        input.Theme,
		"framework":    input.Framework,
	}

	we, err := s.temporalClient.ExecuteWorkflow(ctx, options, workflow.DynamicPipelineWorkflow, websiteDeploymentPipelineDefinition(), workflowInputs)
	if err != nil {
		return deploymentStartResult{}, err
	}

	message := "Website deployment initiated via remote VoltAgent plan."
	if planningSource == planningSourceEmbedded {
		message = "Website deployment initiated via embedded fallback."
	}

	return deploymentStartResult{
		Status:         "started",
		WorkflowID:     we.GetID(),
		RunID:          we.GetRunID(),
		CurrentStep:    "INIT",
		Message:        message,
		PlanningSource: planningSource,
		Data: map[string]string{
			"project_name": input.ProjectName,
			"template":     firstNonEmptyString(input.Template, input.Framework),
			"framework":    input.Framework,
			"theme":        input.Theme,
		},
	}, nil
}

func websiteDeploymentPipelineDefinition() domain.PipelineDefinition {
	return domain.PipelineDefinition{
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
}

func websiteDeploymentInputFromArgs(args map[string]interface{}) (WebsiteDeploymentInput, error) {
	if args == nil {
		return WebsiteDeploymentInput{}, fmt.Errorf("deployment arguments are required")
	}

	return WebsiteDeploymentInput{
		ProjectName: stringValue(args["project_name"]),
		Prompt:      stringValue(args["prompt"]),
		Template:    stringValue(args["template"]),
		Theme:       stringValue(args["theme"]),
		Framework:   stringValue(args["framework"]),
	}, nil
}

func normalizeWebsiteDeploymentInput(input WebsiteDeploymentInput) (WebsiteDeploymentInput, error) {
	normalized := input
	normalized.ProjectName = strings.TrimSpace(normalized.ProjectName)
	normalized.Prompt = strings.TrimSpace(normalized.Prompt)
	normalized.Template = strings.TrimSpace(normalized.Template)
	normalized.Theme = strings.TrimSpace(normalized.Theme)
	normalized.Framework = strings.TrimSpace(normalized.Framework)

	if normalized.ProjectName == "" {
		return WebsiteDeploymentInput{}, fmt.Errorf("project_name is required")
	}
	if normalized.Framework == "" {
		normalized.Framework = firstNonEmptyString(normalized.Template, "svelte")
	}
	if normalized.Theme == "" {
		normalized.Theme = "minimal"
	}
	if normalized.Template == "" {
		normalized.Template = normalized.Framework
	}
	if normalized.Prompt == "" {
		normalized.Prompt = fmt.Sprintf("Deploy a %s website for project %s.", normalized.Framework, normalized.ProjectName)
	}

	return normalized, nil
}

func mergeWebsiteDeploymentInput(base WebsiteDeploymentInput, workflowInput map[string]any) WebsiteDeploymentInput {
	merged := base
	if workflowInput == nil {
		return merged
	}

	if value := stringValue(workflowInput["project_name"]); value != "" {
		merged.ProjectName = value
	}
	if value := stringValue(workflowInput["prompt"]); value != "" {
		merged.Prompt = value
	}
	if value := stringValue(workflowInput["theme"]); value != "" {
		merged.Theme = value
	}
	if value := stringValue(workflowInput["framework"]); value != "" {
		merged.Framework = value
		if merged.Template == "" {
			merged.Template = value
		}
	}

	normalized, err := normalizeWebsiteDeploymentInput(merged)
	if err != nil {
		return merged
	}
	return normalized
}

func (s *VoltAgentService) remotePlanningEnabled() bool {
	if s.cfg == nil {
		return s.remoteClient != nil
	}
	return s.cfg.VoltAgent.Enabled
}

func (s *VoltAgentService) embeddedFallbackEnabled() bool {
	if s.cfg == nil {
		return true
	}
	return s.cfg.VoltAgent.UseEmbeddedFallback
}

func (s *VoltAgentService) contractVersion() string {
	if s.cfg == nil || strings.TrimSpace(s.cfg.VoltAgent.ContractVersion) == "" {
		return voltagentclient.DefaultContractVersion
	}
	return s.cfg.VoltAgent.ContractVersion
}

func stringValue(value interface{}) string {
	if text, ok := value.(string); ok {
		return strings.TrimSpace(text)
	}
	return ""
}

func firstNonEmptyString(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}
