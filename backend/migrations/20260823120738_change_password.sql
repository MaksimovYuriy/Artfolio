-- +goose Up
ALTER TABLE users
    RENAME TO admin_keys;

ALTER TABLE admin_keys
    RENAME COLUMN password_hash TO key_hash;

ALTER TABLE admin_keys
    ALTER COLUMN key_hash TYPE BYTEA
        USING decode(key_hash, 'hex'),
    ADD CONSTRAINT admin_keys_key_hash_length
        CHECK (octet_length(key_hash) = 32);

CREATE UNIQUE INDEX admin_keys_only_one
    ON admin_keys ((TRUE));

-- +goose Down
DROP INDEX admin_keys_only_one;

ALTER TABLE admin_keys
    DROP CONSTRAINT admin_keys_key_hash_length,
    ALTER COLUMN key_hash TYPE VARCHAR(255)
        USING encode(key_hash, 'hex');

ALTER TABLE admin_keys
    RENAME COLUMN key_hash TO password_hash;

ALTER TABLE admin_keys
    RENAME TO users;
