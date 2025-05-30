-- name: CreateBlogCategories :one
INSERT INTO blog_categories (
  name, status
) VALUES (
  $1, $2
)
RETURNING *;

-- name: GetBlogCategories :one
SELECT * FROM blog_categories
WHERE id = $1 LIMIT 1;

-- name: ListBlogCategories :many
SELECT * FROM blog_categories
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateBlogCategories :exec
UPDATE blog_categories
SET 
  name = $2,
  status = $3
WHERE id = $1;

-- name: DeleteBlogCategories :exec
DELETE FROM blog_categories
WHERE id = $1;