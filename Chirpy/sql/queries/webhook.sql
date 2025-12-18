-- name: UpdateRed :exec
UPDATE users set is_chirpy_red = TRUE
WHERE id = $1;
