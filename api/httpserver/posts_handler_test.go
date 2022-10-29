package httpserver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Jonss/posterr/pkg/post"
	post_mock "github.com/Jonss/posterr/pkg/post/mock"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
)

func TestFetchPosts(t *testing.T) {
	testCases := []struct {
		name           string
		endpoint       string
		wantStatusCode int
		buildStubs     func(service *post_mock.MockService)
		isErrorWant    bool
		errorMessage   string
		wantResponse   FetchPostsResponse
	}{
		{
			name:           "should get an empty list of posts",
			endpoint:       "",
			wantStatusCode: http.StatusOK,
			buildStubs: func(service *post_mock.MockService) {
				service.EXPECT().
					FetchPosts(gomock.Any(), gomock.Any()).
					Times(1).
					Return(post.FetchPostResponse{}, nil)
			},
			wantResponse: FetchPostsResponse{},
		},
		{
			name:           "should get an error when user param 'only_mine' is set to true and 'user_id' is unset",
			endpoint:       "?only_mine=true",
			wantStatusCode: http.StatusBadRequest,
			buildStubs: func(service *post_mock.MockService) {
			},
			isErrorWant:  true,
			errorMessage: ErrUserIDRequired.Error(),
		},
		{
			name:           "should get a list of posts with 5 items",
			endpoint:       "",
			wantStatusCode: http.StatusOK,
			buildStubs: func(service *post_mock.MockService) {
				service.EXPECT().
					FetchPosts(gomock.Any(), gomock.Any()).
					Times(1).
					Return(fetchPosts(5), nil)
			},
			wantResponse: wantResponse(5),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// setup
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			postServiceMock := post_mock.NewMockService(ctrl)
			tc.buildStubs(postServiceMock)

			srv := NewHttpServer(
				mux.NewRouter(), fakeConfig, Services{postServiceMock})
			srv.Start()
			// end setup

			r := httptest.NewRequest(http.MethodGet, "/api/posts"+tc.endpoint, nil)

			w := httptest.NewRecorder()
			srv.router.ServeHTTP(w, r)

			got := w.Result()

			if tc.wantStatusCode != got.StatusCode {
				t.Fatalf("GET /api/posts . status code. want %d, got %d", tc.wantStatusCode, got.StatusCode)
			}

			if tc.isErrorWant {
				var response ErrorResponses
				err := json.NewDecoder(got.Body).Decode(&response)
				if err != nil {
					t.Fatalf("unexpected error decoding error response. error=(%v)", err)
				}

				if tc.errorMessage != response.Errors[0].Message {
					t.Errorf("error.Message want '%v', got '%v'", tc.errorMessage, response.Errors[0].Message)
				}
			} else {
				var response FetchPostsResponse
				err := json.NewDecoder(got.Body).Decode(&response)
				if err != nil {
					t.Fatalf("unexpected error decoding success response. error=(%v)", err)
				}
				if len(tc.wantResponse.FetchPostResponses) != len(response.FetchPostResponses) {
					t.Errorf("len(FetchResponses) want %v, got %v", len(tc.wantResponse.FetchPostResponses), len(response.FetchPostResponses))
				}

				if tc.wantResponse.HasNext != response.HasNext {
					t.Errorf("response.hasNext want %v, got %v", tc.wantResponse.HasNext, response.HasNext)
				}

				if tc.wantResponse.HasPrev != response.HasPrev {
					t.Errorf("response.hasPrev want %v, got %v", tc.wantResponse.HasPrev, response.HasPrev)
				}
			}
		})
	}
}

func fetchPosts(times int) post.FetchPostResponse {
	posts := make([]post.FetchPost, times)
	for i := 0; i < times; i++ {
		message := fmt.Sprintf("A pleasent post %d", i)
		posts[i] = post.FetchPost{
			Post: post.Post{
				Message:  &message,
				Username: "snarf",
				Type:     post.Original,
			},
		}
	}
	return post.FetchPostResponse{
		Posts: posts,
	}
}

func wantResponse(times int) FetchPostsResponse {
	fetchPostResponses := make([]FetchPostResponse, times)
	for i := 0; i < times; i++ {
		message := fmt.Sprintf("A pleasent post %d", i)
		fetchPostResponses[i] = FetchPostResponse{
			Post: Post{
				ID:      int64(i + 1),
				Message: &message,
				Type:    post.Original,
			},
		}
	}

	return FetchPostsResponse{
		FetchPostResponses: fetchPostResponses,
		HasNext:            false,
		HasPrev:            false,
	}
}
