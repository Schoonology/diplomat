package parsers

// PlainTextParser parses all provided text as-is.
type PlainTextParser struct{}

// Parse parses all the lines received over the provided channel, parsing
// them into Paragraphs it sends over the returned channel.
func (m *PlainTextParser) Parse(lines chan string) chan Paragraph {
	c := make(chan Paragraph)

	go func() {
		paragraph := Paragraph{}

		for line := range lines {
			paragraph.Body = append(paragraph.Body, line)
		}

		if len(paragraph.Body) > 0 {
			paragraph.LineNumber = 1
			c <- paragraph
		}

		close(c)
	}()

	return c
}
