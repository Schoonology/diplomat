package parsers

// A SpecParser is capable of parsing all lines in `body`.
type SpecParser interface {
	Parse(chan string, chan error) chan []string
}
