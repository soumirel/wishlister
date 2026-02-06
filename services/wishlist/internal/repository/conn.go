package repository

import (
	"context"

	"github.com/soumirel/wishlister/services/wishlist/pkg/postgres"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Conn interface {
	Querier
	Begin(context.Context) (Tx, error)
	Release()
}

type ConnFactory interface {
	GetConn(context.Context) (Conn, error)
}

type conn struct {
	*pgxpool.Conn
}

func (c *conn) Begin(ctx context.Context) (Tx, error) {
	pgxTx, err := c.Conn.Begin(ctx)
	if err != nil {
		return nil, err
	}
	return pgxTx, nil
}

type connFactory struct {
	postgresClient postgres.Client
}

func (f *connFactory) GetConn(ctx context.Context) (Conn, error) {
	pgxPoolConn, err := f.postgresClient.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	return &conn{
		pgxPoolConn,
	}, nil
}
