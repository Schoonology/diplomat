package runners

import "github.com/testdouble/http-assertion-tool/parsers"

// A SpecRunner runs the entirety of a Spec in a given order, returning the
// complete results.
type SpecRunner interface {
	Run(*parsers.Spec) (*Result, error)
}

// TestResult is a container for the name and diff of a completed Test.
type TestResult struct {
	Name string
	Diff string
}

// Result is a set of TestResults.
type Result struct {
	Results []TestResult
}
