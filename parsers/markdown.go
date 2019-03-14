package parsers

import (
	"strings"
)

// The Markdown parser parses all lines inside of code fences (```).
type Markdown struct{}

const (
	inRichText = iota
	inCodeFence
)

// Parse parses all the lines received over the provided channel, parsing
// them into Tests it sends over the returned channel.
// It sends any errors encountered over the errors channel.
func (m *Markdown) Parse(lines chan string, errors chan error) chan Spec {
	c := make(chan Spec)

	go func() {
		spec := Spec{}
		mode := inRichText

		for line := range lines {
			trimmedLine := strings.TrimSpace(line)

			switch {
			case strings.HasPrefix(trimmedLine, "```"):
				if mode == inRichText {
					mode = inCodeFence
				} else if mode == inCodeFence {
					c <- spec
					spec = Spec{}
					mode = inRichText
				} else {
					mode = inRichText
				}
			case mode == inCodeFence:
				spec.Body = append(spec.Body, line)
			case strings.HasPrefix(trimmedLine, "#"):
				spec.Name = strings.TrimSpace(strings.SplitN(trimmedLine, " ", 2)[1])
			}
		}

		if len(spec.Body) != 0 {
			c <- spec
		}

		close(c)
	}()

	return c
}
