-- name: CreatePost :one
INSERT INTO posts (
 user_id,
 text
) VALUES (
 $1, $2
) RETURNING *;

-- name: GetPost :one
SELECT * FROM posts
WHERE id = $1 LIMIT 1;

-- name: ListPosts :many
SELECT * FROM posts
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdatePost :one
UPDATE posts
  set text = $2
WHERE id = $1
RETURNING *;

-- name: DeletePost :exec
DELETE FROM posts
WHERE id = $1;