// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	runners "github.com/schoonology/diplomat/runners"
	mock "github.com/stretchr/testify/mock"
)

// ResultsFormatter is an autogenerated mock type for the ResultsFormatter type
type ResultsFormatter struct {
	mock.Mock
}

// FormatAll provides a mock function with given fields: _a0, _a1
func (_m *ResultsFormatter) FormatAll(_a0 chan runners.TestResult, _a1 chan error) chan string {
	ret := _m.Called(_a0, _a1)

	var r0 chan string
	if rf, ok := ret.Get(0).(func(chan runners.TestResult, chan error) chan string); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(chan string)
		}
	}

	return r0
}
