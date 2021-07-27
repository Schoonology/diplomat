package formatters

import "github.com/schoonology/diplomat/runners"

// A ResultsFormatter is responsible for formatting readable output for a test result.
type ResultsFormatter interface {
	FormatAll(chan runners.TestResult, chan error) chan string
}
