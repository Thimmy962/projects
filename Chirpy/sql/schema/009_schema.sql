-- +goose up
ALTER TABLE users
ADD COLUMN is_chirpy_red BOOLEAN DEFAULT FALSE;

-- +goose down
ALTER TABLE users
DROP COLUMN IF EXISTS is_chirpy_red;
