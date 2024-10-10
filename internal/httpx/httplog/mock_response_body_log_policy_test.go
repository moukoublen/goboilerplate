// Code generated by mockery v2.46.2. DO NOT EDIT.

package httplog

import (
	http "net/http"

	mock "github.com/stretchr/testify/mock"
)

// MockResponseBodyLogPolicy is an autogenerated mock type for the ResponseBodyLogPolicy type
type MockResponseBodyLogPolicy struct {
	mock.Mock
}

type MockResponseBodyLogPolicy_Expecter struct {
	mock *mock.Mock
}

func (_m *MockResponseBodyLogPolicy) EXPECT() *MockResponseBodyLogPolicy_Expecter {
	return &MockResponseBodyLogPolicy_Expecter{mock: &_m.Mock}
}

// Execute provides a mock function with given fields: r
func (_m *MockResponseBodyLogPolicy) Execute(r *http.Response) bool {
	ret := _m.Called(r)

	if len(ret) == 0 {
		panic("no return value specified for Execute")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func(*http.Response) bool); ok {
		r0 = rf(r)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// MockResponseBodyLogPolicy_Execute_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Execute'
type MockResponseBodyLogPolicy_Execute_Call struct {
	*mock.Call
}

// Execute is a helper method to define mock.On call
//   - r *http.Response
func (_e *MockResponseBodyLogPolicy_Expecter) Execute(r interface{}) *MockResponseBodyLogPolicy_Execute_Call {
	return &MockResponseBodyLogPolicy_Execute_Call{Call: _e.mock.On("Execute", r)}
}

func (_c *MockResponseBodyLogPolicy_Execute_Call) Run(run func(r *http.Response)) *MockResponseBodyLogPolicy_Execute_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*http.Response))
	})
	return _c
}

func (_c *MockResponseBodyLogPolicy_Execute_Call) Return(_a0 bool) *MockResponseBodyLogPolicy_Execute_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockResponseBodyLogPolicy_Execute_Call) RunAndReturn(run func(*http.Response) bool) *MockResponseBodyLogPolicy_Execute_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockResponseBodyLogPolicy creates a new instance of MockResponseBodyLogPolicy. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockResponseBodyLogPolicy(t interface {
	mock.TestingT
	Cleanup(func())
},
) *MockResponseBodyLogPolicy {
	mock := &MockResponseBodyLogPolicy{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}