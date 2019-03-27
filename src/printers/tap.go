package printers

import (
	"fmt"

	"github.com/testdouble/diplomat/errors"
	"github.com/testdouble/diplomat/runners"
)

// Tap provides a TAP-conforming Printer.
type Tap struct{}

// Print sends TAP-conforming test results to STDOUT.
func (t *Tap) Print(results chan runners.TestResult, errorChannel chan error) {
	fmt.Println("TAP version 13")

	go func() {
		defer close(errorChannel)

		idx := 1
		for result := range results {
			if result.Err != nil {
				// TODO:(bam) - replace these defers with a string builder
				// the output is appearing in reverse order.
				defer errors.Display(result.Err)
				defer fmt.Printf("%s:\n", result.Name)
				fmt.Printf("not ok %d %s\n", idx, result.Name)
				idx++
				errorChannel <- result.Err
				continue
			}

			failed := len(result.Diff) > 0
			status := "ok"
			if failed {
				status = "not ok"
			}

			fmt.Printf("%s %d %s\n", status, idx, result.Name)
			idx++

			if failed {
				defer fmt.Printf("%s:\n%s\n", result.Name, result.Diff)
				// fmt.Print(result.Diff)
			}
		}

		fmt.Println()
	}()
}
