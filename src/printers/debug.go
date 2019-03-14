package printers

import (
	"fmt"

	"github.com/testdouble/diplomat/runners"
)

// Debug defines an unfiltered Printer.
type Debug struct{}

// Print prints all output, unfiltered.
func (t *Debug) Print(results chan runners.TestResult, errors chan error) {
	go func() {
		for result := range results {
			fmt.Printf("%v\n%v\n", result.Name, result.Diff)
		}

		close(errors)
	}()
}
