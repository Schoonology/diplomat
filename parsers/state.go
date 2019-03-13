package parsers

import (
	"strconv"
	"strings"

	"github.com/testdouble/diplomat/errors"
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
	finalizer   testFinalizer
	tests       chan Test
}

func newParserState() parserState {
	return parserState{
		mode:  modeAwaitingRequest,
		tests: make(chan Test),
	}
}

func (s *parserState) pushCurrentTest() error {
	if s.currentTest.Request == nil {
		return &errors.MissingRequest{}
	}
	if s.currentTest.Response == nil {
		return &errors.MissingResponse{}
	}

	if s.finalizer != nil {
		s.finalizer(&s.currentTest)
	}

	s.tests <- s.currentTest
	s.currentTest = Test{}
	return nil
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
			err := s.pushCurrentTest()
			if err != nil {
				return err
			}
		}

		if isResponseMode(s.mode) || s.mode == modeAwaitingRequest {
			s.mode = modeInRequestHeaders
			request, err := requestFromLine(trimmedLine)
			if err != nil {
				return err
			}

			s.currentTest.Request = request
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
			response, err := responseFromLine(trimmedLine)
			if err != nil {
				return err
			}

			s.currentTest.Response = response
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
		} else if s.mode == modeAwaitingRequest || s.mode == modeInResponseBody {
			// This can happen if an expected response body contains a response
			// indicator (e.g. `<`) as the first character. This explicitly
			// disallows that, preventing unexpected issues with missing requests.
			return &errors.MissingRequest{}
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

func (s *parserState) finalize() (err error) {
	if s.mode != modeAwaitingRequest {
		err = s.pushCurrentTest()
	}

	close(s.tests)

	return err
}

func requestFromLine(line string) (*http.Request, error) {
	pieces := strings.Split(line, " ")

	if len(pieces) < 2 {
		return nil, &errors.BadRequestLine{Line: line}
	}

	return http.NewRequest(pieces[0], pieces[1]), nil
}

func headerFromLine(line string) (string, string, error) {
	pieces := strings.Split(line, ":")

	if len(pieces) < 2 {
		return "", "", &errors.BadHeader{Header: line}
	}

	return strings.TrimSpace(pieces[0]), strings.TrimSpace(pieces[1]), nil
}

func responseFromLine(line string) (*http.Response, error) {
	pieces := strings.Split(line, " ")

	if len(pieces) < 3 {
		return nil, &errors.BadResponseStatus{Line: line}
	}

	code, _ := strconv.Atoi(pieces[1])

	return http.NewResponse(code, strings.Join(pieces[2:], " ")), nil
}
