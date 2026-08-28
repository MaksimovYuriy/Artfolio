//go:build integration

package integrationtests

import (
	"bytes"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/maksimovyuriy/artfolio/backend/internal/repo"
	sessionrepo "github.com/maksimovyuriy/artfolio/backend/internal/repo/session"
)

func TestSessionRepositoryLifecycle(t *testing.T) {
	database := openTestDatabase(t, "TRUNCATE admin_sessions, admin_keys RESTART IDENTITY CASCADE")
	ctx := context.Background()
	keyHash := bytes.Repeat([]byte{1}, 32)
	var actorID int64
	if err := database.QueryRowContext(ctx,
		"INSERT INTO admin_keys (key_hash) VALUES ($1) RETURNING id",
		keyHash,
	).Scan(&actorID); err != nil {
		t.Fatalf("create admin key: %v", err)
	}

	repository := sessionrepo.NewRepo(database)
	tokenHash := bytes.Repeat([]byte{2}, 32)
	sessionID, err := repository.Create(ctx, actorID, tokenHash, time.Now().Add(time.Hour))
	if err != nil {
		t.Fatalf("create session: %v", err)
	}

	authenticated, err := repository.FindActive(ctx, tokenHash, time.Now())
	if err != nil {
		t.Fatalf("find active session: %v", err)
	}
	if authenticated.ID != sessionID || authenticated.ActorID != actorID {
		t.Fatalf("authenticated session = %+v, want ID %d and actor ID %d", authenticated, sessionID, actorID)
	}

	revoked, err := repository.Revoke(ctx, tokenHash, time.Now())
	if err != nil {
		t.Fatalf("revoke session: %v", err)
	}
	if revoked != authenticated {
		t.Fatalf("revoked session = %+v, want %+v", revoked, authenticated)
	}

	if _, err := repository.FindActive(ctx, tokenHash, time.Now()); !errors.Is(err, repo.ErrNotFound) {
		t.Fatalf("find revoked session error = %v, want repository not found", err)
	}
}
