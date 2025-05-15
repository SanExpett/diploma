package mocks

import (
	"context"

	proto "github.com/SanExpett/diploma/internal/session/proto"

	"google.golang.org/grpc"
)

//go:generate mockgen -source=user_client.go -destination=user_client_mock.go -package=mocks
type UsersClient interface {
	CreateUser(ctx context.Context, in *proto.CreateUserRequest, opts ...grpc.CallOption) (*proto.CreateUserResponse, error)
	RemoveUser(ctx context.Context, in *proto.RemoveUserRequest, opts ...grpc.CallOption) (*proto.RemoveUserResponse, error)
	HasUser(ctx context.Context, in *proto.HasUserRequest, opts ...grpc.CallOption) (*proto.HasUserResponse, error)
	GetUser(ctx context.Context, in *proto.GetUserRequest, opts ...grpc.CallOption) (*proto.GetUserResponse, error)
	ChangeUserPassword(ctx context.Context, in *proto.ChangeUserPasswordRequest, opts ...grpc.CallOption) (*proto.ChangeUserPasswordResponse, error)
	ChangeUserName(ctx context.Context, in *proto.ChangeUserNameRequest, opts ...grpc.CallOption) (*proto.ChangeUserNameResponse, error)
	GetUserDataByUuid(ctx context.Context, in *proto.GetUserDataByUuidRequest, opts ...grpc.CallOption) (*proto.GetUserDataByUuidResponse, error)
	GetUserPreview(ctx context.Context, in *proto.GetUserPreviewRequest, opts ...grpc.CallOption) (*proto.GetUserPreviewResponse, error)
	ChangeUserPasswordByUuid(ctx context.Context, in *proto.ChangeUserPasswordByUuidRequest, opts ...grpc.CallOption) (*proto.ChangeUserPasswordByUuidResponse, error)
	ChangeUserNameByUuid(ctx context.Context, in *proto.ChangeUserNameByUuidRequest, opts ...grpc.CallOption) (*proto.ChangeUserNameByUuidResponse, error)
	ChangeUserAvatarByUuid(ctx context.Context, in *proto.ChangeUserAvatarByUuidRequest, opts ...grpc.CallOption) (*proto.ChangeUserAvatarByUuidResponse, error)
	HasSubscription(ctx context.Context, in *proto.HasSubscriptionRequest, opts ...grpc.CallOption) (*proto.HasSubscriptionResponse, error)
	GetSubscriptions(ctx context.Context, in *proto.GetSubscriptionsRequest, opts ...grpc.CallOption) (*proto.GetSubscriptionsResponse, error)
	PaySubscription(ctx context.Context, in *proto.PaySubscriptionRequest, opts ...grpc.CallOption) (*proto.PaySubscriptionResponse, error)
}
