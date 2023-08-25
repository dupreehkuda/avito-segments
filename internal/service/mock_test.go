// Code generated by MockGen. DO NOT EDIT.
// Source: service.go

// Package service is a generated GoMock package.
package service

import (
	context "context"
	reflect "reflect"

	models "github.com/dupreehkuda/avito-segments/internal/models"
	gomock "go.uber.org/mock/gomock"
)

// MockRepository is a mock of Repository interface.
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository.
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance.
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// SegmentAdd mocks base method.
func (m *MockRepository) SegmentAdd(ctx context.Context, segment models.Segment) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SegmentAdd", ctx, segment)
	ret0, _ := ret[0].(error)
	return ret0
}

// SegmentAdd indicates an expected call of SegmentAdd.
func (mr *MockRepositoryMockRecorder) SegmentAdd(ctx, segment interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SegmentAdd", reflect.TypeOf((*MockRepository)(nil).SegmentAdd), ctx, segment)
}

// SegmentDelete mocks base method.
func (m *MockRepository) SegmentDelete(ctx context.Context, tag string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SegmentDelete", ctx, tag)
	ret0, _ := ret[0].(error)
	return ret0
}

// SegmentDelete indicates an expected call of SegmentDelete.
func (mr *MockRepositoryMockRecorder) SegmentDelete(ctx, tag interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SegmentDelete", reflect.TypeOf((*MockRepository)(nil).SegmentDelete), ctx, tag)
}

// SegmentGet mocks base method.
func (m *MockRepository) SegmentGet(ctx context.Context, tag string) (*models.Segment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SegmentGet", ctx, tag)
	ret0, _ := ret[0].(*models.Segment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SegmentGet indicates an expected call of SegmentGet.
func (mr *MockRepositoryMockRecorder) SegmentGet(ctx, tag interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SegmentGet", reflect.TypeOf((*MockRepository)(nil).SegmentGet), ctx, tag)
}
