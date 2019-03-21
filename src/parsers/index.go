package parsers

// A SpecParser is capable of parsing all lines in `body`.
type SpecParser interface {
	Parse(chan string, chan error) chan Spec
}

// Spec contains a name and a body (array of lines) representing a test specification.
type Spec struct {
	Name       string
	Body       []string
	LineNumber int
}
