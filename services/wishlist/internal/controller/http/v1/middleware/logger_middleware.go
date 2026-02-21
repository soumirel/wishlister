package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid/v5"
	"github.com/soumirel/wishlister/pkg/logger"
	"go.uber.org/zap"
)

func LoggerMiddleware() gin.HandlerFunc {
	mwLogger := logger.L().
		Named("http")
	return func(ctx *gin.Context) {
		requestID := uuid.Must(uuid.NewV7())
		reqLogger := mwLogger.With(
			zap.String("method", ctx.Request.Method),
			zap.String("uri", ctx.Request.RequestURI),
			zap.String("ip", ctx.ClientIP()),
			zap.String("request_id", requestID.String()),
		)
		ctx.Request = ctx.Request.WithContext(
			logger.WithContext(ctx, reqLogger),
		)
		ctx.Next()

		status := ctx.Writer.Status()
		reqLogger = reqLogger.With(
			zap.Int("status", status),
		)
		switch {
		case status >= http.StatusInternalServerError:
			reqLogger.Error(
				"http_request_invalid",
				zap.Error(ctx.Errors.Last()),
			)
		case status >= http.StatusBadRequest:
			reqLogger.Warn(
				"http_request_abort",
				zap.Error(ctx.Errors.Last()),
			)
		default:
			reqLogger.Info(
				"http_request",
			)
		}
	}
}
