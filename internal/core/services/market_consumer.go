package services

import (
	"context"
	"encoding/json"
	"log"
	"strings"

	"nexus-super-node-v3/internal/adapters/gateway"
	"nexus-super-node-v3/internal/adapters/redpanda"
	"nexus-super-node-v3/internal/config"
	"nexus-super-node-v3/internal/core/domain"
	"nexus-super-node-v3/internal/workflow"

	"github.com/twmb/franz-go/pkg/kgo"
	"go.temporal.io/sdk/client"
)

// MarketConsumer listens to Redpanda for market analysis results
type MarketConsumer struct {
	redpandaClient *redpanda.Client
	wsHandler      *gateway.WebSocketHandler
	config         *config.Config
	temporalClient client.Client
}

// NewMarketConsumer creates a new MarketConsumer
func NewMarketConsumer(rp *redpanda.Client, ws *gateway.WebSocketHandler, cfg *config.Config, tc client.Client) *MarketConsumer {
	return &MarketConsumer{
		redpandaClient: rp,
		wsHandler:      ws,
		config:         cfg,
		temporalClient: tc,
	}
}

// Start begins the consumption loop in a background goroutine
func (s *MarketConsumer) Start(ctx context.Context) {
	log.Println("Starting MarketConsumer...")
	go s.redpandaClient.ConsumeLoop(ctx, s.handleMessage)
}

func (s *MarketConsumer) handleMessage(ctx context.Context, record *kgo.Record) {
	// Log raw message for debugging
	log.Printf("Raw Redpanda Message: %s", string(record.Value))

	var result domain.MarketAnalysisResult
	if err := json.Unmarshal(record.Value, &result); err != nil {
		log.Printf("Error unmarshalling market analysis: %v", err)
		return
	}

	log.Printf("Received Market Analysis: %+v", result)

	// Broadcast to WebSocket (Portal SPA)
	if s.wsHandler != nil {
		s.wsHandler.Broadcast(record.Value)
	}

	// Send to Logs (Matrix notification removed)
	icon := "⚪️"
	if strings.ToUpper(result.Strategy) == "BUY" {
		icon = "🟢"
	} else if strings.ToUpper(result.Strategy) == "SELL" {
		icon = "🔴"
	}

	log.Printf("%s Market Analysis - Strategy: %s, Pick: %s, Risk: %s",
		icon, result.Strategy, result.TopPick, result.RiskLevel)

	// Trigger AI Agent Workflow if signal is significant
	if strings.ToUpper(result.Strategy) == "BUY" || strings.ToUpper(result.Strategy) == "SELL" {
		log.Println("🚨 SIGNAL DETECTED! Triggering AI Agent Workflow via Dynamic Pipeline...")

		userID := "auto-agent"
		workflowOptions := client.StartWorkflowOptions{
			ID:        domain.GetCryptoAnalysisPipelineID(userID),
			TaskQueue: "handoff-task-queue",
		}

		inputs := map[string]interface{}{
			"user_id":    userID,
			"time_frame": "1h",
		}

		// Execute Dynamic Pipeline Workflow
		we, err := s.temporalClient.ExecuteWorkflow(ctx, workflowOptions, workflow.DynamicPipelineWorkflow, domain.CryptoAnalysisPipeline, inputs)
		if err != nil {
			log.Printf("Failed to trigger dynamic workflow: %v", err)
		} else {
			log.Printf("Dynamic Workflow started: ID=%s RunID=%s", we.GetID(), we.GetRunID())
		}
	}
}
