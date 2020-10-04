// Code generated by MockGen. DO NOT EDIT.
// Source: legacy/api.go

// Package mock_legacy is a generated GoMock package.
package mock_legacy

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	entity "maga-auctions/entity"
	legacy "maga-auctions/legacy"
	reflect "reflect"
)

// MockAPI is a mock of API interface
type MockAPI struct {
	ctrl     *gomock.Controller
	recorder *MockAPIMockRecorder
}

// MockAPIMockRecorder is the mock recorder for MockAPI
type MockAPIMockRecorder struct {
	mock *MockAPI
}

// NewMockAPI creates a new mock instance
func NewMockAPI(ctrl *gomock.Controller) *MockAPI {
	mock := &MockAPI{ctrl: ctrl}
	mock.recorder = &MockAPIMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockAPI) EXPECT() *MockAPIMockRecorder {
	return m.recorder
}

// Get mocks base method
func (m *MockAPI) Get(ctx context.Context) ([]legacy.VehicleLegacy, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx)
	ret0, _ := ret[0].([]legacy.VehicleLegacy)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get
func (mr *MockAPIMockRecorder) Get(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockAPI)(nil).Get), ctx)
}

// Create mocks base method
func (m *MockAPI) Create(ctx context.Context, vehicle *entity.Vehicle) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, vehicle)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create
func (mr *MockAPIMockRecorder) Create(ctx, vehicle interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockAPI)(nil).Create), ctx, vehicle)
}

// Update mocks base method
func (m *MockAPI) Update(ctx context.Context, vehicle *entity.Vehicle) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, vehicle)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update
func (mr *MockAPIMockRecorder) Update(ctx, vehicle interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockAPI)(nil).Update), ctx, vehicle)
}

// Delete mocks base method
func (m *MockAPI) Delete(ctx context.Context, id int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete
func (mr *MockAPIMockRecorder) Delete(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockAPI)(nil).Delete), ctx, id)
}
