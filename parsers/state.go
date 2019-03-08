package parsers

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/testdouble/diplomat/http"
)

const (
	modeAwaitingRequest = iota
	modeInRequestHeaders
	modeInRequestBody
	modeInResponseHeaders
	modeInResponseBody
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
		mode: modeAwaitingRequest,
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
	return mode == modeInRequestHeaders || mode == modeInRequestBody
}

func isResponseMode(mode int) bool {
	return mode == modeInResponseHeaders || mode == modeInResponseBody
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

		if isResponseMode(s.mode) || s.mode == modeAwaitingRequest {
			s.mode = modeInRequestHeaders
			s.currentTest.Request = requestFromLine(trimmedLine)
			return nil
		} else if s.mode == modeInRequestHeaders && len(trimmedLine) == 0 {
			s.mode = modeInRequestBody
			return nil
		} else if s.mode == modeInRequestHeaders {
			key, value, err := headerFromLine(trimmedLine)
			if err != nil {
				return err
			}

			s.currentTest.Request.Headers[key] = value
			return nil
		}
	}

	if isResponseMetadataIndicator(char) {
		if isRequestMode(s.mode) {
			s.mode = modeInResponseHeaders
			s.currentTest.Response = responseFromLine(trimmedLine)
			return nil
		} else if s.mode == modeInResponseHeaders && len(trimmedLine) == 0 {
			s.mode = modeInResponseBody
			return nil
		} else if s.mode == modeInResponseHeaders {
			key, value, err := headerFromLine(trimmedLine)
			if err != nil {
				return err
			}

			s.currentTest.Response.Headers[key] = value
			return nil
		}
	}

	if isRequestMode(s.mode) {
		s.mode = modeInRequestBody
		s.currentTest.Request.Body = append(s.currentTest.Request.Body, []byte(line+"\n")...)
		return nil
	}

	if isResponseMode(s.mode) {
		s.mode = modeInResponseBody
		s.currentTest.Response.Body = append(s.currentTest.Response.Body, []byte(line+"\n")...)
		return nil
	}

	return nil
}

func (s *parserState) finalize() (*Spec, error) {
	if s.mode != modeAwaitingRequest {
		s.pushCurrentTest()
	}

	return s.spec, nil
}

// TODO(schoon) - Handle badly-formatted lines.
func requestFromLine(line string) *http.Request {
	pieces := strings.Split(line, " ")

	return http.NewRequest(pieces[0], pieces[1])
}

// TODO(schoon) - Handle badly-formatted lines.
func headerFromLine(line string) (string, string, error) {
	pieces := strings.Split(line, ":")

	if len(pieces) < 2 {
		return "", "", fmt.Errorf("badly formatted header: %v", line)
	}

	return strings.TrimSpace(pieces[0]), strings.TrimSpace(pieces[1]), nil
}

// TODO(schoon) - Handle badly-formatted lines.
func responseFromLine(line string) *http.Response {
	pieces := strings.Split(line, " ")
	code, _ := strconv.Atoi(pieces[1])

	return http.NewResponse(code, strings.Join(pieces[2:], " "))
}
