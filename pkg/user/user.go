package user

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrUserNotFound = errors.New("error user not found")
)

type FetchUser struct {
	UserID int64
	Username string
	CreatedAt time.Time
}

func (s *service) FetchUser(ctx context.Context, username string) (*FetchUser, error) {
	user, err := s.queries.FindUserByUsername(ctx, username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &FetchUser{
		UserID: user.ID,
		Username: user.Username,
		CreatedAt: user.CreatedAt,
	}, nil
}
