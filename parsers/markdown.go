package parsers

import (
	"strings"
)

// The Markdown parser parses all lines inside of code fences (```).
type Markdown struct {
	plainText PlainTextParser
}

const (
	inRichText = iota
	inCodeFence
)

// Parse parses all the lines received over the provided channel, parsing
// them into Tests it sends over the returned channel.
// It sends any errors encountered over the errors channel.
func (m *Markdown) Parse(lines chan string, errors chan error) chan Test {
	state := newParserState()

	go func() {
		mode := inRichText
		thisTestName := ""
		nextTestName := ""
		state.finalizer = func(test *Test) {
			if len(thisTestName) == 0 {
				fallbackTestName(test)
			} else {
				test.Name = thisTestName
			}

			thisTestName = nextTestName
			nextTestName = ""
		}

		for line := range lines {
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
					errors <- err
					close(state.tests)
					return
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

		err := state.finalize()
		if err != nil {
			errors <- err
		}
	}()

	return state.tests
}
