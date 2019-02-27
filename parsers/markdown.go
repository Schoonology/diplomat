package parsers

import "github.com/testdouble/http-assertion-tool/loaders"

type MarkdownParser struct{}

func (m *MarkdownParser) Parse(*loaders.Body) (*Spec, error) {
	return nil, nil
}
