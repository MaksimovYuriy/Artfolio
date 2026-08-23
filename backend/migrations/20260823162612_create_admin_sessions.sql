-- +goose Up
CREATE TABLE admin_sessions (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    admin_key_id BIGINT NOT NULL REFERENCES admin_keys(id) ON DELETE CASCADE,
    token_hash BYTEA NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    revoked_at TIMESTAMPTZ,
    CONSTRAINT admin_sessions_token_hash_length
        CHECK (octet_length(token_hash) = 32),
    CONSTRAINT admin_sessions_expiration
        CHECK (expires_at > created_at),
    CONSTRAINT admin_sessions_revocation
        CHECK (revoked_at IS NULL OR revoked_at >= created_at)
);

-- +goose Down
DROP TABLE IF EXISTS admin_sessions;
