package repository

import (
	"context"
	"log"
	"time"

	"github.com/soumirel/wishlister/wishlist/pkg/postgres"

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

func InitPostgresClient(migrationsScript string, refreshScript string) *postgresClient {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	pgClient, err := connect(ctx)
	if err != nil {
		log.Fatal("connection to db failed")
	}
	err = pgClient.Ping(ctx)
	if err != nil {
		log.Fatal("db ping failed:", err.Error())
	}
	err = execScript(ctx, pgClient, migrationsScript)
	if err != nil {
		log.Fatal("db migrations script failed:", err.Error())
	}
	err = execScript(ctx, pgClient, refreshScript)
	if err != nil {
		log.Fatal("db refresh data script failed:", err.Error())
	}
	return &postgresClient{
		pgClient,
	}
}

func connect(ctx context.Context) (*postgres.Client, error) {
	postgresClient, err := postgres.NewClient(ctx, postgres.DbConfig{
		Host:     "localhost",
		Port:     "5432",
		Password: "wishlister",
		User:     "wishlister",
		Database: "wishlister",
	})
	if err != nil {
		return nil, err
	}
	return postgresClient, nil
}

func execScript(ctx context.Context, querier Querier, migrationsScript string) error {
	_, err := querier.Exec(ctx, migrationsScript)
	if err != nil {
		return err
	}
	return nil
}
