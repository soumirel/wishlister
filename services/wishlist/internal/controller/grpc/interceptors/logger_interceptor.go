package interceptors

import (
	"context"
	"strings"

	"github.com/gofrs/uuid/v5"
	"github.com/soumirel/wishlister/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func LoggerUnaryInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	requestID := uuid.Must(uuid.NewV7())
	reqLogger := logger.L().With(
		zap.String("grpc.service", serviceFromFullMethod(info.FullMethod)),
		zap.String("grpc.method", methodFromFullMethod(info.FullMethod)),
		zap.String("request_id", requestID.String()),
	)
	ctx = logger.WithContext(ctx, reqLogger)

	resp, err := handler(ctx, req)

	code := codes.OK.String()
	statusCode := 0

	if st, ok := status.FromError(err); ok {
		code = st.Code().String()
		statusCode = int(st.Code())
	}

	switch {
	case statusCode >= 16:
		reqLogger.Error("grpc_request_failed",
			zap.String("grpc.code", code),
			zap.Error(err),
		)
	case statusCode >= 10:
		reqLogger.Warn("grpc_request_invalid",
			zap.String("grpc.code", code),
		)
	default:
		reqLogger.Info("grpc_request",
			zap.String("grpc.code", code),
		)
	}

	return resp, err
}

func serviceFromFullMethod(fullMethod string) string {
	parts := strings.Split(fullMethod, "/")
	if len(parts) >= 2 {
		return parts[1]
	}
	return "unknown"
}

func methodFromFullMethod(fullMethod string) string {
	parts := strings.Split(fullMethod, "/")
	if len(parts) == 3 {
		return parts[2]
	}
	return fullMethod
}
