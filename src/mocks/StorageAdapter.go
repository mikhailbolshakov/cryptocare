// Code generated by mockery 2.14.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// StorageAdapter is an autogenerated mock type for the StorageAdapter type
type StorageAdapter struct {
	mock.Mock
}

// Close provides a mock function with given fields: ctx
func (_m *StorageAdapter) Close(ctx context.Context) error {
	ret := _m.Called(ctx)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Init provides a mock function with given fields: ctx, cfg
func (_m *StorageAdapter) Init(ctx context.Context, cfg interface{}) error {
	ret := _m.Called(ctx, cfg)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, interface{}) error); ok {
		r0 = rf(ctx, cfg)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewStorageAdapter interface {
	mock.TestingT
	Cleanup(func())
}

// NewStorageAdapter creates a new instance of StorageAdapter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewStorageAdapter(t mockConstructorTestingTNewStorageAdapter) *StorageAdapter {
	mock := &StorageAdapter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}