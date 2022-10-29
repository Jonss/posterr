package db

import "context"

type AppQuerier interface {
	Querier
	FetchPosts(ctx context.Context, arg FetchPostsParams) (FetchPosts, error)
	FetchPost(ctx context.Context, postID int64) (*FetchPost, error)
}

var _ AppQuerier = (*Queries)(nil)
