// Code generated by mockery v2.20.2. DO NOT EDIT.

package mocks

import (
	context "context"

	models "github.com/AleksandrVishniakov/distributed-calculator/auth/app/internal/domain/models"
	mock "github.com/stretchr/testify/mock"
)

// UserProvider is an autogenerated mock type for the UserProvider type
type UserProvider struct {
	mock.Mock
}

// User provides a mock function with given fields: ctx, login
func (_m *UserProvider) User(ctx context.Context, login string) (*models.User, error) {
	ret := _m.Called(ctx, login)

	var r0 *models.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*models.User, error)); ok {
		return rf(ctx, login)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *models.User); ok {
		r0 = rf(ctx, login)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, login)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewUserProvider interface {
	mock.TestingT
	Cleanup(func())
}

// NewUserProvider creates a new instance of UserProvider. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewUserProvider(t mockConstructorTestingTNewUserProvider) *UserProvider {
	mock := &UserProvider{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
