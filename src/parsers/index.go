package parsers

import "github.com/testdouble/diplomat/loaders"

// A ParagraphParser is capable of parsing a stream of lines into chunks.
type ParagraphParser interface {
	Parse([]string) []Paragraph
}

// A ParseDelegator chooses the correct parser to handle an input
type ParseDelegator interface {
	ParseAll(chan loaders.File) chan Paragraph
}

// Paragraph contains a name and a body (array of lines) representing a test specification.
type Paragraph struct {
	Name       string
	Body       []string
	LineNumber int
}
