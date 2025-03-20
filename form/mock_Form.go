// Code generated by mockery v2.53.3. DO NOT EDIT.

package form

import mock "github.com/stretchr/testify/mock"

// MockForm is an autogenerated mock type for the Form type
type MockForm struct {
	mock.Mock
}

type MockForm_Expecter struct {
	mock *mock.Mock
}

func (_m *MockForm) EXPECT() *MockForm_Expecter {
	return &MockForm_Expecter{mock: &_m.Mock}
}

// Run provides a mock function with no fields
func (_m *MockForm) Run() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Run")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockForm_Run_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Run'
type MockForm_Run_Call struct {
	*mock.Call
}

// Run is a helper method to define mock.On call
func (_e *MockForm_Expecter) Run() *MockForm_Run_Call {
	return &MockForm_Run_Call{Call: _e.mock.On("Run")}
}

func (_c *MockForm_Run_Call) Run(run func()) *MockForm_Run_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockForm_Run_Call) Return(_a0 error) *MockForm_Run_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockForm_Run_Call) RunAndReturn(run func() error) *MockForm_Run_Call {
	_c.Call.Return(run)
	return _c
}

// SetOptions provides a mock function with given fields: stringOptions
func (_m *MockForm) SetOptions(stringOptions []string) {
	_m.Called(stringOptions)
}

// MockForm_SetOptions_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetOptions'
type MockForm_SetOptions_Call struct {
	*mock.Call
}

// SetOptions is a helper method to define mock.On call
//   - stringOptions []string
func (_e *MockForm_Expecter) SetOptions(stringOptions interface{}) *MockForm_SetOptions_Call {
	return &MockForm_SetOptions_Call{Call: _e.mock.On("SetOptions", stringOptions)}
}

func (_c *MockForm_SetOptions_Call) Run(run func(stringOptions []string)) *MockForm_SetOptions_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].([]string))
	})
	return _c
}

func (_c *MockForm_SetOptions_Call) Return() *MockForm_SetOptions_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockForm_SetOptions_Call) RunAndReturn(run func([]string)) *MockForm_SetOptions_Call {
	_c.Run(run)
	return _c
}

// SetSelections provides a mock function with given fields: selections
func (_m *MockForm) SetSelections(selections *[]string) {
	_m.Called(selections)
}

// MockForm_SetSelections_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetSelections'
type MockForm_SetSelections_Call struct {
	*mock.Call
}

// SetSelections is a helper method to define mock.On call
//   - selections *[]string
func (_e *MockForm_Expecter) SetSelections(selections interface{}) *MockForm_SetSelections_Call {
	return &MockForm_SetSelections_Call{Call: _e.mock.On("SetSelections", selections)}
}

func (_c *MockForm_SetSelections_Call) Run(run func(selections *[]string)) *MockForm_SetSelections_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*[]string))
	})
	return _c
}

func (_c *MockForm_SetSelections_Call) Return() *MockForm_SetSelections_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockForm_SetSelections_Call) RunAndReturn(run func(*[]string)) *MockForm_SetSelections_Call {
	_c.Run(run)
	return _c
}

// NewMockForm creates a new instance of MockForm. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockForm(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockForm {
	mock := &MockForm{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
