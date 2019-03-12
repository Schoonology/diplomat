package main_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	main "github.com/testdouble/diplomat"
	"github.com/testdouble/diplomat/http"
	"github.com/testdouble/diplomat/loaders"
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

	lastSpecChannel := make(chan *parsers.Spec)

	transforms := []transforms.TransformStream{
		func(specChannel chan *parsers.Spec, errorChannel chan error) (chan *parsers.Spec, chan error) {
			updatedSpecChannel := make(chan *parsers.Spec)
			go func() {
				spec := <-specChannel
				spec.Tests = append(spec.Tests, parsers.Test{Request: &http.Request{Method: "One"}})
				updatedSpecChannel <- spec
			}()
			return updatedSpecChannel, errorChannel
		},
		func(specChannel chan *parsers.Spec, errorChannel chan error) (chan *parsers.Spec, chan error) {
			go func() {
				spec := <-specChannel
				spec.Tests = append(spec.Tests, parsers.Test{Request: &http.Request{Method: "Two"}})
				lastSpecChannel <- spec
			}()
			return lastSpecChannel, errorChannel
		},
	}

	body := new(loaders.Body)
	spec := new(parsers.Spec)
	result := new(runners.Result)

	bodyChannel := make(chan *loaders.Body)
	errorChannel := make(chan error)

	go func() {
		bodyChannel <- body
	}()

	specChannel := make(chan *parsers.Spec)
	specErrorChannel := make(chan error)

	go func() {
		specChannel <- spec
	}()

	resultChannel := make(chan *runners.Result)
	resultErrorChannel := make(chan error)

	go func() {
		resultChannel <- result
	}()

	loader.On("Stream", "test-file").Return(bodyChannel, errorChannel)
	parser.On("Stream", bodyChannel).Return(specChannel, specErrorChannel)
	runner.On("Stream", lastSpecChannel, specErrorChannel).Return(resultChannel, resultErrorChannel)
	// runner.On("Run", spec).Return(result, nil)
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
