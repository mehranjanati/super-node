package mcp

// Tool represents a single tool that can be executed.
type Tool struct {
	Name        string
	Description string
	Action      func(...interface{}) (interface{}, error)
}

// ToolBelt is a collection of tools.
type ToolBelt struct {
	Name  string
	Tools []Tool
}

// MCPServerConfig defines how to connect to an external MCP server
type MCPServerConfig struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Type        string   `json:"type"` // "stdio" or "sse"
	Command     string   `json:"command"`
	Args        []string `json:"args"`
	URL         string   `json:"url"` // for sse
	Environment []string `json:"env"`
}
