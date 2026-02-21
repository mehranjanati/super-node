package rivet

import (
	"context"
	"fmt"
	"nexus-super-node-v3/internal/ports"
)

// RivetService is the service for the Rivet gRPC client.
type RivetService struct {
	engine ports.RivetEngine
}

// NewRivetService creates a new RivetService.
func NewRivetService(engine ports.RivetEngine) *RivetService {
	return &RivetService{
		engine: engine,
	}
}

// ExecuteGraph executes a graph on the Rivet server.
func (s *RivetService) ExecuteGraph(ctx context.Context, graphID string, inputs map[string]interface{}) (map[string]interface{}, error) {
	if s.engine == nil {
		return nil, fmt.Errorf("rivet engine is not initialized")
	}
	return s.engine.ExecuteGraph(ctx, graphID, inputs)
}
