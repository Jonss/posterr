// Code generated by sqlc. DO NOT EDIT.
// source: users.sql

package db

import (
	"context"
)

const findUserByUsername = `-- name: FindUserByUsername :one
SELECT id, username, created_at FROM users WHERE username = $1
`

func (q *Queries) FindUserByUsername(ctx context.Context, username string) (User, error) {
	row := q.db.QueryRowContext(ctx, findUserByUsername, username)
	var i User
	err := row.Scan(&i.ID, &i.Username, &i.CreatedAt)
	return i, err
}

const seedUser = `-- name: SeedUser :one
INSERT INTO users(username) VALUES ($1) returning id
`

func (q *Queries) SeedUser(ctx context.Context, username string) (int64, error) {
	row := q.db.QueryRowContext(ctx, seedUser, username)
	var id int64
	err := row.Scan(&id)
	return id, err
}
