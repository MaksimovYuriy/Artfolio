-- +goose Up
CREATE TABLE artist_social_links (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    handle VARCHAR(256) NOT NULL
        CHECK (btrim(handle) <> ''),
    artist_profile_id BIGINT NOT NULL
        REFERENCES artist_profiles(id) ON DELETE CASCADE,
    platform VARCHAR(32) NOT NULL
        CHECK (platform IN ('telegram', 'instagram', 'vk', 'behance')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE (artist_profile_id, platform)
);

-- +goose Down
DROP TABLE IF EXISTS artist_social_links;
