package loaders

import (
	"io/ioutil"
	"strings"
)

// FileLoader loads all lines in a file.
type FileLoader struct{}

// Load returns all lines from the provided `filename`.
func (l *FileLoader) Load(filename string) (*Body, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return &Body{
		Lines: strings.Split(string(bytes), "\n"),
	}, nil
}

// Stream loads a file and sends it to a channel.
func (l *FileLoader) Stream(filename string) (chan *Body, chan error) {
	c := make(chan *Body)
	e := make(chan error)

	go func() {
		file, err := l.Load(filename)
		if err != nil {
			e <- err
		}

		c <- file
	}()

	return c, e
}
