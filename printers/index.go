package printers

import "github.com/testdouble/diplomat/runners"

// A ResultsPrinter defines a method to output all provided test results
// in a meaningful way.
type ResultsPrinter interface {
	Streamer
	Print(*runners.Result) error
}

// A Streamer executes ResultsPrinter on test results through a stream.
type Streamer interface {
	Stream(chan *runners.Result, chan error) (chan int, chan error)
}
