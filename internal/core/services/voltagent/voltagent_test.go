package voltagent

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.temporal.io/sdk/client"

	"nexus-super-node-v3/internal/adapters/voltagentclient"
	"nexus-super-node-v3/internal/config"
)

// MockTemporalClient is a basic mock for Temporal client
type MockTemporalClient struct {
	mock.Mock
	client.Client // Embed to satisfy the interface implicitly for unmocked methods
}

func (m *MockTemporalClient) ExecuteWorkflow(ctx context.Context, options client.StartWorkflowOptions, workflow interface{}, args ...interface{}) (client.WorkflowRun, error) {
	callArgs := m.Called(ctx, options, workflow, args)
	return callArgs.Get(0).(client.WorkflowRun), callArgs.Error(1)
}

// MockWorkflowRun mocks a Temporal WorkflowRun
type MockWorkflowRun struct {
	mock.Mock
}

func (m *MockWorkflowRun) GetID() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockWorkflowRun) GetRunID() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockWorkflowRun) Get(ctx context.Context, valuePtr interface{}) error {
	args := m.Called(ctx, valuePtr)
	return args.Error(0)
}

func (m *MockWorkflowRun) GetWithOptions(ctx context.Context, valuePtr interface{}, options client.WorkflowRunGetOptions) error {
	args := m.Called(ctx, valuePtr, options)
	return args.Error(0)
}

func TestVoltAgentService_HealthAndFallbackValidation(t *testing.T) {
	// Scenario:
	// 1. Remote VoltAgent is up. Health check succeeds.
	// 2. Remote Planning succeeds -> planningSourceRemote
	// 3. Remote VoltAgent goes down. Health check fails.
	// 4. Remote Planning fails -> embedded fallback -> planningSourceEmbedded
	// 5. Remote VoltAgent comes back up. Health check succeeds again.

	remoteStatus := "ok"
	remotePlanStatus := "ok"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/health" {
			if remoteStatus == "ok" {
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"status":"ok","service":"voltagent","contract_version":"v1alpha1"}`))
			} else {
				w.WriteHeader(http.StatusServiceUnavailable)
			}
			return
		}

		if r.URL.Path == "/plan" {
			if remotePlanStatus == "ok" {
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{
					"status":"ok",
					"contract_version":"v1alpha1",
					"request_id":"test-req-1",
					"plan":{
						"intent":"deploy_website",
						"kind":"workflow",
						"execution_target":"go-temporal",
						"workflow":{
							"name":"Website Deployment",
							"action":"start_dynamic_pipeline",
							"input":{"project_name":"test-site","framework":"svelte"}
						},
						"artifacts":{"planner":"remote"}
					}
				}`))
			} else {
				w.WriteHeader(http.StatusBadGateway)
				_, _ = w.Write([]byte(`{
					"status":"error",
					"error":{
						"code":"VOLTAGENT_PLAN_GENERATION_FAILED",
						"message":"planner failed"
					}
				}`))
			}
			return
		}

		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	// Setup cfg
	cfg := &config.Config{
		VoltAgent: config.VoltAgentConfig{
			Enabled:             true,
			UseEmbeddedFallback: true,
			ContractVersion:     "v1alpha1",
		},
	}

	remoteClient := voltagentclient.NewClient(server.URL, 2*time.Second)
	mockTemporal := new(MockTemporalClient)

	svc := NewVoltAgentService(cfg, nil, mockTemporal, nil, remoteClient)

	// Mock temporal workflow execution for remote plan
	mockRun := new(MockWorkflowRun)
	mockRun.On("GetID").Return("wf-remote-123")
	mockRun.On("GetRunID").Return("run-remote-123")

	// Mock temporal workflow execution for fallback plan
	mockRunFallback := new(MockWorkflowRun)
	mockRunFallback.On("GetID").Return("wf-fallback-456")
	mockRunFallback.On("GetRunID").Return("run-fallback-456")

	mockTemporal.On("ExecuteWorkflow", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(mockRun, nil).Once()

	ctx := context.Background()
	meta := RequestMetadata{
		RequestID: "req-e2e-1",
		Source:    "test-suite",
	}

	t.Run("1. Initial Health Check - Remote UP", func(t *testing.T) {
		resp, err := svc.Health(ctx, meta)
		require.NoError(t, err)
		assert.Equal(t, "ok", resp.Status)
	})

	t.Run("2. Remote Planning Success", func(t *testing.T) {
		input := WebsiteDeploymentInput{
			ProjectName: "test-site",
			Prompt:      "A simple test site",
		}
		result, err := svc.StartWebsiteDeployment(ctx, input, meta)
		require.NoError(t, err)
		assert.Equal(t, planningSourceRemote, result["planning_source"])
		assert.Equal(t, "wf-remote-123", result["workflow_id"])
	})

	t.Run("3. Remote VoltAgent Goes Down", func(t *testing.T) {
		remoteStatus = "down"
		remotePlanStatus = "error"

		_, err := svc.Health(ctx, meta)
		require.Error(t, err)
	})

	t.Run("4. Remote Planning Fails -> Embedded Fallback", func(t *testing.T) {
		mockTemporal.On("ExecuteWorkflow", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(mockRunFallback, nil).Once()

		input := WebsiteDeploymentInput{
			ProjectName: "test-site",
			Prompt:      "A simple test site",
		}
		result, err := svc.StartWebsiteDeployment(ctx, input, meta)
		require.NoError(t, err)
		assert.Equal(t, planningSourceEmbedded, result["planning_source"])
		assert.Equal(t, "wf-fallback-456", result["workflow_id"])
	})

	t.Run("5. Remote VoltAgent Recovers -> Health OK", func(t *testing.T) {
		remoteStatus = "ok"
		resp, err := svc.Health(ctx, meta)
		require.NoError(t, err)
		assert.Equal(t, "ok", resp.Status)
	})
}
