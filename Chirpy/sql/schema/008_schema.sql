-- +goose up
ALTER TABLE users
ALTER COLUMN email SET NOT NULL;

CREATE TABLE IF NOT exists blacklist (
    token text PRIMARY KEY
);


-- +goose down
ALTER TABLE refresh_tokens
DROP COLUMN revoked_at;

ALTER TABLE refresh_tokens
DROP COLUMN expires_at;