package user

import (
	"context"

	"github.com/Jonss/posterr/db"
)

type service struct {
	queries db.AppQuerier
}

func NewUservice(q db.AppQuerier) *service {
	return &service{queries: q}
}

type Service interface {
	FetchUser(ctx context.Context, username string) (*FetchUser, error)
}

var _ Service = (*service)(nil)
