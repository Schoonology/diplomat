package parsers

import (
	"fmt"

	"github.com/testdouble/diplomat/loaders"
)

// PlainTextParser parses all provided text as-is.
type PlainTextParser struct{}

func fallbackTestName(test *Test) {
	test.Name = fmt.Sprintf(
		"%s %s -> %d",
		test.Request.Method,
		test.Request.Path,
		test.Response.StatusCode)
}

// Parse parses all lines in `body`.
func (m *PlainTextParser) Parse(body *loaders.Body) (*Spec, error) {
	state := newParserState()
	state.finalizer = fallbackTestName

	for _, line := range body.Lines {
		err := state.addLine(line)
		if err != nil {
			return nil, err
		}
	}

	return state.finalize()
}
