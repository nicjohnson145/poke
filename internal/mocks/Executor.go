// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	internal "github.com/nicjohnson145/poke/internal"
	mock "github.com/stretchr/testify/mock"
)

// Executor is an autogenerated mock type for the Executor type
type Executor struct {
	mock.Mock
}

type Executor_Expecter struct {
	mock *mock.Mock
}

func (_m *Executor) EXPECT() *Executor_Expecter {
	return &Executor_Expecter{mock: &_m.Mock}
}

// Execute provides a mock function with given fields: call
func (_m *Executor) Execute(call internal.Call) (*internal.ExecuteResult, error) {
	ret := _m.Called(call)

	var r0 *internal.ExecuteResult
	var r1 error
	if rf, ok := ret.Get(0).(func(internal.Call) (*internal.ExecuteResult, error)); ok {
		return rf(call)
	}
	if rf, ok := ret.Get(0).(func(internal.Call) *internal.ExecuteResult); ok {
		r0 = rf(call)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*internal.ExecuteResult)
		}
	}

	if rf, ok := ret.Get(1).(func(internal.Call) error); ok {
		r1 = rf(call)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Executor_Execute_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Execute'
type Executor_Execute_Call struct {
	*mock.Call
}

// Execute is a helper method to define mock.On call
//   - call internal.Call
func (_e *Executor_Expecter) Execute(call interface{}) *Executor_Execute_Call {
	return &Executor_Execute_Call{Call: _e.mock.On("Execute", call)}
}

func (_c *Executor_Execute_Call) Run(run func(call internal.Call)) *Executor_Execute_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(internal.Call))
	})
	return _c
}

func (_c *Executor_Execute_Call) Return(_a0 *internal.ExecuteResult, _a1 error) *Executor_Execute_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Executor_Execute_Call) RunAndReturn(run func(internal.Call) (*internal.ExecuteResult, error)) *Executor_Execute_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewExecutor interface {
	mock.TestingT
	Cleanup(func())
}

// NewExecutor creates a new instance of Executor. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewExecutor(t mockConstructorTestingTNewExecutor) *Executor {
	mock := &Executor{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
