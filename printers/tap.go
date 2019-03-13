package printers

import (
	"fmt"

	"github.com/testdouble/diplomat/runners"
)

// Tap provides a TAP-conforming Printer.
type Tap struct{}

// Print sends TAP-conforming test results to STDOUT.
func (t *Tap) Print(results chan runners.TestResult, errors chan error) {
	fmt.Println("TAP version 13")

	go func() {
		idx := 0
		for result := range results {
			failed := len(result.Diff) > 0
			status := "ok"
			if failed {
				status = "not ok"
			}

			fmt.Printf("%s %d %s\n", status, idx, result.Name)
			idx++

			if failed {
				fmt.Print(result.Diff)
			}
		}

		close(errors)
	}()
}
