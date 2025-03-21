// Code generated by mockery v2.53.3. DO NOT EDIT.

package git

import mock "github.com/stretchr/testify/mock"

// MockGit is an autogenerated mock type for the Git type
type MockGit struct {
	mock.Mock
}

type MockGit_Expecter struct {
	mock *mock.Mock
}

func (_m *MockGit) EXPECT() *MockGit_Expecter {
	return &MockGit_Expecter{mock: &_m.Mock}
}

// AddWorktree provides a mock function with given fields: _a0, _a1, _a2, _a3, _a4
func (_m *MockGit) AddWorktree(_a0 string, _a1 bool, _a2 string, _a3 string, _a4 string) error {
	ret := _m.Called(_a0, _a1, _a2, _a3, _a4)

	if len(ret) == 0 {
		panic("no return value specified for AddWorktree")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, bool, string, string, string) error); ok {
		r0 = rf(_a0, _a1, _a2, _a3, _a4)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockGit_AddWorktree_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AddWorktree'
type MockGit_AddWorktree_Call struct {
	*mock.Call
}

// AddWorktree is a helper method to define mock.On call
//   - _a0 string
//   - _a1 bool
//   - _a2 string
//   - _a3 string
//   - _a4 string
func (_e *MockGit_Expecter) AddWorktree(_a0 interface{}, _a1 interface{}, _a2 interface{}, _a3 interface{}, _a4 interface{}) *MockGit_AddWorktree_Call {
	return &MockGit_AddWorktree_Call{Call: _e.mock.On("AddWorktree", _a0, _a1, _a2, _a3, _a4)}
}

func (_c *MockGit_AddWorktree_Call) Run(run func(_a0 string, _a1 bool, _a2 string, _a3 string, _a4 string)) *MockGit_AddWorktree_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(bool), args[2].(string), args[3].(string), args[4].(string))
	})
	return _c
}

func (_c *MockGit_AddWorktree_Call) Return(_a0 error) *MockGit_AddWorktree_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockGit_AddWorktree_Call) RunAndReturn(run func(string, bool, string, string, string) error) *MockGit_AddWorktree_Call {
	_c.Call.Return(run)
	return _c
}

// Clone provides a mock function with given fields: name
func (_m *MockGit) Clone(name string) (string, error) {
	ret := _m.Called(name)

	if len(ret) == 0 {
		panic("no return value specified for Clone")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (string, error)); ok {
		return rf(name)
	}
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(name)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockGit_Clone_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Clone'
type MockGit_Clone_Call struct {
	*mock.Call
}

// Clone is a helper method to define mock.On call
//   - name string
func (_e *MockGit_Expecter) Clone(name interface{}) *MockGit_Clone_Call {
	return &MockGit_Clone_Call{Call: _e.mock.On("Clone", name)}
}

func (_c *MockGit_Clone_Call) Run(run func(name string)) *MockGit_Clone_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockGit_Clone_Call) Return(_a0 string, _a1 error) *MockGit_Clone_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockGit_Clone_Call) RunAndReturn(run func(string) (string, error)) *MockGit_Clone_Call {
	_c.Call.Return(run)
	return _c
}

// CloneBare provides a mock function with given fields: _a0, _a1
func (_m *MockGit) CloneBare(_a0 string, _a1 string) error {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for CloneBare")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockGit_CloneBare_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CloneBare'
type MockGit_CloneBare_Call struct {
	*mock.Call
}

// CloneBare is a helper method to define mock.On call
//   - _a0 string
//   - _a1 string
func (_e *MockGit_Expecter) CloneBare(_a0 interface{}, _a1 interface{}) *MockGit_CloneBare_Call {
	return &MockGit_CloneBare_Call{Call: _e.mock.On("CloneBare", _a0, _a1)}
}

func (_c *MockGit_CloneBare_Call) Run(run func(_a0 string, _a1 string)) *MockGit_CloneBare_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(string))
	})
	return _c
}

func (_c *MockGit_CloneBare_Call) Return(_a0 error) *MockGit_CloneBare_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockGit_CloneBare_Call) RunAndReturn(run func(string, string) error) *MockGit_CloneBare_Call {
	_c.Call.Return(run)
	return _c
}

