// Code generated by mockery v2.46.2. DO NOT EDIT.

package httplog

import mock "github.com/stretchr/testify/mock"

// MockHeaderMatcher is an autogenerated mock type for the HeaderMatcher type
type MockHeaderMatcher struct {
	mock.Mock
}

type MockHeaderMatcher_Expecter struct {
	mock *mock.Mock
}

func (_m *MockHeaderMatcher) EXPECT() *MockHeaderMatcher_Expecter {
	return &MockHeaderMatcher_Expecter{mock: &_m.Mock}
}

// Match provides a mock function with given fields: key, values
func (_m *MockHeaderMatcher) Match(key string, values []string) bool {
	ret := _m.Called(key, values)

	if len(ret) == 0 {
		panic("no return value specified for Match")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func(string, []string) bool); ok {
		r0 = rf(key, values)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// MockHeaderMatcher_Match_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Match'
type MockHeaderMatcher_Match_Call struct {
	*mock.Call
}

// Match is a helper method to define mock.On call
//   - key string
//   - values []string
func (_e *MockHeaderMatcher_Expecter) Match(key interface{}, values interface{}) *MockHeaderMatcher_Match_Call {
	return &MockHeaderMatcher_Match_Call{Call: _e.mock.On("Match", key, values)}
}

func (_c *MockHeaderMatcher_Match_Call) Run(run func(key string, values []string)) *MockHeaderMatcher_Match_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].([]string))
	})
	return _c
}

func (_c *MockHeaderMatcher_Match_Call) Return(_a0 bool) *MockHeaderMatcher_Match_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockHeaderMatcher_Match_Call) RunAndReturn(run func(string, []string) bool) *MockHeaderMatcher_Match_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockHeaderMatcher creates a new instance of MockHeaderMatcher. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockHeaderMatcher(t interface {
	mock.TestingT
	Cleanup(func())
},
) *MockHeaderMatcher {
	mock := &MockHeaderMatcher{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
