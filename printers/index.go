package printers

import "github.com/testdouble/http-assertion-tool/runners"

// A ResultsPrinter defines a method to output all provided test results
// in a meaningful way.
type ResultsPrinter interface {
	Print(*runners.Result) error
}
