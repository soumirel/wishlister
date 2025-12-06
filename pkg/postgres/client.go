package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type txKey int

const (
	ctxTxKey txKey = iota
)

var (
	ErrNestedTx = errors.New("tx already started")
	ErrNoCtxTx  = errors.New("tx not started")
)

type PgxQueryExecutor interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
}

type postgresClient struct {
	pgxPool *pgxpool.Pool
}

func NewPostgresClient(ctx context.Context, dsn string) (*postgresClient, error) {
	dbPool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}
	return &postgresClient{
		pgxPool: dbPool,
	}, nil
}

func (c *postgresClient) activeTx(ctx context.Context) (Tx, bool) {
	tx, ok := ctx.Value(ctxTxKey).(Tx)
	if ok {
		return tx, true
	}
	return nil, false
}

func (c *postgresClient) executor(ctx context.Context) PgxQueryExecutor {
	tx, ok := c.activeTx(ctx)
	if ok {
		return tx
	}
	return c.pgxPool
}

func (c *postgresClient) BeginCtxTx(ctx context.Context) (context.Context, error) {
	_, ok := c.activeTx(ctx)
	if ok {
		return nil, ErrNestedTx
	}
	pgxTx, err := c.pgxPool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	txCtx := context.WithValue(ctx, ctxTxKey, tx{pgxTx})
	return txCtx, nil
}

func (c *postgresClient) CommitCtxTx(ctx context.Context) error {
	tx, ok := c.activeTx(ctx)
	if !ok {
		return ErrNoCtxTx
	}
	err := tx.Commit(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (c *postgresClient) RollbackCtxTx(ctx context.Context) error {
	tx, ok := c.activeTx(ctx)
	if !ok {
		return ErrNoCtxTx
	}
	err := tx.Rollback(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (c *postgresClient) Query(ctx context.Context, sql string, args ...any) (Rows, error) {
	pgxRows, err := c.executor(ctx).Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	return &rows{pgxRows}, err
}

func (c *postgresClient) QueryRow(ctx context.Context, sql string, args ...any) Row {
	pgxRow := c.executor(ctx).QueryRow(ctx, sql, args...)
	return row{pgxRow}
}

func (c *postgresClient) Exec(ctx context.Context, sql string, args ...any) (CommandTag, error) {
	pgxCommandTag, err := c.executor(ctx).Exec(ctx, sql, args...)
	if err != nil {
		return CommandTag{}, err
	}
	return CommandTag{&pgxCommandTag}, nil
}

func (c *postgresClient) Ping(ctx context.Context) error {
	err := c.pgxPool.Ping(ctx)
	if err != nil {
		return nil
	}
	return nil
}
