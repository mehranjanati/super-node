package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/twmb/franz-go/pkg/kgo"
	"go.temporal.io/api/workflowservice/v1"
	"go.temporal.io/sdk/client"
)

// TestEndToEnd_ServicesSuite runs all E2E tests for Super Node services
func TestEndToEnd_ServicesSuite(t *testing.T) {
	// 1. Test Redpanda -> Temporal Flow (Core Logic)
	t.Run("Redpanda_Temporal_Integration", TestEndToEnd_MarketSignalTrigger)

	// 2. Test Hasura GraphQL Engine
	t.Run("Hasura_Connectivity", TestEndToEnd_Hasura)

	// 3. Test Matrix Synapse
	t.Run("Matrix_Connectivity", TestEndToEnd_Matrix)

	// 4. Test Rivet Service
	t.Run("Rivet_Connectivity", TestEndToEnd_Rivet)

	// 5. Test LiveKit
	t.Run("LiveKit_Connectivity", TestEndToEnd_LiveKit)

	// 6. Test TiDB
	t.Run("TiDB_Connectivity", TestEndToEnd_TiDB)
}

func TestEndToEnd_MarketSignalTrigger(t *testing.T) {
	// 1. Setup Clients
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Redpanda Client
	rpClient, err := kgo.NewClient(
		kgo.SeedBrokers("localhost:9092"),
		kgo.AllowAutoTopicCreation(),
	)
	require.NoError(t, err, "Failed to create Redpanda client")
	defer rpClient.Close()

	// Temporal Client
	temporalClient, err := client.Dial(client.Options{
		HostPort: "localhost:7233",
	})
	require.NoError(t, err, "Failed to create Temporal client")
	defer temporalClient.Close()

	// 2. Prepare Test Data
	topic := "nexus_tasks" // Must match config.yaml
	signalID := fmt.Sprintf("e2e-test-%d", time.Now().Unix())
	payload := map[string]interface{}{
		"strategy":   "BUY",
		"top_pick":   "BTC",
		"risk_level": "HIGH",
		"reasoning":  "E2E Test Trigger " + signalID,
	}
	payloadBytes, _ := json.Marshal(payload)

	// 3. Produce Message
	t.Logf("Producing message to topic %s: %s", topic, string(payloadBytes))
	record := &kgo.Record{
		Topic: topic,
		Value: payloadBytes,
	}
	result := rpClient.ProduceSync(ctx, record)
	require.NoError(t, result.FirstErr(), "Failed to produce message to Redpanda")

	// 4. Verify Temporal Workflow Execution
	t.Log("Waiting for Temporal workflow to start...")

	// Poll for workflow
	found := false
	var executionID string

	// Retry loop to find the workflow
	for i := 0; i < 10; i++ {
		// List open workflows
		resp, err := temporalClient.ListWorkflow(ctx, &workflowservice.ListWorkflowExecutionsRequest{
			Namespace: "default",
			Query:     "ExecutionStatus='Running'",
		})

		if err != nil {
			t.Logf("Error listing workflows: %v", err)
			time.Sleep(1 * time.Second)
			continue
		}

		for _, exec := range resp.Executions {
			// Check WorkflowType.
			// In internal/workflow/crypto_workflow.go, the function is CryptoWorkflow.
			if exec.Type.Name == "CryptoWorkflow" {
				// Verify it's recent (optional, but good practice)
				executionID = exec.Execution.WorkflowId
				found = true
				break
			}
		}

		if found {
			break
		}
		time.Sleep(1 * time.Second)
	}

	// 5. Assertions
	if !found {
		t.Log("⚠️  Workflow not found. Listing all running workflows for debugging:")
		// List again to show what IS running
		resp, _ := temporalClient.ListWorkflow(ctx, &workflowservice.ListWorkflowExecutionsRequest{
			Namespace: "default",
			Query:     "ExecutionStatus='Running'",
		})
		if resp != nil {
			for _, exec := range resp.Executions {
				t.Logf("- ID: %s, Type: %s, StartTime: %v", exec.Execution.WorkflowId, exec.Type.Name, exec.StartTime)
			}
		}

		t.Log("Is the Super Node running? Run: 'go run cmd/nexus-super-node/main.go cmd/nexus-super-node/providers.go' in another terminal.")
		t.Fail() // Fail but allow logs
	} else {
		t.Logf("✅ Successfully verified workflow execution: %s", executionID)
		assert.True(t, found, "Workflow should have been started by the market signal")
	}
}

func TestEndToEnd_Hasura(t *testing.T) {
	url := "http://localhost:8080/v1/version"
	t.Logf("Checking Hasura at %s", url)

	req, _ := http.NewRequest("GET", url, nil)
	// Add Admin Secret if needed, though version might be public
	// req.Header.Add("X-Hasura-Admin-Secret", "myadminsecretkey")

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)

	if err != nil {
		t.Fatalf("Failed to connect to Hasura: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	t.Logf("Hasura Version Response: %s", string(body))

	assert.Equal(t, 200, resp.StatusCode, "Hasura should return 200 OK")
}

func TestEndToEnd_Matrix(t *testing.T) {
	url := "http://localhost:8008/_matrix/client/versions"
	t.Logf("Checking Matrix Synapse at %s", url)

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)

	if err != nil {
		t.Fatalf("Failed to connect to Matrix: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	t.Logf("Matrix Versions Response: %s", string(body))

	assert.Equal(t, 200, resp.StatusCode, "Matrix should return 200 OK")
}

func TestEndToEnd_Rivet(t *testing.T) {
	address := "localhost:50051"
	t.Logf("Checking Rivet Service at %s", address)

	conn, err := net.DialTimeout("tcp", address, 5*time.Second)
	if err != nil {
		t.Fatalf("Failed to connect to Rivet Service: %v", err)
	}
	defer conn.Close()
	t.Log("✅ Successfully connected to Rivet TCP port")
}

func TestEndToEnd_LiveKit(t *testing.T) {
	// LiveKit usually runs HTTP/WebSocket on 7880
	url := "http://localhost:7880"
	t.Logf("Checking LiveKit Server at %s", url)

	// Simple TCP check first
	conn, err := net.DialTimeout("tcp", "localhost:7880", 5*time.Second)
	if err != nil {
		t.Fatalf("Failed to connect to LiveKit port 7880: %v", err)
	}
	conn.Close()
	t.Log("✅ Successfully connected to LiveKit TCP port")
}

func TestEndToEnd_TiDB(t *testing.T) {
	dsn := "root:@tcp(127.0.0.1:4001)/test"
	t.Logf("Checking TiDB at %s", dsn)

	db, err := sql.Open("mysql", dsn)
	require.NoError(t, err)
	defer db.Close()

	err = db.Ping()
	if err != nil {
		t.Fatalf("Failed to ping TiDB: %v", err)
	}
	t.Log("✅ Successfully pinged TiDB")
}
