// Code generated by mockery 2.14.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	er "github.com/mikhailbolshakov/cryptocare/src/kit/er"
)

// AppErrBuilder is an autogenerated mock type for the AppErrBuilder type
type AppErrBuilder struct {
	mock.Mock
}

// Business provides a mock function with given fields:
func (_m *AppErrBuilder) Business() er.AppErrBuilder {
	ret := _m.Called()

	var r0 er.AppErrBuilder
	if rf, ok := ret.Get(0).(func() er.AppErrBuilder); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(er.AppErrBuilder)
		}
	}

	return r0
}

// C provides a mock function with given fields: ctx
func (_m *AppErrBuilder) C(ctx context.Context) er.AppErrBuilder {
	ret := _m.Called(ctx)

	var r0 er.AppErrBuilder
	if rf, ok := ret.Get(0).(func(context.Context) er.AppErrBuilder); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(er.AppErrBuilder)
		}
	}

	return r0
}

// Err provides a mock function with given fields:
func (_m *AppErrBuilder) Err() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// F provides a mock function with given fields: fields
func (_m *AppErrBuilder) F(fields er.FF) er.AppErrBuilder {
	ret := _m.Called(fields)

	var r0 er.AppErrBuilder
	if rf, ok := ret.Get(0).(func(er.FF) er.AppErrBuilder); ok {
		r0 = rf(fields)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(er.AppErrBuilder)
		}
	}

	return r0
}

// GrpcSt provides a mock function with given fields: status
func (_m *AppErrBuilder) GrpcSt(status uint32) er.AppErrBuilder {
	ret := _m.Called(status)

	var r0 er.AppErrBuilder
	if rf, ok := ret.Get(0).(func(uint32) er.AppErrBuilder); ok {
		r0 = rf(status)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(er.AppErrBuilder)
		}
	}

	return r0
}

// HttpSt provides a mock function with given fields: status
func (_m *AppErrBuilder) HttpSt(status uint32) er.AppErrBuilder {
	ret := _m.Called(status)

	var r0 er.AppErrBuilder
	if rf, ok := ret.Get(0).(func(uint32) er.AppErrBuilder); ok {
		r0 = rf(status)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(er.AppErrBuilder)
		}
	}

	return r0
}

// Panic provides a mock function with given fields:
func (_m *AppErrBuilder) Panic() er.AppErrBuilder {
	ret := _m.Called()

	var r0 er.AppErrBuilder
	if rf, ok := ret.Get(0).(func() er.AppErrBuilder); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(er.AppErrBuilder)
		}
	}

	return r0
}

// System provides a mock function with given fields:
func (_m *AppErrBuilder) System() er.AppErrBuilder {
	ret := _m.Called()

	var r0 er.AppErrBuilder
	if rf, ok := ret.Get(0).(func() er.AppErrBuilder); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(er.AppErrBuilder)
		}
	}

	return r0
}

// Type provides a mock function with given fields: t
func (_m *AppErrBuilder) Type(t string) er.AppErrBuilder {
	ret := _m.Called(t)

	var r0 er.AppErrBuilder
	if rf, ok := ret.Get(0).(func(string) er.AppErrBuilder); ok {
		r0 = rf(t)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(er.AppErrBuilder)
		}
	}

	return r0
}

type mockConstructorTestingTNewAppErrBuilder interface {
	mock.TestingT
	Cleanup(func())
}

// NewAppErrBuilder creates a new instance of AppErrBuilder. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewAppErrBuilder(t mockConstructorTestingTNewAppErrBuilder) *AppErrBuilder {
	mock := &AppErrBuilder{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}