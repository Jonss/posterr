package httpserver

import (
	"net/http"

	"github.com/Jonss/posterr/pkg/user"
	"github.com/Jonss/posterr/pkg/utils"
	"github.com/gorilla/mux"
)

// TODO:
// test
type FetchUserResponse struct {
	Username   string `json:"username"`
	DateJoined string `json:"dateJoined"`
	PostsCount int64  `json:"postsCount"`
}

func (s *HttpServer) FetchUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		username := mux.Vars(r)["username"]

		fetchUser, err := s.services.UserService.FetchUser(ctx, username)
		if err != nil {
			if err == user.ErrUserNotFound {
				apiResponse(w, http.StatusNotFound, NewErrorResponses(ErrorResponse{Message: err.Error()}))
				return
			}
			apiResponse(w, http.StatusInternalServerError, NewErrorResponses(ErrorResponse{Message: "unexpected error"}))
			return
		}

		postsCount, err := s.services.PostService.CountPosts(ctx, fetchUser.UserID)
		if err != nil {
			apiResponse(w, http.StatusInternalServerError, NewErrorResponses(ErrorResponse{Message: "unexpected error"}))
			return
		}

		response := FetchUserResponse{
			Username:   fetchUser.Username,
			DateJoined: utils.ResponseFormatDate(fetchUser.CreatedAt),
			PostsCount: postsCount,
		}

		apiResponse(w, http.StatusOK, response)
	}
}