// CreateTempBranch provides a mock function with given fields: path
func (_m *MockGit) CreateTempBranch(path string) error {
	ret := _m.Called(path)

	if len(ret) == 0 {
		panic("no return value specified for CreateTempBranch")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(path)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockGit_CreateTempBranch_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateTempBranch'
type MockGit_CreateTempBranch_Call struct {
	*mock.Call
}

// CreateTempBranch is a helper method to define mock.On call
//   - path string
func (_e *MockGit_Expecter) CreateTempBranch(path interface{}) *MockGit_CreateTempBranch_Call {
	return &MockGit_CreateTempBranch_Call{Call: _e.mock.On("CreateTempBranch", path)}
}

func (_c *MockGit_CreateTempBranch_Call) Run(run func(path string)) *MockGit_CreateTempBranch_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockGit_CreateTempBranch_Call) Return(_a0 error) *MockGit_CreateTempBranch_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockGit_CreateTempBranch_Call) RunAndReturn(run func(string) error) *MockGit_CreateTempBranch_Call {
	_c.Call.Return(run)
	return _c
}

// DeleteBranch provides a mock function with given fields: branch, path
func (_m *MockGit) DeleteBranch(branch string, path string) error {
	ret := _m.Called(branch, path)

	if len(ret) == 0 {
		panic("no return value specified for DeleteBranch")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(branch, path)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockGit_DeleteBranch_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteBranch'
type MockGit_DeleteBranch_Call struct {
	*mock.Call
}

// DeleteBranch is a helper method to define mock.On call
//   - branch string
//   - path string
func (_e *MockGit_Expecter) DeleteBranch(branch interface{}, path interface{}) *MockGit_DeleteBranch_Call {
	return &MockGit_DeleteBranch_Call{Call: _e.mock.On("DeleteBranch", branch, path)}
}

func (_c *MockGit_DeleteBranch_Call) Run(run func(branch string, path string)) *MockGit_DeleteBranch_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(string))
	})
	return _c
}

func (_c *MockGit_DeleteBranch_Call) Return(_a0 error) *MockGit_DeleteBranch_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockGit_DeleteBranch_Call) RunAndReturn(run func(string, string) error) *MockGit_DeleteBranch_Call {
	_c.Call.Return(run)
	return _c
}

// FetchOrigin provides a mock function with given fields: branch, path
func (_m *MockGit) FetchOrigin(branch string, path string) error {
	ret := _m.Called(branch, path)

	if len(ret) == 0 {
		panic("no return value specified for FetchOrigin")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(branch, path)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockGit_FetchOrigin_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FetchOrigin'
type MockGit_FetchOrigin_Call struct {
	*mock.Call
}

// FetchOrigin is a helper method to define mock.On call
//   - branch string
//   - path string
func (_e *MockGit_Expecter) FetchOrigin(branch interface{}, path interface{}) *MockGit_FetchOrigin_Call {
	return &MockGit_FetchOrigin_Call{Call: _e.mock.On("FetchOrigin", branch, path)}
}

func (_c *MockGit_FetchOrigin_Call) Run(run func(branch string, path string)) *MockGit_FetchOrigin_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(string))
	})
	return _c
}

func (_c *MockGit_FetchOrigin_Call) Return(_a0 error) *MockGit_FetchOrigin_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockGit_FetchOrigin_Call) RunAndReturn(run func(string, string) error) *MockGit_FetchOrigin_Call {
	_c.Call.Return(run)
	return _c
}

// GetLocalBranches provides a mock function with given fields: _a0
func (_m *MockGit) GetLocalBranches(_a0 string) ([]string, error) {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for GetLocalBranches")
	}

	var r0 []string
	var r1 error
	if rf, ok := ret.Get(0).(func(string) ([]string, error)); ok {
		return rf(_a0)
	}
	if rf, ok := ret.Get(0).(func(string) []string); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockGit_GetLocalBranches_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetLocalBranches'
type MockGit_GetLocalBranches_Call struct {
	*mock.Call
}

// GetLocalBranches is a helper method to define mock.On call
//   - _a0 string
func (_e *MockGit_Expecter) GetLocalBranches(_a0 interface{}) *MockGit_GetLocalBranches_Call {
	return &MockGit_GetLocalBranches_Call{Call: _e.mock.On("GetLocalBranches", _a0)}
}

func (_c *MockGit_GetLocalBranches_Call) Run(run func(_a0 string)) *MockGit_GetLocalBranches_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockGit_GetLocalBranches_Call) Return(_a0 []string, _a1 error) *MockGit_GetLocalBranches_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockGit_GetLocalBranches_Call) RunAndReturn(run func(string) ([]string, error)) *MockGit_GetLocalBranches_Call {
	_c.Call.Return(run)
	return _c
}

