package parsers

import (
	"errors"
	"strconv"
	"strings"

	"github.com/testdouble/http-assertion-tool/http"
)

const (
	ModeAwaitingRequest = iota
	ModeInRequestHeaders
	ModeAwaitingRequestBody
	ModeInResponseHeaders
	ModeInResponseBody
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
	if s.mode == ModeInResponseBody {
		s.currentTest.Response.Body = append(s.currentTest.Response.Body, []byte(line+"\n")...)
		return nil
	}

	switch {
	case len(line) > 0 && line[0] == '>':
		line = strings.TrimSpace(line[1:])

		if s.mode == ModeAwaitingRequest || s.mode == ModeInResponseHeaders || s.mode == ModeInResponseBody {
			if s.mode == ModeInResponseHeaders || s.mode == ModeInResponseBody {
				s.pushCurrentTest()
			}

			s.currentTest.Request = RequestFromLine(line)
			s.mode = ModeInRequestHeaders
		} else if s.mode == ModeInRequestHeaders && len(line) == 0 {
			s.mode = ModeAwaitingRequestBody
		} else if s.mode == ModeInRequestHeaders {
			key, value, err := HeaderFromLine(line)
			if err != nil {
				return err
			}

			s.currentTest.Request.Headers[key] = value
		}
	case len(line) > 0 && line[0] == '<':
		line := strings.TrimSpace(line[1:])

		if s.mode == ModeInRequestHeaders || s.mode == ModeAwaitingRequestBody {
			s.currentTest.Response = ResponseFromLine(line)
			s.mode = ModeInResponseHeaders
		} else if s.mode == ModeInResponseHeaders && len(line) == 0 {
			s.mode = ModeInResponseBody
		} else if s.mode == ModeInResponseHeaders {
			key, value, err := HeaderFromLine(line)
			if err != nil {
				return err
			}

			s.currentTest.Response.Headers[key] = value
		}
	default:
		if s.mode == ModeAwaitingRequest {
			return nil
		} else if s.mode == ModeInResponseHeaders || s.mode == ModeInResponseBody {
			s.mode = ModeInResponseBody
			s.currentTest.Response.Body = append(s.currentTest.Response.Body, []byte(line+"\n")...)
		} else {
			return errors.New("invalid formatting")
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
func HeaderFromLine(line string) (string, string, error) {
	pieces := strings.Split(line, ":")

	if len(pieces) < 2 {
		return "", "", errors.New("badly formatted header")
	}

	return strings.TrimSpace(pieces[0]), strings.TrimSpace(pieces[1]), nil
}

// TODO(schoon) - Handle badly-formatted lines.
func ResponseFromLine(line string) *http.Response {
	pieces := strings.Split(line, " ")
	code, _ := strconv.Atoi(pieces[1])

	return http.NewResponse(code, strings.Join(pieces[2:], " "))
}
