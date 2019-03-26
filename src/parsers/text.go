package parsers

// PlainTextParser parses all provided text as-is.
type PlainTextParser struct{}

// Parse returns a slice of Paragraphs from a single file
func (m *PlainTextParser) Parse(file []string) []Paragraph {
	paragraphs := make([]Paragraph, 0)
	paragraph := Paragraph{}

	for _, line := range file {
		paragraph.Body = append(paragraph.Body, line)
	}

	if len(paragraph.Body) > 0 {
		paragraph.LineNumber = 1
		paragraphs = append(paragraphs, paragraph)
	}

	return paragraphs
}
