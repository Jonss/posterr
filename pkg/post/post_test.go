package post_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/Jonss/posterr/db"
	mock_db "github.com/Jonss/posterr/db/mock"
	"github.com/Jonss/posterr/pkg/post"
	"github.com/Jonss/posterr/pkg/utils"
	"github.com/golang/mock/gomock"
)

func TestFetchPosts(t *testing.T) {
	testCases := []struct {
		name        string
		buildStubs  func(querier *mock_db.MockAppQuerier)
		wantTypes   []post.PostType
		wantErr     error
		wantHasNext bool
		wantHasPrev bool
	}{
		{
			name:      "should fetch Posts and validate types",
			wantTypes: []post.PostType{post.Original, post.Reposting, post.QuotePost},
			buildStubs: func(querier *mock_db.MockAppQuerier) {
				querier.EXPECT().
					FetchPosts(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.FetchPosts{
						Posts:   fetchPosts(),
						HasNext: false,
						HasPrev: false,
					}, nil)
			},
		},
		{
			name:      "should fetch empty list",
			wantTypes: []post.PostType{post.Original, post.Reposting, post.QuotePost},
			buildStubs: func(querier *mock_db.MockAppQuerier) {
				querier.EXPECT().
					FetchPosts(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.FetchPosts{
						Posts:   []db.FetchPost{},
						HasNext: false,
						HasPrev: false,
					}, nil)
			},
		},
		{
			name:      "should have hasNext and HasPrev",
			wantTypes: []post.PostType{post.Original, post.Reposting, post.QuotePost},
			buildStubs: func(querier *mock_db.MockAppQuerier) {
				querier.EXPECT().
					FetchPosts(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.FetchPosts{
						Posts:   []db.FetchPost{},
						HasNext: true,
						HasPrev: true,
					}, nil)
			},
			wantHasNext: true,
			wantHasPrev: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			querier := mock_db.NewMockAppQuerier(ctrl)
			tc.buildStubs(querier)

			service := post.NewPostService(querier)
			got, err := service.FetchPosts(context.Background(), post.FetchPostParams{})
			if err != nil && err != tc.wantErr {
				t.Fatalf("unexpected error fetching posts. error=(%v)", err)
			}

			if tc.wantErr == nil {
				for i, got := range got.Posts {
					if tc.wantTypes[i] != got.Post.Type {
						t.Fatalf("FetchPost.Type = %v, want %v", got.Post.Type, tc.wantTypes[i])
					}
				}

				if got.HasNext != tc.wantHasNext {
					t.Fatalf("FetchPost.HasNext = %v, want %v", got.HasNext, tc.wantHasNext)
				}

				if got.HasPrev != tc.wantHasPrev {
					t.Fatalf("FetchPost.HasPrev = %v, want %v", got.HasPrev, tc.wantHasPrev)
				}
			}
		})
	}
}

