package loaders

import (
	"io/ioutil"
	"strings"
)

type FileLoader struct{}

func (l *FileLoader) Load(filename string) (*Body, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return &Body{
		Lines: strings.Split(string(bytes), "\n"),
	}, nil
}
