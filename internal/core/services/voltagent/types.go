package voltagent

import (
	"nexus-super-node-v3/internal/adapters/ai"
	"nexus-super-node-v3/internal/adapters/voltagentclient"
	"nexus-super-node-v3/internal/config"
	"nexus-super-node-v3/internal/core/services/mcp"

	"go.temporal.io/sdk/client"
)

// VoltAgentManifest represents the tools and configuration shared with VoltAgent
type VoltAgentManifest struct {
	Version string          `json:"version"`
	Tools   []VoltAgentTool `json:"tools"`
}

// VoltAgentTool is a tool definition compatible with VoltAgent's expectation
type VoltAgentTool struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Parameters  interface{} `json:"parameters"` // JSON Schema
}

type RequestMetadata struct {
	RequestID     string
	CorrelationID string
	Source        string
}

type WebsiteDeploymentInput struct {
	ProjectName string
	Prompt      string
	Template    string
	Theme       string
	Framework   string
}

// VoltAgentService handles communication with the VoltAgent reasoning engine
type VoltAgentService struct {
	mcpSvc         *mcp.MCPService
	temporalClient client.Client
	aiClient       *ai.OpenAIClient
	cfg            *config.Config
	remoteClient   *voltagentclient.Client
}

func NewVoltAgentService(cfg *config.Config, mcpSvc *mcp.MCPService, temporalClient client.Client, aiClient *ai.OpenAIClient, remoteClient *voltagentclient.Client) *VoltAgentService {
	return &VoltAgentService{
		mcpSvc:         mcpSvc,
		temporalClient: temporalClient,
		aiClient:       aiClient,
		cfg:            cfg,
		remoteClient:   remoteClient,
	}
}
