-- name: CreatePost :one
INSERT INTO posts(
    content, user_id, original_post_id
) VALUES (
    $1, $2, $3
) RETURNING *;
