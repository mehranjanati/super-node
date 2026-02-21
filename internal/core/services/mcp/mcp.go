package mcp

import (
	"fmt"
	"sync"
)

// MCPService is the service for the MCP Aggregator Kernel.
type MCPService struct {
	ToolBelts      map[string]ToolBelt
	DynamicServers map[string]MCPServerConfig
	ActiveProxies  map[string]*StdioProxy
	mu             sync.Mutex
}

// NewMCPService creates a new MCPService.
func NewMCPService() *MCPService {
	return &MCPService{
		ToolBelts:      make(map[string]ToolBelt),
		DynamicServers: make(map[string]MCPServerConfig),
		ActiveProxies:  make(map[string]*StdioProxy),
	}
}

// RegisterDynamicServer registers an external MCP server
func (s *MCPService) RegisterDynamicServer(config MCPServerConfig) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.DynamicServers[config.ID] = config

	// If it's stdio, we can try to start it and discover tools
	if config.Type == "stdio" {
		proxy := NewStdioProxy(config)
		if err := proxy.Start(); err != nil {
			return fmt.Errorf("failed to start mcp server %s: %w", config.ID, err)
		}
		s.ActiveProxies[config.ID] = proxy

		// Discover tools from the dynamic server
		go func() {
			result, err := proxy.ExecuteCall("tools/list", nil)
			if err != nil {
				fmt.Printf("failed to discover tools for %s: %v\n", config.ID, err)
				return
			}
			fmt.Printf("Discovered tools for %s: %s\n", config.ID, string(result))
		}()
	}

	return nil
}

// ListDynamicServers returns all registered dynamic servers
func (s *MCPService) ListDynamicServers() []MCPServerConfig {
	s.mu.Lock()
	defer s.mu.Unlock()

	servers := make([]MCPServerConfig, 0, len(s.DynamicServers))
	for _, config := range s.DynamicServers {
		servers = append(servers, config)
	}
	return servers
}

// AddToolBelt adds a tool belt to the MCP service.
func (s *MCPService) AddToolBelt(toolBelt ToolBelt) {
	s.ToolBelts[toolBelt.Name] = toolBelt
}

// GetToolBelt retrieves a tool belt from the MCP service.
func (s *MCPService) GetToolBelt(name string) (ToolBelt, error) {
	toolBelt, ok := s.ToolBelts[name]
	if !ok {
		return ToolBelt{}, fmt.Errorf("tool belt %s not found", name)
	}
	return toolBelt, nil
}

// ExecuteTool executes a tool from a tool belt.
func (s *MCPService) ExecuteTool(toolBeltName, toolName string, args ...interface{}) (interface{}, error) {
	// 1. Try local tool belts first
	toolBelt, err := s.GetToolBelt(toolBeltName)
	if err == nil {
		for _, tool := range toolBelt.Tools {
			if tool.Name == toolName {
				return tool.Action(args...)
			}
		}
	}

	// 2. Try dynamic servers
	s.mu.Lock()
	proxy, ok := s.ActiveProxies[toolBeltName]
	s.mu.Unlock()

	if ok {
		// Prepare parameters for MCP call
		params := map[string]interface{}{
			"name":      toolName,
			"arguments": args,
		}

		// In MCP protocol, executing a tool is 'tools/call'
		result, err := proxy.ExecuteCall("tools/call", params)
		if err != nil {
			return nil, fmt.Errorf("error calling dynamic tool %s/%s: %w", toolBeltName, toolName, err)
		}
		return result, nil
	}

	return nil, fmt.Errorf("tool %s not found in tool belt %s", toolName, toolBeltName)
}
