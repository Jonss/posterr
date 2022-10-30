package post

import (
	"context"

	"github.com/Jonss/posterr/db"
)

type service struct {
	queries db.AppQuerier
}

func NewPostService(q db.AppQuerier) *service {
	return &service{queries: q}
}

type Service interface {
	FetchPosts(ctx context.Context, arg FetchPostParams) (FetchPostResponse, error)
	CountDailyPosts(ctx context.Context, userID int64) error
	CreatePost(ctx context.Context, arg CreatePostParams) (*CreatePostResponse, error)
	CountPosts(ctx context.Context, userID int64) (int64, error)
}

var _ Service = (*service)(nil)
