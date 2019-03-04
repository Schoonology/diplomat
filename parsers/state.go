package parsers

import (
	"strconv"
	"strings"

	"github.com/testdouble/http-assertion-tool/http"
)

const (
	ModeAwaitingRequest = iota
	ModeInRequestHeaders
	ModeAwaitingRequestBody
	ModeInResponseHeaders
	ModeAwaitingResponseBody
)

type testFinalizer func(test *Test)

type parserState struct {
	currentTest Test
	mode        int
	spec        *Spec
	finalizer   testFinalizer
}

func newParserState() parserState {
	spec := new(Spec)
	spec.Tests = make([]Test, 0)

	return parserState{
		mode: ModeAwaitingRequest,
		spec: spec,
	}
}

func (s *parserState) pushCurrentTest() {
	if s.finalizer != nil {
		s.finalizer(&s.currentTest)
	}

	s.spec.Tests = append(s.spec.Tests, s.currentTest)
	s.currentTest = Test{}
}

func (s *parserState) addLine(line string) error {
	switch {
	case len(line) == 0:
		return nil
	case line[0] == '>':
		line = strings.TrimSpace(line[1:])

		if s.mode == ModeAwaitingRequest || s.mode == ModeInResponseHeaders || s.mode == ModeAwaitingResponseBody {
			if s.mode == ModeInResponseHeaders || s.mode == ModeAwaitingResponseBody {
				s.pushCurrentTest()
			}

			s.currentTest.Request = RequestFromLine(line)
			s.mode = ModeInRequestHeaders
		} else if s.mode == ModeInRequestHeaders && len(line) == 0 {
			s.mode = ModeAwaitingRequestBody
		} else if s.mode == ModeInRequestHeaders {
			key, value := HeaderFromLine(line)
			s.currentTest.Request.Headers[key] = value
		}
	case line[0] == '<':
		line = strings.TrimSpace(line[1:])

		if s.mode == ModeInRequestHeaders || s.mode == ModeAwaitingRequestBody {
			s.currentTest.Response = ResponseFromLine(line)
			s.mode = ModeInResponseHeaders
		} else if s.mode == ModeInResponseHeaders && len(line) == 0 {
			s.mode = ModeAwaitingResponseBody
		} else if s.mode == ModeInResponseHeaders {
			key, value := HeaderFromLine(line)
			s.currentTest.Response.Headers[key] = value
		}
	}

	return nil
}

func (s *parserState) finalize() (*Spec, error) {
	if s.mode != ModeAwaitingRequest {
		s.pushCurrentTest()
	}

	return s.spec, nil
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
