package interceptors

import (
	"context"

	"github.com/soumirel/wishlister/services/telegram-bot/internal/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	userIdMdKey = "user_id"
)

func AuthUnaryInterceptor(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		md = metadata.New(nil)
	} else {
		md = md.Copy()
	}

	au, ok := auth.FromCtx(ctx)
	if ok {
		md.Set(userIdMdKey, au.UserID)
	}

	ctxWithMD := metadata.NewOutgoingContext(ctx, md)
	return invoker(ctxWithMD, method, req, reply, cc, opts...)
}
