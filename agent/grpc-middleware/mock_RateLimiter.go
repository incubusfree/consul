// Code generated by mockery v2.12.0. DO NOT EDIT.

package middleware

import (
	testing "testing"

	rate "github.com/hashicorp/consul/agent/consul/rate"
	mock "github.com/stretchr/testify/mock"
)

// MockRateLimiter is an autogenerated mock type for the RateLimiter type
type MockRateLimiter struct {
	mock.Mock
}

// Allow provides a mock function with given fields: _a0
func (_m *MockRateLimiter) Allow(_a0 rate.Operation) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(rate.Operation) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewMockRateLimiter creates a new instance of MockRateLimiter. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockRateLimiter(t testing.TB) *MockRateLimiter {
	mock := &MockRateLimiter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
