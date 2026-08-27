-- +goose Up
ALTER TABLE artist_profiles
    DROP COLUMN avatar_url;

-- +goose Down
ALTER TABLE artist_profiles
    ADD COLUMN avatar_url TEXT;
