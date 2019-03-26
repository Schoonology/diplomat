package loaders

import (
	"io/ioutil"
	"strings"
)

// FileLoader loads all lines in a file.
type FileLoader struct{}

// Load retrieves all bytes from a specified file.
func (l *FileLoader) Load(filename string) ([]byte, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

// LoadAll loads all files from a stream, and sends the bytes through the output channel.
// If an error is encountered, it stops the process and sends an error into the error stream.
func (l *FileLoader) LoadAll(filenames chan string, errors chan error) chan []string {
	files := make(chan []string)

	go func() {
		defer close(files)

		for filename := range filenames {
			bytes, err := l.Load(filename)
			if err != nil {
				errors <- err
				close(errors)
				return
			}

			files <- strings.Split(string(bytes), "\n")
		}
	}()

	return files
}
