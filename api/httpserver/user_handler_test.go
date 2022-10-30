package httpserver

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	post_mock "github.com/Jonss/posterr/pkg/post/mock"
	"github.com/Jonss/posterr/pkg/user"
	user_mock "github.com/Jonss/posterr/pkg/user/mock"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
)

func TestFetchUser(t *testing.T) {
	testCases := []struct {
		name           string
		username       string
		wantStatusCode int
		buildStubs     func(postService *post_mock.MockService, userService *user_mock.MockService)
		isErrorWant    bool
		wantErr        error
		wantResponse   FetchUserResponse
	}{
		{
			name:           "should fetch a user",
			username:       "brain",
			wantStatusCode: http.StatusOK,
			buildStubs: func(postService *post_mock.MockService, userService *user_mock.MockService) {
				date := time.Date(2021, time.May, 2, 0, 0, 0, 0, time.UTC)
				userService.EXPECT().
					FetchUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&user.FetchUser{
						UserID:    1,
						Username:  "brain",
						CreatedAt: date,
					}, nil)

				postService.EXPECT().
					CountPosts(gomock.Any(), gomock.Any()).
					Times(1).
					Return(int64(42), nil)
			},
			wantResponse: FetchUserResponse{
				Username:   "brain",
				DateJoined: "May 2, 2021",
				PostsCount: 42,
			},
		},
		{
			name:           "should get an error when user does not exists",
			username:       "brainiac",
			wantStatusCode: http.StatusNotFound,
			buildStubs: func(postService *post_mock.MockService, userService *user_mock.MockService) {
				userService.EXPECT().
					FetchUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, user.ErrUserNotFound)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// setup
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			postServiceMock := post_mock.NewMockService(ctrl)
			userServiceMock := user_mock.NewMockService(ctrl)
			tc.buildStubs(postServiceMock, userServiceMock)

			srv := NewHttpServer(
				mux.NewRouter(), fakeConfig, Services{postServiceMock, userServiceMock})
			srv.Start()
			// end setup

			r := httptest.NewRequest(http.MethodGet, "/api/users/"+tc.username, nil)

			w := httptest.NewRecorder()
			srv.router.ServeHTTP(w, r)

			got := w.Result()

			if tc.wantStatusCode != got.StatusCode {
				t.Fatalf("GET /api/users/%s . status code. want %d, got %d", tc.username, tc.wantStatusCode, got.StatusCode)
			}

			if tc.isErrorWant {
				var response ErrorResponses
				err := json.NewDecoder(got.Body).Decode(&response)
				if err != nil {
					t.Fatalf("unexpected error decoding error response. error=(%v)", err)
				}

				if tc.wantErr.Error() != response.Errors[0].Message {
					t.Errorf("error.Message want '%v', got '%v'", tc.wantErr.Error(), response.Errors[0].Message)
				}
			} else {
				var response FetchUserResponse
				err := json.NewDecoder(got.Body).Decode(&response)
				if err != nil {
					t.Fatalf("unexpected error decoding success response. error=(%v)", err)
				}

				if tc.wantResponse.Username != response.Username {
					t.Errorf("response.Username want %v, got %v", tc.wantResponse.Username, response.Username)
				}

				if tc.wantResponse.DateJoined != response.DateJoined {
					t.Errorf("response.DateJoined want %v, got %v", tc.wantResponse.DateJoined, response.DateJoined)
				}
				if tc.wantResponse.PostsCount != response.PostsCount {
					t.Errorf("response.PostsCount want %v, got %v", tc.wantResponse.PostsCount, response.PostsCount)
				}
			}
		})
	}

}
