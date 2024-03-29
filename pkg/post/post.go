package post

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Jonss/posterr/db"
	"github.com/Jonss/posterr/pkg/utils"
)

const (
	maxPostsDaily = 5
)

var (
	ErrMaxPostsDailyAchieved = fmt.Errorf("error user is not allowed to post more than %d times within a day", maxPostsDaily)
	ErrRepost                = errors.New("error user cannot repost an existing repost")
	ErrQuotePost             = errors.New("error user cannot quote an existing quote-post")
)

var (
	today         = time.Now()
	startOfTheDay = time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, time.UTC)
	endOfTheDay   = time.Date(today.Year(), today.Month(), today.Day(), 23, 59, 59, 59, time.UTC)
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
				Message:   utils.NullStrToPointer(p.Post.Content),
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

func (s *service) CountDailyPosts(ctx context.Context, userId int64) error {
	count, err := s.queries.CountPostsInRange(ctx, db.CountPostsInRangeParams{
		UserID:      userId,
		CreatedAt:   startOfTheDay,
		CreatedAt_2: endOfTheDay,
	})
	if err != nil {
		return err
	}
	if count >= maxPostsDaily {
		return ErrMaxPostsDailyAchieved
	}

	return nil
}

type CreatePostParams struct {
	UserID         int64
	Message        *string
	OriginalPostID *int64
}

type CreatePostResponse struct {
	ID             int64
	Content        *string
	OriginalPostID *int64
	CreatedAt      time.Time
}

func (arg CreatePostParams) getType() PostType {
	hasMessage := arg.Message != nil
	hasOriginalPost := arg.OriginalPostID != nil

	if hasOriginalPost && !hasMessage {
		return Reposting
	}
	if hasOriginalPost && hasMessage {
		return QuotePost
	}
	return Original
}

func (s *service) CreatePost(ctx context.Context, arg CreatePostParams) (*CreatePostResponse, error) {
	err := handleOriginalPost(ctx, s.queries, arg)
	if err != nil {
		return nil, err
	}

	dbPost, err := s.queries.CreatePost(ctx, db.CreatePostParams{
		UserID:         arg.UserID,
		Content:        utils.StrPtrToNullStr(arg.Message),
		OriginalPostID: utils.Int64PtrToNullInt64(arg.OriginalPostID),
	})
	if err != nil {
		return nil, err
	}

	return &CreatePostResponse{
		ID:             dbPost.ID,
		Content:        utils.NullStrToPointer(dbPost.Content),
		CreatedAt:      dbPost.CreatedAt,
		OriginalPostID: utils.NullInt64ToInt64Ptr(dbPost.OriginalPostID),
	}, nil
}

func buildOriginalPost(p *db.OriginalPost) *Post {
	if p == nil {
		return nil
	}
	return &Post{
		ID:        utils.NullInt64ToInt64(p.ID),
		Message:   utils.NullStrToPointer(p.Content),
		Username:  *utils.NullStrToPointer(p.Username),
		CreatedAt: p.CreatedAt.Time,
	}
}

func handleOriginalPost(ctx context.Context, q db.AppQuerier, arg CreatePostParams) error {
	if arg.OriginalPostID != nil {
		originalPost, err := q.FetchPost(ctx, *arg.OriginalPostID)
		if err != nil {
			return err
		}
		if originalPost != nil {
			newPostType := arg.getType()
			if newPostType == getType(*originalPost) {
				if newPostType == Reposting {
					return ErrRepost
				}
				return ErrQuotePost
			}
		}
	}
	return nil
}

func (s *service) CountPosts(ctx context.Context, userID int64) (int64, error) {
	postsCount, err := s.queries.CountPosts(ctx, userID)
	if err != nil {
		return 0, err
	}
	return postsCount, nil
}
