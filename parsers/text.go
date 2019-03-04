package parsers

import (
	"fmt"

	"github.com/testdouble/http-assertion-tool/loaders"
)

type PlainTextParser struct{}

func fallbackTestName(test *Test) {
	test.Name = fmt.Sprintf(
		"%s %s -> %d",
		test.Request.Method,
		test.Request.Path,
		test.Response.StatusCode)
}

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
