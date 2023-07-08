// Code generated by mockery 2.14.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	domain "github.com/mikhailbolshakov/cryptocare/src/domain"

	service "github.com/mikhailbolshakov/cryptocare/src/service"
)

// ArbitrageService is an autogenerated mock type for the ArbitrageService type
type ArbitrageService struct {
	mock.Mock
}

// GetProfitableChain provides a mock function with given fields: ctx, chainId
func (_m *ArbitrageService) GetProfitableChain(ctx context.Context, chainId string) (*domain.ProfitableChain, error) {
	ret := _m.Called(ctx, chainId)

	var r0 *domain.ProfitableChain
	if rf, ok := ret.Get(0).(func(context.Context, string) *domain.ProfitableChain); ok {
		r0 = rf(ctx, chainId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.ProfitableChain)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, chainId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetProfitableChains provides a mock function with given fields: ctx, rq
func (_m *ArbitrageService) GetProfitableChains(ctx context.Context, rq *domain.GetProfitableChainsRequest) (*domain.GetProfitableChainsResponse, error) {
	ret := _m.Called(ctx, rq)

	var r0 *domain.GetProfitableChainsResponse
	if rf, ok := ret.Get(0).(func(context.Context, *domain.GetProfitableChainsRequest) *domain.GetProfitableChainsResponse); ok {
		r0 = rf(ctx, rq)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.GetProfitableChainsResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *domain.GetProfitableChainsRequest) error); ok {
		r1 = rf(ctx, rq)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Init provides a mock function with given fields: cfg
func (_m *ArbitrageService) Init(cfg *service.Config) {
	_m.Called(cfg)
}

// RunCalculationBackground provides a mock function with given fields: ctx
func (_m *ArbitrageService) RunCalculationBackground(ctx context.Context) error {
	ret := _m.Called(ctx)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// StopCalculation provides a mock function with given fields: ctx
func (_m *ArbitrageService) StopCalculation(ctx context.Context) error {
	ret := _m.Called(ctx)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewArbitrageService interface {
	mock.TestingT
	Cleanup(func())
}

// NewArbitrageService creates a new instance of ArbitrageService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewArbitrageService(t mockConstructorTestingTNewArbitrageService) *ArbitrageService {
	mock := &ArbitrageService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
