package main

import (
	"context"
	"log"
	"time"

	"go.temporal.io/sdk/client"
	"nexus-super-node-v3/internal/workflow"
)

func main() {
	// Connect to Temporal
	c, err := client.Dial(client.Options{
		HostPort: "localhost:7233",
	})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	workflowOptions := client.StartWorkflowOptions{
		ID:        "crypto-trading-workflow-001",
		TaskQueue: "handoff-task-queue",
	}

	input := workflow.CryptoWorkflowInput{
		UserID:    "user_123",
		TimeFrame: "daily",
	}

	log.Println("Starting Crypto Trading Workflow...")
	we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, workflow.CryptoWorkflow, input)
	if err != nil {
		log.Fatalln("Unable to execute workflow", err)
	}

	log.Println("Started workflow", "WorkflowID", we.GetID(), "RunID", we.GetRunID())

	// Wait for a bit to let the analysis finish
	log.Println("Waiting for AI analysis and User Notification...")
	time.Sleep(10 * time.Second)

	// Simulate User Approval
	log.Println("Simulating User Approval (Signal)...")
	err = c.SignalWorkflow(context.Background(), we.GetID(), we.GetRunID(), "approve_trade", workflow.CryptoWorkflowSignal{
		Approved: true,
		Reason:   "Looks good to me!",
	})
	if err != nil {
		log.Fatalln("Error signaling workflow", err)
	}

	// Get Result
	var result string
	err = we.Get(context.Background(), &result)
	if err != nil {
		log.Fatalln("Unable get workflow result", err)
	}

	log.Println("Workflow Result:", result)
}
