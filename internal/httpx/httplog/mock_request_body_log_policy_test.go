// Code generated by mockery v2.46.2. DO NOT EDIT.

package httplog

import (
	http "net/http"

	mock "github.com/stretchr/testify/mock"
)

// MockRequestBodyLogPolicy is an autogenerated mock type for the RequestBodyLogPolicy type
type MockRequestBodyLogPolicy struct {
	mock.Mock
}

type MockRequestBodyLogPolicy_Expecter struct {
	mock *mock.Mock
}

func (_m *MockRequestBodyLogPolicy) EXPECT() *MockRequestBodyLogPolicy_Expecter {
	return &MockRequestBodyLogPolicy_Expecter{mock: &_m.Mock}
}

// Execute provides a mock function with given fields: r
func (_m *MockRequestBodyLogPolicy) Execute(r *http.Request) bool {
	ret := _m.Called(r)

	if len(ret) == 0 {
		panic("no return value specified for Execute")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func(*http.Request) bool); ok {
		r0 = rf(r)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// MockRequestBodyLogPolicy_Execute_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Execute'
type MockRequestBodyLogPolicy_Execute_Call struct {
	*mock.Call
}

// Execute is a helper method to define mock.On call
//   - r *http.Request
func (_e *MockRequestBodyLogPolicy_Expecter) Execute(r interface{}) *MockRequestBodyLogPolicy_Execute_Call {
	return &MockRequestBodyLogPolicy_Execute_Call{Call: _e.mock.On("Execute", r)}
}

func (_c *MockRequestBodyLogPolicy_Execute_Call) Run(run func(r *http.Request)) *MockRequestBodyLogPolicy_Execute_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*http.Request))
	})
	return _c
}

func (_c *MockRequestBodyLogPolicy_Execute_Call) Return(_a0 bool) *MockRequestBodyLogPolicy_Execute_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockRequestBodyLogPolicy_Execute_Call) RunAndReturn(run func(*http.Request) bool) *MockRequestBodyLogPolicy_Execute_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockRequestBodyLogPolicy creates a new instance of MockRequestBodyLogPolicy. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockRequestBodyLogPolicy(t interface {
	mock.TestingT
	Cleanup(func())
},
) *MockRequestBodyLogPolicy {
	mock := &MockRequestBodyLogPolicy{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
