-- name: SeedUser :one
INSERT INTO users(id, username) VALUES ($1, $2) returning id;