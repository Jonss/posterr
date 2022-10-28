-- name: CreatePost :one
INSERT INTO posts(
    content, user_id, original_post_id
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: SeedPost :one
INSERT INTO posts(
    content, user_id, original_post_id, created_at
) VALUES (
    $1, $2, $3, $4
) RETURNING *;
