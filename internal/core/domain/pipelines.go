package domain

import (
	"fmt"
	"time"
)

// Standard Pipeline Definitions for Reuse
var (
	// CryptoAnalysisPipeline is the dynamic definition for AI-driven crypto trading
	CryptoAnalysisPipeline = PipelineDefinition{
		ID:      "crypto-analysis-v1",
		Name:    "AI Crypto Analysis & Trade",
		Version: "1.0.0",
		Inputs:  []string{"user_id", "time_frame"},
		Steps: []PipelineStep{
			{
				ID:           "log_start",
				ActivityName: "LogToTerminalWrapper",
				Args: map[string]interface{}{
					"message": "Starting market analysis for user: {{user_id}}",
				},
			},
			{
				ID:           "fetch_data",
				ActivityName: "FetchMarketDataActivityWrapper",
				Args:         map[string]interface{}{},
				ResultKey:    "market_data",
			},
			{
				ID:           "analyze",
				ActivityName: "AnalyzeMarketActivityWrapper",
				Args: map[string]interface{}{
					"market_data": "{{market_data}}",
				},
				ResultKey: "analysis",
			},
			{
				ID:           "log_analysis",
				ActivityName: "LogToTerminalWrapper",
				Args: map[string]interface{}{
					"message": "Analysis complete. Waiting for user approval...",
				},
			},
			{
				ID:         "wait_approval",
				WaitSignal: "approve_trade",
				ResultKey:  "approval_data",
			},
			{
				ID:           "execute_trade",
				ActivityName: "ExecuteTradeActivityWrapper",
				Args: map[string]interface{}{
					"token":  "{{analysis.top_pick}}",
					"amount": 100.0,
				},
				ResultKey: "trade_result",
			},
		},
	}

	// HumanHandoffPipeline is the dynamic definition for operator escalation
	HumanHandoffPipeline = PipelineDefinition{
		ID:      "human-handoff-v1",
		Name:    "Human Handoff Protocol",
		Version: "1.0.0",
		Inputs:  []string{"room_id", "target_user_id"},
		Steps: []PipelineStep{
			{
				ID:           "log_start",
				ActivityName: "LogToTerminalWrapper",
				Args: map[string]interface{}{
					"message": "Initiating handoff for room: {{room_id}}",
				},
			},
			{
				ID:           "send_event",
				ActivityName: "SendHandoffEventWrapper",
				Args: map[string]interface{}{
					"RoomID":       "{{room_id}}",
					"TargetUserID": "{{target_user_id}}",
				},
			},
			{
				ID:         "wait_acceptance",
				WaitSignal: "AcceptCall",
				Timeout:    "30s",
			},
		},
	}
)

// GetCryptoAnalysisPipelineID returns a unique ID for a crypto analysis run
func GetCryptoAnalysisPipelineID(userID string) string {
	return fmt.Sprintf("crypto-analysis-%s-%d", userID, time.Now().Unix())
}

// GetHandoffPipelineID returns a unique ID for a handoff run
func GetHandoffPipelineID(roomID string) string {
	return fmt.Sprintf("handoff-%s", roomID)
}
