-- +goose Up
-- Existing rows may share the default position 0. Assign deterministic,
-- contiguous positions before enforcing uniqueness.
WITH ordered_artworks AS (
    SELECT
        id,
        ROW_NUMBER() OVER (
            ORDER BY position ASC, created_at DESC, id DESC
        ) - 1 AS normalized_position
    FROM artworks
)
UPDATE artworks
SET position = ordered_artworks.normalized_position::INTEGER
FROM ordered_artworks
WHERE artworks.id = ordered_artworks.id;

ALTER TABLE artworks
    ADD CONSTRAINT artworks_position_unique
    UNIQUE (position)
    DEFERRABLE INITIALLY DEFERRED;

DROP INDEX artworks_order_idx;

CREATE INDEX artworks_published_order_idx
    ON artworks (is_published, position, id);

-- +goose Down
DROP INDEX artworks_published_order_idx;

ALTER TABLE artworks
    DROP CONSTRAINT artworks_position_unique;

CREATE INDEX artworks_order_idx
    ON artworks (position ASC, created_at DESC, id DESC);
