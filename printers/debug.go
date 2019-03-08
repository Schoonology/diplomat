package printers

import (
	"fmt"

	"github.com/testdouble/http-assertion-tool/runners"
)

// Debug defines an unfiltered Printer.
type Debug struct{}

// Print prints all output, unfiltered.
func (t *Debug) Print(result *runners.Result) error {
	for _, result := range result.Results {
		fmt.Printf("%v\n%v\n", result.Name, result.Diff)
	}

	return nil
}
