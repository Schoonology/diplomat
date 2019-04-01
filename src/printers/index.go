package printers

import "github.com/testdouble/diplomat/runners"

// A ResultsPrinter defines a method to output all provided test results
// in a meaningful way.
type ResultsPrinter interface {
	Print(chan runners.TestResult, chan error)
}
