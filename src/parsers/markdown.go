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

// Parse returns a slice of Paragraphs from a single file
func (m *Markdown) Parse(file []string) []Paragraph {
	paragraphs := make([]Paragraph, 0)

	paragraph := Paragraph{
		LineNumber: 1,
	}
	mode := inRichText
	lineNumber := 1

	for _, line := range file {
		trimmedLine := strings.TrimSpace(line)

		switch {
		case strings.HasPrefix(trimmedLine, "```"):
			if mode == inRichText {
				mode = inCodeFence
				paragraph.LineNumber = lineNumber + 1
			} else if mode == inCodeFence {
				paragraphs = append(paragraphs, paragraph)
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
		paragraphs = append(paragraphs, paragraph)
	}

	return paragraphs
}

// ParseAll parses all the lines received over the provided channel, parsing
// them into Paragraphs it sends over the returned channel.
func (m *Markdown) ParseAll(files chan []string) chan Paragraph {
	c := make(chan Paragraph)

	go func() {
		for file := range files {
			paragraphs := m.Parse(file)

			for _, paragraph := range paragraphs {
				c <- paragraph
			}
		}

		close(c)
	}()

	return c
}
