// Code generated by mockery 2.14.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// RouteSetter is an autogenerated mock type for the RouteSetter type
type RouteSetter struct {
	mock.Mock
}

// Set provides a mock function with given fields:
func (_m *RouteSetter) Set() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewRouteSetter interface {
	mock.TestingT
	Cleanup(func())
}

// NewRouteSetter creates a new instance of RouteSetter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewRouteSetter(t mockConstructorTestingTNewRouteSetter) *RouteSetter {
	mock := &RouteSetter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
