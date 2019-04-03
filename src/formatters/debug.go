package formatters

import (
	"fmt"

	"github.com/testdouble/diplomat/errors"
	"github.com/testdouble/diplomat/runners"
)

// Debug defines an unfiltered Printer.
type Debug struct{}

// Format returns all output, unfiltered.
func (d *Debug) Format(result runners.TestResult) (string, error) {
	if result.Err != nil {
		return errors.Format(result.Err), result.Err
	}

	return fmt.Sprintf("%v\n%v\n", result.Name, result.Diff), nil
}

// FormatAll formates all test results in a channel.
func (d *Debug) FormatAll(results chan runners.TestResult, errorChannel chan error) chan string {
	c := make(chan string)

	go func() {
		defer close(c)

		for result := range results {
			formattedResult, err := d.Format(result)
			if err != nil {
				errorChannel <- err
			}

			c <- formattedResult
		}
	}()

	return c
}
