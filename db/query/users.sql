-- name: SeedUser :one
INSERT INTO users(username) VALUES ($1) returning id;

-- name: FindUserByUsername :one    
SELECT id, username, created_at FROM users WHERE username = $1;