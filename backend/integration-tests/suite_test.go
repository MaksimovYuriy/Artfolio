//go:build integration

package integrationtests

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose"
)

const defaultTestDatabaseURL = "postgres://artfolio_test:artfolio_test@127.0.0.1:55432/artfolio_test?sslmode=disable"

func openTestDatabase(t *testing.T, truncateSQL string) *sql.DB {
	t.Helper()
	databaseURL := os.Getenv("ARTFOLIO_TEST_DATABASE_URL")
	if databaseURL == "" {
		databaseURL = defaultTestDatabaseURL
	}
	database, err := sql.Open("pgx", databaseURL)
	if err != nil {
		t.Fatalf("open test database: %v", err)
	}
	t.Cleanup(func() { _ = database.Close() })

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := database.PingContext(ctx); err != nil {
		t.Fatalf("connect to test database: %v", err)
	}
	lockConnection, err := database.Conn(ctx)
	if err != nil {
		t.Fatalf("reserve integration connection: %v", err)
	}
	if _, err := lockConnection.ExecContext(ctx, "SELECT pg_advisory_lock(91724001)"); err != nil {
		_ = lockConnection.Close()
		t.Fatalf("lock integration database: %v", err)
	}
	t.Cleanup(func() {
		_, _ = lockConnection.ExecContext(context.Background(), "SELECT pg_advisory_unlock(91724001)")
		_ = lockConnection.Close()
	})

	if err := goose.SetDialect("postgres"); err != nil {
		t.Fatalf("set migration dialect: %v", err)
	}
	if err := goose.Up(database, filepath.Join("..", "migrations")); err != nil {
		t.Fatalf("apply migrations: %v", err)
	}
	if _, err := database.ExecContext(ctx, truncateSQL); err != nil {
		t.Fatalf("clean integration database: %v", err)
	}
	return database
}
