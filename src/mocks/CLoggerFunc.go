// Code generated by mockery 2.14.0. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	log "github.com/mikhailbolshakov/cryptocare/src/kit/log"
)

// CLoggerFunc is an autogenerated mock type for the CLoggerFunc type
type CLoggerFunc struct {
	mock.Mock
}

// Execute provides a mock function with given fields:
func (_m *CLoggerFunc) Execute() log.CLogger {
	ret := _m.Called()

	var r0 log.CLogger
	if rf, ok := ret.Get(0).(func() log.CLogger); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(log.CLogger)
		}
	}

	return r0
}

type mockConstructorTestingTNewCLoggerFunc interface {
	mock.TestingT
	Cleanup(func())
}

// NewCLoggerFunc creates a new instance of CLoggerFunc. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewCLoggerFunc(t mockConstructorTestingTNewCLoggerFunc) *CLoggerFunc {
	mock := &CLoggerFunc{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
