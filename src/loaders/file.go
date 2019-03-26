package loaders

import (
	"io/ioutil"
	"path/filepath"
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
func (l *FileLoader) LoadAll(filenames chan string, errors chan error) chan File {
	files := make(chan File)

	go func() {
		defer close(files)

		for filename := range filenames {
			bytes, err := l.Load(filename)
			if err != nil {
				errors <- err
				close(errors)
				return
			}

			files <- File{
				Name: filepath.Base(filename),
				Body: strings.Split(string(bytes), "\n"),
			}
		}
	}()

	return files
}
