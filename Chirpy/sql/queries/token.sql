-- name: CreateRefreshToken :exec
INSERT INTO refresh_tokens (token, created_at, updated_at, user_id)
VALUES (
    $1, NOW(), NOW(), $2
);

-- name: GetUserIDFromRefreshToken :one
SELECT user_id from refresh_tokens
WHERE token = $1;

-- name: DeleteRefreshToken :exec
DELETE FROM refresh_tokens
WHERE token = $1;