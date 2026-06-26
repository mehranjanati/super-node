package gateway

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"nexus-super-node-v3/internal/adapters/voltagentclient"
	"nexus-super-node-v3/internal/core/domain"
	"nexus-super-node-v3/internal/core/services/mcp"
	"nexus-super-node-v3/internal/core/services/voltagent"
	workflowpkg "nexus-super-node-v3/internal/workflow"

	"github.com/stretchr/testify/assert"
)

type MockEventProducer struct {
	ProducedEvents []struct {
		Key   []byte
		Value []byte
	}
}

type MockAppDataRepository struct {
	Data map[string][]byte
}

func (m *MockAppDataRepository) GetAppData(ctx context.Context, id string) (*domain.AppData, error) {
	if m == nil || m.Data == nil {
		return nil, nil
	}
	value, ok := m.Data[id]
	if !ok {
		return nil, nil
	}
	return &domain.AppData{ID: id, Data: value}, nil
}

func (m *MockAppDataRepository) CreateAppData(ctx context.Context, id string, data []byte) error {
	if m.Data == nil {
		m.Data = map[string][]byte{}
	}
	m.Data[id] = data
	return nil
}

func (m *MockAppDataRepository) UpsertAppData(ctx context.Context, id string, data []byte) error {
	if m.Data == nil {
		m.Data = map[string][]byte{}
	}
	m.Data[id] = data
	return nil
}

func (m *MockEventProducer) Produce(ctx context.Context, key, value []byte) error {
	m.ProducedEvents = append(m.ProducedEvents, struct {
		Key   []byte
		Value []byte
	}{Key: key, Value: value})
	return nil
}

type MockAgentService struct {
	ListAgentsFunc func(ctx context.Context, ownerID string) ([]*domain.Agent, error)
}

func (m *MockAgentService) CreateAgent(ctx context.Context, agent *domain.Agent) error {
	return nil
}

func (m *MockAgentService) GetAgent(ctx context.Context, id string) (*domain.Agent, error) {
	return nil, nil
}

func (m *MockAgentService) ListAgents(ctx context.Context, ownerID string) ([]*domain.Agent, error) {
	if m.ListAgentsFunc != nil {
		return m.ListAgentsFunc(ctx, ownerID)
	}
	return nil, nil
}

func (m *MockAgentService) UpdateAgent(ctx context.Context, agent *domain.Agent) error {
	return nil
}

func (m *MockAgentService) DeleteAgent(ctx context.Context, id string) error {
	return nil
}

func (m *MockAgentService) DeployAgent(ctx context.Context, id string) error {
	return nil
}

func (m *MockAgentService) PauseAgent(ctx context.Context, id string) error {
	return nil
}

