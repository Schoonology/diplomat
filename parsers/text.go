package parsers

import (
	"strings"

	"github.com/testdouble/http-assertion-tool/loaders"
)

type PlainTextParser struct{}

func (m *PlainTextParser) Parse(body *loaders.Body) (*Spec, error) {
	request := new(strings.Builder)
	response := new(strings.Builder)
	for _, line := range body.Lines {
		if len(line) == 0 {
			continue
		} else if line[0] == '>' {
			if len(line) > 1 {
				request.WriteString(line[2:])
			} else {
				request.WriteString(line[1:])
			}
			request.WriteString("\r\n")
		} else if line[0] == '<' {
			if len(line) > 1 {
				response.WriteString(line[2:])
			} else {
				response.WriteString(line[1:])
			}
			response.WriteString("\r\n")
		}
	}

	spec := Spec{
		Tests: []Test{Test{
			Request:  request.String(),
			Response: response.String(),
		}},
	}

	return &spec, nil
}
