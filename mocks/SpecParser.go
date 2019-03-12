// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import loaders "github.com/testdouble/diplomat/loaders"
import mock "github.com/stretchr/testify/mock"
import parsers "github.com/testdouble/diplomat/parsers"

// SpecParser is an autogenerated mock type for the SpecParser type
type SpecParser struct {
	mock.Mock
}

// Parse provides a mock function with given fields: _a0
func (_m *SpecParser) Parse(_a0 *loaders.Body) (*parsers.Spec, error) {
	ret := _m.Called(_a0)

	var r0 *parsers.Spec
	if rf, ok := ret.Get(0).(func(*loaders.Body) *parsers.Spec); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*parsers.Spec)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*loaders.Body) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Stream provides a mock function with given fields: _a0
func (_m *SpecParser) Stream(_a0 chan *loaders.Body) (chan *parsers.Spec, chan error) {
	ret := _m.Called(_a0)

	var r0 chan *parsers.Spec
	if rf, ok := ret.Get(0).(func(chan *loaders.Body) chan *parsers.Spec); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(chan *parsers.Spec)
		}
	}

	var r1 chan error
	if rf, ok := ret.Get(1).(func(chan *loaders.Body) chan error); ok {
		r1 = rf(_a0)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(chan error)
		}
	}

	return r0, r1
}
