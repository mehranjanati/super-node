package mcp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"sync"
)

// StdioProxy manages an external MCP server process via stdio
type StdioProxy struct {
	Config MCPServerConfig
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	stdout io.ReadCloser
	scanner *bufio.Scanner
	mu      sync.Mutex
	requestID int
}

// NewStdioProxy creates a new proxy for a stdio-based MCP server
func NewStdioProxy(config MCPServerConfig) *StdioProxy {
	return &StdioProxy{
		Config: config,
	}
}

// Start launches the MCP server process
func (p *StdioProxy) Start() error {
	p.cmd = exec.Command(p.Config.Command, p.Config.Args...)
	
	var err error
	p.stdin, err = p.cmd.StdinPipe()
	if err != nil {
		return err
	}
	
	p.stdout, err = p.cmd.StdoutPipe()
	if err != nil {
		return err
	}
	
	p.scanner = bufio.NewScanner(p.stdout)
	
	if err := p.cmd.Start(); err != nil {
		return err
	}
	
	// Initial MCP Handshake (initialize request) should happen here in a full implementation
	return nil
}

// ExecuteCall sends a JSON-RPC call to the MCP server
func (p *StdioProxy) ExecuteCall(method string, params interface{}) (json.RawMessage, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	
	p.requestID++
	request := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      p.requestID,
		"method":  method,
		"params":  params,
	}
	
	data, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	
	if _, err := fmt.Fprintln(p.stdin, string(data)); err != nil {
		return nil, err
	}
	
	if !p.scanner.Scan() {
		return nil, fmt.Errorf("failed to read response from MCP server")
	}
	
	var response struct {
		ID     interface{}     `json:"id"`
		Result json.RawMessage `json:"result"`
		Error  interface{}     `json:"error"`
	}
	
	if err := json.Unmarshal(p.scanner.Bytes(), &response); err != nil {
		return nil, err
	}
	
	if response.Error != nil {
		return nil, fmt.Errorf("mcp server error: %v", response.Error)
	}
	
	return response.Result, nil
}

// Stop terminates the MCP server process
func (p *StdioProxy) Stop() error {
	if p.cmd != nil && p.cmd.Process != nil {
		return p.cmd.Process.Kill()
	}
	return nil
}
