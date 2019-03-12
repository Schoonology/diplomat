package runners

import "github.com/testdouble/diplomat/parsers"

// A SpecRunner runs the entirety of a Spec in a given order, returning the
// complete results.
type SpecRunner interface {
	Streamer
	Run(*parsers.Spec) (*Result, error)
}

// A Streamer executes the SpecRunner on Specs in a channel.
type Streamer interface {
	Stream(chan *parsers.Spec, chan error) (chan *Result, chan error)
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
