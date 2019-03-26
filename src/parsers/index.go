package parsers

// A ParagraphParser is capable of parsing a stream of lines into chunks.
type ParagraphParser interface {
	Parse([]string) []Paragraph
	ParseAll(chan []string) chan Paragraph
}

// Paragraph contains a name and a body (array of lines) representing a test specification.
type Paragraph struct {
	Name       string
	Body       []string
	LineNumber int
}
