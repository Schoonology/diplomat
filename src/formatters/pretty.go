package formatters

import (
	"fmt"
	"strings"

	"github.com/testdouble/diplomat/colors"
	"github.com/testdouble/diplomat/errors"
	"github.com/testdouble/diplomat/runners"
)

// Pretty structures test results into a human-readable format.
type Pretty struct {
	Colorizer colors.Colorizer
}

// Format prints each test's name and a symbol indicating pass or fail.
// Failing tests will include an error or diff below the test name.
func (p *Pretty) Format(results chan runners.TestResult, errorChannel chan error) chan string {
	c := make(chan string)

	go func() {
		defer close(c)

		for result := range results {
			builder := strings.Builder{}

			if result.Err != nil {
				if result.Name != "" {
					builder.WriteString(p.Colorizer.Paint(fmt.Sprintf("✗ %s\n", result.Name), colors.Red))
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

			builder.WriteString(p.Colorizer.Paint(fmt.Sprintf("%s %s\n", symbol, result.Name), color))
			builder.WriteString(result.Diff)
			c <- builder.String()
		}
	}()

	return c
}
