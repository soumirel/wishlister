package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type Tx interface {
	Querier
	Commit(context.Context) error
	Rollback(context.Context) error
}

type tx struct {
	pgx.Tx
}
