-- name: CreateImage :one
INSERT INTO images (
 post_id,
 image_path
) VALUES (
 $1, $2
) RETURNING *;

-- name: GetImage :one
SELECT * FROM images
WHERE id = $1 LIMIT 1;

-- name: ListImages :many
SELECT * FROM images
WHERE post_id = $1
ORDER BY post_id
LIMIT $2
OFFSET $3;

-- name: UpdateImage :one
UPDATE images
  set post_id = $2,
  image_path = $3
WHERE id = $1
RETURNING *;

-- name: DeleteImage :exec
DELETE FROM posts
WHERE id = $1;