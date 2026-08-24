-- +goose Up
CREATE TABLE artworks (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	title VARCHAR(256) NOT NULL,
	description TEXT NOT NULL DEFAULT '',
	technique VARCHAR(256) NOT NULL DEFAULT '',
	year SMALLINT CHECK (year BETWEEN 0 AND 9999),
	image_key TEXT NOT NULL CHECK (image_key <> ''),
	image_alt VARCHAR(256) NOT NULL DEFAULT '',
	position INTEGER NOT NULL DEFAULT 0 CHECK (position >= 0),
	is_published BOOLEAN NOT NULL DEFAULT FALSE,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX artworks_order_idx
	ON artworks (position ASC, created_at DESC, id DESC);

-- +goose Down
DROP TABLE IF EXISTS artworks;
