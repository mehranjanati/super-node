package mcp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMCPService(t *testing.T) {
	mcpService := NewMCPService()
	assert.NotNil(t, mcpService)
	assert.NotNil(t, mcpService.ToolBelts)
}

func TestAddAndGetToolBelt(t *testing.T) {
	mcpService := NewMCPService()
	sampleToolBelt := NewSampleToolBelt()
	mcpService.AddToolBelt(sampleToolBelt)

	retrievedToolBelt, err := mcpService.GetToolBelt("sample")
	assert.NoError(t, err)
	assert.NotNil(t, retrievedToolBelt)
	assert.Equal(t, "sample", retrievedToolBelt.Name)
}

func TestGetToolBeltNotFound(t *testing.T) {
	mcpService := NewMCPService()
	_, err := mcpService.GetToolBelt("non_existent_tool_belt")
	assert.Error(t, err)
}

func TestExecuteTool(t *testing.T) {
	mcpService := NewMCPService()
	sampleToolBelt := NewSampleToolBelt()
	mcpService.AddToolBelt(sampleToolBelt)

	result, err := mcpService.ExecuteTool("sample", "hello")
	assert.NoError(t, err)
	assert.Equal(t, "Hello, world!", result)
}

func TestExecuteToolNotFound(t *testing.T) {
	mcpService := NewMCPService()
	_, err := mcpService.ExecuteTool("non_existent_tool_belt", "hello", nil)
	assert.Error(t, err)
}
