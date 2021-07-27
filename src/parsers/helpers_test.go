package parsers_test

import (
	"testing"

	"github.com/schoonology/diplomat/parsers"
	"github.com/stretchr/testify/assert"
)

func streamBody(body []string) chan string {
	lines := make(chan string)

	go func() {
		for _, line := range body {
			lines <- line
		}

		close(lines)
	}()

	return lines
}

func assertTest(t *testing.T, expected parsers.Paragraph, paragraphs chan parsers.Paragraph) {
	select {
	case paragraph := <-paragraphs:
		assert.Equal(t, expected, paragraph)
	}
}
