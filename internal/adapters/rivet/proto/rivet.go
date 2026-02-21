package proto

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/structpb"
)

type RivetServiceClient interface {
	ExecuteGraph(ctx context.Context, in *ExecuteGraphRequest, opts ...grpc.CallOption) (*ExecuteGraphResponse, error)
}

type rivetServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewRivetServiceClient(cc grpc.ClientConnInterface) RivetServiceClient {
	return &rivetServiceClient{cc}
}

func (c *rivetServiceClient) ExecuteGraph(ctx context.Context, in *ExecuteGraphRequest, opts ...grpc.CallOption) (*ExecuteGraphResponse, error) {
	// Mock implementation or just return nil since this is a stub
	return &ExecuteGraphResponse{}, nil
}

type ExecuteGraphRequest struct {
	GraphId        string
	Inputs         *structpb.Struct
	ToolbeltJson   string
	ProjectContent string
}

type ExecuteGraphResponse struct {
	Outputs *structpb.Struct
}
