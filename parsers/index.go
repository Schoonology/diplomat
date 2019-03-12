package parsers

import (
	"github.com/testdouble/diplomat/http"
	"github.com/testdouble/diplomat/loaders"
)

// A SpecParser is capable of parsing all lines in `body`.
type SpecParser interface {
	Streamer
	Parse(*loaders.Body) (*Spec, error)
}

// A Streamer parses objects via a stream.
type Streamer interface {
	Stream(chan *loaders.Body) (chan *Spec, chan error)
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
