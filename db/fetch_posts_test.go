package db

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"
)

func TestFetchPosts(t *testing.T) {
	querier, tearDown := newDbTestSetup(t)
	defer tearDown()

	ctx := context.Background()

	page3ID, page3Posts := userWithThreePosts(context.Background(), querier, t, "page_3")
	threePostsID, threePosts := userWithThreePosts(ctx, querier, t, "three_posts")
	postsInTwoYearsID, postsInTwoYears := userWithPostsOn2021And2022(ctx, querier, t)
	tenPostsID, tenPosts := userWithTenPosts(ctx, querier, t)
	
	testCases := []struct {
		name        string
		arg         FetchPostsParams
		wantPostsQty   int
		wantHasNext bool
		wantHasPrev bool
		wantPosts []Post
	}{
		{
			name: "should fetch empty posts lists when user has no post and want to see its posts",
			arg: FetchPostsParams{
				UserID:        userWithoutPosts(ctx, querier, t),
				IsOnlyMyPosts: true,
			},
			wantPostsQty: 0,
		},
		{
			name: "should fetch 3 posts with original_post",
			arg: FetchPostsParams{
				UserID:        threePostsID,
				Page:          0,
				Size:          5,
				IsOnlyMyPosts: true,
			},
			wantPostsQty: 3,
			wantPosts: threePosts,
		},
		{
			name: "should fetch posts within 2022 only",
			arg: FetchPostsParams{
				UserID:        postsInTwoYearsID,
				Page:          0,
				Size:          5,
				StartDate:     parsedDatePtr("2022-01-01", t),
				EndDate:       parsedDatePtr("2022-12-31", t),
				IsOnlyMyPosts: true,
			},
			wantPostsQty: 1,
			wantPosts: postsInTwoYears,
		},
		{
			name: "should fetch no posts when page is 3 and size is 3",
			arg: FetchPostsParams{
				UserID:        page3ID,
				Page:          3,
				Size:          3,
				IsOnlyMyPosts: true,
			},
			wantPostsQty:   0,
			wantHasPrev: true,
			wantPosts: page3Posts,
		},
		{
			name: "should fetch posts and contain hasNext as true",
			arg: FetchPostsParams{
				UserID:        tenPostsID,
				Page:          0,
				Size:          5,
				IsOnlyMyPosts: true,
			},
			wantPostsQty:   5,
			wantHasNext: true,
			wantPosts: tenPosts[5:], // subslice to check latest posts
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := querier.FetchPosts(ctx, tc.arg)
			if err != nil {
				t.Fatalf("unexpected error=(%v)", err)
			}
			if tc.wantPostsQty != len(got.Posts) {
				t.Fatalf("len(Posts) want %v, got %v", tc.wantPosts, len(got.Posts))
			}

			if tc.wantHasNext != got.HasNext {
				t.Fatalf("HasNext want %v, got %v", tc.wantHasNext, got.HasNext)
			}

			if tc.wantHasPrev != got.HasPrev {
				t.Fatalf("HasNext want %v, got %v", tc.wantHasPrev, got.HasPrev)
			}

			start := len(got.Posts)-1 // starts in the end because the query orders desc by id
			end := 0

			for start > end {
				if tc.wantPosts[start].ID != got.Posts[end].Post.ID {
					t.Errorf("Post.ID want %v, got %v", tc.wantPosts[start].ID, got.Posts[end].Post.ID)
				}
				if tc.wantPosts[start].Content != got.Posts[end].Post.Content {
					t.Errorf("Post.Content want %v, got %v", tc.wantPosts[start].Content, got.Posts[end].Post.Content)
				}
				if tc.wantPosts[start].UserID != got.Posts[end].Post.UserID {
					t.Errorf("Post.UserID want %v, got %v", tc.wantPosts[start].UserID, got.Posts[end].Post.UserID)
				}
				if tc.wantPosts[start].OriginalPostID != got.Posts[end].Post.OriginalPostID {
					t.Errorf("Post.UserID want %v, got %v", tc.wantPosts[start].OriginalPostID, got.Posts[end].Post.OriginalPostID)
				}
				if tc.wantPosts[start].OriginalPostID != got.Posts[end].Post.OriginalPostID {
					t.Errorf("Post.UserID want %v, got %v", tc.wantPosts[start].OriginalPostID, got.Posts[end].Post.OriginalPostID)
				}
				if tc.wantPosts[start].CreatedAt != got.Posts[end].Post.CreatedAt {
					t.Errorf("Post.UserID want %v, got %v", tc.wantPosts[start].OriginalPostID, got.Posts[end].Post.OriginalPostID)
				}
				start--
				end++
			}
		})
	}
}

func userWithThreePosts(ctx context.Context, querier *Queries, t *testing.T, username string) (int64, []Post) {
	userID, err := querier.SeedUser(ctx, username)
	if err != nil {
		t.Fatalf("error seeding user. error=(%v)", err)
	}

	size := 3
	posts := make([]Post, size)
	for i := 0; i < len(posts); i++ {
		originalPostId := sql.NullInt64{}
		if i > 0 {
			originalPostId = sql.NullInt64{Int64: posts[i-1].ID, Valid: true}
		}

		post, err := querier.CreatePost(ctx, CreatePostParams{
			Content:        sql.NullString{String: fmt.Sprintf("Original post %d", i), Valid: true},
			UserID:         userID,
			OriginalPostID: originalPostId,
		})
		if err != nil {
			t.Fatalf("error creating post. error=(%v). %v", err, originalPostId)
		}
		posts[i] = post
	}
	return userID, posts
}

func userWithoutPosts(ctx context.Context, querier *Queries, t *testing.T) int64 {
	userID, err := querier.SeedUser(ctx, "without_posts")
	if err != nil {
		t.Fatalf("error seeding user. error=(%v)", err)
	}
	return userID
}

func userWithPostsOn2021And2022(ctx context.Context, querier *Queries, t *testing.T) (int64, []Post) {
	userID, err := querier.SeedUser(ctx, "p_with_date")
	if err != nil {
		t.Fatalf("error seeding user. error=(%v)", err)
	}

	postDates := []time.Time{
		parsedDate("2021-05-02", t),
		parsedDate("2021-12-31", t),
		parsedDate("2022-05-02", t),
	}
	posts := make([]Post, len(postDates))
	for i, pd := range postDates {
		post, err := querier.SeedPost(ctx, SeedPostParams{
			Content:   sql.NullString{String: fmt.Sprintf("a post %d", i), Valid: true},
			UserID:    userID,
			CreatedAt: pd,
		})

		if err != nil {
			t.Fatalf("error creating post. error=(%v).", err)
		}
		posts[i] = post
	}
	return userID, posts
}

func userWithTenPosts(ctx context.Context, querier *Queries, t *testing.T) (int64, []Post) {
	userID, err := querier.SeedUser(ctx, "ten_posts")
	if err != nil {
		t.Fatalf("error seeding user. error=(%v)", err)
	}

	posts := make([]Post, 10)
	for i := 0; i < len(posts); i++ {
		post, err := querier.CreatePost(ctx, CreatePostParams{
			Content: sql.NullString{String: fmt.Sprintf("Original post %d", i), Valid: true},
			UserID:  userID,
		})
		if err != nil {
			t.Fatalf("error creating post. error=(%v).", err)
		}
		posts[i] = post
	}
	return userID, posts
}
