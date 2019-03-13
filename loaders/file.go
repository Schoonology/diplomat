package loaders

import (
	"io/ioutil"
	"strings"
)

// FileLoader loads all lines in a file.
type FileLoader struct{}

// Load returns all lines from the provided `filename`.
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
