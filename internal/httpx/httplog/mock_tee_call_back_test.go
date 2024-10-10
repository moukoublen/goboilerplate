// Code generated by mockery v2.46.2. DO NOT EDIT.

package httplog

import (
	bytes "bytes"

	mock "github.com/stretchr/testify/mock"
)

// MockTeeCallBack is an autogenerated mock type for the TeeCallBack type
type MockTeeCallBack struct {
	mock.Mock
}

type MockTeeCallBack_Expecter struct {
	mock *mock.Mock
}

func (_m *MockTeeCallBack) EXPECT() *MockTeeCallBack_Expecter {
	return &MockTeeCallBack_Expecter{mock: &_m.Mock}
}

// Execute provides a mock function with given fields: readErr, closeErr, buf
func (_m *MockTeeCallBack) Execute(readErr error, closeErr error, buf *bytes.Buffer) {
	_m.Called(readErr, closeErr, buf)
}

// MockTeeCallBack_Execute_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Execute'
type MockTeeCallBack_Execute_Call struct {
	*mock.Call
}

// Execute is a helper method to define mock.On call
//   - readErr error
//   - closeErr error
//   - buf *bytes.Buffer
func (_e *MockTeeCallBack_Expecter) Execute(readErr interface{}, closeErr interface{}, buf interface{}) *MockTeeCallBack_Execute_Call {
	return &MockTeeCallBack_Execute_Call{Call: _e.mock.On("Execute", readErr, closeErr, buf)}
}

func (_c *MockTeeCallBack_Execute_Call) Run(run func(readErr error, closeErr error, buf *bytes.Buffer)) *MockTeeCallBack_Execute_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(error), args[1].(error), args[2].(*bytes.Buffer))
	})
	return _c
}

func (_c *MockTeeCallBack_Execute_Call) Return() *MockTeeCallBack_Execute_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockTeeCallBack_Execute_Call) RunAndReturn(run func(error, error, *bytes.Buffer)) *MockTeeCallBack_Execute_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockTeeCallBack creates a new instance of MockTeeCallBack. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockTeeCallBack(t interface {
	mock.TestingT
	Cleanup(func())
},
) *MockTeeCallBack {
	mock := &MockTeeCallBack{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}