package parsers

// PlainTextParser parses all provided text as-is.
type PlainTextParser struct{}

// Parse parses all the lines received over the provided channel, parsing
// them into Specs it sends over the returned channel.
// It sends any errors encountered over the errors channel.
func (m *PlainTextParser) Parse(lines chan string, errors chan error) chan Spec {
	c := make(chan Spec)

	go func() {
		spec := Spec{}

		for line := range lines {
			spec.Body = append(spec.Body, line)
		}

		if len(spec.Body) > 0 {
			c <- spec
		}

		close(c)
	}()

	return c
}
