// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	loaders "github.com/schoonology/diplomat/loaders"
	mock "github.com/stretchr/testify/mock"

	parsers "github.com/schoonology/diplomat/parsers"
)

// ParseDelegator is an autogenerated mock type for the ParseDelegator type
type ParseDelegator struct {
	mock.Mock
}

// ParseAll provides a mock function with given fields: _a0
func (_m *ParseDelegator) ParseAll(_a0 chan loaders.File) chan parsers.Paragraph {
	ret := _m.Called(_a0)

	var r0 chan parsers.Paragraph
	if rf, ok := ret.Get(0).(func(chan loaders.File) chan parsers.Paragraph); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(chan parsers.Paragraph)
		}
	}

	return r0
}
