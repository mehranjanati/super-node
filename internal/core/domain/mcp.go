package domain

// MCPTool represents a single tool that can be executed.

type MCPTool struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	Pattern     string `json:"pattern"`
}

// MCPServer represents a server that provides a set of tools.

type MCPServer struct {
	ID    string      `json:"id"`
	Type  string      `json:"type"`
	Tools []*MCPTool `json:"tools"`
}
