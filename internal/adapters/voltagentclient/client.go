package voltagentclient

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"
)

const (
	DefaultContractVersion = "v1alpha1"

	CodeInvalidRequest             = "VOLTAGENT_INVALID_REQUEST"
	CodeUnsupportedContractVersion = "VOLTAGENT_UNSUPPORTED_CONTRACT_VERSION"
	CodeUnsupportedIntent          = "VOLTAGENT_UNSUPPORTED_INTENT"
	CodePlanInvalid                = "VOLTAGENT_PLAN_INVALID"
	CodePlanGenerationFailed       = "VOLTAGENT_PLAN_GENERATION_FAILED"
	CodeInternalError              = "VOLTAGENT_INTERNAL_ERROR"
	CodeUnavailable                = "VOLTAGENT_UNAVAILABLE"
	CodeTimeout                    = "VOLTAGENT_TIMEOUT"
	CodeBadResponse                = "VOLTAGENT_BAD_RESPONSE"
)

// Client is the internal HTTP client for Go -> voltagent-service calls.
type Client struct {
	baseURL       string
	callerService string
	httpClient    *http.Client
}

// RequestOptions carries tracing headers across internal service calls.
type RequestOptions struct {
	RequestID     string
	CorrelationID string
}

type HealthResponse struct {
	Status          string       `json:"status"`
	Service         string       `json:"service"`
	ContractVersion string       `json:"contract_version"`
	RequestID       string       `json:"request_id"`
	Checks          HealthChecks `json:"checks"`
}

type HealthChecks struct {
	API     string `json:"api"`
	Planner string `json:"planner"`
}

type PlanRequest struct {
	ContractVersion string         `json:"contract_version"`
	Intent          string         `json:"intent"`
	Input           map[string]any `json:"input"`
	Context         map[string]any `json:"context,omitempty"`
}

type PlanResponse struct {
	Status          string         `json:"status"`
	ContractVersion string         `json:"contract_version"`
	RequestID       string         `json:"request_id"`
	Plan            *ExecutionPlan `json:"plan"`
}

type ExecutionPlan struct {
	Intent          string         `json:"intent"`
	Kind            string         `json:"kind"`
	ExecutionTarget string         `json:"execution_target"`
	Workflow        WorkflowPlan   `json:"workflow"`
	Artifacts       map[string]any `json:"artifacts"`
	Warnings        []string       `json:"warnings"`
}

type WorkflowPlan struct {
	Name   string         `json:"name"`
	Action string         `json:"action"`
	Input  map[string]any `json:"input"`
}

type ErrorEnvelope struct {
	Status          string      `json:"status"`
	ContractVersion string      `json:"contract_version"`
	RequestID       string      `json:"request_id"`
	Error           ErrorDetail `json:"error"`
}

type ErrorDetail struct {
	Code      string         `json:"code"`
	Message   string         `json:"message"`
	Retryable bool           `json:"retryable"`
	Details   map[string]any `json:"details,omitempty"`
}

// APIError is the canonical error type returned by the VoltAgent client.
type APIError struct {
	StatusCode int
	RequestID  string
	Code       string
	Message    string
	Retryable  bool
	Details    map[string]any
	Cause      error
}

func (e *APIError) Error() string {
	if e == nil {
		return ""
	}

	base := fmt.Sprintf("%s: %s", e.Code, e.Message)
	if e.StatusCode > 0 {
		base = fmt.Sprintf("%s (status %d)", base, e.StatusCode)
	}
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", base, e.Cause)
	}
	return base
}

func (e *APIError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Cause
}

// NewClient creates a new VoltAgent client with sane defaults.
func NewClient(baseURL string, timeout time.Duration) *Client {
	if timeout <= 0 {
		timeout = 5 * time.Second
	}

	return &Client{
		baseURL:       strings.TrimRight(baseURL, "/"),
		callerService: "go-gateway",
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

// Health checks whether voltagent-service is reachable and ready.
func (c *Client) Health(ctx context.Context, opts RequestOptions) (*HealthResponse, error) {
	req, err := c.newRequest(ctx, http.MethodGet, "/health", nil, opts)
	if err != nil {
		return nil, err
	}

	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read health response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, c.decodeError(resp.StatusCode, body)
	}

	var payload HealthResponse
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, c.badResponse("health response is not valid JSON", err)
	}
	if payload.Status != "ok" || payload.Service == "" {
		return nil, c.badResponse("health response does not match the expected contract", nil)
	}

	return &payload, nil
}

// Plan requests a structured execution plan from voltagent-service.
func (c *Client) Plan(ctx context.Context, payload *PlanRequest, opts RequestOptions) (*PlanResponse, error) {
	if payload == nil {
		return nil, &APIError{
			Code:      CodeInvalidRequest,
			Message:   "plan request payload is required",
			Retryable: false,
		}
	}

	if payload.ContractVersion == "" {
		payload.ContractVersion = DefaultContractVersion
	}

	req, err := c.newRequest(ctx, http.MethodPost, "/plan", payload, opts)
	if err != nil {
		return nil, err
	}

	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read plan response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, c.decodeError(resp.StatusCode, body)
	}

	var result PlanResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, c.badResponse("plan response is not valid JSON", err)
	}
	if result.Status != "ok" || result.Plan == nil {
		return nil, c.badResponse("plan response does not match the expected contract", nil)
	}

	return &result, nil
}

func (c *Client) newRequest(ctx context.Context, method, path string, payload any, opts RequestOptions) (*http.Request, error) {
	var body io.Reader

	if payload != nil {
		data, err := json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request payload: %w", err)
		}
		body = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Caller-Service", c.callerService)
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if opts.RequestID != "" {
		req.Header.Set("X-Request-Id", opts.RequestID)
	}
	if opts.CorrelationID != "" {
		req.Header.Set("X-Correlation-Id", opts.CorrelationID)
	}

	return req, nil
}

func (c *Client) do(req *http.Request) (*http.Response, error) {
	resp, err := c.httpClient.Do(req)
	if err == nil {
		return resp, nil
	}

	if errors.Is(err, context.DeadlineExceeded) {
		return nil, &APIError{
			Code:      CodeTimeout,
			Message:   "VoltAgent request timed out",
			Retryable: true,
			Cause:     err,
		}
	}

	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return nil, &APIError{
			Code:      CodeTimeout,
			Message:   "VoltAgent request timed out",
			Retryable: true,
			Cause:     err,
		}
	}

	return nil, &APIError{
		Code:      CodeUnavailable,
		Message:   "VoltAgent service is unavailable",
		Retryable: true,
		Cause:     err,
	}
}

func (c *Client) decodeError(statusCode int, body []byte) error {
	var envelope ErrorEnvelope
	if err := json.Unmarshal(body, &envelope); err != nil {
		return c.badResponse("error response is not valid JSON", err)
	}
	if envelope.Error.Code == "" || envelope.Error.Message == "" {
		return c.badResponse("error response does not match the expected contract", nil)
	}

	return &APIError{
		StatusCode: statusCode,
		RequestID:  envelope.RequestID,
		Code:       envelope.Error.Code,
		Message:    envelope.Error.Message,
		Retryable:  envelope.Error.Retryable,
		Details:    envelope.Error.Details,
	}
}

func (c *Client) badResponse(message string, cause error) error {
	return &APIError{
		Code:      CodeBadResponse,
		Message:   message,
		Retryable: false,
		Cause:     cause,
	}
}
