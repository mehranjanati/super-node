package services

import (
	"context"
	"encoding/json"
	"testing"

	"nexus-super-node-v3/internal/core/domain"

	"github.com/stretchr/testify/mock"
	"github.com/twmb/franz-go/pkg/kgo"
	"go.temporal.io/sdk/mocks"
)

func TestMarketConsumer_HandleMessage_TriggerWorkflow(t *testing.T) {
	// Mock Temporal Client
	mockTemporal := new(mocks.Client)

	// Create MarketConsumer with mock Temporal client
	// We don't need Redpanda client or WebSocket handler for this test
	consumer := &MarketConsumer{
		temporalClient: mockTemporal,
	}

	// Prepare input message
	analysisResult := domain.MarketAnalysisResult{
		Strategy:  "BUY",
		TopPick:   "BTC",
		RiskLevel: "LOW",
	}
	value, _ := json.Marshal(analysisResult)

	record := &kgo.Record{
		Value: value,
	}

	// Expectations
	mockWorkflowRun := new(mocks.WorkflowRun)
	mockWorkflowRun.On("GetID").Return("test-workflow-id")
	mockWorkflowRun.On("GetRunID").Return("test-run-id")

	// We expect ExecuteWorkflow to be called with specific arguments
	mockTemporal.On("ExecuteWorkflow",
		mock.Anything,                 // context
		mock.Anything,                 // options (client.StartWorkflowOptions)
		mock.Anything,                 // workflow function (cannot compare funcs)
		domain.CryptoAnalysisPipeline, // pipeline definition
		mock.Anything,                 // inputs (map[string]interface{})
	).Return(mockWorkflowRun, nil)

	// Call private method directly (since we are in same package)
	consumer.handleMessage(context.Background(), record)

	// Verify
	mockTemporal.AssertExpectations(t)
}

func TestMarketConsumer_HandleMessage_NoTrigger(t *testing.T) {
	// Mock Temporal Client
	mockTemporal := new(mocks.Client)

	consumer := &MarketConsumer{
		temporalClient: mockTemporal,
	}

	// Prepare input message with HOLD strategy
	analysisResult := domain.MarketAnalysisResult{
		Strategy:  "HOLD",
		TopPick:   "BTC",
		RiskLevel: "LOW",
	}
	value, _ := json.Marshal(analysisResult)

	record := &kgo.Record{
		Value: value,
	}

	// Expectations: ExecuteWorkflow should NOT be called
	// No expectations set on mockTemporal implies no calls allowed if we used strict mock,
	// but testify mocks allow no calls by default unless configured.
	// To be sure, we can explicitly assert not called or just not set expectation.

	// Call
	consumer.handleMessage(context.Background(), record)

	// Verify
	mockTemporal.AssertNotCalled(t, "ExecuteWorkflow")
}
