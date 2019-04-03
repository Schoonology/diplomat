package loggers

import "fmt"

// StandardOutput writes strings to stdout.
type StandardOutput struct{}

// Print prints strings to stdout.
func (s *StandardOutput) Print(str string) {
	fmt.Print(str)
}

// PrintAll prints all strings in a channel.
func (s *StandardOutput) PrintAll(output chan string, errors chan error) {
	go func() {
		defer close(errors)

		for str := range output {
			s.Print(str)
		}
	}()
}
