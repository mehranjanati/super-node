package ports

import (
	"context"
)

// WasmRunner is the port for the Wasm compute engine.

type WasmRunner interface {
	// RunModule executes a function in a Wasm module.
	RunModule(ctx context.Context, moduleID string, functionName string, params []uint64) ([]uint64, error)
}
