// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	loaders "github.com/schoonology/diplomat/loaders"
	mock "github.com/stretchr/testify/mock"
)

// Loader is an autogenerated mock type for the Loader type
type Loader struct {
	mock.Mock
}

// Load provides a mock function with given fields: _a0
func (_m *Loader) Load(_a0 string) ([]byte, error) {
	ret := _m.Called(_a0)

	var r0 []byte
	if rf, ok := ret.Get(0).(func(string) []byte); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// LoadAll provides a mock function with given fields: _a0, _a1
func (_m *Loader) LoadAll(_a0 chan string, _a1 chan error) chan loaders.File {
	ret := _m.Called(_a0, _a1)

	var r0 chan loaders.File
	if rf, ok := ret.Get(0).(func(chan string, chan error) chan loaders.File); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(chan loaders.File)
		}
	}

	return r0
}
