package printers

import "github.com/testdouble/diplomat/runners"

// A ResultsPrinter defines a method to output all provided test results
// in a meaningful way.
type ResultsPrinter interface {
	Print(chan runners.TestResult, chan error)
}

// A PrinterOptions object contains additional configuration options for a printer
type PrinterOptions struct {
	Color bool
}
