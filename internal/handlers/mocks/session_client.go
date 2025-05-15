package mocks

import (
	"context"

	proto "github.com/SanExpett/diploma/internal/session/proto"
	"google.golang.org/grpc"
)

//go:generate mockgen -source=session_client.go -destination=sessions_client_mock.go -package=mocks
type SessionsClient interface {
	Add(ctx context.Context, in *proto.AddRequest, opts ...grpc.CallOption) (*proto.AddResponse, error)
	DeleteSession(ctx context.Context, in *proto.DeleteSessionRequest, opts ...grpc.CallOption) (*proto.DeleteSessionResponse, error)
	Update(ctx context.Context, in *proto.UpdateRequest, opts ...grpc.CallOption) (*proto.UpdateRequestResponse, error)
	CheckVersion(ctx context.Context, in *proto.CheckVersionRequest, opts ...grpc.CallOption) (*proto.CheckVersionResponse, error)
	GetVersion(ctx context.Context, in *proto.GetVersionRequest, opts ...grpc.CallOption) (*proto.GetVersionResponse, error)
	HasSession(ctx context.Context, in *proto.HasSessionRequest, opts ...grpc.CallOption) (*proto.HasSessionResponse, error)
	CheckAllUserSessionTokens(ctx context.Context, in *proto.CheckAllUserSessionTokensRequest, opts ...grpc.CallOption) (*proto.CheckAllUserSessionTokensResponse, error)
}