func TestCreatePosts(t *testing.T) {
	testCases := []struct {
		name        string
		buildStubs  func(querier *mock_db.MockAppQuerier)
		arg         post.CreatePostParams
		isErrorWant bool
		wantErr     error
	}{
		{
			name: "should create a post",
			buildStubs: func(querier *mock_db.MockAppQuerier) {
				querier.EXPECT().
					CreatePost(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Post{
						ID:        1,
						Content:   utils.StrToNullStr("a post"),
						CreatedAt: time.Now(),
					}, nil)
			},
			arg: post.CreatePostParams{
				UserID:  1,
				Message: utils.StrToPointer("a post"),
			},
		},
		{
			name: "should return an error when original post is a repost",
			buildStubs: func(querier *mock_db.MockAppQuerier) {
				originalPost := &db.OriginalPost{
					ID:      utils.Int64ToNullInt64(1),
					Content: utils.StrToNullStr("this is an original post"),
				}

				querier.EXPECT().
					FetchPost(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&db.FetchPost{
						Post: db.DetailedPost{
							Post: db.Post{
								UserID:         1,
								ID:             2,
								OriginalPostID: originalPost.ID,
							},
						},
						OriginalPost: originalPost,
					}, nil)
			},
			arg: post.CreatePostParams{
				UserID:         1,
				OriginalPostID: utils.Int64ToPtr(2),
			},
			isErrorWant: true,
			wantErr:     post.ErrRepost,
		},
		{
			name: "should return an error when original post is a quote-post",
			buildStubs: func(querier *mock_db.MockAppQuerier) {
				originalPost := &db.OriginalPost{
					ID:      utils.Int64ToNullInt64(1),
					Content: utils.StrToNullStr("this is an original post"),
				}

				querier.EXPECT().
					FetchPost(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&db.FetchPost{
						Post: db.DetailedPost{
							Post: db.Post{
								UserID:         1,
								ID:             2,
								OriginalPostID: utils.Int64ToNullInt64(1),
								Content:        utils.StrToNullStr("this is a quote post"),
							},
						},
						OriginalPost: originalPost,
					}, nil)
			},
			arg: post.CreatePostParams{
				UserID:         1,
				OriginalPostID: utils.Int64ToPtr(2),
				Message:        utils.StrToPointer("this is sparta!"),
			},
			isErrorWant: true,
			wantErr:     post.ErrQuotePost,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			querier := mock_db.NewMockAppQuerier(ctrl)
			tc.buildStubs(querier)

			service := post.NewPostService(querier)
			got, err := service.CreatePost(context.Background(), tc.arg)
			if err != nil && !tc.isErrorWant {
				t.Fatalf("unexpected error creating post. error=(%v)", err)
			}

			if tc.isErrorWant {
				if err != tc.wantErr {
					t.Errorf("Err got %v want %v", err, tc.wantErr)
				}
			} else {
				now := time.Now()
				if got.CreatedAt.After(now) {
					t.Errorf("Post.CreatedAt got %v want %v", got.CreatedAt, now)
				}

				if got.Content != nil && tc.arg.Message != nil && *got.Content != *tc.arg.Message {
					t.Errorf("Post.Content got %v want %v", *got.Content, *tc.arg.Message)
				}
			}
		})
	}
}

func buildPost(
	username string,
	content sql.NullString,
	ID int64,
	userID int64) db.FetchPost {
	return db.FetchPost{Post: db.DetailedPost{
		Post: db.Post{
			ID:        ID,
			UserID:    userID,
			Content:   content,
			CreatedAt: time.Now(),
		},
		Username: username,
	}}
}

func fetchPosts() []db.FetchPost {
	firstPost := buildPost("aemon", utils.StrToNullStr("Am I Aemon? Need to know!"), 1, 1)
	secondPost := buildPost("drogon", sql.NullString{}, 2, 2)
	thirdPost := buildPost("vyserion", utils.StrToNullStr("that's a good question"), 3, 3)

	return []db.FetchPost{
		{
			Post:         firstPost.Post,
			OriginalPost: nil,
		},
		{
			Post: secondPost.Post,
			OriginalPost: &db.OriginalPost{
				ID:             utils.Int64ToNullInt64(firstPost.Post.ID),
				Content:        firstPost.OriginalPost.Content,
				UserID:         utils.Int64ToNullInt64(firstPost.Post.UserID),
				OriginalPostID: firstPost.OriginalPost.ID,
				CreatedAt:      firstPost.OriginalPost.CreatedAt,
				Username:       utils.StrToNullStr(firstPost.Post.Username),
			},
		},
		{
			Post: thirdPost.Post,
			OriginalPost: &db.OriginalPost{
				ID:             utils.Int64ToNullInt64(firstPost.Post.ID),
				Content:        firstPost.OriginalPost.Content,
				UserID:         utils.Int64ToNullInt64(firstPost.Post.UserID),
				OriginalPostID: firstPost.OriginalPost.ID,
				CreatedAt:      firstPost.OriginalPost.CreatedAt,
				Username:       utils.StrToNullStr(firstPost.Post.Username),
			},
		},
	}
}
