package repository

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	connString = "postgres://wishlister:wishlister@localhost:5432/wishlister?sslmode=disable"
)

func InitDB() *pgxpool.Pool {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	pool, err := connect(ctx)
	if err != nil {
		log.Fatal("connection to db failed")
	}
	err = pool.Ping(ctx)
	if err != nil {
		log.Fatal("db ping failed:", err.Error())
	}
	err = applyMigrations(ctx, pool)
	if err != nil {
		log.Fatal("db migrations failed:", err.Error())
	}
	return pool
}

func connect(ctx context.Context) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, err
	}
	return pool, nil
}

func applyMigrations(ctx context.Context, db *pgxpool.Pool) error {
	query := `
	CREATE TABLE IF NOT EXISTS users(
		id text NOT NULL,
		name text NOT NULL,
		CONSTRAINT users_pk PRIMARY KEY(id)
	);

	CREATE TABLE IF NOT EXISTS wishes(
		id text NOT NULL,
		user_id text NOT NULL,
		name text NOT NULL,
		CONSTRAINT wishes_pk PRIMARY KEY(id),
		CONSTRAINT user_id_fk FOREIGN KEY(user_id) REFERENCES users(id),
		CONSTRAINT user_id_wish_id_uq UNIQUE (user_id, id)
	)
	`
	_, err := db.Exec(ctx, query)
	if err != nil {
		return err
	}
	return nil
}
