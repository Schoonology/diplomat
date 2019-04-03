package printers

import (
	"fmt"
	"strings"

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

		for result := range results {
			builder := strings.Builder{}

			if result.Err != nil {
				if result.Name != "" {
					builder.WriteString(t.Colorizer.Paint(fmt.Sprintf("✗ %s\n", result.Name), colors.Red))
				}
				builder.WriteString(errors.Format(result.Err))
				errorChannel <- result.Err
				c <- builder.String()
				continue
			}

			color := colors.Green
			symbol := "✓"
			if len(result.Diff) != 0 {
				color = colors.Red
				symbol = "✗"
				errorChannel <- errors.NewAssertionError(result.Diff)
			}

			builder.WriteString(t.Colorizer.Paint(fmt.Sprintf("%s %s\n", symbol, result.Name), color))
			builder.WriteString(result.Diff)
			c <- builder.String()
		}
	}()

	return c
}
