package repository

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/soumirel/wishlister/wishlist/pkg/postgres"
)

func setupTestDB(t *testing.T) Querier {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	migrationsBytes, err := os.ReadFile("../../db/init/init.sql")
	if err != nil {
		t.Fatalf("failed to read migrations file: %v", err)
	}

	refreshScriptBytes, err := os.ReadFile("../../db/refresh_test_data.sql")
	if err != nil {
		t.Fatalf("failed to read refresh script file: %v", err)
	}

	postgresClient, err := postgres.NewClient(ctx, postgres.DbConfig{
		Host:     "localhost",
		Port:     "5432",
		Password: "wishlister",
		User:     "wishlister",
		Database: "wishlister",
	})
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}

	// Run migrations
	_, err = postgresClient.Exec(ctx, string(migrationsBytes))
	if err != nil {
		t.Fatalf("failed to run migrations: %v", err)
	}

	// Refresh test data
	_, err = postgresClient.Exec(ctx, string(refreshScriptBytes))
	if err != nil {
		t.Fatalf("failed to refresh test data: %v", err)
	}

	return postgresClient
}
