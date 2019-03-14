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
func (m *Markdown) Parse(lines chan string, errors chan error) chan []string {
	c := make(chan []string)

	go func() {
		spec := []string{}
		mode := inRichText
		// thisTestName := ""
		// nextTestName := ""
		// state.finalizer = func(test *Test) {
		// 	if len(thisTestName) == 0 {
		// 		fallbackTestName(test)
		// 	} else {
		// 		test.Name = thisTestName
		// 	}

		// 	thisTestName = nextTestName
		// 	nextTestName = ""
		// }

		for line := range lines {
			trimmedLine := strings.TrimSpace(line)

			switch {
			case strings.HasPrefix(trimmedLine, "```"):
				if mode == inRichText {
					mode = inCodeFence
				} else if mode == inCodeFence {
					c <- spec
					spec = []string{}
					mode = inRichText
				} else {
					mode = inRichText
				}
			case mode == inCodeFence:
				spec = append(spec, line)
				// case strings.HasPrefix(trimmedLine, "#"):
				// 	name := strings.TrimSpace(strings.SplitN(trimmedLine, " ", 2)[1])

				// 	if state.mode == modeAwaitingRequest {
				// 		thisTestName = name
				// 	} else {
				// 		nextTestName = name
				// 	}
			}
		}

		close(c)

		// err := state.finalize()
		// if err != nil {
		// 	errors <- err
		// }
	}()

	return c
}
