package printers

import (
	"fmt"

	"github.com/testdouble/diplomat/runners"
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

// Stream applies Print to results via a channel.
func (t *Debug) Stream(resultChannel chan *runners.Result, errorChannel chan error) (chan int, chan error) {
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
