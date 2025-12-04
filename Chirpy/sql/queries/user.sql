-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
    gen_random_uuid(), NOW(), NOW(), $1, $2
)
RETURNING *;


-- name: DeleteUsers :exec
DELETE FROM users;


-- name: GetUserPassword :one
SELECT hashed_password FROM users
WHERE email = $1;

-- name: GetUser :one
SELECT id, created_at, updated_at, email FROM users
WHERE email = $1 AND hashed_password = $2;