package builders

import (
	"github.com/testdouble/diplomat/http"
)

// A SpecBuilder constructs a test from an array of lines.
type SpecBuilder interface {
	Build([]string) (Test, error)
	BuildAll(chan []string, chan error) chan Test
}

// Test contains a name, request, and expected response.
type Test struct {
	Name     string
	Request  *http.Request
	Response *http.Response
}
