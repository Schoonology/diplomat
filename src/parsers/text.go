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

// ParseAll parses all the lines received over the provided channel, parsing
// them into Paragraphs it sends over the returned channel.
func (m *PlainTextParser) ParseAll(files chan []string) chan Paragraph {
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
