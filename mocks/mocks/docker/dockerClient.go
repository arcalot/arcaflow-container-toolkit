// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/arcalot/arcaflow-plugin-image-builder/internal/docker (interfaces: DockerClient)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	io "io"
	reflect "reflect"

	types "github.com/docker/docker/api/types"
	gomock "github.com/golang/mock/gomock"
)

// MockDockerClient is a mock of DockerClient interface.
type MockDockerClient struct {
	ctrl     *gomock.Controller
	recorder *MockDockerClientMockRecorder
}

// MockDockerClientMockRecorder is the mock recorder for MockDockerClient.
type MockDockerClientMockRecorder struct {
	mock *MockDockerClient
}

// NewMockDockerClient creates a new mock instance.
func NewMockDockerClient(ctrl *gomock.Controller) *MockDockerClient {
	mock := &MockDockerClient{ctrl: ctrl}
	mock.recorder = &MockDockerClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDockerClient) EXPECT() *MockDockerClientMockRecorder {
	return m.recorder
}

// ImageBuild mocks base method.
func (m *MockDockerClient) ImageBuild(arg0 context.Context, arg1 io.Reader, arg2 types.ImageBuildOptions) (types.ImageBuildResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ImageBuild", arg0, arg1, arg2)
	ret0, _ := ret[0].(types.ImageBuildResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ImageBuild indicates an expected call of ImageBuild.
func (mr *MockDockerClientMockRecorder) ImageBuild(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ImageBuild", reflect.TypeOf((*MockDockerClient)(nil).ImageBuild), arg0, arg1, arg2)
}

// ImagePush mocks base method.
func (m *MockDockerClient) ImagePush(arg0 context.Context, arg1 string, arg2 types.ImagePushOptions) (io.ReadCloser, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ImagePush", arg0, arg1, arg2)
	ret0, _ := ret[0].(io.ReadCloser)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ImagePush indicates an expected call of ImagePush.
func (mr *MockDockerClientMockRecorder) ImagePush(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ImagePush", reflect.TypeOf((*MockDockerClient)(nil).ImagePush), arg0, arg1, arg2)
}

// ImageTag mocks base method.
func (m *MockDockerClient) ImageTag(arg0 context.Context, arg1, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ImageTag", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// ImageTag indicates an expected call of ImageTag.
func (mr *MockDockerClientMockRecorder) ImageTag(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ImageTag", reflect.TypeOf((*MockDockerClient)(nil).ImageTag), arg0, arg1, arg2)
}
