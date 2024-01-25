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
WHERE user_id = $1
ORDER BY user_id
LIMIT $2
OFFSET $3;

-- name: UpdatePost :one
UPDATE posts
  set user_id = $2,
  text = $3
WHERE id = $1
RETURNING *;

-- name: DeletePost :exec
DELETE FROM posts
WHERE id = $1;