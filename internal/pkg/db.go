package pkg

import (
	"context"
	"wishlister/pkg/postgres"
)

type DbClient interface {
	CtxTransactor
	DbExecutor

	Ping(ctx context.Context) error
}

type DbExecutor interface {
	Exec(ctx context.Context, sql string, args ...any) (postgres.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (postgres.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) postgres.Row
}

type CtxTransactor interface {
	BeginCtxTx(ctx context.Context) (context.Context, error)
	CommitCtxTx(ctx context.Context) error
	RollbackCtxTx(ctx context.Context) error
}
