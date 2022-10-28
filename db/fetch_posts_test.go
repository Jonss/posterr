package db

import (
	"context"
	"database/sql"
	"testing"
	"time"
)

func TestFetchPosts(t *testing.T) {
	querier, tearDown := newDbTestSetup(t)
	defer tearDown()

	ctx := context.Background()
	
	testCases := []struct{
		name string
		arg FetchPostsParams
		wantPosts int
		wantHasNext bool
		wantHasPrev bool
	}{
		{
			name: "should fetch empty posts lists when user has no post and want to see its posts",
			arg: FetchPostsParams{
				UserID: userWithoutPosts(ctx, querier, t),
				IsOnlyMyPosts: true,
			},
			wantPosts: 0,
		},
		{
			name: "should fetch 3 posts with original_post",
			arg: FetchPostsParams{
				UserID: userWithThreePosts(ctx, querier, t, "three_posts"),
				Page: 0,
				Size: 5,
				IsOnlyMyPosts: true,
			},
			wantPosts: 3,
		},
		{
			name: "should fetch posts within 2022 only",
			arg: FetchPostsParams{
				UserID: userWithPostsOn2021And2022(ctx, querier, t),
				Page: 0,
				Size: 5,
				StartDate: parsedDatePtr("2022-01-01", t),
				EndDate: parsedDatePtr("2022-12-31", t),
				IsOnlyMyPosts: true,
			},
			wantPosts: 1,
		},
		{
			name: "should fetch no posts when page is 3 and size is 3",
			arg: FetchPostsParams{
				UserID: userWithThreePosts(context.Background(), querier, t, "page_1"),
				Page: 3,
				Size: 3,
				IsOnlyMyPosts: true,
			},
			wantPosts: 0,
			wantHasPrev: true,
		},
		{
			name: "should fetch posts and contain hasNext as true",
			arg: FetchPostsParams{
				UserID: userWithTenPosts(ctx, querier, t),
				Page: 0,
				Size: 5,
				IsOnlyMyPosts: true,
			},
			wantPosts: 5,
			wantHasNext: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := querier.FetchPosts(ctx, tc.arg)
			if err != nil {
				t.Fatalf("unexpected error=(%v)", err)
			}
			if tc.wantPosts != len(got.Posts) {
				t.Fatalf("len(Posts) want %v, got %v", tc.wantPosts, len(got.Posts))
			}

			if tc.wantHasNext != got.HasNext {
				t.Fatalf("HasNext want %v, got %v", tc.wantHasNext, got.HasNext)
			}
			
			if tc.wantHasPrev != got.HasPrev {
				t.Fatalf("HasNext want %v, got %v", tc.wantHasPrev, got.HasPrev)
			}
		})
	}
}

func userWithThreePosts(ctx context.Context, querier *Queries, t *testing.T, username string) int64 {
	userID, err := querier.SeedUser(ctx, username)
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
	return userID
}

func userWithoutPosts(ctx context.Context, querier *Queries, t *testing.T) int64 {
	userID, err := querier.SeedUser(ctx, "without_posts")
	if err != nil {
		t.Fatalf("error seeding user. error=(%v)", err)
	}
	return userID
}

func userWithPostsOn2021And2022(ctx context.Context, querier *Queries, t *testing.T) int64 {
	userID, err := querier.SeedUser(ctx,"p_with_date")
	if err != nil {
		t.Fatalf("error seeding user. error=(%v)", err)
	}

	postDates := []time.Time{
		parsedDate("2021-05-02", t),
		parsedDate("2021-12-31", t),
		parsedDate("2022-05-02", t),
	}
	for _, pd := range postDates {
		_, err := querier.SeedPost(ctx,SeedPostParams{
			Content: sql.NullString{String: "a post on ", Valid: true},
			UserID: userID,
			CreatedAt: pd,
		})
		
		if err != nil {
			t.Fatalf("error creating post. error=(%v).", err)
		}
	}
	return userID
}

func userWithTenPosts(ctx context.Context, querier *Queries, t *testing.T) int64 {
	userID, err := querier.SeedUser(ctx, "five_posts")
	if err != nil {
		t.Fatalf("error seeding user. error=(%v)", err)
	}

	postIds := make([]int64, 10)
	for i := 0; i < len(postIds); i++ {

		_, err := querier.CreatePost(ctx, CreatePostParams{
			Content: sql.NullString{String: "Original post", Valid: true},
			UserID: userID,
		})
		if err != nil {
			t.Fatalf("error creating post. error=(%v).", err)
		}
	}
	return userID
}