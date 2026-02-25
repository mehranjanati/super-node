package main

import (
	"context"
	"log"

	"go.uber.org/fx"

	"nexus-super-node-v3/internal/adapters/benthos"
	"nexus-super-node-v3/internal/adapters/gateway"
	"nexus-super-node-v3/internal/adapters/mlops"
	"nexus-super-node-v3/internal/adapters/podman"
	"nexus-super-node-v3/internal/adapters/redpanda"
	"nexus-super-node-v3/internal/adapters/temporal"
	"nexus-super-node-v3/internal/config"
	"nexus-super-node-v3/internal/core/domain"
	"nexus-super-node-v3/internal/core/services"
	"nexus-super-node-v3/internal/core/services/agent"
	"nexus-super-node-v3/internal/core/services/mcp"
	"nexus-super-node-v3/internal/core/services/rivet"
	"nexus-super-node-v3/internal/core/services/voltagent"
	"nexus-super-node-v3/internal/core/services/wazero"
	"nexus-super-node-v3/internal/core/services/worker"

	// "nexus-super-node-v3/internal/migrations"
	"nexus-super-node-v3/internal/ports"
	"nexus-super-node-v3/internal/workflow"
)

func main() {
	// 1. Load config early to determine role
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	log.Printf("Starting Nexus Super Node in ROLE: %s", cfg.Role)

	// Base Modules (Common)
	baseModules := fx.Options(
		fx.Provide(
			func() *config.Config { return cfg },
			context.Background,
			// newMigrationConfig,
			newTiDBRepository,
			agent.NewAgentService,
			newFinanceService,
			newOpenClawClient,
			newRedpandaClient,
			func(c *redpanda.Client) ports.EventProducer { return c },
			temporal.NewClient,
			newOpenAIClient,
			newSocialService, // DPIN Social Feed
			func(cfg *config.Config) *benthos.Client {
				return benthos.NewClient(cfg.Benthos.APIURL)
			},
		),
		// MLOps Collector (Common for now, or move to specific roles)
		fx.Provide(func() *mlops.Collector {
			return mlops.NewCollector("./data/mlops")
		}),
	)

	// Role-Specific Modules
	var roleModules fx.Option

	switch cfg.Role {
	case "api":
		roleModules = fx.Options(
			fx.Provide(
				gateway.NewEchoGateway,
				services.NewMCPAggregator,
				newRivetEngine,
				newWebSocketHandler,
				newChatService,
			),
			fx.Invoke(func(gateway *gateway.EchoGateway, ws *gateway.WebSocketHandler) {
				gateway.RegisterWebSocketRoutes(ws)
				go gateway.Start(context.Background())
			}),
			mcp.Module,
			fx.Invoke(registerConfiguredMCPServers),
		)

	case "worker":
		roleModules = fx.Options(
			fx.Provide(
				worker.NewWorkerService,
				newRivetEngine,
			),
			fx.Invoke(workflow.RegisterWorker),
			rivet.Module,
			voltagent.Module,
			wazero.Module,
			fx.Invoke(podman.ManageInfrastructure),
		)

	case "consumer":
		roleModules = fx.Options(
			fx.Provide(
				newMarketConsumer,
				// Consumer might need WebSocketHandler if it pushes directly,
				// or it should push to Redis/Redpanda for API to pick up.
				// For now, assuming it needs it for notification (but ideally should be decoupled).
				newWebSocketHandler,
				services.NewMCPAggregator, // MarketConsumer might use MCP
			),
			fx.Invoke(func(mc *services.MarketConsumer) {
				mc.Start(context.Background())
			}),
			fx.Invoke(registerSampleToolBelt), // If consumer uses tools
		)

	case "monolith":
		fallthrough
	default:
		// All modules
		roleModules = fx.Options(
			fx.Provide(
				gateway.NewEchoGateway,
				services.NewMCPAggregator,
				newRivetEngine,
				worker.NewWorkerService,
				newWebSocketHandler,
				newMarketConsumer,
				newChatService,
			),
			fx.Invoke(func(gateway *gateway.EchoGateway, ws *gateway.WebSocketHandler) {
				gateway.RegisterWebSocketRoutes(ws)
				go gateway.Start(context.Background())
			}),
			fx.Invoke(func(mc *services.MarketConsumer) {
				mc.Start(context.Background())
			}),
			fx.Invoke(workflow.RegisterWorker),
			mcp.Module,
			rivet.Module,
			voltagent.Module,
			wazero.Module,
			fx.Invoke(registerSampleToolBelt),
			fx.Invoke(registerRedpandaServer),
			fx.Invoke(registerConfiguredMCPServers),
			fx.Invoke(podman.ManageInfrastructure),
		)
	}

	app := fx.New(
		baseModules,
		roleModules,
	)

	app.Run()

	if err := app.Err(); err != nil {
		log.Fatal(err)
	}
}

func registerRedpandaServer(router ports.MCPRouter, client *redpanda.Client) {
	tools, err := client.GetToolBelt(context.Background())
	if err != nil {
		log.Printf("failed to get redpanda toolbelt (non-critical): %v", err)
		return
	}

	server := &domain.MCPServer{
		ID:    "redpanda",
		Type:  "http",
		Tools: tools,
	}

	router.RegisterServer(server)
}

func registerSampleToolBelt(mcpService *mcp.MCPService) {
	mcpService.AddToolBelt(mcp.NewSampleToolBelt())
	mcpService.AddToolBelt(mcp.NewCryptoToolBelt())
}

func registerConfiguredMCPServers(mcpService *mcp.MCPService, cfg *config.Config) {
	for _, serverCfg := range cfg.MCP.Servers {
		// Convert config.MCPServerConfig to mcp.MCPServerConfig
		mcpCfg := mcp.MCPServerConfig{
			ID:          serverCfg.ID,
			Name:        serverCfg.Name,
			Type:        serverCfg.Type,
			Command:     serverCfg.Command,
			Args:        serverCfg.Args,
			URL:         serverCfg.URL,
			Environment: serverCfg.Environment,
		}

		if err := mcpService.RegisterDynamicServer(mcpCfg); err != nil {
			log.Printf("Failed to register configured MCP server %s: %v", serverCfg.Name, err)
		} else {
			log.Printf("Registered configured MCP server: %s", serverCfg.Name)
		}
	}
}
