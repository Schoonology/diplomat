package parsers

import (
	"fmt"
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

// Parse parses all the lines received over the provided channel, parsing
// them into Tests it sends over the returned channel.
// It sends any errors encountered over the errors channel.
func (m *PlainTextParser) Parse(lines chan string, errors chan error) chan Test {
	state := newParserState()
	state.finalizer = fallbackTestName

	go func() {
		for line := range lines {
			err := state.addLine(line)
			if err != nil {
				errors <- err
				close(state.tests)
				return
			}
		}

		err := state.finalize()
		if err != nil {
			errors <- err
			return
		}
	}()

	return state.tests
}
