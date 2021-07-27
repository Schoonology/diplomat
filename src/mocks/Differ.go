// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	http "github.com/schoonology/diplomat/http"
	mock "github.com/stretchr/testify/mock"
)

// Differ is an autogenerated mock type for the Differ type
type Differ struct {
	mock.Mock
}

// Diff provides a mock function with given fields: _a0, _a1
func (_m *Differ) Diff(_a0 *http.Response, _a1 *http.Response) (string, error) {
	ret := _m.Called(_a0, _a1)

	var r0 string
	if rf, ok := ret.Get(0).(func(*http.Response, *http.Response) string); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*http.Response, *http.Response) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
