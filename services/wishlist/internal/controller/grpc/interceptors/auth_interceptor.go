package interceptors

import (
	"context"

	"github.com/soumirel/wishlister/services/wishlist/internal/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var (
	errMissingMetadata = status.Errorf(codes.InvalidArgument, "no incoming metadata in rpc context")
)

func AuthUnaryInterceptor(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errMissingMetadata
	}

	userIdMd := md.Get("user_id")
	if len(userIdMd) == 1 {
		userId := userIdMd[0]

		ctx = auth.NewCtx(ctx, auth.Auth{
			UserID: userId,
		})
	}

	resp, err := handler(ctx, req)

	return resp, err
}
