// Code generated by MockGen. DO NOT EDIT.
// Source: vec-node/internal/store (interfaces: WorkflowStore,QueueStore,QueueSizeStore)

// Package store_mock is a generated GoMock package.
package store_mock

import (
	context "context"
	reflect "reflect"
	store "vec-node/internal/store"

	gomock "github.com/golang/mock/gomock"
)

// MockWorkflowStore is a mock of WorkflowStore interface.
type MockWorkflowStore struct {
	ctrl     *gomock.Controller
	recorder *MockWorkflowStoreMockRecorder
}

// MockWorkflowStoreMockRecorder is the mock recorder for MockWorkflowStore.
type MockWorkflowStoreMockRecorder struct {
	mock *MockWorkflowStore
}

// NewMockWorkflowStore creates a new mock instance.
func NewMockWorkflowStore(ctrl *gomock.Controller) *MockWorkflowStore {
	mock := &MockWorkflowStore{ctrl: ctrl}
	mock.recorder = &MockWorkflowStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockWorkflowStore) EXPECT() *MockWorkflowStoreMockRecorder {
	return m.recorder
}

// CompleteWorkflow mocks base method.
func (m *MockWorkflowStore) CompleteWorkflow(arg0 context.Context, arg1 int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CompleteWorkflow", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// CompleteWorkflow indicates an expected call of CompleteWorkflow.
func (mr *MockWorkflowStoreMockRecorder) CompleteWorkflow(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CompleteWorkflow", reflect.TypeOf((*MockWorkflowStore)(nil).CompleteWorkflow), arg0, arg1)
}

// GetWorkflowByID mocks base method.
func (m *MockWorkflowStore) GetWorkflowByID(arg0 context.Context, arg1 int) (*store.Workflow, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetWorkflowByID", arg0, arg1)
	ret0, _ := ret[0].(*store.Workflow)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetWorkflowByID indicates an expected call of GetWorkflowByID.
func (mr *MockWorkflowStoreMockRecorder) GetWorkflowByID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetWorkflowByID", reflect.TypeOf((*MockWorkflowStore)(nil).GetWorkflowByID), arg0, arg1)
}

// GetWorkflows mocks base method.
func (m *MockWorkflowStore) GetWorkflows(arg0 context.Context, arg1 *store.WorkflowFilter) ([]store.Workflow, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetWorkflows", arg0, arg1)
	ret0, _ := ret[0].([]store.Workflow)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetWorkflows indicates an expected call of GetWorkflows.
func (mr *MockWorkflowStoreMockRecorder) GetWorkflows(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetWorkflows", reflect.TypeOf((*MockWorkflowStore)(nil).GetWorkflows), arg0, arg1)
}

// SaveWorkflow mocks base method.
func (m *MockWorkflowStore) SaveWorkflow(arg0 context.Context, arg1 *store.Workflow) (*store.Workflow, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveWorkflow", arg0, arg1)
	ret0, _ := ret[0].(*store.Workflow)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SaveWorkflow indicates an expected call of SaveWorkflow.
func (mr *MockWorkflowStoreMockRecorder) SaveWorkflow(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveWorkflow", reflect.TypeOf((*MockWorkflowStore)(nil).SaveWorkflow), arg0, arg1)
}

// StartWorkflow mocks base method.
func (m *MockWorkflowStore) StartWorkflow(arg0 context.Context, arg1 int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StartWorkflow", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// StartWorkflow indicates an expected call of StartWorkflow.
func (mr *MockWorkflowStoreMockRecorder) StartWorkflow(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StartWorkflow", reflect.TypeOf((*MockWorkflowStore)(nil).StartWorkflow), arg0, arg1)
}

// UpdateWorkflow mocks base method.
func (m *MockWorkflowStore) UpdateWorkflow(arg0 context.Context, arg1 *store.Workflow) (*store.Workflow, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateWorkflow", arg0, arg1)
	ret0, _ := ret[0].(*store.Workflow)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateWorkflow indicates an expected call of UpdateWorkflow.
func (mr *MockWorkflowStoreMockRecorder) UpdateWorkflow(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateWorkflow", reflect.TypeOf((*MockWorkflowStore)(nil).UpdateWorkflow), arg0, arg1)
}

// MockQueueStore is a mock of QueueStore interface.
type MockQueueStore struct {
	ctrl     *gomock.Controller
	recorder *MockQueueStoreMockRecorder
}

// MockQueueStoreMockRecorder is the mock recorder for MockQueueStore.
type MockQueueStoreMockRecorder struct {
	mock *MockQueueStore
}

// NewMockQueueStore creates a new mock instance.
func NewMockQueueStore(ctrl *gomock.Controller) *MockQueueStore {
	mock := &MockQueueStore{ctrl: ctrl}
	mock.recorder = &MockQueueStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockQueueStore) EXPECT() *MockQueueStoreMockRecorder {
	return m.recorder
}

// CompleteWorkflowInQueue mocks base method.
func (m *MockQueueStore) CompleteWorkflowInQueue(arg0 context.Context, arg1 int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CompleteWorkflowInQueue", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// CompleteWorkflowInQueue indicates an expected call of CompleteWorkflowInQueue.
func (mr *MockQueueStoreMockRecorder) CompleteWorkflowInQueue(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CompleteWorkflowInQueue", reflect.TypeOf((*MockQueueStore)(nil).CompleteWorkflowInQueue), arg0, arg1)
}

// Enqueue mocks base method.
func (m *MockQueueStore) Enqueue(arg0 context.Context, arg1 int) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Enqueue", arg0, arg1)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Enqueue indicates an expected call of Enqueue.
func (mr *MockQueueStoreMockRecorder) Enqueue(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Enqueue", reflect.TypeOf((*MockQueueStore)(nil).Enqueue), arg0, arg1)
}

// GetCount mocks base method.
func (m *MockQueueStore) GetCount(arg0 context.Context) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCount", arg0)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCount indicates an expected call of GetCount.
func (mr *MockQueueStoreMockRecorder) GetCount(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCount", reflect.TypeOf((*MockQueueStore)(nil).GetCount), arg0)
}

// GetQueueByID mocks base method.
func (m *MockQueueStore) GetQueueByID(arg0 context.Context, arg1 int) (*store.Queue, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetQueueByID", arg0, arg1)
	ret0, _ := ret[0].(*store.Queue)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetQueueByID indicates an expected call of GetQueueByID.
func (mr *MockQueueStoreMockRecorder) GetQueueByID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetQueueByID", reflect.TypeOf((*MockQueueStore)(nil).GetQueueByID), arg0, arg1)
}

// GetQueueStatus mocks base method.
func (m *MockQueueStore) GetQueueStatus(arg0 context.Context) ([]store.Queue, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetQueueStatus", arg0)
	ret0, _ := ret[0].([]store.Queue)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetQueueStatus indicates an expected call of GetQueueStatus.
func (mr *MockQueueStoreMockRecorder) GetQueueStatus(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetQueueStatus", reflect.TypeOf((*MockQueueStore)(nil).GetQueueStatus), arg0)
}

// IsEmpty mocks base method.
func (m *MockQueueStore) IsEmpty(arg0 context.Context) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsEmpty", arg0)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsEmpty indicates an expected call of IsEmpty.
func (mr *MockQueueStoreMockRecorder) IsEmpty(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsEmpty", reflect.TypeOf((*MockQueueStore)(nil).IsEmpty), arg0)
}

// IsSpaceAvailable mocks base method.
func (m *MockQueueStore) IsSpaceAvailable(arg0 context.Context) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsSpaceAvailable", arg0)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsSpaceAvailable indicates an expected call of IsSpaceAvailable.
func (mr *MockQueueStoreMockRecorder) IsSpaceAvailable(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsSpaceAvailable", reflect.TypeOf((*MockQueueStore)(nil).IsSpaceAvailable), arg0)
}

// PeekInProcess mocks base method.
func (m *MockQueueStore) PeekInProcess(arg0 context.Context) (*store.Queue, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PeekInProcess", arg0)
	ret0, _ := ret[0].(*store.Queue)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PeekInProcess indicates an expected call of PeekInProcess.
func (mr *MockQueueStoreMockRecorder) PeekInProcess(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PeekInProcess", reflect.TypeOf((*MockQueueStore)(nil).PeekInProcess), arg0)
}

// PeekQueued mocks base method.
func (m *MockQueueStore) PeekQueued(arg0 context.Context) (*store.Queue, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PeekQueued", arg0)
	ret0, _ := ret[0].(*store.Queue)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PeekQueued indicates an expected call of PeekQueued.
func (mr *MockQueueStoreMockRecorder) PeekQueued(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PeekQueued", reflect.TypeOf((*MockQueueStore)(nil).PeekQueued), arg0)
}

// ProcessWorkflowInQueue mocks base method.
func (m *MockQueueStore) ProcessWorkflowInQueue(arg0 context.Context, arg1 int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ProcessWorkflowInQueue", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// ProcessWorkflowInQueue indicates an expected call of ProcessWorkflowInQueue.
func (mr *MockQueueStoreMockRecorder) ProcessWorkflowInQueue(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProcessWorkflowInQueue", reflect.TypeOf((*MockQueueStore)(nil).ProcessWorkflowInQueue), arg0, arg1)
}

// MockQueueSizeStore is a mock of QueueSizeStore interface.
type MockQueueSizeStore struct {
	ctrl     *gomock.Controller
	recorder *MockQueueSizeStoreMockRecorder
}

// MockQueueSizeStoreMockRecorder is the mock recorder for MockQueueSizeStore.
type MockQueueSizeStoreMockRecorder struct {
	mock *MockQueueSizeStore
}

// NewMockQueueSizeStore creates a new mock instance.
func NewMockQueueSizeStore(ctrl *gomock.Controller) *MockQueueSizeStore {
	mock := &MockQueueSizeStore{ctrl: ctrl}
	mock.recorder = &MockQueueSizeStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockQueueSizeStore) EXPECT() *MockQueueSizeStoreMockRecorder {
	return m.recorder
}

// GetQueueSize mocks base method.
func (m *MockQueueSizeStore) GetQueueSize(arg0 context.Context) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetQueueSize", arg0)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetQueueSize indicates an expected call of GetQueueSize.
func (mr *MockQueueSizeStoreMockRecorder) GetQueueSize(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetQueueSize", reflect.TypeOf((*MockQueueSizeStore)(nil).GetQueueSize), arg0)
}

// SetQueueSize mocks base method.
func (m *MockQueueSizeStore) SetQueueSize(arg0 context.Context, arg1 int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetQueueSize", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetQueueSize indicates an expected call of SetQueueSize.
func (mr *MockQueueSizeStoreMockRecorder) SetQueueSize(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetQueueSize", reflect.TypeOf((*MockQueueSizeStore)(nil).SetQueueSize), arg0, arg1)
}

// UpdateQueueSize mocks base method.
func (m *MockQueueSizeStore) UpdateQueueSize(arg0 context.Context, arg1 int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateQueueSize", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateQueueSize indicates an expected call of UpdateQueueSize.
func (mr *MockQueueSizeStoreMockRecorder) UpdateQueueSize(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateQueueSize", reflect.TypeOf((*MockQueueSizeStore)(nil).UpdateQueueSize), arg0, arg1)
}
