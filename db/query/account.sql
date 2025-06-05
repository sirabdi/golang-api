-- name: CreateAccount :one
INSERT INTO accounts (
  username, password_hash, role, owner, balance, currency, profile_picture
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
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
  username = $2,
  password_hash = $3,
  role = $4,
  owner = $5,
  balance = $6,
  currency = $7,
  profile_picture = $8
WHERE id = $1;

-- name: DeleteAccount :exec
DELETE FROM accounts
WHERE id = $1;