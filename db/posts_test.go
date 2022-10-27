package db

import (
	"context"
	"database/sql"
	"testing"
)

func TestCreatePost(t *testing.T) {
	querier, tearDown := newDbTestSetup(t)
	defer tearDown()

	ctx := context.Background()
	userID, err := querier.SeedUser(ctx, SeedUserParams{Username: "jupiter"})
	if err != nil {
		t.Fatalf("error seeding user. error=(%v)", err)
	}

	originalPost, err := querier.CreatePost(ctx, CreatePostParams{
		Content: sql.NullString{String: "Original post", Valid: true},
		UserID: userID,
	})
	if err != nil {
		t.Fatalf("error creating originalPost. error=(%v)", err)
	}

	testCases := []struct{
		name string
		post CreatePostParams
 		wantErr error
	}{
		{
			name: "should create a post",
			post: CreatePostParams{
				Content: sql.NullString{String: "Hello, world! This is my first post", Valid: true},	
				UserID: userID,
			},
		},
		{
			name: "should create a post without content",
			post: CreatePostParams{
				Content: sql.NullString{},	
				UserID: userID,
				OriginalPostID: sql.NullInt64{Int64: originalPost.ID, Valid: true},
			},
		},
		{
			name: "should create a post with original content",
			post: CreatePostParams{
				Content: sql.NullString{},	
				UserID: userID,
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