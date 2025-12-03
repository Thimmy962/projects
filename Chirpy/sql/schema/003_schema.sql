-- +goose Up
ALTER TABLE users
ADD CONSTRAINT users_email_unique UNIQUE (email);

-- +goose Down
ALTER TABLE users
DROP CONSTRAINT IF EXISTS users_email_unique;
