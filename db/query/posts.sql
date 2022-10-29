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

-- name: CountPosts :one
SELECT count(1) FROM posts
WHERE user_id = $1
AND created_at BETWEEN $2 AND $3;