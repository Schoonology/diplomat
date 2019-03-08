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
