package db

import "context"

type AppQuerier interface {
	Querier
	FetchPosts(ctx context.Context, arg FetchPostsParams) (FetchPosts, error)
}

var _ AppQuerier = (*Queries)(nil)
