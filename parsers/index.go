package parsers

import (
	"github.com/testdouble/diplomat/http"
)

// A SpecParser is capable of parsing all lines in `body`.
type SpecParser interface {
	Parse(chan string, chan error) chan Test
}

// Test contains a name, request, and expected response.
type Test struct {
	Name     string
	Request  *http.Request
	Response *http.Response
}
