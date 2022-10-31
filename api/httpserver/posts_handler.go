package httpserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/Jonss/posterr/pkg/post"
	"github.com/Jonss/posterr/pkg/utils"
)

var (
	ErrUserIDRequired = errors.New("error user_id is required and only_mine param is set")
)

var (
	defaultPageSize                   = 5
	messageMaxLength                  = 777
	messageLengthErrorMessage         = fmt.Sprintf("message must be at maximum %d characters in length", messageMaxLength)
	messageOrOriginalPostIdIsRequired = "error message or originalPostId is required"
)

type Post struct {
	ID       int64         `json:"id"`
	Message  *string       `json:"message"`
	Type     post.PostType `json:"type"`
	Username string        `json:"username"`
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

func (s *HttpServer) FetchPosts() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		params, err := fetchPostsParams(r.URL.Query())
		if err != nil {
			apiResponse(w, http.StatusBadRequest, NewErrorResponses(ErrorResponse{Message: err.Error()}))
			return
		}
		fetchPostsResponse, err := s.services.PostService.FetchPosts(ctx, params)
		if err != nil {
			apiResponse(w, http.StatusInternalServerError, NewErrorResponses(ErrorResponse{Message: "unexpected error"}))
			return
		}

		fetchPosts := fetchPostsResponse.Posts
		posts := make([]FetchPostResponse, len(fetchPosts))
		for i, fp := range fetchPosts {
			posts[i] = FetchPostResponse{
				Post: Post{
					ID:       fp.Post.ID,
					Message:  fp.Post.Message,
					Type:     fp.Post.Type,
					Username: fp.Post.Username,
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

type CreatePostRequest struct {
	UserID         int64   `json:"userId" validate:"required"`
	Message        *string `json:"message"`
	OriginalPostID *int64  `json:"originalPostId"`
}

type CreatePostResponse struct {
	ID             int64     `json:"id"`
	Message        *string   `json:"message"`
	OriginalPostID *int64    `json:"originalPostId"`
	CreatedAt      time.Time `json:"createdAt"`
}

func (s *HttpServer) CreatePost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var req CreatePostRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			apiResponse(w, http.StatusBadRequest, nil)
			return
		}

		err := s.restValidator.Validator.Struct(req)
		if err != nil {
			validateRequestBody(err, w, s.restValidator.Translator)
			return
		}

		if req.Message == nil && req.OriginalPostID == nil {
			apiResponse(w, http.StatusBadRequest, NewErrorResponses(NewErrorResponse(messageOrOriginalPostIdIsRequired)))
			return
		}

		if req.Message != nil && len(*req.Message) > messageMaxLength {
			apiResponse(w, http.StatusBadRequest, NewErrorResponses(NewErrorResponse(messageLengthErrorMessage)))
			return
		}

		// checks if user can post in the day
		err = s.services.PostService.CountDailyPosts(ctx, req.UserID)
		if err != nil {
			apiResponse(w, http.StatusUnprocessableEntity, NewErrorResponses(NewErrorResponse(err.Error())))
			return
		}

		cpResponse, err := s.services.PostService.CreatePost(ctx, post.CreatePostParams{
			Message:        req.Message,
			UserID:         int64(req.UserID),
			OriginalPostID: req.OriginalPostID,
		})

		if err != nil {
			// TODO: handle errors when original_post and user does not exists
			if err == post.ErrQuotePost || err == post.ErrRepost {
				apiResponse(w, http.StatusUnprocessableEntity, NewErrorResponses(NewErrorResponse(err.Error())))
				return
			}

			apiResponse(w, http.StatusInternalServerError, NewErrorResponses(NewErrorResponse(err.Error())))
			return
		}

		response := CreatePostResponse{
			ID:             cpResponse.ID,
			CreatedAt:      cpResponse.CreatedAt,
			Message:        cpResponse.Content,
			OriginalPostID: cpResponse.OriginalPostID,
		}

		apiResponse(w, http.StatusCreated, response)
	}
}

func fetchPostsParams(values url.Values) (post.FetchPostParams, error) {
	endDate, err := utils.ParseStringToDate(values, "end_date")
	if err != nil {
		return post.FetchPostParams{}, err
	}

	startDate, err := utils.ParseStringToDate(values, "start_date")
	if err != nil {
		return post.FetchPostParams{}, err
	}

	onlyMine, err := utils.ParseStringToBool(values, "only_mine")
	if err != nil {
		return post.FetchPostParams{}, err
	}

	size, err := utils.ParseStringToInt(values, "size")
	if err != nil {
		return post.FetchPostParams{}, err
	}
	if size == 0 {
		size = defaultPageSize
	}

	page, err := utils.ParseStringToInt(values, "page")
	if err != nil {
		return post.FetchPostParams{}, err
	}

	userID, err := utils.ParseStringToInt(values, "user_id")
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
