// Code generated by mockery 2.14.0. DO NOT EDIT.

package mocks

import (
	context "context"

	auth "github.com/mikhailbolshakov/cryptocare/src/kit/auth"

	http "net/http"

	mock "github.com/stretchr/testify/mock"
)

// ResourcePolicyManager is an autogenerated mock type for the ResourcePolicyManager type
type ResourcePolicyManager struct {
	mock.Mock
}

// GetRequestedResources provides a mock function with given fields: ctx, routeId, r
func (_m *ResourcePolicyManager) GetRequestedResources(ctx context.Context, routeId string, r *http.Request) ([]*auth.AuthorizationResource, error) {
	ret := _m.Called(ctx, routeId, r)

	var r0 []*auth.AuthorizationResource
	if rf, ok := ret.Get(0).(func(context.Context, string, *http.Request) []*auth.AuthorizationResource); ok {
		r0 = rf(ctx, routeId, r)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*auth.AuthorizationResource)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, *http.Request) error); ok {
		r1 = rf(ctx, routeId, r)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RegisterResourceMapping provides a mock function with given fields: routeId, policies
func (_m *ResourcePolicyManager) RegisterResourceMapping(routeId string, policies ...auth.ResourcePolicy) {
	_va := make([]interface{}, len(policies))
	for _i := range policies {
		_va[_i] = policies[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, routeId)
	_ca = append(_ca, _va...)
	_m.Called(_ca...)
}

type mockConstructorTestingTNewResourcePolicyManager interface {
	mock.TestingT
	Cleanup(func())
}

// NewResourcePolicyManager creates a new instance of ResourcePolicyManager. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewResourcePolicyManager(t mockConstructorTestingTNewResourcePolicyManager) *ResourcePolicyManager {
	mock := &ResourcePolicyManager{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
