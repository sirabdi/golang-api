-- name: CreateAccount :one
INSERT INTO accounts (
  owner, balance, currency, profile_picture
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: GetAccount :one
SELECT * FROM accounts
WHERE id = $1 LIMIT 1;

-- name: ListAccounts :many
SELECT * FROM accounts
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateAccount :exec
UPDATE accounts
SET 
  owner = $2,
  balance = $3,
  currency = $4,
  profile_picture = $5
WHERE id = $1;

-- name: DeleteAccount :exec
DELETE FROM accounts
WHERE id = $1;