package workflow

import (
	"nexus-super-node-v3/internal/core/domain"
	"nexus-super-node-v3/internal/core/services/mcp"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"go.temporal.io/sdk/testsuite"
)

type UnitTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite
}

func TestUnitTestSuite(t *testing.T) {
	suite.Run(t, new(UnitTestSuite))
}

func (s *UnitTestSuite) TestDynamicPipelineWorkflow_CryptoAnalysis() {
	env := s.NewTestWorkflowEnvironment()

	// 1. Setup MCP Service with Crypto ToolBelt
	mcpSvc := mcp.NewMCPService()
	mcpSvc.AddToolBelt(mcp.NewCryptoToolBelt())

	// 2. Setup Activities
	cryptoActivities := &CryptoActivities{
		MCPService: mcpSvc,
		// MLOpsCollector: nil (optional)
	}
	handoffActivities := &HandoffActivities{
		LiveKitAPIKey:    "test-key",
		LiveKitAPISecret: "test-secret",
	}

	// We don't need WebsiteActivities for this test, so we can pass nil or a dummy
	// dynamic_activities.go checks specific fields? No, it just calls methods on fields.
	// But DynamicActivities struct has WebsiteActivities field.
	// Let's create a dummy or nil if not used.
	// Since CryptoAnalysisPipeline doesn't use WebsiteActivities, nil is fine.

	dynamicActivities := &DynamicActivities{
		CryptoActivities:  cryptoActivities,
		HandoffActivities: handoffActivities,
		WebsiteActivities: nil,
	}

	// 3. Register Activities
	env.RegisterActivity(dynamicActivities)

	// 4. Register Delayed Callback to simulate user approval
	env.RegisterDelayedCallback(func() {
		env.SignalWorkflow("approve_trade", map[string]interface{}{
			"approved": true,
			"user_id":  "test-user-123",
		})
	}, 2*time.Second) // Wait 2 seconds (virtual time) then approve

	// 5. Execute Workflow
	inputs := map[string]interface{}{
		"user_id":    "test-user-123",
		"time_frame": "1h",
	}

	env.ExecuteWorkflow(DynamicPipelineWorkflow, domain.CryptoAnalysisPipeline, inputs)

	// 6. Assertions
	s.True(env.IsWorkflowCompleted())
	s.NoError(env.GetWorkflowError())

	var result map[string]interface{}
	env.GetWorkflowResult(&result)

	// Verify "analysis" key exists
	s.Contains(result, "analysis")
	// Note: In workflow result, it might be map[string]interface{} due to JSON serialization boundary
	// or preserved if internal. Temporal test environment usually preserves types for local execution?
	// Actually, GetWorkflowResult unmarshals into the pointer.
	// But since DynamicPipelineWorkflow returns map[string]interface{},
	// the value in map might be map[string]interface{} if it went through serialization.
	// However, TestWorkflowEnvironment runs in-memory.
	// Let's just check keys.

	// If it fails type assertion, we can check it as map
	_, ok := result["analysis"].(mcp.CryptoAnalysisResult)
	if !ok {
		// It might be a map
		analysisMap, okMap := result["analysis"].(map[string]interface{})
		s.True(okMap, "Result 'analysis' should be convertible to map or struct")
		if okMap {
			// s.NotEmpty(analysisMap["strategy"]) // Use lower case keys if map
			// Check if strategy exists (case sensitive?)
			// Usually JSON unmarshals to map[string]interface{} with original keys.
			// mcp.CryptoAnalysisResult uses json tags like "strategy".
			// So key should be "strategy".
			_, hasStrategy := analysisMap["strategy"]
			s.True(hasStrategy, "Strategy field missing in analysis map")
		}
	} else {
		// s.NotEmpty(analysis.Strategy)
	}

	// Verify "trade_result"
	s.Contains(result, "trade_result")
	tradeRes, _ := result["trade_result"].(string)
	s.Contains(tradeRes, "Successfully swapped")
}

func (s *UnitTestSuite) TestDummyDeploymentWorkflow() {
	env := s.NewTestWorkflowEnvironment()

	input := DummyDeploymentInput{
		ProjectName: "status-demo",
		Template:    "svelte",
	}

	env.ExecuteWorkflow(DummyDeploymentWorkflow, input)

	s.True(env.IsWorkflowCompleted())
	s.NoError(env.GetWorkflowError())

	var result DummyDeploymentStatus
	env.GetWorkflowResult(&result)

	s.Equal("COMPLETED", result.Status)
	s.Equal("DONE", result.CurrentStep)
	s.NotNil(result.Artifacts)
	if result.Artifacts != nil {
		s.Equal("status-demo", result.Artifacts.ProjectName)
		s.Contains(result.Artifacts.URL, "status-demo")
	}
	s.Contains(result.Logs, "Deployment finished successfully.")
}
