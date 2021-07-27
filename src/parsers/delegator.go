package parsers

import "github.com/schoonology/diplomat/loaders"

// A Delegator is responsible for communicating between a stream of files and the appropriate parser.
type Delegator struct{}

// ParseAll parses all the lines received over the provided channel, parsing
// them into Paragraphs it sends over the returned channel.
func (d *Delegator) ParseAll(files chan loaders.File) chan Paragraph {
	c := make(chan Paragraph)

	go func() {
		for file := range files {
			parser := GetParserFromFileName(file.Name)

			paragraphs := parser.Parse(file.Body)

			for _, paragraph := range paragraphs {
				paragraph.FileName = file.Name
				c <- paragraph
			}
		}

		close(c)
	}()

	return c
}
