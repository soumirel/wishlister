package logger

import (
	"context"
	"sync"

	"go.uber.org/zap"
)

type loggerCtxKey struct{}

var (
	globalLogger *zap.Logger
	once         *sync.Once
)

func init() {
	once = &sync.Once{}
}

func Init(defaultFields map[string]any) *zap.Logger {
	once.Do(func() {
		cfg := zap.NewDevelopmentConfig()
		cfg.InitialFields = defaultFields
		globalLogger = zap.Must(cfg.Build()).Named("app")
		zap.ReplaceGlobals(globalLogger)
	})
	return L()
}

func L() *zap.Logger {
	return zap.L()
}

func FromContext(ctx context.Context) *zap.Logger {
	logger, ok := ctx.Value(loggerCtxKey{}).(*zap.Logger)
	if ok {
		return logger
	}
	return L()
}

func WithContext(ctx context.Context, l *zap.Logger) context.Context {
	return context.WithValue(ctx, loggerCtxKey{}, l)
}
