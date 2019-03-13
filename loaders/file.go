package loaders

import (
	"io/ioutil"
	"strings"
)

// FileLoader loads all lines in a file.
type FileLoader struct{}

// Load sends all lines from the provided `filename` along the returned
// channel, closing the channel once all lines have been sent.
// Load sends any errors to the provided error channel.
func (l *FileLoader) Load(filename string, errors chan error) chan string {
	lines := make(chan string)

	go func() {
		bytes, err := ioutil.ReadFile(filename)
		if err != nil {
			errors <- err
			return
		}

		for _, line := range strings.Split(string(bytes), "\n") {
			lines <- line
		}

		close(lines)
	}()

	return lines
}
