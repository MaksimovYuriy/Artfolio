//go:build integration

package integrationtests

import (
	"bytes"
	"context"
	"testing"
	"time"

	sessionrepo "github.com/maksimovyuriy/artfolio/backend/internal/repo/session"
)

func TestSessionRepositoryLifecycle(t *testing.T) {
	database := openTestDatabase(t, "TRUNCATE admin_sessions, admin_keys RESTART IDENTITY CASCADE")
	ctx := context.Background()
	keyHash := bytes.Repeat([]byte{1}, 32)
	var adminKeyID int64
	if err := database.QueryRowContext(ctx,
		"INSERT INTO admin_keys (key_hash) VALUES ($1) RETURNING id",
		keyHash,
	).Scan(&adminKeyID); err != nil {
		t.Fatalf("create admin key: %v", err)
	}

	repository := sessionrepo.NewRepo(database)
	tokenHash := bytes.Repeat([]byte{2}, 32)
	if err := repository.Create(ctx, adminKeyID, tokenHash, time.Now().Add(time.Hour)); err != nil {
		t.Fatalf("create session: %v", err)
	}

	active, err := repository.ExistsActive(ctx, tokenHash, time.Now())
	if err != nil {
		t.Fatalf("check active session: %v", err)
	}
	if !active {
		t.Fatal("new session is not active")
	}

	if err := repository.Revoke(ctx, tokenHash, time.Now()); err != nil {
		t.Fatalf("revoke session: %v", err)
	}
	active, err = repository.ExistsActive(ctx, tokenHash, time.Now())
	if err != nil {
		t.Fatalf("check revoked session: %v", err)
	}
	if active {
		t.Fatal("revoked session is still active")
	}
}
