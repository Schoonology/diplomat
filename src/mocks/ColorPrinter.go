// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import colors "github.com/testdouble/diplomat/colors"
import mock "github.com/stretchr/testify/mock"

// ColorPrinter is an autogenerated mock type for the ColorPrinter type
type ColorPrinter struct {
	mock.Mock
}

// Print provides a mock function with given fields: str, color
func (_m *ColorPrinter) Print(str string, color colors.Color) {
	_m.Called(str, color)
}
