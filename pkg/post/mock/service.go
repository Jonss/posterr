// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/Jonss/posterr/pkg/post (interfaces: Service)

// Package mock_post is a generated GoMock package.
package mock_post

import (
	context "context"
	reflect "reflect"

	post "github.com/Jonss/posterr/pkg/post"
	gomock "github.com/golang/mock/gomock"
)

// MockService is a mock of Service interface.
type MockService struct {
	ctrl     *gomock.Controller
	recorder *MockServiceMockRecorder
}

// MockServiceMockRecorder is the mock recorder for MockService.
type MockServiceMockRecorder struct {
	mock *MockService
}

// NewMockService creates a new mock instance.
func NewMockService(ctrl *gomock.Controller) *MockService {
	mock := &MockService{ctrl: ctrl}
	mock.recorder = &MockServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockService) EXPECT() *MockServiceMockRecorder {
	return m.recorder
}

// FetchPosts mocks base method.
func (m *MockService) FetchPosts(arg0 context.Context, arg1 post.FetchPostParams) ([]post.FetchPost, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FetchPosts", arg0, arg1)
	ret0, _ := ret[0].([]post.FetchPost)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FetchPosts indicates an expected call of FetchPosts.
func (mr *MockServiceMockRecorder) FetchPosts(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FetchPosts", reflect.TypeOf((*MockService)(nil).FetchPosts), arg0, arg1)
}