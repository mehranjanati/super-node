package ports

import (
	"context"

	"nexus-super-node-v3/internal/core/domain"
)

// MCPRouter is the port for the MCP Aggregator Kernel.

type MCPRouter interface {
	// GetToolBelt returns the aggregated list of all available tools.
	GetToolBelt(ctx context.Context) ([]*domain.MCPTool, error)

	// RouteToolCall routes a tool call to the appropriate MCP server.
	RouteToolCall(ctx context.Context, toolID string, inputs map[string]interface{}) (map[string]interface{}, error)

	// RegisterServer registers a new MCP server with the aggregator.
	RegisterServer(server *domain.MCPServer)
}
