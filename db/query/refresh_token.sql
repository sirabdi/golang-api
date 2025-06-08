-- name: CreateRefreshToken :one
INSERT INTO refresh_token (
  account_id, token_refresh
) VALUES (
  $1, $2
)
RETURNING *;

-- name: GetRefreshToken :one
SELECT * FROM refresh_token
WHERE account_id = $1 LIMIT 1;

-- name: DeleteRefreshToken :exec
DELETE FROM refresh_token
WHERE account_id = $1;