package formatters

import (
	"fmt"
	"strings"

	"github.com/schoonology/diplomat/errors"
	"github.com/schoonology/diplomat/runners"
)

// Tap structures test results to be TAP-compliant.
type Tap struct{}

// Format converts a test result into a TAP-compliant strings.
func (t *Tap) Format(result runners.TestResult, idx int) (string, error) {
	if result.Err != nil {
		output := fmt.Sprintf("not ok %d %s\n", idx, result.Name)
		err := fmt.Errorf("%s:\n%v", result.Name, errors.Format(result.Err))

		return output, err
	}

	failed := len(result.Diff) > 0
	status := "ok"
	if failed {
		status = "not ok"
	}

	var err error
	output := fmt.Sprintf("%s %d %s\n", status, idx, result.Name)

	if failed {
		err = fmt.Errorf("%s:\n%s", result.Name, result.Diff)
	}

	return output, err
}

// FormatAll formats all test results in a channel.
func (t *Tap) FormatAll(results chan runners.TestResult, errorChannel chan error) chan string {
	c := make(chan string)

	go func() {
		defer close(c)

		errorBuilder := strings.Builder{}
		errorBuilder.WriteString("\n")

		c <- "TAP version 13\n"

		idx := 1

		for result := range results {
			formattedResult, err := t.Format(result, idx)
			if err != nil {
				errorChannel <- err
				errorBuilder.WriteString(fmt.Sprintf("%v\n", err.Error()))
			}

			idx++

			c <- formattedResult
		}

		c <- errorBuilder.String()
	}()

	return c
}
