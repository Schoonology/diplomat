package printers

import (
	"fmt"

	"github.com/testdouble/diplomat/colors"
	"github.com/testdouble/diplomat/errors"
	"github.com/testdouble/diplomat/runners"
)

// Pretty defines a formatted printer.
type Pretty struct {
	Colorizer colors.Colorizer
}

// Print prints all output, unfiltered.
func (t *Pretty) Print(results chan runners.TestResult, errorChannel chan error) chan string {
	c := make(chan string)

	go func() {
		defer close(c)
		defer close(errorChannel)

		for result := range results {
			if result.Err != nil {
				if result.Name != "" {
					fmt.Print(t.Colorizer.Paint(fmt.Sprintf("✗ %s\n", result.Name), colors.Red))
				}
				fmt.Print(errors.Format(result.Err))
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

			fmt.Print(t.Colorizer.Paint(fmt.Sprintf("%s %s\n", symbol, result.Name), color))
			fmt.Println(result.Diff)
		}
	}()

	return c
}
