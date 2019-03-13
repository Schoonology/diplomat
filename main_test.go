package main_test

import (
	"testing"

	main "github.com/testdouble/diplomat"
	"github.com/testdouble/diplomat/mocks"
	"github.com/testdouble/diplomat/parsers"
	"github.com/testdouble/diplomat/runners"
	"github.com/testdouble/diplomat/transforms"
)

func TestEngineStart(t *testing.T) {
	loader := &mocks.Loader{}
	parser := &mocks.SpecParser{}
	runner := &mocks.SpecRunner{}
	printer := &mocks.ResultsPrinter{}
	transforms := []transforms.Transform{}

	errorChannel := make(chan error)
	bodyChannel := make(chan string)
	testChannel := make(chan parsers.Test)
	resultChannel := make(chan runners.TestResult)

	loader.On("Load", "test-file", errorChannel).Return(bodyChannel)
	parser.On("Parse", bodyChannel, errorChannel).Return(testChannel)
	runner.On("Run", testChannel, errorChannel).Return(resultChannel)
	printer.On("Print", resultChannel, errorChannel).Return()

	subject := main.Engine{
		Loader:     loader,
		Parser:     parser,
		Transforms: transforms,
		Runner:     runner,
		Printer:    printer,
	}

	subject.Start("test-file", errorChannel)

	loader.AssertExpectations(t)
	parser.AssertExpectations(t)
	runner.AssertExpectations(t)
	printer.AssertExpectations(t)
}
