package httpserver

import (
	"errors"
	"net/http"
	"net/url"

	"github.com/Jonss/posterr/pkg/post"
	"github.com/Jonss/posterr/pkg/strings"
)

var (
	ErrUserIDRequired = errors.New("error user_id is required and only_mine param is set")
)

var (
	defaultPageSize = 5
)

type Post struct {
	ID      int64         `json:"id"`
	Message *string       `json:"message"`
	Type    post.PostType `json:"type"`
}

type OriginalPost struct {
	ID      int64   `json:"id"`
	Message *string `json:"message"`
}

type FetchPostResponse struct {
	Post
	OriginalPost *OriginalPost `json:"originalPost"`
}

type FetchPostsResponse struct {
	FetchPostResponses []FetchPostResponse `json:"content"`
	HasNext            bool                `json:"hasNext"`
	HasPrev            bool                `json:"hasPrev"`
}

// start_date=2022-05-02&end_date=2022-05-20&page=0&size=10&only_mine=true
func (s *HttpServer) FetchPosts() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		params, err := fetchPostsParams(r.URL.Query())
		if err != nil {
			apiResponse(w, http.StatusBadRequest, err)
			return
		}
		fetchPostsResponse, err := s.services.PostService.FetchPosts(ctx, params)
		if err != nil {
			apiResponse(w, http.StatusInternalServerError, err)
			return
		}

		fetchPosts := fetchPostsResponse.Posts
		posts := make([]FetchPostResponse, len(fetchPosts))
		for i, fp := range fetchPosts {
			posts[i] = FetchPostResponse{
				Post: Post{
					ID:      fp.Post.ID,
					Message: fp.Post.Message,
					Type:    fp.Post.Type,
				},
				OriginalPost: newOriginalPost(fp.OriginalPost),
			}
		}

		response := FetchPostsResponse{
			FetchPostResponses: posts,
			HasNext:            fetchPostsResponse.HasNext,
			HasPrev:            fetchPostsResponse.HasPrev,
		}

		apiResponse(w, http.StatusOK, response)
	}
}

func fetchPostsParams(values url.Values) (post.FetchPostParams, error) {
	endDate, err := strings.ParseStringToDate(values, "end_date")
	if err != nil {
		return post.FetchPostParams{}, err
	}

	startDate, err := strings.ParseStringToDate(values, "start_date")
	if err != nil {
		return post.FetchPostParams{}, err
	}

	onlyMine, err := strings.ParseStringToBool(values, "only_mine")
	if err != nil {
		return post.FetchPostParams{}, err
	}

	size, err := strings.ParseStringToInt(values, "size")
	if err != nil {
		return post.FetchPostParams{}, err
	}
	if size == 0 {
		size = defaultPageSize
	}

	page, err := strings.ParseStringToInt(values, "page")
	if err != nil {
		return post.FetchPostParams{}, err
	}

	userID, err := strings.ParseStringToInt(values, "user_id")
	if err != nil {
		return post.FetchPostParams{}, err
	}

	if onlyMine && userID == 0 {
		return post.FetchPostParams{}, ErrUserIDRequired
	}

	return post.FetchPostParams{
		UserID:        int64(userID),
		IsOnlyMyPosts: onlyMine,
		Size:          size,
		Page:          page,
		StartDate:     startDate,
		EndDate:       endDate,
	}, nil
}

func newOriginalPost(op *post.Post) *OriginalPost {
	if op == nil {
		return nil
	}
	return &OriginalPost{
		ID:      op.ID,
		Message: op.Message,
	}
}
