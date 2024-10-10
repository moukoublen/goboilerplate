// Code generated by mockery v2.46.2. DO NOT EDIT.

package httplog

import (
	bytes "bytes"
	http "net/http"

	mock "github.com/stretchr/testify/mock"
)

// MockResponseWriterWrapper is an autogenerated mock type for the ResponseWriterWrapper type
type MockResponseWriterWrapper struct {
	mock.Mock
}

type MockResponseWriterWrapper_Expecter struct {
	mock *mock.Mock
}

func (_m *MockResponseWriterWrapper) EXPECT() *MockResponseWriterWrapper_Expecter {
	return &MockResponseWriterWrapper_Expecter{mock: &_m.Mock}
}

// Buffer provides a mock function with given fields:
func (_m *MockResponseWriterWrapper) Buffer() *bytes.Buffer {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Buffer")
	}

	var r0 *bytes.Buffer
	if rf, ok := ret.Get(0).(func() *bytes.Buffer); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*bytes.Buffer)
		}
	}

	return r0
}

// MockResponseWriterWrapper_Buffer_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Buffer'
type MockResponseWriterWrapper_Buffer_Call struct {
	*mock.Call
}

// Buffer is a helper method to define mock.On call
func (_e *MockResponseWriterWrapper_Expecter) Buffer() *MockResponseWriterWrapper_Buffer_Call {
	return &MockResponseWriterWrapper_Buffer_Call{Call: _e.mock.On("Buffer")}
}

func (_c *MockResponseWriterWrapper_Buffer_Call) Run(run func()) *MockResponseWriterWrapper_Buffer_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockResponseWriterWrapper_Buffer_Call) Return(_a0 *bytes.Buffer) *MockResponseWriterWrapper_Buffer_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockResponseWriterWrapper_Buffer_Call) RunAndReturn(run func() *bytes.Buffer) *MockResponseWriterWrapper_Buffer_Call {
	_c.Call.Return(run)
	return _c
}

// BytesWritten provides a mock function with given fields:
func (_m *MockResponseWriterWrapper) BytesWritten() int {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for BytesWritten")
	}

	var r0 int
	if rf, ok := ret.Get(0).(func() int); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int)
	}

	return r0
}

// MockResponseWriterWrapper_BytesWritten_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'BytesWritten'
type MockResponseWriterWrapper_BytesWritten_Call struct {
	*mock.Call
}

// BytesWritten is a helper method to define mock.On call
func (_e *MockResponseWriterWrapper_Expecter) BytesWritten() *MockResponseWriterWrapper_BytesWritten_Call {
	return &MockResponseWriterWrapper_BytesWritten_Call{Call: _e.mock.On("BytesWritten")}
}

func (_c *MockResponseWriterWrapper_BytesWritten_Call) Run(run func()) *MockResponseWriterWrapper_BytesWritten_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockResponseWriterWrapper_BytesWritten_Call) Return(_a0 int) *MockResponseWriterWrapper_BytesWritten_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockResponseWriterWrapper_BytesWritten_Call) RunAndReturn(run func() int) *MockResponseWriterWrapper_BytesWritten_Call {
	_c.Call.Return(run)
	return _c
}

// Header provides a mock function with given fields:
func (_m *MockResponseWriterWrapper) Header() http.Header {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Header")
	}

	var r0 http.Header
	if rf, ok := ret.Get(0).(func() http.Header); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(http.Header)
		}
	}

	return r0
}

// MockResponseWriterWrapper_Header_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Header'
type MockResponseWriterWrapper_Header_Call struct {
	*mock.Call
}

// Header is a helper method to define mock.On call
func (_e *MockResponseWriterWrapper_Expecter) Header() *MockResponseWriterWrapper_Header_Call {
	return &MockResponseWriterWrapper_Header_Call{Call: _e.mock.On("Header")}
}

