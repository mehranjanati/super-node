package voltagentclient

import (
	"context"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPlanPropagatesMetadataAndDefaultsContractVersion(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/plan", r.URL.Path)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "application/json", r.Header.Get("Accept"))
		assert.Equal(t, "go-gateway", r.Header.Get("X-Caller-Service"))
		assert.Equal(t, "req-123", r.Header.Get("X-Request-Id"))
		assert.Equal(t, "corr-456", r.Header.Get("X-Correlation-Id"))

		var payload PlanRequest
		require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
		assert.Equal(t, DefaultContractVersion, payload.ContractVersion)
		assert.Equal(t, "deploy_website", payload.Intent)
		assert.Equal(t, "demo-site", payload.Input["project_name"])
		assert.Equal(t, "internal_tools_deploy", payload.Context["source"])

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"status":"ok",
			"contract_version":"v1alpha1",
			"request_id":"plan-req-1",
			"plan":{
				"intent":"deploy_website",
				"kind":"workflow",
				"execution_target":"go-temporal",
				"workflow":{
					"name":"Website Deployment",
					"action":"start_dynamic_pipeline",
					"input":{"project_name":"demo-site","framework":"svelte"}
				},
				"artifacts":{"planner":"remote"},
				"warnings":["using default theme"]
			}
		}`))
	}))
	defer server.Close()

	client := NewClient(server.URL, time.Second)
	resp, err := client.Plan(context.Background(), &PlanRequest{
		Intent: "deploy_website",
		Input: map[string]any{
			"project_name": "demo-site",
		},
		Context: map[string]any{
			"source": "internal_tools_deploy",
		},
	}, RequestOptions{
		RequestID:     "req-123",
		CorrelationID: "corr-456",
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, resp.Plan)
	assert.Equal(t, "ok", resp.Status)
	assert.Equal(t, "deploy_website", resp.Plan.Intent)
	assert.Equal(t, "workflow", resp.Plan.Kind)
	assert.Equal(t, "start_dynamic_pipeline", resp.Plan.Workflow.Action)
	assert.Equal(t, "svelte", resp.Plan.Workflow.Input["framework"])
}

func TestPlanReturnsAPIErrorFromErrorEnvelope(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
		_, _ = w.Write([]byte(`{
			"status":"error",
			"contract_version":"v1alpha1",
			"request_id":"remote-err-1",
			"error":{
				"code":"VOLTAGENT_PLAN_GENERATION_FAILED",
				"message":"planner failed",
				"retryable":true,
				"details":{"stage":"planner"}
			}
		}`))
	}))
	defer server.Close()

	client := NewClient(server.URL, time.Second)
	resp, err := client.Plan(context.Background(), &PlanRequest{
		Intent: "deploy_website",
		Input:  map[string]any{"project_name": "demo-site"},
	}, RequestOptions{})

	require.Nil(t, resp)
	var apiErr *APIError
	require.ErrorAs(t, err, &apiErr)
	assert.Equal(t, http.StatusBadGateway, apiErr.StatusCode)
	assert.Equal(t, "remote-err-1", apiErr.RequestID)
	assert.Equal(t, CodePlanGenerationFailed, apiErr.Code)
	assert.Equal(t, "planner failed", apiErr.Message)
	assert.True(t, apiErr.Retryable)
	assert.Equal(t, "planner", apiErr.Details["stage"])
}

func TestHealthRejectsInvalidContractResponse(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"ok","service":"","checks":{"api":"ok","planner":"ok"}}`))
	}))
	defer server.Close()

	client := NewClient(server.URL, time.Second)
	resp, err := client.Health(context.Background(), RequestOptions{})

	require.Nil(t, resp)
	var apiErr *APIError
	require.ErrorAs(t, err, &apiErr)
	assert.Equal(t, CodeBadResponse, apiErr.Code)
	assert.Contains(t, apiErr.Message, "expected contract")
}

func TestDoClassifiesTimeoutErrors(t *testing.T) {
	t.Parallel()

	client := NewClient("http://example.test", time.Second)
	client.httpClient = &http.Client{
		Timeout: time.Second,
		Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
			return nil, timeoutError{err: errors.New("simulated timeout")}
		}),
	}

	_, err := client.Health(context.Background(), RequestOptions{})
	var apiErr *APIError
	require.ErrorAs(t, err, &apiErr)
	assert.Equal(t, CodeTimeout, apiErr.Code)
	assert.True(t, apiErr.Retryable)
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

type timeoutError struct {
	err error
}

func (e timeoutError) Error() string {
	return e.err.Error()
}

func (timeoutError) Timeout() bool {
	return true
}

func (timeoutError) Temporary() bool {
	return true
}

var _ net.Error = timeoutError{}
