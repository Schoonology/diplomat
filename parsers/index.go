package parsers

import (
	"github.com/testdouble/http-assertion-tool/http"
	"github.com/testdouble/http-assertion-tool/loaders"
)

type SpecParser interface {
	Parse(*loaders.Body) (*Spec, error)
}

type Spec struct {
	Tests []Test
}

type Test struct {
	Request  *http.Request
	Response *http.Response
}
