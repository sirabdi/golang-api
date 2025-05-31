-- name: CreateBlogs :one
INSERT INTO blogs (
  title, description_article, category_id, account_id, banner_image
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING *;

-- name: GetBlogs :one
SELECT * FROM blogs
WHERE id = $1 LIMIT 1;

-- name: ListBlogss :many
SELECT * FROM blogs
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateBlogs :exec
UPDATE blogs
SET 
  title = $2,
  description_article = $3,
  category_id = $4,
  account_id = $5,
  banner_image = $6
WHERE id = $1;

-- name: DeleteBlogs :exec
DELETE FROM blogs
WHERE id = $1;