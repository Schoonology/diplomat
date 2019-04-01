package printers

import (
	"fmt"

	"github.com/testdouble/diplomat/colors"
	"github.com/testdouble/diplomat/errors"
	"github.com/testdouble/diplomat/runners"
)

// Pretty defines a formatted printer.
type Pretty struct {
	Options PrinterOptions
}

// Print prints all output, unfiltered.
func (t *Pretty) Print(results chan runners.TestResult, errorChannel chan error) {
	colorizer := colors.DefaultColorizer(t.Options.Color)

	go func() {
		defer close(errorChannel)

		for result := range results {
			if result.Err != nil {
				if result.Name != "" {
					colorizer.Print(fmt.Sprintf("✗ %s\n", result.Name), colors.Red)
				}
				errors.Display(result.Err)
				errorChannel <- result.Err
				continue
			}

			color := colors.Green
			symbol := "✓"
			if len(result.Diff) != 0 {
				color = colors.Red
				symbol = "✗"
				errorChannel <- errors.NewAssertionError(result.Diff)
			}

			colorizer.Print(fmt.Sprintf("%s %s\n", symbol, result.Name), color)
			fmt.Println(result.Diff)
		}
	}()
}