// GetRemoteBranches provides a mock function with given fields: _a0
func (_m *MockGit) GetRemoteBranches(_a0 string) ([]string, error) {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for GetRemoteBranches")
	}

	var r0 []string
	var r1 error
	if rf, ok := ret.Get(0).(func(string) ([]string, error)); ok {
		return rf(_a0)
	}
	if rf, ok := ret.Get(0).(func(string) []string); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockGit_GetRemoteBranches_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetRemoteBranches'
type MockGit_GetRemoteBranches_Call struct {
	*mock.Call
}

// GetRemoteBranches is a helper method to define mock.On call
//   - _a0 string
func (_e *MockGit_Expecter) GetRemoteBranches(_a0 interface{}) *MockGit_GetRemoteBranches_Call {
	return &MockGit_GetRemoteBranches_Call{Call: _e.mock.On("GetRemoteBranches", _a0)}
}

func (_c *MockGit_GetRemoteBranches_Call) Run(run func(_a0 string)) *MockGit_GetRemoteBranches_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockGit_GetRemoteBranches_Call) Return(_a0 []string, _a1 error) *MockGit_GetRemoteBranches_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockGit_GetRemoteBranches_Call) RunAndReturn(run func(string) ([]string, error)) *MockGit_GetRemoteBranches_Call {
	_c.Call.Return(run)
	return _c
}

// GetRepoName provides a mock function with given fields: path
func (_m *MockGit) GetRepoName(path string) (string, error) {
	ret := _m.Called(path)

	if len(ret) == 0 {
		panic("no return value specified for GetRepoName")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (string, error)); ok {
		return rf(path)
	}
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(path)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(path)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockGit_GetRepoName_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetRepoName'
type MockGit_GetRepoName_Call struct {
	*mock.Call
}

// GetRepoName is a helper method to define mock.On call
//   - path string
func (_e *MockGit_Expecter) GetRepoName(path interface{}) *MockGit_GetRepoName_Call {
	return &MockGit_GetRepoName_Call{Call: _e.mock.On("GetRepoName", path)}
}

func (_c *MockGit_GetRepoName_Call) Run(run func(path string)) *MockGit_GetRepoName_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockGit_GetRepoName_Call) Return(_a0 string, _a1 error) *MockGit_GetRepoName_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockGit_GetRepoName_Call) RunAndReturn(run func(string) (string, error)) *MockGit_GetRepoName_Call {
	_c.Call.Return(run)
	return _c
}

// GetWorktrees provides a mock function with no fields
func (_m *MockGit) GetWorktrees() ([]string, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetWorktrees")
	}

	var r0 []string
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]string, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []string); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockGit_GetWorktrees_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetWorktrees'
type MockGit_GetWorktrees_Call struct {
	*mock.Call
}

// GetWorktrees is a helper method to define mock.On call
func (_e *MockGit_Expecter) GetWorktrees() *MockGit_GetWorktrees_Call {
	return &MockGit_GetWorktrees_Call{Call: _e.mock.On("GetWorktrees")}
}

func (_c *MockGit_GetWorktrees_Call) Run(run func()) *MockGit_GetWorktrees_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockGit_GetWorktrees_Call) Return(_a0 []string, _a1 error) *MockGit_GetWorktrees_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockGit_GetWorktrees_Call) RunAndReturn(run func() ([]string, error)) *MockGit_GetWorktrees_Call {
	_c.Call.Return(run)
	return _c
}

