package main_test

import (
	"testing"

	main "github.com/testdouble/http-assertion-tool"
	"github.com/testdouble/http-assertion-tool/loaders"
	"github.com/testdouble/http-assertion-tool/mocks"
	"github.com/testdouble/http-assertion-tool/parsers"
	"github.com/testdouble/http-assertion-tool/runners"
)

func TestEngineStart(t *testing.T) {
	loader := &mocks.Loader{}
	parser := &mocks.SpecParser{}
	runner := &mocks.SpecRunner{}
	printer := &mocks.ResultsPrinter{}

	body := new(loaders.Body)
	spec := new(parsers.Spec)
	result := new(runners.Result)

	loader.On("Load", "test-file").Return(body, nil)
	parser.On("Parse", body).Return(spec, nil)
	runner.On("Run", spec).Return(result, nil)
	printer.On("Print", result).Return(nil)

	subject := main.Engine{
		Loader:  loader,
		Parser:  parser,
		Runner:  runner,
		Printer: printer,
	}

	err := subject.Start("test-file")
	if err != nil {
		t.Fatalf("Failed with: %v", err)
	}

	loader.AssertExpectations(t)
	parser.AssertExpectations(t)
	runner.AssertExpectations(t)
	printer.AssertExpectations(t)
}
