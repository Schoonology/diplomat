package builders

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/testdouble/diplomat/errors"
	"github.com/testdouble/diplomat/http"
	"github.com/testdouble/diplomat/parsers"
)

// The State builder converts a set of lines into a test.
type State struct{}

const (
	modeAwaitingRequest = iota
	modeInRequestHeaders
	modeInRequestBody
	modeInResponseHeaders
	modeInResponseBody
)

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

	if len(pieces) < 2 {
		return nil, &errors.BadResponseStatus{Line: line}
	}

	code, err := strconv.Atoi(pieces[0])
	remainder := pieces[1:]

	if err != nil {
		// TODO(schoon) - Validate this response HTTP version rather than
		// dropping it.
		code, _ = strconv.Atoi(pieces[1])
		remainder = pieces[2:]
	}

	return http.NewResponse(code, strings.Join(remainder, " ")), nil
}

func fallbackTestName(test Test) string {
	return fmt.Sprintf(
		"%s %s -> %d",
		test.Request.Method,
		test.Request.Path,
		test.Response.StatusCode)
}

// Build parses a set of lines into a test.
func (s *State) Build(paragraph parsers.Paragraph) (Test, error) {
	mode := modeAwaitingRequest
	test := Test{
		Name: paragraph.Name,
	}

	for _, line := range paragraph.Body {
		char := byte(0)
		trimmedLine := ""

		if len(line) > 0 {
			char = line[0]
			trimmedLine = strings.TrimSpace(line[1:])
		}

		if isRequestMetadataIndicator(char) {
			if isResponseMode(mode) {
				return test, errors.NewBuildError(paragraph, &errors.MisplacedRequest{})
			}
			if mode == modeAwaitingRequest {
				mode = modeInRequestHeaders
				request, err := requestFromLine(trimmedLine)
				if err != nil {
					return test, errors.NewBuildError(paragraph, err)
				}

				test.Request = request
			} else if mode == modeInRequestHeaders && len(trimmedLine) == 0 {
				mode = modeInRequestBody
			} else if mode == modeInRequestHeaders {
				key, value, err := headerFromLine(trimmedLine)
				if err != nil {
					return test, errors.NewBuildError(paragraph, err)
				}

				test.Request.Headers[key] = value
			}
		} else if isResponseMetadataIndicator(char) {
			if isRequestMode(mode) {
				mode = modeInResponseHeaders
				response, err := responseFromLine(trimmedLine)
				if err != nil {
					return test, errors.NewBuildError(paragraph, err)
				}

				test.Response = response
			} else if mode == modeInResponseHeaders && len(trimmedLine) == 0 {
				mode = modeInResponseBody
			} else if mode == modeInResponseHeaders {
				key, value, err := headerFromLine(trimmedLine)
				if err != nil {
					return test, errors.NewBuildError(paragraph, err)
				}

				test.Response.Headers[key] = value
			} else if mode == modeAwaitingRequest || mode == modeInResponseBody {
				// This can happen if an expected response body contains a response
				// indicator (e.g. `<`) as the first character. This explicitly
				// disallows that, preventing unexpected issues with missing requests.
				return test, errors.NewBuildError(paragraph, &errors.MissingRequest{})
			}
		} else if isRequestMode(mode) {
			mode = modeInRequestBody
			test.Request.Body = append(test.Request.Body, []byte(line+"\n")...)
		} else if isResponseMode(mode) {
			mode = modeInResponseBody
			test.Response.Body = append(test.Response.Body, []byte(line+"\n")...)
		}
	}

	if test.Request == nil {
		return test, errors.NewBuildError(paragraph, &errors.MissingRequest{})
	}
	if test.Response == nil {
		return test, errors.NewBuildError(paragraph, &errors.MissingResponse{})
	}
	if len(test.Name) == 0 {
		test.Name = fallbackTestName(test)
	}

	return test, nil
}

// BuildAll returns the results of running Build on all string arrays in a channel.
func (s *State) BuildAll(paragraphs chan parsers.Paragraph) chan Test {
	tests := make(chan Test)

	go func() {
		defer close(tests)

		for paragraph := range paragraphs {
			test, err := s.Build(paragraph)
			if err != nil {
				test.Err = err
				tests <- test
			} else {
				tests <- test
			}
		}
	}()

	return tests
}
