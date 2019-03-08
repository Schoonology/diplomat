package parsers

import (
	"github.com/testdouble/http-assertion-tool/http"
	"github.com/testdouble/http-assertion-tool/loaders"
)

// A SpecParser is capable of parsing all lines in `body`.
type SpecParser interface {
	Parse(*loaders.Body) (*Spec, error)
}

// Spec contains a set of tests.
type Spec struct {
	Tests []Test
}

// Test contains a name, request, and expected response.
type Test struct {
	Name     string
	Request  *http.Request
	Response *http.Response
}
