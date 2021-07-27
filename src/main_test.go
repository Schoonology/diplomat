package main_test

import (
	"testing"

	main "github.com/schoonology/diplomat"
	"github.com/schoonology/diplomat/builders"
	"github.com/schoonology/diplomat/loaders"
	"github.com/schoonology/diplomat/mocks"
	"github.com/schoonology/diplomat/parsers"
	"github.com/schoonology/diplomat/runners"
	"github.com/schoonology/diplomat/transforms"
)

func TestEngineStart(t *testing.T) {
	loader := &mocks.Loader{}
	parser := &mocks.ParseDelegator{}
	builder := &mocks.SpecBuilder{}
	runner := &mocks.SpecRunner{}
	formatter := &mocks.ResultsFormatter{}
	logger := &mocks.Logger{}
	firstTransformer := &mocks.Transformer{}
	secondTransformer := &mocks.Transformer{}

	filenameChannel := make(chan string)
	errorChannel := make(chan error)
	fileChannel := make(chan loaders.File)
	paragraphChannel := make(chan parsers.Paragraph)
	testChannel := make(chan builders.Test)
	firstTransformerChannel := make(chan builders.Test)
	secondTransformerChannel := make(chan builders.Test)
	resultChannel := make(chan runners.TestResult)
	outputChannel := make(chan string)

	loader.On("LoadAll", filenameChannel, errorChannel).Return(fileChannel)
	parser.On("ParseAll", fileChannel).Return(paragraphChannel)
	builder.On("BuildAll", paragraphChannel).Return(testChannel)
	firstTransformer.On("TransformAll", testChannel).Return(firstTransformerChannel)
	secondTransformer.On("TransformAll", firstTransformerChannel).Return(secondTransformerChannel)
	runner.On("RunAll", secondTransformerChannel).Return(resultChannel)
	formatter.On("FormatAll", resultChannel, errorChannel).Return(outputChannel)
	logger.On("PrintAll", outputChannel, errorChannel).Return()

	subject := main.Engine{
		Loader:  loader,
		Parser:  parser,
		Builder: builder,
		Transforms: []transforms.Transformer{
			firstTransformer,
			secondTransformer,
		},
		Runner:    runner,
		Formatter: formatter,
		Logger:    logger,
	}

	subject.Start(filenameChannel, errorChannel)

	loader.AssertExpectations(t)
	parser.AssertExpectations(t)
	builder.AssertExpectations(t)
	firstTransformer.AssertExpectations(t)
	secondTransformer.AssertExpectations(t)
	runner.AssertExpectations(t)
	formatter.AssertExpectations(t)
	logger.AssertExpectations(t)
}
