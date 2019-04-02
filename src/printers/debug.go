package printers

import (
	"fmt"

	"github.com/testdouble/diplomat/errors"
	"github.com/testdouble/diplomat/runners"
)

// Debug defines an unfiltered Printer.
type Debug struct{}

// Print prints all output, unfiltered.
func (t *Debug) Print(results chan runners.TestResult, errorChannel chan error) chan string {
	c := make(chan string)

	go func() {
		defer close(c)
		defer close(errorChannel)

		for result := range results {
			if result.Err != nil {
				errors.Display(result.Err)
				errorChannel <- result.Err
				continue
			}

			fmt.Printf("%v\n%v\n", result.Name, result.Diff)
		}
	}()

	return c
}
