package printers

import (
	"fmt"

	"github.com/testdouble/diplomat/errors"
	"github.com/testdouble/diplomat/runners"

	"github.com/logrusorgru/aurora"
)

// Pretty defines a formatted printer.
type Pretty struct {
	Options PrinterOptions
}

// Print prints all output, unfiltered.
func (t *Pretty) Print(results chan runners.TestResult, errorChannel chan error) {
	var au aurora.Aurora
	au = aurora.NewAurora(t.Options.Color)

	go func() {
		defer close(errorChannel)

		for result := range results {
			if result.Err != nil {
				if result.Name != "" {
					fmt.Println(au.Red(fmt.Sprintf("✗ %s", result.Name)))
				}
				errors.Display(result.Err)
				errorChannel <- result.Err
				continue
			}

			colorFunction := au.Green
			symbol := "✓"
			if len(result.Diff) != 0 {
				colorFunction = au.Red
				symbol = "✗"
				errorChannel <- errors.NewAssertionError(result.Diff)
			}

			fmt.Println(colorFunction(fmt.Sprintf("%s %s", symbol, result.Name)))
			fmt.Println(result.Diff)
		}
	}()
}
