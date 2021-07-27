package formatters

import (
	"fmt"
	"strings"

	"github.com/schoonology/diplomat/colors"
	"github.com/schoonology/diplomat/errors"
	"github.com/schoonology/diplomat/runners"
)

// Pretty structures test results into a human-readable format.
type Pretty struct {
	Colorizer colors.Colorizer
}

// Format prints a test's name and a symbol indicating pass or fail.
// Failing tests will include an error or diff below the test name.
func (p *Pretty) Format(result runners.TestResult) (string, error) {
	builder := strings.Builder{}

	if result.Err != nil {
		if result.Name != "" {
			builder.WriteString(p.Colorizer.Paint(fmt.Sprintf("✗ %s\n", result.Name), colors.Red))
		}
		builder.WriteString(errors.Format(result.Err))
		return builder.String(), result.Err
	}

	var err error
	color := colors.Green
	symbol := "✓"
	if len(result.Diff) != 0 {
		err = errors.NewAssertionError(result.Diff)
		color = colors.Red
		symbol = "✗"
	}

	builder.WriteString(p.Colorizer.Paint(fmt.Sprintf("%s %s\n", symbol, result.Name), color))
	builder.WriteString(result.Diff)

	return builder.String(), err
}

// FormatAll formats all test results in a channel.
func (p *Pretty) FormatAll(results chan runners.TestResult, errorChannel chan error) chan string {
	c := make(chan string)

	go func() {
		defer close(c)

		for result := range results {
			formattedResult, err := p.Format(result)
			if err != nil {
				errorChannel <- err
			}

			c <- formattedResult
		}
	}()

	return c
}
