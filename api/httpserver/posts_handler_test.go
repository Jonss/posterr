package httpserver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Jonss/posterr/pkg/post"
	post_mock "github.com/Jonss/posterr/pkg/post/mock"
	"github.com/Jonss/posterr/pkg/utils"
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

func TestCreatePost(t *testing.T) {
	testCases := []struct {
		name              string
		requestBody       string
		wantStatusCode    int
		buildStubs        func(service *post_mock.MockService)
		isErrorWant       bool
		wantErrorResponse ErrorResponses
	}{
		{
			name: "should create a post",
			requestBody: `
			{
				"user_id": 1,
				"message": "Ahoy, World!",
				"originalPostId": null
			}`,
			wantStatusCode: http.StatusCreated,
			buildStubs: func(service *post_mock.MockService) {
				service.EXPECT().
					CountDailyPosts(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil)

				service.EXPECT().
					CreatePost(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&post.CreatePostResponse{
						ID:             1,
						Content:        utils.StrToPointer("Ahoy, World!"),
						OriginalPostID: nil,
					}, nil)
			},
		},
		{
			name:           "should get errors when request body is empty",
			requestBody:    `{}`,
			wantStatusCode: http.StatusBadRequest,
			buildStubs:     func(service *post_mock.MockService) {},
			isErrorWant:    true,
			wantErrorResponse: NewErrorResponses(
				NewErrorResponse("userid is a required field"),
				NewErrorResponse("message is a required field"),
				NewErrorResponse("originalpostid is a required field"),
			),
		},
		{
			name: "should get errors when originalPostId length is above 777 characters",
			requestBody: `{
				"user_id": 1,
				"message": "Lorem ipsum dolor sit amet, consectetuer adipiscing elit. Aenean commodo ligula eget dolor. Aenean massa. Cum sociis natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus. Donec quam felis, ultricies nec, pellentesque eu, pretium quis, sem. Nulla consequat massa quis enim. Donec pede justo, fringilla vel, aliquet nec, vulputate eget, arcu. In enim justo, rhoncus ut, imperdiet a, venenatis vitae, justo. Nullam dictum felis eu pede mollis pretium. Integer tincidunt. Cras dapibus. Vivamus elementum semper nisi. Aenean vulputate eleifend tellus. Aenean leo ligula, porttitor eu, consequat vitae, eleifend ac, enim. Aliquam lorem ante, dapibus in, viverra quis, feugiat a, tellus. Phasellus viverra nulla ut metus varius laoreet. Quisque rutrum. Aenean imp8",
				"originalPostId": null
			}`,
			wantStatusCode: http.StatusBadRequest,
			buildStubs:     func(service *post_mock.MockService) {},
			isErrorWant:    true,
			wantErrorResponse: NewErrorResponses(
				NewErrorResponse("message must be at maximum 777 characters in length"),
			),
		},
		{
			name: "should get an error when user already created 5 posts within a day",
			requestBody: `
			{
				"user_id": 1,
				"message": "Ahoy, World!",
				"originalPostId": null
			}`,
			wantStatusCode: http.StatusUnprocessableEntity,
			buildStubs: func(service *post_mock.MockService) {
				service.EXPECT().
					CountDailyPosts(gomock.Any(), gomock.Any()).
					Times(1).
					Return(post.ErrMaxPostsDailyAchieved)
			},
			isErrorWant: true,
			wantErrorResponse: NewErrorResponses(
				NewErrorResponse("error user is not allowed to post more than 5 times within a day")),
		},
		{
			name: "should get an error original post is a reposting",
			requestBody: `
			{
				"user_id": 1,
				"message": "Ahoy, World!",
				"originalPostId": 1
			}`,
			wantStatusCode: http.StatusUnprocessableEntity,
			buildStubs: func(service *post_mock.MockService) {
				service.EXPECT().
					CountDailyPosts(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil)

				service.EXPECT().
					CreatePost(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, post.ErrRepost)
			},
			isErrorWant: true,
			wantErrorResponse: NewErrorResponses(
				NewErrorResponse("error user cannot repost an existing repost")),
		},
		{
			name: "should get an error original post is a quote-post",
			requestBody: `
			{
				"user_id": 1,
				"message": "Ahoy, World!",
				"originalPostId": 1
			}`,
			wantStatusCode: http.StatusUnprocessableEntity,
			buildStubs: func(service *post_mock.MockService) {
				service.EXPECT().
					CountDailyPosts(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil)

				service.EXPECT().
					CreatePost(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, post.ErrQuotePost)
			},
			isErrorWant: true,
			wantErrorResponse: NewErrorResponses(
				NewErrorResponse("error user cannot quote an existing quote-post")),
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

			r := httptest.NewRequest(http.MethodPost, "/api/posts", bytes.NewBuffer([]byte(tc.requestBody)))

			w := httptest.NewRecorder()
			srv.router.ServeHTTP(w, r)

			got := w.Result()

			if tc.wantStatusCode != got.StatusCode {
				t.Fatalf("POST /api/posts . status code. want %d, got %d", tc.wantStatusCode, got.StatusCode)
			}

			if tc.isErrorWant {
				var response ErrorResponses
				err := json.NewDecoder(got.Body).Decode(&response)
				if err != nil {
					t.Fatalf("unexpected error decoding error response. error=(%v)", err)
				}

				for i, errResp := range response.Errors {
					fmt.Println(errResp.Message)
					if tc.wantErrorResponse.Errors[i].Message != errResp.Message {
						t.Errorf("POST /api/posts. want error message %s, got %s", tc.wantErrorResponse.Errors[i].Message, errResp.Message)
					}
				}
			} else {
				var response CreatePostRequest
				err := json.NewDecoder(got.Body).Decode(&response)
				if err != nil {
					t.Fatalf("unexpected error decoding success response. error=(%v)", err)
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
