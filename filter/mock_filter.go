// Code generated by mockery v2.53.3. DO NOT EDIT.

package filter

import (
	worktreeobj "github.com/garrettkrohn/treekanga/worktreeObj"
	mock "github.com/stretchr/testify/mock"
)

// MockFilter is an autogenerated mock type for the Filter type
type MockFilter struct {
	mock.Mock
}

type MockFilter_Expecter struct {
	mock *mock.Mock
}

func (_m *MockFilter) EXPECT() *MockFilter_Expecter {
	return &MockFilter_Expecter{mock: &_m.Mock}
}

// GetBranchMatchList provides a mock function with given fields: _a0, _a1
func (_m *MockFilter) GetBranchMatchList(_a0 []string, _a1 []worktreeobj.WorktreeObj) []worktreeobj.WorktreeObj {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for GetBranchMatchList")
	}

	var r0 []worktreeobj.WorktreeObj
	if rf, ok := ret.Get(0).(func([]string, []worktreeobj.WorktreeObj) []worktreeobj.WorktreeObj); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]worktreeobj.WorktreeObj)
		}
	}

	return r0
}

// MockFilter_GetBranchMatchList_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetBranchMatchList'
type MockFilter_GetBranchMatchList_Call struct {
	*mock.Call
}

// GetBranchMatchList is a helper method to define mock.On call
//   - _a0 []string
//   - _a1 []worktreeobj.WorktreeObj
func (_e *MockFilter_Expecter) GetBranchMatchList(_a0 interface{}, _a1 interface{}) *MockFilter_GetBranchMatchList_Call {
	return &MockFilter_GetBranchMatchList_Call{Call: _e.mock.On("GetBranchMatchList", _a0, _a1)}
}

func (_c *MockFilter_GetBranchMatchList_Call) Run(run func(_a0 []string, _a1 []worktreeobj.WorktreeObj)) *MockFilter_GetBranchMatchList_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].([]string), args[1].([]worktreeobj.WorktreeObj))
	})
	return _c
}

func (_c *MockFilter_GetBranchMatchList_Call) Return(_a0 []worktreeobj.WorktreeObj) *MockFilter_GetBranchMatchList_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockFilter_GetBranchMatchList_Call) RunAndReturn(run func([]string, []worktreeobj.WorktreeObj) []worktreeobj.WorktreeObj) *MockFilter_GetBranchMatchList_Call {
	_c.Call.Return(run)
	return _c
}

// GetBranchNoMatchList provides a mock function with given fields: _a0, _a1
func (_m *MockFilter) GetBranchNoMatchList(_a0 []string, _a1 []worktreeobj.WorktreeObj) []worktreeobj.WorktreeObj {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for GetBranchNoMatchList")
	}

	var r0 []worktreeobj.WorktreeObj
	if rf, ok := ret.Get(0).(func([]string, []worktreeobj.WorktreeObj) []worktreeobj.WorktreeObj); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]worktreeobj.WorktreeObj)
		}
	}

	return r0
}

// MockFilter_GetBranchNoMatchList_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetBranchNoMatchList'
type MockFilter_GetBranchNoMatchList_Call struct {
	*mock.Call
}

// GetBranchNoMatchList is a helper method to define mock.On call
//   - _a0 []string
//   - _a1 []worktreeobj.WorktreeObj
func (_e *MockFilter_Expecter) GetBranchNoMatchList(_a0 interface{}, _a1 interface{}) *MockFilter_GetBranchNoMatchList_Call {
	return &MockFilter_GetBranchNoMatchList_Call{Call: _e.mock.On("GetBranchNoMatchList", _a0, _a1)}
}

func (_c *MockFilter_GetBranchNoMatchList_Call) Run(run func(_a0 []string, _a1 []worktreeobj.WorktreeObj)) *MockFilter_GetBranchNoMatchList_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].([]string), args[1].([]worktreeobj.WorktreeObj))
	})
	return _c
}

func (_c *MockFilter_GetBranchNoMatchList_Call) Return(_a0 []worktreeobj.WorktreeObj) *MockFilter_GetBranchNoMatchList_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockFilter_GetBranchNoMatchList_Call) RunAndReturn(run func([]string, []worktreeobj.WorktreeObj) []worktreeobj.WorktreeObj) *MockFilter_GetBranchNoMatchList_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockFilter creates a new instance of MockFilter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockFilter(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockFilter {
	mock := &MockFilter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
