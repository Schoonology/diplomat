package loggers

import "fmt"

// StandardOutput writes strings to stdout.
type StandardOutput struct{}

// Print prints strings to stdout.
func (s *StandardOutput) Print(str string) {
	fmt.Print(str)
}
