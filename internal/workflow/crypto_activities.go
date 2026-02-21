package workflow

import (
	"context"
	"fmt"
	"log"

	"nexus-super-node-v3/internal/adapters/mlops"
	"nexus-super-node-v3/internal/core/domain"
	"nexus-super-node-v3/internal/core/services/mcp"
)

// CryptoActivities holds the dependencies for the crypto trading activities
type CryptoActivities struct {
	MCPService     *mcp.MCPService
	MLOpsCollector *mlops.Collector
}

// LogThoughtActivity logs a thought step
func (a *CryptoActivities) LogThoughtActivity(ctx context.Context, roomID, agentID, thought string) error {
	log.Printf("Thought Log: [%s] %s", agentID, thought)
	return nil
}

// FetchMarketDataActivity fetches the top 10 coins
func (a *CryptoActivities) FetchMarketDataActivity(ctx context.Context) ([]map[string]interface{}, error) {
	log.Println("Activity: Fetching Market Data...")
	result, err := a.MCPService.ExecuteTool("crypto-agent", "fetch_top_coins")
	if err != nil {
		return nil, err
	}
	// Convert result to the expected type
	// In production, robust type assertions/unmarshalling needed
	return result.([]map[string]interface{}), nil
}

// AnalyzeMarketActivity uses the simulated AI to generate a strategy
func (a *CryptoActivities) AnalyzeMarketActivity(ctx context.Context, marketData []map[string]interface{}) (mcp.CryptoAnalysisResult, error) {
	log.Println("Activity: Analyzing Market with AI...")
	result, err := a.MCPService.ExecuteTool("crypto-agent", "analyze_market", marketData)
	if err != nil {
		return mcp.CryptoAnalysisResult{}, err
	}
	return result.(mcp.CryptoAnalysisResult), nil
}

// NotifyUserActivity simulates sending a notification (SIP/Matrix/Telegram)
func (a *CryptoActivities) NotifyUserActivity(ctx context.Context, userID string, analysis mcp.CryptoAnalysisResult) error {
	msg := fmt.Sprintf("🚨 STRATEGY ALERT: %s on %s. Reasoning: %s. Reply with 'approve' to execute.",
		analysis.Strategy, analysis.TopPick, analysis.Reasoning)

	log.Printf("Activity: NOTIFICATION SENT to %s. Message: %s", userID, msg)

	// Simulating SIP Call notification
	log.Printf("Activity: SIP CALL initiated to user %s to announce strategy...", userID)

	return nil
}

// CaptureTrainingSampleActivity saves the interaction for MLOps
func (a *CryptoActivities) CaptureTrainingSampleActivity(ctx context.Context, sample domain.TrainingSample) error {
	if a.MLOpsCollector == nil {
		log.Println("Warning: MLOps Collector not configured")
		return nil
	}
	return a.MLOpsCollector.LogInteraction(sample)
}

// ExecuteTradeActivity simulates a DEX swap
func (a *CryptoActivities) ExecuteTradeActivity(ctx context.Context, token string, amount float64) (string, error) {
	log.Printf("Activity: Executing Trade for %s...", token)
	result, err := a.MCPService.ExecuteTool("crypto-agent", "execute_dex_swap", "USDT", token, amount)
	if err != nil {
		return "", err
	}
	return result.(string), nil
}
