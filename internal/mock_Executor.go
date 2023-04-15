// Code generated by mockery v2.23.2. DO NOT EDIT.

package internal

import mock "github.com/stretchr/testify/mock"

// MockExecutor is an autogenerated mock type for the Executor type
type MockExecutor struct {
	mock.Mock
}

type MockExecutor_Expecter struct {
	mock *mock.Mock
}

func (_m *MockExecutor) EXPECT() *MockExecutor_Expecter {
	return &MockExecutor_Expecter{mock: &_m.Mock}
}

// Execute provides a mock function with given fields: call
func (_m *MockExecutor) Execute(call Call) (*ExecuteResult, error) {
	ret := _m.Called(call)

	var r0 *ExecuteResult
	var r1 error
	if rf, ok := ret.Get(0).(func(Call) (*ExecuteResult, error)); ok {
		return rf(call)
	}
	if rf, ok := ret.Get(0).(func(Call) *ExecuteResult); ok {
		r0 = rf(call)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ExecuteResult)
		}
	}

	if rf, ok := ret.Get(1).(func(Call) error); ok {
		r1 = rf(call)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockExecutor_Execute_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Execute'
type MockExecutor_Execute_Call struct {
	*mock.Call
}

// Execute is a helper method to define mock.On call
//   - call Call
func (_e *MockExecutor_Expecter) Execute(call interface{}) *MockExecutor_Execute_Call {
	return &MockExecutor_Execute_Call{Call: _e.mock.On("Execute", call)}
}

func (_c *MockExecutor_Execute_Call) Run(run func(call Call)) *MockExecutor_Execute_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(Call))
	})
	return _c
}

func (_c *MockExecutor_Execute_Call) Return(_a0 *ExecuteResult, _a1 error) *MockExecutor_Execute_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockExecutor_Execute_Call) RunAndReturn(run func(Call) (*ExecuteResult, error)) *MockExecutor_Execute_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewMockExecutor interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockExecutor creates a new instance of MockExecutor. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockExecutor(t mockConstructorTestingTNewMockExecutor) *MockExecutor {
	mock := &MockExecutor{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
