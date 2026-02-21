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
	app := fx.New(
		fx.Provide(
			context.Background,
			newConfig,
			// newMigrationConfig,
			newTiDBRepository,
			newFinanceService,
			newOpenClawClient,
			gateway.NewEchoGateway,
			services.NewMCPAggregator,
			newRedpandaClient,
			func(c *redpanda.Client) ports.EventProducer { return c },
			newRivetEngine,
			worker.NewWorkerService,
			temporal.NewClient,
			newWebSocketHandler, // Register WS Handler
			newMarketConsumer,
			newOpenAIClient,
			newChatService,   // Unified Chat (Matrix + LiveKit)
			newSocialService, // DPIN Social Feed
			func(cfg *config.Config) *benthos.Client {
				return benthos.NewClient(cfg.Benthos.APIURL)
			},
		),
		// fx.Invoke(migrations.Run),
		fx.Invoke(func(gateway *gateway.EchoGateway, ws *gateway.WebSocketHandler) {
			gateway.RegisterWebSocketRoutes(ws) // Register Routes
			go gateway.Start(context.Background())
		}),
		fx.Invoke(func(mc *services.MarketConsumer) {
			mc.Start(context.Background())
		}),

		// MLOps Collector
		fx.Provide(func() *mlops.Collector {
			return mlops.NewCollector("./data/mlops")
		}),

		// Register Worker and Tools
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
