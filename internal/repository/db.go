package repository

import (
	"context"
	"log"
	"time"
	"wishlister/internal/pkg"
	"wishlister/pkg/postgres"
)

const (
	connString = "postgres://wishlister:wishlister@localhost:5432/wishlister?sslmode=disable"
)

type db struct {
	pkg.DbClient
}

func InitDbClient(migrationsScript string, refreshScript string) pkg.DbClient {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	db, err := connect(ctx)
	if err != nil {
		log.Fatal("connection to db failed")
	}
	err = db.Ping(ctx)
	if err != nil {
		log.Fatal("db ping failed:", err.Error())
	}
	err = db.execScript(ctx, migrationsScript)
	if err != nil {
		log.Fatal("db migrations script failed:", err.Error())
	}
	err = db.execScript(ctx, refreshScript)
	if err != nil {
		log.Fatal("db refresh data script failed:", err.Error())
	}
	return db
}

func connect(ctx context.Context) (*db, error) {
	dbClient, err := postgres.NewPostgresClient(ctx, connString)
	if err != nil {
		return nil, err
	}
	return &db{dbClient}, nil
}

func (d *db) execScript(ctx context.Context, migrationsScript string) error {
	_, err := d.Exec(ctx, migrationsScript)
	if err != nil {
		return err
	}
	return nil
}