func (_c *MockResponseWriterWrapper_Header_Call) Run(run func()) *MockResponseWriterWrapper_Header_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockResponseWriterWrapper_Header_Call) Return(_a0 http.Header) *MockResponseWriterWrapper_Header_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockResponseWriterWrapper_Header_Call) RunAndReturn(run func() http.Header) *MockResponseWriterWrapper_Header_Call {
	_c.Call.Return(run)
	return _c
}

// Status provides a mock function with given fields:
func (_m *MockResponseWriterWrapper) Status() int {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Status")
	}

	var r0 int
	if rf, ok := ret.Get(0).(func() int); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int)
	}

	return r0
}

// MockResponseWriterWrapper_Status_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Status'
type MockResponseWriterWrapper_Status_Call struct {
	*mock.Call
}

// Status is a helper method to define mock.On call
func (_e *MockResponseWriterWrapper_Expecter) Status() *MockResponseWriterWrapper_Status_Call {
	return &MockResponseWriterWrapper_Status_Call{Call: _e.mock.On("Status")}
}

func (_c *MockResponseWriterWrapper_Status_Call) Run(run func()) *MockResponseWriterWrapper_Status_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockResponseWriterWrapper_Status_Call) Return(_a0 int) *MockResponseWriterWrapper_Status_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockResponseWriterWrapper_Status_Call) RunAndReturn(run func() int) *MockResponseWriterWrapper_Status_Call {
	_c.Call.Return(run)
	return _c
}

// Write provides a mock function with given fields: _a0
func (_m *MockResponseWriterWrapper) Write(_a0 []byte) (int, error) {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for Write")
	}

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func([]byte) (int, error)); ok {
		return rf(_a0)
	}
	if rf, ok := ret.Get(0).(func([]byte) int); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func([]byte) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockResponseWriterWrapper_Write_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Write'
type MockResponseWriterWrapper_Write_Call struct {
	*mock.Call
}

// Write is a helper method to define mock.On call
//   - _a0 []byte
func (_e *MockResponseWriterWrapper_Expecter) Write(_a0 interface{}) *MockResponseWriterWrapper_Write_Call {
	return &MockResponseWriterWrapper_Write_Call{Call: _e.mock.On("Write", _a0)}
}

func (_c *MockResponseWriterWrapper_Write_Call) Run(run func(_a0 []byte)) *MockResponseWriterWrapper_Write_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].([]byte))
	})
	return _c
}

func (_c *MockResponseWriterWrapper_Write_Call) Return(_a0 int, _a1 error) *MockResponseWriterWrapper_Write_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockResponseWriterWrapper_Write_Call) RunAndReturn(run func([]byte) (int, error)) *MockResponseWriterWrapper_Write_Call {
	_c.Call.Return(run)
	return _c
}

// WriteHeader provides a mock function with given fields: statusCode
func (_m *MockResponseWriterWrapper) WriteHeader(statusCode int) {
	_m.Called(statusCode)
}

// MockResponseWriterWrapper_WriteHeader_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'WriteHeader'
type MockResponseWriterWrapper_WriteHeader_Call struct {
	*mock.Call
}

// WriteHeader is a helper method to define mock.On call
//   - statusCode int
func (_e *MockResponseWriterWrapper_Expecter) WriteHeader(statusCode interface{}) *MockResponseWriterWrapper_WriteHeader_Call {
	return &MockResponseWriterWrapper_WriteHeader_Call{Call: _e.mock.On("WriteHeader", statusCode)}
}

func (_c *MockResponseWriterWrapper_WriteHeader_Call) Run(run func(statusCode int)) *MockResponseWriterWrapper_WriteHeader_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(int))
	})
	return _c
}

func (_c *MockResponseWriterWrapper_WriteHeader_Call) Return() *MockResponseWriterWrapper_WriteHeader_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockResponseWriterWrapper_WriteHeader_Call) RunAndReturn(run func(int)) *MockResponseWriterWrapper_WriteHeader_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockResponseWriterWrapper creates a new instance of MockResponseWriterWrapper. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockResponseWriterWrapper(t interface {
	mock.TestingT
	Cleanup(func())
},
) *MockResponseWriterWrapper {
	mock := &MockResponseWriterWrapper{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}