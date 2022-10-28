package db

import (
	"context"
	"database/sql"
	"testing"
	"time"
)

func TestCreatePost(t *testing.T) {
	querier, tearDown := newDbTestSetup(t)
	defer tearDown()

	ctx := context.Background()
	userID, err := querier.SeedUser(ctx, "jupiter")
	if err != nil {
		t.Fatalf("error seeding user. error=(%v)", err)
	}

	originalPost, err := querier.CreatePost(ctx, CreatePostParams{
		Content: sql.NullString{String: "Original post", Valid: true},
		UserID:  userID,
	})
	if err != nil {
		t.Fatalf("error creating originalPost. error=(%v)", err)
	}

	testCases := []struct {
		name    string
		post    CreatePostParams
		wantErr error
	}{
		{
			name: "should create a post",
			post: CreatePostParams{
				Content: sql.NullString{String: "Hello, world! This is my first post", Valid: true},
				UserID:  userID,
			},
		},
		{
			name: "should create a post without content",
			post: CreatePostParams{
				Content:        sql.NullString{},
				UserID:         userID,
				OriginalPostID: sql.NullInt64{Int64: originalPost.ID, Valid: true},
			},
		},
		{
			name: "should create a post with original content",
			post: CreatePostParams{
				Content: sql.NullString{},
				UserID:  userID,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := querier.CreatePost(ctx, tc.post)
			if err != nil {
				t.Fatalf("unexpected error creating post. error=(%v)", err)
			}
			if got.Content != tc.post.Content {
				t.Fatalf("Post.Content = %v, want %v", got.Content, tc.post.Content)
			}
			if got.UserID != tc.post.UserID {
				t.Fatalf("Post.userID = %v, want %v", got.UserID, tc.post.UserID)
			}
			if got.OriginalPostID != tc.post.OriginalPostID {
				t.Fatalf("Post.OriginalPostID = %v, want %v", got.OriginalPostID, tc.post.OriginalPostID)
			}
		})
	}
}

func TestCountPosts(t *testing.T) {
	querier, tearDown := newDbTestSetup(t)
	defer tearDown()

	ctx := context.Background()
	tenDaysAgo := time.Now().AddDate(0, 0, -10)

	testCases := []struct{
		name string
		username string
		times int
		wantTimes int64
		date *time.Time
	} {
		{
			name: "should count 10 posts today",
			username: "10posts",
			times: 10,
			wantTimes: 10,
		},
		{
			name: "should count 0 posts today when user has old posts",
			username: "0poststoday",
			times: 10,
			wantTimes: 0,
			date: &tenDaysAgo,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			userID := buildPosts(ctx, t, querier, tc.times, tc.username, tc.date)
			
			today := time.Now()
			local := time.Local
			start := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, local)
			end := time.Date(today.Year(), today.Month(), today.Day(), 23, 59, 59, 59, local)
			
			got, err := querier.CountPosts(ctx, CountPostsParams{
				UserID: userID,
				CreatedAt: start,
				CreatedAt_2: end,
			})
			if err != nil {
				t.Fatalf("unexpected error. error=(%v)", err)
			}
			if got != tc.wantTimes {
				t.Fatalf("CountPosts() want %v, got %v", tc.wantTimes, got)
			}
		})
	}
}

func buildPosts(ctx context.Context, t *testing.T,  querier *Queries, times int, username string, date *time.Time) int64 {
	userID, err := querier.SeedUser(ctx, username)
	if err != nil {
		t.Fatalf("error seeding user. error=(%v)", err)
	}
	if date == nil {
		now := time.Now()
		date = &now
	}

	for i := 0; i < times; i++ {
		_, err := querier.SeedPost(ctx, SeedPostParams{
			Content: sql.NullString{String: "A post", Valid: true},
			UserID:  userID,
			CreatedAt: *date,
		})
		
		if err != nil {
			t.Fatalf("error creating originalPost. error=(%v)", err)
		}
	}
	

	return userID
}
