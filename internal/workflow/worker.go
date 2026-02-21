package workflow

import (
	"context"
	"log"
	"os"

	"nexus-super-node-v3/internal/adapters/ai"
	"nexus-super-node-v3/internal/adapters/benthos"
	"nexus-super-node-v3/internal/adapters/git"
	"nexus-super-node-v3/internal/adapters/mlops"
	"nexus-super-node-v3/internal/core/services/mcp"
	"nexus-super-node-v3/internal/core/services/wazero"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.uber.org/fx"
)

func RegisterWorker(lc fx.Lifecycle, c client.Client, mcpSvc *mcp.MCPService, mlopsCollector *mlops.Collector, wazeroSvc *wazero.WazeroService, aiClient *ai.OpenAIClient, benthosClient *benthos.Client) {
	// Task Queue name should be consistent with what starter uses
	w := worker.New(c, "handoff-task-queue", worker.Options{})

	activities := &HandoffActivities{
		LiveKitAPIKey:    os.Getenv("LIVEKIT_API_KEY"),
		LiveKitAPISecret: os.Getenv("LIVEKIT_API_SECRET"),
	}

	cryptoActivities := &CryptoActivities{
		MCPService:     mcpSvc,
		MLOpsCollector: mlopsCollector,
	}

	websiteActivities := &WebsiteActivities{
		GitProvider:   git.NewMockGitClient("dummy-token"),
		WazeroSvc:     wazeroSvc,
		AIClient:      aiClient,
		BenthosClient: benthosClient,
	}

	// Register Dynamic Pipeline Workflow
	w.RegisterWorkflow(DynamicPipelineWorkflow)

	// Register Dynamic Activities (which wrap the legacy ones)
	dynamicActivities := &DynamicActivities{
		WebsiteActivities: websiteActivities,
		CryptoActivities:  cryptoActivities,
		HandoffActivities: activities,
	}
	w.RegisterActivity(dynamicActivities)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Println("Starting Temporal Worker...")
			go func() {
				if err := w.Run(worker.InterruptCh()); err != nil {
					log.Printf("Worker stopped with error: %v", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Println("Stopping Temporal Worker...")
			w.Stop()
			return nil
		},
	})
}
