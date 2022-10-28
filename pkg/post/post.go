package post

import (
	"context"
	"time"

	"github.com/Jonss/posterr/db"
	"github.com/Jonss/posterr/pkg/strings"
)

type FetchPostParams struct {
	UserID        int64
	IsOnlyMyPosts bool
	Size          int
	Page          int
	StartDate     *time.Time
	EndDate       *time.Time
}

type PostType string

const (
	Original  PostType = "ORIGINAL"
	Reposting PostType = "REPOSTING"
	QuotePost PostType = "QUOTE_POST"
)

type Post struct {
	ID        int64
	Username  string
	Message   *string
	CreatedAt time.Time
	Type      PostType
}

func getType(fp db.FetchPost) PostType {
	hasContent := fp.Post.Content.Valid
	hasOriginalPost := fp.OriginalPost != nil
	if hasOriginalPost && !hasContent {
		return Reposting
	}
	if hasOriginalPost && hasContent {
		return QuotePost
	}
	return Original
}

type FetchPost struct {
	Post         Post
	OriginalPost *Post
}

type FetchPostResponse struct {
	Posts   []FetchPost
	HasNext bool
	HasPrev bool
}

func (s *service) FetchPosts(ctx context.Context, arg FetchPostParams) (FetchPostResponse, error) {
	dbPosts, err := s.queries.FetchPosts(ctx, db.FetchPostsParams{
		UserID:        arg.UserID,
		IsOnlyMyPosts: arg.IsOnlyMyPosts,
		Size:          arg.Size,
		Page:          arg.Page,
		StartDate:     arg.StartDate,
		EndDate:       arg.EndDate,
	})
	if err != nil {
		return FetchPostResponse{}, err
	}
	posts := dbPosts.Posts
	fetchPosts := make([]FetchPost, len(posts))
	for i, p := range posts {
		fetchPosts[i] = FetchPost{
			Post: Post{
				ID:        p.Post.ID,
				Message:   strings.NullStrToPointer(p.Post.Content),
				Username:  p.Post.Username,
				CreatedAt: p.Post.CreatedAt,
				Type:      getType(p),
			},
			OriginalPost: buildOriginalPost(p.OriginalPost),
		}
	}

	response := FetchPostResponse{
		HasNext: dbPosts.HasNext,
		HasPrev: dbPosts.HasPrev,
		Posts:   fetchPosts,
	}
	return response, nil
}

func buildOriginalPost(p *db.DetailedPost) *Post {
	if p == nil {
		return nil
	}
	return &Post{
		ID:        p.Post.ID,
		Message:   strings.NullStrToPointer(p.Post.Content),
		Username:  p.Username,
		CreatedAt: p.Post.CreatedAt,
		Type:      originalPostType(p),
	}
}

func originalPostType(op *db.DetailedPost) PostType {
	if op == nil {
		return Original
	}
	originalPost := *op
	fp := db.FetchPost{
		Post: originalPost,
	}
	return getType(fp)
}
