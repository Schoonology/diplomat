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
// them into Paragraphs it sends over the returned channel.
func (m *Markdown) Parse(lines chan string) chan Paragraph {
	c := make(chan Paragraph)

	go func() {
		paragraph := Paragraph{
			LineNumber: 1,
		}
		mode := inRichText
		lineNumber := 1

		for line := range lines {
			trimmedLine := strings.TrimSpace(line)

			switch {
			case strings.HasPrefix(trimmedLine, "```"):
				if mode == inRichText {
					mode = inCodeFence
					paragraph.LineNumber = lineNumber + 1
				} else if mode == inCodeFence {
					c <- paragraph
					paragraph = Paragraph{}
					mode = inRichText
				} else {
					mode = inRichText
				}
			case mode == inCodeFence:
				paragraph.Body = append(paragraph.Body, line)
			case strings.HasPrefix(trimmedLine, "#"):
				paragraph.Name = strings.TrimSpace(strings.SplitN(trimmedLine, " ", 2)[1])
			}

			lineNumber++
		}

		if len(paragraph.Body) != 0 {
			c <- paragraph
		}

		close(c)
	}()

	return c
}
