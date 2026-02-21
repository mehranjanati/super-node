package wazero

import (
	"context"
	"fmt"
	"log"
	"nexus-super-node-v3/internal/ports"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

// WazeroService is the service for the Wazero runner.
type WazeroService struct {
	runtime wazero.Runtime
}

// Ensure WazeroService implements ports.WasmRunner
var _ ports.WasmRunner = (*WazeroService)(nil)

// NewWazeroService creates a new WazeroService.
func NewWazeroService() *WazeroService {
	ctx := context.Background()
	r := wazero.NewRuntime(ctx)

	// Instantiate WASI, which is required by many languages (TinyGo, Rust, etc.)
	if _, err := wasi_snapshot_preview1.Instantiate(ctx, r); err != nil {
		log.Printf("Failed to instantiate WASI: %v", err)
	}

	return &WazeroService{
		runtime: r,
	}
}

// RunModule implements ports.WasmRunner.
// It executes a function in a Wasm module provided as bytes or loaded from a store.
// For simplicity, this implementation expects the moduleBytes to be provided via context or a separate store.
func (s *WazeroService) RunModule(ctx context.Context, moduleID string, functionName string, params []uint64) ([]uint64, error) {
	// In a real implementation, we would fetch moduleBytes from a registry/store using moduleID.
	// For now, we'll assume moduleID is a placeholder and we need a way to get the bytes.
	return nil, fmt.Errorf("RunModule not fully implemented: need module store for %s", moduleID)
}

// ExecuteModule executes a function from a Wasm module.
// Currently supports simple void functions or returning basic integers.
// For complex types, we need memory manipulation (omitted for brevity).
func (s *WazeroService) ExecuteModule(ctx context.Context, moduleBytes []byte, functionName string) (uint64, error) {
	// Compile the module
	compiled, err := s.runtime.CompileModule(ctx, moduleBytes)
	if err != nil {
		return 0, fmt.Errorf("compile failed: %w", err)
	}

	// Instantiate the module
	mod, err := s.runtime.InstantiateModule(ctx, compiled, wazero.NewModuleConfig().WithStdout(nil).WithStderr(nil))
	if err != nil {
		return 0, fmt.Errorf("instantiate failed: %w", err)
	}
	defer mod.Close(ctx)

	// Export the function
	fn := mod.ExportedFunction(functionName)
	if fn == nil {
		return 0, fmt.Errorf("function %s not found", functionName)
	}

	// Call the function
	results, err := fn.Call(ctx)
	if err != nil {
		return 0, fmt.Errorf("execution failed: %w", err)
	}

	if len(results) > 0 {
		return results[0], nil
	}
	return 0, nil
}
