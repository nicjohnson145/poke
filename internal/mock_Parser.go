// Code generated by mockery v2.23.2. DO NOT EDIT.

package internal

import mock "github.com/stretchr/testify/mock"

// MockParser is an autogenerated mock type for the Parser type
type MockParser struct {
	mock.Mock
}

type MockParser_Expecter struct {
	mock *mock.Mock
}

func (_m *MockParser) EXPECT() *MockParser_Expecter {
	return &MockParser_Expecter{mock: &_m.Mock}
}

// Parse provides a mock function with given fields: path
func (_m *MockParser) Parse(path string) (SequenceMap, error) {
	ret := _m.Called(path)

	var r0 SequenceMap
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (SequenceMap, error)); ok {
		return rf(path)
	}
	if rf, ok := ret.Get(0).(func(string) SequenceMap); ok {
		r0 = rf(path)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(SequenceMap)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(path)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockParser_Parse_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Parse'
type MockParser_Parse_Call struct {
	*mock.Call
}

// Parse is a helper method to define mock.On call
//   - path string
func (_e *MockParser_Expecter) Parse(path interface{}) *MockParser_Parse_Call {
	return &MockParser_Parse_Call{Call: _e.mock.On("Parse", path)}
}

func (_c *MockParser_Parse_Call) Run(run func(path string)) *MockParser_Parse_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockParser_Parse_Call) Return(_a0 SequenceMap, _a1 error) *MockParser_Parse_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockParser_Parse_Call) RunAndReturn(run func(string) (SequenceMap, error)) *MockParser_Parse_Call {
	_c.Call.Return(run)
	return _c
}

// ParseSingleSequence provides a mock function with given fields: path
func (_m *MockParser) ParseSingleSequence(path string) (Sequence, error) {
	ret := _m.Called(path)

	var r0 Sequence
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (Sequence, error)); ok {
		return rf(path)
	}
	if rf, ok := ret.Get(0).(func(string) Sequence); ok {
		r0 = rf(path)
	} else {
		r0 = ret.Get(0).(Sequence)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(path)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockParser_ParseSingleSequence_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ParseSingleSequence'
type MockParser_ParseSingleSequence_Call struct {
	*mock.Call
}

// ParseSingleSequence is a helper method to define mock.On call
//   - path string
func (_e *MockParser_Expecter) ParseSingleSequence(path interface{}) *MockParser_ParseSingleSequence_Call {
	return &MockParser_ParseSingleSequence_Call{Call: _e.mock.On("ParseSingleSequence", path)}
}

func (_c *MockParser_ParseSingleSequence_Call) Run(run func(path string)) *MockParser_ParseSingleSequence_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockParser_ParseSingleSequence_Call) Return(_a0 Sequence, _a1 error) *MockParser_ParseSingleSequence_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockParser_ParseSingleSequence_Call) RunAndReturn(run func(string) (Sequence, error)) *MockParser_ParseSingleSequence_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewMockParser interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockParser creates a new instance of MockParser. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockParser(t mockConstructorTestingTNewMockParser) *MockParser {
	mock := &MockParser{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
