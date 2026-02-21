package gateway

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"nexus-super-node-v3/internal/core/services/mcp"

	"github.com/stretchr/testify/assert"
)

type MockEventProducer struct {
	ProducedEvents []struct {
		Key   []byte
		Value []byte
	}
}

func (m *MockEventProducer) Produce(ctx context.Context, key, value []byte) error {
	m.ProducedEvents = append(m.ProducedEvents, struct {
		Key   []byte
		Value []byte
	}{Key: key, Value: value})
	return nil
}

func TestWorkflowRun(t *testing.T) {
	// Setup
	mockProducer := &MockEventProducer{}
	gateway := NewEchoGateway(nil, nil, mockProducer)
	gateway.setupWorkflowRoutes()
	e := gateway.echo

	// Request Body
	workflowPayload := map[string]interface{}{
		"graph_id": "test-graph-123",
		"inputs": map[string]interface{}{
			"prompt": "hello rivet",
		},
	}
	bodyBytes, _ := json.Marshal(workflowPayload)

	// Request
	req := httptest.NewRequest(http.MethodPost, "/workflows/run", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	// Assertions
	assert.Equal(t, http.StatusOK, rec.Code)

	// Check Redpanda Production
	assert.Len(t, mockProducer.ProducedEvents, 1)
	assert.Equal(t, "workflow-run", string(mockProducer.ProducedEvents[0].Key))

	var producedPayload map[string]interface{}
	err := json.Unmarshal(mockProducer.ProducedEvents[0].Value, &producedPayload)
	assert.NoError(t, err)
	assert.Equal(t, "test-graph-123", producedPayload["graph_id"])
}

func TestGetAgents(t *testing.T) {
	// Setup
	gateway := NewEchoGateway(nil, nil, nil)
	gateway.setupAgentRoutes()
	e := gateway.echo

	// Request
	req := httptest.NewRequest(http.MethodGet, "/agents", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	// Assertions
	assert.Equal(t, http.StatusOK, rec.Code)

	var agents []map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &agents)
	assert.NoError(t, err)
	assert.Len(t, agents, 2)
	assert.Equal(t, "Wasm Agent V1", agents[0]["name"])
	assert.Equal(t, "Overlord Agent", agents[1]["name"])
}

func TestChatStream(t *testing.T) {
	// Setup
	gateway := NewEchoGateway(nil, nil, nil)
	gateway.setupChatRoutes()
	e := gateway.echo

	// Request Body
	requestBody := map[string]interface{}{
		"messages": []map[string]string{
			{"role": "user", "content": "Hello World"},
		},
		"model": "gpt-4o",
	}
	bodyBytes, _ := json.Marshal(requestBody)

	// Request
	req := httptest.NewRequest(http.MethodPost, "/api/chat", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	// Assertions
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "text/plain; charset=utf-8", rec.Header().Get("Content-Type"))
	assert.Equal(t, "chunked", rec.Header().Get("Transfer-Encoding"))

	// Check content
	// The mock implementation returns chunks: "Hello", " ", "from", " ", "Super", " ", "Node", "!\n", response
	responseBody := rec.Body.String()
	assert.Contains(t, responseBody, "Hello from Super Node!")
	assert.Contains(t, responseBody, "Hello World")
}

func TestToolExecute(t *testing.T) {
	// Setup MCP Service with a test tool
	mcpSvc := mcp.NewMCPService()
	testTool := mcp.Tool{
		Name:        "echo",
		Description: "Echoes the input",
		Action: func(args ...interface{}) (interface{}, error) {
			if len(args) > 0 {
				return args[0], nil
			}
			return "", nil
		},
	}
	mcpSvc.AddToolBelt(mcp.ToolBelt{
		Name:  "test-belt",
		Tools: []mcp.Tool{testTool},
	})

	// Setup Gateway
	gateway := NewEchoGateway(nil, mcpSvc, nil)
	gateway.setupToolRoutes()
	e := gateway.echo

	// Request Body
	requestBody := map[string]interface{}{
		"tool_belt_name": "test-belt",
		"tool_name":      "echo",
		"args":           []interface{}{"hello mcp"},
	}
	bodyBytes, _ := json.Marshal(requestBody)

	// Request
	req := httptest.NewRequest(http.MethodPost, "/tools/execute", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	// Assertions
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, "\"hello mcp\"", rec.Body.String())
}
