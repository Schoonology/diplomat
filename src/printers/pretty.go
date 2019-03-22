package printers

import (
	"fmt"

	"github.com/testdouble/diplomat/errors"
	"github.com/testdouble/diplomat/runners"
)

// Pretty defines a formatted printer.
type Pretty struct{}

// Print prints all output, unfiltered.
func (t *Pretty) Print(results chan runners.TestResult, errorChannel chan error) {
	go func() {
		defer close(errorChannel)

		for result := range results {
			if result.Err != nil {
				if result.Name != "" {
					fmt.Printf("✗ %s\n", result.Name)
				}
				errors.Display(result.Err)
				errorChannel <- result.Err
				continue
			}

			fmt.Printf("✓ %s\n%s", result.Name, result.Diff)
		}
	}()
}
