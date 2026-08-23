-- +goose Up
CREATE UNIQUE INDEX artist_profiles_singleton_idx
    ON artist_profiles ((true));

-- +goose Down
DROP INDEX IF EXISTS artist_profiles_singleton_idx;
