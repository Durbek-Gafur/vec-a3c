// Code generated by MockGen. DO NOT EDIT.
// Source: vec-node/internal/queue (interfaces: Queue)

// Package queueu_mock is a generated GoMock package.
package queueu_mock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockQueue is a mock of Queue interface.
type MockQueue struct {
	ctrl     *gomock.Controller
	recorder *MockQueueMockRecorder
}

// MockQueueMockRecorder is the mock recorder for MockQueue.
type MockQueueMockRecorder struct {
	mock *MockQueue
}

// NewMockQueue creates a new mock instance.
func NewMockQueue(ctrl *gomock.Controller) *MockQueue {
	mock := &MockQueue{ctrl: ctrl}
	mock.recorder = &MockQueueMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockQueue) EXPECT() *MockQueueMockRecorder {
	return m.recorder
}

// StartPeriodicCheck mocks base method.
func (m *MockQueue) StartPeriodicCheck(arg0 context.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "StartPeriodicCheck", arg0)
}

// StartPeriodicCheck indicates an expected call of StartPeriodicCheck.
func (mr *MockQueueMockRecorder) StartPeriodicCheck(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StartPeriodicCheck", reflect.TypeOf((*MockQueue)(nil).StartPeriodicCheck), arg0)
}
