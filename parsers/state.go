package parsers

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/testdouble/http-assertion-tool/http"
)

const (
	ModeAwaitingRequest = iota
	ModeInRequestHeaders
	ModeInRequestBody
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

func isRequestMetadataIndicator(char byte) bool {
	return char == '>'
}

func isResponseMetadataIndicator(char byte) bool {
	return char == '<'
}

func isRequestMode(mode int) bool {
	return mode == ModeInRequestHeaders || mode == ModeInRequestBody
}

func isResponseMode(mode int) bool {
	return mode == ModeInResponseHeaders || mode == ModeInResponseBody
}

func (s *parserState) addLine(line string) error {
	// fmt.Printf("Adding line: %s\n", line)

	char := byte(0)
	trimmedLine := ""

	if len(line) > 0 {
		char = line[0]
		trimmedLine = strings.TrimSpace(line[1:])
	}

	if isRequestMetadataIndicator(char) {
		if isResponseMode(s.mode) {
			s.pushCurrentTest()
		}

		if isResponseMode(s.mode) || s.mode == ModeAwaitingRequest {
			s.mode = ModeInRequestHeaders
			s.currentTest.Request = RequestFromLine(trimmedLine)
			return nil
		} else if s.mode == ModeInRequestHeaders && len(trimmedLine) == 0 {
			s.mode = ModeInRequestBody
			return nil
		} else if s.mode == ModeInRequestHeaders {
			key, value, err := HeaderFromLine(trimmedLine)
			if err != nil {
				return err
			}

			s.currentTest.Request.Headers[key] = value
			return nil
		}
	}

	if isResponseMetadataIndicator(char) {
		if isRequestMode(s.mode) {
			s.mode = ModeInResponseHeaders
			s.currentTest.Response = ResponseFromLine(trimmedLine)
			return nil
		} else if s.mode == ModeInResponseHeaders && len(trimmedLine) == 0 {
			s.mode = ModeInResponseBody
			return nil
		} else if s.mode == ModeInResponseHeaders {
			key, value, err := HeaderFromLine(trimmedLine)
			if err != nil {
				return err
			}

			s.currentTest.Response.Headers[key] = value
			return nil
		}
	}

	if isRequestMode(s.mode) {
		s.mode = ModeInRequestBody
		s.currentTest.Request.Body = append(s.currentTest.Request.Body, []byte(line+"\n")...)
		return nil
	}

	if isResponseMode(s.mode) {
		s.mode = ModeInResponseBody
		s.currentTest.Response.Body = append(s.currentTest.Response.Body, []byte(line+"\n")...)
		return nil
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
		return "", "", fmt.Errorf("badly formatted header: %v", line)
	}

	return strings.TrimSpace(pieces[0]), strings.TrimSpace(pieces[1]), nil
}

// TODO(schoon) - Handle badly-formatted lines.
func ResponseFromLine(line string) *http.Response {
	pieces := strings.Split(line, " ")
	code, _ := strconv.Atoi(pieces[1])

	return http.NewResponse(code, strings.Join(pieces[2:], " "))
}
