package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Client struct {
	*pgxpool.Pool
}

func NewClient(ctx context.Context, dbCfg DbConfig, opts ...Option) (*Client, error) {
	o := &defaultOptions

	o.db = dbCfg

	for _, opt := range opts {
		err := opt.apply(o)
		if err != nil {
			return nil, err
		}
	}

	pgxPoolConfig, err := o.toPgxPoolConfig()
	if err != nil {
		return nil, err
	}
	pgxPool, err := pgxpool.NewWithConfig(ctx, pgxPoolConfig)
	if err != nil {
		return nil, err
	}
	return &Client{
		pgxPool,
	}, nil
}