func TestWorkflowRun(t *testing.T) {
	// Setup
	mockProducer := &MockEventProducer{}
	gateway := NewEchoGateway(nil, nil, nil, nil, nil, nil, nil, nil, mockProducer, nil)
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
	mockAgentSvc := &MockAgentService{
		ListAgentsFunc: func(ctx context.Context, ownerID string) ([]*domain.Agent, error) {
			now := time.Now()
			return []*domain.Agent{
				{ID: "agent-1", Name: "Wasm Agent V1", OwnerID: ownerID, CreatedAt: now, UpdatedAt: now},
				{ID: "agent-2", Name: "Overlord Agent", OwnerID: ownerID, CreatedAt: now, UpdatedAt: now},
			}, nil
		},
	}
	gateway := NewEchoGateway(nil, nil, nil, nil, nil, nil, nil, mockAgentSvc, nil, nil)
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
	t.Skip("Streaming route now depends on a configured VoltAgent AI client")
	// Setup
	gateway := NewEchoGateway(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
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
	gateway := NewEchoGateway(nil, nil, mcpSvc, nil, nil, nil, nil, nil, nil, nil)
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

func TestInternalHealthAggregatesVoltAgentStatus(t *testing.T) {
	voltAgentServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/health", r.URL.Path)

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"status":"ok",
			"service":"voltagent-service",
			"contract_version":"v1alpha1",
			"checks":{"api":"ok","planner":"ok"}
		}`))
	}))
	defer voltAgentServer.Close()

	remoteClient := voltagentclient.NewClient(voltAgentServer.URL, time.Second)
	voltSvc := voltagent.NewVoltAgentService(nil, nil, nil, nil, remoteClient)

	gateway := NewEchoGateway(nil, nil, nil, voltSvc, nil, nil, nil, nil, nil, nil)
	gateway.setupHealthRoutes()

	req := httptest.NewRequest(http.MethodGet, "/internal/health", nil)
	rec := httptest.NewRecorder()
	gateway.echo.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var payload map[string]any
	err := json.Unmarshal(rec.Body.Bytes(), &payload)
	assert.NoError(t, err)
	assert.Equal(t, "ok", payload["status"])

	dependencies, ok := payload["dependencies"].(map[string]any)
	assert.True(t, ok)

	voltagentDependency, ok := dependencies["voltagent_service"].(map[string]any)
	assert.True(t, ok)
	assert.Equal(t, "ok", voltagentDependency["status"])
	assert.Equal(t, "voltagent-service", voltagentDependency["service"])
	assert.Equal(t, "v1alpha1", voltagentDependency["contract_version"])
}

func TestInternalHealthReturnsDegradedWhenVoltAgentIsUnavailable(t *testing.T) {
	remoteClient := voltagentclient.NewClient("http://127.0.0.1:1", 100*time.Millisecond)
	voltSvc := voltagent.NewVoltAgentService(nil, nil, nil, nil, remoteClient)

	gateway := NewEchoGateway(nil, nil, nil, voltSvc, nil, nil, nil, nil, nil, nil)
	gateway.setupHealthRoutes()

	req := httptest.NewRequest(http.MethodGet, "/internal/health", nil)
	rec := httptest.NewRecorder()
	gateway.echo.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusServiceUnavailable, rec.Code)

	var payload map[string]any
	err := json.Unmarshal(rec.Body.Bytes(), &payload)
	assert.NoError(t, err)
	assert.Equal(t, "degraded", payload["status"])

	dependencies, ok := payload["dependencies"].(map[string]any)
	assert.True(t, ok)

	voltagentDependency, ok := dependencies["voltagent_service"].(map[string]any)
	assert.True(t, ok)
	assert.Equal(t, "unavailable", voltagentDependency["status"])
	assert.Equal(t, "VOLTAGENT_UNAVAILABLE", voltagentDependency["error_code"])
}

func TestWorkflowReadModelRoutes(t *testing.T) {
	record := domain.WorkflowExecutionRecord{
		WorkflowID:  "dummy-deploy-demo-1",
		Kind:        "deployment",
		Name:        "demo-site",
		Status:      "RUNNING",
		CurrentStep: "BUILDING",
		Logs:        []string{"Dummy deployment workflow started.", "Building deployment artifact."},
		Artifacts: &domain.WorkflowExecutionArtifacts{
			ProjectName:    "demo-site",
			Template:       "svelte",
			PlanningSource: "remote_voltagent",
			Message:        "Dummy deployment workflow started.",
		},
		CreatedAt: time.Now().Add(-time.Minute),
		UpdatedAt: time.Now(),
	}

	recordPayload, err := json.Marshal(record)
	assert.NoError(t, err)

	indexPayload, err := json.Marshal(map[string][]string{
		"ids": {record.WorkflowID},
	})
	assert.NoError(t, err)

	repo := &MockAppDataRepository{
		Data: map[string][]byte{
			workflowExecutionIndexKey:               indexPayload,
			workflowExecutionKey(record.WorkflowID): recordPayload,
		},
	}

	gateway := NewEchoGateway(nil, repo, nil, nil, nil, nil, nil, nil, nil, nil)
	gateway.setupWorkflowRoutes()
	e := gateway.echo

	req := httptest.NewRequest(http.MethodGet, "/workflows", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "dummy-deploy-demo-1")
	assert.Contains(t, rec.Body.String(), "demo-site")
	assert.Contains(t, rec.Body.String(), "remote_voltagent")

	req = httptest.NewRequest(http.MethodGet, "/logs", nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "Building deployment artifact.")
	assert.Contains(t, rec.Body.String(), "TEMPORAL_DEPLOY")
}

func TestPersistWorkflowStatusKeepsPlanningSource(t *testing.T) {
	repo := &MockAppDataRepository{}
	gateway := NewEchoGateway(nil, repo, nil, nil, nil, nil, nil, nil, nil, nil)

	err := gateway.persistWorkflowStart(context.Background(), "wf-123", "deployment", "demo-site", "INIT", []string{"started"}, &domain.WorkflowExecutionArtifacts{
		ProjectName:    "demo-site",
		Template:       "svelte",
		PlanningSource: "embedded_fallback",
		Message:        "Website deployment workflow started.",
	})
	assert.NoError(t, err)

	record, err := gateway.persistWorkflowStatus(context.Background(), "wf-123", &workflowpkg.DummyDeploymentStatus{
		Status:      "COMPLETED",
		CurrentStep: "DONE",
		Logs:        []string{"Deployment finished successfully."},
		Artifacts: &workflowpkg.DummyDeploymentArtifacts{
			ProjectName: "demo-site",
			Template:    "svelte",
			URL:         "https://demo.example.com",
			Message:     "Deployment finished successfully.",
		},
	})
	assert.NoError(t, err)
	if assert.NotNil(t, record) && assert.NotNil(t, record.Artifacts) {
		assert.Equal(t, "embedded_fallback", record.Artifacts.PlanningSource)
		assert.Equal(t, "https://demo.example.com", record.Artifacts.URL)
	}
}
