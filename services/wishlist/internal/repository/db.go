package repository

import (
	"context"
	"log"
	"time"

	"github.com/soumirel/wishlister/services/wishlist/pkg/postgres"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Querier interface {
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type Pinger interface {
	Ping(ctx context.Context) error
}

type postgresClient struct {
	*postgres.Client
}

func (c *postgresClient) GetConn(ctx context.Context) (Conn, error) {
	pgxConn, err := c.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	return &conn{
		pgxConn,
	}, nil
}

func InitPostgresClient(dbCfg postgres.DbConfig, migrationsScript string) *postgresClient {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	pgClient, err := postgres.NewClient(ctx, dbCfg)
	if err != nil {
		log.Fatal("connection to db failed:", err)
	}
	err = pgClient.Ping(ctx)
	if err != nil {
		log.Fatal("db ping failed:", err.Error())
	}
	err = execScript(ctx, pgClient, migrationsScript)
	if err != nil {
		log.Fatal("db migrations script failed:", err.Error())
	}
	return &postgresClient{
		pgClient,
	}
}

func execScript(ctx context.Context, querier Querier, migrationsScript string) error {
	_, err := querier.Exec(ctx, migrationsScript)
	if err != nil {
		return err
	}
	return nil
}
