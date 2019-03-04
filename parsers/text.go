package parsers

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/testdouble/http-assertion-tool/http"
	"github.com/testdouble/http-assertion-tool/loaders"
)

const (
	AwaitingRequest = iota
	InRequestHeaders
	AwaitingRequestBody
	InResponseHeaders
	AwaitingResponseBody
)

type PlainTextParser struct{}

type parserState struct {
	currentTest Test
	mode        int
	spec        *Spec
}

func newParserState() parserState {
	spec := new(Spec)
	spec.Tests = make([]Test, 0)

	return parserState{
		mode: AwaitingRequest,
		spec: spec,
	}
}

func (s *parserState) pushCurrentTest() {
	s.currentTest.Name = fmt.Sprintf(
		"%s %s -> %d", s.currentTest.Request.Method, s.currentTest.Request.Path, s.currentTest.Response.StatusCode)

	s.spec.Tests = append(s.spec.Tests, s.currentTest)
	s.currentTest = Test{}
}

func (s *parserState) addLine(line string) error {
	switch {
	case len(line) == 0:
		return nil
	case line[0] == '>':
		line = strings.TrimSpace(line[1:])

		if s.mode == AwaitingRequest || s.mode == InResponseHeaders || s.mode == AwaitingResponseBody {
			if s.mode == InResponseHeaders || s.mode == AwaitingResponseBody {
				s.pushCurrentTest()
			}

			s.currentTest.Request = RequestFromLine(line)
			s.mode = InRequestHeaders
		} else if s.mode == InRequestHeaders && len(line) == 0 {
			s.mode = AwaitingRequestBody
		} else if s.mode == InRequestHeaders {
			key, value := HeaderFromLine(line)
			s.currentTest.Request.Headers[key] = value
		}
	case line[0] == '<':
		line = strings.TrimSpace(line[1:])

		if s.mode == InRequestHeaders || s.mode == AwaitingRequestBody {
			s.currentTest.Response = ResponseFromLine(line)
			s.mode = InResponseHeaders
		} else if s.mode == InResponseHeaders && len(line) == 0 {
			s.mode = AwaitingResponseBody
		} else if s.mode == InResponseHeaders {
			key, value := HeaderFromLine(line)
			s.currentTest.Response.Headers[key] = value
		}
	}

	return nil
}

func (s *parserState) finalize() (*Spec, error) {
	if s.mode != AwaitingRequest {
		s.pushCurrentTest()
	}

	return s.spec, nil
}

func (m *PlainTextParser) Parse(body *loaders.Body) (*Spec, error) {
	state := newParserState()

	for _, line := range body.Lines {
		err := state.addLine(line)
		if err != nil {
			return nil, err
		}
	}

	return state.finalize()
}

// TODO(schoon) - Handle badly-formatted lines.
func RequestFromLine(line string) *http.Request {
	pieces := strings.Split(line, " ")

	return http.NewRequest(pieces[0], pieces[1])
}

// TODO(schoon) - Handle badly-formatted lines.
func HeaderFromLine(line string) (string, string) {
	pieces := strings.Split(line, ":")

	return strings.TrimSpace(pieces[0]), strings.TrimSpace(pieces[1])
}

// TODO(schoon) - Handle badly-formatted lines.
func ResponseFromLine(line string) *http.Response {
	pieces := strings.Split(line, " ")
	code, _ := strconv.Atoi(pieces[1])

	return http.NewResponse(code, strings.Join(pieces[2:], " "))
}
