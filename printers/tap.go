package printers

import (
	"fmt"

	"github.com/testdouble/diplomat/runners"
)

// Tap provides a TAP-conforming Printer.
type Tap struct{}

// Print sends TAP-conforming test results to STDOUT.
func (t *Tap) Print(result *runners.Result) error {
	fmt.Println("TAP version 13")
	fmt.Printf("1..%d\n", len(result.Results))

	for idx, result := range result.Results {
		failed := len(result.Diff) > 0
		status := "ok"
		if failed {
			status = "not ok"
		}

		fmt.Printf("%s %d %s\n", status, idx, result.Name)

		if failed {
			fmt.Print(result.Diff)
		}
	}

	return nil
}

// Stream applies Print to results via a channel.
func (t *Tap) Stream(resultChannel chan *runners.Result, errorChannel chan error) (chan int, chan error) {
	c := make(chan int)
	e := make(chan error)

	go func() {
		var result *runners.Result
		select {
		case result = <-resultChannel:
		case err := <-errorChannel:
			e <- err
			return
		}

		err := t.Print(result)
		if err != nil {
			e <- err
		}

		c <- 0
	}()

	return c, e
}
