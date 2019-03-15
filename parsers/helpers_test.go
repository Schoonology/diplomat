package parsers_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/testdouble/diplomat/parsers"
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

func assertTest(t *testing.T, expected parsers.Spec, specs chan parsers.Spec, errors chan error) {
	select {
	case err := <-errors:
		assert.FailNow(t, "Error should not exist.", err)
	case spec := <-specs:
		assert.Equal(t, expected, spec)
	}
}
