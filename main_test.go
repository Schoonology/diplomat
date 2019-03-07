package main_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	main "github.com/testdouble/http-assertion-tool"
	"github.com/testdouble/http-assertion-tool/http"
	"github.com/testdouble/http-assertion-tool/loaders"
	"github.com/testdouble/http-assertion-tool/mocks"
	"github.com/testdouble/http-assertion-tool/parsers"
	"github.com/testdouble/http-assertion-tool/runners"
	"github.com/testdouble/http-assertion-tool/transforms"
)

func TestEngineStart(t *testing.T) {
	loader := &mocks.Loader{}
	parser := &mocks.SpecParser{}
	runner := &mocks.SpecRunner{}
	printer := &mocks.ResultsPrinter{}
	transforms := []transforms.Transform{
		func(spec *parsers.Spec) error {
			spec.Tests = append(spec.Tests, parsers.Test{Request: &http.Request{Method: "One"}})
			return nil
		},
		func(spec *parsers.Spec) error {
			spec.Tests = append(spec.Tests, parsers.Test{Request: &http.Request{Method: "Two"}})
			return nil
		},
	}

	body := new(loaders.Body)
	spec := new(parsers.Spec)
	result := new(runners.Result)

	loader.On("Load", "test-file").Return(body, nil)
	parser.On("Parse", body).Return(spec, nil)
	runner.On("Run", spec).Return(result, nil)
	printer.On("Print", result).Return(nil)

	subject := main.Engine{
		Loader:     loader,
		Parser:     parser,
		Transforms: transforms,
		Runner:     runner,
		Printer:    printer,
	}

	err := subject.Start("test-file")
	if err != nil {
		t.Fatalf("Failed with: %v", err)
	}

	loader.AssertExpectations(t)
	parser.AssertExpectations(t)
	runner.AssertExpectations(t)
	printer.AssertExpectations(t)

	assert.Equal(t, len(spec.Tests), 2)
}
