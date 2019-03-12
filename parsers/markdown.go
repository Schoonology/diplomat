package parsers

import (
	"strings"

	"github.com/testdouble/diplomat/loaders"
)

// The Markdown parser parses all lines inside of code fences (```).
type Markdown struct {
	plainText PlainTextParser
}

const (
	inRichText = iota
	inCodeFence
)

// Parse parses all lines in `body`.
func (m *Markdown) Parse(body *loaders.Body) (*Spec, error) {
	mode := inRichText
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
			if mode == inRichText {
				mode = inCodeFence
			} else {
				mode = inRichText
			}
		case mode == inCodeFence:
			err := state.addLine(line)
			if err != nil {
				return nil, err
			}
		case strings.HasPrefix(trimmedLine, "#"):
			name := strings.TrimSpace(strings.SplitN(trimmedLine, " ", 2)[1])

			if state.mode == modeAwaitingRequest {
				thisTestName = name
			} else {
				nextTestName = name
			}
		}
	}

	return state.finalize()
}

// Stream parses a streamed body
func (m *Markdown) Stream(bodyChannel chan *loaders.Body, errorChannel chan error) (chan *Spec, chan error) {
	c := make(chan *Spec)
	e := make(chan error)

	go func() {
		var body *loaders.Body
		select {
		case body = <-bodyChannel:
		case err := <-errorChannel:
			e <- err
			return
		}

		spec, err := m.Parse(body)
		if err != nil {
			e <- err
		}

		c <- spec
	}()

	return c, e
}
