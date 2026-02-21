package mcp

import "go.uber.org/fx"

// Module is the Fx module for the MCP service.
var Module = fx.Module("mcp",
	fx.Provide(NewMCPService),
)
