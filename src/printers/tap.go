package printers

import (
	"fmt"
	"strings"

	"github.com/testdouble/diplomat/errors"
	"github.com/testdouble/diplomat/runners"
)

// Tap provides a TAP-conforming Printer.
type Tap struct{}

// Print sends TAP-conforming test results to STDOUT.
func (t *Tap) Print(results chan runners.TestResult, errorChannel chan error) chan string {
	c := make(chan string)

	go func() {
		defer close(c)

		errorBuilder := strings.Builder{}
		errorBuilder.WriteString("\n")

		c <- "TAP version 13\n"

		idx := 1
		for result := range results {
			builder := strings.Builder{}

			if result.Err != nil {
				builder.WriteString(fmt.Sprintf("not ok %d %s\n", idx, result.Name))
				errorBuilder.WriteString(fmt.Sprintf("%s:\n", result.Name))
				errorBuilder.WriteString(errors.Format(result.Err))
				idx++
				c <- builder.String()
				errorChannel <- result.Err
				continue
			}

			failed := len(result.Diff) > 0
			status := "ok"
			if failed {
				status = "not ok"
			}

			builder.WriteString(fmt.Sprintf("%s %d %s\n", status, idx, result.Name))
			idx++

			if failed {
				errorBuilder.WriteString(fmt.Sprintf("%s:\n%s\n", result.Name, result.Diff))
			}

			c <- builder.String()
		}

		c <- errorBuilder.String()
	}()

	return c
}
