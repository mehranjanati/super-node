package ports

import (
	"context"
)

// RivetEngine defines the interface for executing Rivet workflows.
type RivetEngine interface {
	// ExecuteGraph executes a Rivet graph/workflow.
	// graphID: The ID of the graph to execute.
	// inputs: A map of inputs for the graph.
	// Returns the outputs of the graph or an error.
	ExecuteGraph(ctx context.Context, graphID string, inputs map[string]interface{}) (map[string]interface{}, error)
}
