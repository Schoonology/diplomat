package builders

import "github.com/testdouble/diplomat/parsers"

// A SpecBuilder constructs a test from an array of lines.
type SpecBuilder interface {
	Build([]string) parsers.Test
}
