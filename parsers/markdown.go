package parsers

import (
	"strings"

	"github.com/testdouble/http-assertion-tool/loaders"
)

type Markdown struct {
	plainText PlainTextParser
}

const (
	InRichText = iota
	InCodeFence
)

func (m *Markdown) Parse(body *loaders.Body) (*Spec, error) {
	mode := InRichText
	thisTestName := ""
	nextTestName := ""
	state := newParserState()
	state.finalizer = func(test *Test) {
		if len(thisTestName) == 0 {
			fallbackTestName(test)
		} else {
			test.Name = thisTestName
		}

		thisTestName = nextTestName
		nextTestName = ""
	}

	for _, line := range body.Lines {
		trimmedLine := strings.TrimSpace(line)

		switch {
		case strings.HasPrefix(trimmedLine, "```"):
			if mode == InRichText {
				mode = InCodeFence
			} else {
				mode = InRichText
			}
		case mode == InCodeFence:
			state.addLine(line)
		case strings.HasPrefix(trimmedLine, "#"):
			name := strings.TrimSpace(strings.SplitN(trimmedLine, " ", 2)[1])

			if state.mode == ModeAwaitingRequest {
				thisTestName = name
			} else {
				nextTestName = name
			}
		}
	}

	return state.finalize()
}
