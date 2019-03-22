package builders

import (
	"github.com/testdouble/diplomat/http"
	"github.com/testdouble/diplomat/parsers"
)

// A SpecBuilder constructs a test from an array of lines.
type SpecBuilder interface {
	Build(parsers.Spec) (Test, error)
	BuildAll(chan parsers.Spec, chan error) chan Test
}

// Test contains a name, request, and expected response.
type Test struct {
	Name     string
	Request  *http.Request
	Response *http.Response
	Err      error
}
