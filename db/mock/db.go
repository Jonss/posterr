// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/Jonss/posterr/db (interfaces: AppQuerier)

// Package mock_db is a generated GoMock package.
package mock_db

import (
	context "context"
	reflect "reflect"

	db "github.com/Jonss/posterr/db"
	gomock "github.com/golang/mock/gomock"
)

// MockAppQuerier is a mock of AppQuerier interface.
type MockAppQuerier struct {
	ctrl     *gomock.Controller
	recorder *MockAppQuerierMockRecorder
}

// MockAppQuerierMockRecorder is the mock recorder for MockAppQuerier.
type MockAppQuerierMockRecorder struct {
	mock *MockAppQuerier
}

// NewMockAppQuerier creates a new mock instance.
func NewMockAppQuerier(ctrl *gomock.Controller) *MockAppQuerier {
	mock := &MockAppQuerier{ctrl: ctrl}
	mock.recorder = &MockAppQuerierMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAppQuerier) EXPECT() *MockAppQuerierMockRecorder {
	return m.recorder
}

// CreatePost mocks base method.
func (m *MockAppQuerier) CreatePost(arg0 context.Context, arg1 db.CreatePostParams) (db.Post, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreatePost", arg0, arg1)
	ret0, _ := ret[0].(db.Post)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreatePost indicates an expected call of CreatePost.
func (mr *MockAppQuerierMockRecorder) CreatePost(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreatePost", reflect.TypeOf((*MockAppQuerier)(nil).CreatePost), arg0, arg1)
}

// FetchPosts mocks base method.
func (m *MockAppQuerier) FetchPosts(arg0 context.Context, arg1 db.FetchPostsParams) (db.FetchPosts, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FetchPosts", arg0, arg1)
	ret0, _ := ret[0].(db.FetchPosts)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FetchPosts indicates an expected call of FetchPosts.
func (mr *MockAppQuerierMockRecorder) FetchPosts(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FetchPosts", reflect.TypeOf((*MockAppQuerier)(nil).FetchPosts), arg0, arg1)
}

// SeedPost mocks base method.
func (m *MockAppQuerier) SeedPost(arg0 context.Context, arg1 db.SeedPostParams) (db.Post, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SeedPost", arg0, arg1)
	ret0, _ := ret[0].(db.Post)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SeedPost indicates an expected call of SeedPost.
func (mr *MockAppQuerierMockRecorder) SeedPost(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SeedPost", reflect.TypeOf((*MockAppQuerier)(nil).SeedPost), arg0, arg1)
}

// SeedUser mocks base method.
func (m *MockAppQuerier) SeedUser(arg0 context.Context, arg1 string) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SeedUser", arg0, arg1)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SeedUser indicates an expected call of SeedUser.
func (mr *MockAppQuerierMockRecorder) SeedUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SeedUser", reflect.TypeOf((*MockAppQuerier)(nil).SeedUser), arg0, arg1)
}
