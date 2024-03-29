package builders

import (
	"github.com/schoonology/diplomat/http"
	"github.com/schoonology/diplomat/parsers"
)

// A SpecBuilder constructs a test from a Paragraph.
type SpecBuilder interface {
	Build(parsers.Paragraph) (Test, error)
	BuildAll(chan parsers.Paragraph) chan Test
}

// Test contains a name, request, and expected response.
type Test struct {
	Name     string
	Request  *http.Request
	Response *http.Response
	Err      error
}
