package rivet

import (
	"context"
	"encoding/json"
	"fmt"

	proto "nexus-super-node-v3/internal/adapters/rivet/proto"
	"nexus-super-node-v3/internal/ports"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/structpb"
)

// grpcClient is the implementation of the RivetEngine port.

type grpcClient struct {
	conn   *grpc.ClientConn
	router ports.MCPRouter
}

// NewGRPCClient creates a new gRPC client for the Rivet engine.

func NewGRPCClient(ctx context.Context, target string, router ports.MCPRouter) (ports.RivetEngine, error) {
	conn, err := grpc.Dial(target, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return &grpcClient{
		conn:   conn,
		router: router,
	}, nil
}

// ExecuteGraph executes a Rivet graph.

func (c *grpcClient) ExecuteGraph(ctx context.Context, graphID string, inputs map[string]interface{}) (map[string]interface{}, error) {
	toolbelt, err := c.router.GetToolBelt(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get toolbelt: %w", err)
	}

	toolbeltBytes, err := json.Marshal(toolbelt)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal toolbelt: %w", err)
	}

	pbInputs, err := structpb.NewStruct(inputs)
	if err != nil {
		return nil, fmt.Errorf("failed to convert inputs to structpb: %w", err)
	}

	client := proto.NewRivetServiceClient(c.conn)

	resp, err := client.ExecuteGraph(ctx, &proto.ExecuteGraphRequest{
		GraphId:        graphID,
		Inputs:         pbInputs,
		ToolbeltJson:   string(toolbeltBytes),
		ProjectContent: "", // Default to empty, server will use local file
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute graph: %w", err)
	}

	return resp.Outputs.AsMap(), nil
}
