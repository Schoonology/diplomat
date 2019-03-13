package runners

import (
	"github.com/testdouble/diplomat/parsers"
)

// A SpecRunner runs the entirety of a Spec in a given order, emitting test
// results to the returned channel.
type SpecRunner interface {
	Run(chan parsers.Test, chan error) chan TestResult
}

// TestResult is a container for the name and diff of a completed Test.
type TestResult struct {
	Name string
	Diff string
}
