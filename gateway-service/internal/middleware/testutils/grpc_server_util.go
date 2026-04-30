package testutils

import (
	"context"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	commonv1 "github.com/ai-api-gateway/api/gen/common/v1"
	providerv1 "github.com/ai-api-gateway/api/gen/provider/v1"
)

// CreateProviderServer creates a test gRPC server for provider service
func CreateProviderServer() (providerv1.ProviderServiceClient, net.Listener, *grpc.Server, error) {
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		return nil, nil, nil, err
	}

	server := grpc.NewServer()
	providerv1.RegisterProviderServiceServer(server, &mockProviderServer{})

	go func() {
		// Server.Serve will return error when listener is closed, which is expected
		server.Serve(listener)
	}()

	conn, err := grpc.Dial(listener.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		server.Stop()
		return nil, nil, nil, err
	}

	client := providerv1.NewProviderServiceClient(conn)
	return client, listener, server, nil
}

type mockProviderServer struct {
	providerv1.UnimplementedProviderServiceServer
}

func (m *mockProviderServer) ForwardRequest(ctx context.Context, req *providerv1.ForwardRequestRequest) (*providerv1.ForwardRequestResponse, error) {
	return &providerv1.ForwardRequestResponse{
		ResponseBody: []byte(`{"choices":[{"message":{"content":"test response"}}]}`),
		TokenCounts: &commonv1.TokenCounts{
			PromptTokens:     100,
			CompletionTokens: 50,
		},
	}, nil
}

func (m *mockProviderServer) StreamRequest(req *providerv1.StreamRequestRequest, stream providerv1.ProviderService_StreamRequestServer) error {
	// Send multiple chunks with token accumulation
	chunks := []*providerv1.ProviderChunk{
		{ChunkData: []byte(`{"choices":[{"delta":{"content":"Hello"}}]}`), AccumulatedTokens: &commonv1.TokenCounts{PromptTokens: 10, CompletionTokens: 5}},
		{ChunkData: []byte(`{"choices":[{"delta":{"content":" world"}}]}`), AccumulatedTokens: &commonv1.TokenCounts{PromptTokens: 10, CompletionTokens: 10}},
		{ChunkData: []byte(`{"choices":[{"delta":{"content":"!"}}]}`), AccumulatedTokens: &commonv1.TokenCounts{PromptTokens: 10, CompletionTokens: 15}, Done: true},
	}

	for _, chunk := range chunks {
		if err := stream.Send(chunk); err != nil {
			return err
		}
	}
	return nil
}