// GitCommonDir provides a mock function with given fields: name
func (_m *MockGit) GitCommonDir(name string) (bool, string, error) {
	ret := _m.Called(name)

	if len(ret) == 0 {
		panic("no return value specified for GitCommonDir")
	}

	var r0 bool
	var r1 string
	var r2 error
	if rf, ok := ret.Get(0).(func(string) (bool, string, error)); ok {
		return rf(name)
	}
	if rf, ok := ret.Get(0).(func(string) bool); ok {
		r0 = rf(name)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(string) string); ok {
		r1 = rf(name)
	} else {
		r1 = ret.Get(1).(string)
	}

	if rf, ok := ret.Get(2).(func(string) error); ok {
		r2 = rf(name)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// MockGit_GitCommonDir_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GitCommonDir'
type MockGit_GitCommonDir_Call struct {
	*mock.Call
}

// GitCommonDir is a helper method to define mock.On call
//   - name string
func (_e *MockGit_Expecter) GitCommonDir(name interface{}) *MockGit_GitCommonDir_Call {
	return &MockGit_GitCommonDir_Call{Call: _e.mock.On("GitCommonDir", name)}
}

func (_c *MockGit_GitCommonDir_Call) Run(run func(name string)) *MockGit_GitCommonDir_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockGit_GitCommonDir_Call) Return(_a0 bool, _a1 string, _a2 error) *MockGit_GitCommonDir_Call {
	_c.Call.Return(_a0, _a1, _a2)
	return _c
}

func (_c *MockGit_GitCommonDir_Call) RunAndReturn(run func(string) (bool, string, error)) *MockGit_GitCommonDir_Call {
	_c.Call.Return(run)
	return _c
}

// PullBranch provides a mock function with given fields: url
func (_m *MockGit) PullBranch(url string) error {
	ret := _m.Called(url)

	if len(ret) == 0 {
		panic("no return value specified for PullBranch")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(url)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockGit_PullBranch_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'PullBranch'
type MockGit_PullBranch_Call struct {
	*mock.Call
}

// PullBranch is a helper method to define mock.On call
//   - url string
func (_e *MockGit_Expecter) PullBranch(url interface{}) *MockGit_PullBranch_Call {
	return &MockGit_PullBranch_Call{Call: _e.mock.On("PullBranch", url)}
}

func (_c *MockGit_PullBranch_Call) Run(run func(url string)) *MockGit_PullBranch_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockGit_PullBranch_Call) Return(_a0 error) *MockGit_PullBranch_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockGit_PullBranch_Call) RunAndReturn(run func(string) error) *MockGit_PullBranch_Call {
	_c.Call.Return(run)
	return _c
}

// RemoveWorktree provides a mock function with given fields: _a0
func (_m *MockGit) RemoveWorktree(_a0 string) (string, error) {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for RemoveWorktree")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (string, error)); ok {
		return rf(_a0)
	}
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockGit_RemoveWorktree_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RemoveWorktree'
type MockGit_RemoveWorktree_Call struct {
	*mock.Call
}

// RemoveWorktree is a helper method to define mock.On call
//   - _a0 string
func (_e *MockGit_Expecter) RemoveWorktree(_a0 interface{}) *MockGit_RemoveWorktree_Call {
	return &MockGit_RemoveWorktree_Call{Call: _e.mock.On("RemoveWorktree", _a0)}
}

func (_c *MockGit_RemoveWorktree_Call) Run(run func(_a0 string)) *MockGit_RemoveWorktree_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockGit_RemoveWorktree_Call) Return(_a0 string, _a1 error) *MockGit_RemoveWorktree_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockGit_RemoveWorktree_Call) RunAndReturn(run func(string) (string, error)) *MockGit_RemoveWorktree_Call {
	_c.Call.Return(run)
	return _c
}

// ShowTopLevel provides a mock function with given fields: name
func (_m *MockGit) ShowTopLevel(name string) (bool, string, error) {
	ret := _m.Called(name)

	if len(ret) == 0 {
		panic("no return value specified for ShowTopLevel")
	}

	var r0 bool
	var r1 string
	var r2 error
	if rf, ok := ret.Get(0).(func(string) (bool, string, error)); ok {
		return rf(name)
	}
	if rf, ok := ret.Get(0).(func(string) bool); ok {
		r0 = rf(name)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(string) string); ok {
		r1 = rf(name)
	} else {
		r1 = ret.Get(1).(string)
	}

	if rf, ok := ret.Get(2).(func(string) error); ok {
		r2 = rf(name)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// MockGit_ShowTopLevel_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ShowTopLevel'
type MockGit_ShowTopLevel_Call struct {
	*mock.Call
}

// ShowTopLevel is a helper method to define mock.On call
//   - name string
func (_e *MockGit_Expecter) ShowTopLevel(name interface{}) *MockGit_ShowTopLevel_Call {
	return &MockGit_ShowTopLevel_Call{Call: _e.mock.On("ShowTopLevel", name)}
}

func (_c *MockGit_ShowTopLevel_Call) Run(run func(name string)) *MockGit_ShowTopLevel_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockGit_ShowTopLevel_Call) Return(_a0 bool, _a1 string, _a2 error) *MockGit_ShowTopLevel_Call {
	_c.Call.Return(_a0, _a1, _a2)
	return _c
}

func (_c *MockGit_ShowTopLevel_Call) RunAndReturn(run func(string) (bool, string, error)) *MockGit_ShowTopLevel_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockGit creates a new instance of MockGit. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockGit(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockGit {
	mock := &MockGit{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
