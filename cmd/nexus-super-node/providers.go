package main

import (
	"context"

	"nexus-super-node-v3/internal/adapters/ai"
	"nexus-super-node-v3/internal/adapters/gateway"
	"nexus-super-node-v3/internal/adapters/openclaw"
	"nexus-super-node-v3/internal/adapters/persistence/tidb"
	"nexus-super-node-v3/internal/adapters/redpanda"
	rivetAdapter "nexus-super-node-v3/internal/adapters/rivet" // Import Gateway
	"nexus-super-node-v3/internal/config"
	"nexus-super-node-v3/internal/core/services"
	"nexus-super-node-v3/internal/core/services/chat"
	"nexus-super-node-v3/internal/core/services/finance"
	"nexus-super-node-v3/internal/core/services/social"
	"nexus-super-node-v3/internal/core/services/wazero"

	// "nexus-super-node-v3/internal/migrations"
	"nexus-super-node-v3/internal/ports"

	"go.temporal.io/sdk/client"
)

func newConfig() (*config.Config, error) {
	return config.LoadConfig()
}

// func newMigrationConfig(cfg *config.Config) *migrations.Config {
// 	return &migrations.Config{
// 		PostgresURL: cfg.Postgres.URL,
// 	}
// }

func newTiDBRepository(ctx context.Context, cfg *config.Config) (ports.UserRepository, ports.FinanceRepository, ports.AppDataRepository, ports.SocialRepository, ports.AgentRepository, error) {
	return tidb.NewTiDBRepository(ctx, cfg.TiDB.DSN)
}

func newFinanceService(repo ports.FinanceRepository) ports.FinanceService {
	return finance.NewFinanceService(repo)
}

func newRivetEngine(cfg *config.Config, router ports.MCPRouter) (ports.RivetEngine, error) {
	// Use configuration for Rivet Service URL
	return rivetAdapter.NewGRPCClient(context.Background(), cfg.Rivet.ServiceURL, router)
}

func newRedpandaClient(cfg *config.Config) (*redpanda.Client, error) {
	return redpanda.NewClient(cfg)
}

func newWebSocketHandler() *gateway.WebSocketHandler {
	return gateway.NewWebSocketHandler()
}

func newMarketConsumer(rp *redpanda.Client, ws *gateway.WebSocketHandler, cfg *config.Config, tc client.Client) *services.MarketConsumer {
	return services.NewMarketConsumer(rp, ws, cfg, tc)
}

func newOpenAIClient(cfg *config.Config) *ai.OpenAIClient {
	return ai.NewOpenAIClient(cfg.OpenAI.APIKey)
}

func newOpenClawClient(cfg *config.Config) *openclaw.Client {
	return openclaw.NewClient(cfg.OpenClaw.GatewayURL, cfg.OpenClaw.AuthSecret)
}

func newChatService(cfg *config.Config, claw *openclaw.Client) ports.ChatService {
	return chat.NewUnifiedChatService(cfg, claw)
}

func newSocialService(repo ports.SocialRepository, rp *redpanda.Client, wasm *wazero.WazeroService) ports.SocialService {
	return social.NewSocialService(repo, rp, wasm)
}
