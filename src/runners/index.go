package runners

import "github.com/testdouble/diplomat/builders"

// A SpecRunner runs the entirety of a Spec in a given order, emitting test
// results to the returned channel.
type SpecRunner interface {
	Run(builders.Test) (TestResult, error)
	RunAll(chan builders.Test) chan TestResult
}

// TestResult is a container for the name and diff of a completed Test.
type TestResult struct {
	Name string
	Diff string
	Err  error
}
