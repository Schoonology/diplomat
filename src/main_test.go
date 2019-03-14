package main_test

import (
	"testing"

	main "github.com/testdouble/diplomat"
	"github.com/testdouble/diplomat/builders"
	"github.com/testdouble/diplomat/mocks"
	"github.com/testdouble/diplomat/parsers"
	"github.com/testdouble/diplomat/runners"
	"github.com/testdouble/diplomat/transforms"
)

func TestEngineStart(t *testing.T) {
	loader := &mocks.Loader{}
	parser := &mocks.SpecParser{}
	builder := &mocks.SpecBuilder{}
	runner := &mocks.SpecRunner{}
	printer := &mocks.ResultsPrinter{}
	firstTransformer := &mocks.Transformer{}
	secondTransformer := &mocks.Transformer{}

	errorChannel := make(chan error)
	bodyChannel := make(chan string)
	specChannel := make(chan parsers.Spec)
	testChannel := make(chan builders.Test)
	firstTransformerChannel := make(chan builders.Test)
	secondTransformerChannel := make(chan builders.Test)
	resultChannel := make(chan runners.TestResult)

	loader.On("Load", "test-file", errorChannel).Return(bodyChannel)
	parser.On("Parse", bodyChannel, errorChannel).Return(specChannel)
	builder.On("BuildAll", specChannel, errorChannel).Return(testChannel)
	firstTransformer.On("TransformAll", testChannel, errorChannel).Return(firstTransformerChannel)
	secondTransformer.On("TransformAll", firstTransformerChannel, errorChannel).Return(secondTransformerChannel)
	runner.On("RunAll", secondTransformerChannel, errorChannel).Return(resultChannel)
	printer.On("Print", resultChannel, errorChannel).Return()

	subject := main.Engine{
		Loader:  loader,
		Parser:  parser,
		Builder: builder,
		Transforms: []transforms.Transformer{
			firstTransformer,
			secondTransformer,
		},
		Runner:  runner,
		Printer: printer,
	}

	subject.Start("test-file", errorChannel)

	loader.AssertExpectations(t)
	parser.AssertExpectations(t)
	runner.AssertExpectations(t)
	printer.AssertExpectations(t)
}
