package services

import (
	"context"
	"fmt"

	"nexus-super-node-v3/internal/core/domain"
	"nexus-super-node-v3/internal/core/services/mcp"
	"nexus-super-node-v3/internal/ports"
)

// mcpAggregator is the implementation of the MCPRouter port.

type mcpAggregator struct {
	servers []*domain.MCPServer
	mcpSvc  *mcp.MCPService
}

// NewMCPAggregator creates a new MCP aggregator.

func NewMCPAggregator(mcpSvc *mcp.MCPService) ports.MCPRouter {
	return &mcpAggregator{
		servers: make([]*domain.MCPServer, 0),
		mcpSvc:  mcpSvc,
	}
}

// GetToolBelt returns the aggregated list of all available tools.

func (a *mcpAggregator) GetToolBelt(ctx context.Context) ([]*domain.MCPTool, error) {
	var toolbelt []*domain.MCPTool
	for _, server := range a.servers {
		toolbelt = append(toolbelt, server.Tools...)
	}
	return toolbelt, nil
}

// RouteToolCall routes a tool call to the appropriate MCP server.

func (a *mcpAggregator) RouteToolCall(ctx context.Context, toolID string, inputs map[string]interface{}) (map[string]interface{}, error) {
	for _, server := range a.servers {
		for _, tool := range server.Tools {
			if tool.ID == toolID {
				// This is a simplified execution model. A real implementation would need to
				// handle different server types (e.g., HTTP, gRPC, WASM) and their specific
				// invocation methods.
				result, err := a.mcpSvc.ExecuteTool(server.ID, tool.ID, toInterfaceSlice(inputs)...)
				if err != nil {
					return nil, fmt.Errorf("error executing tool %s on server %s: %w", toolID, server.ID, err)
				}
				return map[string]interface{}{"result": result}, nil
			}
		}
	}
	return nil, fmt.Errorf("tool %s not found in any registered MCP server", toolID)
}

func (a *mcpAggregator) RegisterServer(server *domain.MCPServer) {
	a.servers = append(a.servers, server)
}

func toInterfaceSlice(m map[string]interface{}) []interface{} {
	var s []interface{}
	for _, v := range m {
		s = append(s, v)
	}
	return s
}
