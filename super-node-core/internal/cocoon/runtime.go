package cocoon

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

type Transaction struct {
	ID           string  `json:"transaction_id"`
	UserID       string  `json:"user_id"`
	Amount       float64 `json:"amount"`
	SourceWallet string  `json:"source_wallet"`
	Status       string  `json:"status"`
}

// Orchestrator manages the workflow execution
type Orchestrator struct {
	// Mocking DB and Event Bus connections
	DBConnection     string
	RedpandaProducer string
	WasmRuntime      wazero.Runtime
}

func NewOrchestrator() *Orchestrator {
	ctx := context.Background()
	r := wazero.NewRuntime(ctx)

	// Instantiate WASI (required for TinyGo)
	wasi_snapshot_preview1.MustInstantiate(ctx, r)

	return &Orchestrator{
		DBConnection:     "Hasura_GraphQL_Client_Mock",
		RedpandaProducer: "Redpanda_Producer_Mock",
		WasmRuntime:      r,
	}
}

// ExecuteWorkflow runs the specific financial transaction scenario
func (o *Orchestrator) ExecuteWorkflow(ctx context.Context, tx *Transaction) error {
	log.Printf("[Orchestrator] Starting workflow for Tx: %s", tx.ID)

	// Step 1: Call Hasura (Mock) - Create Pending
	o.mockHasuraCreate(tx)

	// Step 2: Call Wasm (Mock Agent: fetch_user_wallet)
	// In a real scenario, this would load another Wasm module.
	// For simplicity, we skip to the internal logic.

	// Step 3: Submit to Redpanda (Mock)
	o.mockRedpandaSubmit(tx)

	// Step 4 & 5: Cocoon Runtime - Run Risk Score & AML
	// We will run the AML check using the actual Wasm logic we defined (loading from disk)

	riskScore := o.runRiskScore(tx) // CPU Heavy
	log.Printf("[Cocoon] Risk Score: %f", riskScore)

	amlSafe, err := o.runAMLCheckWasm(ctx, tx) // GPU Heavy (Wasm)
	if err != nil {
		log.Printf("[Cocoon] AML Check Failed: %v", err)
		return err
	}
	log.Printf("[Cocoon] AML Safe: %v", amlSafe)

	// Step 6: Decision Node (Rivet Logic)
	if riskScore < 0.8 && amlSafe {
		// Approved
		tx.Status = "COMPLETED"
		o.mockHasuraUpdate(tx)
		log.Printf("[Orchestrator] Transaction APPROVED")
	} else {
		// Rejected
		tx.Status = "REJECTED"
		o.mockHasuraUpdate(tx)
		o.mockNotificationEmail(tx.UserID, "Transaction Rejected", "Your transaction was flagged.")
		log.Printf("[Orchestrator] Transaction REJECTED")
	}

	return nil
}

// runAMLCheckWasm loads and executes the compiled TinyGo binary
func (o *Orchestrator) runAMLCheckWasm(ctx context.Context, tx *Transaction) (bool, error) {
	// Load compiled Wasm (assuming it's compiled to 'agents/aml_check.wasm')
	// For this demo, we assume the binary exists. In production, we'd compile/fetch it.
	wasmBytes, err := os.ReadFile("agents/aml_check.wasm")
	if err != nil {
		// Fallback for demo if file doesn't exist
		log.Println("[Wasm] Wasm file not found, simulating SAFE result")
		return true, nil
	}

	// Create module config with host functions
	config := wazero.NewModuleConfig().WithStdout(os.Stdout).WithStderr(os.Stderr)

	// Define Host Functions (FFI)
	builder := o.WasmRuntime.NewHostModuleBuilder("env")
	builder.NewFunctionBuilder().
		WithFunc(func() {
			fmt.Println("[Host-GPU] GPU Acceleration Requested by Wasm Agent...")
		}).
		Export("request_gpu_acceleration")

	_, err = builder.Instantiate(ctx)
	if err != nil {
		return false, err
	}

	// Instantiate the Wasm module
	mod, err := o.WasmRuntime.InstantiateWithConfig(ctx, wasmBytes, config)
	if err != nil {
		return false, err
	}
	defer mod.Close(ctx)

	// Call 'allocate' to write strings to memory
	allocate := mod.ExportedFunction("allocate")
	performAnalysis := mod.ExportedFunction("perform_analysis")

	// Pass Wallet Address
	walletStr := tx.SourceWallet
	results, err := allocate.Call(ctx, uint64(len(walletStr)))
	if err != nil {
		return false, err
	}
	walletPtr := results[0]
	mod.Memory().Write(uint32(walletPtr), []byte(walletStr))

	// Pass Txn ID
	txnStr := tx.ID
	results, err = allocate.Call(ctx, uint64(len(txnStr)))
	if err != nil {
		return false, err
	}
	txnPtr := results[0]
	mod.Memory().Write(uint32(txnPtr), []byte(txnStr))

	// Call main logic
	// perform_analysis(walletPtr, walletSize, txnPtr, txnSize)
	res, err := performAnalysis.Call(ctx, walletPtr, uint64(len(walletStr)), txnPtr, uint64(len(txnStr)))
	if err != nil {
		return false, err
	}

	// 1 = Safe, 0 = Suspicious
	return res[0] == 1, nil
}

// --- Mocks ---

func (o *Orchestrator) runRiskScore(tx *Transaction) float64 {
	// Simulate ML model
	if tx.Amount > 10000 {
		return 0.9 // High risk
	}
	return 0.1 // Low risk
}

func (o *Orchestrator) mockHasuraCreate(tx *Transaction) {
	log.Println("[Hasura] Creating PENDING transaction...")
}

func (o *Orchestrator) mockHasuraUpdate(tx *Transaction) {
	log.Printf("[Hasura] Updating transaction status to %s...", tx.Status)
}

func (o *Orchestrator) mockRedpandaSubmit(tx *Transaction) {
	log.Println("[Redpanda] Task submitted to 'supernode_internal_tasks'...")
}

func (o *Orchestrator) mockNotificationEmail(user, subject, body string) {
	log.Printf("[Email] Sending to %s: %s", user, subject)
}
