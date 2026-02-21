package worker

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRivetEngine
type MockRivetEngine struct {
	mock.Mock
}

func (m *MockRivetEngine) ExecuteGraph(ctx context.Context, graphID string, inputs map[string]interface{}) (map[string]interface{}, error) {
	args := m.Called(ctx, graphID, inputs)
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func TestProcessTask_Workflow(t *testing.T) {
	// Setup
	mockRivet := new(MockRivetEngine)
	workerSvc := NewWorkerService(nil, nil, mockRivet)

	// Define expected behavior
	expectedOutputs := map[string]interface{}{"result": "success"}
	mockRivet.On("ExecuteGraph", mock.Anything, "test-graph-123", mock.MatchedBy(func(inputs map[string]interface{}) bool {
		return inputs["prompt"] == "hello rivet"
	})).Return(expectedOutputs, nil)

	// Payload (simulating Redpanda message value)
	task := WorkflowTask{
		Type:    "workflow",
		GraphID: "test-graph-123",
		Inputs: map[string]interface{}{
			"prompt": "hello rivet",
		},
	}
	payload, _ := json.Marshal(task)

	// Execute
	err := workerSvc.ProcessTask(context.Background(), payload)

	// Verify
	assert.NoError(t, err)
	mockRivet.AssertExpectations(t)
}
