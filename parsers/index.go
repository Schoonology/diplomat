package parsers

import "github.com/testdouble/http-assertion-tool/loaders"

type SpecParser interface {
	Parse(*loaders.File) (*Spec, error)
}

type Spec struct{}
