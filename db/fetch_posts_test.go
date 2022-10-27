package db

import (
	"context"
	"database/sql"
	"testing"
)

func TestFetchPosts(t *testing.T) {
	querier, tearDown := newDbTestSetup(t)
	defer tearDown()

	ctx := context.Background()
	userID, err := querier.SeedUser(ctx, SeedUserParams{Username: "jupiter"})
	if err != nil {
		t.Fatalf("error seeding user. error=(%v)", err)
	}

	postIds := make([]int64, 3)
	for i := 0; i < len(postIds); i++ {
		originalPostId := sql.NullInt64{}
		if i > 0 {
			originalPostId = sql.NullInt64{Int64: postIds[i-1], Valid: true}
		}

		post, err := querier.CreatePost(ctx, CreatePostParams{
			Content: sql.NullString{String: "Original post", Valid: true},
			UserID: userID,
			OriginalPostID: originalPostId,
		})
		if err != nil {
			t.Fatalf("error creating post. error=(%v). %v", err, originalPostId)
		}
		postIds[i] = post.ID
	}

	testCases := []struct{
		name string
		arg FetchPostsParams
		wantPosts int
	}{
		{
			name: "should fetch posts with page and size",
			wantPosts: 0,
		},
		{
			name: "should fetch 3 posts with original_post",
			arg: FetchPostsParams{
				UserID: userID,
				Page: 0,
				Size: 5,
			},
			wantPosts: len(postIds),
		},
		{
			name: "should fetch posts with start_date and end_data",
		},
		{
			name: "should fetch posts with page and size",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			posts, err := querier.FetchPosts(ctx, tc.arg)
			if err != nil {
				t.Fatalf("unexpected error=(%v)", err)
			}
			if tc.wantPosts != len(posts) {
				t.Fatalf("want %v, got %v", tc.wantPosts, len(posts))
			}
		})
	}
}